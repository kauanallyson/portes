package codegen

import (
	"sort"
	"strings"

	"github.com/kauanallyson/portes/internal/ast"
)

// reserved are names a Portes variable cannot keep in C: keywords (C99,
// C23 and GNU, since plain cc defaults to gnu17 or gnu23), the functions
// and macros the generated code uses or includes, and main. Names that
// start with _ are reserved too; see isReserved.
var reserved = map[string]bool{}

func init() {
	for _, name := range []string{
		// C99 keywords
		"auto", "break", "case", "char", "const", "continue", "default", "do",
		"double", "else", "enum", "extern", "float", "for", "goto", "if",
		"inline", "int", "long", "register", "restrict", "return", "short",
		"signed", "sizeof", "static", "struct", "switch", "typedef", "union",
		"unsigned", "void", "volatile", "while",
		// C23 and GNU keywords
		"alignas", "alignof", "constexpr", "nullptr", "static_assert",
		"thread_local", "typeof", "typeof_unqual", "asm",
		// GNU predefined macros
		"linux", "unix", "i386",
		// used by the generated code
		"main", "printf", "scanf", "strcpy", "strcmp",
		// <stdbool.h>
		"bool", "true", "false",
		// macros of <stdio.h>
		"NULL", "EOF", "BUFSIZ", "FILENAME_MAX", "FOPEN_MAX", "L_tmpnam", "TMP_MAX",
		"SEEK_SET", "SEEK_CUR", "SEEK_END", "stdin", "stdout", "stderr", "P_tmpdir", "L_ctermid",
		"getc", "putc", "getchar", "putchar",
	} {
		reserved[name] = true
	}
}

// names maps Portes identifiers to C identifiers, so that a variable
// called int becomes int_ (and _x becomes v_x) without clashing with
// another variable.
type names struct {
	c    map[string]string // Portes name -> C name
	used map[string]bool   // every C name taken so far
}

func newNames(vars map[string]ast.Type) *names {
	n := &names{c: map[string]string{}, used: map[string]bool{}}
	// Names that are fine in C are kept as they are and claimed first,
	// so a renamed variable never takes the name of another one.
	var clashing []string
	for v := range vars {
		if isReserved(v) {
			clashing = append(clashing, v)
			continue
		}
		n.c[v], n.used[v] = v, true
	}
	sort.Strings(clashing)
	for _, v := range clashing {
		if strings.HasPrefix(v, "_") {
			n.c[v] = n.fresh("v" + v)
		} else {
			n.c[v] = n.fresh(v + "_")
		}
	}
	return n
}

// isReserved reports whether name cannot be used as is in C. Every name
// starting with _ is: C keeps _X and __x for itself (_Bool, __LINE__...).
func isReserved(name string) bool {
	return reserved[name] || strings.HasPrefix(name, "_")
}

// of returns the C name of a Portes variable.
func (n *names) of(v string) string { return n.c[v] }

// fresh returns an unused C name based on base (base, base_, base__...)
// and claims it.
func (n *names) fresh(base string) string {
	name := base
	for n.used[name] || isReserved(name) {
		name += "_"
	}
	n.used[name] = true
	return name
}
