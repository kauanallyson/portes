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
        exit(1);
    }

    PortesString outputpath = portes_string_from_cstr("output.txt");
    FILE* output = portes_open_file(outputpath, "w");
    if (!output) {
        fprintf(stderr, "couldn't create output file\n");
        exit(1);
    }

    PortesTokens tokens = {0};
    tokens = portes_generate_tokens(source, tokens);

    fprintf(output, "%-20s | %-15s | %s\n", "VALUE", "KIND", "LINE:COL");
    fprintf(output, "--------------------------------------------------\n");
    for (size_t i = 0; i < tokens.amount; ++i) {
        PortesToken t = tokens.items[i];

        fprintf(output, "%-20.*s | %-15s | %zu:%zu\n", (int)t.value.amount,
                t.value.items, portes_kind_to_cstr(t.kind), t.line, t.column);
    }

    printf("File '%s' created and written successfully.\n", outputpath.items);

    portes_tokens_free(&tokens);
    portes_string_free(&sourcepath);
    portes_close_file(output);
    portes_close_file(source);
    return 0;
}