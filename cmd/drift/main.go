package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SamyRai/bank-data/internal/drift"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: drift_report <old_dataset.csv> <new_dataset.csv> <country_code>")
		os.Exit(1)
	}
	oldPath := os.Args[1]
	newPath := os.Args[2]
	countryCode := os.Args[3]

	oldData, err := drift.ParseCSV(oldPath, countryCode)
	if err != nil {
		fmt.Printf("Error reading old dataset: %v\n", err)
		os.Exit(1)
	}

	newData, err := drift.ParseCSV(newPath, countryCode)
	if err != nil {
		fmt.Printf("Error reading new dataset: %v\n", err)
		os.Exit(1)
	}

	report := drift.CalculateDrift(oldData, newData)

	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("Error generating report JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("drift_report.json", reportJSON, 0600); err != nil {
		fmt.Printf("Error writing report file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Drift report written to drift_report.json")
}
