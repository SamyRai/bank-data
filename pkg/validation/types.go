// Package validation provides shared types and interfaces for validation across the project.
package validation

import "sync"

// Validator is a generic interface for validating any input type and returning a result.
type Validator[T any, R any] interface {
	Validate(input T) R
}

// CrossFieldValidator is an interface for validating relationships between multiple fields or conditional rules.
type CrossFieldValidator[T any, R any] interface {
	ValidateFields(input T) R
}

// ValidationResult is a standard result type for all validators.
type ValidationResult struct {
	Input   any            // The input value validated
	Valid   bool           // Whether the input is valid
	Error   error          // Error if validation failed
	Details map[string]any // Additional details about validation
}

// ValidationRegistry allows registration and lookup of multiple validators by name.
type ValidationRegistry struct {
	validators map[string]any
	mu         sync.RWMutex
}

// NewValidationRegistry creates a new ValidationRegistry instance.
func NewValidationRegistry() *ValidationRegistry {
	return &ValidationRegistry{
		validators: make(map[string]any),
	}
}

// Register adds a validator to the registry under the given name.
func (r *ValidationRegistry) Register(name string, v any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.validators[name] = v
}

// Get retrieves a validator by name from the registry.
func (r *ValidationRegistry) Get(name string) any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.validators[name]
}
