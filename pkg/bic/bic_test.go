package bic_test

import (
	"testing"

	"github.com/SamyRai/bank-data/pkg/bic"
)

func TestValidate(t *testing.T) {
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
		{"DEUTDEFGHJK", false, "not in directory, but valid format"},
		{"DEUTUS33", false, "invalid country code"},
		{"DEUTDEFF@#", false, "invalid chars"},
		{"DEUTDE", false, "too short"},
		{"DEUTDEFF1234", false, "too long"},
	}
	for _, c := range cases {
		err := bic.Validate(c.bic)
		if (err == nil) != c.ok {
			t.Errorf("%s: Validate(%q) = %v, want ok=%v", c.name, c.bic, err, c.ok)
		}
	}
}
