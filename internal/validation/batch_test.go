package validation

import (
	"errors"
	"testing"
)

type testBatchItem struct {
	ID    int
	Value string
}

type testResult struct {
	ID    int
	Valid bool
	Error error
}

type testValidator struct{}

func (v *testValidator) Validate(item testBatchItem) testResult {
	if item.Value == "invalid" {
		return testResult{ID: item.ID, Valid: false, Error: errors.New("invalid value")}
	}
	return testResult{ID: item.ID, Valid: true, Error: nil}
}

func TestValidateBatch(t *testing.T) {
	validator := &BatchValidator[testBatchItem, testResult]{
		Validator: &testValidator{},
		Workers:   2,
	}
	items := []testBatchItem{
		{ID: 1, Value: "valid"},
		{ID: 2, Value: "invalid"},
		{ID: 3, Value: "valid"},
		{ID: 4, Value: "invalid"},
	}

	results := validator.ValidateBatch(items)

	var invalidResults []testResult
	for _, res := range results {
		if !res.Valid {
			invalidResults = append(invalidResults, res)
		}
	}

	if len(invalidResults) != 2 {
		t.Fatalf("expected 2 validation errors, got %d", len(invalidResults))
	}

	if invalidResults[0].ID != 2 {
		t.Errorf("expected error for ID 2, got ID %d", invalidResults[0].ID)
	}
	if invalidResults[1].ID != 4 {
		t.Errorf("expected error for ID 4, got ID %d", invalidResults[1].ID)
	}
}
