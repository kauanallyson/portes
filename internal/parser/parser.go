package parser

import (
	"fmt"

	"github.com/kauanallyson/portes/internal/ast"
	"github.com/kauanallyson/portes/internal/token"
)

// Error is a syntax error at a given position of the source.
type Error struct {
	Pos token.Pos
	Msg string
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Pos, e.Msg) }

type Parser struct {
	tokens []token.Token
	pos    int
}

// Parse builds the program AST. It stops at the first syntax error.
func Parse(tokens []token.Token) (prog *ast.Program, err error) {
	p := &Parser{tokens: tokens}
	defer func() {
		if r := recover(); r != nil {
			perr, ok := r.(*Error)
			if !ok {
				panic(r)
			}
			prog, err = nil, perr
		}
	}()
	return p.parseProgram(), nil
}

/* --- HELPERS --- */

func (p *Parser) cur() token.Token { return p.tokens[p.pos] }

func (p *Parser) at(kinds ...token.Kind) bool {
	for _, k := range kinds {
		if p.cur().Kind == k {
			return true
		}
	}
	return false
}

func (p *Parser) next() token.Token {
	tk := p.cur()
	if tk.Kind != token.EOF {
		p.pos++
	}
	return tk
}

func (p *Parser) fail(pos token.Pos, format string, args ...any) {
	panic(&Error{Pos: pos, Msg: fmt.Sprintf(format, args...)})
}

func (p *Parser) expect(kind token.Kind) token.Token {
	if !p.at(kind) {
		p.fail(p.cur().Pos, "esperado %s, encontrado %s", describe(kind), found(p.cur()))
	}
	return p.next()
}

func describe(k token.Kind) string {
	switch k {
	case token.IDENTIFIER:
		return "um identificador"
	case token.EOF:
		return "fim do arquivo"
	}
	return fmt.Sprintf("'%s'", k)
}

func found(tk token.Token) string {
	switch tk.Kind {
	case token.EOF:
		return "fim do arquivo"
	case token.STRING_LIT:
		return fmt.Sprintf("%q", tk.Value)
	case token.UNKNOWN:
		return fmt.Sprintf("símbolo desconhecido '%s'", tk.Value)
	}
	return fmt.Sprintf("'%s'", tk.Value)
}

func (p *Parser) parseIdent() *ast.Ident {
	tk := p.expect(token.IDENTIFIER)
	return &ast.Ident{At: tk.Pos, Name: tk.Value}
}

func (p *Parser) parseIdentList() []*ast.Ident {
	ids := []*ast.Ident{p.parseIdent()}
	for p.at(token.COMMA) {
		p.next()
		ids = append(ids, p.parseIdent())
	}
	return ids
}

/* --- PROGRAM & DECLARATIONS --- */

func (p *Parser) parseProgram() *ast.Program {
	start := p.expect(token.ALGORITMO)
	prog := &ast.Program{At: start.Pos}
	prog.Name = p.parseIdent()
	p.expect(token.SEMICOLON)

	if p.at(token.VARIAVEIS) {
		p.next()
		for p.at(token.IDENTIFIER) {
			prog.Vars = append(prog.Vars, p.parseVarDecl())
		}
	}

	p.expect(token.INICIO)
	prog.Body = p.parseBlock(token.FIM)
	p.expect(token.FIM)
	p.expect(token.SEMICOLON)
	p.expect(token.EOF)
	return prog
}

func (p *Parser) parseVarDecl() *ast.VarDecl {
	decl := &ast.VarDecl{At: p.cur().Pos}
	decl.Names = p.parseIdentList()
	p.expect(token.COLON)
	decl.Type = p.parseType()
	p.expect(token.SEMICOLON)
	return decl
}

func (p *Parser) parseType() ast.Type {
	tk := p.next()
	switch tk.Kind {
	case token.INTEIRO:
		return ast.Inteiro
	case token.REAL:
		return ast.Real
	case token.STRING:
		return ast.String
	case token.LOGICO:
		return ast.Logico
	}
	p.fail(tk.Pos, "esperado um tipo (INTEIRO, REAL, STRING ou LOGICO), encontrado %s", found(tk))
	return 0
}

/* --- STATEMENTS --- */

