package iban_test

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	internal "github.com/SamyRai/bank-data/internal/iban"
	"github.com/SamyRai/bank-data/pkg/iban"
)

func FuzzIBANValidate(f *testing.F) {
	// Seed with valid and invalid IBANs from testdata
	f.Add("DE89370400440532013000")      // valid
	f.Add("GB82WEST12345698765432")      // valid
	f.Add("FR1420041010050500013M02606") // valid
	f.Add("INVALIDIBAN123")              // invalid
	f.Add("")                            // invalid

	// Add more from testdata/iban_test_vectors.txt
	file, err := os.Open("../../testdata/iban_test_vectors.txt")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 || line[0] == '#' {
				continue
			}
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				f.Add(parts[0])
			}
		}
		file.Close()
	}

	// Add random fuzz cases
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		length := rand.Intn(40) // up to 40 chars
		iban := make([]byte, length)
		for j := range iban {
			// Random A-Z, 0-9, and some invalid chars
			c := rand.Intn(40)
			if c < 26 {
				iban[j] = byte('A' + c)
			} else if c < 36 {
				iban[j] = byte('0' + c - 26)
			} else {
				iban[j] = byte('!') // invalid
			}
		}
		f.Add(string(iban))
	}

	service := iban.NewService(
		internal.NewValidator(),
		internal.NewParser(),
		internal.NewDetector(),
		nil, // pass nil for registry
	)
	f.Fuzz(func(t *testing.T, s string) {
		_ = service.Validate(s) // Should not panic or crash
	})
}
