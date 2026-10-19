// Package check does the semantic analysis of a parsed program: every
// variable must be declared once, and every expression gets a type.
package check

import (
	"fmt"

	"github.com/kauanallyson/portes/internal/ast"
	"github.com/kauanallyson/portes/internal/token"
)

// Info is what the code generator needs to know about a checked program.
type Info struct {
	Vars  map[string]ast.Type   // declared variables
	Types map[ast.Expr]ast.Type // type of every well-typed expression
}

type checker struct {
	info *Info
	errs []error
}

// Check validates the program and returns its type information.
// Errors come in source order; if there are any, Info is incomplete.
func Check(prog *ast.Program) (*Info, []error) {
	c := &checker{info: &Info{
		Vars:  map[string]ast.Type{},
		Types: map[ast.Expr]ast.Type{},
	}}
	for _, decl := range prog.Vars {
		for _, id := range decl.Names {
			if _, dup := c.info.Vars[id.Name]; dup {
				c.errorf(id.Pos(), "variável '%s' já declarada", id.Name)
				continue
			}
			c.info.Vars[id.Name] = decl.Type
		}
	}
	c.stmts(prog.Body)
	return c.info, c.errs
}

func (c *checker) errorf(pos token.Pos, format string, args ...any) {
	c.errs = append(c.errs, &token.Error{Pos: pos, Msg: fmt.Sprintf(format, args...)})
}

// lookup returns the type of a variable, reporting it if undeclared.
func (c *checker) lookup(id *ast.Ident) (ast.Type, bool) {
	t, ok := c.info.Vars[id.Name]
	if !ok {
		c.errorf(id.Pos(), "variável '%s' não declarada", id.Name)
	}
	return t, ok
}

/* --- STATEMENTS --- */

func (c *checker) stmts(stmts []ast.Stmt) {
	for _, s := range stmts {
		c.stmt(s)
	}
}

func (c *checker) stmt(s ast.Stmt) {
	switch s := s.(type) {
	case *ast.AssignStmt:
		target, targetOK := c.lookup(s.Target)
		value, valueOK := c.expr(s.Value)
		if targetOK && valueOK && !assignable(target, value) {
			c.errorf(s.Pos(), "não é possível atribuir %s a '%s' (%s)", value, s.Target.Name, target)
		}
	case *ast.ReadStmt:
		for _, id := range s.Targets {
			if t, ok := c.lookup(id); ok && t == ast.Logico {
				c.errorf(id.Pos(), "LEIA não aceita a variável LOGICO '%s'", id.Name)
			}
		}
	case *ast.WriteStmt:
		for _, a := range s.Args {
			c.expr(a)
		}
	case *ast.IfStmt:
		c.want("condição do SE", s.Cond, ast.Logico)
		c.stmts(s.Then)
		c.stmts(s.Else)
	case *ast.WhileStmt:
		c.want("condição do ENQUANTO", s.Cond, ast.Logico)
		c.stmts(s.Body)
	case *ast.ForStmt:
		c.forStmt(s)
	default:
		panic(fmt.Sprintf("check: comando desconhecido %T", s))
	}
}

func (c *checker) forStmt(s *ast.ForStmt) {
	if t, ok := c.lookup(s.Var); ok && t != ast.Inteiro {
		c.errorf(s.Var.Pos(), "variável do PARA deve ser INTEIRO, '%s' é %s", s.Var.Name, t)
	}
	c.want("início do PARA", s.From, ast.Inteiro)
	c.want("fim do PARA", s.To, ast.Inteiro)
	if s.Step != nil {
		// The step must be known at compile time: its sign decides
		// whether the generated loop counts up (<=) or down (>=).
		if step, ok := ast.ConstInt(s.Step); !ok || step == 0 {
			c.errorf(s.Step.Pos(), "PASSO deve ser um número inteiro constante diferente de zero")
		}
	}
	c.stmts(s.Body)
}

// want checks that e has type t. what names e in the error message.
func (c *checker) want(what string, e ast.Expr, t ast.Type) {
	if got, ok := c.expr(e); ok && got != t {
		c.errorf(e.Pos(), "%s deve ser %s, encontrado %s", what, t, got)
	}
}

// assignable reports whether a value of type v can be stored in a variable
// of type t. The only implicit conversion is INTEIRO -> REAL.
func assignable(t, v ast.Type) bool {
	return t == v || t == ast.Real && v == ast.Inteiro
}

/* --- EXPRESSIONS --- */

// expr returns the type of e. ok is false when e has an error, which was
// already reported; callers then skip their own checks to avoid cascades.
func (c *checker) expr(e ast.Expr) (t ast.Type, ok bool) {
	defer func() {
		if ok {
			c.info.Types[e] = t
		}
	}()
	switch e := e.(type) {
	case *ast.Ident:
		return c.lookup(e)
	case *ast.IntLit:
		return ast.Inteiro, true
	case *ast.RealLit:
		return ast.Real, true
	case *ast.StringLit:
		return ast.String, true
	case *ast.BoolLit:
		return ast.Logico, true
	case *ast.ParenExpr:
		return c.expr(e.X)
	case *ast.UnaryExpr:
		x, ok := c.expr(e.X)
		if !ok {
			return 0, false
		}
		if t, ok := unaryType(e.Op, x); ok {
			return t, true
		}
		c.errorf(e.Pos(), "operador '%s' não se aplica a %s", e.Op, x)
	case *ast.BinaryExpr:
		l, leftOK := c.expr(e.Left)
		r, rightOK := c.expr(e.Right)
		if !leftOK || !rightOK {
			return 0, false
		}
		if t, ok := binaryType(e.Op, l, r); ok {
			return t, true
		}
		c.errorf(e.Pos(), "operador '%s' não se aplica a %s e %s", e.Op, l, r)
	}
	return 0, false
}

func isNumber(t ast.Type) bool { return t == ast.Inteiro || t == ast.Real }

func unaryType(op token.Kind, x ast.Type) (ast.Type, bool) {
	switch op {
	case token.MINUS:
		return x, isNumber(x)
	case token.NAO:
		return ast.Logico, x == ast.Logico
	}
	return 0, false
}

func binaryType(op token.Kind, l, r ast.Type) (ast.Type, bool) {
	switch op {
	case token.PLUS, token.MINUS, token.MULTIPLICATION, token.DIVISION:
		// As in C, INTEIRO / INTEIRO is an integer division.
		if !isNumber(l) || !isNumber(r) {
			return 0, false
		}
		if l == ast.Real || r == ast.Real {
			return ast.Real, true
		}
		return ast.Inteiro, true
	case token.MOD:
		return ast.Inteiro, l == ast.Inteiro && r == ast.Inteiro
	case token.E, token.OU:
		return ast.Logico, l == ast.Logico && r == ast.Logico
	case token.EQUAL, token.DIFFERENT:
		return ast.Logico, l == r || isNumber(l) && isNumber(r)
	case token.L_THAN, token.G_THAN, token.LE_THAN, token.GE_THAN:
		return ast.Logico, isNumber(l) && isNumber(r) || l == ast.String && r == ast.String
	}
	return 0, false
}