func (p *Parser) parseBlock(terminators ...token.Kind) []ast.Stmt {
	var stmts []ast.Stmt
	for !p.at(terminators...) {
		if p.at(token.EOF) {
			p.fail(p.cur().Pos, "fim do arquivo inesperado, esperado %s", describe(terminators[0]))
		}
		if p.at(token.FIM, token.FIMSE, token.SENAO) {
			p.fail(p.cur().Pos, "esperado %s, encontrado %s", describe(terminators[0]), found(p.cur()))
		}
		stmts = append(stmts, p.parseStmt())
	}
	return stmts
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.cur().Kind {
	case token.IDENTIFIER:
		return p.parseAssign()
	case token.LEIA:
		return p.parseRead()
	case token.ESCREVA:
		return p.parseWrite()
	case token.SE:
		return p.parseIf()
	}
	p.fail(p.cur().Pos, "comando inválido: %s", found(p.cur()))
	return nil
}

func (p *Parser) parseAssign() *ast.AssignStmt {
	target := p.parseIdent()
	p.expect(token.ASSIGN)
	value := p.parseExpr()
	p.expect(token.SEMICOLON)
	return &ast.AssignStmt{At: target.Pos(), Target: target, Value: value}
}

func (p *Parser) parseRead() *ast.ReadStmt {
	start := p.expect(token.LEIA)
	p.expect(token.L_PAREN)
	targets := p.parseIdentList()
	p.expect(token.R_PAREN)
	p.expect(token.SEMICOLON)
	return &ast.ReadStmt{At: start.Pos, Targets: targets}
}

func (p *Parser) parseWrite() *ast.WriteStmt {
	start := p.expect(token.ESCREVA)
	stmt := &ast.WriteStmt{At: start.Pos}
	p.expect(token.L_PAREN)
	if !p.at(token.R_PAREN) {
		stmt.Args = append(stmt.Args, p.parseExpr())
		for p.at(token.COMMA) {
			p.next()
			stmt.Args = append(stmt.Args, p.parseExpr())
		}
	}
	p.expect(token.R_PAREN)
	p.expect(token.SEMICOLON)
	return stmt
}

func (p *Parser) parseIf() *ast.IfStmt {
	start := p.expect(token.SE)
	stmt := &ast.IfStmt{At: start.Pos}
	stmt.Cond = p.parseExpr()
	p.expect(token.ENTAO)
	stmt.Then = p.parseBlock(token.FIMSE, token.SENAO)
	if p.at(token.SENAO) {
		p.next()
		stmt.Else = p.parseBlock(token.FIMSE)
		if stmt.Else == nil {
			stmt.Else = []ast.Stmt{}
		}
	}
	p.expect(token.FIMSE)
	if p.at(token.SEMICOLON) {
		p.next()
	}
	return stmt
}

/* --- EXPRESSIONS --- */

func (p *Parser) parseExpr() ast.Expr {
	left := p.parseAdditive()
	if p.at(token.EQUAL, token.DIFFERENT, token.L_THAN, token.G_THAN, token.LE_THAN, token.GE_THAN) {
		op := p.next()
		right := p.parseAdditive()
		left = &ast.BinaryExpr{At: op.Pos, Op: op.Kind, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseAdditive() ast.Expr {
	left := p.parseTerm()
	for p.at(token.PLUS, token.MINUS) {
		op := p.next()
		right := p.parseTerm()
		left = &ast.BinaryExpr{At: op.Pos, Op: op.Kind, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseTerm() ast.Expr {
	left := p.parseUnary()
	for p.at(token.MULTIPLICATION, token.DIVISION) {
		op := p.next()
		right := p.parseUnary()
		left = &ast.BinaryExpr{At: op.Pos, Op: op.Kind, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseUnary() ast.Expr {
	if p.at(token.MINUS) {
		op := p.next()
		return &ast.UnaryExpr{At: op.Pos, Op: op.Kind, X: p.parseUnary()}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Expr {
	tk := p.next()
	switch tk.Kind {
	case token.IDENTIFIER:
		return &ast.Ident{At: tk.Pos, Name: tk.Value}
	case token.INT_LIT:
		return &ast.IntLit{At: tk.Pos, Value: tk.Value}
	case token.REAL_LIT:
		return &ast.RealLit{At: tk.Pos, Value: tk.Value}
	case token.STRING_LIT:
		return &ast.StringLit{At: tk.Pos, Value: tk.Value}
	case token.VERDADEIRO:
		return &ast.BoolLit{At: tk.Pos, Value: true}
	case token.FALSO:
		return &ast.BoolLit{At: tk.Pos, Value: false}
	case token.L_PAREN:
		x := p.parseExpr()
		p.expect(token.R_PAREN)
		return &ast.ParenExpr{At: tk.Pos, X: x}
	}
	p.fail(tk.Pos, "esperado uma expressão, encontrado %s", found(tk))
	return nil
}
