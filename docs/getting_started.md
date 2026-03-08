# Getting Started (5 Minutes)

This guide walks through a minimal integration of `bank-data` using the canonical `pkg/financial` API.

## 1. Install

```sh
go get github.com/SamyRai/bank-data
```

## 2. Validate and Parse

```go
package main

import (
	"fmt"

	"github.com/SamyRai/bank-data/pkg/financial"
)

func main() {
	svc := financial.NewService()

	report, err := svc.Validate("DE89370400440532013000", financial.IdentifierIBAN)
	if err != nil {
		panic(err)
	}
	fmt.Printf("type=%s normalized=%s valid=%v\n", report.Type, report.Normalized, report.Valid)

	parsed, err := svc.Parse("US0378331005", financial.IdentifierISIN)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ISIN country=%s nsin=%s\n", parsed.Fields["country_code"], parsed.Fields["nsin"])
}
```

## 3. Auto-detect Unknown Inputs

```go
typ, err := svc.Detect("529900T8BM49AURSDO78")
// typ == financial.IdentifierLEI
```

## 4. Batch Validation

```go
inputs := []string{"DE89370400440532013000", "US0378331005", "invalid"}
reports := svc.ValidateBatch(context.Background(), inputs, "")
```

## 5. Compatibility Packages

`pkg/iban`, `pkg/bic`, and `pkg/sepa` remain available for legacy integrations.
For all new work, use `pkg/financial`.

## 6. Optional Local VoP Matching

```go
import (
	"context"
	"fmt"

	"github.com/SamyRai/bank-data/pkg/vop"
)

matcher := vop.NewMatcher()
resp, err := matcher.Verify(context.Background(), vop.MatchRequest{
	SuppliedName: "Acme Payments Ltd",
	ExpectedName: "ACME PAYMENTS LIMITED",
})
if err != nil {
	panic(err)
}
fmt.Println(resp.Category, resp.Score)
```
