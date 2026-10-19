// Package codegen translates a checked program into C.
package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kauanallyson/portes/internal/ast"
	"github.com/kauanallyson/portes/internal/check"
	"github.com/kauanallyson/portes/internal/token"
)

// StringSize is the size of the char array behind a STRING variable.
const StringSize = 256

// stringLib declares strcpy and strcmp, used for STRING.
const stringLib = "string.h"

// cTypes says how each Portes type is declared, read and written in C.
var cTypes = map[ast.Type]struct {
	name   string // C type
	header string // header that defines name, if any
	zero   string // initial value
	scanf  string // conversion for LEIA ("" when not accepted)
	printf string // conversion for ESCREVA
}{
	ast.Inteiro: {"int", "", "0", "%d", "%d"},
	// %.15g keeps 7.5 as 7.5 and avoids 1.23457e+06 for 1234567.5.
	ast.Real:   {"double", "", "0", "%lf", "%.15g"},
	ast.String: {"char", "", `""`, fmt.Sprintf("%%%ds", StringSize-1), "%s"},
	ast.Logico: {"bool", "stdbool.h", "false", "", "%s"},
}

type generator struct {
	info     *check.Info
	sb       strings.Builder
	depth    int
	includes map[string]bool
	names    *names
}

// Generate returns the C source of prog. info must come from a Check
// of the same program that returned no errors.
func Generate(prog *ast.Program, info *check.Info) string {
	g := &generator{info: info, includes: map[string]bool{"stdio.h": true}, depth: 1, names: newNames(info.Vars)}
	for _, v := range prog.Vars {
		g.varDecl(v)
	}
	if len(prog.Vars) > 0 && len(prog.Body) > 0 {
		g.sb.WriteString("\n")
	}
	g.stmts(prog.Body)
	g.line("return 0;")
	body := g.sb.String()

	var out strings.Builder
	headers := make([]string, 0, len(g.includes))
	for h := range g.includes {
		headers = append(headers, h)
	}
	sort.Strings(headers)
	for _, h := range headers {
		fmt.Fprintf(&out, "#include <%s>\n", h)
	}
	fmt.Fprintf(&out, "\n/* ALGORITMO %s */\nint main(void) {\n%s}\n", prog.Name.Name, body)
	return out.String()
}

func (g *generator) line(format string, args ...any) {
	g.sb.WriteString(strings.Repeat("    ", g.depth))
	fmt.Fprintf(&g.sb, format, args...)
	g.sb.WriteString("\n")
}

func (g *generator) block(stmts []ast.Stmt) {
	g.depth++
	g.stmts(stmts)
	g.depth--
}

/* --- DECLARATIONS --- */

func (g *generator) varDecl(v *ast.VarDecl) {
	ct := cTypes[v.Type]
	if ct.header != "" {
		g.includes[ct.header] = true
	}
	names := make([]string, len(v.Names))
	for i, id := range v.Names {
		name := g.names.of(id.Name)
		if v.Type == ast.String {
			name = fmt.Sprintf("%s[%d]", name, StringSize)
		}
		names[i] = name + " = " + ct.zero
	}
	g.line("%s %s;", ct.name, strings.Join(names, ", "))
}

/* --- STATEMENTS --- */

func (g *generator) stmts(stmts []ast.Stmt) {
	for _, s := range stmts {
		g.stmt(s)
	}
}

func (g *generator) stmt(s ast.Stmt) {
	switch s := s.(type) {
	case *ast.AssignStmt:
		if g.info.Vars[s.Target.Name] == ast.String {
			g.includes[stringLib] = true
			g.line("strcpy(%s, %s);", g.names.of(s.Target.Name), g.expr(s.Value))
			return
		}
		g.line("%s = %s;", g.names.of(s.Target.Name), g.expr(s.Value))
	case *ast.ReadStmt:
		for _, id := range s.Targets {
			t, name := g.info.Vars[id.Name], g.names.of(id.Name)
			if t != ast.String { // a char array is already a pointer
				name = "&" + name
			}
			g.line(`scanf("%s", %s);`, cTypes[t].scanf, name)
		}
	case *ast.WriteStmt:
		g.write(s)
	case *ast.IfStmt:
		g.line("if (%s) {", g.expr(s.Cond))
		g.block(s.Then)
		if s.Else != nil {
			g.line("} else {")
			g.block(s.Else)
		}
		g.line("}")
	case *ast.WhileStmt:
		g.line("while (%s) {", g.expr(s.Cond))
		g.block(s.Body)
		g.line("}")
	case *ast.ForStmt:
		g.forStmt(s)
	default:
		panic(fmt.Sprintf("codegen: comando desconhecido %T", s))
	}
}

