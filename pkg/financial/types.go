package financial

// IdentifierType represents the type of financial identifier detected or parsed.
type IdentifierType string

const (
	// TypeUnknown indicates the identifier type could not be determined.
	TypeUnknown IdentifierType = "unknown"
	// TypeIBAN indicates an International Bank Account Number.
	TypeIBAN IdentifierType = "iban"
	// TypeBIC indicates a Bank Identifier Code (SWIFT).
	TypeBIC IdentifierType = "bic"
	// TypeSEPACreditorID indicates a SEPA Creditor Identifier.
	TypeSEPACreditorID IdentifierType = "sepa_creditor_id"
)

// ValidationReport provides a standardized result across all identifier types.
type ValidationReport struct {
	Valid          bool           // True if the identifier is fully valid
	IdentifierType IdentifierType // The detected type of the identifier
	Normalized     string         // The normalized (cleaned) version of the input
	Error          error          // The validation error, if any
}

// ParsedIdentifier provides a canonical representation of a parsed identifier.
type ParsedIdentifier struct {
	IdentifierType IdentifierType    // The detected type of the identifier
	Normalized     string            // The normalized string
	Components     map[string]string // Key-value pairs of extracted components (e.g., "CountryCode", "BankCode")
}

// Suggestion provides a possible valid identifier if the input was slightly malformed (e.g. whitespace, case).
type Suggestion struct {
	Suggested string // The suggested valid identifier
	Reason    string // Why the suggestion was made (e.g. "Removed whitespace")
}

// Service is the canonical v1 facade for all financial identifier operations.
type Service interface {
	// Validate checks an identifier against all supported types and returns a unified report.
	Validate(input string) ValidationReport

	// Parse deconstructs an identifier into its canonical components.
	Parse(input string) (*ParsedIdentifier, error)

	// Detect identifies the type of the given financial string.
	Detect(input string) IdentifierType
}
