package iban

import (
	"strings"

	"github.com/SamyRai/bank-data/internal/countrymeta"
	ibanid "github.com/SamyRai/bank-data/internal/identifiers/iban"
)

type parser struct{ module *ibanid.Module }

type detector struct{ module *ibanid.Module }

func NewParser() Parser { return &parser{module: ibanid.New()} }

func (p *parser) Parse(ibanStr string) (*IBANInfo, error) {
	n := p.module.Normalize(ibanStr)
	fields, err := p.module.Parse(n)
	if err != nil {
		return nil, toIBANError(n, err)
	}
	return &IBANInfo{
		CountryCode:   fields["country_code"],
		BankCode:      fields["bank_code"],
		AccountNumber: fields["account_number"],
		CheckDigits:   fields["check_digits"],
		Raw:           fields["raw"],
	}, nil
}

func NewDetector() Detector { return &detector{module: ibanid.New()} }

func (d *detector) Detect(ibanStr string) (*IBANStructure, error) {
	n := d.module.Normalize(ibanStr)
	if len(n) < 2 {
		return nil, &IBANError{Code: ErrCodeWrongLength, Field: "length", Message: "IBAN is too short", Value: ibanStr}
	}
	meta, ok := countrymeta.Registry[n[:2]]
	if !ok {
		return nil, &IBANError{Code: ErrCodeUnsupportedCountry, Field: "country", Message: "IBAN country code is not supported", Value: n[:2]}
	}
	return &IBANStructure{
		CountryCode: meta.Country,
		Length:      meta.Length,
		Structure:   ibanid.BuildStructure(meta),
	}, nil
}

func buildIBANStructureString(meta countrymeta.Meta) string {
	// Backward-compatible helper retained for tests/callers in this package.
	return ibanid.BuildStructure(meta)
}

func validateIBANChecksum(ibanStr string) bool {
	n := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	if len(n) < 4 || len(n) > 34 {
		return false
	}
	var buf [34]byte
	copy(buf[:], n[4:])
	copy(buf[len(n)-4:], n[:4])
	rem := 0
	for i := 0; i < len(n); i++ {
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
	return rem == 1
}
