package sepa

import (
	"testing"
)

func TestCreditorIDValidator_Validate(t *testing.T) {
	v := NewCreditorIDValidator()

	cases := []struct {
		id    string
		valid bool
		name  string
	}{
		{"DE98ZZZ09999999999", true, "valid German SCI"},
		{"AT80ZZZ09999999999", true, "valid Austrian SCI"},
		{"FR41ZZZ09999999999", true, "valid French SCI"},
		{"DE98ZZZ0999999999X", false, "invalid char"},
		{"DE99ZZZ09999999999", false, "invalid checksum"},
		{"DEXX", false, "too short"},
	}

	for _, c := range cases {
		res := v.Validate(c.id)
		if res.Valid != c.valid {
			t.Errorf("%s: Validate(%q) = %v, want valid=%v", c.name, c.id, res.Error, c.valid)
		}
	}
}
