# Development Guide

## Prerequisites

- Go 1.26+
- `make` (optional)

## Standard Checks

```sh
# Unit tests
go test ./...

# Vet
go vet ./...

# Race checks for core packages
go test -race ./pkg/financial ./pkg/iban ./pkg/bic ./pkg/sepa ./pkg/validation

# Fuzz smoke
go test -run=^$ -fuzz=FuzzIBANValidate -fuzztime=3s ./pkg/iban
```

## Metadata Registry Workflow

Country metadata is generated from `datasets/iban-registry.txt`.

```sh
# regenerate internal/countrymeta/registry.go
go generate ./...

# run focused IBAN checks after regeneration
go test ./pkg/iban
```

## Contribution Rules

1. New functionality should target `pkg/financial` and `internal/identifiers/*`.
2. Keep module ownership clear (normalize/detect/validate/parse per identifier).
3. Avoid introducing `internal/*` types in public package signatures.
4. Keep docs aligned with actual code behavior and CI checks.
