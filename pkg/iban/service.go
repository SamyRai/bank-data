// Package iban provides compatibility APIs for IBAN-only workflows.
package iban

import (
	"context"
	"iter"

	corevalidation "github.com/SamyRai/bank-data/internal/core/validation"
	"github.com/SamyRai/bank-data/pkg/bank"
)

// IBANInfo holds parsed IBAN details.
type IBANInfo struct {
	CountryCode   string
	BankCode      string
	BranchCode    string
	AccountNumber string
	CheckDigits   string
	BIC           string
	BankName      string
	Raw           string
}

// IBANStructure describes country-level IBAN format metadata.
type IBANStructure struct {
	CountryCode string
	Length      int
	Structure   string
}

// Service is a compatibility facade for IBAN operations.
type Service struct {
	Validator  Validator
	Parser     Parser
	Detector   Detector
	bankLookup BankLookup
}

// NewService creates a Service. The fourth parameter is now bank enrichment lookup.
func NewService(v Validator, p Parser, d Detector, bankLookup BankLookup) *Service {
	if v == nil {
		v = NewValidator()
	}
	if p == nil {
		p = NewParser()
	}
	if d == nil {
		d = NewDetector()
	}
	return &Service{Validator: v, Parser: p, Detector: d, bankLookup: bankLookup}
}

func (s *Service) Validate(ibanStr string) error { return s.Validator.Validate(ibanStr) }

// ValidateBatch validates many IBANs concurrently.
func (s *Service) ValidateBatch(ctx context.Context, inputs []string) []error {
	engine := corevalidation.NewEngine[string, error](validatorAdapter{s.Validator}, 0)
	return engine.ValidateBatch(ctx, inputs)
}

// StreamValidate validates IBANs from an iterator.
func (s *Service) StreamValidate(ctx context.Context, inputs iter.Seq[string]) iter.Seq[error] {
	engine := corevalidation.NewEngine[string, error](validatorAdapter{s.Validator}, 0)
	return engine.StreamValidate(ctx, inputs)
}

func (s *Service) Parse(ibanStr string) (*IBANInfo, error) {
	info, err := s.Parser.Parse(ibanStr)
	if err != nil {
		return nil, err
	}
	if s.bankLookup != nil {
		if bi, ok := s.bankLookup.LookupBank(info.CountryCode, info.BankCode); ok {
			info.BIC = bi.BIC
			info.BankName = bi.BankName
		}
	}
	return info, nil
}

func (s *Service) Detect(ibanStr string) (*IBANStructure, error) { return s.Detector.Detect(ibanStr) }

func (s *Service) LookupBank(country, bankCode string) (*bank.BankInfo, error) {
	if s.bankLookup == nil {
		return nil, ErrBankInfoNotFound
	}
	info, ok := s.bankLookup.LookupBank(country, bankCode)
	if !ok {
		return nil, ErrBankInfoNotFound
	}
	return info, nil
}

// EnrichWithBankInfo attaches bank metadata using a public lookup seam.
func EnrichWithBankInfo(info *IBANInfo, lookup BankLookup) (*bank.BankInfo, error) {
	if info == nil {
		return nil, ErrNilIBANInfo
	}
	if lookup == nil {
		return nil, ErrBankInfoNotFound
	}
	bankInfo, ok := lookup.LookupBank(info.CountryCode, info.BankCode)
	if !ok {
		return nil, ErrBankInfoNotFound
	}
	return bankInfo, nil
}

func (s *Service) Must() *MustService { return &MustService{Service: s} }

type MustService struct{ *Service }

func (s *MustService) Validate(ibanStr string) {
	if err := s.Service.Validate(ibanStr); err != nil {
		panic(err)
	}
}

func (s *MustService) Parse(ibanStr string) *IBANInfo {
	info, err := s.Service.Parse(ibanStr)
	if err != nil {
		panic(err)
	}
	return info
}

func (s *MustService) Detect(ibanStr string) *IBANStructure {
	st, err := s.Service.Detect(ibanStr)
	if err != nil {
		panic(err)
	}
	return st
}

type validatorAdapter struct{ Validator }

func (a validatorAdapter) Validate(input string) error {
	return a.Validator.Validate(input)
}
