// internal/sepa/sepa_validator.go
package sepa

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// -----------------------------------------------------------------------------
// Public Errors
// -----------------------------------------------------------------------------
var (
	ErrInvalidFormat   = errors.New("invalid SEPA Creditor ID format")
	ErrInvalidLength   = errors.New("invalid SEPA Creditor ID length")
	ErrInvalidChecksum = errors.New("invalid SEPA Creditor ID MOD-97 checksum")
	ErrUnknownCountry  = errors.New("unknown SEPA Creditor ID country code")
)

// -----------------------------------------------------------------------------
// CreditorIDValidator
// -----------------------------------------------------------------------------
//
// Validation steps:
//  1. Normalize (trim, strip internal spaces, upper-case).
//  2. Optional ISO country whitelist (strict mode).
//  3. Generic pattern check (CC2 + check2 + business3 + national[1..28], <=35 total).
//  4. Country-specific override pattern (if present).
//  5. MOD-97 checksum.
//
// Notes:
// * Country-specific regexes are *tight* (e.g., numeric-only national part for DE).
// * If a country has no specific regex, the generic pattern governs.
// * Strict mode adds a whitelist so "XX..." fails early with ErrUnknownCountry.
type CreditorIDValidator struct {
	countrySpecific     map[string]*regexp.Regexp // ISO-3166-1 → strict pattern
	allowUnknownCountry bool                      // backward-compatible default = true
}

// NewDefaultCreditorIDValidator loads the most common euro-zone patterns.
// Backward-compatible: unknown countries are allowed (generic pattern + checksum).
func NewDefaultCreditorIDValidator() *CreditorIDValidator {
	return &CreditorIDValidator{
		countrySpecific: map[string]*regexp.Regexp{
			// Germany (DE): 18 chars, digits in national part
			"DE": regexp.MustCompile(`^DE[0-9]{2}[A-Z0-9]{3}[0-9]{11}$`),
			// France (FR): 13 chars, 6-digit NNE
			"FR": regexp.MustCompile(`^FR[0-9]{2}[A-Z0-9]{3}[0-9]{6}$`),
			// Spain (ES): 16 chars
			"ES": regexp.MustCompile(`^ES[0-9]{2}[A-Z0-9]{3}[A-Z0-9]{9}$`),
			// Portugal (PT): 13 chars
			"PT": regexp.MustCompile(`^PT[0-9]{2}[A-Z0-9]{3}[0-9]{6}$`),
			// Italy (IT): 23 chars (BIC + padding)
			"IT": regexp.MustCompile(`^IT[0-9]{2}[A-Z0-9]{3}[A-Z0-9]{18}$`),
			// Netherlands (NL): 18 chars
			"NL": regexp.MustCompile(`^NL[0-9]{2}[A-Z0-9]{3}[0-9]{11}$`),
		},
		allowUnknownCountry: true,
	}
}

// NewStrictCreditorIDValidator is opt-in stricter constructor.
// Only listed ISO country codes are allowed; others error w/ ErrUnknownCountry.
func NewStrictCreditorIDValidator(isoCountries []string) *CreditorIDValidator {
	v := NewDefaultCreditorIDValidator()
	v.allowUnknownCountry = false
	if len(isoCountries) > 0 {
		// Build regex that allows *only* these countries in the generic fallback.
		v.buildCountryRestrictedGeneric(isoCountries)
	}
	return v
}

// buildCountryRestrictedGeneric rebuilds genericRegex (package-level) to inline
// an alternation of allowed country codes. Safe because users rarely change after init.
func (v *CreditorIDValidator) buildCountryRestrictedGeneric(countries []string) {
	// Defensive copy + sanitize
	up := make([]string, 0, len(countries))
	for _, c := range countries {
		c = strings.ToUpper(strings.TrimSpace(c))
		if len(c) == 2 && isoCountryRegex.MatchString(c) {
			up = append(up, c)
		}
	}
	if len(up) == 0 {
		return
	}
	// (?:(DE|FR|...))[0-9]{2}...
	pat := fmt.Sprintf(`^(?:%s)[0-9]{2}[A-Z0-9]{3}[A-Z0-9]{1,28}$`, strings.Join(up, "|"))
	genericRegex = regexp.MustCompile(pat)
}

