package financial

import (
	"strings"

	"github.com/SamyRai/bank-data/pkg/bic"
	"github.com/SamyRai/bank-data/pkg/iban"
	"github.com/SamyRai/bank-data/pkg/sepa"
)

// canonicalService implements the Service interface for v1.
type canonicalService struct {
	ibanService *iban.Service
	bicVal      *bic.Validator
	sepaVal     *sepa.CreditorIDValidator
}

// NewService constructs a standard financial data service using underlying components.
func NewService() Service {
	// Re-use IBAN components
	svc := iban.NewService(nil, nil, nil, nil)
	return &canonicalService{
		ibanService: svc,
		bicVal:      bic.NewValidator(),
		sepaVal:     sepa.NewCreditorIDValidator(),
	}
}

// Detect returns the type of the financial identifier based on structure and lengths.
func (s *canonicalService) Detect(input string) IdentifierType {
	norm := strings.ToUpper(strings.ReplaceAll(input, " ", ""))

	// Fast structural checks without full validation

	// 1. SEPA Creditor ID: 8-35 chars, typically starts with CC + mod97 + 3 chars business code
	// Check if it's a completely valid SEPA ID.
	if len(norm) >= 8 {
		if res := s.sepaVal.Validate(norm); res.Valid {
			return TypeSEPACreditorID
		}
	}

	// 2. IBAN: 15-34 chars, starts with 2 letters and 2 digits
	// First check if it's a completely valid IBAN.
	if len(norm) >= 15 && len(norm) <= 34 {
		if err := s.ibanService.Validate(norm); err == nil {
			return TypeIBAN
		}
	}

	// 3. BIC: 8 or 11 chars, typical pattern AAAA BB CC [DDD]
	if len(norm) == 8 || len(norm) == 11 {
		if res := s.bicVal.Validate(norm); res.Valid {
			return TypeBIC
		}
	}

	// 4. If nothing is completely valid, fallback to structure checks to return the best-effort type
	// This ensures we get the right error for structurally sound but invalid inputs.

	// Structure: SEPA Creditor ID format check (ignoring checksum)
	if len(norm) >= 8 {
		if norm[0] >= 'A' && norm[0] <= 'Z' && norm[1] >= 'A' && norm[1] <= 'Z' &&
			norm[2] >= '0' && norm[2] <= '9' && norm[3] >= '0' && norm[3] <= '9' {

			res := s.sepaVal.Validate(norm)
			if res.Error == nil || strings.Contains(res.Error.Error(), "checksum") {
				return TypeSEPACreditorID
			}
		}
	}

	// Structure: IBAN detection.
	if len(norm) >= 15 && len(norm) <= 34 {
		if norm[0] >= 'A' && norm[0] <= 'Z' && norm[1] >= 'A' && norm[1] <= 'Z' &&
			norm[2] >= '0' && norm[2] <= '9' && norm[3] >= '0' && norm[3] <= '9' {
			if _, err := s.ibanService.Detect(norm); err == nil {
				return TypeIBAN
			}
		}
	}

	// Structure: BIC format check
	if len(norm) == 8 || len(norm) == 11 {
		if _, err := bic.Parse(norm); err == nil {
			return TypeBIC
		}
	}

	return TypeUnknown
}

// Validate provides a consolidated validation report for the input string.
func (s *canonicalService) Validate(input string) ValidationReport {
    norm := strings.ToUpper(strings.ReplaceAll(input, " ", ""))
    idType := s.Detect(norm)

    report := ValidationReport{
        Valid:          false,
        IdentifierType: idType,
        Normalized:     norm,
        Error:          nil,
    }

    switch idType {
    case TypeIBAN:
        err := s.ibanService.Validate(norm)
        report.Valid = err == nil
        report.Error = err
    case TypeBIC:
        res := s.bicVal.Validate(norm)
        report.Valid = res.Valid
        report.Error = res.Error
    case TypeSEPACreditorID:
        res := s.sepaVal.Validate(norm)
        report.Valid = res.Valid
        report.Error = res.Error
    default:
        // Try to run validations to get *why* it failed, if it looks like one
        // Just return the IBAN error by default if it's unknown
        err := s.ibanService.Validate(norm)
        report.Error = err
    }

    return report
}

// Parse extracts components from a recognized financial identifier.
func (s *canonicalService) Parse(input string) (*ParsedIdentifier, error) {
    norm := strings.ToUpper(strings.ReplaceAll(input, " ", ""))
    idType := s.Detect(norm)

    parsed := &ParsedIdentifier{
        IdentifierType: idType,
        Normalized:     norm,
        Components:     make(map[string]string),
    }

    switch idType {
    case TypeIBAN:
        info, err := s.ibanService.Parse(norm)
        if err != nil {
            return nil, err
        }
        parsed.Components["CountryCode"] = info.CountryCode
        parsed.Components["BankCode"] = info.BankCode
        parsed.Components["BranchCode"] = info.BranchCode
        parsed.Components["AccountNumber"] = info.AccountNumber
        parsed.Components["CheckDigits"] = info.CheckDigits
        parsed.Components["BIC"] = info.BIC
        parsed.Components["BankName"] = info.BankName

    case TypeBIC:
        info, err := bic.Parse(norm)
        if err != nil {
            return nil, err
        }
        parsed.Components["Institution"] = info.Institution
        parsed.Components["Country"] = info.Country
        parsed.Components["Location"] = info.Location
        parsed.Components["Branch"] = info.Branch

    case TypeSEPACreditorID:
        // SEPA Creditor ID: AA 00 BBB CCCCCCCCCCC
        // CC, mod-97, Business Code, National ID
        if len(norm) < 8 {
            return nil, s.sepaVal.Validate(norm).Error
        }
        parsed.Components["CountryCode"] = norm[0:2]
        parsed.Components["CheckDigits"] = norm[2:4]
        parsed.Components["BusinessCode"] = norm[4:7]
        parsed.Components["NationalID"] = norm[7:]

    default:
        // Try parsing as IBAN to get an error
        _, err := s.ibanService.Parse(norm)
        return nil, err
    }

    return parsed, nil
}
