// Package main provides a CLI tool to generate client methods from API schemas.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/murasame29/unifi-go-sdk/internal/clientgen"
)

func main() {
	// CLI flags
	input := flag.String("input", "", "Input directory containing API schema JSON files (required)")
	output := flag.String("output", "", "Output directory for generated client code (required)")

	flag.Parse()

	// Validate required flags
	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  clientgen -input <DIR> -output <DIR>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		fmt.Fprintln(os.Stderr, "  -input   Input directory containing API schema JSON files (required)")
		fmt.Fprintln(os.Stderr, "  -output  Output directory for generated client code (required)")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Example:")
		fmt.Fprintln(os.Stderr, "  clientgen -input ./pkg -output ./pkg")
		os.Exit(1)
	}

	// Create generator with options
	gen, err := clientgen.New(
		clientgen.WithInputDir(*input),
		clientgen.WithOutputDir(*output),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create generator: %v\n", err)
		os.Exit(1)
	}

	// Run generation
	if err := gen.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate client code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Client code generated successfully.\n")
	fmt.Printf("  Input:  %s\n", *input)
	fmt.Printf("  Output: %s\n", *output)
}
