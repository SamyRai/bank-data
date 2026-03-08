package sepa

import "testing"

func FuzzSEPACreditorValidate(f *testing.F) {
	f.Add("DE98ZZZ09999999999")
	f.Add("DE99ZZZ09999999999")
	f.Add("INVALID")
	f.Add("")

	v := NewCreditorIDValidator()
	f.Fuzz(func(t *testing.T, s string) {
		_ = v.Validate(s)
	})
}
