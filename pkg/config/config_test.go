// Package config provides configuration for the UniFi SDK.
package config

import (
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("returns config with default values", func(t *testing.T) {
		cfg := New()

		if cfg.UserAgent != DefaultUserAgent {
			t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, DefaultUserAgent)
		}
		if cfg.Timeout != DefaultTimeout {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, DefaultTimeout)
		}
		if cfg.APIKey != "" {
			t.Errorf("APIKey = %q, want empty string", cfg.APIKey)
		}
		if cfg.BaseURL != "" {
			t.Errorf("BaseURL = %q, want empty string", cfg.BaseURL)
		}
		if cfg.HTTPClient != nil {
			t.Errorf("HTTPClient = %v, want nil", cfg.HTTPClient)
		}
	})
}

func TestConfigAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sets API key",
			input:    "test-api-key",
			expected: "test-api-key",
		},
		{
			name:     "trims whitespace from API key",
			input:    "  test-api-key  ",
			expected: "test-api-key",
		},
		{
			name:     "handles empty API key",
			input:    "",
			expected: "",
		},
		{
			name:     "handles whitespace only",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New()
			opt := ConfigAPIKey(tt.input)
			opt(&cfg)

			if cfg.APIKey != tt.expected {
				t.Errorf("APIKey = %q, want %q", cfg.APIKey, tt.expected)
			}
		})
	}
}

func TestConfigBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sets base URL",
			input:    "https://api.example.com",
			expected: "https://api.example.com",
		},
		{
			name:     "removes trailing slash",
			input:    "https://api.example.com/",
			expected: "https://api.example.com",
		},
		{
			name:     "handles URL without trailing slash",
			input:    "https://api.example.com/v1",
			expected: "https://api.example.com/v1",
		},
		{
			name:     "removes trailing slash from path",
			input:    "https://api.example.com/v1/",
			expected: "https://api.example.com/v1",
		},
		{
			name:     "handles empty URL",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New()
			opt := ConfigBaseURL(tt.input)
			opt(&cfg)

			if cfg.BaseURL != tt.expected {
				t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, tt.expected)
			}
		})
	}
}

func TestConfigTimeout(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected time.Duration
	}{
		{
			name:     "sets custom timeout",
			input:    60 * time.Second,
			expected: 60 * time.Second,
		},
		{
			name:     "sets zero timeout",
			input:    0,
			expected: 0,
		},
		{
			name:     "sets millisecond timeout",
			input:    500 * time.Millisecond,
			expected: 500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New()
			opt := ConfigTimeout(tt.input)
			opt(&cfg)

			if cfg.Timeout != tt.expected {
				t.Errorf("Timeout = %v, want %v", cfg.Timeout, tt.expected)
			}
		})
	}
}

func TestConfigUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sets custom user agent",
			input:    "custom-agent/2.0",
			expected: "custom-agent/2.0",
		},
		{
			name:     "sets empty user agent",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New()
			opt := ConfigUserAgent(tt.input)
			opt(&cfg)

			if cfg.UserAgent != tt.expected {
				t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, tt.expected)
			}
		})
	}
}

func TestConfigHTTPClient(t *testing.T) {
	t.Run("sets custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 120 * time.Second,
		}

		cfg := New()
		opt := ConfigHTTPClient(customClient)
		opt(&cfg)

		if cfg.HTTPClient != customClient {
			t.Errorf("HTTPClient = %v, want %v", cfg.HTTPClient, customClient)
		}
	})

	t.Run("sets nil HTTP client", func(t *testing.T) {
		cfg := New()
		cfg.HTTPClient = &http.Client{} // Set a client first
		opt := ConfigHTTPClient(nil)
		opt(&cfg)

		if cfg.HTTPClient != nil {
			t.Errorf("HTTPClient = %v, want nil", cfg.HTTPClient)
		}
	})
}

func TestInit(t *testing.T) {
	t.Run("applies multiple options", func(t *testing.T) {
		cfg := New()
		opts := []ConfigOption{
			ConfigAPIKey("my-api-key"),
			ConfigBaseURL("https://api.example.com/"),
			ConfigTimeout(45 * time.Second),
			ConfigUserAgent("test-agent/1.0"),
		}

		err := cfg.Init(opts)
		if err != nil {
			t.Fatalf("Init() error = %v, want nil", err)
		}

		if cfg.APIKey != "my-api-key" {
			t.Errorf("APIKey = %q, want %q", cfg.APIKey, "my-api-key")
		}
		if cfg.BaseURL != "https://api.example.com" {
			t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://api.example.com")
		}
		if cfg.Timeout != 45*time.Second {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, 45*time.Second)
		}
		if cfg.UserAgent != "test-agent/1.0" {
			t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, "test-agent/1.0")
		}
	})

	t.Run("creates default HTTP client when not provided", func(t *testing.T) {
		cfg := New()
		cfg.Timeout = 60 * time.Second

		err := cfg.Init(nil)
		if err != nil {
			t.Fatalf("Init() error = %v, want nil", err)
		}

		if cfg.HTTPClient == nil {
			t.Error("HTTPClient = nil, want non-nil")
		}
		if cfg.HTTPClient.Timeout != 60*time.Second {
			t.Errorf("HTTPClient.Timeout = %v, want %v", cfg.HTTPClient.Timeout, 60*time.Second)
		}
	})

	t.Run("preserves custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 120 * time.Second,
		}

		cfg := New()
		opts := []ConfigOption{
			ConfigHTTPClient(customClient),
		}

		err := cfg.Init(opts)
		if err != nil {
			t.Fatalf("Init() error = %v, want nil", err)
		}

		if cfg.HTTPClient != customClient {
			t.Errorf("HTTPClient = %v, want %v", cfg.HTTPClient, customClient)
		}
	})

	t.Run("applies options in order", func(t *testing.T) {
		cfg := New()
		opts := []ConfigOption{
			ConfigAPIKey("first-key"),
			ConfigAPIKey("second-key"),
		}

		err := cfg.Init(opts)
		if err != nil {
			t.Fatalf("Init() error = %v, want nil", err)
		}

		if cfg.APIKey != "second-key" {
			t.Errorf("APIKey = %q, want %q (last option should win)", cfg.APIKey, "second-key")
		}
	})

	t.Run("handles empty options slice", func(t *testing.T) {
		cfg := New()
		err := cfg.Init([]ConfigOption{})
		if err != nil {
			t.Fatalf("Init() error = %v, want nil", err)
		}

		// Should still create default HTTP client
		if cfg.HTTPClient == nil {
			t.Error("HTTPClient = nil, want non-nil")
		}
	})
}
