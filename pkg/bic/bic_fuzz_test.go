package bic_test

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/SamyRai/bank-data/pkg/bic"
)

func FuzzBICValidate(f *testing.F) {
	// Seed with valid and invalid BICs
	f.Add("DEUTDEFF")    // valid
	f.Add("DEUTDEFF500") // valid
	f.Add("NEDSZAJJ")    // inactive
	f.Add("DEUTDEFF5")   // invalid length
	f.Add("DEUTDEFFF0")  // invalid length
	f.Add("DEUTDEFGHJK") // not in directory, valid format
	f.Add("DEUTUS33")    // invalid country code
	f.Add("DEUTDEFF@#")  // invalid chars
	f.Add("")            // empty

	// Add more from testdata/bic_test_vectors.txt
	file, err := os.Open("../../testdata/bic_test_vectors.txt")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 || line[0] == '#' {
				continue
			}
			parts := strings.Split(line, ",")
			if len(parts) >= 1 {
				f.Add(parts[0])
			}
		}
		if err := file.Close(); err != nil {
			fmt.Printf("error closing bic_test_vectors.txt: %v\n", err)
			// TODO: log error closing bic_test_vectors.txt, priority: low, effort: 1m
		}
	}

	// Add random fuzz cases using crypto/rand
	for i := 0; i < 100; i++ {
		length := 8
		b := make([]byte, 1)
		if _, err := rand.Read(b); err == nil {
			length += int(b[0]) % 4 // 8 to 11 chars
		}
		bic := make([]byte, length)
		for j := range bic {
			b := make([]byte, 1)
			if _, err := rand.Read(b); err == nil {
				c := int(b[0]) % 40
				if c < 26 {
					bic[j] = byte('A' + c)
				} else if c < 36 {
					bic[j] = byte('0' + c - 26)
				} else {
					bic[j] = byte('@') // invalid
				}
			} else {
				bic[j] = 'X'
			}
		}
		f.Add(string(bic))
	}

	f.Fuzz(func(_ *testing.T, s string) {
		_ = bic.Validate(s) // Should not panic or crash
	})
}
