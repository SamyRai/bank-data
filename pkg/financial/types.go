package financial

import "fmt"

// IdentifierType represents a supported financial identifier domain.
type IdentifierType string

const (
	IdentifierIBAN         IdentifierType = "IBAN"
	IdentifierBIC          IdentifierType = "BIC"
	IdentifierSEPACreditor IdentifierType = "SEPA_CREDITOR_ID"
	IdentifierLEI          IdentifierType = "LEI"
	IdentifierISIN         IdentifierType = "ISIN"
	IdentifierPAN          IdentifierType = "PAN"
	IdentifierVAT          IdentifierType = "VAT"
	IdentifierNationalAccountUK IdentifierType = "NATIONAL_ACCOUNT_UK"
)

// ValidationError contains machine-friendly and human-friendly failure information.
type ValidationError struct {
	Type    IdentifierType
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Type, e.Code, e.Message)
}

// ValidationReport is the canonical validation output for public APIs.
type ValidationReport struct {
	Type       IdentifierType
	Input      string
	Normalized string
	Valid      bool
	Error      *ValidationError
}

// ParsedIdentifier is a generic parsed representation.
type ParsedIdentifier struct {
	Type       IdentifierType
	Normalized string
	Fields     map[string]string
}

// Suggestion describes a candidate correction for invalid input.
type Suggestion struct {
	Type      IdentifierType
	Candidate string
	Reason    string
}
