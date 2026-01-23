#ifndef PORTES_H_
#define PORTES_H_

#include <ctype.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* --- TYPES & STRUCTURES --- */

typedef struct PortesString {
    const char* items;
    size_t capacity;
    size_t amount;
} PortesString;
typedef enum PortesTokenKind {
    PORTES_ALGORITMO,       // ALGORITMO
    PORTES_VARIAVEIS,       // VARIAVEIS
    PORTES_INICIO,          // INICIO
    PORTES_FIM,             // FIM
    PORTES_SE,              // SE
    PORTES_ENTAO,           // ENTAO
    PORTES_SENAO,           // SENAO
    PORTES_FIMSE,           // FIMSE
    PORTES_INTEIRO,         // INTEIRO
    PORTES_REAL,            // REAL
    PORTES_LOGICO,          // LOGICO
    PORTES_VERDADEIRO,      // VERDADEIRO
    PORTES_FALSO,           // FALSO
    PORTES_CARACTERE,       // 'a-Z'
    PORTES_STRING,          // "[a-Z]"
    PORTES_IDENTIFIER,      // VARIABLE NAME
    PORTES_LEIA,            // LEIA (VARIABLE NAME)
    PORTES_ESCREVA,         // ESCREVA (STRING)
    PORTES_NAO,             // NAO
    PORTES_E,               // E
    PORTES_OU,              // OU
    PORTES_POT,             // POT (x,y)
    PORTES_RAD,             // RAD (x)
    PORTES_MOD,             // x MOD y
    PORTES_DIV,             // x DIV y
    PORTES_PLUS,            // +
    PORTES_MINUS,           // -
    PORTES_MULTIPLICATION,  // *
    PORTES_DIVISION,        // /
    PORTES_EQUALS,          // =
    PORTES_L_THAN,          //
    PORTES_G_THAN,          // >
    PORTES_EL_THAN,         // <=
    PORTES_EG_THAN,         // >=
    PORTES_DIFFERENT,       // <>
    PORTES_SEMICOLON,       // ;
    PORTES_COLON,           // :
    PORTES_COMMA,           // ,
    PORTES_L_PAREN,         // (
    PORTES_R_PAREN,         // )
    PORTES_ASSIGN,          // <-
    PORTES_S_QUOTE,         // '
    PORTES_D_QUOTE,         // "
    PORTES_EOF,             // EOF
    PORTES_UNKWOWN          // UNRECOGNIZED TOKEN
} PortesTokenKind;

typedef struct PortesToken {
    PortesTokenKind kind;
    PortesString value;
    size_t line;
    size_t column;
} PortesToken;

typedef struct PortesTokens {
    PortesToken* items;
    size_t capacity;
    size_t amount;
} PortesTokens;

/* --- MACROS --- */

#define portes_da_append(xs, x)                                              \
    do {                                                                     \
        if ((xs).amount >= (xs).capacity) {                                  \
            size_t new_cap = ((xs).capacity == 0) ? 256 : (xs).capacity * 2; \
            void* new_ptr =                                                  \
                realloc((xs).items, new_cap * sizeof(*(xs).items));          \
            if (new_ptr == NULL) {                                           \
                fprintf(stderr, "Out of memory\n");                          \
                exit(1);                                                     \
            }                                                                \
            (xs).items = new_ptr;                                            \
            (xs).capacity = new_cap;                                         \
        }                                                                    \
        (xs).items[(xs).amount++] = (x);                                     \
    } while (0)

/* --- FUNCTION PROTOTYPES --- */

// PORTES STRING METHODS
PortesString portes_string_from_cstr(const char* cstr);
void portes_tokens_free(PortesTokens* ptks);
bool portes_string_cmp(PortesString str1, PortesString str2);
void portes_string_append_char(PortesString* str, char c);
void portes_string_free(PortesString* str);

// FILE MANAGEMENT
FILE* portes_open_file(PortesString filepath, const char* mode);
void portes_close_file(FILE* f);
void portes_print_file(FILE* f);
int portes_file_size(PortesString filepath);

// LEXER
int portes_peek(FILE* f);
const char* portes_kind_to_cstr(PortesTokenKind kind);
PortesToken portes_set_token_kind(PortesToken* token, PortesTokenKind kind,
                                  PortesString value);
PortesTokenKind portes_get_keyword_kind(PortesString str);
PortesTokens portes_generate_tokens(FILE* f, PortesTokens ptks);

/* --- IMPLEMENTATION --- */

#ifdef PORTES_IMPLEMENTATION

PortesString portes_string_from_cstr(const char* cstr) {
    PortesString str = {0};
    if (cstr == NULL) return str;
    for (size_t i = 0; cstr[i] != '\0'; ++i) {
        portes_string_append_char(&str, cstr[i]);
    }
    return str;
}

