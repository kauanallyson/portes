package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kauanallyson/portes/internal/ast"
	"github.com/kauanallyson/portes/internal/check"
	"github.com/kauanallyson/portes/internal/codegen"
	"github.com/kauanallyson/portes/internal/lexer"
	"github.com/kauanallyson/portes/internal/parser"
	"github.com/kauanallyson/portes/internal/token"
)

func main() {
	showTokens := flag.Bool("tokens", false, "imprime a lista de tokens")
	showAST := flag.Bool("ast", false, "imprime a árvore sintática")
	output := flag.String("o", "", "grava o código C neste arquivo (padrão: saída padrão)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: portes [-tokens] [-ast] [-o saida.c] arquivo.ptes\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	path := flag.Arg(0)
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portes: %v\n", err)
		os.Exit(1)
	}

	if *showTokens || *showAST {
		tokens := lexer.Tokenize(string(src))
		if *showTokens {
			printTokens(os.Stdout, tokens)
		}
		if *showAST {
			prog, err := parser.Parse(tokens)
			if err != nil {
				fail(path, []error{err})
			}
			ast.Fprint(os.Stdout, prog)
		}
		return
	}

	csrc, errs := compile(string(src))
	if len(errs) > 0 {
		fail(path, errs)
	}
	if *output == "" {
		fmt.Print(csrc)
		return
	}
	if err := os.WriteFile(*output, []byte(csrc), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "portes: %v\n", err)
		os.Exit(1)
	}
}

// compile translates Portes source into C. On a syntax error it returns
// just that error; otherwise it returns every semantic error found.
func compile(src string) (string, []error) {
	prog, err := parser.Parse(lexer.Tokenize(src))
	if err != nil {
		return "", []error{err}
	}
	info, errs := check.Check(prog)
	if len(errs) > 0 {
		return "", errs
	}
	return codegen.Generate(prog, info), nil
}

func fail(path string, errs []error) {
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "%s:%v\n", path, err)
	}
	os.Exit(1)
}

func printTokens(w io.Writer, tokens []token.Token) {
	fmt.Fprintf(w, "%-15s | %-12s | %s\n", "VALUE", "KIND", "LINE:COL")
	fmt.Fprintln(w, "------------------------------------------")
	for _, t := range tokens {
		fmt.Fprintf(w, "%-15s | %-12s | %s\n", t.Value, t.Kind, t.Pos)
	}
}
