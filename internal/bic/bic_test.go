package bic

import (
	"testing"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
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
		err := Validate(c.bic)
		if (err == nil) != c.ok {
			t.Errorf("%s: Validate(%q) = %v, want ok=%v", c.name, c.bic, err, c.ok)
		}
	}
}

func TestValidate_WithBankBICMap(t *testing.T) {
	// Minimal test map for DE: bankCode is first 4 letters of BIC
	m := bicmap.BankBICMap{
		"DE": {
			"PBNK": bicmap.BankBICEntry{BankCode: "PBNK", BIC: "PBNKDEFF", BankName: "Postbank Ndl Deutsche Bank", Country: "DE"},
		},
	}
	SetBankBICMap(m)

	cases := []struct {
		bic  string
		ok   bool
		name string
	}{
		{"PBNKDEFF", true, "valid BIC from map"},
		{"XXXXDEFF", false, "bank code not in map"},
		{"PBNKDEFF1", false, "invalid length (9)"},
		{"PBNKDEFFXXX", true, "valid 11-char BIC from map (branch)"},
		{"PBNKDEFG", false, "valid format, not in map or directory"},
	}
	for _, c := range cases {
		err := Validate(c.bic)
		if (err == nil) != c.ok {
			t.Errorf("%s: Validate(%q) = %v, want ok=%v", c.name, c.bic, err, c.ok)
		}
	}
}
