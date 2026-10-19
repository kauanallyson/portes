package lexer

import (
	"testing"

	"github.com/kauanallyson/portes/internal/token"
)

func tk(kind token.Kind, value string, line, column int) token.Token {
	return token.Token{Kind: kind, Value: value, Pos: token.Pos{Line: line, Column: column}}
}

func TestTokenizeProgram(t *testing.T) {
	src := "ALGORITMO main;\r\n" +
		"VARIAVEIS\n" +
		"  soma : INTEIRO;\n" +
		"INICIO\n" +
		"  soma <- n1 + 2.5; // comentario\n" +
		"FIM;\n"

	want := []token.Token{
		tk(token.ALGORITMO, "ALGORITMO", 1, 1),
		tk(token.IDENTIFIER, "main", 1, 11),
		tk(token.SEMICOLON, ";", 1, 15),
		tk(token.VARIAVEIS, "VARIAVEIS", 2, 1),
		tk(token.IDENTIFIER, "soma", 3, 3),
		tk(token.COLON, ":", 3, 8),
		tk(token.INTEIRO, "INTEIRO", 3, 10),
		tk(token.SEMICOLON, ";", 3, 17),
		tk(token.INICIO, "INICIO", 4, 1),
		tk(token.IDENTIFIER, "soma", 5, 3),
		tk(token.ASSIGN, "<-", 5, 8),
		tk(token.IDENTIFIER, "n1", 5, 11),
		tk(token.PLUS, "+", 5, 14),
		tk(token.REAL_LIT, "2.5", 5, 16),
		tk(token.SEMICOLON, ";", 5, 19),
		tk(token.FIM, "FIM", 6, 1),
		tk(token.SEMICOLON, ";", 6, 4),
		tk(token.EOF, "", 7, 1),
	}

	got := Tokenize(src)
	if len(got) != len(want) {
		t.Fatalf("got %d tokens, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestOperatorsAndLiterals(t *testing.T) {
	tests := []struct {
		src   string
		kind  token.Kind
		value string
	}{
		{"<-", token.ASSIGN, "<-"},
		{"<>", token.DIFFERENT, "<>"},
		{"<=", token.LE_THAN, "<="},
		{"<", token.L_THAN, "<"},
		{">=", token.GE_THAN, ">="},
		{">", token.G_THAN, ">"},
		{"=", token.EQUAL, "="},
		{"42", token.INT_LIT, "42"},
		{"3.14", token.REAL_LIT, "3.14"},
		{`"Olá, mundo"`, token.STRING_LIT, "Olá, mundo"},
		{`"sem fim`, token.UNKNOWN, `"sem fim`},
		{"VERDADEIRO", token.VERDADEIRO, "VERDADEIRO"},
		{"LOGICO", token.LOGICO, "LOGICO"},
		{"ENQUANTO", token.ENQUANTO, "ENQUANTO"},
		{"FACA", token.FACA, "FACA"},
		{"FIMENQUANTO", token.FIMENQUANTO, "FIMENQUANTO"},
		{"PARA", token.PARA, "PARA"},
		{"DE", token.DE, "DE"},
		{"ATE", token.ATE, "ATE"},
		{"PASSO", token.PASSO, "PASSO"},
		{"FIMPARA", token.FIMPARA, "FIMPARA"},
		{"E", token.E, "E"},
		{"OU", token.OU, "OU"},
		{"NAO", token.NAO, "NAO"},
		{"MOD", token.MOD, "MOD"},
		{"Enquanto", token.IDENTIFIER, "Enquanto"},
		{"e", token.IDENTIFIER, "e"},
		{"fimse", token.IDENTIFIER, "fimse"},
		{"@", token.UNKNOWN, "@"},
	}
	for _, tt := range tests {
		got := New(tt.src).Next()
		if got.Kind != tt.kind || got.Value != tt.value {
			t.Errorf("%q: got %s %q, want %s %q", tt.src, got.Kind, got.Value, tt.kind, tt.value)
		}
	}
}
