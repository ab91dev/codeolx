.PHONY: build run

build:
	@go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api

run: build
	@./bin/api
