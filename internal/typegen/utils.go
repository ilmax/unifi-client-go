package typegen

import (
	"regexp"
	"strings"
)

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

	// Handle numeric field names (e.g., "5", "2.4")
	if len(s) > 0 && (s[0] >= '0' && s[0] <= '9') {
		s = "N" + strings.ReplaceAll(s, ".", "_")
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

// toEnumConstName converts enum values like "AUTHORIZE_GUEST_ACCESS" to "AuthorizeGuestAccess".
func toEnumConstName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", " ")
	parts := strings.Fields(s)

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

// toGoType converts a schema type to a Go type.
func toGoType(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))

	if strings.Contains(t, "object") {
		return "json.RawMessage"
	}

	switch t {
	case "string":
		return "string"
	case "integer", "int", "int64", "number":
		return "int64"
	case "float", "float64", "double":
		return "float64"
	case "boolean", "bool":
		return "bool"
	case "array":
		return "[]json.RawMessage"
	default:
		return "json.RawMessage"
	}
}

// toGoTypeWithFieldName determines Go type based on field name patterns.
func toGoTypeWithFieldName(baseType, fieldName string, required bool) string {
	lowerName := strings.ToLower(fieldName)

	// Timestamp fields -> time.Time or *time.Time
	timestampSuffixes := []string{"at", "_at"}
	for _, suffix := range timestampSuffixes {
		if strings.HasSuffix(lowerName, suffix) {
			if required {
				return "time.Time"
			}
			return "*time.Time"
		}
	}

	// Percentage/ratio fields -> float64
	if strings.HasSuffix(lowerName, "pct") ||
		strings.HasSuffix(lowerName, "percent") ||
		strings.HasSuffix(lowerName, "ratio") {
		return "float64"
	}

	// Load average fields -> float64
	if strings.Contains(lowerName, "loadaverage") {
		return "float64"
	}

	// Frequency fields -> float64
	if strings.Contains(lowerName, "frequencyghz") ||
		(strings.Contains(lowerName, "frequency") && strings.Contains(lowerName, "ghz")) {
		return "float64"
	}

	return baseType
}

// endpointToFilename converts an endpoint name to a valid filename.
func endpointToFilename(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	parts := re.Split(name, -1)

	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, strings.ToLower(part))
		}
	}

	return strings.Join(result, "_")
}

// deriveEndpointName extracts endpoint name from URL.
func deriveEndpointName(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if idx := strings.Index(lastPart, "?"); idx != -1 {
			lastPart = lastPart[:idx]
		}
		if lastPart != "" {
			return toStructNameFromEndpoint(lastPart)
		}
	}
	return "API"
}

// toStructNameFromEndpoint converts an endpoint path to a struct name.
func toStructNameFromEndpoint(endpoint string) string {
	name := endpoint
	name = strings.TrimPrefix(name, "Execute ")
	name = strings.TrimPrefix(name, "Get ")
	name = strings.TrimPrefix(name, "List ")
	name = strings.TrimPrefix(name, "Create ")
	name = strings.TrimPrefix(name, "Update ")
	name = strings.TrimPrefix(name, "Delete ")
	name = strings.TrimSuffix(name, " Action")

	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	parts := re.Split(name, -1)
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}
	return result.String()
}

// extractCategoryFromPath extracts the category (resource name) from an API path.
// Examples:
//   - /api/v1/sites/{siteId}/clients -> "clients"
//   - /api/v1/sites/{siteId}/clients/{clientId} -> "clients"
//   - /api/v1/sites/{siteId}/vouchers -> "vouchers"
//   - /api/v1/sites -> "sites"
//   - /api/v1/sites/{siteId}/devices/{deviceId}/actions -> "devices"
func extractCategoryFromPath(path string) string {
	if path == "" {
		return ""
	}

	// Split path into segments
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) == 0 {
		return ""
	}

	// Find the last non-parameter segment that represents a resource collection
	// Parameters are identified by {param} or :param format
	var category string
	for i := len(segments) - 1; i >= 0; i-- {
		segment := segments[i]

		// Skip empty segments
		if segment == "" {
			continue
		}

		// Skip path parameters (e.g., {siteId}, :id)
		if isPathParameter(segment) {
			continue
		}

		// Skip common non-resource segments
		if isNonResourceSegment(segment) {
			continue
		}

		// Found a resource segment
		category = strings.ToLower(segment)
		break
	}

	return category
}

// isPathParameter checks if a path segment is a parameter placeholder.
// Supports formats: {param}, :param
func isPathParameter(segment string) bool {
	// Check for {param} format
	if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
		return true
	}
	// Check for :param format (must have at least one character after colon)
	if strings.HasPrefix(segment, ":") && len(segment) > 1 {
		return true
	}
	return false
}

// isNonResourceSegment checks if a segment is a non-resource path component.
// These are typically API version prefixes or action suffixes.
func isNonResourceSegment(segment string) bool {
	lowerSegment := strings.ToLower(segment)

	// API version prefixes (e.g., "api", "v1", "v2")
	if lowerSegment == "api" {
		return true
	}
	if matched, _ := regexp.MatchString(`^v\d+(\.\d+)*$`, lowerSegment); matched {
		return true
	}

	// Common action suffixes that are not resource names
	actionSuffixes := []string{
		"actions", "action", "execute", "batch", "bulk",
		"export", "import", "search", "count", "stats",
	}
	for _, suffix := range actionSuffixes {
		if lowerSegment == suffix {
			return true
		}
	}

	return false
}