void portes_tokens_free(PortesTokens* ptks) {
    if (ptks->items) {
        for (size_t i = 0; i < ptks->amount; ++i) {
            portes_string_free(&ptks->items[i].value);
        }
        free(ptks->items);
        ptks->items = NULL;
    }
    ptks->amount = 0;
    ptks->capacity = 0;
}

bool portes_string_cmp(PortesString str1, PortesString str2) {
    if (str1.amount != str2.amount) return false;
    for (size_t i = 0; i < str1.amount; ++i) {
        if (str1.items[i] != str2.items[i]) return false;
    }
    return true;
}

void portes_string_append_char(PortesString* str, char c) {
    if (str->amount >= str->capacity) {
        if (str->capacity == 0)
            str->capacity = 256;
        else
            str->capacity *= 2;
        str->items =
            realloc((void*)str->items, str->capacity * sizeof(*str->items));
    }
    ((char*)str->items)[str->amount++] = c;
}

void portes_string_free(PortesString* str) {
    if (str->items) {
        free((void*)str->items);
        str->items = NULL;
        str->capacity = 0;
        str->amount = 0;
    }
}

FILE* portes_open_file(PortesString filepath, const char* mode) {
    char* cstr = malloc(filepath.amount + 1);
    memcpy(cstr, filepath.items, filepath.amount);
    cstr[filepath.amount] = '\0';

    FILE* fp = fopen(cstr, mode);
    free(cstr);

    if (fp == NULL) {
        fprintf(stderr, "Error loading the file\n");
    }
    return fp;
}

void portes_close_file(FILE* f) {
    if (f) fclose(f);
}

const char* portes_kind_to_cstr(PortesTokenKind kind) {
    switch (kind) {
        case PORTES_ALGORITMO:
            return "ALGORITMO";
        case PORTES_VARIAVEIS:
            return "VARIAVEIS";
        case PORTES_INICIO:
            return "INICIO";
        case PORTES_FIM:
            return "FIM";
        case PORTES_SE:
            return "SE";
        case PORTES_ENTAO:
            return "ENTAO";
        case PORTES_SENAO:
            return "SENAO";
        case PORTES_FIMSE:
            return "FIMSE";
        case PORTES_INTEIRO:
            return "INTEIRO";
        case PORTES_REAL:
            return "REAL";
        case PORTES_LOGICO:
            return "LOGICO";
        case PORTES_VERDADEIRO:
            return "VERDADEIRO";
        case PORTES_FALSO:
            return "FALSO";
        case PORTES_IDENTIFIER:
            return "IDENTIFIER";
        case PORTES_LEIA:
            return "LEIA";
        case PORTES_ESCREVA:
            return "ESCREVA";
        case PORTES_PLUS:
            return "PLUS";
        case PORTES_MINUS:
            return "MINUS";
        case PORTES_MULTIPLICATION:
            return "MULT";
        case PORTES_DIVISION:
            return "DIV";
        case PORTES_ASSIGN:
            return "ASSIGN";
        case PORTES_SEMICOLON:
            return "SEMICOLON";
        case PORTES_COLON:
            return "COLON";
        case PORTES_COMMA:
            return "COMMA";
        case PORTES_L_PAREN:
            return "L_PAREN";
        case PORTES_R_PAREN:
            return "R_PAREN";
        case PORTES_EOF:
            return "EOF";
        default:
            return "UNKNOWN";
    }
}

void portes_print_file(FILE* f) {
    if (!f) return;
    int c;
    while ((c = fgetc(f)) != EOF) {
        putchar(c);
    }
}

int portes_file_size(PortesString filepath) {
    char* cstr = malloc(filepath.amount + 1);
    memcpy(cstr, filepath.items, filepath.amount);
    cstr[filepath.amount] = '\0';

    FILE* fp = fopen(cstr, "rb");
    free(cstr);

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
    return c == EOF ? EOF : ungetc(c, f);
}

PortesToken portes_set_token_kind(PortesToken* token, PortesTokenKind kind,
                                  PortesString value) {
    token->kind = kind;
    token->value = value;
    return *token;
}

