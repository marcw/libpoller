all: test

test:
	@go version
	@go build -race ./...
	@go test -v ./...
	@go vet ./...