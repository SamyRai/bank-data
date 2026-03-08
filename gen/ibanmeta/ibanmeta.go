// Package ibanmeta provides functionality to generate Go code for the IBAN country registry
// from the official SWIFT IBAN registry CSV file.
// It reads the CSV file and generates a Go map containing country-specific IBAN metadata.
package ibanmeta

import (
	"fmt"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var isoCountryCode = regexp.MustCompile(`^[A-Z]{2}$`)

type Meta struct {
	Country      string // ISO code (e.g., "DE")
	Name         string // Country name (e.g., "Germany")
	Length       int
	Regex        string
	BankStart    int
	BankEnd      int
	AccountStart int
	AccountEnd   int
}

// GenerateRegistry reads the TXT and returns the Go code for the registry as []byte, using go/ast for code generation.
// txtPath must be a trusted, relative path to a registry file.
func GenerateRegistry(txtPath string) ([]byte, error) {
	if strings.HasPrefix(txtPath, "/") || !strings.HasSuffix(txtPath, ".txt") {
		return nil, fmt.Errorf("invalid registry file path: %s", txtPath)
	}
	data, err := os.ReadFile(txtPath) // #nosec G304
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("not enough lines in registry file")
	}
	// Parse as transposed table
	var table [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var fields []string
		if strings.Contains(line, "\t") {
			fields = strings.Split(line, "\t")
		} else {
			fields = strings.Split(line, ",")
		}
		table = append(table, fields)
	}
	if len(table) < 2 || len(table[0]) < 2 {
		return nil, fmt.Errorf("table too small or malformed")
	}
	// Find row indices for required fields using robust normalization
	normalize := func(s string) string {
		return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\u200b", ""), "\ufeff", "")))
	}
	rowIdx := map[string]int{}
	for i, row := range table {
		if len(row) == 0 {
			continue
		}
		label := normalize(row[0])
		switch label {
		case "name of country":
			rowIdx["name"] = i
		case "iban prefix country code (iso 3166)":
			rowIdx["iso"] = i
		case "iban structure":
			rowIdx["structure"] = i
		case "iban length":
			rowIdx["length"] = i
		case "bban structure":
			rowIdx["bban"] = i
		case "bank identifier position within the bban":
			rowIdx["bankstart"] = i
		case "bank identifier pattern":
			rowIdx["bankend"] = i
		case "account number position within the bban":
			rowIdx["acctstart"] = i
		case "account number pattern":
			rowIdx["acctend"] = i
		}
	}
	var metas []Meta
	for col := 1; col < len(table[0]); col++ {
		iso := ""
		name := ""
		length := 0
		regex := ""
		bankStart := 0
		bankEnd := 0
		acctStart := 0
		acctEnd := 0
		if idx, ok := rowIdx["iso"]; ok && len(table[idx]) > col {
			iso = strings.TrimSpace(table[idx][col])
		}
		if !isoCountryCode.MatchString(iso) {
			continue // skip non-country columns
		}
		if idx, ok := rowIdx["name"]; ok && len(table[idx]) > col {
			name = strings.TrimSpace(table[idx][col])
		}
		if idx, ok := rowIdx["length"]; ok && len(table[idx]) > col {
			length, _ = strconv.Atoi(strings.TrimSpace(table[idx][col]))
		}
		if idx, ok := rowIdx["structure"]; ok && len(table[idx]) > col {
			// Convert SWIFT/IBAN regex to Go regex
			regex = swiftToGoRegex(strings.TrimSpace(table[idx][col]))
		}
		// Always use BBAN structure parsing for field positions
		bbanIdx, bbanOk := rowIdx["bban"]
		if bbanOk && len(table[bbanIdx]) > col {
			bbanStruct := strings.TrimSpace(table[bbanIdx][col])
			// Use both manual and regex tokenizers for robustness
			tokensManual := parseBBANTokensManual(bbanStruct)
			tokensRegex := parseBBANTokensRegex(bbanStruct)
			if len(tokensManual) != len(tokensRegex) || strings.Join(tokensManual, ",") != strings.Join(tokensRegex, ",") {
				fmt.Printf("TODO: [BBAN PARSE DIFF] %s: manual vs regex tokens mismatch: manual=%v regex=%v\n", iso, tokensManual, tokensRegex)
			}
			// Use manual tokens for field positions
			pos := 1
			for idx, tok := range tokensManual {
				num := 0
				for k := 0; k < len(tok) && tok[k] >= '0' && tok[k] <= '9'; k++ {
					num = num*10 + int(tok[k]-'0')
				}
				if idx == 0 {
					bankStart = pos
					bankEnd = pos + num - 1
				}
				if idx == 1 {
					acctStart = pos
					acctEnd = pos + num - 1
				}
				pos += num
			}
		}
		metas = append(metas, Meta{
			Country:      iso,
			Name:         name,
			Length:       length,
			Regex:        regex,
			BankStart:    bankStart,
			BankEnd:      bankEnd,
			AccountStart: acctStart,
			AccountEnd:   acctEnd,
		})
	}
	// Sort metas by Country code for deterministic, readable output
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Country < metas[j].Country
	})

	// Instead of AST, generate the new Go file as a string
	var buf strings.Builder
	buf.WriteString("package countrymeta\n\n")
	buf.WriteString("import (\n\t\"regexp\"\n)\n\n")
	buf.WriteString("func init() {\n\tRegistry = map[string]Meta{\n")
	for _, m := range metas {
		buf.WriteString("\t\t\"" + m.Country + "\": {\n")
		buf.WriteString("\t\t\tCountry: \"" + m.Country + "\",\n")
		buf.WriteString("\t\t\tName: \"" + strings.ReplaceAll(m.Name, "\"", "\\\"") + "\",\n")
		buf.WriteString(fmt.Sprintf("\t\t\tLength: %d,\n", m.Length))
		if m.Regex != "" {
			buf.WriteString("\t\t\tRegex: regexp.MustCompile(\"" + m.Regex + "\"),\n")
		} else {
			buf.WriteString("\t\t\tRegex: nil,\n")
		}
		buf.WriteString(fmt.Sprintf("\t\t\tBankStart: %d,\n", m.BankStart))
		buf.WriteString(fmt.Sprintf("\t\t\tBankEnd: %d,\n", m.BankEnd))
		buf.WriteString(fmt.Sprintf("\t\t\tAccountStart: %d,\n", m.AccountStart))
		buf.WriteString(fmt.Sprintf("\t\t\tAccountEnd: %d,\n", m.AccountEnd))
		buf.WriteString("\t\t},\n")
	}
	buf.WriteString("\t}\n}\n")

	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to format generated code: %v", err)
	}

	return formatted, nil
}

