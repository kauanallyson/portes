#define PORTES_IMPLEMENTATION
#include "portes.h"

int main(int argc, const char** argv) {
    if (argc != 2) {
        printf("Usage: ./portes [filepath]\n");
        exit(1);
    }

    const char* filepath = argv[1];
    FILE* file = portes_open_file(filepath);
    portes_print_file(file);
    portes_close_file(file);
    return 0;
}