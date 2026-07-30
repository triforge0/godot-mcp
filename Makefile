.PHONY: build run doctor version test

BINARY := bin/godot-mcp

build:
	go build -o $(BINARY) ./cmd/godot-mcp

run: build
	./$(BINARY) start

doctor: build
	./$(BINARY) doctor

version: build
	./$(BINARY) version

test:
	go test ./...
