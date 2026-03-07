# API Reference: IBAN Package

The `iban` package provides a unified interface for validating, parsing, and detecting International Bank Account Numbers.

## The IBAN Service

The `Service` is the central entry point. For optimal performance, initialize it once and reuse it across your application.

```go
import "github.com/SamyRai/bank-data/pkg/iban"

// Initialize with default components by passing nil
svc := iban.NewService(nil, nil, nil, nil)
```

### Methods

#### `Validate(ibanStr string) error`
Validates an IBAN string for character set, length, and MOD-97 checksum. Returns a typed `*iban.IBANError` if invalid.

#### `Parse(ibanStr string) (*IBANInfo, error)`
Deconstructs a normalized IBAN into its constituent fields (Country Code, Bank Code, Account Number, etc.).

#### `Detect(ibanStr string) (*IBANStructure, error)`
Retrieves metadata for a specific country's IBAN format, including the expected length and structure template.

---

## Errors and Status Codes

The package uses a structured `IBANError` type for precise error handling.

| Error Code | Constant | Description |
| :--- | :--- | :--- |
| `invalid_chars` | `ErrCodeInvalidChars` | Input contains characters outside of [A-Z0-9]. |
| `wrong_length` | `ErrCodeWrongLength` | Length does not match country meta or limits. |
| `checksum_failed` | `ErrCodeChecksum` | The MOD-97 checksum validation failed. |
| `unsupported_country` | `ErrCodeUnsupportedCountry` | The ISO country code is not in the registry. |

### Idiomatic Error Handling
```go
err := svc.Validate(input)
if err != nil {
	var ibanErr *iban.IBANError
	if errors.As(err, &ibanErr) {
		fmt.Printf("Validation failed: %s (Code: %s)\n", ibanErr.Message, ibanErr.Code)
	}
}
```

---

## Advanced Usage

### The `Must` Wrapper
Use the `Must()` wrapper in environments where an invalid IBAN should trigger a panic (e.g., CLI tools or during migration initialization).

```go
// Panics if the IBAN is invalid
svc.Must().Validate("DE89370400440532013000")

// Returns IBANInfo directly, panics on error
info := svc.Must().Parse("DE89370400440532013000")
```

### Bank Information Enrichment
Combine `IBANInfo` with a `BankBICMap` to retrieve detailed bank metadata.

```go
// Enrich parsed info with bank/branch details
bankInfo, err := iban.EnrichWithBankInfo(info, bicMap)
```
