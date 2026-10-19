package token

import "fmt"

type Kind int

const (
	// Specials
	UNKNOWN Kind = iota
	EOF

	// Literals
	IDENTIFIER
	INT_LIT    // 67
	REAL_LIT   // 3.14
	STRING_LIT // "text"

	// Operators and Symbols
	PLUS           // +
	MINUS          // -
	MULTIPLICATION // *
	DIVISION       // /
	ASSIGN         // <-
	EQUAL          // =
	DIFFERENT      // <>
	L_THAN         // <
	G_THAN         // >
	LE_THAN        // <=
	GE_THAN        // >=
	SEMICOLON      // ;
	COLON          // :
	COMMA          // ,
	L_PAREN        // (
	R_PAREN        // )

	keywordBegin
	ALGORITMO
	VARIAVEIS
	INICIO
	FIM
	SE
	ENTAO
	SENAO
	FIMSE
	INTEIRO
	REAL
	STRING
	LOGICO
	VERDADEIRO
	FALSO
	LEIA
	ESCREVA
	ENQUANTO
	FACA
	FIMENQUANTO
	PARA
	DE
	ATE
	PASSO
	FIMPARA
	E
	OU
	NAO
	MOD
	keywordEnd
)

var names = [...]string{
	UNKNOWN: "UNKNOWN",
	EOF:     "EOF",

	IDENTIFIER: "IDENTIFIER",
	INT_LIT:    "INT_LIT",
	REAL_LIT:   "REAL_LIT",
	STRING_LIT: "STRING_LIT",

	PLUS:           "+",
	MINUS:          "-",
	MULTIPLICATION: "*",
	DIVISION:       "/",
	ASSIGN:         "<-",
	EQUAL:          "=",
	DIFFERENT:      "<>",
	L_THAN:         "<",
	G_THAN:         ">",
	LE_THAN:        "<=",
	GE_THAN:        ">=",
	SEMICOLON:      ";",
	COLON:          ":",
	COMMA:          ",",
	L_PAREN:        "(",
	R_PAREN:        ")",

	ALGORITMO:  "ALGORITMO",
	VARIAVEIS:  "VARIAVEIS",
	INICIO:     "INICIO",
	FIM:        "FIM",
	SE:         "SE",
	ENTAO:      "ENTAO",
	SENAO:      "SENAO",
	FIMSE:      "FIMSE",
	INTEIRO:    "INTEIRO",
	REAL:       "REAL",
	STRING:     "STRING",
	LOGICO:     "LOGICO",
	VERDADEIRO: "VERDADEIRO",
	FALSO:      "FALSO",
	LEIA:       "LEIA",
	ESCREVA:    "ESCREVA",

	ENQUANTO:    "ENQUANTO",
	FACA:        "FACA",
	FIMENQUANTO: "FIMENQUANTO",
	PARA:        "PARA",
	DE:          "DE",
	ATE:         "ATE",
	PASSO:       "PASSO",
	FIMPARA:     "FIMPARA",
	E:           "E",
	OU:          "OU",
	NAO:         "NAO",
	MOD:         "MOD",
}

func (k Kind) String() string {
	if k >= 0 && int(k) < len(names) && names[k] != "" {
		return names[k]
	}
	return fmt.Sprintf("Kind(%d)", int(k))
}

var keywords = func() map[string]Kind {
	m := make(map[string]Kind, keywordEnd-keywordBegin)
	for k := keywordBegin + 1; k < keywordEnd; k++ {
		m[names[k]] = k
	}
	return m
}()

func Lookup(word string) Kind {
	if k, ok := keywords[word]; ok {
		return k
	}
	return IDENTIFIER
}

type Pos struct {
	Line, Column int
}

func (p Pos) String() string { return fmt.Sprintf("%d:%d", p.Line, p.Column) }

// IsBlockEnd reports whether k closes a block of statements.
func IsBlockEnd(k Kind) bool {
	switch k {
	case FIM, FIMSE, SENAO, FIMENQUANTO, FIMPARA:
		return true
	}
	return false
}

// Error is a syntax or semantic error at a given position of the source.
type Error struct {
	Pos Pos
	Msg string
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %s", e.Pos, e.Msg) }

type Token struct {
	Kind  Kind
	Value string
	Pos   Pos
}
