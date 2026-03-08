# Architecture Overview

`bank-data` now follows a single public-facade model: all identifier workflows go through `pkg/financial`.

## High-Level Layout

### 1. Public API
- `pkg/financial`: canonical API for `Detect`, `Validate`, `Parse`, batch, and stream operations.
- `pkg/iban`, `pkg/bic`, `pkg/sepa`: compatibility wrappers for legacy integrations.

### 2. Identifier Modules
- `internal/identifiers/{iban,bic,sepa,lei,isin,pan,vat}`
- Each module owns four responsibilities only:
  - normalize input
  - detect candidate shape
  - validate semantics/checksum
  - parse structured fields

### 3. Shared Core Runtime
- `internal/core/validation`: generic concurrent engine used by public facades.
- `internal/countrymeta`: generated IBAN country metadata.
- `internal/bankregistry`: explicit bank enrichment registry with source metadata.

## Design Principles

- Single ownership per module (SRP): one identifier module, one reason to change.
- Typed boundaries: no `map[string]any`-style public validator registry.
- Deterministic behavior: no implicit cwd-based dataset loading in core validation paths.
- Compatibility by adapter: old package entrypoints wrap current internals rather than owning new logic.
