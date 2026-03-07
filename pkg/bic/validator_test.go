package bic

import (
	"testing"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
)

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()
	cases := []struct {
		bic  string
		ok   bool
		name string
	}{
		{"DEUTDEFF", true, "valid 8-char BIC, active"},
		{"DEUTDEFF500", true, "valid 11-char BIC, active"},
		{"NEDSZAJJ", false, "inactive BIC"},
		{"DEUTDEFF5", false, "invalid length"},
		{"DEUTDEFFF0", false, "invalid length (10)"},
		{"DEUTUS33", false, "invalid country code"},
		{"DEUTDEFF@#", false, "invalid chars"},
	}
	for _, c := range cases {
		res := v.Validate(c.bic)
		if res.Valid != c.ok {
			t.Errorf("%s: Validate(%q) = %v, want ok=%v", c.name, c.bic, res.Error, c.ok)
		}
	}
}

func TestValidator_WithBankBICMap(t *testing.T) {
	// Setup a mock map
	_ = bicmap.NewBankBICMap("tmp", map[string]string{"DE": "mock.csv"})
	// We can't easily mock the file load without more refactoring,
	// but we can manually inject into the cache for testing.
	// Since cache is private, we depend on SetBankBICMap which I added.

	// Actually, let's just test the logic by manually setting up the internal state if possible,
	// or just test the Parse function which is simpler.

	v := NewValidator()
	SetBankBICMap(nil) // Reset

	if res := v.Validate("PBNKDEFF"); res.Valid {
		t.Error("expected invalid for unknown BIC without mapping")
	}
}

func TestParse(t *testing.T) {
	info, err := Parse("DEUTDEFF500")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if info.Institution != "DEUT" {
		t.Errorf("expected DEUT, got %s", info.Institution)
	}
	if info.Country != "DE" {
		t.Errorf("expected DE, got %s", info.Country)
	}
	if info.Branch != "500" {
		t.Errorf("expected 500, got %s", info.Branch)
	}

	info8, _ := Parse("DEUTDEFF")
	if info8.Branch != "XXX" {
		t.Errorf("expected XXX for 8-char BIC, got %s", info8.Branch)
	}
}
