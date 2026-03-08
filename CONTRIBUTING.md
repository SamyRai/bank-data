# Contributing

Thanks for contributing to `bank-data`.

## Prerequisites

- Go 1.26+
- A clean working tree before submitting a PR

## Development Workflow

```sh
# Run full test suite
go test ./...

# Run vet
go vet ./...

# Run race checks on core public packages
go test -race ./pkg/financial ./pkg/iban ./pkg/bic ./pkg/sepa ./pkg/validation
```

## Architecture Rules

- New features should target `pkg/financial` and `internal/identifiers/*`.
- Keep one identifier per module with clear `Normalize` / `DetectCandidate` / `Validate` / `Parse` ownership.
- Do not expose `internal/*` package types in public package signatures.
- Prefer small, behavior-focused commits with updated tests.

## Pull Requests

- Include tests for any behavior change.
- Update docs when public behavior or API changes.
- Keep PR descriptions explicit about breaking changes and migration impact.
