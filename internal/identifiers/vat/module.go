package vat

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/SamyRai/bank-data/internal/identifiers"
)

var (
	frRegex = regexp.MustCompile(`^FR[0-9A-Z]{2}[0-9]{9}$`)
	nlRegex = regexp.MustCompile(`^NL[0-9]{9}B[0-9]{2}$`)
	deRegex = regexp.MustCompile(`^DE[0-9]{9}$`)
	itRegex = regexp.MustCompile(`^IT[0-9]{11}$`)
	esNIF   = regexp.MustCompile(`^ES[0-9]{8}[A-Z]$`)
)

// Module validates and parses VAT IDs (initial EU subset).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	n := strings.ToUpper(strings.TrimSpace(input))
	replacer := strings.NewReplacer(" ", "", "-", "", ".", "", "_", "")
	return replacer.Replace(n)
}

func (m *Module) DetectCandidate(normalized string) bool {
	return len(normalized) >= 4 && normalized[0] >= 'A' && normalized[0] <= 'Z' && normalized[1] >= 'A' && normalized[1] <= 'Z'
}

func (m *Module) Validate(normalized string) error {
	if !m.DetectCandidate(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "VAT must start with country prefix"}
	}
	country := normalized[:2]
	body := normalized[2:]

	switch country {
	case "DE":
		if !deRegex.MatchString(normalized) {
			return &identifiers.ValidationError{Code: "invalid_format", Message: "DE VAT format must be DE + 9 digits"}
		}
		if !validateDE(body) {
			return &identifiers.ValidationError{Code: "checksum_failed", Message: "DE VAT checksum is invalid"}
		}
	case "FR":
		if !frRegex.MatchString(normalized) {
			return &identifiers.ValidationError{Code: "invalid_format", Message: "FR VAT format is invalid"}
		}
		if !validateFR(body) {
			return &identifiers.ValidationError{Code: "checksum_failed", Message: "FR VAT checksum is invalid"}
		}
	case "NL":
		if !nlRegex.MatchString(normalized) {
			return &identifiers.ValidationError{Code: "invalid_format", Message: "NL VAT format is invalid"}
		}
		if !validateNL(body[:9]) {
			return &identifiers.ValidationError{Code: "checksum_failed", Message: "NL VAT checksum is invalid"}
		}
	case "IT":
		if !itRegex.MatchString(normalized) {
			return &identifiers.ValidationError{Code: "invalid_format", Message: "IT VAT format must be IT + 11 digits"}
		}
		if !validateIT(body) {
			return &identifiers.ValidationError{Code: "checksum_failed", Message: "IT VAT checksum is invalid"}
		}
	case "ES":
		if !esNIF.MatchString(normalized) {
			return &identifiers.ValidationError{Code: "invalid_format", Message: "ES VAT currently supports NIF form ES + 8 digits + letter"}
		}
		if !validateESNIF(body) {
			return &identifiers.ValidationError{Code: "checksum_failed", Message: "ES VAT checksum is invalid"}
		}
	default:
		return &identifiers.ValidationError{Code: "unsupported_country", Message: fmt.Sprintf("VAT country %s is not supported yet", country)}
	}

	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	return map[string]string{
		"country_code": normalized[:2],
		"body":         normalized[2:],
		"raw":          normalized,
	}, nil
}

func validateDE(body string) bool {
	product := 10
	for i := 0; i < 8; i++ {
		d := int(body[i] - '0')
		sum := (d + product) % 10
		if sum == 0 {
			sum = 10
		}
		product = (2 * sum) % 11
	}
	check := 11 - product
	if check == 10 || check == 11 {
		check = 0
	}
	return check == int(body[8]-'0')
}

func validateFR(body string) bool {
	key := body[:2]
	siren := body[2:]
	if key[0] < '0' || key[0] > '9' || key[1] < '0' || key[1] > '9' {
		// Some legal forms use alphanumeric keys; accept format-only for those.
		return true
	}
	var sirenNum int
	for i := 0; i < len(siren); i++ {
		sirenNum = sirenNum*10 + int(siren[i]-'0')
	}
	expected := (12 + 3*(sirenNum%97)) % 97
	keyNum := int(key[0]-'0')*10 + int(key[1]-'0')
	return expected == keyNum
}

func validateNL(firstNine string) bool {
	weights := []int{9, 8, 7, 6, 5, 4, 3, 2, -1}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(firstNine[i]-'0') * weights[i]
	}
	return sum%11 == 0
}

func validateIT(body string) bool {
	sum := 0
	for i := 0; i < 10; i++ {
		d := int(body[i] - '0')
		if i%2 == 0 {
			sum += d
		} else {
			d *= 2
			if d > 9 {
				d -= 9
			}
			sum += d
		}
	}
	check := (10 - (sum % 10)) % 10
	return check == int(body[10]-'0')
}

func validateESNIF(body string) bool {
	number := body[:8]
	letter := body[8]
	var n int
	for i := 0; i < len(number); i++ {
		n = n*10 + int(number[i]-'0')
	}
	letters := "TRWAGMYFPDXBNJZSQVHLCKE"
	return letters[n%23] == letter
}
