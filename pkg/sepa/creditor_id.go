// Package sepa provides validation for SEPA-specific identifiers like Creditor ID.
package sepa

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/SamyRai/bank-data/pkg/validation"
)

// sciRegex matches the structural format of a SEPA Creditor ID.
var sciRegex = regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{3}[A-Z0-9]{1,28}$`)

// CreditorIDValidator implements validation.Validator for SEPA Creditor IDs.
type CreditorIDValidator struct{}

// NewCreditorIDValidator returns a new validator for SEPA Creditor IDs.
func NewCreditorIDValidator() *CreditorIDValidator {
	return &CreditorIDValidator{}
}

// Validate checks if the input is a valid SEPA Creditor ID.
func (v *CreditorIDValidator) Validate(input string) validation.ValidationResult {
	res := validation.ValidationResult{
		Input: input,
		Valid: true,
	}

	norm := strings.ToUpper(strings.ReplaceAll(input, " ", ""))
	if !sciRegex.MatchString(norm) {
		res.Valid = false
		res.Error = fmt.Errorf("invalid SEPA Creditor ID format")
		return res
	}

	// MOD-97 check
	if !validateSCIChecksum(norm) {
		res.Valid = false
		res.Error = fmt.Errorf("invalid SEPA Creditor ID checksum")
		return res
	}

	return res
}

// validateSCIChecksum implements ISO 7064 MOD-97-10 for SEPA Creditor ID.
func validateSCIChecksum(id string) bool {
	// 1. Remove positions 5-7 (Business Code)
	// Actually, the standard says to ignore them for checksum?
	// "the Creditor Business Code (positions 5 to 7) is omitted for the purpose of the check digit calculation"
	stripped := id[0:4] + id[7:]

	// 2. Rearrange: move country code and check digits to the end
	rearranged := stripped[4:] + stripped[0:4]

	// 3. Convert letters to digits
	var sb strings.Builder
	for _, r := range rearranged {
		if r >= '0' && r <= '9' {
			sb.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			sb.WriteString(fmt.Sprintf("%d", int(r-'A'+10)))
		}
	}

	// 4. Calculate MOD-97
	return mod97(sb.String()) == 1
}

func mod97(numStr string) int {
	rem := 0
	for _, r := range numStr {
		rem = (rem*10 + int(r-'0')) % 97
	}
	return rem
}
