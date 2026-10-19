package parser

import (
	"strings"
	"testing"

	"github.com/kauanallyson/portes/internal/ast"
	"github.com/kauanallyson/portes/internal/lexer"
)

func parse(t *testing.T, src string) *ast.Program {
	t.Helper()
	prog, err := Parse(lexer.Tokenize(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return prog
}

func TestParseProgram(t *testing.T) {
	src := `ALGORITMO main;
VARIAVEIS
  n1, n2 : INTEIRO;
  ok : LOGICO;
INICIO
  LEIA(n1, n2);
  ok <- VERDADEIRO;
  SE n1 + n2 * 2 >= 10 ENTAO
    ESCREVA("grande", n1);
  SENAO
    n1 <- -(n1 - 1) / 2;
  FIMSE
  ESCREVA();
FIM;`

	var sb strings.Builder
	ast.Fprint(&sb, parse(t, src))

	want := `ALGORITMO main
  VARIAVEIS
    n1, n2 : INTEIRO
    ok : LOGICO
  INICIO
    LEIA n1, n2
    ok <- VERDADEIRO
    SE (>= (+ n1 (* n2 2)) 10)
      ENTAO
        ESCREVA "grande", n1
      SENAO
        n1 <- (/ (- (- n1 1)) 2)
    ESCREVA
`
	if got := sb.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestLeftAssociativity(t *testing.T) {
	prog := parse(t, "ALGORITMO a; INICIO x <- 1 - 2 - 3; FIM;")
	got := ast.ExprString(prog.Body[0].(*ast.AssignStmt).Value)
	if want := "(- (- 1 2) 3)"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestSyntaxErrors(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"ALGORITMO a INICIO FIM;", "1:13: esperado ';', encontrado 'INICIO'"},
		{"ALGORITMO a; VARIAVEIS x : TEXTO; INICIO FIM;", "1:28: esperado um tipo (INTEIRO, REAL, STRING ou LOGICO), encontrado 'TEXTO'"},
		{"ALGORITMO a; INICIO x <- ; FIM;", "1:26: esperado uma expressão, encontrado ';'"},
		{"ALGORITMO a; INICIO x <- 1 FIM;", "1:28: esperado ';', encontrado 'FIM'"},
		{"ALGORITMO a; INICIO SE x ENTAO FIM;", "1:32: esperado 'FIMSE', encontrado 'FIM'"},
		{"ALGORITMO a; INICIO SE x ENTAO", "1:31: fim do arquivo inesperado, esperado 'FIMSE'"},
		{"ALGORITMO a; INICIO 42; FIM;", "1:21: comando inválido: '42'"},
		{"ALGORITMO a; INICIO x <- @; FIM;", "1:26: esperado uma expressão, encontrado símbolo desconhecido '@'"},
		{"ALGORITMO a; INICIO FIM; x", "1:26: esperado fim do arquivo, encontrado 'x'"},
	}
	for _, tt := range tests {
		_, err := Parse(lexer.Tokenize(tt.src))
		if err == nil {
			t.Errorf("%q: expected error %q, got none", tt.src, tt.want)
			continue
		}
		if err.Error() != tt.want {
			t.Errorf("%q:\n got  %q\n want %q", tt.src, err.Error(), tt.want)
		}
	}
}
