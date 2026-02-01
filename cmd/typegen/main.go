// Package main provides a CLI tool to generate Go types from UniFi API documentation.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/murasame29/unifi-client-go/internal/typegen"
)

func main() {
	// Single endpoint mode
	url := flag.String("url", "", "UniFi API documentation URL to scrape (single endpoint mode)")
	output := flag.String("output", "", "Output file path (default: stdout)")
	pkg := flag.String("package", "network", "Package name for generated types")

	// Batch mode - discover and generate all endpoints
	discover := flag.String("discover", "", "Base URL to discover all API endpoints (batch mode)")
	outputDir := flag.String("output-dir", "", "Output directory for generated files (batch mode)")
	workers := flag.Int("workers", 4, "Number of parallel workers (batch mode)")

	// List mode - just list discovered endpoints
	list := flag.Bool("list", false, "List discovered endpoints without generating (use with -discover)")

	flag.Parse()

	gen := typegen.New()
	defer gen.Close()

	// Batch mode: discover and generate all endpoints
	if *discover != "" {
		if *list {
			// Just list endpoints
			endpoints, err := gen.DiscoverEndpoints(*discover)
			if err != nil {
				log.Fatalf("Failed to discover endpoints: %v", err)
			}

			fmt.Printf("Discovered %d endpoints:\n\n", len(endpoints))
			currentCategory := ""
			for _, ep := range endpoints {
				if ep.Category != currentCategory {
					currentCategory = ep.Category
					if currentCategory != "" {
						fmt.Printf("\n[%s]\n", currentCategory)
					}
				}
				fmt.Printf("  - %s\n    %s\n", ep.Name, ep.URL)
			}
			return
		}

		// Generate all endpoints
		if *outputDir == "" {
			log.Fatal("Output directory (-output-dir) is required for batch mode")
		}

		results, err := gen.GenerateAll(*discover, *pkg, *outputDir, *workers)
		if err != nil {
			log.Fatalf("Failed to generate types: %v", err)
		}

		// Print summary
		var succeeded, failed int
		for _, r := range results {
			if r.Error != nil {
				failed++
				fmt.Printf("FAILED: %s - %v\n", r.Endpoint.Name, r.Error)
			} else {
				succeeded++
			}
		}

		fmt.Printf("\nGeneration complete: %d succeeded, %d failed\n", succeeded, failed)
		return
	}

	// Single endpoint mode
	if *url == "" {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  Single endpoint:")
		fmt.Fprintln(os.Stderr, "    typegen -url <URL> [-output <file>] [-package <name>]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  Batch mode (discover and generate all):")
		fmt.Fprintln(os.Stderr, "    typegen -discover <BASE_URL> -output-dir <DIR> [-package <name>] [-workers <N>]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  List endpoints:")
		fmt.Fprintln(os.Stderr, "    typegen -discover <BASE_URL> -list")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  typegen -url https://developer.ui.com/network/v10.0.162/executeconnectedclientaction")
		fmt.Fprintln(os.Stderr, "  typegen -discover https://developer.ui.com/network/v10.0.162 -output-dir ./generated")
		fmt.Fprintln(os.Stderr, "  typegen -discover https://developer.ui.com/network/v10.0.162 -list")
		os.Exit(1)
	}

	result, err := gen.Generate(*url, *pkg)
	if err != nil {
		log.Fatalf("Failed to generate types: %v", err)
	}

	if *output != "" {
		if err := os.WriteFile(*output, []byte(result), 0644); err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Generated types written to %s\n", *output)
	} else {
		fmt.Println(result)
	}
}
