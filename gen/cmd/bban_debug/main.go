package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

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

func main() {
	data, err := os.ReadFile("iban-registry.txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	var bbanRow, countryRow int
	for i, line := range lines {
		if strings.HasPrefix(line, "BBAN structure") {
			bbanRow = i
		}
		if strings.HasPrefix(line, "Name of country") {
			countryRow = i
		}
	}
	if bbanRow == 0 || countryRow == 0 {
		panic("Could not find BBAN structure or country row")
	}
	bbanFields := strings.Split(lines[bbanRow], "\t")
	countryFields := strings.Split(lines[countryRow], "\t")
	for col, country := range countryFields {
		country = strings.TrimSpace(country)
		if country == "" {
			continue
		}
		fmt.Printf("\n=== %s BBAN structure ===\n", country)
		bban := ""
		if col < len(bbanFields) {
			bban = bbanFields[col]
		}
		fmt.Println("Raw BBAN struct:", bban)
		// Try regex-based parsing
		tokensRegex := parseBBANTokensRegex(bban)
		fmt.Printf("[Regex] Parsed tokens: %v\n", tokensRegex)
		// Try manual parsing
		tokensManual := parseBBANTokensManual(bban)
		fmt.Printf("[Manual] Parsed tokens: %v\n", tokensManual)
		// Compare results
		if len(tokensRegex) != len(tokensManual) || strings.Join(tokensRegex, ",") != strings.Join(tokensManual, ",") {
			fmt.Println("[DIFF] Tokenization mismatch between regex and manual!")
		}
		// Compute BBAN field positions (manual)
		pos := 1
		for idx, tok := range tokensManual {
			num := 0
			for k := 0; k < len(tok) && tok[k] >= '0' && tok[k] <= '9'; k++ {
				num = num*10 + int(tok[k]-'0')
			}
			fmt.Printf("[Manual] Token %d: %s, start: %d, end: %d\n", idx, tok, pos, pos+num-1)
			pos += num
		}
		// Compute BBAN field positions (regex)
		pos = 1
		for idx, tok := range tokensRegex {
			num := 0
			for k := 0; k < len(tok) && tok[k] >= '0' && tok[k] <= '9'; k++ {
				num = num*10 + int(tok[k]-'0')
			}
			fmt.Printf("[Regex] Token %d: %s, start: %d, end: %d\n", idx, tok, pos, pos+num-1)
			pos += num
		}
		fmt.Println("--- End ---")
	}
}
