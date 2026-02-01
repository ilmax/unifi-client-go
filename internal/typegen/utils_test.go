package typegen

import "testing"

func TestExtractCategoryFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "clients endpoint with site parameter",
			path:     "/api/v1/sites/{siteId}/clients",
			expected: "clients",
		},
		{
			name:     "single client endpoint",
			path:     "/api/v1/sites/{siteId}/clients/{clientId}",
			expected: "clients",
		},
		{
			name:     "vouchers endpoint",
			path:     "/api/v1/sites/{siteId}/vouchers",
			expected: "vouchers",
		},
		{
			name:     "devices endpoint",
			path:     "/api/v1/sites/{siteId}/devices",
			expected: "devices",
		},
		{
			name:     "single device endpoint",
			path:     "/api/v1/sites/{siteId}/devices/{deviceId}",
			expected: "devices",
		},
		{
			name:     "sites endpoint",
			path:     "/api/v1/sites",
			expected: "sites",
		},
		{
			name:     "single site endpoint",
			path:     "/api/v1/sites/{siteId}",
			expected: "sites",
		},
		{
			name:     "device actions endpoint",
			path:     "/api/v1/sites/{siteId}/devices/{deviceId}/actions",
			expected: "devices",
		},
		{
			name:     "nested resource",
			path:     "/api/v1/sites/{siteId}/clients/{clientId}/sessions",
			expected: "sessions",
		},
		{
			name:     "colon parameter format",
			path:     "/api/v1/sites/:siteId/clients",
			expected: "clients",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "root path",
			path:     "/",
			expected: "",
		},
		{
			name:     "api only",
			path:     "/api",
			expected: "",
		},
		{
			name:     "version only",
			path:     "/api/v1",
			expected: "",
		},
		{
			name:     "path without leading slash",
			path:     "api/v1/sites/{siteId}/clients",
			expected: "clients",
		},
		{
			name:     "path with trailing slash",
			path:     "/api/v1/sites/{siteId}/clients/",
			expected: "clients",
		},
		{
			name:     "batch action endpoint",
			path:     "/api/v1/sites/{siteId}/devices/batch",
			expected: "devices",
		},
		{
			name:     "export action endpoint",
			path:     "/api/v1/sites/{siteId}/clients/export",
			expected: "clients",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCategoryFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("extractCategoryFromPath(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsPathParameter(t *testing.T) {
	tests := []struct {
		segment  string
		expected bool
	}{
		{"{siteId}", true},
		{"{clientId}", true},
		{":id", true},
		{":siteId", true},
		{"clients", false},
		{"sites", false},
		{"api", false},
		{"v1", false},
		{"{}", true},
		{":", false},
	}

	for _, tt := range tests {
		t.Run(tt.segment, func(t *testing.T) {
			result := isPathParameter(tt.segment)
			if result != tt.expected {
				t.Errorf("isPathParameter(%q) = %v, want %v", tt.segment, result, tt.expected)
			}
		})
	}
}

func TestIsNonResourceSegment(t *testing.T) {
	tests := []struct {
		segment  string
		expected bool
	}{
		{"api", true},
		{"API", true},
		{"v1", true},
		{"v2", true},
		{"v1.0", true},
		{"v9.1.120", true},
		{"actions", true},
		{"batch", true},
		{"export", true},
		{"clients", false},
		{"sites", false},
		{"vouchers", false},
		{"devices", false},
	}

	for _, tt := range tests {
		t.Run(tt.segment, func(t *testing.T) {
			result := isNonResourceSegment(tt.segment)
			if result != tt.expected {
				t.Errorf("isNonResourceSegment(%q) = %v, want %v", tt.segment, result, tt.expected)
			}
		})
	}
}
