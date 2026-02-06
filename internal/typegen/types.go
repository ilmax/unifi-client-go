// Package typegen provides functionality to scrape UniFi API documentation and generate Go types.
package typegen

// APIEndpoint represents a discovered API endpoint from navigation.
type APIEndpoint struct {
	Name     string // Display name from navigation
	URL      string // Full URL to the endpoint documentation
	Category string // Category/group from navigation (e.g., "Clients", "Sites")
}

// GenerateResult holds the result of generating types for an endpoint.
type GenerateResult struct {
	Endpoint APIEndpoint
	Schema   *APISchema
	Code     string
	Error    error
}

// APISchema represents the extracted API schema.
type APISchema struct {
	Endpoint    string
	Method      string
	Path        string
	Description string
	PathParams  []Property
	Request     *SchemaObject
	Response    *SchemaObject
	Category    string // Category extracted from path (e.g., "clients", "sites", "vouchers")
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
	Children    []Property // Nested properties for object types
	IsArray     bool       // Whether this is an array of the type
}

// propertyRowInfo holds extracted info from a property row for hierarchy building.
type propertyRowInfo struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Enum        []string
	IsObject    bool
	IsArray     bool
	Depth       int // Indentation depth (0 = top level)
}
