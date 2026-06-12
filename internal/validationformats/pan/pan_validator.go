// Package cardpan provides validation for Card Number (PAN)
package cardpan

import "errors"

var (
	ErrInvalidPANFormat = errors.New("invalid card number format")
	ErrInvalidPANLength = errors.New("invalid card number length")
	ErrInvalidLuhn      = errors.New("invalid card number (Luhn check failed)")
)

// Validator provides card number validation interface.
type Validator interface {
	Validate(card string) error
}

// SimpleValidator is a basic implementation.
type SimpleValidator struct{}

func (v *SimpleValidator) Validate(_ string) error {
	// TODO: Implement format, length, and Luhn validation (priority: high, effort: 0.5d)
	return nil
}
