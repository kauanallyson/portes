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

type Token struct {
	Kind  Kind
	Value string
	Pos   Pos
}
