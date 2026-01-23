# Portes - v0.0.2

Transpilador de Português Estruturado (Portes) para C.
Esse projeto tenta facilitar o estudo de alunos que estão aprendendo Portes e querem migrar para C.

+ Como rodar o projeto:

``` bash
make # compila o projeto
./bin/portes caminho/do/arquivo.pe # para compilar o arquivo pe em um arquivo C
```

## ESTADO ATUAL

Por enquanto o portes apenas tokeniza os arquivos .pe.

Exemplo de entrada (.pe):

``` pe
ALGORITMO main;
VARIAVEIS
  n1, n2, n3, n4, soma : INTEIRO;
INICIO
  LEIA(n1, n2, n3, n4);
  soma <- n1 + n2 + n3 + n4;
  ESCREVA(soma);
FIM;
```

Saída da Tokenização:

``` txt
VALUE           | KIND                 | LINE:COL
--------------------------------------------------
ALGORITMO       | PORTES_ALGORITMO     | 1:1
main            | PORTES_IDENTIFIER    | 1:11
;               | PORTES_SEMICOLON     | 1:15
VARIAVEIS       | PORTES_VARIAVEIS     | 2:1
n1              | PORTES_IDENTIFIER    | 3:3
,               | PORTES_COMMA         | 3:5
n2              | PORTES_IDENTIFIER    | 3:7
,               | PORTES_COMMA         | 3:9
n3              | PORTES_IDENTIFIER    | 3:11
,               | PORTES_COMMA         | 3:13
n4              | PORTES_IDENTIFIER    | 3:15
,               | PORTES_COMMA         | 3:17
soma            | PORTES_IDENTIFIER    | 3:19
:               | PORTES_COLON         | 3:24
INTEIRO         | PORTES_INTEIRO       | 3:26
;               | PORTES_SEMICOLON     | 3:33
INICIO          | PORTES_INICIO        | 4:1
LEIA            | PORTES_LEIA          | 5:3
(               | PORTES_L_PAREN       | 5:7
n1              | PORTES_IDENTIFIER    | 5:8
,               | PORTES_COMMA         | 5:10
n2              | PORTES_IDENTIFIER    | 5:12
,               | PORTES_COMMA         | 5:14
n3              | PORTES_IDENTIFIER    | 5:16
,               | PORTES_COMMA         | 5:18
n4              | PORTES_IDENTIFIER    | 5:20
)               | PORTES_R_PAREN       | 5:22
;               | PORTES_SEMICOLON     | 5:23
soma            | PORTES_IDENTIFIER    | 6:3
<-              | PORTES_ASSIGN        | 6:8
n1              | PORTES_IDENTIFIER    | 6:11
+               | PORTES_PLUS          | 6:14
n2              | PORTES_IDENTIFIER    | 6:16
+               | PORTES_PLUS          | 6:19
n3              | PORTES_IDENTIFIER    | 6:21
+               | PORTES_PLUS          | 6:24
n4              | PORTES_IDENTIFIER    | 6:26
;               | PORTES_SEMICOLON     | 6:28
ESCREVA         | PORTES_ESCREVA       | 7:3
(               | PORTES_L_PAREN       | 7:10
soma            | PORTES_IDENTIFIER    | 7:11
)               | PORTES_R_PAREN       | 7:15
;               | PORTES_SEMICOLON     | 7:16
FIM             | PORTES_FIM           | 8:1
;               | PORTES_SEMICOLON     | 8:4
                | PORTES_EOF           | 9:0
```

Todos os tokens estão sendo mapeados corretamente. Próximo passo é parsear esses tokens e criar um programa .c que use essas exatas declarações.

## Boas práticas usadas

+ X-Macros
+ Single-Header Library Pattern
+ Arrays Dinâmicos
