# bank-data

A Go package for robust IBAN validation, parsing, and structure detection. Designed with clean architecture, SOLID principles, and extensibility in mind for future banking and financial data types.

[![Go Reference](https://pkg.go.dev/badge/github.com/SamyRai/bank-data.svg)](https://pkg.go.dev/github.com/SamyRai/bank-data)
[![Test & Coverage](https://github.com/SamyRai/bank-data/actions/workflows/test.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/test.yml)
[![Lint](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/SamyRai/bank-data)](https://goreportcard.com/report/github.com/SamyRai/bank-data)

## Installation

```sh
go get github.com/SamyRai/bank-data
```

## Quick Start

### IBAN Validation

```go
import (
    "fmt"
    "log"

    "github.com/SamyRai/bank-data/pkg/iban"
    "github.com/SamyRai/bank-data/internal/iban"
)

func main() {
    // Create a new IBAN service
    service := iban.NewService(
        internal.NewValidator(),
        internal.NewParser(),
        internal.NewDetector(),
        nil, // or provide a custom *validation.ValidationRegistry
    )

    // Validate an IBAN
    err := service.Validate("DE89370400440532013000")
    if err != nil {
        log.Fatalf("IBAN validation failed: %v", err)
    }
    fmt.Println("IBAN is valid!")
}
```

### BIC Validation

```go
import (
    "fmt"
    "log"

    "github.com/SamyRai/bank-data/pkg/bic"
)

func main() {
    // Validate a BIC
    err := bic.Validate("DEUTDEFF")
    if err != nil {
        log.Fatalf("BIC validation failed: %v", err)
    }
    fmt.Println("BIC is valid!")
}
```

### IBAN Parsing and Detection

```go
import (
    "fmt"
    "log"

    "github.com/SamyRai/bank-data/pkg/iban"
    "github.com/SamyRai/bank-data/internal/iban"
)

func main() {
    // Create a new IBAN service
    service := iban.NewService(
        internal.NewValidator(),
        internal.NewParser(),
        internal.NewDetector(),
        nil,
    )

    // Parse an IBAN
    info, err := service.Parse("DE89370400440532013000")
    if err != nil {
        log.Fatalf("IBAN parsing failed: %v", err)
    }
    fmt.Printf("Parsed IBAN: %+v\n", info)

    // Detect the structure of an IBAN
    structure, err := service.Detect("DE89370400440532013000")
    if err != nil {
        log.Fatalf("IBAN detection failed: %v", err)
    }
    fmt.Printf("IBAN Structure: %+v\n", structure)
}
```

## Features

- **IBAN Validation**: Format and checksum validation for a wide range of countries.
- **BIC Validation**: Format and country code validation for SWIFT Business Identifier Codes.
- **IBAN Parsing**: Extracts country code, bank code, account number, and more.
- **IBAN Structure Detection**: Identifies country, length, and structure of IBANs.
- **Typed Error Handling**: Domain-specific, typed errors for all validation outcomes.
- **Logging**: Structured logging for validation and parsing events.
- **Batch and Fuzz Testing**: High test coverage, including edge cases and fuzzing.
- **Cross-Field and Conditional Validation**: Validate relationships between multiple fields (e.g., account and country) and support conditional rules. Easily register cross-field validators in the ValidationRegistry for use in the service layer or for custom business logic.

## API Documentation

The public API of this library is exposed through the `pkg` directory.

- **`pkg/iban`**: Provides the main service for IBAN validation, parsing, and detection. See the [Go documentation](https://pkg.go.dev/github.com/SamyRai/bank-data/pkg/iban) for more details.
- **`pkg/bic`**: Provides functions for BIC validation. See the [Go documentation](https://pkg.go.dev/github.com/SamyRai/bank-data/pkg/bic) for more details.

## Error Handling

The library returns typed errors, allowing you to handle different error conditions programmatically.

```go
import (
    "errors"
    "fmt"
    "log"

    "github.com/SamyRai/bank-data/pkg/iban"
    "github.com/SamyRai/bank-data/internal/iban"
)

func main() {
    service := iban.NewService(
        internal.NewValidator(),
        internal.NewParser(),
        internal.NewDetector(),
        nil,
    )

    err := service.Validate("invalid-iban")
    if err != nil {
        var ibanErr *iban.IBANError
        if errors.As(err, &ibanErr) {
            fmt.Printf("IBAN validation failed with code: %s, message: %s\n", ibanErr.Code, ibanErr.Message)
        } else {
            log.Fatalf("An unexpected error occurred: %v", err)
        }
    }
}
```

## Project Structure

- `pkg/`: Public API and domain types for `iban` and `bic`.
- `internal/`: Internal implementation of the business logic.
- `gen/`: Code generation tools for automating registry updates.
- `datasets/`: Data files used for generation and testing.
- `testdata/`: Test vectors for IBAN validation.

## Testing

Run all tests (including fuzz and table-driven tests):

```sh
make test
```

### Code Generation

This project uses `go:generate` to automate the creation of the IBAN country registry. To run the generator, use the following command:

```sh
go generate ./...
```

## Roadmap & TODOs

See [FUTURE_FEATURES.md](FUTURE_FEATURES.md) and [TODO.md](TODO.md) for planned features, progress tracking, and priorities.

## Contributing

Contributions are welcome! Please see the documentation and TODOs for guidance on project structure and priorities.

## License

[MIT License](LICENSE) (add this file if not present)
