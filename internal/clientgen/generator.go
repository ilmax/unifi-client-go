// Package clientgen generates client methods from API schemas.
package clientgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func sanitizeGoPackageName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return "common"
	}

	repl := strings.NewReplacer("-", "_", " ", "_", ".", "_", "/", "_")
	name = repl.Replace(name)

	reInvalid := regexp.MustCompile(`[^a-z0-9_]+`)
	name = reInvalid.ReplaceAllString(name, "")

	reUnderscores := regexp.MustCompile(`_+`)
	name = reUnderscores.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	if name == "" {
		return "common"
	}
	if name[0] >= '0' && name[0] <= '9' {
		name = "pkg_" + name
	}
	return name
}

// APISchema represents the extracted API schema for client generation.
// This mirrors the structure from typegen package.
type APISchema struct {
	Endpoint    string
	Method      string
	Path        string
	Description string
	PathParams  []Property
	Request     *SchemaObject
	Response    *SchemaObject
	Category    string
}

// SchemaObject represents a JSON schema object.
type SchemaObject struct {
	Name       string
	Properties []Property
}

// Property represents a property in a schema.
type Property struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Enum        []string
	Children    []Property
	IsArray     bool
}

// MethodInfo holds information for generating a client method.
type MethodInfo struct {
	Name           string   // Method name (e.g., ListClients, GetClient)
	Description    string   // Method description
	HTTPMethod     string   // HTTP method (GET, POST, PUT, DELETE)
	Path           string   // API path with placeholders
	PathParams     []string // Path parameter names
	HasRequestBody bool     // Whether the method has a request body
	RequestType    string   // Request type name
	ResponseType   string   // Response type name
	EndpointName   string   // Original endpoint name for type references
}

// Generator generates client code from API schemas.
type Generator struct {
	methodTemplate *template.Template
	clientTemplate *template.Template
	inputDir       string
	outputDir      string
}

// Option is a functional option for configuring the Generator.
type Option func(*Generator)

// WithInputDir sets the input directory for reading API schemas.
func WithInputDir(dir string) Option {
	return func(g *Generator) {
		g.inputDir = dir
	}
}

// WithOutputDir sets the output directory for generated client code.
func WithOutputDir(dir string) Option {
	return func(g *Generator) {
		g.outputDir = dir
	}
}

// templateFuncs returns the template functions used in code generation.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"toCamelCase":  toCamelCase,
		"toPascalCase": toPascalCase,
		"formatPath": func(params []string, path string) string {
			return formatPathWithParams(path, params)
		},
	}
}

