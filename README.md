# bank-data

A Go (Golang) package for financial and banking data processing, starting with IBAN validation, parsing, and detection.

[![Go Reference](https://pkg.go.dev/badge/github.com/SamyRai/bank-data.svg)](https://pkg.go.dev/github.com/SamyRai/bank-data)
[![Test & Coverage](https://github.com/SamyRai/bank-data/actions/workflows/test.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/test.yml)
[![Lint](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml)

## Features

- IBAN validation (format and checksum)
- IBAN parsing (extract country, bank code, account number, etc.)
- IBAN structure detection (country, length, structure)

## Structure

- `pkg/iban`: Public interfaces and types
- `internal/iban`: Implementations and tests

## Usage

Import the interfaces from `pkg/iban` and use the implementations from `internal/iban`.

## Testing

```sh
go test ./...
```

## TODO

See [TODO.md](TODO.md)
