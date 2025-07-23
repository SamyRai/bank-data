# bank-data

A production-ready Go package for robust IBAN validation, parsing, and structure detection. Built with clean architecture, SOLID principles, and extensibility for future banking and financial data types.

[![Go Reference](https://pkg.go.dev/badge/github.com/SamyRai/bank-data.svg)](https://pkg.go.dev/github.com/SamyRai/bank-data)
[![Test & Coverage](https://github.com/SamyRai/bank-data/actions/workflows/test.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/test.yml)
[![Lint](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml)

## Overview

bank-data provides:

- Secure, standards-compliant IBAN validation (format, length, MOD-97 checksum)
- Country-specific parsing and enrichment (e.g., DE BLZ→BIC)
- Typed error handling and structured logging
- High test coverage (unit, fuzz, edge cases)
- Extensible design for future banking data types

## Installation

```sh
go get github.com/SamyRai/bank-data
```

## Quick Start

Basic validation:

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

Advanced usage (parsing, enrichment):

```go
info, err := service.Parse("DK5000400440116243") // Denmark example
if err == nil {
    fmt.Println("Bank code:", info.BankCode)
    // Enrich with BIC lookup, etc.
}
```

## Features

- **IBAN Validation**: Format, length, and checksum for all major countries (including Denmark)
- **IBAN Parsing**: Extracts country code, bank code, account number, and check digits
- **Bank Directory Lookup**: Enriches parsed IBANs with BIC and bank name (DE, extensible)
- **Typed Error Handling**: Domain-specific errors for all validation outcomes
- **Structured Logging**: For validation and parsing events
- **Batch & Fuzz Testing**: High coverage, edge cases, and fuzzing
- **Extensible Validation**: Register custom validators for new data types

## Supported Countries

- All major EU/SEPA countries (DE, DK, FR, IT, ES, NL, BE, AT, CH, SE, etc.)
- Denmark (DK) is fully supported and tested
- Easily extensible via countrymeta.Registry

## Project Structure

- `pkg/iban`: Public interfaces, types, and main service API
- `internal/iban`: Core validation, parsing, and detection logic
- `internal/validation`: Shared validation types and registry
- `testdata/`: Test vectors for IBAN validation

## Testing

Run all tests (unit, fuzz, edge cases):

```sh
go test ./...
```

## CI & Security

- Automated tests, linting, security scans (GoSec, CodeQL, dependency review)
- Multi-platform, multi-version CI via GitHub Actions
- See `.github/workflows/` for details

## Roadmap & TODOs

See [FUTURE_FEATURES.md](FUTURE_FEATURES.md) and [TODO.md](TODO.md) for planned features and progress tracking.

## Contributing

Contributions are welcome! See documentation and TODOs for guidance on project structure and priorities.

## License

This project is licensed under the [MIT License](LICENSE).
