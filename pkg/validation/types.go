// Package validation exposes typed validation primitives.
package validation

import "sync"

// Validator validates an input and returns a result.
type Validator[T any, R any] interface {
	Validate(input T) R
}

// CrossFieldValidator validates interdependent values.
type CrossFieldValidator[T any, R any] interface {
	ValidateFields(input T) R
}

// ValidationResult is a typed, serializable validation outcome.
type ValidationResult struct {
	Input   string
	Valid   bool
	Code    string
	Message string
}

// Registry stores validators with typed keys/values.
type Registry[K comparable, V any] struct {
	mu         sync.RWMutex
	validators map[K]V
}

func NewRegistry[K comparable, V any]() *Registry[K, V] {
	return &Registry[K, V]{validators: make(map[K]V)}
}

func (r *Registry[K, V]) Register(key K, v V) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.validators[key] = v
}

func (r *Registry[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.validators[key]
	return v, ok
}
