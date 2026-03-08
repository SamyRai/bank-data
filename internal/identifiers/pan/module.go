package pan

import (
	"strings"

	"github.com/SamyRai/bank-data/internal/identifiers"
)

// Module validates and parses card PAN values.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	replacer := strings.NewReplacer(" ", "", "-", "")
	return replacer.Replace(strings.TrimSpace(input))
}

func (m *Module) DetectCandidate(normalized string) bool {
	if len(normalized) < 12 || len(normalized) > 19 {
		return false
	}
	for i := 0; i < len(normalized); i++ {
		if normalized[i] < '0' || normalized[i] > '9' {
			return false
		}
	}
	return true
}

func (m *Module) Validate(normalized string) error {
	if !m.DetectCandidate(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "PAN must be 12-19 digits"}
	}
	if !luhn(normalized) {
		return &identifiers.ValidationError{Code: "checksum_failed", Message: "PAN checksum is invalid"}
	}
	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	last4 := normalized
	if len(normalized) > 4 {
		last4 = normalized[len(normalized)-4:]
	}
	iin := normalized
	if len(normalized) >= 6 {
		iin = normalized[:6]
	}
	return map[string]string{
		"network": detectNetwork(normalized),
		"iin":     iin,
		"last4":   last4,
		"raw":     normalized,
	}, nil
}

func luhn(num string) bool {
	sum := 0
	double := false
	for i := len(num) - 1; i >= 0; i-- {
		d := int(num[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}

func detectNetwork(pan string) string {
	if strings.HasPrefix(pan, "4") {
		return "visa"
	}
	if len(pan) >= 2 {
		prefix2 := pan[:2]
		if prefix2 == "34" || prefix2 == "37" {
			return "amex"
		}
	}
	if len(pan) >= 2 {
		if pan[:2] >= "51" && pan[:2] <= "55" {
			return "mastercard"
		}
	}
	if len(pan) >= 4 {
		if pan[:4] >= "2221" && pan[:4] <= "2720" {
			return "mastercard"
		}
	}
	if len(pan) >= 2 && pan[:2] >= "56" && pan[:2] <= "69" {
		return "maestro"
	}
	if strings.HasPrefix(pan, "50") {
		return "maestro"
	}
	return "unknown"
}
