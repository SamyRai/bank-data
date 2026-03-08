package financial

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"sort"
	"strings"
	"sync"

	corevalidation "github.com/SamyRai/bank-data/internal/core/validation"
	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/identifiers"
	bicid "github.com/SamyRai/bank-data/internal/identifiers/bic"
	ibanid "github.com/SamyRai/bank-data/internal/identifiers/iban"
	isinid "github.com/SamyRai/bank-data/internal/identifiers/isin"
	leiid "github.com/SamyRai/bank-data/internal/identifiers/lei"
	nationalid "github.com/SamyRai/bank-data/internal/identifiers/nationalaccount"
	panid "github.com/SamyRai/bank-data/internal/identifiers/pan"
	sepaid "github.com/SamyRai/bank-data/internal/identifiers/sepa"
	vatid "github.com/SamyRai/bank-data/internal/identifiers/vat"
	"github.com/SamyRai/bank-data/pkg/bank"
)

// BankEnricher optionally enriches parsed IBAN records with bank metadata.
type BankEnricher interface {
	LookupBank(countryCode, bankCode string) (*bank.BankInfo, bool)
}

// Registry stores modules keyed by IdentifierType.
type Registry struct {
	mu      sync.RWMutex
	modules map[IdentifierType]identifiers.Module
}

func NewRegistry() *Registry {
	return &Registry{modules: make(map[IdentifierType]identifiers.Module)}
}

func (r *Registry) Register(t IdentifierType, mod identifiers.Module) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.modules[t] = mod
}

func (r *Registry) Get(t IdentifierType) (identifiers.Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	mod, ok := r.modules[t]
	return mod, ok
}

// Service is the canonical public entrypoint for financial identifiers.
type Service struct {
	registry      *Registry
	detectionRank map[IdentifierType]int
	bankEnricher  BankEnricher
}

// Option configures a Service.
type Option func(*Service)

func WithBankEnricher(enricher BankEnricher) Option {
	return func(s *Service) { s.bankEnricher = enricher }
}

func WithRegistry(reg *Registry) Option {
	return func(s *Service) {
		if reg != nil {
			s.registry = reg
		}
	}
}

// NewService creates a Service with default modules.
func NewService(opts ...Option) *Service {
	s := &Service{
		registry: NewRegistry(),
		detectionRank: map[IdentifierType]int{
			IdentifierIBAN:         100,
			IdentifierLEI:          95,
			IdentifierISIN:         90,
			IdentifierSEPACreditor: 85,
			IdentifierBIC:          80,
			IdentifierNationalAccountUK: 75,
			IdentifierVAT:          70,
			IdentifierPAN:          60,
		},
	}
	for _, opt := range opts {
		opt(s)
	}

	// Register default modules only if absent.
	s.ensureDefault(IdentifierIBAN, ibanid.New())
	s.ensureDefault(IdentifierBIC, bicid.New())
	s.ensureDefault(IdentifierSEPACreditor, sepaid.New())
	s.ensureDefault(IdentifierLEI, leiid.New())
	s.ensureDefault(IdentifierISIN, isinid.New())
	s.ensureDefault(IdentifierPAN, panid.New())
	s.ensureDefault(IdentifierVAT, vatid.New())
	s.ensureDefault(IdentifierNationalAccountUK, nationalid.New())

	return s
}

func (s *Service) ensureDefault(t IdentifierType, mod identifiers.Module) {
	if _, ok := s.registry.Get(t); ok {
		return
	}
	s.registry.Register(t, mod)
}

// Detect auto-detects the identifier type based on successful validation.
func (s *Service) Detect(input string) (IdentifierType, error) {
	t, _, err := s.detectWithNormalized(input)
	return t, err
}

