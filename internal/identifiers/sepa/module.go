package sepa

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/SamyRai/bank-data/internal/identifiers"
)

var sciRegex = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{3}[A-Z0-9]{1,28}$`)

// Module validates and parses SEPA Creditor IDs.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(input), " ", ""))
}

func (m *Module) DetectCandidate(normalized string) bool {
	if len(normalized) < 7 || len(normalized) > 35 {
		return false
	}
	return normalized[0] >= 'A' && normalized[0] <= 'Z' && normalized[1] >= 'A' && normalized[1] <= 'Z'
}

func (m *Module) Validate(normalized string) error {
	if !sciRegex.MatchString(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "SEPA Creditor ID format is invalid"}
	}
	if !validateChecksum(normalized) {
		return &identifiers.ValidationError{Code: "checksum_failed", Message: "SEPA Creditor ID checksum is invalid"}
	}
	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	return map[string]string{
		"country_code":  normalized[0:2],
		"check_digits":  normalized[2:4],
		"business_code": normalized[4:7],
		"identifier":    normalized[7:],
		"raw":           normalized,
	}, nil
}

func validateChecksum(id string) bool {
	stripped := id[:4] + id[7:]
	rearranged := stripped[4:] + stripped[:4]

	var sb strings.Builder
	for _, r := range rearranged {
		switch {
		case r >= '0' && r <= '9':
			sb.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			sb.WriteString(fmt.Sprintf("%d", int(r-'A'+10)))
		default:
			return false
		}
	}
	return mod97(sb.String()) == 1
}

func mod97(num string) int {
	rem := 0
	for _, r := range num {
		rem = (rem*10 + int(r-'0')) % 97
	}
	return rem
}
