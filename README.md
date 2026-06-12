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
    internal "github.com/SamyRai/bank-data/internal/validationformats/iban"
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

- **Expanded Country Registry**: Supports all major EU/SEPA and a broad set of IBAN countries via `countrymeta.Registry`.
- **Robust IBAN Validation**: Streaming MOD-97 algorithm, format, length, and checksum validation for all supported countries.
- **Comprehensive Testing**: Table-driven, fuzz, and edge case tests for reliability and coverage.
- **Unified Public API**: Single facade (`iban.Service`) with dependency injection for validation, parsing, and enrichment.
- **Typed Error Handling**: Domain-specific, exported error types for all validation and parsing outcomes.
- **Security & Compliance**: Input capped at 34 characters, constant-time check digit comparison, and thorough edge case coverage.
- **Bank Directory Enrichment**: IBAN parsing extracts bank code and supports enrichment (e.g., DE BLZ→BIC lookup).
- **Documentation & CI**: Production-ready documentation, MIT license, and multi-platform/multi-version CI workflows.
- **Refactored Abstractions**: Removed redundant IBANRegistrySource in favor of general encoder/decoder for maintainability.

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

## Contribution Guidelines

- Follow SOLID and clean architecture principles.
- Use dependency injection and keep cyclomatic complexity low (≤ 10 per function).
- Group code by feature, then by type.
- Add comprehensive unit tests for all new code.
- Use TODO comments for improvements, with context, priority, and estimated effort.
- Update TODO.md after each milestone.
- Validate all changes with tests and CI before submitting PRs.

## License

This project is licensed under the [MIT License](LICENSE).

## Architecture & Design Principles

- **Clean Architecture**: Separation of concerns, interface-driven, and feature-first organization.
- **SOLID Principles**: Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion.
- **Dependency Injection**: All core services accept dependencies via constructors for loose coupling and testability.
- **Error Handling**: Domain-specific error types for validation, parsing, and enrichment. All errors are typed and documented.
- **Logging & Monitoring**: Structured logging for all validation and parsing events. Extendable for external monitoring.
