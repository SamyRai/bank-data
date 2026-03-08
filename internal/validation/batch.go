// Package validation provides deprecated compatibility wrappers.
package validation

import (
	"context"

	corevalidation "github.com/SamyRai/bank-data/internal/core/validation"
)

// BatchValidator is kept for backward compatibility and delegates to core engine.
type BatchValidator[T any, R any] struct {
	Validator corevalidation.Validator[T, R]
	Workers   int
}

func (b *BatchValidator[T, R]) ValidateBatch(inputs []T) []R {
	engine := corevalidation.NewEngine[T, R](b.Validator, b.Workers)
	return engine.ValidateBatch(context.Background(), inputs)
}
