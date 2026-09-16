bin:
	@go build -o bin/api ./cmd/api/main.go

run: bin
	@./bin/api
