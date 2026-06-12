// Package vat provides validation for VAT numbers
package vat

import "errors"

var (
	ErrInvalidVATFormat   = errors.New("invalid VAT number format")
	ErrInvalidVATLength   = errors.New("invalid VAT number length")
	ErrInvalidVATChecksum = errors.New("invalid VAT number checksum")
)

// Validator provides VAT number validation interface.
type Validator interface {
	Validate(vat string) error
}

// SimpleValidator is a basic implementation.
type SimpleValidator struct{}

func (v *SimpleValidator) Validate(_ string) error {
	// TODO: Implement format, length, and checksum validation (priority: high, effort: 0.5d)
	return nil
}
