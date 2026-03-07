# Architecture Overview

The `bank-data` library is designed for high-integrity financial data processing, adhering to clean architecture and modular design.

## Design Philosophy

- **Clean Architecture**: Separation between public interfaces/implementations (`pkg/`) and internal meta-engines (`internal/`).
- **Extensibility**: Uses Registry and Strategy patterns to support new financial data types (BIC, VAT, etc.) without breaking changes.
- **Precision**: Algorithmic accuracy over simplicity, using streaming MOD-97 for IBAN validation to avoid numeric overflows.

---

## Component Breakdown

### 1. Public API & Core (`pkg/`)
- **`pkg/iban`**: Primary entry point. Hosts the `Service` facade, core validation logic (including the streaming checksum), and the BBAN parser.
- **`pkg/validation`**: Generic interfaces and the `ValidationRegistry` for type-agnostic processing.
- **`pkg/bank`**: Shared types for bank and branch information.

### 2. Implementation Layer (`internal/`)
- **`internal/countrymeta`**: The metadata engine, containing format rules for 80+ countries.
- **`internal/bic/map`**: Logic for mapping national bank codes to global BICs.
- **`internal/log`**: A dependency-free, structured logger used throughout the library.

### 3. Tooling and Generation (`gen/`)
- **`gen/cmd/gen_registry`**: Automated tool that transforms dataset specifications into Go metadata, ensuring the registry remains up-to-date with minimal manual effort.

---

## Technical Patterns

### Registry Pattern
The `ValidationRegistry` decouples the library from specific data types. Clients can register custom validators, allowing the core service to handle diverse financial formats (IBAN, BIC, VAT, etc.) through a unified interface.

### Strategy Pattern
Operations like validation and parsing are defined via interfaces (e.g., `iban.Validator`). This enables easy swapping of implementations—for example, switching from a single IBAN validator to a high-concurrency batch validator—without affecting client code.

### Service Facade
The `iban.Service` simplifies complex workflows by orchestrating multiple internal components (parsers, registries, detectors) into a single, cohesive API.
