package iban

import (
	"strings"

	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/log"
)

// parser implements the Parser interface for IBAN parsing.
type parser struct{}

// NewParser returns a new IBAN Parser implementing the Parser interface.
func NewParser() Parser {
	return &parser{}
}

// Parse extracts IBANInfo from the given IBAN string. Returns IBANInfo and error if parsing fails.
func (p *parser) Parse(ibanStr string) (*IBANInfo, error) {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	log.Debug("Parsing IBAN", log.Fields{"iban": ibanStrNorm, "operation": "parse"})
	if len(ibanStrNorm) < 4 {
		err := *ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN parse failed: wrong length", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		err := *ErrUnsupportedCountry
		err.Value = country
		log.Warn("IBAN parse failed: unsupported country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	if len(ibanStrNorm) != meta.Length {
		err := *ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN parse failed: wrong length for country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	checkDigits := ibanStrNorm[2:4]
	bankCode := ""
	accountNumber := ""
	// BBAN starts at position 4 in IBAN
	bbanOffset := 4
	if meta.BankStart > 0 && meta.BankEnd > 0 && meta.BankEnd > meta.BankStart {
		start := bbanOffset + (meta.BankStart - 1)
		end := bbanOffset + meta.BankEnd
		if end <= len(ibanStrNorm) && start < end {
			bankCode = ibanStrNorm[start:end]
		}
	}
	if meta.AccountStart > 0 && meta.AccountEnd > 0 && meta.AccountEnd > meta.AccountStart {
		start := bbanOffset + (meta.AccountStart - 1)
		end := bbanOffset + meta.AccountEnd
		if end <= len(ibanStrNorm) && start < end {
			accountNumber = ibanStrNorm[start:end]
		}
	}
	log.Info("IBAN parsed successfully", log.Fields{"iban": ibanStrNorm, "operation": "parse"})
	return &IBANInfo{
		CountryCode:   country,
		BankCode:      bankCode,
		BranchCode:    "",
		AccountNumber: accountNumber,
		CheckDigits:   checkDigits,
		Raw:           ibanStrNorm,
	}, nil
}

// detector implements the Detector interface.
type detector struct{}

// NewDetector returns a new IBAN Detector.
func NewDetector() Detector {
	return &detector{}
}

// Detect returns IBAN structure metadata for a given IBAN string.
func (d *detector) Detect(ibanStr string) (*IBANStructure, error) {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	log.Debug("Detecting IBAN structure", log.Fields{"iban": ibanStrNorm, "operation": "detect"})
	if len(ibanStrNorm) < 2 {
		err := *ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN detect failed: wrong length", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		err := *ErrUnsupportedCountry
		err.Value = country
		log.Warn("IBAN detect failed: unsupported country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	structure := buildIBANStructureString(meta)
	return &IBANStructure{
		CountryCode: meta.Country,
		Length:      meta.Length,
		Structure:   structure, // See structureLegend for meaning
	}, nil
}

// buildIBANStructureString builds the structure representation (e.g., CCKKBBBB).
func buildIBANStructureString(meta countrymeta.Meta) string {
	structure := ""
	if len(meta.Country) == 2 {
		structure += "CC"
	}
	structure += "KK"
	for i := 4; i < meta.Length; i++ {
		switch {
		case i >= meta.BankStart && i < meta.BankEnd:
			structure += "B"
		case i >= meta.AccountStart && i < meta.AccountEnd:
			structure += "A"
		default:
			structure += "X"
		}
	}
	return structure
}
