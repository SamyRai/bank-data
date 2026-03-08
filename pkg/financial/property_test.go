package financial

import (
	"strings"
	"testing"

	"pgregory.net/rapid"
)

func TestProperty_NormalizationIdempotence(t *testing.T) {
	svc := NewService()
	types := []IdentifierType{
		IdentifierIBAN,
		IdentifierBIC,
		IdentifierSEPACreditor,
		IdentifierLEI,
		IdentifierISIN,
		IdentifierPAN,
		IdentifierVAT,
		IdentifierNationalAccountUK,
	}

	rapid.Check(t, func(rt *rapid.T) {
		input := rapid.StringMatching(`[A-Za-z0-9 _\-\.]{0,48}`).Draw(rt, "input")
		for _, typ := range types {
			r1, _ := svc.Validate(input, typ)
			r2, _ := svc.Validate(r1.Normalized, typ)
			if r1.Normalized != r2.Normalized {
				rt.Fatalf("normalization not idempotent for %s: %q != %q", typ, r1.Normalized, r2.Normalized)
			}
		}
	})
}

func TestProperty_ParseValidateConsistency(t *testing.T) {
	svc := NewService()
	validByType := map[IdentifierType][]string{
		IdentifierIBAN:         {"DE89370400440532013000", "GB82WEST12345698765432"},
		IdentifierBIC:          {"DEUTDEFF", "DEUTDEFF500"},
		IdentifierSEPACreditor: {"DE98ZZZ09999999999"},
		IdentifierLEI:          {"529900T8BM49AURSDO78"},
		IdentifierISIN:         {"US0378331005"},
		IdentifierPAN:          {"4111111111111111", "5555555555554444"},
		IdentifierVAT:          {"DE136695976", "FR40303265045", "NL123456782B12", "IT12345678903", "ES12345678Z"},
		IdentifierNationalAccountUK: {"20-00-00 55779911"},
	}

	rapid.Check(t, func(rt *rapid.T) {
		keys := make([]IdentifierType, 0, len(validByType))
		for k := range validByType {
			keys = append(keys, k)
		}
		typ := keys[rapid.IntRange(0, len(keys)-1).Draw(rt, "type_idx")]
		samples := validByType[typ]
		in := samples[rapid.IntRange(0, len(samples)-1).Draw(rt, "sample_idx")]

		rep, err := svc.Validate(in, typ)
		if err != nil || !rep.Valid {
			rt.Fatalf("valid sample rejected: type=%s input=%q err=%v rep=%+v", typ, in, err, rep)
		}
		parsed, err := svc.Parse(in, typ)
		if err != nil {
			rt.Fatalf("parse failed for valid sample: type=%s input=%q err=%v", typ, in, err)
		}
		if parsed.Normalized != rep.Normalized {
			rt.Fatalf("normalized mismatch parse/validate: %q vs %q", parsed.Normalized, rep.Normalized)
		}
	})
}

func TestProperty_ChecksumMutationRejection(t *testing.T) {
	svc := NewService()
	targets := []struct {
		typ   IdentifierType
		input string
	}{
		{IdentifierIBAN, "DE89370400440532013000"},
		{IdentifierLEI, "529900T8BM49AURSDO78"},
		{IdentifierISIN, "US0378331005"},
		{IdentifierPAN, "4111111111111111"},
		{IdentifierVAT, "DE136695976"},
	}

	rapid.Check(t, func(rt *rapid.T) {
		t := targets[rapid.IntRange(0, len(targets)-1).Draw(rt, "target_idx")]
		mutated := mutateLastAlnum(t.input)
		rep, err := svc.Validate(mutated, t.typ)
		if err == nil || rep.Valid {
			rt.Fatalf("checksum mutation unexpectedly valid: type=%s original=%q mutated=%q rep=%+v", t.typ, t.input, mutated, rep)
		}
	})
}

func mutateLastAlnum(s string) string {
	if s == "" {
		return s
	}
	last := s[len(s)-1]
	replacement := byte('0')
	switch {
	case last >= '0' && last <= '9':
		replacement = byte('0' + ((last-'0'+1)%10))
	case last >= 'A' && last <= 'Z':
		replacement = byte('A' + ((last-'A'+1)%26))
	case last >= 'a' && last <= 'z':
		replacement = byte('a' + ((last-'a'+1)%26))
	default:
		replacement = 'X'
	}
	return strings.TrimSuffix(s, string(last)) + string(replacement)
}
