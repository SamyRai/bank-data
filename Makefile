# Makefile for bank-data

.PHONY: all build lint test fmt

all: build

build:
	go build ./...

lint:
	golangci-lint run --timeout=5m

test:
	go test ./...

fmt:
	gofmt -s -w .
