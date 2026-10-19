package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// run transpiles src, compiles the C with cc and runs it with the given stdin.
func run(t *testing.T, src, stdin string) string {
	t.Helper()
	cmd := exec.Command(build(t, src))
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("program failed: %v", err)
	}
	return string(out)
}

// build transpiles src and compiles the C with cc, returning the binary path.
func build(t *testing.T, src string) string {
	t.Helper()
	cc, err := exec.LookPath("cc")
	if err != nil {
		t.Skip("cc not found")
	}
	csrc, errs := compile(src)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	dir := t.TempDir()
	cfile, bin := filepath.Join(dir, "prog.c"), filepath.Join(dir, "prog")
	if err := os.WriteFile(cfile, []byte(csrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(cc, "-Wall", "-Wextra", "-Werror", "-o", bin, cfile).CombinedOutput(); err != nil {
		t.Fatalf("cc failed: %v\n%s\nC source:\n%s", err, out, csrc)
	}
	return bin
}

// TestExamplesCompile keeps every program in examples/ valid.
func TestExamplesCompile(t *testing.T) {
	files, err := filepath.Glob("examples/*.ptes")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no examples found")
	}
	for _, f := range files {
		t.Run(filepath.Base(f), func(t *testing.T) {
			src, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			build(t, string(src))
		})
	}
}

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name, src, stdin, want string
	}{
		{
			name: "soma",
			src: `ALGORITMO main;
VARIAVEIS
  n1, n2, n3, n4, soma : INTEIRO;
INICIO
  LEIA(n1, n2, n3, n4);
  soma <- n1 + n2 + n3 + n4;
  ESCREVA(soma);
FIM;`,
			stdin: "1 2 3 4\n",
			want:  "10\n",
		},
		{
			name: "media e strings",
			src: `ALGORITMO media;
VARIAVEIS
  n1, n2 : INTEIRO;
  m : REAL;
  nome : STRING;
INICIO
  LEIA(nome, n1, n2);
  m <- (n1 + n2) / 2.0;
  SE m >= 7 ENTAO
    ESCREVA(nome, ": aprovado com ", m, " (100%)");
  SENAO
    ESCREVA(nome, ": reprovado");
  FIMSE
  SE nome = "Ana" ENTAO
    nome <- "Ana \A";
  FIMSE
  ESCREVA(nome, " ", 7 / 2, " ", 1 - 2 - 3, " ", -(-n1), " ", n1 = 8);
FIM;`,
			stdin: "Ana 8 7\n",
			want:  "Ana: aprovado com 7.5 (100%)\nAna \\A 3 -4 8 VERDADEIRO\n",
		},
		{
			name: "lacos e logica",
			src: `ALGORITMO lacos;
VARIAVEIS
  i, n, fat : INTEIRO;
  achou : LOGICO;
INICIO
  LEIA(n);
  fat <- 1;
  PARA i DE 2 ATE n FACA
    fat <- fat * i;
  FIMPARA
  ESCREVA(n, "! = ", fat);
  PARA i DE 10 ATE 1 PASSO -3 FACA
    ESCREVA(i);
  FIMPARA
  i <- 0;
  achou <- FALSO;
  ENQUANTO NAO achou E i < 100 FACA
    i <- i + 1;
    SE i MOD 7 = 0 E i MOD 5 = 0 ENTAO
      achou <- VERDADEIRO;
    FIMSE
  FIMENQUANTO
  ESCREVA(i, " ", achou, " ", NAO achou OU i > 50);
FIM;`,
			stdin: "5\n",
			want:  "5! = 120\n10\n7\n4\n1\n35 VERDADEIRO FALSO\n",
		},
		{
			name: "nomes reservados do C",
			src: `ALGORITMO main;
VARIAVEIS
  int, int_, printf, EOF, ate_i, i : INTEIRO;
  bool : LOGICO;
INICIO
  int <- 1;
  int_ <- 2;
  printf <- int + int_;
  EOF <- printf * 10;
  bool <- EOF > 0;
  ate_i <- 7;
  PARA i DE 1 ATE int_ FACA
    ESCREVA(int, " ", int_, " ", printf, " ", EOF, " ", bool, " ", ate_i);
  FIMPARA
FIM;`,
			want: "1 2 3 30 VERDADEIRO 7\n1 2 3 30 VERDADEIRO 7\n",
		},
		{
			name: "nomes reservados do C com _ e do GNU/C23",
			src: `ALGORITMO a;
VARIAVEIS
  _x, v_x, __LINE__, __func__, _IOFBF, _Alignas : INTEIRO;
  typeof, asm, P_tmpdir, linux, unix, nullptr, constexpr, static_assert, alignof : INTEIRO;
INICIO
  _x <- 1; v_x <- 2; __LINE__ <- 3; __func__ <- 4; _IOFBF <- 5; _Alignas <- 6;
  typeof <- 7; asm <- 8; P_tmpdir <- 9; linux <- 10; unix <- 11;
  nullptr <- 12; constexpr <- 13; static_assert <- 14; alignof <- 15;
  ESCREVA(_x + v_x + __LINE__ + __func__ + _IOFBF + _Alignas + typeof + asm + P_tmpdir
    + linux + unix + nullptr + constexpr + static_assert + alignof);
FIM;`,
			want: "120\n",
		},
		{
			name: "limite do PARA avaliado uma vez",
			src: `ALGORITMO a;
VARIAVEIS
  i, n : INTEIRO;
INICIO
  n <- 3;
  PARA i DE 1 ATE n FACA
    n <- n + 1;
    ESCREVA(i, " ", n);
  FIMPARA
  PARA i DE n ATE n - 2 PASSO -1 FACA
    ESCREVA(i);
  FIMPARA
FIM;`,
			want: "1 4\n2 5\n3 6\n6\n5\n4\n",
		},
		{
			name: "REAL sem notacao cientifica",
			src: `ALGORITMO a;
VARIAVEIS
  x : REAL;
INICIO
  x <- 1234567.5;
  ESCREVA(x, " ", 7 / 2.0, " ", 0.1 + 0.2, " ", 2.0);
FIM;`,
			want: "1234567.5 3.5 0.3 2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := run(t, tt.src, tt.stdin); got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestCompileErrors(t *testing.T) {
	_, errs := compile("ALGORITMO a; INICIO x <- 1; y <- 2; FIM;")
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %v", errs)
	}
	_, errs = compile("ALGORITMO a INICIO FIM;")
	if len(errs) != 1 || errs[0].Error() != "1:13: esperado ';', encontrado 'INICIO'" {
		t.Fatalf("unexpected errors: %v", errs)
	}
}