var (
	// Generic EPC rule: country(2) + check(2) + business(3) + national(1-28) ≤ 35
	genericRegex    = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{3}[A-Z0-9]{1,28}$`)
	isoCountryRegex = regexp.MustCompile(`^[A-Z]{2}$`)
)

// Validate performs full CI validation.
func (v *CreditorIDValidator) Validate(raw string) error {
	id := Normalize(raw)

	// Hard sanity on length first (covers empty, crazy long)
	if l := len(id); l < 8 || l > 35 {
		return ErrInvalidLength
	}

	cc := id[:2]
	if !v.allowUnknownCountry {
		if !isoCountryRegex.MatchString(cc) {
			return ErrUnknownCountry
		}
		// If you want to enforce membership in a known list, you could add a set here.
	}

	// Generic structure (guards against short/long/illegal chars quickly)
	if !genericRegex.MatchString(id) {
		return ErrInvalidFormat
	}

	// Country-specific strictness (overrides generic; yields ErrInvalidFormat)
	if pat, ok := v.countrySpecific[cc]; ok && !pat.MatchString(id) {
		return ErrInvalidFormat
	}

	// MOD-97 checksum
	if !mod97(id) {
		return ErrInvalidChecksum
	}
	return nil
}

// -----------------------------------------------------------------------------
// Helpers / Utilities
// -----------------------------------------------------------------------------

// Normalize collapses spaces, trims, upper-cases.
func Normalize(s string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(s), " ", ""))
}

// ComputeCheckDigits returns the 2-digit MOD-97 check for the (country,business,national) payload.
func ComputeCheckDigits(country, business, national string) (string, error) {
	country = strings.ToUpper(strings.TrimSpace(country))
	business = strings.ToUpper(strings.TrimSpace(business))
	national = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(national), " ", ""))

	if !isoCountryRegex.MatchString(country) {
		return "", ErrUnknownCountry
	}
	if !regexp.MustCompile(`^[A-Z0-9]{3}$`).MatchString(business) {
		return "", fmt.Errorf("business segment must be exactly 3 alphanumerics")
	}
	if l := len(national); l < 1 || l > 28 || !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(national) {
		return "", fmt.Errorf("invalid national segment")
	}

	// Build with placeholder check digits "00"
	base := country + "00" + business + national
	return calcCheckDigits(base), nil
}

// BuildCI builds a fully validated CI from segments (computes checksum).
func BuildCI(country, business, national string) (string, error) {
	cd, err := ComputeCheckDigits(country, business, national)
	if err != nil {
		return "", err
	}
	ci := strings.ToUpper(strings.TrimSpace(country)) + cd + strings.ToUpper(strings.TrimSpace(business)) +
		strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(national), " ", ""))
	return ci, nil
}

// calcCheckDigits expects a 4+ char string that already has "00" in positions 3-4.
func calcCheckDigits(base string) string {
	// Move first 4 chars behind the payload (like mod97)
	rearr := base[4:] + base[:4]
	rem := 0
	for _, r := range rearr {
		var val int
		switch {
		case '0' <= r && r <= '9':
			val = int(r - '0')
		case 'A' <= r && r <= 'Z':
			val = int(r-'A') + 10
		default:
			// should never happen because input already validated
			return "00"
		}

		// Append digit(s) of val to the ongoing remainder
		if val < 10 {
			rem = (rem*10 + val) % 97
		} else {
			rem = (rem*100 + val) % 97 // two digits
		}
	}
	return fmt.Sprintf("%02d", 98-rem)
}

// mod97 implements ISO-7064 MOD-97-10, streaming (low alloc).
func mod97(ci string) bool {
	// ci assumed normalized
	rearr := ci[4:] + ci[:4]

	rem := 0
	for _, r := range rearr {
		var val int
		switch {
		case '0' <= r && r <= '9':
			val = int(r - '0')
			rem = (rem*10 + val) % 97
		case 'A' <= r && r <= 'Z':
			val = int(r-'A') + 10 // two digits 10..35
			rem = (rem*100 + val) % 97
		default:
			return false
		}
	}
	return rem == 1
}
