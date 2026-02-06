package clientgen

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if gen == nil {
		t.Fatal("New() returned nil generator")
	}
	if gen.methodTemplate == nil {
		t.Error("methodTemplate is nil")
	}
	if gen.clientTemplate == nil {
		t.Error("clientTemplate is nil")
	}
}

func TestGenerateMethodName(t *testing.T) {
	tests := []struct {
		name         string
		endpointName string
		httpMethod   string
		want         string
	}{
		{
			name:         "GET list endpoint",
			endpointName: "List Clients",
			httpMethod:   "GET",
			want:         "ListClients",
		},
		{
			name:         "GET single endpoint",
			endpointName: "Get Client",
			httpMethod:   "GET",
			want:         "GetClient",
		},
		{
			name:         "POST create endpoint",
			endpointName: "Create Voucher",
			httpMethod:   "POST",
			want:         "CreateVoucher",
		},
		{
			name:         "PUT update endpoint",
			endpointName: "Update Device",
			httpMethod:   "PUT",
			want:         "UpdateDevice",
		},
		{
			name:         "DELETE endpoint",
			endpointName: "Delete Voucher",
			httpMethod:   "DELETE",
			want:         "DeleteVoucher",
		},
		{
			name:         "GET endpoint with plural name",
			endpointName: "Clients",
			httpMethod:   "GET",
			want:         "ListClients",
		},
		{
			name:         "GET endpoint with singular name",
			endpointName: "Client",
			httpMethod:   "GET",
			want:         "GetClient",
		},
		{
			name:         "POST execute action",
			endpointName: "Execute Command",
			httpMethod:   "POST",
			want:         "ExecuteCommand",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateMethodName(tt.endpointName, tt.httpMethod)
			if got != tt.want {
				t.Errorf("generateMethodName(%q, %q) = %q, want %q", tt.endpointName, tt.httpMethod, got, tt.want)
			}
		})
	}
}

func TestSanitizeStructName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple name",
			input: "Client",
			want:  "Client",
		},
		{
			name:  "name with spaces",
			input: "List Clients",
			want:  "ListClients",
		},
		{
			name:  "name with special characters",
			input: "Get-Client_Info",
			want:  "GetClientInfo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeStructName(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeStructName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single word",
			input: "client",
			want:  "Client",
		},
		{
			name:  "underscore separated",
			input: "list_clients",
			want:  "ListClients",
		},
		{
			name:  "hyphen separated",
			input: "list-clients",
			want:  "ListClients",
		},
		{
			name:  "space separated",
			input: "list clients",
			want:  "ListClients",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toPascalCase(tt.input)
			if got != tt.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single word",
			input: "client",
			want:  "client",
		},
		{
			name:  "underscore separated",
			input: "site_id",
			want:  "siteId",
		},
		{
			name:  "hyphen separated",
			input: "client-id",
			want:  "clientId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatPathWithParams(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		params []string
		want   string
	}{
		{
			name:   "no params",
			path:   "/api/v1/sites",
			params: nil,
			want:   "/api/v1/sites",
		},
		{
			name:   "single param",
			path:   "/api/v1/sites/{siteId}/clients",
			params: []string{"siteId"},
			want:   "/api/v1/sites/%s/clients",
		},
		{
			name:   "multiple params",
			path:   "/api/v1/sites/{siteId}/clients/{clientId}",
			params: []string{"siteId", "clientId"},
			want:   "/api/v1/sites/%s/clients/%s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPathWithParams(tt.path, tt.params)
			if got != tt.want {
				t.Errorf("formatPathWithParams(%q, %v) = %q, want %q", tt.path, tt.params, got, tt.want)
			}
		})
	}
}

