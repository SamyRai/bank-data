package iban

import "testing"

func FuzzIBANParseAndDetect(f *testing.F) {
	f.Add("DE89370400440532013000")
	f.Add("GB82WEST12345698765432")
	f.Add("INVALID")
	f.Add("")

	svc := NewService(nil, nil, nil, nil)
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = svc.Parse(s)
		_, _ = svc.Detect(s)
	})
}
