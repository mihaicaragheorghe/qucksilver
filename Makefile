run:
	@go run cmd/main.go

build:
	@go build -o bin/ cmd/main.go

test:
	@go test ./...

lint:
	@golangci-lint run