// forStmt counts up to To when PASSO is positive and down to it otherwise.
// Like in Portes, a non-constant To is evaluated only once, before the loop.
func (g *generator) forStmt(s *ast.ForStmt) {
	i, step := g.names.of(s.Var.Name), 1
	if s.Step != nil {
		step, _ = ast.ConstInt(s.Step)
	}
	to := g.expr(s.To)
	if _, ok := ast.ConstInt(s.To); !ok {
		limit := g.names.fresh("ate_" + i)
		g.line("int %s = %s;", limit, to)
		to = limit
	}
	cmp := "<="
	if step < 0 {
		cmp = ">="
	}
	var inc string
	switch {
	case step == 1:
		inc = i + "++"
	case step == -1:
		inc = i + "--"
	case step > 0:
		inc = fmt.Sprintf("%s += %d", i, step)
	default:
		inc = fmt.Sprintf("%s -= %d", i, -step)
	}
	g.line("for (%s = %s; %s %s %s; %s) {", i, g.expr(s.From), i, cmp, to, inc)
	g.block(s.Body)
	g.line("}")
}

// write turns ESCREVA into a single printf. String literals go straight
// into the format; every other argument gets a conversion by its type.
func (g *generator) write(s *ast.WriteStmt) {
	var format strings.Builder
	var args []string
	for _, a := range s.Args {
		if lit, ok := a.(*ast.StringLit); ok {
			format.WriteString(strings.ReplaceAll(escape(lit.Value), "%", "%%"))
			continue
		}
		t := g.info.Types[a]
		format.WriteString(cTypes[t].printf)
		if t == ast.Logico {
			args = append(args, g.operand(a, precTernary)+` ? "VERDADEIRO" : "FALSO"`)
		} else {
			args = append(args, g.expr(a))
		}
	}
	format.WriteString(`\n`)
	g.line("printf(%s);", strings.Join(append([]string{`"` + format.String() + `"`}, args...), ", "))
}

/* --- EXPRESSIONS --- */

// C precedence levels, from loosest to tightest.
const (
	precTernary = iota + 1
	precOr
	precAnd
	precEquality
	precRelational
	precAdditive
	precMultiplicative
	precUnary
	precPrimary
)

var unaryOps = map[token.Kind]string{
	token.MINUS: "-",
	token.NAO:   "!",
}

var binaryOps = map[token.Kind]struct {
	c    string
	prec int
}{
	token.PLUS:           {"+", precAdditive},
	token.MINUS:          {"-", precAdditive},
	token.MULTIPLICATION: {"*", precMultiplicative},
	token.DIVISION:       {"/", precMultiplicative},
	token.MOD:            {"%", precMultiplicative},
	token.E:              {"&&", precAnd},
	token.OU:             {"||", precOr},
	token.EQUAL:          {"==", precEquality},
	token.DIFFERENT:      {"!=", precEquality},
	token.L_THAN:         {"<", precRelational},
	token.G_THAN:         {">", precRelational},
	token.LE_THAN:        {"<=", precRelational},
	token.GE_THAN:        {">=", precRelational},
}

func prec(e ast.Expr) int {
	switch e := e.(type) {
	case *ast.BinaryExpr:
		return binaryOps[e.Op].prec
	case *ast.UnaryExpr:
		return precUnary
	}
	return precPrimary
}

// operand renders e, in parentheses if it binds looser than minPrec.
func (g *generator) operand(e ast.Expr, minPrec int) string {
	if prec(e) < minPrec {
		return "(" + g.expr(e) + ")"
	}
	return g.expr(e)
}

func (g *generator) expr(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.Ident:
		return g.names.of(e.Name)
	case *ast.IntLit:
		return e.Value
	case *ast.RealLit:
		return e.Value
	case *ast.StringLit:
		return `"` + escape(e.Value) + `"`
	case *ast.BoolLit:
		g.includes[cTypes[ast.Logico].header] = true
		if e.Value {
			return "true"
		}
		return "false"
	case *ast.ParenExpr:
		return "(" + g.expr(e.X) + ")"
	case *ast.UnaryExpr:
		// Unary operands are always wrapped, so "- -x" never becomes "--x".
		return unaryOps[e.Op] + g.operand(e.X, precPrimary)
	case *ast.BinaryExpr:
		op := binaryOps[e.Op]
		if g.info.Types[e.Left] == ast.String {
			// Strings are compared with strcmp: a < b becomes strcmp(a, b) < 0.
			g.includes[stringLib] = true
			return fmt.Sprintf("strcmp(%s, %s) %s 0", g.expr(e.Left), g.expr(e.Right), op.c)
		}
		// Left-associative: the right side needs parentheses on a tie.
		left, right := g.operand(e.Left, op.prec), g.operand(e.Right, op.prec+1)
		if e.Op == token.OU {
			// Not needed by C, but makes a || (b && c) explicit (and keeps -Wall quiet).
			left, right = g.orOperand(e.Left, left), g.orOperand(e.Right, right)
		}
		return fmt.Sprintf("%s %s %s", left, op.c, right)
	}
	panic(fmt.Sprintf("codegen: expressão desconhecida %T", e))
}

func (g *generator) orOperand(e ast.Expr, rendered string) string {
	if b, ok := e.(*ast.BinaryExpr); ok && b.Op == token.E {
		return "(" + rendered + ")"
	}
	return rendered
}

// escape makes s safe inside a C string literal.
func escape(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\t", `\t`, "\r", `\r`)
	return r.Replace(s)
}
