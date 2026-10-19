package ast

import "github.com/kauanallyson/portes/internal/token"

type Node interface {
	Pos() token.Pos
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Base struct {
	At token.Pos
}

func (b Base) Pos() token.Pos { return b.At }

// ALGORITMO <Name>; VARIAVEIS <Vars> INICIO <Body> FIM;
type Program struct {
	Base
	Name *Ident
	Vars []*VarDecl
	Body []Stmt
}

type Type int

const (
	Inteiro Type = iota
	Real
	String
	Logico
)

func (t Type) String() string {
	switch t {
	case Inteiro:
		return "INTEIRO"
	case Real:
		return "REAL"
	case String:
		return "STRING"
	case Logico:
		return "LOGICO"
	}
	return "TIPO?"
}

// n1, n2, soma : INTEIRO;
type VarDecl struct {
	Base
	Names []*Ident
	Type  Type
}

/* --- STATEMENTS --- */

// soma <- n1 + n2;
type AssignStmt struct {
	Base
	Target *Ident
	Value  Expr
}

// LEIA(n1, n2);
type ReadStmt struct {
	Base
	Targets []*Ident
}

// ESCREVA("soma: ", soma);
type WriteStmt struct {
	Base
	Args []Expr
}

// SE <Cond> ENTAO <Then> [SENAO <Else>] FIMSE
type IfStmt struct {
	Base
	Cond Expr
	Then []Stmt
	Else []Stmt
}

func (*AssignStmt) stmtNode() {}
func (*ReadStmt) stmtNode()   {}
func (*WriteStmt) stmtNode()  {}
func (*IfStmt) stmtNode()     {}

/* --- EXPRESSIONS --- */

type Ident struct {
	Base
	Name string
}

type IntLit struct {
	Base
	Value string
}

type RealLit struct {
	Base
	Value string
}

type StringLit struct {
	Base
	Value string
}

type BoolLit struct {
	Base
	Value bool
}

// -x
type UnaryExpr struct {
	Base
	Op token.Kind
	X  Expr
}

// a + b, a < b, ...
type BinaryExpr struct {
	Base
	Op    token.Kind
	Left  Expr
	Right Expr
}

type ParenExpr struct {
	Base
	X Expr
}

func (*Ident) exprNode()      {}
func (*IntLit) exprNode()     {}
func (*RealLit) exprNode()    {}
func (*StringLit) exprNode()  {}
func (*BoolLit) exprNode()    {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*ParenExpr) exprNode()  {}
