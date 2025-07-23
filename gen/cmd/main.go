package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "gen_iban_meta":
		runSubcommand("gen_iban_meta", args)
	case "download":
		runSubcommand("download", args)
	case "bban_debug":
		runSubcommand("bban_debug", args)
	case "pipeline":
		runSubcommand("pipeline", args)
	case "gen_bin":
		runSubcommand("gen_bin", args)
	case "help", "--help", "-h":
		help()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		help()
		os.Exit(1)
	}
}

func runSubcommand(name string, args []string) {
	bin := fmt.Sprintf("%s/%s/main.go", "./cmd", name)
	if _, err := os.Stat(bin); err != nil {
		fmt.Fprintf(os.Stderr, "Command '%s' not found.\n", name)
		os.Exit(1)
	}
	// Lazy dispatch: run the subcommand as a separate process
	// In production, build each command as a binary and exec it here
	fmt.Printf("Dispatching to %s with args %v\n", name, args)
	// Example: os.Exec or go run
	// For now, just print usage
}

func help() {
	fmt.Println("Usage: gen <command> [args]")
	fmt.Println("Available commands:")
	fmt.Println("  gen_iban_meta   Generate IBAN meta Go source from CSV")
	fmt.Println("  download        Download datasets")
	fmt.Println("  bban_debug      Debug BBAN parsing")
	fmt.Println("  pipeline        Run pipeline operations")
	fmt.Println("  gen_bin         Generate binary assets from datasets")
	fmt.Println("  help            Show this help message")
}
