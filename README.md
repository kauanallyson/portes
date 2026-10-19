# Portes - v0.1

Transpilador de Português Estruturado (Portes) para C.
Esse projeto tenta facilitar o estudo de alunos que estão aprendendo Portes e querem migrar para C.

+ Como rodar o projeto (requer Go 1.27+):

``` bash
go build .                              # compila para ./portes
go test ./...                           # roda os testes
./portes caminho/arquivo.ptes           # imprime o código C gerado
./portes -o saida.c arquivo.ptes        # grava o C em saida.c
cc -o prog saida.c && ./prog            # compila e roda o C gerado
./portes -ast caminho/arquivo.ptes      # imprime a árvore sintática (AST)
./portes -tokens caminho/arquivo.ptes   # imprime os tokens
```

## Estrutura

``` txt
internal/
  token/   tipos de token e palavras-chave
  lexer/   texto .ptes -> lista de tokens
  ast/     nós da árvore sintática + impressão
  parser/  tokens -> AST (descida recursiva; gramática abaixo)
  check/   análise semântica: variáveis declaradas e tipos das expressões
  codegen/ AST verificada -> código C
main.go    ponto de entrada (main_test.go compila e roda exemplos com cc)
```

## ESTADO ATUAL

Pipeline completo: lexer → parser → verificação semântica → geração de C.

Exemplo de entrada (.ptes):

``` ptes
ALGORITMO main;
VARIAVEIS
  n1, n2, n3, n4, soma : INTEIRO;
INICIO
  LEIA(n1, n2, n3, n4);
  soma <- n1 + n2 + n3 + n4;
  ESCREVA(soma);
FIM;
```

Saída de `./portes main.ptes`:

``` c
#include <stdio.h>

/* ALGORITMO main */
int main(void) {
    int n1 = 0, n2 = 0, n3 = 0, n4 = 0, soma = 0;

    scanf("%d", &n1);
    scanf("%d", &n2);
    scanf("%d", &n3);
    scanf("%d", &n4);
    soma = n1 + n2 + n3 + n4;
    printf("%d\n", soma);
    return 0;
}
```

Saída de `./portes -ast main.ptes` (expressões em notação prefixa, mostrando a precedência):

``` txt
ALGORITMO main
  VARIAVEIS
    n1, n2, n3, n4, soma : INTEIRO
  INICIO
    LEIA n1, n2, n3, n4
    soma <- (+ (+ (+ n1 n2) n3) n4)
    ESCREVA soma
```

Erros de sintaxe apontam linha e coluna:

``` txt
main.ptes:7:3: esperado ';', encontrado 'ESCREVA'   # faltou o ';' no fim da linha 6
```

Erros semânticos também, todos de uma vez:

``` txt
main.ptes:6:3: variável 'total' não declarada
main.ptes:7:3: não é possível atribuir STRING a 'n1' (INTEIRO)
```

## O que a linguagem suporta hoje

+ Tipos: `INTEIRO`, `REAL`, `STRING`, `LOGICO`
+ Comandos: atribuição (`<-`), `LEIA(...)`, `ESCREVA(...)`, `SE ... ENTAO ... [SENAO ...] FIMSE`,
  `ENQUANTO ... FACA ... FIMENQUANTO`, `PARA i DE a ATE b [PASSO n] FACA ... FIMPARA`
+ Expressões: `+ - * / MOD`, `-` unário, `= <> < > <= >=`, `E OU NAO`, parênteses
+ Literais: inteiros, reais (`3.14`), strings (`"texto"`), `VERDADEIRO`/`FALSO`
+ Comentários de linha: `// ...`

## Gramática (EBNF)

Notação: `{ x }` = zero ou mais vezes, `[ x ]` = opcional, `|` = alternativa, `"x"` = texto literal.

