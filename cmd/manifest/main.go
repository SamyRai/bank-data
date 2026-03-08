package main

import (
	"fmt"
	"os"

	"github.com/SamyRai/bank-data/internal/manifest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validate_manifest <path/to/manifest.json>")
		os.Exit(1)
	}
	manifestPath := os.Args[1]

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Printf("Error reading manifest: %v\n", err)
		os.Exit(1)
	}

	errs := manifest.ValidateManifest(data, manifest.CalculateChecksum)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		os.Exit(1)
	}

	fmt.Println("Manifest validation successful")
}
