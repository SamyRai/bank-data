package iban

import (
	"testing"

	"github.com/SamyRai/bank-data/internal/log"
)

func TestMain(m *testing.M) {
	log.SetMinLevel(log.LevelError)
	m.Run()
}

func BenchmarkValidateIBANChecksum(b *testing.B) {
	input := "DE37040044053201300089" // rearranged for checksum check
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateIBANChecksum(input)
	}
}

func BenchmarkCharacterCheck(b *testing.B) {
	input := "DE89370400440532013000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(input); j++ {
			r := input[j]
			_ = (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z')
		}
	}
}
