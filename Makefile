.PHONY: build run

build:
	@go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api

run: build
	@./bin/api

build-migrate:
	@go build -o bin/migrate ./cmd/migrate

migrate-up: build-migrate
	@./bin/migrate up

migrate-down: build-migrate
	@./bin/migrate down