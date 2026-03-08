package validation

import "fmt"

// Error is a typed validation error.
type Error struct {
	Code    string
	Field   string
	Message string
}

func (e *Error) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s %s: %s", e.Code, e.Field, e.Message)
}

func NewError(code, field, message string) *Error {
	return &Error{Code: code, Field: field, Message: message}
}
