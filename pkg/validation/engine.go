package validation

import (
	"context"
	"iter"
	"runtime"
	"sync"
)

// Engine is a high-performance validation engine for processing batches of data.
// It leverages Go 1.26 optimizations like the Green Tea GC and fast small-object allocation.
type Engine[T any, R any] struct {
	validator Validator[T, R]
	workers   int
}

// NewEngine creates a new validation engine with the given validator.
// If workers is <= 0, it defaults to runtime.NumCPU().
func NewEngine[T any, R any](v Validator[T, R], workers int) *Engine[T, R] {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &Engine[T, R]{
		validator: v,
		workers:   workers,
	}
}

// ValidateBatch processes a slice of inputs in parallel and returns results in the same order.
func (e *Engine[T, R]) ValidateBatch(ctx context.Context, inputs []T) []R {
	n := len(inputs)
	results := make([]R, n)
	if n == 0 {
		return results
	}

	var wg sync.WaitGroup
	jobs := make(chan int, e.workers)

	// Start workers
	for w := 0; w < e.workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					results[i] = e.validator.Validate(inputs[i])
				}
			}
		}()
	}

	// Distribute work
loop:
	for i := 0; i < n; i++ {
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

// StreamValidate processes an iterator of inputs and returns an iterator of results.
// This is memory-efficient for large datasets as it doesn't require loading all inputs at once.
func (e *Engine[T, R]) StreamValidate(ctx context.Context, inputIter iter.Seq[T]) iter.Seq[R] {
	return func(yield func(R) bool) {
		// Use a dedicated context for this stream to handle early exit from yield
		streamCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		resultsChan := make(chan R, e.workers*2)
		var wg sync.WaitGroup

		// Start the producer
		go func() {
			semaphore := make(chan struct{}, e.workers)
			for input := range inputIter {
				select {
				case <-streamCtx.Done():
					goto end
				case semaphore <- struct{}{}:
					wg.Add(1)
					go func(in T) {
						defer wg.Done()
						defer func() { <-semaphore }()

						res := e.validator.Validate(in)
						select {
						case <-streamCtx.Done():
						case resultsChan <- res:
						}
					}(input)
				}
			}
		end:
			wg.Wait()
			close(resultsChan)
		}()

		// Consumer
		for res := range resultsChan {
			if !yield(res) {
				return
			}
		}
	}
}
