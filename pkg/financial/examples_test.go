package financial_test

import (
	"fmt"

	"github.com/SamyRai/bank-data/pkg/financial"
)

func ExampleService_validateIBAN() {
	svc := financial.NewService()
	report, err := svc.Validate("DE89370400440532013000", financial.IdentifierIBAN)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true IBAN true
}

func ExampleService_validateBIC() {
	svc := financial.NewService()
	report, err := svc.Validate("DEUTDEFF", financial.IdentifierBIC)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true BIC true
}

func ExampleService_validateSEPACreditorID() {
	svc := financial.NewService()
	report, err := svc.Validate("DE98ZZZ09999999999", financial.IdentifierSEPACreditor)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true SEPA_CREDITOR_ID true
}

func ExampleService_validateLEI() {
	svc := financial.NewService()
	report, err := svc.Validate("529900T8BM49AURSDO78", financial.IdentifierLEI)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true LEI true
}

func ExampleService_validateISIN() {
	svc := financial.NewService()
	report, err := svc.Validate("US0378331005", financial.IdentifierISIN)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true ISIN true
}

func ExampleService_validatePAN() {
	svc := financial.NewService()
	report, err := svc.Validate("4111111111111111", financial.IdentifierPAN)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true PAN true
}

func ExampleService_validateVAT() {
	svc := financial.NewService()
	report, err := svc.Validate("DE136695976", financial.IdentifierVAT)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true VAT true
}

func ExampleService_validateNationalAccountUK() {
	svc := financial.NewService()
	report, err := svc.Validate("20-00-00 55779911", financial.IdentifierNationalAccountUK)
	fmt.Println(err == nil, report.Type, report.Valid)
	// Output: true NATIONAL_ACCOUNT_UK true
}
