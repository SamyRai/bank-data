package iban_test

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"testing"

	internal "github.com/SamyRai/bank-data/internal/validationformats/iban"
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
		if err := file.Close(); err != nil {
			fmt.Printf("error closing iban_test_vectors.txt: %v\n", err)
			// TODO: log error closing iban_test_vectors.txt, priority: low, effort: 1m
		}
	}

	// Add random fuzz cases using crypto/rand
	for i := 0; i < 100; i++ {
		length := 0
		b := make([]byte, 1)
		if _, err := rand.Read(b); err == nil {
			length = int(b[0]) % 40 // up to 40 chars
		}
		iban := make([]byte, length)
		for j := range iban {
			b := make([]byte, 1)
			if _, err := rand.Read(b); err == nil {
				c := int(b[0]) % 40
				if c < 26 {
					iban[j] = byte('A' + c)
				} else if c < 36 {
					iban[j] = byte('0' + c - 26)
				} else {
					iban[j] = byte('!') // invalid
				}
			} else {
				iban[j] = 'X'
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
	f.Fuzz(func(_ *testing.T, s string) {
		_ = service.Validate(s) // Should not panic or crash
	})
}
