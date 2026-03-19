package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilmax/unifi-client-go/internal/openapipatch"
)

func main() {
	inPath := flag.String("in", "", "path to the input OpenAPI JSON file")
	outPath := flag.String("out", "", "path to the output OpenAPI JSON file")
	flag.Parse()

	if *inPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: openapi-preprocess -in <input.json> -out <output.json>")
		os.Exit(2)
	}

	in, err := os.ReadFile(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input spec: %v\n", err)
		os.Exit(1)
	}

	out, err := openapipatch.RewriteDiscriminatorUnions(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rewrite discriminator unions: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outPath, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write output spec: %v\n", err)
		os.Exit(1)
	}
}
