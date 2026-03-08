package iban

import (
	"strings"
	"testing"

	"github.com/SamyRai/bank-data/internal/countrymeta"
)

func TestBuildStructure_ConformsForAllCountries(t *testing.T) {
	if len(countrymeta.Registry) == 0 {
		t.Fatalf("country registry is empty")
	}
	for cc, meta := range countrymeta.Registry {
		s := BuildStructure(meta)
		if len(s) != meta.Length {
			t.Fatalf("%s: structure length %d != meta length %d", cc, len(s), meta.Length)
		}
		if !strings.HasPrefix(s, "CCKK") {
			t.Fatalf("%s: structure must start with CCKK, got %s", cc, s)
		}
		for i := 0; i < len(s); i++ {
			ch := s[i]
			switch ch {
			case 'C', 'K', 'B', 'A', 'X':
			default:
				t.Fatalf("%s: invalid structure symbol %q at pos %d", cc, ch, i)
			}
		}
		if meta.Length < 4 {
			t.Fatalf("%s: invalid country length < 4", cc)
		}
		bbanLen := meta.Length - 4
		if meta.BankStart < 0 || meta.BankEnd < 0 || meta.AccountStart < 0 || meta.AccountEnd < 0 {
			t.Fatalf("%s: negative BBAN ranges: %+v", cc, meta)
		}
		if meta.BankEnd > bbanLen || meta.AccountEnd > bbanLen {
			t.Fatalf("%s: BBAN range exceeds length=%d: %+v", cc, bbanLen, meta)
		}
	}
}