// swiftToGoRegex converts SWIFT/IBAN regex notation to Go regex.
// This is a best-effort conversion for the most common patterns.
// TODO: Handle all edge cases and document limitations (priority: medium, effort: 2h)
var classMap = map[string]string{
	"n": "[0-9]",
	"a": "[A-Z]",
	"c": "[A-Za-z0-9]",
}

func swiftToGoRegex(swift string) string {
	var result strings.Builder
	result.WriteString("^")
	// e.g. DE2!n8!n10!n → [DE][0-9]{2}[0-9]{8}[0-9]{10}
	prefix := swift[:2]
	result.WriteString(prefix)
	pattern := swift[2:]
	// Find all N!X groups (e.g. 2!n, 8!n, 10!n)
	re := regexp.MustCompile(`(\d+)!([nac])`)
	matches := re.FindAllStringSubmatch(pattern, -1)
	for _, match := range matches {
		count := match[1]
		class := classMap[match[2]]
		result.WriteString(class + "{" + count + "}")
	}
	result.WriteString("$")
	return result.String()
}

// parseBBANTokensRegex uses regex to extract tokens like 8!n, 10!n, etc.
func parseBBANTokensRegex(bban string) []string {
	re := regexp.MustCompile(`\d+![nac]`)
	return re.FindAllString(bban, -1)
}

// parseBBANTokensManual is the original manual parser
func parseBBANTokensManual(bban string) []string {
	var tokens []string
	for i := 0; i < len(bban); {
		j := i
		for j < len(bban) && bban[j] >= '0' && bban[j] <= '9' {
			j++
		}
		if j > i && j+2 <= len(bban) && bban[j] == '!' {
			tok := bban[i : j+2]
			tokens = append(tokens, tok)
			i = j + 2
		} else {
			i++
		}
	}
	return tokens
}
