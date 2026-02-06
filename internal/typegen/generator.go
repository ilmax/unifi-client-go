package typegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Generator scrapes UniFi API documentation and generates Go types.
type Generator struct {
	browser   *rod.Browser
	hashDepth map[string]int
	nextDepth int
}

// New creates a new Generator instance.
func New() *Generator {
	l := launcher.New().Headless(true).NoSandbox(true)
	browser := rod.New().ControlURL(l.MustLaunch()).MustConnect()
	return &Generator{
		browser:   browser,
		hashDepth: make(map[string]int),
		nextDepth: 0,
	}
}

// Close releases browser resources.
func (g *Generator) Close() {
	if g.browser != nil {
		g.browser.MustClose()
	}
}

// Generate scrapes the given URL and generates Go types.
func (g *Generator) Generate(url, pkgName string) (string, error) {
	schema, err := g.scrape(url)
	if err != nil {
		return "", fmt.Errorf("failed to scrape: %w", err)
	}

	return generateCode(schema, pkgName), nil
}

// GenerateAll generates types for all discovered endpoints in parallel.
// When outputDir is specified, files are organized by category:
//   - pkg/{category}/{endpoint_name}.go for each endpoint
//   - pkg/common/types.go for shared types
func (g *Generator) GenerateAll(baseURL, pkgName, outputDir string, workers int) ([]GenerateResult, error) {
	endpoints, err := g.DiscoverEndpoints(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover endpoints: %w", err)
	}

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints discovered")
	}

	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Generate common types file in pkg/common/ directory
	if outputDir != "" {
		commonDir := filepath.Join(outputDir, "common")
		if err := os.MkdirAll(commonDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create common directory: %w", err)
		}
		commonCode := generateCommonTypesCode("common")
		commonPath := filepath.Join(commonDir, "types.go")
		if err := os.WriteFile(commonPath, []byte(commonCode), 0644); err != nil {
			return nil, fmt.Errorf("failed to write common types: %w", err)
		}
		fmt.Printf("Written: %s\n", commonPath)
	}

	if workers <= 0 {
		workers = 4
	}
	if workers > len(endpoints) {
		workers = len(endpoints)
	}

	fmt.Printf("Generating types for %d endpoints with %d workers\n", len(endpoints), workers)

	jobs := make(chan APIEndpoint, len(endpoints))
	results := make(chan GenerateResult, len(endpoints))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			g.worker(workerID, jobs, results, pkgName, outputDir)
		}(i)
	}

	for _, ep := range endpoints {
		jobs <- ep
	}
	close(jobs)

	wg.Wait()
	close(results)

	var allResults []GenerateResult
	for result := range results {
		allResults = append(allResults, result)
	}

	// Write schema JSON files for clientgen.
	// clientgen expects one or more files ending in *_schema.json containing APISchema objects.
	if outputDir != "" {
		schemasByCategory := make(map[string][]APISchema)
		for _, r := range allResults {
			if r.Error != nil || r.Schema == nil {
				continue
			}

			schema := *r.Schema
			category := strings.TrimSpace(schema.Category)
			if category == "" {
				category = extractCategoryFromPath(schema.Path)
			}
			if category == "" {
				category = "common"
			}
			schema.Category = category
			schemasByCategory[category] = append(schemasByCategory[category], schema)
		}

		for category := range schemasByCategory {
			sort.Slice(schemasByCategory[category], func(i, j int) bool {
				return schemasByCategory[category][i].Endpoint < schemasByCategory[category][j].Endpoint
			})
		}

		for category, schemas := range schemasByCategory {
			categoryDir := filepath.Join(outputDir, category)
			if err := os.MkdirAll(categoryDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create schema output directory %s: %w", categoryDir, err)
			}

			data, err := json.MarshalIndent(schemas, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("failed to marshal schemas for category %s: %w", category, err)
			}

			schemaPath := filepath.Join(categoryDir, "types_schema.json")
			if err := os.WriteFile(schemaPath, append(data, '\n'), 0644); err != nil {
				return nil, fmt.Errorf("failed to write schema file %s: %w", schemaPath, err)
			}
			fmt.Printf("Written: %s\n", schemaPath)
		}
	}

	// Print summary of categories
	if outputDir != "" {
		categories := make(map[string]int)
		for _, result := range allResults {
			if result.Error == nil && result.Endpoint.Category != "" {
				categories[result.Endpoint.Category]++
			}
		}
		fmt.Println("\nGeneration summary by category:")
		for cat, count := range categories {
			fmt.Printf("  %s: %d endpoints\n", cat, count)
		}
	}

	return allResults, nil
}

