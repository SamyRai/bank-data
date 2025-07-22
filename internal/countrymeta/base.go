package countrymeta

import "regexp"

type Meta struct {
	Country      string
	Name         string
	Length       int
	Regex        *regexp.Regexp
	BankStart    int
	BankEnd      int
	AccountStart int
	AccountEnd   int
}

// Registry is populated at runtime by generated code in registry.go
var Registry map[string]Meta
