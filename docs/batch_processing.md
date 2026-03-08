# Batch Processing

`pkg/financial` provides concurrent batch and streaming validation through a shared core engine.

## Batch Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/SamyRai/bank-data/pkg/financial"
)

func main() {
	svc := financial.NewService()
	inputs := []string{"DE89370400440532013000", "US0378331005", "invalid"}

	reports := svc.ValidateBatch(context.Background(), inputs, "")
	for _, r := range reports {
		fmt.Printf("type=%s valid=%v normalized=%s err=%v\n", r.Type, r.Valid, r.Normalized, r.Error)
	}
}
```

## Streaming Example

```go
seq := func(yield func(string) bool) {
	for _, v := range inputs {
		if !yield(v) {
			return
		}
	}
}

for report := range svc.StreamValidate(context.Background(), seq, "") {
	// consume report
}
```

## Notes

- Batch results preserve input order.
- Stream results are emitted as validations complete.
- Use identifier hints for stricter routing when input type is known ahead of time.
