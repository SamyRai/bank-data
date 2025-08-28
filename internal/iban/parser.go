// Package iban implements IBAN parsing and detection logic.
package iban

import (
	"strings"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/bank"
	"github.com/SamyRai/bank-data/pkg/iban"
)

// parser implements the iban.Parser interface for IBAN parsing.
type parser struct{}

// NewParser returns a new IBAN Parser implementing the Parser interface.
func NewParser() iban.Parser {
	return &parser{}
}

// Parse extracts IBANInfo from the given IBAN string. Returns IBANInfo and error if parsing fails.
func (p *parser) Parse(ibanStr string) (*iban.IBANInfo, error) {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	log.Debug("Parsing IBAN", log.Fields{"iban": ibanStrNorm, "operation": "parse"})
	if len(ibanStrNorm) < 4 {
		err := *iban.ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN parse failed: wrong length", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		err := *iban.ErrUnsupportedCountry
		err.Value = country
		log.Warn("IBAN parse failed: unsupported country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	if len(ibanStrNorm) != meta.Length {
		err := *iban.ErrWrongLength
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
	return &iban.IBANInfo{
		CountryCode:   country,
		BankCode:      bankCode,
		BranchCode:    "",
		AccountNumber: accountNumber,
		CheckDigits:   checkDigits,
		Raw:           ibanStrNorm,
	}, nil
}

// EnrichWithBankInfo enriches a parsed IBANInfo with BankInfo using the provided mapping.
func (p *parser) EnrichWithBankInfo(info *iban.IBANInfo, bicMap bicmap.BankBICMap) (*bank.BankInfo, error) {
	if info == nil {
		return nil, ErrNilIBANInfo
	}
	bankInfo, ok := bicMap.LookupBankInfo(info.CountryCode, info.BankCode)
	if !ok {
		return nil, ErrBankInfoNotFound
	}
	return bankInfo, nil
}

// detector implements the iban.Detector interface.
type detector struct{}

// NewDetector returns a new IBAN Detector.
func NewDetector() iban.Detector {
	return &detector{}
}

// Detect returns IBAN structure metadata for a given IBAN string.
func (d *detector) Detect(ibanStr string) (*iban.IBANStructure, error) {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	log.Debug("Detecting IBAN structure", log.Fields{"iban": ibanStrNorm, "operation": "detect"})
	if len(ibanStrNorm) < 2 {
		err := *iban.ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN detect failed: wrong length", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		err := *iban.ErrUnsupportedCountry
		err.Value = country
		log.Warn("IBAN detect failed: unsupported country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return nil, &err
	}
	structure := buildIBANStructureString(meta)
	return &iban.IBANStructure{
		CountryCode: meta.CountryCode,
		Length:      meta.Length,
		Structure:   structure, // See structureLegend for meaning
	}, nil
}

// structureLegend describes the meaning of each character in the structure string.
// Example: DEkkbbbbbbbbcccccccccc
//
//	D/E = country code, k = check digits, b = bank code, c = account number
//	Structure legend: C=country, K=check digits, B=bank code, A=account number, X=other/unknown
func buildIBANStructureString(meta countrymeta.IBANMeta) string {
	// Start with country code and check digits
	structure := ""
	if len(meta.CountryCode) == 2 {
		structure += "CC" // C = country code
	}
	structure += "KK" // K = check digits
	// Fill the rest with B (bank), A (account), or X (other)
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

// ErrNilIBANInfo is returned if IBANInfo is nil.
var ErrNilIBANInfo = &iban.IBANError{Code: "nil_iban_info", Message: "IBANInfo is nil"}

// ErrBankInfoNotFound is returned if no bank info is found for the IBAN.
var ErrBankInfoNotFound = &iban.IBANError{Code: "bank_info_not_found", Message: "No bank info found for IBAN"}
