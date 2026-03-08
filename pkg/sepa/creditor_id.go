// Package sepa provides compatibility APIs for SEPA Creditor IDs.
package sepa

import (
	"fmt"
	"regexp"
	"strings"
)

var sciRegex = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{3}[A-Z0-9]{1,28}$`)

type ValidationResult struct {
	Input   string
	Valid   bool
	Code    string
	Message string
}

type CreditorIDValidator struct{}

func NewCreditorIDValidator() *CreditorIDValidator { return &CreditorIDValidator{} }

func (v *CreditorIDValidator) Validate(input string) ValidationResult {
	res := ValidationResult{Input: input, Valid: true}
	norm := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(input), " ", ""))
	if !sciRegex.MatchString(norm) {
		res.Valid = false
		res.Code = "invalid_format"
		res.Message = "invalid SEPA Creditor ID format"
		return res
	}
	if !validateSCIChecksum(norm) {
		res.Valid = false
		res.Code = "checksum_failed"
		res.Message = "invalid SEPA Creditor ID checksum"
		return res
	}
	return res
}

func validateSCIChecksum(id string) bool {
	stripped := id[0:4] + id[7:]
	rearranged := stripped[4:] + stripped[0:4]
	var sb strings.Builder
	for _, r := range rearranged {
		if r >= '0' && r <= '9' {
			sb.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			sb.WriteString(fmt.Sprintf("%d", int(r-'A'+10)))
		} else {
			return false
		}
	}
	return mod97(sb.String()) == 1
}

func mod97(numStr string) int {
	rem := 0
	for _, r := range numStr {
		rem = (rem*10 + int(r-'0')) % 97
	}
	return rem
}
