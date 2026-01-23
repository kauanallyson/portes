CC = gcc
CFLAGS = -Wall -Wextra -Werror -g -std=c11
SRC = ./src/portes.c
OUT = ./bin/portes

default: setup
	$(CC) $(SRC) $(CFLAGS) -o $(OUT)

setup:
	mkdir -p bin

clean:
	rm -rf bin