// DiscoverEndpoints scrapes the navigation to find all API endpoint URLs.
func (g *Generator) DiscoverEndpoints(baseURL string) ([]APIEndpoint, error) {
	fmt.Println("Discovering endpoints from:", baseURL)
	page := g.browser.MustPage(baseURL)
	defer page.MustClose()

	page.MustWaitLoad()
	time.Sleep(3 * time.Second)

	var endpoints []APIEndpoint

	g.expandNavigation(page)

	navLinks, err := page.Elements("nav a[href], aside a[href], [class*='Navigation'] a[href], [class*='Sidebar'] a[href]")
	if err != nil {
		return nil, fmt.Errorf("failed to find navigation links: %w", err)
	}

	fmt.Printf("Found %d navigation links\n", len(navLinks))

	seenURLs := make(map[string]bool)
	currentCategory := ""

	for _, link := range navLinks {
		href, err := link.Attribute("href")
		if err != nil || href == nil || *href == "" {
			continue
		}

		url := *href
		if strings.HasPrefix(url, "#") || strings.Contains(url, "javascript:") {
			continue
		}

		if strings.HasPrefix(url, "/") {
			parts := strings.SplitN(baseURL, "/", 4)
			if len(parts) >= 3 {
				url = parts[0] + "//" + parts[2] + url
			}
		}

		if seenURLs[url] {
			continue
		}

		if !isEndpointURL(url) {
			text := strings.TrimSpace(link.MustText())
			if text != "" && !strings.Contains(strings.ToLower(text), "overview") {
				currentCategory = text
			}
			continue
		}

		seenURLs[url] = true

		name := strings.TrimSpace(link.MustText())
		if name == "" {
			name = deriveEndpointName(url)
		}

		endpoints = append(endpoints, APIEndpoint{
			Name:     name,
			URL:      url,
			Category: currentCategory,
		})
	}

	fmt.Printf("Discovered %d API endpoints\n", len(endpoints))
	return endpoints, nil
}

// expandNavigation expands all collapsed navigation sections.
func (g *Generator) expandNavigation(page *rod.Page) {
	maxIterations := 3

	for i := 0; i < maxIterations; i++ {
		expandables, _ := page.Elements("[class*='NavItem'] button, [class*='Expand'], [class*='Toggle'], [aria-expanded='false']")

		clicked := 0
		for _, el := range expandables {
			visible, _ := el.Visible()
			if visible {
				el.Click("left", 1)
				clicked++
				time.Sleep(100 * time.Millisecond)
			}
		}

		if clicked == 0 {
			break
		}

		time.Sleep(300 * time.Millisecond)
	}
}

// worker processes endpoints from the jobs channel.
func (g *Generator) worker(id int, jobs <-chan APIEndpoint, results chan<- GenerateResult, pkgName, outputDir string) {
	l := launcher.New().Headless(true).NoSandbox(true)
	browser := rod.New().ControlURL(l.MustLaunch()).MustConnect()
	defer browser.MustClose()

	workerGen := &Generator{
		browser:   browser,
		hashDepth: make(map[string]int),
		nextDepth: 0,
	}

	for endpoint := range jobs {
		fmt.Printf("[Worker %d] Processing: %s\n", id, endpoint.Name)

		// Generate code with category-based package name
		schema, err := workerGen.scrape(endpoint.URL)
		var code string
		if err == nil {
			// Use category from schema if available, otherwise use default pkgName
			categoryPkgName := pkgName
			if schema.Category != "" {
				categoryPkgName = schema.Category
			}
			code = generateCode(schema, categoryPkgName)
		}

		result := GenerateResult{
			Endpoint: endpoint,
			Schema:   schema,
			Code:     code,
			Error:    err,
		}

		// Update endpoint category from schema
		if err == nil && schema.Category != "" {
			result.Endpoint.Category = schema.Category
		}

		if err == nil && outputDir != "" {
			// Determine output directory based on category
			targetDir := outputDir
			if schema.Category != "" {
				targetDir = filepath.Join(outputDir, schema.Category)
			}

			// Create category directory if it doesn't exist
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				result.Error = fmt.Errorf("failed to create category directory %s: %w", targetDir, err)
			} else {
				filename := endpointToFilename(endpoint.Name) + ".go"
				filePath := filepath.Join(targetDir, filename)
				if writeErr := os.WriteFile(filePath, []byte(code), 0644); writeErr != nil {
					result.Error = fmt.Errorf("failed to write file: %w", writeErr)
				} else {
					fmt.Printf("[Worker %d] Written: %s (category: %s)\n", id, filePath, schema.Category)
				}
			}
		}

		results <- result
	}
}
