package financial

import "testing"

func FuzzService_DetectValidateParse(f *testing.F) {
	svc := NewService()
	seeds := []string{
		"DE89370400440532013000",
		"DEUTDEFF",
		"DE98ZZZ09999999999",
		"529900T8BM49AURSDO78",
		"US0378331005",
		"4111111111111111",
		"DE136695976",
		"20-00-00 55779911",
		"INVALID",
		"",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		_, _ = svc.Detect(input)
		_, _ = svc.Validate(input, "")
		_, _ = svc.Parse(input, "")

		hints := []IdentifierType{
			IdentifierIBAN,
			IdentifierBIC,
			IdentifierSEPACreditor,
			IdentifierLEI,
			IdentifierISIN,
			IdentifierPAN,
			IdentifierVAT,
			IdentifierNationalAccountUK,
		}
		for _, hint := range hints {
			_, _ = svc.Validate(input, hint)
			_, _ = svc.Parse(input, hint)
		}
	})
}
