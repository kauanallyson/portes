package lexer

import (
	"strings"
	"unicode"

	"github.com/kauanallyson/portes/internal/token"
)

type Lexer struct {
	src    []rune
	pos    int
	line   int
	column int
}

func New(src string) *Lexer {
	return &Lexer{src: []rune(src), line: 1}
}

// Tokenize lexes the whole source. The last token is always EOF.
func Tokenize(src string) []token.Token {
	l := New(src)
	var tks []token.Token
	for {
		tk := l.Next()
		tks = append(tks, tk)
		if tk.Kind == token.EOF {
			return tks
		}
	}
}

func (l *Lexer) atEnd() bool { return l.pos >= len(l.src) }

func (l *Lexer) peek() rune {
	if l.atEnd() {
		return 0
	}
	return l.src[l.pos]
}

func (l *Lexer) peekNext() rune {
	if l.pos+1 >= len(l.src) {
		return 0
	}
	return l.src[l.pos+1]
}

func (l *Lexer) advance() rune {
	c := l.src[l.pos]
	l.pos++
	if c == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
	return c
}

func (l *Lexer) skipIgnored() {
	for !l.atEnd() {
		switch {
		case unicode.IsSpace(l.peek()):
			l.advance()
		case l.peek() == '/' && l.peekNext() == '/':
			for !l.atEnd() && l.peek() != '\n' {
				l.advance()
			}
		default:
			return
		}
	}
}

func (l *Lexer) Next() token.Token {
	l.skipIgnored()

	start := token.Pos{Line: l.line, Column: l.column + 1}
	tok := func(kind token.Kind, value string) token.Token {
		return token.Token{Kind: kind, Value: value, Pos: start}
	}

	if l.atEnd() {
		return tok(token.EOF, "")
	}

	c := l.advance()
	switch c {
	case '+':
		return tok(token.PLUS, "+")
	case '-':
		return tok(token.MINUS, "-")
	case '*':
		return tok(token.MULTIPLICATION, "*")
	case '/':
		return tok(token.DIVISION, "/")
	case '=':
		return tok(token.EQUAL, "=")
	case ';':
		return tok(token.SEMICOLON, ";")
	case ':':
		return tok(token.COLON, ":")
	case ',':
		return tok(token.COMMA, ",")
	case '(':
		return tok(token.L_PAREN, "(")
	case ')':
		return tok(token.R_PAREN, ")")
	case '<':
		switch l.peek() {
		case '-':
			l.advance()
			return tok(token.ASSIGN, "<-")
		case '>':
			l.advance()
			return tok(token.DIFFERENT, "<>")
		case '=':
			l.advance()
			return tok(token.LE_THAN, "<=")
		}
		return tok(token.L_THAN, "<")
	case '>':
		if l.peek() == '=' {
			l.advance()
			return tok(token.GE_THAN, ">=")
		}
		return tok(token.G_THAN, ">")
	case '"':
		return l.readString(start)
	}

	switch {
	case isLetter(c):
		return l.readWord(c, start)
	case isDigit(c):
		return l.readNumber(c, start)
	}
	return tok(token.UNKNOWN, string(c))
}

// readString reads until the closing quote. The quotes are not part of the value.
func (l *Lexer) readString(start token.Pos) token.Token {
	var sb strings.Builder
	for !l.atEnd() && l.peek() != '"' {
		sb.WriteRune(l.advance())
	}
	if l.atEnd() {
		return token.Token{Kind: token.UNKNOWN, Value: `"` + sb.String(), Pos: start}
	}
	l.advance() // closing quote
	return token.Token{Kind: token.STRING_LIT, Value: sb.String(), Pos: start}
}

func (l *Lexer) readWord(first rune, start token.Pos) token.Token {
	var sb strings.Builder
	sb.WriteRune(first)
	for isLetter(l.peek()) || isDigit(l.peek()) {
		sb.WriteRune(l.advance())
	}
	word := sb.String()
	return token.Token{Kind: token.Lookup(word), Value: word, Pos: start}
}

// readNumber reads an integer (42) or a real (3.14).
func (l *Lexer) readNumber(first rune, start token.Pos) token.Token {
	var sb strings.Builder
	sb.WriteRune(first)
	kind := token.INT_LIT
	for isDigit(l.peek()) {
		sb.WriteRune(l.advance())
	}
	if l.peek() == '.' && isDigit(l.peekNext()) {
		kind = token.REAL_LIT
		sb.WriteRune(l.advance())
		for isDigit(l.peek()) {
			sb.WriteRune(l.advance())
		}
	}
	return token.Token{Kind: kind, Value: sb.String(), Pos: start}
}

func isLetter(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func isDigit(c rune) bool { return c >= '0' && c <= '9' }
