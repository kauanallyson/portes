package ast

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Fprint writes an indented view of the program, with expressions shown
// in prefix form so precedence is explicit: soma <- (+ (+ n1 n2) n3)
func Fprint(w io.Writer, p *Program) {
	fmt.Fprintf(w, "ALGORITMO %s\n", p.Name.Name)
	if len(p.Vars) > 0 {
		fmt.Fprintln(w, "  VARIAVEIS")
		for _, v := range p.Vars {
			names := make([]string, len(v.Names))
			for i, id := range v.Names {
				names[i] = id.Name
			}
			fmt.Fprintf(w, "    %s : %s\n", strings.Join(names, ", "), v.Type)
		}
	}
	fmt.Fprintln(w, "  INICIO")
	printStmts(w, p.Body, 2)
}

func printStmts(w io.Writer, stmts []Stmt, depth int) {
	for _, s := range stmts {
		printStmt(w, s, depth)
	}
}

func printStmt(w io.Writer, s Stmt, depth int) {
	indent := strings.Repeat("  ", depth)
	switch s := s.(type) {
	case *AssignStmt:
		fmt.Fprintf(w, "%s%s <- %s\n", indent, s.Target.Name, ExprString(s.Value))
	case *ReadStmt:
		names := make([]string, len(s.Targets))
		for i, id := range s.Targets {
			names[i] = id.Name
		}
		fmt.Fprintf(w, "%sLEIA %s\n", indent, strings.Join(names, ", "))
	case *WriteStmt:
		args := make([]string, len(s.Args))
		for i, a := range s.Args {
			args[i] = ExprString(a)
		}
		fmt.Fprintf(w, "%s%s\n", indent, strings.TrimSpace("ESCREVA "+strings.Join(args, ", ")))
	case *IfStmt:
		fmt.Fprintf(w, "%sSE %s\n", indent, ExprString(s.Cond))
		fmt.Fprintf(w, "%s  ENTAO\n", indent)
		printStmts(w, s.Then, depth+2)
		if s.Else != nil {
			fmt.Fprintf(w, "%s  SENAO\n", indent)
			printStmts(w, s.Else, depth+2)
		}
	case *WhileStmt:
		fmt.Fprintf(w, "%sENQUANTO %s\n", indent, ExprString(s.Cond))
		fmt.Fprintf(w, "%s  FACA\n", indent)
		printStmts(w, s.Body, depth+2)
	case *ForStmt:
		fmt.Fprintf(w, "%sPARA %s DE %s ATE %s", indent, s.Var.Name, ExprString(s.From), ExprString(s.To))
		if s.Step != nil {
			fmt.Fprintf(w, " PASSO %s", ExprString(s.Step))
		}
		fmt.Fprintf(w, "\n%s  FACA\n", indent)
		printStmts(w, s.Body, depth+2)
	default:
		fmt.Fprintf(w, "%s<stmt desconhecido %T>\n", indent, s)
	}
}

// ExprString renders an expression in prefix form.
func ExprString(e Expr) string {
	switch e := e.(type) {
	case *Ident:
		return e.Name
	case *IntLit:
		return e.Value
	case *RealLit:
		return e.Value
	case *StringLit:
		return strconv.Quote(e.Value)
	case *BoolLit:
		if e.Value {
			return "VERDADEIRO"
		}
		return "FALSO"
	case *UnaryExpr:
		return fmt.Sprintf("(%s %s)", e.Op, ExprString(e.X))
	case *BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", e.Op, ExprString(e.Left), ExprString(e.Right))
	case *ParenExpr:
		return ExprString(e.X)
	}
	return fmt.Sprintf("<expr desconhecida %T>", e)
}
