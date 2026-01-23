#define PORTES_IMPLEMENTATION
#include "portes.h"

int main(int argc, const char** argv) {
    if (argc != 2) {
        printf("Usage: ./portes [filepath]\n");
        return 1;
    }

    PortesString sourcepath = portes_string_from_cstr(argv[1]);
    FILE* source = portes_open_file(sourcepath, "r");
    if (!source) {
        fprintf(stderr, "couldn't locate source code file\n");
        portes_string_free(&sourcepath);
        exit(1);
    }

    FILE* output = portes_open_file(portes_string_from_cstr("output.txt"), "w");
    if (!output) {
        fprintf(stderr, "couldn't create output file\n");
        portes_close_file(source);
        portes_string_free(&sourcepath);
        exit(1);
    }

    PortesTokens tokens = {0};
    tokens = portes_generate_tokens(source, tokens);

    fprintf(output, "%-15s | %-20s | %s\n", "VALUE", "KIND", "LINE:COL");
    fprintf(output, "--------------------------------------------------\n");
    for (size_t i = 0; i < tokens.amount; ++i) {
        PortesToken t = tokens.items[i];
        fprintf(output, "%-15.*s | %-20s | %zu:%zu\n", (int)t.value.amount,
                t.value.items, portes_kind_to_cstr(t.kind), t.line, t.column);
    }

    printf("File 'output.txt' created and written successfully.\n");

    portes_tokens_free(&tokens);
    portes_string_free(&sourcepath);
    portes_close_file(output);
    portes_close_file(source);

    return 0;
}