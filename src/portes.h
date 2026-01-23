#ifndef PORTES_H_
#define PORTES_H_

#include <ctype.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define PORTES_TOKEN_LIST              \
    /* Keywords */                     \
    X(PORTES_ALGORITMO, "ALGORITMO")   \
    X(PORTES_VARIAVEIS, "VARIAVEIS")   \
    X(PORTES_INICIO, "INICIO")         \
    X(PORTES_FIM, "FIM")               \
    X(PORTES_SE, "SE")                 \
    X(PORTES_ENTAO, "ENTAO")           \
    X(PORTES_SENAO, "SENAO")           \
    X(PORTES_FIMSE, "FIMSE")           \
    X(PORTES_INTEIRO, "INTEIRO")       \
    X(PORTES_REAL, "REAL")             \
    X(PORTES_STRING, "STRING")         \
    X(PORTES_LOGICO, "LOGICO")         \
    X(PORTES_VERDADEIRO, "VERDADEIRO") \
    X(PORTES_FALSO, "FALSO")           \
    X(PORTES_LEIA, "LEIA")             \
    X(PORTES_ESCREVA, "ESCREVA")       \
    /* Operators & Symbols */          \
    X(PORTES_PLUS, "+")                \
    X(PORTES_MINUS, "-")               \
    X(PORTES_MULTIPLICATION, "*")      \
    X(PORTES_DIVISION, "/")            \
    X(PORTES_ASSIGN, "<-")             \
    X(PORTES_SEMICOLON, ";")           \
    X(PORTES_COLON, ":")               \
    X(PORTES_COMMA, ",")               \
    X(PORTES_L_PAREN, "(")             \
    X(PORTES_R_PAREN, ")")             \
    X(PORTES_L_THAN, "<")              \
    X(PORTES_G_THAN, ">")              \
    X(PORTES_LE_THAN, "<=")            \
    X(PORTES_GE_THAN, ">=")            \
    X(PORTES_DIFFERENT, "<>")          \
    /* Literals / Specials */          \
    X(PORTES_IDENTIFIER, "IDENTIFIER") \
    X(PORTES_EOF, "EOF")               \
    X(PORTES_UNKWOWN, "UNKNOWN")

/* --- TYPES & STRUCTURES --- */

typedef struct PortesString {
    const char* items;
    size_t capacity;
    size_t amount;
} PortesString;

typedef enum PortesTokenKind {
#define X(kind, name) kind,
    PORTES_TOKEN_LIST
#undef X
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

// PORTES STRING METHODS

PortesString portes_string_from_cstr(const char* cstr) {
    PortesString result = {0};
    if (!cstr) return result;

    size_t len = strlen(cstr);
    result.items = malloc(len);
    if (result.items) {
        memcpy((void*)result.items, cstr, len);
        result.amount = len;
        result.capacity = len;
    }
    return result;
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

// FILE MANAGEMENT

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

// LEXER

int portes_peek(FILE* f) {
    const int c = getc(f);
    return c == EOF ? EOF : ungetc(c, f);
}

const char* portes_kind_to_cstr(PortesTokenKind kind) {
    switch (kind) {
#define X(kind, name) \
    case kind:        \
        return #kind;
        PORTES_TOKEN_LIST
#undef X
        default:
            return "UNKNOWN";
    }
}

PortesToken portes_set_token_kind(PortesToken* token, PortesTokenKind kind,
                                  PortesString value) {
    token->kind = kind;
    token->value = value;
    return *token;
}

PortesTokenKind portes_get_keyword_kind(PortesString str) {
#define X(kind, name)                             \
    if (str.amount == strlen(name) &&             \
        memcmp(str.items, name, str.amount) == 0) \
        return kind;

    X(PORTES_ALGORITMO, "ALGORITMO")
    X(PORTES_VARIAVEIS, "VARIAVEIS")
    X(PORTES_INICIO, "INICIO")
    X(PORTES_FIM, "FIM")
    X(PORTES_INTEIRO, "INTEIRO")
    X(PORTES_REAL, "REAL")
    X(PORTES_LEIA, "LEIA")
    X(PORTES_ESCREVA, "ESCREVA")
    X(PORTES_SE, "SE")
    X(PORTES_ENTAO, "ENTAO")
    X(PORTES_SENAO, "SENAO")
    X(PORTES_FIMSE, "FIMSE")
#undef X

    return PORTES_IDENTIFIER;
}

PortesTokens portes_generate_tokens(FILE* f, PortesTokens ptks) {
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
                    token.kind = PORTES_LE_THAN;
                    token.value = portes_string_from_cstr("<=");
                } else {
                    token.kind = PORTES_L_THAN;
                    token.value = portes_string_from_cstr("<");
                }
            } break;
            case '"': {
                PortesString str_val = {0};
                // We don't include the opening quote in the value
                while (portes_peek(f) != '"' && portes_peek(f) != EOF) {
                    char next_c = (char)fgetc(f);
                    column++;
                    if (next_c == '\n') {
                        line++;
                        column = 0;
                    }
                    portes_string_append_char(&str_val, next_c);
                }

                if (fgetc(f) == '"') {  // Consume closing quote
                    column++;
                    token.kind = PORTES_STRING;
                    token.value = str_val;
                } else {
                    // Handle unclosed string error
                    token.kind = PORTES_UNKWOWN;
                    token.value = str_val;
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