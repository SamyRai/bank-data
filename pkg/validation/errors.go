package validation

import "fmt"

// ValidationError is a domain-specific error for validation failures.
type ValidationError struct {
	Field   string
	Message string
	Value   any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s (value: %v)", e.Field, e.Message, e.Value)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string, value any) error {
	return &ValidationError{Field: field, Message: message, Value: value}
}
