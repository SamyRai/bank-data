package bic

import "testing"

func FuzzBICValidateAndParse(f *testing.F) {
	f.Add("DEUTDEFF")
	f.Add("DEUTDEFF500")
	f.Add("INVALID")
	f.Add("")

	v := NewValidator()
	f.Fuzz(func(_ *testing.T, s string) {
		_ = v.Validate(s)
		_, _ = Parse(s)
	})
}
