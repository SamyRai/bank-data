// internal/sepa/sepa_validator_test.go
package sepa

import (
	"errors"
	"strings"
	"testing"
)

// goodDE() / goodFR() help produce internally consistent samples.
// They guard against test rot if patterns change.
func goodDE() string {
	ci, err := BuildCI("DE", "ZZZ", "00000000000") // 11-digit national part
	if err != nil {
		panic(err)
	}
	return ci
}

func goodFR() string {
	ci, err := BuildCI("FR", "ZZZ", "123456") // 6-digit national part
	if err != nil {
		panic(err)
	}
	return ci
}

func TestCreditorIDValidator_DefaultValidate(t *testing.T) {
	v := NewDefaultCreditorIDValidator()

	de := goodDE()
	fr := goodFR()

	// --- happy paths ---
	valid := []string{
		de,                  // canonical DE
		strings.ToLower(de), // case-insensitive
		insertSpaces(de),    // embedded spaces ok
		fr,                  // canonical FR
	}

	for _, id := range valid {
		if err := v.Validate(id); err != nil {
			t.Errorf("expected %q to be valid, got %v", id, err)
		}
	}

	// --- unhappy paths (default = allow unknown country) ---
	tests := []struct {
		id      string
		wantErr error
	}{
		{de[:len(de)-1], ErrInvalidFormat},                 // too short (fails DE pattern)
		{de + "0", ErrInvalidFormat},                       // too long (fails DE pattern)
		{de[:len(de)-1] + "O", ErrInvalidFormat},           // alpha in numeric national part
		{replaceCheckDigits(de, "00"), ErrInvalidChecksum}, // bad checksum
		{"PT00ZZZ12345", ErrInvalidFormat},                 // PT too short
		{"XX12ZZZ123456", ErrInvalidChecksum},              // unknown CC allowed → format ok, checksum decides
		{"", ErrInvalidLength},                             // empty
		{"   ", ErrInvalidLength},                          // whitespace
	}

	for _, tc := range tests {
		if err := v.Validate(tc.id); err == nil || !errors.Is(err, tc.wantErr) {
			t.Errorf("%q: expected %v, got %v", tc.id, tc.wantErr, err)
		}
	}
}

func TestCreditorIDValidator_Strict(t *testing.T) {
	// Only DE + FR allowed
	v := NewStrictCreditorIDValidator([]string{"DE", "FR"})

	de := goodDE()
	if err := v.Validate(de); err != nil {
		t.Fatalf("strict: expected %q valid: %v", de, err)
	}
	// Unknown country now fails earlier
	if err := v.Validate("ES00ZZZ123456789"); !errors.Is(err, ErrInvalidFormat) && !errors.Is(err, ErrUnknownCountry) {
		t.Fatalf("strict: expected unknown/format err, got %v", err)
	}
}

// -----------------------------------------------------------------------------
// Helper funcs for tests
// -----------------------------------------------------------------------------

// insertSpaces inserts spaces every 4 chars for readability.
func insertSpaces(s string) string {
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && i%4 == 0 {
			out = append(out, ' ')
		}
		out = append(out, c)
	}
	return string(out)
}

// replaceCheckDigits swaps the 2-digit checksum at pos 2..4
func replaceCheckDigits(ci, cd string) string {
	if len(ci) < 4 {
		return ci
	}
	return ci[:2] + cd + ci[4:]
}
