package validation

import (
	"context"
	"iter"
	"runtime"
	"sync"
)

// Validator validates an input and returns an output.
type Validator[T any, R any] interface {
	Validate(input T) R
}

// Engine is a generic concurrent validation engine used by public facades.
type Engine[T any, R any] struct {
	validator Validator[T, R]
	workers   int
}

// NewEngine constructs a new engine. workers<=0 means NumCPU.
func NewEngine[T any, R any](v Validator[T, R], workers int) *Engine[T, R] {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &Engine[T, R]{validator: v, workers: workers}
}

// ValidateBatch validates inputs concurrently and preserves input order.
func (e *Engine[T, R]) ValidateBatch(ctx context.Context, inputs []T) []R {
	results := make([]R, len(inputs))
	if len(inputs) == 0 {
		return results
	}

	jobs := make(chan int, e.workers)
	var wg sync.WaitGroup

	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					results[idx] = e.validator.Validate(inputs[idx])
				}
			}
		}()
	}

loop:
	for i := range inputs {
		select {
		case <-ctx.Done():
			break loop
		case jobs <- i:
		}
	}

	close(jobs)
	wg.Wait()
	return results
}

// StreamValidate validates items from an iterator and yields results as they complete.
func (e *Engine[T, R]) StreamValidate(ctx context.Context, seq iter.Seq[T]) iter.Seq[R] {
	return func(yield func(R) bool) {
		streamCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		sem := make(chan struct{}, e.workers)
		out := make(chan R, e.workers*2)
		var wg sync.WaitGroup

		go func() {
			for in := range seq {
				select {
				case <-streamCtx.Done():
					goto done
				case sem <- struct{}{}:
					wg.Add(1)
					go func(input T) {
						defer wg.Done()
						defer func() { <-sem }()
						res := e.validator.Validate(input)
						select {
						case <-streamCtx.Done():
						case out <- res:
						}
					}(in)
				}
			}
		done:
			wg.Wait()
			close(out)
		}()

		for res := range out {
			if !yield(res) {
				cancel()
				return
			}
		}
	}
}
