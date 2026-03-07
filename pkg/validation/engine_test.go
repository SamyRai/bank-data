package validation_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/SamyRai/bank-data/pkg/validation"
)

type mockValidator struct{}

func (v *mockValidator) Validate(input string) string {
	return "validated:" + input
}

func TestEngine_ValidateBatch(t *testing.T) {
	v := &mockValidator{}
	engine := validation.NewEngine[string, string](v, 4)
	inputs := []string{"a", "b", "c", "d", "e"}
	results := engine.ValidateBatch(context.Background(), inputs)

	expected := []string{"validated:a", "validated:b", "validated:c", "validated:d", "validated:e"}
	if !slices.Equal(results, expected) {
		t.Errorf("expected %v, got %v", expected, results)
	}
}

func TestEngine_StreamValidate(t *testing.T) {
	v := &mockValidator{}
	engine := validation.NewEngine[string, string](v, 2)

	inputSeq := func(yield func(string) bool) {
		for i := 0; i < 5; i++ {
			if !yield(fmt.Sprintf("%d", i)) {
				return
			}
		}
	}

	results := []string{}
	for res := range engine.StreamValidate(context.Background(), inputSeq) {
		results = append(results, res)
	}

	// Order might be different in streaming due to parallelism, but here we check content
	// Wait, Engine.StreamValidate currently doesn't preserve order?
	// The implementation uses a resultsChan where workers push. So order is NOT preserved.
	// This is fine for streaming if not specified otherwise.
	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
	slices.Sort(results)
	expected := []string{"validated:0", "validated:1", "validated:2", "validated:3", "validated:4"}
	if !slices.Equal(results, expected) {
		t.Errorf("expected %v, got %v", expected, results)
	}
}
