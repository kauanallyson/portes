# Portes - v0.0.3

Transpilador de Português Estruturado (Portes) para C.
Esse projeto tenta facilitar o estudo de alunos que estão aprendendo Portes e querem migrar para C.

+ Como rodar o projeto (requer Go 1.27+):

``` bash
go build .                              # compila para ./portes
go test ./...                           # roda os testes
./portes caminho/arquivo.ptes           # imprime a árvore sintática (AST)
./portes -tokens caminho/arquivo.ptes   # imprime os tokens
```

## Estrutura

``` txt
internal/
  token/   tipos de token e palavras-chave
  lexer/   texto .ptes -> lista de tokens
  ast/     nós da árvore sintática + impressão
  parser/  tokens -> AST (descida recursiva; gramática no topo de parser.go)
main.go    ponto de entrada
```

## ESTADO ATUAL

Lexer e parser prontos; o próximo passo é gerar o código C a partir da AST.

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

Saída de `./portes main.ptes` (expressões em notação prefixa, mostrando a precedência):

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

## O que a linguagem suporta hoje

+ Tipos: `INTEIRO`, `REAL`, `STRING`, `LOGICO`
+ Comandos: atribuição (`<-`), `LEIA(...)`, `ESCREVA(...)`, `SE ... ENTAO ... [SENAO ...] FIMSE`
+ Expressões: `+ - * /`, `-` unário, `= <> < > <= >=`, parênteses
+ Literais: inteiros, reais (`3.14`), strings (`"texto"`), `VERDADEIRO`/`FALSO`
+ Comentários de linha: `// ...`

## Gramática (EBNF)

Notação: `{ x }` = zero ou mais vezes, `[ x ]` = opcional, `|` = alternativa, `"x"` = texto literal.

``` ebnf
program  = "ALGORITMO" IDENT ";" [ "VARIAVEIS" { varDecl } ] "INICIO" { stmt } "FIM" ";" EOF ;
varDecl  = IDENT { "," IDENT } ":" type ";" ;
type     = "INTEIRO" | "REAL" | "STRING" | "LOGICO" ;

stmt     = assign | read | write | if ;
assign   = IDENT "<-" expr ";" ;
read     = "LEIA" "(" IDENT { "," IDENT } ")" ";" ;
write    = "ESCREVA" "(" [ expr { "," expr } ] ")" ";" ;
if       = "SE" expr "ENTAO" { stmt } [ "SENAO" { stmt } ] "FIMSE" [ ";" ] ;

expr     = additive [ ( "=" | "<>" | "<" | ">" | "<=" | ">=" ) additive ] ;
additive = term { ( "+" | "-" ) term } ;
term     = unary { ( "*" | "/" ) unary } ;
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

A ordem das regras de expressão define a precedência: `*` e `/` (em `term`) são avaliados antes de `+` e `-` (em `additive`), que vêm antes das comparações (em `expr`).
