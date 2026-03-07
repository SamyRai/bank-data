# Development Guide

This guide outlines the workflows for contributing to `bank-data`, running tests, and managing the metadata registry.

## Getting Started

### Prerequisites
- Go 1.24+
- `make` (optional)

### Running Tests
We maintain 100% logic coverage using a mix of unit, table-driven, and fuzz tests.

```sh
# Run all tests
go test ./...

# Run IBAN fuzzing
go test -fuzz=FuzzIBANValidate ./pkg/iban
```

---

## Metadata Registry Management

The country metadata in `internal/countrymeta/` is the core of the library's validation logic. It is managed via automated generation to ensure accuracy.

### Updating Country Rules
1. **Modify Specs**: Update the raw definitions in `internal/countrymeta/registry.go` (or the underlying data source if fully automated).
2. **Generate**: Run the generation tool to synchronize metadata.
   ```sh
   go generate ./...
   ```
3. **Validate**: Always run `go test ./pkg/iban` after updates to ensure no regressions in country-specific parsing logic.

---

## Contribution Guidelines

1. **Dependency Discipline**: Core logic must remain dependency-free. Only use the Go standard library.
2. **Interface Driven**: Define public capabilities as interfaces in `pkg/`. While core IBAN logic lives in `pkg/iban` for accessibility, other specialized engines should reside in `internal/`.
3. **Structured Logging**: Use the internal `log` package. Avoid `fmt.Printf` for anything other than CLI tools or debugging.
4. **Error Handling**: Use the typed errors defined in `pkg/iban`. Avoid wrapping errors unless necessary for context.
