GOFMT ?= gofmt -s -w
GOFILES := $(shell find . -name "*.go" -type f)

.PHONY: default
default: test

.PHONY: ci
ci: 
	@go test  -v -race ./...

.PHONY: test
test: 
	@go test  -v -race ./...

.PHONY: fmt
fmt:
	@echo source formatting...
	@$(GOFMT) $(GOFILES)
