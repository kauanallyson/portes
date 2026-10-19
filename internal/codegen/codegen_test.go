package codegen

import (
	"testing"

	"github.com/kauanallyson/portes/internal/check"
	"github.com/kauanallyson/portes/internal/lexer"
	"github.com/kauanallyson/portes/internal/parser"
)

func generate(t *testing.T, src string) string {
	t.Helper()
	prog, err := parser.Parse(lexer.Tokenize(src))
	if err != nil {
		t.Fatalf("unexpected syntax error: %v", err)
	}
	info, errs := check.Check(prog)
	if len(errs) > 0 {
		t.Fatalf("unexpected semantic errors: %v", errs)
	}
	return Generate(prog, info)
}

func TestGenerateProgram(t *testing.T) {
	src := `ALGORITMO media;
VARIAVEIS
  n1, n2 : INTEIRO;
  m : REAL;
  nome : STRING;
  ok : LOGICO;
INICIO
  LEIA(nome, n1, n2);
  m <- (n1 + n2) / 2.0;
  ok <- m >= 7;
  SE ok ENTAO
    ESCREVA("Aprovado: ", nome, " (100%)");
  SENAO
    ESCREVA(nome, " tirou ", m, " -> ", ok);
  FIMSE
  SE nome = "Ana" ENTAO
    nome <- "\a";
  FIMSE
  n1 <- -(-n2) * 2 - (n1 - 1);
  ESCREVA();
FIM;`

	want := `#include <stdbool.h>
#include <stdio.h>
#include <string.h>

/* ALGORITMO media */
int main(void) {
    int n1 = 0, n2 = 0;
    double m = 0;
    char nome[256] = "";
    bool ok = false;

    scanf("%255s", nome);
    scanf("%d", &n1);
    scanf("%d", &n2);
    m = (n1 + n2) / 2.0;
    ok = m >= 7;
    if (ok) {
        printf("Aprovado: %s (100%%)\n", nome);
    } else {
        printf("%s tirou %.15g -> %s\n", nome, m, ok ? "VERDADEIRO" : "FALSO");
    }
    if (strcmp(nome, "Ana") == 0) {
        strcpy(nome, "\\a");
    }
    n1 = -(-n2) * 2 - (n1 - 1);
    printf("\n");
    return 0;
}
`
	if got := generate(t, src); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestIncludesOnlyWhatIsUsed(t *testing.T) {
	got := generate(t, "ALGORITMO a; INICIO FIM;")
	want := `#include <stdio.h>

/* ALGORITMO a */
int main(void) {
    return 0;
}
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGenerateLoopsAndLogic(t *testing.T) {
	src := `ALGORITMO laco;
VARIAVEIS
  i, n : INTEIRO;
  ok : LOGICO;
INICIO
  ENQUANTO NAO ok E i < 10 FACA
    i <- i + 1;
    ok <- i MOD 3 = 0 OU i > 5 E n <> 2;
  FIMENQUANTO
  PARA i DE 1 ATE n FACA
    ESCREVA(i);
  FIMPARA
  PARA i DE 10 ATE 0 PASSO -2 FACA FIMPARA
  PARA i DE 0 ATE 9 PASSO 3 FACA FIMPARA
  PARA i DE n ATE 0 PASSO -1 FACA FIMPARA
  ok <- NAO i = 1 OU NAO (ok);
FIM;`

	want := `#include <stdbool.h>
#include <stdio.h>

/* ALGORITMO laco */
int main(void) {
    int i = 0, n = 0;
    bool ok = false;

    while (!ok && i < 10) {
        i = i + 1;
        ok = i % 3 == 0 || (i > 5 && n != 2);
    }
    int ate_i = n;
    for (i = 1; i <= ate_i; i++) {
        printf("%d\n", i);
    }
    for (i = 10; i >= 0; i -= 2) {
    }
    for (i = 0; i <= 9; i += 3) {
    }
    for (i = n; i >= 0; i--) {
    }
    ok = !(i == 1) || !(ok);
    return 0;
}
`
	if got := generate(t, src); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
