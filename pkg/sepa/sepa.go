// Package sepa provides SEPA Creditor ID validation utilities for banking and financial applications.
package sepa

import "github.com/SamyRai/bank-data/internal/validationformats/sepa"

// ValidateCreditorID validates a SEPA Creditor Identifier using the default validator.
func ValidateCreditorID(id string) error {
	return (&sepa.CreditorIDValidator{}).Validate(id)
}
