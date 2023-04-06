.PHONY: all test fmt pre-commit help

all: fmt

fmt:
	@echo Formatting...
	@go fmt ./...
	@go vet -composites=false ./...

test:
	@go test -v ./...

pre-commit: fmt
	@go mod tidy

