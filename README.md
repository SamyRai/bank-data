# bank-data

A Go (Golang) clean architecture package for financial and banking data processing, starting with IBAN validation, parsing, and detection.

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

- [ ] Add support for more IBAN countries and structures (high, 2d)
- [ ] Improve error types and messages (medium, 1d)
- [ ] Add logging and monitoring (medium, 1d)
- [ ] Optimize performance for large-scale validation (low, 2d)
