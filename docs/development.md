# Development Guide

## Prerequisites

- Go 1.26+
- `make` (optional)

## Standard Checks

```sh
# Docs/release artifact checks
./scripts/check_repo_docs.sh

# Unit tests
go test ./...

# Vet
go vet ./...

# Race checks for core packages
go test -race ./pkg/financial ./pkg/iban ./pkg/bic ./pkg/sepa ./pkg/nationalaccount ./pkg/vop ./pkg/iso20022 ./pkg/validation

# Fuzz smoke
go test -run=^$ -fuzz=FuzzIBANValidate -fuzztime=3s ./pkg/iban
go test -run=^$ -fuzz=FuzzService_DetectValidateParse -fuzztime=2s ./pkg/financial
go test -run=^$ -fuzz=FuzzBICValidateAndParse -fuzztime=2s ./pkg/bic
go test -run=^$ -fuzz=FuzzSEPACreditorValidate -fuzztime=2s ./pkg/sepa
go test -run=^$ -fuzz=FuzzValidateAndParse -fuzztime=2s ./pkg/nationalaccount
go test -run=^$ -fuzz=FuzzMatcherMatch -fuzztime=2s ./pkg/vop
go test -run=^$ -fuzz=FuzzPain001ParseValidate -fuzztime=2s ./pkg/iso20022

# Coverage + benchmark gates
go test -coverprofile=coverage.out ./... && ./scripts/check_coverage.sh
./scripts/check_bench_regression.sh
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
