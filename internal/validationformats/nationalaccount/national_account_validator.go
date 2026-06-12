// Package nationalaccount provides validation for national bank account numbers
package nationalaccount

import "errors"

var (
	ErrInvalidAccountFormat   = errors.New("invalid national account format")
	ErrInvalidAccountLength   = errors.New("invalid national account length")
	ErrInvalidAccountChecksum = errors.New("invalid national account checksum")
)

// Validator provides national account validation interface.
type Validator interface {
	Validate(account string) error
}

// SimpleValidator is a basic implementation.
type SimpleValidator struct{}

func (v *SimpleValidator) Validate(_ string) error {
	// TODO: Implement country-specific format, length, and checksum validation (priority: high, effort: 1d)
	return nil
}
