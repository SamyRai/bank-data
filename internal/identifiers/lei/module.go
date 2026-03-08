package lei

import (
	"fmt"
	"strings"

	"github.com/SamyRai/bank-data/internal/identifiers"
)

// Module validates and parses LEI values.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(input), " ", ""))
}

func (m *Module) DetectCandidate(normalized string) bool {
	if len(normalized) != 20 {
		return false
	}
	for i := 0; i < len(normalized); i++ {
		c := normalized[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}
	return normalized[18] >= '0' && normalized[18] <= '9' && normalized[19] >= '0' && normalized[19] <= '9'
}

func (m *Module) Validate(normalized string) error {
	if !m.DetectCandidate(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "LEI must be 20 uppercase alphanumeric chars"}
	}
	if !isValidLEIChecksum(normalized) {
		return &identifiers.ValidationError{Code: "checksum_failed", Message: "LEI checksum is invalid"}
	}
	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	return map[string]string{
		"lou_prefix":   normalized[0:4],
		"entity_id":    normalized[4:18],
		"check_digits": normalized[18:20],
		"raw":          normalized,
	}, nil
}

func isValidLEIChecksum(lei string) bool {
	rearranged := lei[4:] + lei[:4]
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
	rem := 0
	for _, r := range sb.String() {
		rem = (rem*10 + int(r-'0')) % 97
	}
	return rem == 1
}
