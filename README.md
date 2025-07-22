# bank-data

A Go package for robust IBAN validation, parsing, and structure detection. Designed with clean architecture, SOLID principles, and extensibility in mind for future banking and financial data types.

[![Go Reference](https://pkg.go.dev/badge/github.com/SamyRai/bank-data.svg)](https://pkg.go.dev/github.com/SamyRai/bank-data)
[![Test & Coverage](https://github.com/SamyRai/bank-data/actions/workflows/test.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/test.yml)
[![Lint](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml)

## Installation

```sh
go get github.com/SamyRai/bank-data
```

## Quick Start

```go
import (
    "github.com/SamyRai/bank-data/pkg/iban"
    internal "github.com/SamyRai/bank-data/internal/iban"
)

service := iban.NewService(
    internal.NewValidator(),
    internal.NewParser(),
    internal.NewDetector(),
    nil, // or provide a custom *validation.ValidationRegistry
)

err := service.Validate("DE89370400440532013000")
if err != nil {
    // handle invalid IBAN
}
```

## Features

- **IBAN Validation**: Format and checksum validation for a wide range of countries.
- **IBAN Parsing**: Extracts country code, bank code, account number, and more.
- **IBAN Structure Detection**: Identifies country, length, and structure of IBANs.
- **Typed Error Handling**: Domain-specific, typed errors for all validation outcomes.
- **Logging**: Structured logging for validation and parsing events.
- **Batch and Fuzz Testing**: High test coverage, including edge cases and fuzzing.
- **Cross-Field and Conditional Validation**: Validate relationships between multiple fields (e.g., account and country) and support conditional rules. Easily register cross-field validators in the ValidationRegistry for use in the service layer or for custom business logic.

## Project Structure

- `pkg/iban`: Public interfaces, types, and the main service API (IBAN-specific).
- `internal/iban`: Implementations of IBAN validation, parsing, and detection.
- `internal/validation`: Shared validation types and registry (for future extensibility).
- `testdata/`: Test vectors for IBAN validation.

## Testing

Run all tests (including fuzz and table-driven tests):

```sh
go test ./...
```

## Roadmap & TODOs

See [FUTURE_FEATURES.md](FUTURE_FEATURES.md) and [TODO.md](TODO.md) for planned features, progress tracking, and priorities.

## Contributing

Contributions are welcome! Please see the documentation and TODOs for guidance on project structure and priorities.

## License

[MIT License](LICENSE) (add this file if not present)
