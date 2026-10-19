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

func TestParseLoops(t *testing.T) {
	src := `ALGORITMO laco;
INICIO
  ENQUANTO i < 10 FACA
    i <- i + 1;
  FIMENQUANTO
  PARA i DE 1 ATE n FACA
    ESCREVA(i);
  FIMPARA;
  PARA i DE 10 ATE 0 PASSO -2 FACA
  FIMPARA
FIM;`

	var sb strings.Builder
	ast.Fprint(&sb, parse(t, src))

	want := `ALGORITMO laco
  INICIO
    ENQUANTO (< i 10)
      FACA
        i <- (+ i 1)
    PARA i DE 1 ATE n
      FACA
        ESCREVA i
    PARA i DE 10 ATE 0 PASSO (- 2)
      FACA
`
	if got := sb.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestLogicalPrecedence(t *testing.T) {
	tests := []struct{ src, want string }{
		{"NAO a = b E c OU d MOD 2 = 0", "(OU (E (NAO (= a b)) c) (= (MOD d 2) 0))"},
		{"a OU b OU c", "(OU (OU a b) c)"},
		{"a E b OU c E d", "(OU (E a b) (E c d))"},
		{"NAO NAO a", "(NAO (NAO a))"},
		{"a * b MOD c / d", "(/ (MOD (* a b) c) d)"},
		{"-a MOD 2", "(MOD (- a) 2)"},
	}
	for _, tt := range tests {
		prog := parse(t, "ALGORITMO a; INICIO x <- "+tt.src+"; FIM;")
		if got := ast.ExprString(prog.Body[0].(*ast.AssignStmt).Value); got != tt.want {
			t.Errorf("%s: got %s, want %s", tt.src, got, tt.want)
		}
	}
}

func TestLoopSyntaxErrors(t *testing.T) {
	tests := []struct{ src, want string }{
		{"ALGORITMO a; INICIO ENQUANTO x FACA FIM;", "1:37: esperado 'FIMENQUANTO', encontrado 'FIM'"},
		{"ALGORITMO a; INICIO ENQUANTO x i <- 1; FIMENQUANTO FIM;", "1:32: esperado 'FACA', encontrado 'i'"},
		{"ALGORITMO a; INICIO PARA i 1 ATE 2 FACA FIMPARA FIM;", "1:28: esperado 'DE', encontrado '1'"},
		{"ALGORITMO a; INICIO PARA i DE 1 FACA FIMPARA FIM;", "1:33: esperado 'ATE', encontrado 'FACA'"},
		{"ALGORITMO a; INICIO PARA 1 DE 1 ATE 2 FACA FIMPARA FIM;", "1:26: esperado um identificador, encontrado '1'"},
		{"ALGORITMO a; INICIO SE x ENTAO FIMPARA FIM;", "1:32: esperado 'FIMSE', encontrado 'FIMPARA'"},
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