PortesTokenKind portes_get_keyword_kind(PortesString str) {
    static const struct {
        const char* keyword;
        size_t length;
        PortesTokenKind kind;
    } keywords[] = {
        {"ALGORITMO", 9, PORTES_ALGORITMO},
        {"VARIAVEIS", 9, PORTES_VARIAVEIS},
        {"INICIO", 6, PORTES_INICIO},
        {"FIM", 3, PORTES_FIM},
        {"SE", 2, PORTES_SE},
        {"ENTAO", 5, PORTES_ENTAO},
        {"SENAO", 5, PORTES_SENAO},
        {"FIMSE", 5, PORTES_FIMSE},
        {"INTEIRO", 7, PORTES_INTEIRO},
        {"REAL", 4, PORTES_REAL},
        {"LOGICO", 6, PORTES_LOGICO},
        {"VERDADEIRO", 10, PORTES_VERDADEIRO},
        {"FALSO", 5, PORTES_FALSO},
        {"LEIA", 4, PORTES_LEIA},
        {"ESCREVA", 7, PORTES_ESCREVA},
        {"NAO", 3, PORTES_NAO},
        {"E", 1, PORTES_E},
        {"OU", 2, PORTES_OU},
        {"POT", 3, PORTES_POT},
        {"RAD", 3, PORTES_RAD},
        {"MOD", 3, PORTES_MOD},
        {"DIV", 3, PORTES_DIV},
    };

    for (size_t i = 0; i < sizeof(keywords) / sizeof(keywords[0]); ++i) {
        if (str.amount == keywords[i].length &&
            memcmp(str.items, keywords[i].keyword, str.amount) == 0) {
            return keywords[i].kind;
        }
    }

    return PORTES_IDENTIFIER;
}

PortesTokens portes_generate_tokens(FILE* f, PortesTokens ptks) {
    // Note: Don't zero out ptks if you want to append to an existing list
    int c;
    size_t line = 1;
    size_t column = 0;

    while ((c = fgetc(f)) != EOF) {
        column++;

        if (isspace(c)) {
            if (c == '\n') {
                ++line;
                column = 0;
            }
            continue;
        }

        PortesToken token = {.line = line, .column = column};

        switch (c) {
            case '+':
                token.kind = PORTES_PLUS;
                token.value = portes_string_from_cstr("+");
                break;
            case '-':
                token.kind = PORTES_MINUS;
                token.value = portes_string_from_cstr("-");
                break;
            case '*':
                token.kind = PORTES_MULTIPLICATION;
                token.value = portes_string_from_cstr("*");
                break;
            case '/':
                token.kind = PORTES_DIVISION;
                token.value = portes_string_from_cstr("/");
                break;
            case ';':
                token.kind = PORTES_SEMICOLON;
                token.value = portes_string_from_cstr(";");
                break;
            case ':':
                token.kind = PORTES_COLON;
                token.value = portes_string_from_cstr(":");
                break;
            case ',':
                token.kind = PORTES_COMMA;
                token.value = portes_string_from_cstr(",");
                break;
            case '(':
                token.kind = PORTES_L_PAREN;
                token.value = portes_string_from_cstr("(");
                break;
            case ')':
                token.kind = PORTES_R_PAREN;
                token.value = portes_string_from_cstr(")");
                break;

            case '<': {
                int next = portes_peek(f);
                if (next == '-') {
                    fgetc(f);
                    column++;
                    token.kind = PORTES_ASSIGN;
                    token.value = portes_string_from_cstr("<-");
                } else if (next == '>') {
                    fgetc(f);
                    column++;
                    token.kind = PORTES_DIFFERENT;
                    token.value = portes_string_from_cstr("<>");
                } else if (next == '=') {
                    fgetc(f);
                    column++;
                    token.kind = PORTES_EL_THAN;
                    token.value = portes_string_from_cstr("<=");
                } else {
                    token.kind = PORTES_L_THAN;
                    token.value = portes_string_from_cstr("<");
                }
            } break;

            default:
                if (isalpha(c) || c == '_') {
                    PortesString word = {0};
                    portes_string_append_char(&word, (char)c);
                    while (isalnum(portes_peek(f)) || portes_peek(f) == '_') {
                        portes_string_append_char(&word, (char)fgetc(f));
                        column++;
                    }
                    token.kind = portes_get_keyword_kind(word);
                    token.value = word;
                } else if (isdigit(c)) {
                    PortesString num = {0};
                    portes_string_append_char(&num, (char)c);
                    while (isdigit(portes_peek(f))) {
                        portes_string_append_char(&num, (char)fgetc(f));
                        column++;
                    }
                    token.kind = PORTES_INTEIRO;
                    token.value = num;
                } else {
                    token.kind = PORTES_UNKWOWN;
                    char unknown_str[2] = {(char)c, '\0'};
                    token.value = portes_string_from_cstr(unknown_str);
                }
                break;
        }
        portes_da_append(ptks, token);
    }
    PortesToken eof = {
        .kind = PORTES_EOF, .line = line, .column = column, .value = {0}};
    portes_da_append(ptks, eof);
    return ptks;
}

#endif  // PORTES_IMPLEMENTATION
#endif  // PORTES_H_