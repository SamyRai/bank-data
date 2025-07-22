// Package validation provides shared types and interfaces for validation across the project.
package validation

import "sync"

// Validator is a generic interface for validating any input type and returning a result.
type Validator[T any, R any] interface {
	Validate(input T) R
}

// ValidationResult is a standard result type for all validators.
type ValidationResult struct {
	Input   any
	Valid   bool
	Error   error
	Details map[string]any
}

// ValidationRegistry allows registration and lookup of multiple validators by name.
type ValidationRegistry struct {
	validators map[string]any
	mu         sync.RWMutex
}

func NewValidationRegistry() *ValidationRegistry {
	return &ValidationRegistry{
		validators: make(map[string]any),
	}
}

func (r *ValidationRegistry) Register(name string, v any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.validators[name] = v
}

func (r *ValidationRegistry) Get(name string) any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.validators[name]
}
