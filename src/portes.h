#ifndef PORTES_H_
#define PORTES_H_

#include <ctype.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* --- TYPES & STRUCTURES --- */

typedef enum PortesTokenType {
    PORTES_ALGORITMO,
    PORTES_VARIAVEIS,
    PORTES_INICIO,
    PORTES_FIM,
    PORTES_SE,
    PORTES_SENAO,
    PORTES_FIMSE,
    PORTES_INTEIRO,
    PORTES_REAL,
    PORTES_LOGICO,
    PORTES_CARACTERE,
    PORTES_LEIA,
    PORTES_ESCREVA,
    PORTES_IDENTIFIER,
    PORTES_ASSIGN,
    PORTES_PLUS,
    PORTES_MINUS,
    PORTES_MULT,
    PORTES_DIVISION,
    PORTES_LTHAN,
    PORTES_GTHAN,
    PORTES_SEMICOLON,
    PORTES_COLON,
    PORTES_COMMA,
    PORTES_LPAREN,
    PORTES_RPAREN,
    PORTES_EOF,
    PORTES_UNKWOWN
} PortesTokenType;

typedef struct PortesToken {
    PortesTokenType type;
    const char* value;
    size_t line;
    size_t column;
} PortesToken;

typedef struct PortesTokens {
    PortesToken* items;
    size_t capacity;
    size_t amount;
} PortesTokens;

/* --- MACROS --- */

#define portes_da_append(xs, x)                                           \
    do {                                                                  \
        if ((xs).amount >= (xs).capacity) {                               \
            if ((xs).capacity == 0)                                       \
                (xs).capacity = 256;                                      \
            else                                                          \
                (xs).capacity *= 2;                                       \
            (xs).items =                                                  \
                realloc((xs).items, (xs).capacity * sizeof(*(xs).items)); \
        }                                                                 \
        (xs).items[(xs).amount++] = (x);                                  \
    } while (0)

/* --- FUNCTION PROTOTYPES --- */

// STRING UTILS
void* portes_memset(void* str, int value);
size_t portes_str_size(const char* str);
bool porter_str_eq(const char* str);  // TODO

// FILE MANAGEMENT
FILE* portes_open_file(const char* filepath);
void portes_close_file(FILE* f);
void portes_print_file(FILE* f);
int portes_find_size(const char* filepath);

// LEXER
int portes_peek(FILE* f);
PortesTokens portes_generate_tokens(FILE* f, PortesTokens ptks);

/* --- IMPLEMENTATION --- */

#ifdef PORTES_IMPLEMENTATION

size_t portes_str_size(const char* str) {
    if (str == NULL) return 0;
    size_t size = 0;
    while (str[size] != '\0' && str[size] != '\n') {
        size++;
    }
    return size;
}

FILE* portes_open_file(const char* filepath) {
    FILE* fp = fopen(filepath, "r");
    if (fp == NULL) {
        fprintf(stderr, "Error loading the file\n");
    }
    return fp;
}

void portes_close_file(FILE* f) {
    if (f) fclose(f);
}

void portes_print_file(FILE* f) {
    if (!f) return;
    int c;
    while ((c = fgetc(f)) != EOF) {
        putchar(c);
    }
}

int portes_find_size(const char* filepath) {
    FILE* fp = fopen(filepath, "rb");
    if (fp == NULL) {
        fprintf(stderr, "File not found\n");
        return -1;
    }
    fseek(fp, 0L, SEEK_END);
    int res = (int)ftell(fp);
    fclose(fp);
    return res;
}

int portes_peek(FILE* f) {
    const int c = getc(f);
    c == EOF ? EOF : ungetc(c, f);
}

PortesTokens portes_generate_tokens(FILE* f, PortesTokens ptks) {
    char c = fgetc(f);
    int line = 0;
    while ((c = fgetc(f)) != EOF) {
        switch (c) {
            case '\n':
            case '\t':
            case ' ':
                line++;
                continue;
                // TODO rest of cases

            default:
                break;
        }
    }
}

#endif  // PORTES_IMPLEMENTATION
#endif  // PORTES_H_