func TestIsListEndpoint(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "clients endpoint",
			input: "clients",
			want:  true,
		},
		{
			name:  "list clients",
			input: "list clients",
			want:  true,
		},
		{
			name:  "single client",
			input: "client",
			want:  false,
		},
		{
			name:  "get client",
			input: "get client",
			want:  false,
		},
		{
			name:  "sites endpoint",
			input: "sites",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isListEndpoint(tt.input)
			if got != tt.want {
				t.Errorf("isListEndpoint(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateClientCode(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	schemas := []APISchema{
		{
			Endpoint:    "List Clients",
			Method:      "GET",
			Path:        "/api/v1/sites/{siteId}/clients",
			Description: "Retrieves all clients for a site.",
			PathParams: []Property{
				{Name: "siteId", Type: "string"},
			},
			Response: &SchemaObject{
				Name: "ListClientsResponse",
				Properties: []Property{
					{Name: "data", Type: "array"},
				},
			},
			Category: "clients",
		},
		{
			Endpoint:    "Get Client",
			Method:      "GET",
			Path:        "/api/v1/sites/{siteId}/clients/{clientId}",
			Description: "Retrieves a specific client.",
			PathParams: []Property{
				{Name: "siteId", Type: "string"},
				{Name: "clientId", Type: "string"},
			},
			Response: &SchemaObject{
				Name: "GetClientResponse",
				Properties: []Property{
					{Name: "id", Type: "string"},
				},
			},
			Category: "clients",
		},
	}

	code, err := gen.GenerateClientCode(schemas, "clients", "clients")
	if err != nil {
		t.Fatalf("GenerateClientCode() error = %v", err)
	}

	// Verify generated code contains expected elements
	expectedElements := []string{
		"package clients",
		"type Clients struct",
		"func NewClients(client *http.Client) *Clients",
		"func (c *Clients) ListClients(ctx context.Context, siteId string)",
		"func (c *Clients) GetClient(ctx context.Context, siteId string, clientId string)",
		"*ListClientsResponse",
		"*GetClientResponse",
		`fmt.Sprintf("/api/v1/sites/%s/clients"`,
		`fmt.Sprintf("/api/v1/sites/%s/clients/%s"`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated code does not contain expected element: %q\nGenerated code:\n%s", expected, code)
		}
	}
}

func TestGenerateClientCodeWithRequestBody(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	schemas := []APISchema{
		{
			Endpoint:    "Create Voucher",
			Method:      "POST",
			Path:        "/api/v1/sites/{siteId}/vouchers",
			Description: "Creates a new voucher.",
			PathParams: []Property{
				{Name: "siteId", Type: "string"},
			},
			Request: &SchemaObject{
				Name: "CreateVoucherRequest",
				Properties: []Property{
					{Name: "code", Type: "string"},
					{Name: "duration", Type: "integer"},
				},
			},
			Response: &SchemaObject{
				Name: "CreateVoucherResponse",
				Properties: []Property{
					{Name: "id", Type: "string"},
				},
			},
			Category: "vouchers",
		},
	}

	code, err := gen.GenerateClientCode(schemas, "vouchers", "vouchers")
	if err != nil {
		t.Fatalf("GenerateClientCode() error = %v", err)
	}

	// Verify generated code contains request body parameter
	expectedElements := []string{
		"package vouchers",
		"type Vouchers struct",
		"func (c *Vouchers) CreateVoucher(ctx context.Context, siteId string, req *CreateVoucherRequest)",
		"c.client.Post(ctx, path, req, &result)",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated code does not contain expected element: %q\nGenerated code:\n%s", expected, code)
		}
	}
}

func TestGenerateClientCodeNoPathParams(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	schemas := []APISchema{
		{
			Endpoint:    "List Sites",
			Method:      "GET",
			Path:        "/api/v1/sites",
			Description: "Retrieves all sites.",
			PathParams:  nil,
			Response: &SchemaObject{
				Name: "ListSitesResponse",
				Properties: []Property{
					{Name: "data", Type: "array"},
				},
			},
			Category: "sites",
		},
	}

	code, err := gen.GenerateClientCode(schemas, "sites", "sites")
	if err != nil {
		t.Fatalf("GenerateClientCode() error = %v", err)
	}

	// Verify generated code does not import fmt when not needed
	if strings.Contains(code, `"fmt"`) {
		t.Errorf("Generated code should not import fmt when no path params exist\nGenerated code:\n%s", code)
	}

	// Verify path is set directly without fmt.Sprintf
	if !strings.Contains(code, `path := "/api/v1/sites"`) {
		t.Errorf("Generated code should set path directly without fmt.Sprintf\nGenerated code:\n%s", code)
	}
}

