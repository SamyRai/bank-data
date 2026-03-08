package nationalaccount

import "testing"

func FuzzValidateAndParse(f *testing.F) {
	f.Add("20-00-00 55779911")
	f.Add("20000055779911")
	f.Add("00-00-00 00000000")
	f.Add("INVALID")

	v := NewValidator()
	f.Fuzz(func(t *testing.T, input string) {
		res := v.Validate(input)
		if res.Valid {
			parsed, err := v.Parse(input)
			if err != nil {
				t.Fatalf("valid input failed parse: %q err=%v", input, err)
			}
			if parsed.SortCode == "" || parsed.AccountNumber == "" {
				t.Fatalf("parsed fields missing for valid input: %+v", parsed)
			}
		}
	})
}
