# Makefile for bank-data

.PHONY: all build lint test fmt check-docs check-coverage check-bench

all: build

build:
	go build ./...

lint:
	golangci-lint run --timeout=5m

test:
	go test ./...

fmt:
	gofmt -s -w .

check-docs:
	./scripts/check_repo_docs.sh

check-coverage:
	./scripts/check_coverage.sh

check-bench:
	./scripts/check_bench_regression.sh
