package vop

import "testing"

func BenchmarkMatcher_Match(b *testing.B) {
	m := NewMatcher()
	req := MatchRequest{SuppliedName: "Acme Payments Ltd", ExpectedName: "ACME PAYMENTS LIMITED"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := m.Match(req); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
