default: build

build:
	@go mod tidy
	go build -o bin/test -v cmd/main.go