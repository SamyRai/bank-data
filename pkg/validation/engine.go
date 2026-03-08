package validation

import (
	"context"
	"iter"

	corevalidation "github.com/SamyRai/bank-data/internal/core/validation"
)

// Engine is a public alias to the shared internal validation engine.
type Engine[T any, R any] = corevalidation.Engine[T, R]

func NewEngine[T any, R any](v Validator[T, R], workers int) *Engine[T, R] {
	return corevalidation.NewEngine[T, R](v, workers)
}

func ValidateBatch[T any, R any](ctx context.Context, v Validator[T, R], workers int, inputs []T) []R {
	return NewEngine[T, R](v, workers).ValidateBatch(ctx, inputs)
}

func StreamValidate[T any, R any](ctx context.Context, v Validator[T, R], workers int, inputs iter.Seq[T]) iter.Seq[R] {
	return NewEngine[T, R](v, workers).StreamValidate(ctx, inputs)
}
