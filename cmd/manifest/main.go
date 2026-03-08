package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SamyRai/bank-data/internal/manifest"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validate_manifest <path/to/manifest.json>")
		os.Exit(1)
	}
	manifestPath, err := sanitizeRelativePath(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid manifest path: %v\n", err)
		os.Exit(1)
	}

	// #nosec G304 -- path is sanitized to a repository-relative path before read.
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

func sanitizeRelativePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("path is empty")
	}

	clean := filepath.Clean(trimmed)
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("absolute paths are not allowed")
	}
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path must stay within repository")
	}

	return clean, nil
}
