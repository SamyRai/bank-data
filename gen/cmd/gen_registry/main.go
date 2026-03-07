package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SamyRai/bank-data/gen/ibanmeta"
)

func main() {
	datasetPath := "datasets/iban-registry.txt"
	if len(os.Args) > 1 {
		datasetPath = os.Args[1]
	}
	// The GenerateRegistry function in the ibanmeta package handles the complex parsing
	code, err := ibanmeta.GenerateRegistry(datasetPath)
	if err != nil {
		log.Fatalf("failed to generate registry: %v", err)
	}

	outputPath := "internal/countrymeta/registry.go"
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	if err := os.WriteFile(outputPath, code, 0644); err != nil {
		log.Fatalf("failed to write registry file: %v", err)
	}

	fmt.Printf("Successfully regenerated %s from datasets/iban-registry.txt\n", outputPath)
}
