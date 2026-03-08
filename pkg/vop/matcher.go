// Package vop provides offline Verification of Payee name matching primitives.
package vop

import (
	"context"
	"errors"
	"strings"
	"unicode"
)

// MatchCategory describes VoP response categories.
type MatchCategory string

const (
	CategoryMatch       MatchCategory = "Match"
	CategoryCloseMatch  MatchCategory = "CloseMatch"
	CategoryNoMatch     MatchCategory = "NoMatch"
	CategoryUnavailable MatchCategory = "Unavailable"
)

// Verifier is the integration seam for bank/provider-backed VoP checks.
// Matcher satisfies this interface with a local deterministic implementation.
type Verifier interface {
	Verify(ctx context.Context, req MatchRequest) (MatchResponse, error)
}

// MatchRequest contains payer-provided and bank-expected names.
type MatchRequest struct {
	SuppliedName string
	ExpectedName string
}

// MatchResponse is the local/offline VoP decision.
type MatchResponse struct {
	Category           MatchCategory
	Score              float64
	NormalizedSupplied string
	NormalizedExpected string
}

// Matcher performs local VoP name matching.
type Matcher struct {
	matchThreshold float64
	closeThreshold float64
}

func NewMatcher() *Matcher {
	return &Matcher{matchThreshold: 0.92, closeThreshold: 0.75}
}

// Verify implements Verifier.
func (m *Matcher) Verify(_ context.Context, req MatchRequest) (MatchResponse, error) {
	return m.Match(req)
}

func (m *Matcher) Match(req MatchRequest) (MatchResponse, error) {
	sup := normalizeName(req.SuppliedName)
	exp := normalizeName(req.ExpectedName)
	if sup == "" || exp == "" {
		return MatchResponse{Category: CategoryUnavailable, Score: 0, NormalizedSupplied: sup, NormalizedExpected: exp}, errors.New("both names must be non-empty")
	}
	if sup == exp {
		return MatchResponse{Category: CategoryMatch, Score: 1, NormalizedSupplied: sup, NormalizedExpected: exp}, nil
	}
	d := levenshtein(sup, exp)
	maxLen := len([]rune(sup))
	if len([]rune(exp)) > maxLen {
		maxLen = len([]rune(exp))
	}
	score := 1 - float64(d)/float64(maxLen)

	category := CategoryNoMatch
	if score >= m.matchThreshold {
		category = CategoryMatch
	} else if score >= m.closeThreshold {
		category = CategoryCloseMatch
	}

	return MatchResponse{
		Category:           category,
		Score:              score,
		NormalizedSupplied: sup,
		NormalizedExpected: exp,
	}, nil
}

func normalizeName(s string) string {
	if s == "" {
		return ""
	}
	lower := strings.ToLower(strings.TrimSpace(s))
	replacer := strings.NewReplacer(
		"ä", "a", "ö", "o", "ü", "u", "ß", "ss",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"á", "a", "à", "a", "â", "a", "ã", "a",
		"í", "i", "ì", "i", "ï", "i", "î", "i",
		"ó", "o", "ò", "o", "ô", "o", "õ", "o",
		"ú", "u", "ù", "u", "û", "u",
	)
	n := replacer.Replace(lower)

	// Strip legal suffixes and punctuation.
	legal := []string{" ltd", " limited", " llc", " inc", " plc", " gmbh", " ag", " s a", " sa", " bv"}
	for _, suffix := range legal {
		n = strings.ReplaceAll(n, suffix, "")
	}

	var b strings.Builder
	b.Grow(len(n))
	space := false
	for _, r := range n {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			space = false
			continue
		}
		if unicode.IsSpace(r) {
			if !space {
				b.WriteRune(' ')
				space = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func levenshtein(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)
	if len(ra) == 0 {
		return len(rb)
	}
	if len(rb) == 0 {
		return len(ra)
	}

	prev := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		curr := make([]int, len(rb)+1)
		curr[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 0
			if ra[i-1] != rb[j-1] {
				cost = 1
			}
			curr[j] = min3(
				curr[j-1]+1,
				prev[j]+1,
				prev[j-1]+cost,
			)
		}
		prev = curr
	}
	return prev[len(rb)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
