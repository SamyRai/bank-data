// Package bic provides the public interface for BIC (SWIFT) code validation.
package bic

import internalbic "github.com/SamyRai/bank-data/internal/validationformats/bic"

// Validate validates a BIC (SWIFT) code for format, length, and country code.
// It delegates to the internal/bic/validator.go implementation.
func Validate(bic string) error {
	return internalbic.Validate(bic)
}
