//go:generate go run ../../gen/cmd/gen_registry/main.go

package countrymeta

import "regexp"

// IBANMeta holds metadata for a country's IBAN format.
type IBANMeta struct {
	CountryCode      string
	Length           int
	BBANRegex        *regexp.Regexp
	BankStart        int
	BankEnd          int
	BranchStart      int
	BranchEnd        int
	AccountStart     int
	AccountEnd       int
	Structure        string
}

// Registry is populated at runtime by generated code in registry.go
var Registry map[string]IBANMeta