func (s *Service) detectWithNormalized(input string) (IdentifierType, string, error) {
	s.registry.mu.RLock()
	if len(s.registry.modules) == 0 {
		s.registry.mu.RUnlock()
		return "", "", errors.New("no identifier modules registered")
	}

	type match struct {
		typeID     IdentifierType
		normalized string
		rank       int
		length     int
	}
	matches := make([]match, 0, len(s.registry.modules))

	for t, mod := range s.registry.modules {
		n := mod.Normalize(input)
		if !mod.DetectCandidate(n) {
			continue
		}
		if err := mod.Validate(n); err != nil {
			continue
		}
		matches = append(matches, match{typeID: t, normalized: n, rank: s.detectionRank[t], length: len(n)})
	}
	s.registry.mu.RUnlock()

	if len(matches) == 0 {
		return "", "", &ValidationError{Code: "not_detected", Message: "input does not match any supported identifier"}
	}

	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].rank != matches[j].rank {
			return matches[i].rank > matches[j].rank
		}
		if matches[i].length != matches[j].length {
			return matches[i].length > matches[j].length
		}
		return string(matches[i].typeID) < string(matches[j].typeID)
	})

	return matches[0].typeID, matches[0].normalized, nil
}

// Validate validates input either by explicit hint or auto-detection.
func (s *Service) Validate(input string, hint IdentifierType) (ValidationReport, error) {
	if strings.TrimSpace(input) == "" {
		err := &ValidationError{Type: hint, Code: "empty_input", Message: "input is empty"}
		return ValidationReport{Type: hint, Input: input, Normalized: "", Valid: false, Error: err}, err
	}

	if hint == "" {
		t, n, err := s.detectWithNormalized(input)
		if err != nil {
			ve := toValidationError(t, err)
			return ValidationReport{Type: t, Input: input, Normalized: n, Valid: false, Error: ve}, ve
		}
		return ValidationReport{Type: t, Input: input, Normalized: n, Valid: true, Error: nil}, nil
	}

	mod, ok := s.registry.Get(hint)
	if !ok {
		err := &ValidationError{Type: hint, Code: "unsupported_type", Message: "identifier type is not registered"}
		return ValidationReport{Type: hint, Input: input, Normalized: "", Valid: false, Error: err}, err
	}

	n := mod.Normalize(input)
	if err := mod.Validate(n); err != nil {
		ve := toValidationError(hint, err)
		return ValidationReport{Type: hint, Input: input, Normalized: n, Valid: false, Error: ve}, ve
	}
	return ValidationReport{Type: hint, Input: input, Normalized: n, Valid: true, Error: nil}, nil
}

// Parse validates then parses input using the module for hint/detected type.
func (s *Service) Parse(input string, hint IdentifierType) (ParsedIdentifier, error) {
	report, err := s.Validate(input, hint)
	if err != nil {
		return ParsedIdentifier{}, err
	}

	mod, ok := s.registry.Get(report.Type)
	if !ok {
		return ParsedIdentifier{}, &ValidationError{Type: report.Type, Code: "unsupported_type", Message: "identifier type is not registered"}
	}

	fields, err := mod.Parse(report.Normalized)
	if err != nil {
		return ParsedIdentifier{}, toValidationError(report.Type, err)
	}

	if report.Type == IdentifierIBAN && s.bankEnricher != nil {
		country := fields["country_code"]
		bankCode := fields["bank_code"]
		if bi, ok := s.bankEnricher.LookupBank(country, bankCode); ok {
			fields["bank_name"] = bi.BankName
			fields["bic"] = bi.BIC
		}
	}

	return ParsedIdentifier{Type: report.Type, Normalized: report.Normalized, Fields: fields}, nil
}

