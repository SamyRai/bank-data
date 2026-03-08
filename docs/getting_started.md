# Getting Started with `bank-data`

Welcome to `bank-data`! This guide will help you install and use the library. `bank-data` is a robust Go package designed to provide high-performance validation, parsing, and detection for common financial identifiers like IBAN, BIC, and SEPA Creditor IDs.

The recommended approach is to use the canonical `pkg/financial` package, which unifies all features into a stable, `v1` API contract.

## Installation

Run the following command to download and add the package to your project:

```sh
go get github.com/SamyRai/bank-data
```

## Quick Start Example

This example demonstrates how to use the unified `financial.Service` to validate and parse a mixed set of identifiers.

```go
package main

import (
    "fmt"
    "github.com/SamyRai/bank-data/pkg/financial"
)

func main() {
    // 1. Initialize the Service
    svc := financial.NewService()

    // 2. Sample Data
    identifiers := []string{
        "DE89370400440532013000", // Valid IBAN
        "DEUTDEFFXXX",            // Valid BIC
        "DE98ZZZ09999999999",     // Valid SEPA Creditor ID
        "invalid_data_123",       // Invalid string
    }

    // 3. Process Identifiers
    for _, id := range identifiers {
        report := svc.Validate(id)
        if report.Valid {
            fmt.Printf("✅ Valid %s: %s\n", report.IdentifierType, report.Normalized)

            // Attempt to parse its components
            parsed, err := svc.Parse(id)
            if err == nil {
                fmt.Printf("   Components: %v\n", parsed.Components)
            }
        } else {
            fmt.Printf("❌ Invalid: %s (Type: %s, Error: %v)\n", report.Normalized, report.IdentifierType, report.Error)
        }
    }
}
```

## Moving from Legacy Packages

If you previously used `pkg/iban`, `pkg/bic`, or `pkg/sepa`, we strongly recommend updating your codebase to use `pkg/financial`.

*   The legacy packages expose multiple unique interfaces which can be difficult to manage.
*   The `v1` API Contract in `pkg/financial` is explicitly versioned and provides strong compatibility guarantees. Legacy packages are in a deprecation window and will eventually be removed.
*   For details on migrating your specific function calls, please see the [Migration Matrix in the API Reference](./api_reference.md).

For advanced usage or internal details, check out the [Architecture Document](./architecture.md).
