package countrymeta_test

import (
	"testing"

	"github.com/SamyRai/bank-data/internal/countrymeta"
)

func TestRegistry(t *testing.T) {
	if len(countrymeta.Registry) == 0 {
		t.Fatal("Registry is empty")
	}

	for code, meta := range countrymeta.Registry {
		if code != meta.CountryCode {
			t.Errorf("Mismatched country code for %s: expected %s, got %s", code, code, meta.CountryCode)
		}
		if meta.Length <= 4 {
			t.Errorf("Invalid length for %s: %d", code, meta.Length)
		}
		if meta.BBANRegex == nil {
			t.Errorf("Missing BBAN regex for %s", code)
		}
	}
}
