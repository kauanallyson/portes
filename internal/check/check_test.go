package check

import (
	"strings"
	"testing"

	"github.com/kauanallyson/portes/internal/lexer"
	"github.com/kauanallyson/portes/internal/parser"
)

// errorsOf parses src and returns the semantic errors, one per line.
func errorsOf(t *testing.T, src string) string {
	t.Helper()
	prog, err := parser.Parse(lexer.Tokenize(src))
	if err != nil {
		t.Fatalf("unexpected syntax error: %v", err)
	}
	_, errs := Check(prog)
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}

func TestValidProgram(t *testing.T) {
	src := `ALGORITMO a;
VARIAVEIS
  n : INTEIRO;
  x : REAL;
INICIO
  LEIA(n);
  x <- n * 2;
  ESCREVA("x = ", x);
FIM;`
	if got := errorsOf(t, src); got != "" {
		t.Errorf("expected no errors, got:\n%s", got)
	}
}

func TestDeclarations(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"ALGORITMO a; INICIO x <- 1; FIM;", "1:21: variável 'x' não declarada"},
		{"ALGORITMO a; VARIAVEIS x : INTEIRO; INICIO LEIA(x, y); FIM;", "1:52: variável 'y' não declarada"},
		{"ALGORITMO a; VARIAVEIS x, x : INTEIRO; INICIO FIM;", "1:27: variável 'x' já declarada"},
		{"ALGORITMO a; VARIAVEIS x : INTEIRO; x : REAL; INICIO FIM;", "1:37: variável 'x' já declarada"},
		{
			"ALGORITMO a; INICIO a1 <- 1; ESCREVA(a2); FIM;",
			"1:21: variável 'a1' não declarada\n1:38: variável 'a2' não declarada",
		},
	}
	for _, tt := range tests {
		if got := errorsOf(t, tt.src); got != tt.want {
			t.Errorf("%q:\n got  %q\n want %q", tt.src, got, tt.want)
		}
	}
}

func TestTypes(t *testing.T) {
	const header = "ALGORITMO a; VARIAVEIS n : INTEIRO; x : REAL; s : STRING; ok : LOGICO; INICIO\n"
	tests := []struct {
		stmt string // on line 2
		want string
	}{
		{`x <- n;`, ""},
		{`x <- n / 2 * 1.5 - -x;`, ""},
		{`ok <- n < x;`, ""},
		{`ok <- s = "a";`, ""},
		{`ok <- ok = VERDADEIRO;`, ""},
		{`s <- "texto";`, ""},
		{`ESCREVA(n, s, ok, x + 1);`, ""},
		{`n <- s;`, "2:1: não é possível atribuir STRING a 'n' (INTEIRO)"},
		{`n <- x;`, "2:1: não é possível atribuir REAL a 'n' (INTEIRO)"},
		{`n <- 1 + s;`, "2:8: operador '+' não se aplica a INTEIRO e STRING"},
		{`n <- 1 + s + 2;`, "2:8: operador '+' não se aplica a INTEIRO e STRING"},
		{`ok <- s < n;`, "2:9: operador '<' não se aplica a STRING e INTEIRO"},
		{`ok <- ok < FALSO;`, "2:10: operador '<' não se aplica a LOGICO e LOGICO"},
		{`n <- -s;`, "2:6: operador '-' não se aplica a STRING"},
		{`n <- z + 1;`, "2:6: variável 'z' não declarada"},
		{`SE n ENTAO FIMSE`, "2:4: condição do SE deve ser LOGICO, encontrado INTEIRO"},
		{`SE ok ENTAO n <- s; FIMSE`, "2:13: não é possível atribuir STRING a 'n' (INTEIRO)"},
		{`LEIA(ok);`, "2:6: LEIA não aceita a variável LOGICO 'ok'"},
	}
	for _, tt := range tests {
		if got := errorsOf(t, header+tt.stmt+" FIM;"); got != tt.want {
			t.Errorf("%q:\n got  %q\n want %q", tt.stmt, got, tt.want)
		}
	}
}

func TestLoopsAndLogic(t *testing.T) {
	const header = "ALGORITMO a; VARIAVEIS n : INTEIRO; x : REAL; s : STRING; ok : LOGICO; INICIO\n"
	tests := []struct {
		stmt string // on line 2
		want string
	}{
		{`ok <- ok E n > 1 OU NAO ok;`, ""},
		{`n <- n MOD 2;`, ""},
		{`PARA n DE 1 ATE 10 PASSO -2 FACA ESCREVA(n); FIMPARA`, ""},
		{`ENQUANTO NAO ok FACA ok <- VERDADEIRO; FIMENQUANTO`, ""},
		{`x <- x MOD 2;`, "2:8: operador 'MOD' não se aplica a REAL e INTEIRO"},
		{`ok <- ok E n;`, "2:10: operador 'E' não se aplica a LOGICO e INTEIRO"},
		{`ok <- NAO n;`, "2:7: operador 'NAO' não se aplica a INTEIRO"},
		{`n <- -ok;`, "2:6: operador '-' não se aplica a LOGICO"},
		{`ENQUANTO n FACA FIMENQUANTO`, "2:10: condição do ENQUANTO deve ser LOGICO, encontrado INTEIRO"},
		{`ENQUANTO ok FACA n <- s; FIMENQUANTO`, "2:18: não é possível atribuir STRING a 'n' (INTEIRO)"},
		{`PARA x DE 1 ATE 10 FACA FIMPARA`, "2:6: variável do PARA deve ser INTEIRO, 'x' é REAL"},
		{`PARA z DE 1 ATE 2 FACA FIMPARA`, "2:6: variável 'z' não declarada"},
		{
			`PARA n DE 1.5 ATE s FACA FIMPARA`,
			"2:11: início do PARA deve ser INTEIRO, encontrado REAL\n2:19: fim do PARA deve ser INTEIRO, encontrado STRING",
		},
		{`PARA n DE 1 ATE 10 PASSO n FACA FIMPARA`, "2:26: PASSO deve ser um número inteiro constante diferente de zero"},
		{`PARA n DE 1 ATE 10 PASSO 0 FACA FIMPARA`, "2:26: PASSO deve ser um número inteiro constante diferente de zero"},
		{`PARA n DE 1 ATE 2 FACA n <- x; FIMPARA`, "2:24: não é possível atribuir REAL a 'n' (INTEIRO)"},
	}
	for _, tt := range tests {
		if got := errorsOf(t, header+tt.stmt+" FIM;"); got != tt.want {
			t.Errorf("%q:\n got  %q\n want %q", tt.stmt, got, tt.want)
		}
	}
}
