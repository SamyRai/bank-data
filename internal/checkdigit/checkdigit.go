// Package checkdigit provides IBAN check digit calculation utilities.
package checkdigit

import (
	"errors"
	"strconv"
	"strings"
)

// ComputeIBANCheckDigits calculates the check digits for a given country code and BBAN.
// Returns the two-digit check digits as a string, or an error if input is invalid.
func ComputeIBANCheckDigits(countryCode, bban string) (string, error) {
	countryCode = strings.ToUpper(countryCode)
	if len(countryCode) != 2 {
		return "", ErrInvalidCountryCode
	}
	// Step 1: Form intermediate string: BBAN + country code + "00"
	intermediate := bban + countryCode + "00"
	// Step 2: Normalize and stream MOD-97
	rem := 0
	for _, r := range intermediate {
		switch {
		case r >= '0' && r <= '9':
			rem = (rem*10 + int(r-'0')) % 97
		case r >= 'A' && r <= 'Z':
			v := int(r - 'A' + 10)
			rem = (rem*10 + v/10) % 97
			rem = (rem*10 + v%10) % 97
		default:
			return "", ErrInvalidCharacter
		}
	}
	// Step 3: Compute check digits
	checkDigits := 98 - rem
	if checkDigits < 2 {
		return "", ErrInvalidNumeric // Per industry Q&A, reject < 02
	}
	return leftPad2(checkDigits), nil
}

// leftPad2 pads a number to two digits with leading zero if needed.
func leftPad2(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// Error types
var (
	ErrInvalidCountryCode = errors.New("invalid country code")
	ErrInvalidCharacter   = errors.New("invalid character in BBAN or country code")
	ErrInvalidNumeric     = errors.New("invalid numeric conversion")
)
