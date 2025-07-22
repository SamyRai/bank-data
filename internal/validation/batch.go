// Package validation provides a generic, parallel batch validation engine for any validator.
package validation

import (
	"sync"

	"github.com/SamyRai/bank-data/pkg/validation"
)

// BatchValidator runs validations in parallel using a worker pool.
type BatchValidator[T any, R any] struct {
	Validator validation.Validator[T, R]
	Workers   int // Number of parallel workers
}

// ValidateBatch validates a slice of inputs in parallel and returns results in input order.
func (b *BatchValidator[T, R]) ValidateBatch(inputs []T) []R {
	n := len(inputs)
	results := make([]R, n)
	var wg sync.WaitGroup
	jobs := make(chan int, n)

	// Start workers
	for w := 0; w < b.getWorkers(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = b.Validator.Validate(inputs[i])
			}
		}()
	}

	// Send jobs
	for i := range inputs {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results
}

func (b *BatchValidator[T, R]) getWorkers() int {
	if b.Workers > 0 {
		return b.Workers
	}
	return 4 // default
}
