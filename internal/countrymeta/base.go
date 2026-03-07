package countrymeta

//go:generate go run ../../gen/cmd/gen_registry/main.go ../../datasets/iban-registry.txt registry.go

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
