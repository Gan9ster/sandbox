BINARY=hsrctl

.PHONY: all build test up

all: build

build:
go build -o bin/$(BINARY) ./cmd/hsrctl

test:
go test ./...

up: build
@echo "demo environment not implemented"
