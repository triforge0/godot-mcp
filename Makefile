.PHONY: build run doctor version test test-integration

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

test-integration:
	go test -tags=integration -timeout 10m ./internal/integration/...
