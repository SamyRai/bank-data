package vop

import (
	"context"
	"testing"
)

func FuzzMatcherMatch(f *testing.F) {
	f.Add("ACME LTD", "Acme Limited")
	f.Add("John Doe", "John Doe")
	f.Add("", "")

	m := NewMatcher()
	for _, tc := range []MatchRequest{
		{SuppliedName: "Acme BV", ExpectedName: "ACME B.V."},
		{SuppliedName: "Muller GmbH", ExpectedName: "Mueller GMBH"},
	} {
		f.Add(tc.SuppliedName, tc.ExpectedName)
	}

	f.Fuzz(func(t *testing.T, supplied, expected string) {
		req := MatchRequest{SuppliedName: supplied, ExpectedName: expected}
		res1, _ := m.Match(req)
		res2, _ := m.Verify(context.Background(), req)
		if res1.Category != res2.Category {
			t.Fatalf("verify/match category mismatch: %s vs %s", res1.Category, res2.Category)
		}
		if res1.Score < 0 || res1.Score > 1 {
			t.Fatalf("score out of bounds: %.4f", res1.Score)
		}
	})
}
