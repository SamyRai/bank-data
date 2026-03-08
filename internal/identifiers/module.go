package identifiers

import "fmt"

// ValidationError is a shared internal validation error.
type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Module defines behavior each identifier module must provide.
type Module interface {
	Normalize(input string) string
	DetectCandidate(normalized string) bool
	Validate(normalized string) error
	Parse(normalized string) (map[string]string, error)
}
