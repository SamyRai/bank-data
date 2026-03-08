package iban

import (
	"crypto/subtle"
	"strings"

	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/identifiers"
)

// Module validates and parses IBAN values.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(input), " ", ""))
}

func (m *Module) DetectCandidate(normalized string) bool {
	if len(normalized) < 4 || len(normalized) > 34 {
		return false
	}
	return normalized[0] >= 'A' && normalized[0] <= 'Z' && normalized[1] >= 'A' && normalized[1] <= 'Z'
}

func (m *Module) Validate(normalized string) error {
	if len(normalized) < 4 || len(normalized) > 34 {
		return &identifiers.ValidationError{Code: "wrong_length", Message: "IBAN length must be between 4 and 34"}
	}
	for i := 0; i < len(normalized); i++ {
		c := normalized[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z')) {
			return &identifiers.ValidationError{Code: "invalid_chars", Message: "IBAN must be alphanumeric uppercase"}
		}
	}
	country := normalized[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		return &identifiers.ValidationError{Code: "unsupported_country", Message: "country not found in IBAN registry"}
	}
	if len(normalized) != meta.Length {
		return &identifiers.ValidationError{Code: "wrong_length", Message: "IBAN length does not match country format"}
	}
	if meta.Regex != nil && !meta.Regex.MatchString(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "IBAN does not match country structure"}
	}
	if !validateChecksum(normalized) {
		return &identifiers.ValidationError{Code: "checksum_failed", Message: "IBAN checksum is invalid"}
	}
	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	meta := countrymeta.Registry[normalized[:2]]
	info := map[string]string{
		"country_code": normalized[:2],
		"check_digits": normalized[2:4],
		"raw":          normalized,
	}

	if code, ok := extractRange(normalized, meta.BankStart, meta.BankEnd); ok {
		info["bank_code"] = code
	}
	if account, ok := extractRange(normalized, meta.AccountStart, meta.AccountEnd); ok {
		info["account_number"] = account
	}
	info["structure"] = BuildStructure(meta)
	return info, nil
}

// BuildStructure returns the symbolic IBAN structure (CCKK + BBAN symbols).
func BuildStructure(meta countrymeta.Meta) string {
	var b strings.Builder
	b.Grow(meta.Length)
	b.WriteString("CC")
	b.WriteString("KK")

	bbanLen := meta.Length - 4
	for i := 1; i <= bbanLen; i++ {
		switch {
		case inRange(i, meta.BankStart, meta.BankEnd):
			b.WriteByte('B')
		case inRange(i, meta.AccountStart, meta.AccountEnd):
			b.WriteByte('A')
		default:
			b.WriteByte('X')
		}
	}
	return b.String()
}

func inRange(pos, start, end int) bool {
	if start <= 0 || end <= 0 {
		return false
	}
	return pos >= start && pos <= end
}

func extractRange(iban string, start, end int) (string, bool) {
	if !inRange(start, start, end) {
		return "", false
	}
	left := 4 + (start - 1)
	right := 4 + end
	if left < 0 || right > len(iban) || left >= right {
		return "", false
	}
	return iban[left:right], true
}

func validateChecksum(iban string) bool {
	var buf [34]byte
	copy(buf[:], iban[4:])
	copy(buf[len(iban)-4:], iban[:4])

	rem := 0
	for i := 0; i < len(iban); i++ {
		c := buf[i]
		switch {
		case c >= '0' && c <= '9':
			rem = (rem*10 + int(c-'0')) % 97
		case c >= 'A' && c <= 'Z':
			v := int(c-'A') + 10
			rem = (rem*100 + v) % 97
		default:
			return false
		}
	}
	return subtle.ConstantTimeEq(int32(rem), 1) == 1
}
