package validation_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/validation"
)

func TestMain(m *testing.M) {
	log.SetMinLevel(log.LevelError)
	m.Run()
}

type benchValidator struct {
	delay time.Duration
}

func (v *benchValidator) Validate(input string) bool {
	// Simulate some work (e.g., regex + checksum)
	if v.delay > 0 {
		time.Sleep(v.delay)
	}
	return len(input) > 0
}

func BenchmarkEngine_ValidateBatch_Scaling(b *testing.B) {
	ctx := context.Background()

	// Varying complexity: 0 (instant), 10µs (typical validation), 100µs (heavy)
	delays := []time.Duration{0, 10 * time.Microsecond}
	sizes := []int{100, 1000}
	workerCounts := []int{1, 2, 4, 8, 16}

	for _, delay := range delays {
		v := &benchValidator{delay: delay}
		for _, size := range sizes {
			inputs := make([]string, size)
			for i := 0; i < size; i++ {
				inputs[i] = "test-data-for-validation-engine-benchmarking"
			}

			for _, w := range workerCounts {
				name := fmt.Sprintf("delay=%v/size=%d/workers=%d", delay, size, w)
				b.Run(name, func(b *testing.B) {
					engine := validation.NewEngine[string, bool](v, w)
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_ = engine.ValidateBatch(ctx, inputs)
					}
					b.ReportMetric(float64(b.N*size)/b.Elapsed().Seconds(), "ops/s")
				})
			}
		}
	}
}