func TestSchemaToMethodInfo(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	schema := APISchema{
		Endpoint:    "List Clients",
		Method:      "GET",
		Path:        "/api/v1/sites/{siteId}/clients",
		Description: "Retrieves all clients.",
		PathParams: []Property{
			{Name: "siteId", Type: "string"},
		},
		Request: nil,
		Response: &SchemaObject{
			Name: "ListClientsResponse",
		},
		Category: "clients",
	}

	method := gen.schemaToMethodInfo(schema)

	if method.Name != "ListClients" {
		t.Errorf("method.Name = %q, want %q", method.Name, "ListClients")
	}
	if method.HTTPMethod != "GET" {
		t.Errorf("method.HTTPMethod = %q, want %q", method.HTTPMethod, "GET")
	}
	if method.Path != "/api/v1/sites/{siteId}/clients" {
		t.Errorf("method.Path = %q, want %q", method.Path, "/api/v1/sites/{siteId}/clients")
	}
	if len(method.PathParams) != 1 || method.PathParams[0] != "siteId" {
		t.Errorf("method.PathParams = %v, want [siteId]", method.PathParams)
	}
	if method.HasRequestBody {
		t.Error("method.HasRequestBody should be false")
	}
	if method.ResponseType != "ListClientsResponse" {
		t.Errorf("method.ResponseType = %q, want %q", method.ResponseType, "ListClientsResponse")
	}
}

func TestNewWithOptions(t *testing.T) {
	gen, err := New(
		WithInputDir("/tmp/input"),
		WithOutputDir("/tmp/output"),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if gen == nil {
		t.Fatal("New() returned nil generator")
	}
	if gen.GetInputDir() != "/tmp/input" {
		t.Errorf("GetInputDir() = %q, want %q", gen.GetInputDir(), "/tmp/input")
	}
	if gen.GetOutputDir() != "/tmp/output" {
		t.Errorf("GetOutputDir() = %q, want %q", gen.GetOutputDir(), "/tmp/output")
	}
}

func TestGenerateValidation(t *testing.T) {
	tests := []struct {
		name      string
		inputDir  string
		outputDir string
		wantErr   string
	}{
		{
			name:      "missing input directory",
			inputDir:  "",
			outputDir: "/tmp/output",
			wantErr:   "input directory not specified",
		},
		{
			name:      "missing output directory",
			inputDir:  "/tmp/input",
			outputDir: "",
			wantErr:   "output directory not specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []Option
			if tt.inputDir != "" {
				opts = append(opts, WithInputDir(tt.inputDir))
			}
			if tt.outputDir != "" {
				opts = append(opts, WithOutputDir(tt.outputDir))
			}

			gen, err := New(opts...)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			err = gen.Generate()
			if err == nil {
				t.Fatal("Generate() expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Generate() error = %q, want to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestExtractCategoryFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "clients category",
			path: "pkg/clients/types_schema.json",
			want: "clients",
		},
		{
			name: "sites category",
			path: "/home/user/project/pkg/sites/api_schema.json",
			want: "sites",
		},
		{
			name: "root file",
			path: "schema.json",
			want: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCategoryFromPath(tt.path)
			if got != tt.want {
				t.Errorf("extractCategoryFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestGenerateMethod(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	schema := &APISchema{
		Endpoint:    "List Clients",
		Method:      "GET",
		Path:        "/api/v1/sites/{siteId}/clients",
		Description: "Retrieves all clients for a site.",
		PathParams: []Property{
			{Name: "siteId", Type: "string"},
		},
		Response: &SchemaObject{
			Name: "ListClientsResponse",
			Properties: []Property{
				{Name: "data", Type: "array"},
			},
		},
		Category: "clients",
	}

	code, err := gen.generateMethod(schema)
	if err != nil {
		t.Fatalf("generateMethod() error = %v", err)
	}

	expectedElements := []string{
		"ListClients",
		"ctx context.Context",
		"siteId string",
		"*ListClientsResponse",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(code, expected) {
			t.Errorf("Generated method does not contain expected element: %q\nGenerated code:\n%s", expected, code)
		}
	}
}

func TestGenerateClientCodeEmptySchemas(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = gen.GenerateClientCode([]APISchema{}, "clients", "clients")
	if err == nil {
		t.Fatal("GenerateClientCode() expected error for empty schemas, got nil")
	}
	if !strings.Contains(err.Error(), "no schemas provided") {
		t.Errorf("GenerateClientCode() error = %q, want to contain %q", err.Error(), "no schemas provided")
	}
}
