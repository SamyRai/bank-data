# bank-data

A high-performance Go library for IBAN validation, parsing, and structure detection. Built for financial systems that require precision, reliability, and extensibility.

[![Go Reference](https://pkg.go.dev/badge/github.com/SamyRai/bank-data.svg)](https://pkg.go.dev/github.com/SamyRai/bank-data)
[![Test & Coverage](https://github.com/SamyRai/bank-data/actions/workflows/test.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/test.yml)
[![Lint](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml/badge.svg)](https://github.com/SamyRai/bank-data/actions/workflows/lint.yml)

## Installation

```sh
go get github.com/SamyRai/bank-data
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/SamyRai/bank-data/pkg/iban"
)

func main() {
	// Initialize the service with default components
	svc := iban.NewService(nil, nil, nil, nil)

	// Validate an IBAN
	input := "DE89370400440532013000"
	if err := svc.Validate(input); err != nil {
		fmt.Printf("Invalid IBAN: %v\n", err)
		return
	}

	// Parse IBAN components
	info, _ := svc.Parse(input)
	fmt.Printf("Country: %s, Bank Code: %s\n", info.CountryCode, info.BankCode)
}
```

## Core Features

- **Robust Validation**: Streaming MOD-97 algorithm prevents overflows and ensures constant-time performance.
- **Deep Parsing**: Extracts Country Code, Bank Code, and Account Numbers with precision.
- **Global Support**: Pre-configured with metadata for 80+ countries, easily extensible via registry.
- **Developer First**: Structured logging, typed errors, and zero external dependencies in the core.
- **Security Focused**: Input length capping (34 chars) and strict character set enforcement.

## Project Structure

- `pkg/financial`: The canonical v1 Public API for interacting with financial formats.
- `pkg/iban`, `pkg/bic`, `pkg/sepa`: Legacy compatibility packages (deprecated).
- `internal/`: Component implementations and metadata registry.
- `gen/`: Tooling for automated registry maintenance.

## Documentation

- [Getting Started](docs/getting_started.md)
- [API Reference (v1 API Contract & Migration)](docs/api_reference.md)
- [Architecture Overview](docs/architecture.md)
- [Development Guide](docs/development.md)

## Contributing

Contributions are welcome. See [TODO.md](TODO.md) and [FUTURE_FEATURES.md](FUTURE_FEATURES.md) for current priorities and upcoming developments.

## License

[MIT License](LICENSE)
