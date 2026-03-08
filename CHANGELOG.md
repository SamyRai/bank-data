# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Unified `pkg/financial` API as canonical entrypoint for detection, validation, and parsing.
- Internal modular validators for `IBAN`, `BIC`, `SEPA_CREDITOR_ID`, `LEI`, `ISIN`, `PAN`, and initial `VAT` country support.
- Deterministic bank registry abstraction in `internal/bankregistry`.
- OSS publishing essentials: `LICENSE`, `CONTRIBUTING.md`, root `SECURITY.md`.

### Changed
- Legacy `pkg/iban`, `pkg/bic`, and `pkg/sepa` packages moved to compatibility wrappers.
- CI workflows aligned to `main` pull-request target and `go.mod` toolchain.
- Registry generation pipeline normalized on `datasets/iban-registry.txt` naming.

### Security
- Private vulnerability reporting path documented via GitHub security advisories.
