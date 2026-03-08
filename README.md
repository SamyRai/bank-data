# bank-data

Go library for validating and parsing financial identifiers with a canonical `pkg/financial` API.

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

	"github.com/SamyRai/bank-data/pkg/financial"
)

func main() {
	svc := financial.NewService()

	report, err := svc.Validate("DE89370400440532013000", financial.IdentifierIBAN)
	if err != nil {
		fmt.Printf("invalid: %v\n", err)
		return
	}
	fmt.Printf("type=%s normalized=%s valid=%v\n", report.Type, report.Normalized, report.Valid)

	parsed, _ := svc.Parse("US0378331005", financial.IdentifierISIN)
	fmt.Printf("ISIN country=%s nsin=%s\n", parsed.Fields["country_code"], parsed.Fields["nsin"])
}
```

## Runnable Identifier Examples

```go
svc := financial.NewService()

svc.Validate("DE89370400440532013000", financial.IdentifierIBAN)
svc.Validate("DEUTDEFF", financial.IdentifierBIC)
svc.Validate("DE98ZZZ09999999999", financial.IdentifierSEPACreditor)
svc.Validate("529900T8BM49AURSDO78", financial.IdentifierLEI)
svc.Validate("US0378331005", financial.IdentifierISIN)
svc.Validate("4111111111111111", financial.IdentifierPAN)
svc.Validate("DE136695976", financial.IdentifierVAT)
svc.Validate("20-00-00 55779911", financial.IdentifierNationalAccountUK)
```

See executable examples in [`pkg/financial/examples_test.go`](pkg/financial/examples_test.go).

## Supported Identifier Types

- `IBAN`
- `BIC`
- `SEPA_CREDITOR_ID`
- `LEI`
- `ISIN`
- `PAN`
- `VAT` (initial checks for `DE`, `FR`, `NL`, `IT`, `ES`)
- `NATIONAL_ACCOUNT_UK` (sort code + account number pair)

## Project Structure

- `pkg/financial`: Canonical public API (`Detect`, `Validate`, `Parse`, batch/stream).
- `internal/identifiers/*`: Per-identifier normalization, detection, validation, parsing.
- `internal/core/validation`: Shared concurrent validation engine.
- `internal/bankregistry`: Explicit, deterministic bank enrichment registry.
- `pkg/iban`, `pkg/bic`, `pkg/sepa`: Compatibility packages (legacy entrypoints).

## Documentation

- [Architecture Overview](docs/architecture.md)
- [API Reference](docs/api_reference.md) (includes v1 Compatibility & Deprecation Policy)
- [Development Guide](docs/development.md)
- [Getting Started](docs/getting_started.md)
- [Strategic Development Program](docs/development_program/README.md)

## Contributing

Contributions are welcome. See [TODO.md](TODO.md) and [FUTURE_FEATURES.md](FUTURE_FEATURES.md) for open priorities.

## License

[MIT License](LICENSE)