// New creates a new Generator instance.
func New(opts ...Option) (*Generator, error) {
	funcs := templateFuncs()

	methodTmpl, err := template.New("method").Funcs(funcs).Parse(methodTemplateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse method template: %w", err)
	}

	clientTmpl, err := template.New("client").Funcs(funcs).Parse(clientTemplateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client template: %w", err)
	}

	g := &Generator{
		methodTemplate: methodTmpl,
		clientTemplate: clientTmpl,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g, nil
}

// Generate reads API schemas from the input directory and generates client code.
// It outputs files to pkg/{category}/{category}.go.
func (g *Generator) Generate() error {
	if g.inputDir == "" {
		return fmt.Errorf("input directory not specified")
	}
	if g.outputDir == "" {
		return fmt.Errorf("output directory not specified")
	}

	// Read all schema files from input directory
	schemasByCategory, err := g.readSchemas()
	if err != nil {
		return fmt.Errorf("failed to read schemas: %w", err)
	}

	if len(schemasByCategory) == 0 {
		return fmt.Errorf("no schemas found in %s", g.inputDir)
	}

	// Generate client code for each category
	for category, schemas := range schemasByCategory {
		if err := g.generateCategoryClient(category, schemas); err != nil {
			return fmt.Errorf("failed to generate client for category %s: %w", category, err)
		}
	}

	return nil
}

// readSchemas reads all API schema JSON files from the input directory.
// It returns schemas grouped by category.
func (g *Generator) readSchemas() (map[string][]APISchema, error) {
	schemasByCategory := make(map[string][]APISchema)

	// Walk through the input directory to find schema files
	err := filepath.Walk(g.inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), "_schema.json") {
			return nil
		}

		// Read and parse the schema file
		schemas, err := g.readSchemaFile(path)
		if err != nil {
			return fmt.Errorf("failed to read schema file %s: %w", path, err)
		}

		// Group schemas by category
		for _, schema := range schemas {
			category := schema.Category
			if category == "" {
				category = extractCategoryFromPath(path)
			}
			if category == "" {
				category = "common"
			}
			category = sanitizeGoPackageName(category)
			schemasByCategory[category] = append(schemasByCategory[category], schema)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return schemasByCategory, nil
}

// readSchemaFile reads a single schema JSON file and returns the API schemas.
func (g *Generator) readSchemaFile(path string) ([]APISchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Try to parse as array first
	var schemas []APISchema
	if err := json.Unmarshal(data, &schemas); err == nil {
		return schemas, nil
	}

	// Try to parse as single schema
	var schema APISchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	return []APISchema{schema}, nil
}

// generateCategoryClient generates the client code for a category and writes it to a file.
func (g *Generator) generateCategoryClient(category string, schemas []APISchema) error {
	category = sanitizeGoPackageName(category)
	// Generate the client code
	code, err := g.GenerateClientCode(schemas, category, category)
	if err != nil {
		return fmt.Errorf("failed to generate client code: %w", err)
	}

	// Create the output directory
	outputDir := filepath.Join(g.outputDir, category)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write the client file
	outputPath := filepath.Join(outputDir, category+".go")
	if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write client file: %w", err)
	}

	return nil
}

// extractCategoryFromPath extracts the category from a file path.
// Example: pkg/clients/types_schema.json -> clients
func extractCategoryFromPath(path string) string {
	dir := filepath.Dir(path)
	return filepath.Base(dir)
}

// GenerateClientCode generates Go client code from API schemas.
// It generates a complete client file for a category with all methods.
func (g *Generator) GenerateClientCode(schemas []APISchema, category, pkgName string) (string, error) {
	if len(schemas) == 0 {
		return "", fmt.Errorf("no schemas provided")
	}

	var methods []MethodInfo
	needsFmt := false
	for _, schema := range schemas {
		method := g.schemaToMethodInfo(&schema)
		methods = append(methods, method)
		// Check if any method has path parameters (needs fmt.Sprintf)
		if len(method.PathParams) > 0 {
			needsFmt = true
		}
	}

	clientName := toPascalCase(category)

	data := struct {
		PackageName string
		ClientName  string
		Methods     []MethodInfo
		NeedsFmt    bool
	}{
		PackageName: pkgName,
		ClientName:  clientName,
		Methods:     methods,
		NeedsFmt:    needsFmt,
	}

	var sb strings.Builder
	if err := g.clientTemplate.Execute(&sb, data); err != nil {
		return "", fmt.Errorf("failed to execute client template: %w", err)
	}

	return sb.String(), nil
}

// GenerateMethodCode generates Go code for a single client method.
func (g *Generator) GenerateMethodCode(schema *APISchema) (string, error) {
	method := g.schemaToMethodInfo(schema)

	var sb strings.Builder
	if err := g.methodTemplate.Execute(&sb, method); err != nil {
		return "", fmt.Errorf("failed to execute method template: %w", err)
	}

	return sb.String(), nil
}

