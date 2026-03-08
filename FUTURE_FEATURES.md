# Future Expansion: Financial Data Types

The library's registry-based architecture is designed to support a wide range of financial data types beyond IBANs.

## Planned Validations

| Data Type | Validation Strategy |
| :--- | :--- |
| **BIC/SWIFT** | Regex, length verification, and country/bank lookup. |
| **SEPA Creditor ID** | Format validation including MOD-97 check digits. |
| **National Accounts** | Country-specific rules (e.g., UK Sort Codes, DE BLZ). |
| **VAT Number** | ISO country prefix and national checksum logic. |
| **LEI** | MOD-97-10 (ISO 17442) validation for legal entities. |
| **ISIN** | Luhn checksum and ISO 6166 format validation. |
| **Card (PAN)** | Luhn algorithm and IIN/BIN prefix detection. |

## Current Status (2026-03-08)

- `pkg/financial` is now the canonical public API for detection, validation, and parsing.
- Implemented identifier modules: **IBAN**, **BIC**, **SEPA Creditor ID**, **LEI**, **ISIN**, **PAN**.
- VAT support is implemented for initial country subset: **DE**, **FR**, **NL**, **IT**, **ES**.

---

## Feature Roadmap

### 1. Bank Information Lookups
Extend the current IBAN parser to optionally fetch bank metadata (name, address, BIC) based on the national bank code embedded in the IBAN.

### 2. High-Performance Batching
Provide a dedicated batch processing engine for validating thousands of records with minimal latency, utilizing Go 1.26 concurrency primitives.
See [Batch Processing Architecture](docs/batch_processing.md) for details.

### 3. Localization Support
Dynamic rule updates based on local clearing house requirements (e.g., specific rules for SEPA Instant vs. standard SEPA).
