build:
	@go build -o bin/email-verification-tool

run: build
	@./bin/email-verification-tool

test:
	@go test -v ./...