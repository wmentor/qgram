.PHONY: all

all: gofmt lint test

lint:
	@if [ -x "$$(command -v golangci-lint)" ]; then echo "run linter..." ; golangci-lint run ; else echo "golangci-lint not found (skipped)" ; fi

test:
	@echo "run test..."
	@go clean -testcache
	@go test ./... -cover

gofmt:
	@if [ -x "$$(command -v gofmt)" ]; then echo "run gofmt..." ; gofmt -s -w . ; else echo "gofmt not found (skipped)" ; fi


