package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kauanallyson/portes/internal/ast"
	"github.com/kauanallyson/portes/internal/lexer"
	"github.com/kauanallyson/portes/internal/parser"
	"github.com/kauanallyson/portes/internal/token"
)

func main() {
	showTokens := flag.Bool("tokens", false, "imprime a lista de tokens")
	showAST := flag.Bool("ast", false, "imprime a árvore sintática (padrão)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: portes [-tokens] [-ast] arquivo.ptes\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if !*showTokens {
		*showAST = true
	}

	path := flag.Arg(0)
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portes: %v\n", err)
		os.Exit(1)
	}

	tokens := lexer.Tokenize(string(src))
	if *showTokens {
		printTokens(os.Stdout, tokens)
	}

	if *showAST {
		prog, err := parser.Parse(tokens)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:%v\n", path, err)
			os.Exit(1)
		}
		ast.Fprint(os.Stdout, prog)
	}
}

func printTokens(w io.Writer, tokens []token.Token) {
	fmt.Fprintf(w, "%-15s | %-12s | %s\n", "VALUE", "KIND", "LINE:COL")
	fmt.Fprintln(w, "------------------------------------------")
	for _, t := range tokens {
		fmt.Fprintf(w, "%-15s | %-12s | %s\n", t.Value, t.Kind, t.Pos)
	}
}
