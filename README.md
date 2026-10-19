# Portes - v0.1

Transpilador de Português Estruturado (Portes) para C.
Esse projeto tenta facilitar o estudo de alunos que estão aprendendo Portes e querem migrar para C.

## Como usar

Requer Go 1.27+ e um compilador C (`cc`).

``` bash
go build .                          # compila para ./portes
./portes arquivo.ptes               # imprime o código C gerado
./portes -o saida.c arquivo.ptes    # grava o C em saida.c
cc -o prog saida.c && ./prog        # compila e roda
```

Por exemplo, `./portes examples/02_soma.ptes` imprime o C equivalente a [`02_soma.ptes`](examples/02_soma.ptes).

Erros apontam linha e coluna:

``` txt
main.ptes:6:3: variável 'total' não declarada
main.ptes:7:3: não é possível atribuir STRING a 'n1' (INTEIRO)
```

## Exemplos

A pasta [`examples/`](examples/) tem 15 programas, do mais simples ao mais completo. Os principais:

| Exemplo | O que mostra |
| --- | --- |
| [`02_soma.ptes`](examples/02_soma.ptes) | `LEIA`, aritmética e `ESCREVA` |
| [`03_media.ptes`](examples/03_media.ptes) | `REAL`, `STRING` e `SE`/`SENAO` aninhados |
| [`06_ano_bissexto.ptes`](examples/06_ano_bissexto.ptes) | `MOD` com `E`, `OU` e `NAO` |
| [`07_fatorial.ptes`](examples/07_fatorial.ptes) | laço `PARA` |
| [`09_contagem_regressiva.ptes`](examples/09_contagem_regressiva.ptes) | `PARA` com `PASSO` negativo |
| [`10_fibonacci.ptes`](examples/10_fibonacci.ptes) | laço `ENQUANTO` |
| [`11_primo.ptes`](examples/11_primo.ptes) | `ENQUANTO` com condição composta |
| [`12_senha.ptes`](examples/12_senha.ptes) | comparação de `STRING`s |

## A linguagem

+ Tipos: `INTEIRO`, `REAL`, `STRING`, `LOGICO`
+ Comandos: `<-`, `LEIA(...)`, `ESCREVA(...)`, `SE ... ENTAO ... [SENAO ...] FIMSE`,
  `ENQUANTO ... FACA ... FIMENQUANTO`, `PARA i DE a ATE b [PASSO n] FACA ... FIMPARA`
+ Operadores: `+ - * / MOD`, `= <> < > <= >=`, `E OU NAO`, parênteses
+ Literais: `42`, `3.14`, `"texto"`, `VERDADEIRO`, `FALSO`
+ Comentários: `// ...`

Regras que valem a pena saber:

+ Toda variável é declarada em `VARIAVEIS`, e os tipos precisam bater (só `INTEIRO` → `REAL` é automático).
+ Como em C, `INTEIRO / INTEIRO` é divisão inteira: `7 / 2` dá `3`; use `7 / 2.0`.
+ `NAO a = b` significa `NAO (a = b)`.
+ O `PASSO` do `PARA` é um inteiro constante (padrão 1).
+ `ESCREVA` pula linha no fim; `LEIA` de `STRING` lê uma palavra.

## Tradução para C

| Portes | C |
| --- | --- |
| `INTEIRO`, `REAL`, `LOGICO` | `int`, `double`, `bool` |
| `STRING` | `char nome[256]`, com `strcpy` e `strcmp` |
| `LEIA(a, b)` | um `scanf` por variável |
| `ESCREVA(...)` | um `printf` |
| `ENQUANTO c FACA` | `while (c)` |
| `PARA i DE a ATE b PASSO n` | `for (i = a; i <= b; i += n)` |
| `E`, `OU`, `NAO`, `MOD`, `<>`, `=` | `&&`, `\|\|`, `!`, `%`, `!=`, `==` |
| variável com nome reservado no C (`int`, `printf`...) | ganha `_`: `int_` |
