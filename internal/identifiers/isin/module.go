package isin

import (
	"strings"

	"github.com/SamyRai/bank-data/internal/identifiers"
)

// Module validates and parses ISIN values.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(input), " ", ""))
}

func (m *Module) DetectCandidate(normalized string) bool {
	if len(normalized) != 12 {
		return false
	}
	if normalized[0] < 'A' || normalized[0] > 'Z' || normalized[1] < 'A' || normalized[1] > 'Z' {
		return false
	}
	if normalized[11] < '0' || normalized[11] > '9' {
		return false
	}
	for i := 2; i < 11; i++ {
		c := normalized[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}
	return true
}

func (m *Module) Validate(normalized string) error {
	if !m.DetectCandidate(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "ISIN format is invalid"}
	}
	if !isValidISIN(normalized) {
		return &identifiers.ValidationError{Code: "checksum_failed", Message: "ISIN checksum is invalid"}
	}
	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	return map[string]string{
		"country_code": normalized[0:2],
		"nsin":         normalized[2:11],
		"check_digit":  normalized[11:12],
		"raw":          normalized,
	}, nil
}

func isValidISIN(isin string) bool {
	var digits []int
	for _, r := range isin {
		switch {
		case r >= '0' && r <= '9':
			digits = append(digits, int(r-'0'))
		case r >= 'A' && r <= 'Z':
			v := int(r-'A') + 10
			digits = append(digits, v/10, v%10)
		default:
			return false
		}
	}

	sum := 0
	double := false
	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]
		if double {
			d *= 2
			if d > 9 {
				d = (d / 10) + (d % 10)
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}