``` ebnf
program  = "ALGORITMO" IDENT ";" [ "VARIAVEIS" { varDecl } ] "INICIO" { stmt } "FIM" ";" EOF ;
varDecl  = IDENT { "," IDENT } ":" type ";" ;
type     = "INTEIRO" | "REAL" | "STRING" | "LOGICO" ;

stmt     = assign | read | write | if | while | for ;
assign   = IDENT "<-" expr ";" ;
read     = "LEIA" "(" IDENT { "," IDENT } ")" ";" ;
write    = "ESCREVA" "(" [ expr { "," expr } ] ")" ";" ;
if       = "SE" expr "ENTAO" { stmt } [ "SENAO" { stmt } ] "FIMSE" [ ";" ] ;
while    = "ENQUANTO" expr "FACA" { stmt } "FIMENQUANTO" [ ";" ] ;
for      = "PARA" IDENT "DE" expr "ATE" expr [ "PASSO" expr ] "FACA" { stmt } "FIMPARA" [ ";" ] ;

expr     = and { "OU" and } ;
and      = not { "E" not } ;
not      = "NAO" not | compare ;
compare  = additive [ ( "=" | "<>" | "<" | ">" | "<=" | ">=" ) additive ] ;
additive = term { ( "+" | "-" ) term } ;
term     = unary { ( "*" | "/" | "MOD" ) unary } ;
unary    = "-" unary | primary ;
primary  = INT_LIT | REAL_LIT | STRING_LIT | "VERDADEIRO" | "FALSO" | IDENT | "(" expr ")" ;
```

Tokens (definidos no lexer):

``` ebnf
IDENT      = letra { letra | digito } ;        (* letra = a-z, A-Z ou _ *)
INT_LIT    = digito { digito } ;
REAL_LIT   = digito { digito } "." digito { digito } ;
STRING_LIT = '"' { qualquer caractere exceto '"' } '"' ;
```

A ordem das regras de expressão define a precedência: `*`, `/` e `MOD` (em `term`) são avaliados antes de `+` e `-` (em `additive`), que vêm antes das comparações (em `compare`), depois `NAO`, `E` e por último `OU`. Assim `NAO a = b` significa `NAO (a = b)`.

## Regras de tipo

+ Toda variável deve ser declarada uma única vez em `VARIAVEIS`.
+ `+ - * /` aceitam `INTEIRO` e `REAL`; o resultado é `REAL` se algum lado for `REAL`. Como em C, `INTEIRO / INTEIRO` é divisão inteira (`7 / 2` dá `3`; use `7 / 2.0`).
+ `MOD` só aceita `INTEIRO`; `E`, `OU`, `NAO` só aceitam `LOGICO`.
+ `= <>` comparam valores do mesmo tipo (ou dois números); `< > <= >=` comparam números ou `STRING`s.
+ A condição de `SE` e `ENQUANTO` deve ser `LOGICO`.
+ Atribuição exige o mesmo tipo, exceto `INTEIRO` → `REAL`.
+ `PARA`: a variável, o início e o fim são `INTEIRO`; o `PASSO` (padrão 1) é um inteiro constante diferente de zero.
+ `LEIA` não aceita variáveis `LOGICO`.

## Tradução para C

| Portes | C |
| --- | --- |
| `INTEIRO`, `REAL`, `LOGICO` | `int`, `double`, `bool` (`<stdbool.h>`) |
| `STRING` | `char nome[256]`; atribuição com `strcpy`, comparação com `strcmp` |
| `LEIA(a, b)` | um `scanf` por variável (`%d`, `%lf`, `%255s` — lê uma palavra) |
| `ESCREVA(...)` | um `printf` com os argumentos juntos e `\n` no fim; `REAL` com `%.15g` (`7.5`, `1234567.5`); `LOGICO` sai como `VERDADEIRO`/`FALSO` |
| `ENQUANTO c FACA` | `while (c)` |
| `PARA i DE a ATE b PASSO n` | `for (i = a; i <= b; i += n)` (`>=` e `-=` quando `n` é negativo); se `b` não é constante, ele é calculado uma vez antes do laço (`int ate_i = b;`) |
| `E`, `OU`, `NAO`, `MOD`, `<>`, `=` | `&&`, `\|\|`, `!`, `%`, `!=`, `==` |
| variável com nome reservado no C (`int`, `printf`, `EOF`, `typeof`...) | ganha `_` no fim: `int_`; se começa com `_`, ganha `v` no início: `_x` → `v_x` (sem colidir com outras variáveis) |