// schemaToMethodInfo converts an APISchema to MethodInfo.
func (g *Generator) schemaToMethodInfo(schema *APISchema) MethodInfo {
	methodName := generateMethodName(schema.Endpoint, schema.Method)
	endpointName := sanitizeStructName(schema.Endpoint)

	var pathParams []string
	for _, p := range schema.PathParams {
		pathParams = append(pathParams, p.Name)
	}

	hasRequestBody := schema.Request != nil && len(schema.Request.Properties) > 0
	requestType := ""
	if hasRequestBody {
		requestType = endpointName + "Request"
	}

	responseType := endpointName + "Response"

	return MethodInfo{
		Name:           methodName,
		Description:    schema.Description,
		HTTPMethod:     schema.Method,
		Path:           schema.Path,
		PathParams:     pathParams,
		HasRequestBody: hasRequestBody,
		RequestType:    requestType,
		ResponseType:   responseType,
		EndpointName:   endpointName,
	}
}

// generateMethodName generates a method name from endpoint name and HTTP method.
// Rules:
//   - GET (list): ListXxx
//   - GET (single): GetXxx
//   - POST: CreateXxx
//   - PUT: UpdateXxx
//   - DELETE: DeleteXxx
func generateMethodName(endpointName, httpMethod string) string {
	// Clean up the endpoint name
	name := sanitizeStructName(endpointName)

	// Remove common prefixes that might already be in the endpoint name
	prefixes := []string{"List", "Get", "Create", "Update", "Delete", "Execute"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
			break
		}
	}

	// Determine prefix based on HTTP method and endpoint name pattern
	switch strings.ToUpper(httpMethod) {
	case "GET":
		// Check if it's a list endpoint (typically ends with plural or contains "list")
		lowerEndpoint := strings.ToLower(endpointName)
		if strings.Contains(lowerEndpoint, "list") || isListEndpoint(lowerEndpoint) {
			return "List" + name
		}
		return "Get" + name
	case "POST":
		lowerEndpoint := strings.ToLower(endpointName)
		if strings.Contains(lowerEndpoint, "create") {
			return "Create" + name
		}
		// POST can also be used for actions/execute
		if strings.Contains(lowerEndpoint, "execute") || strings.Contains(lowerEndpoint, "action") {
			return "Execute" + name
		}
		return "Create" + name
	case "PUT":
		return "Update" + name
	case "DELETE":
		return "Delete" + name
	default:
		return name
	}
}

// isListEndpoint checks if an endpoint name suggests a list operation.
func isListEndpoint(name string) bool {
	// Common patterns for list endpoints
	listPatterns := []string{
		"clients", "sites", "vouchers", "devices", "networks",
		"users", "guests", "events", "alerts", "logs",
	}

	lowerName := strings.ToLower(name)
	for _, pattern := range listPatterns {
		if strings.HasSuffix(lowerName, pattern) {
			return true
		}
	}
	return false
}

// sanitizeStructName removes spaces and invalid characters from struct names.
func sanitizeStructName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	parts := re.Split(name, -1)
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}
	return result.String()
}

// toPascalCase converts a string to PascalCase.
func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ".", "_")
	parts := strings.Split(s, "_")

	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}

	return result.String()
}

// toCamelCase converts a string to camelCase.
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(string(pascal[0])) + pascal[1:]
}

// formatPathWithParams converts path parameters from {param} format to Go fmt.Sprintf format.
// Example: /api/v1/sites/{siteId}/clients/{clientId} -> /api/v1/sites/%s/clients/%s
func formatPathWithParams(path string, params []string) string {
	result := path
	for _, param := range params {
		placeholder := "{" + param + "}"
		result = strings.ReplaceAll(result, placeholder, "%s")
	}
	return result
}

// generateMethod generates Go code for a single method from an APISchema.
// This is a helper function that wraps schemaToMethodInfo and template execution.
func (g *Generator) generateMethod(schema *APISchema) (string, error) {
	method := g.schemaToMethodInfo(schema)

	var sb strings.Builder
	if err := g.methodTemplate.Execute(&sb, method); err != nil {
		return "", fmt.Errorf("failed to generate method %s: %w", method.Name, err)
	}

	return sb.String(), nil
}

// GetInputDir returns the configured input directory.
func (g *Generator) GetInputDir() string {
	return g.inputDir
}

// GetOutputDir returns the configured output directory.
func (g *Generator) GetOutputDir() string {
	return g.outputDir
}