// Suggest proposes deterministic correction candidates for invalid inputs.
// Currently implemented for IBAN (adjacent transpositions and single-char substitutions).
func (s *Service) Suggest(input string, hint IdentifierType) ([]Suggestion, error) {
	target := hint
	if target == "" {
		target = IdentifierIBAN
	}
	if target != IdentifierIBAN {
		return nil, &ValidationError{
			Type:    target,
			Code:    "unsupported_type",
			Message: "suggest is currently implemented only for IBAN",
		}
	}

	mod, ok := s.registry.Get(IdentifierIBAN)
	if !ok {
		return nil, &ValidationError{
			Type:    IdentifierIBAN,
			Code:    "unsupported_type",
			Message: "IBAN module is not registered",
		}
	}

	normalized := mod.Normalize(input)
	if err := mod.Validate(normalized); err == nil {
		return []Suggestion{}, nil
	}

	suggestions := make([]Suggestion, 0, 5)
	seen := map[string]struct{}{}
	add := func(candidate, reason string) {
		if _, exists := seen[candidate]; exists {
			return
		}
		if err := mod.Validate(candidate); err != nil {
			return
		}
		seen[candidate] = struct{}{}
		suggestions = append(suggestions, Suggestion{
			Type:      IdentifierIBAN,
			Candidate: candidate,
			Reason:    reason,
		})
	}

	for i := 0; i < len(normalized)-1 && len(suggestions) < 5; i++ {
		if normalized[i] == normalized[i+1] {
			continue
		}
		b := []byte(normalized)
		b[i], b[i+1] = b[i+1], b[i]
		add(string(b), "adjacent_transposition")
	}

	for i := 0; i < len(normalized) && len(suggestions) < 5; i++ {
		c := normalized[i]
		switch {
		case c >= '0' && c <= '9':
			for d := byte('0'); d <= '9' && len(suggestions) < 5; d++ {
				if d == c {
					continue
				}
				b := []byte(normalized)
				b[i] = d
				add(string(b), "single_char_substitution")
			}
		case c >= 'A' && c <= 'Z':
			for d := byte('A'); d <= 'Z' && len(suggestions) < 5; d++ {
				if d == c {
					continue
				}
				b := []byte(normalized)
				b[i] = d
				add(string(b), "single_char_substitution")
			}
		}
	}

	return suggestions, nil
}

// ValidateBatch validates many inputs concurrently.
func (s *Service) ValidateBatch(ctx context.Context, inputs []string, hint IdentifierType) []ValidationReport {
	task := batchValidator{svc: s, hint: hint}
	engine := corevalidation.NewEngine[string, ValidationReport](task, 0)
	return engine.ValidateBatch(ctx, inputs)
}

// StreamValidate validates input stream and yields results as they complete.
func (s *Service) StreamValidate(ctx context.Context, inputs iter.Seq[string], hint IdentifierType) iter.Seq[ValidationReport] {
	task := batchValidator{svc: s, hint: hint}
	engine := corevalidation.NewEngine[string, ValidationReport](task, 0)
	return engine.StreamValidate(ctx, inputs)
}

type batchValidator struct {
	svc  *Service
	hint IdentifierType
}

func (b batchValidator) Validate(input string) ValidationReport {
	rep, err := b.svc.Validate(input, b.hint)
	if err != nil {
		return rep
	}
	return rep
}

func toValidationError(t IdentifierType, err error) *ValidationError {
	if err == nil {
		return nil
	}
	var ve *ValidationError
	if errors.As(err, &ve) {
		if ve.Type == "" {
			ve.Type = t
		}
		return ve
	}
	var ie *identifiers.ValidationError
	if errors.As(err, &ie) {
		return &ValidationError{Type: t, Code: ie.Code, Message: ie.Message}
	}
	return &ValidationError{Type: t, Code: "validation_failed", Message: err.Error()}
}

// MustService panics on API errors and is intended for tooling/fixtures.
type MustService struct{ *Service }

func (s *Service) Must() *MustService { return &MustService{Service: s} }

func (m *MustService) Detect(input string) IdentifierType {
	t, err := m.Service.Detect(input)
	if err != nil {
		panic(err)
	}
	return t
}

func (m *MustService) Validate(input string, hint IdentifierType) ValidationReport {
	report, err := m.Service.Validate(input, hint)
	if err != nil {
		panic(err)
	}
	return report
}

func (m *MustService) Parse(input string, hint IdentifierType) ParsedIdentifier {
	parsed, err := m.Service.Parse(input, hint)
	if err != nil {
		panic(err)
	}
	return parsed
}

// BuildIBANStructure exposes the IBAN symbolic structure for callers that need registry introspection.
func BuildIBANStructure(countryCode string) (string, error) {
	cc := strings.ToUpper(strings.TrimSpace(countryCode))
	meta, ok := countrymeta.Registry[cc]
	if !ok {
		return "", fmt.Errorf("unsupported_country: %s", cc)
	}
	return ibanid.BuildStructure(meta), nil
}
