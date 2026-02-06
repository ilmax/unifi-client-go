// Package http provides HTTP client utilities for the UniFi SDK.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pkgerrors "github.com/ilmax/unifi-client-go/pkg/errors"

	"github.com/ilmax/unifi-client-go/pkg/config"
)

// TestNewClient tests the NewClient function.
func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cfg       config.Config
		wantURL   string
		wantKey   string
		wantAgent string
	}{
		{
			name: "creates client with all config values",
			cfg: config.Config{
				BaseURL:    "https://api.example.com",
				APIKey:     "test-api-key",
				UserAgent:  "test-agent/1.0",
				HTTPClient: &http.Client{Timeout: 10 * time.Second},
			},
			wantURL:   "https://api.example.com",
			wantKey:   "test-api-key",
			wantAgent: "test-agent/1.0",
		},
		{
			name: "creates client with empty values",
			cfg: config.Config{
				BaseURL:    "",
				APIKey:     "",
				UserAgent:  "",
				HTTPClient: &http.Client{},
			},
			wantURL:   "",
			wantKey:   "",
			wantAgent: "",
		},
		{
			name: "creates client with default config",
			cfg: config.Config{
				BaseURL:    "https://unifi.ui.com",
				APIKey:     "my-key",
				UserAgent:  config.DefaultUserAgent,
				HTTPClient: &http.Client{Timeout: config.DefaultTimeout},
			},
			wantURL:   "https://unifi.ui.com",
			wantKey:   "my-key",
			wantAgent: config.DefaultUserAgent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewClient(tt.cfg)

			if client.baseURL != tt.wantURL {
				t.Errorf("baseURL = %q, want %q", client.baseURL, tt.wantURL)
			}
			if client.apiKey != tt.wantKey {
				t.Errorf("apiKey = %q, want %q", client.apiKey, tt.wantKey)
			}
			if client.userAgent != tt.wantAgent {
				t.Errorf("userAgent = %q, want %q", client.userAgent, tt.wantAgent)
			}
			if client.httpClient != tt.cfg.HTTPClient {
				t.Error("httpClient was not set correctly")
			}
		})
	}
}

// TestClient_Get tests the Get method.
func TestClient_Get(t *testing.T) {
	t.Parallel()

	type response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		path           string
		wantResult     *response
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "successful GET request",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response{ID: 1, Name: "test"})
			},
			path:       "/api/test",
			wantResult: &response{ID: 1, Name: "test"},
			wantErr:    false,
		},
		{
			name: "GET request with 404 error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("not found"))
			},
			path:           "/api/notfound",
			wantResult:     nil,
			wantErr:        true,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "GET request with 500 error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			},
			path:           "/api/error",
			wantResult:     nil,
			wantErr:        true,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "GET request with empty response",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			path:       "/api/empty",
			wantResult: nil,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			client := NewClient(config.Config{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				UserAgent:  "test-agent",
				HTTPClient: server.Client(),
			})

			var result response
			var resultPtr *response
			if tt.wantResult != nil {
				resultPtr = &result
			}

			err := client.Get(context.Background(), tt.path, resultPtr)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var apiErr *pkgerrors.APIError
				if errors.As(err, &apiErr) {
					if apiErr.StatusCode != tt.wantStatusCode {
						t.Errorf("status code = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantResult != nil {
				if result.ID != tt.wantResult.ID || result.Name != tt.wantResult.Name {
					t.Errorf("result = %+v, want %+v", result, tt.wantResult)
				}
			}
		})
	}
}

// TestClient_Post tests the Post method.
func TestClient_Post(t *testing.T) {
	t.Parallel()

	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		path           string
		body           interface{}
		wantResult     *response
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "successful POST request",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}
				var req request
				json.NewDecoder(r.Body).Decode(&req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response{ID: 1, Name: req.Name})
			},
			path:       "/api/create",
			body:       request{Name: "new-item"},
			wantResult: &response{ID: 1, Name: "new-item"},
			wantErr:    false,
		},
		{
			name: "POST request with 400 error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("bad request"))
			},
			path:           "/api/create",
			body:           request{Name: ""},
			wantResult:     nil,
			wantErr:        true,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "POST request with nil body",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response{ID: 2, Name: "default"})
			},
			path:       "/api/create",
			body:       nil,
			wantResult: &response{ID: 2, Name: "default"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			client := NewClient(config.Config{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				UserAgent:  "test-agent",
				HTTPClient: server.Client(),
			})

			var result response
			var resultPtr *response
			if tt.wantResult != nil {
				resultPtr = &result
			}

			err := client.Post(context.Background(), tt.path, tt.body, resultPtr)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var apiErr *pkgerrors.APIError
				if errors.As(err, &apiErr) {
					if apiErr.StatusCode != tt.wantStatusCode {
						t.Errorf("status code = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantResult != nil {
				if result.ID != tt.wantResult.ID || result.Name != tt.wantResult.Name {
					t.Errorf("result = %+v, want %+v", result, tt.wantResult)
				}
			}
		})
	}
}

// TestClient_Put tests the Put method.
func TestClient_Put(t *testing.T) {
	t.Parallel()

	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		path           string
		body           interface{}
		wantResult     *response
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "successful PUT request",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("expected PUT method, got %s", r.Method)
				}
				var req request
				json.NewDecoder(r.Body).Decode(&req)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response{ID: 1, Name: req.Name})
			},
			path:       "/api/update/1",
			body:       request{Name: "updated-item"},
			wantResult: &response{ID: 1, Name: "updated-item"},
			wantErr:    false,
		},
		{
			name: "PUT request with 404 error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("not found"))
			},
			path:           "/api/update/999",
			body:           request{Name: "updated"},
			wantResult:     nil,
			wantErr:        true,
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			client := NewClient(config.Config{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				UserAgent:  "test-agent",
				HTTPClient: server.Client(),
			})

			var result response
			var resultPtr *response
			if tt.wantResult != nil {
				resultPtr = &result
			}

			err := client.Put(context.Background(), tt.path, tt.body, resultPtr)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var apiErr *pkgerrors.APIError
				if errors.As(err, &apiErr) {
					if apiErr.StatusCode != tt.wantStatusCode {
						t.Errorf("status code = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantResult != nil {
				if result.ID != tt.wantResult.ID || result.Name != tt.wantResult.Name {
					t.Errorf("result = %+v, want %+v", result, tt.wantResult)
				}
			}
		})
	}
}

// TestClient_Delete tests the Delete method.
func TestClient_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		path           string
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "successful DELETE request",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE method, got %s", r.Method)
				}
				w.WriteHeader(http.StatusNoContent)
			},
			path:    "/api/delete/1",
			wantErr: false,
		},
		{
			name: "DELETE request with 404 error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("not found"))
			},
			path:           "/api/delete/999",
			wantErr:        true,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "DELETE request with 403 error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("forbidden"))
			},
			path:           "/api/delete/protected",
			wantErr:        true,
			wantStatusCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			client := NewClient(config.Config{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				UserAgent:  "test-agent",
				HTTPClient: server.Client(),
			})

			err := client.Delete(context.Background(), tt.path, nil)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var apiErr *pkgerrors.APIError
				if errors.As(err, &apiErr) {
					if apiErr.StatusCode != tt.wantStatusCode {
						t.Errorf("status code = %d, want %d", apiErr.StatusCode, tt.wantStatusCode)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestClient_AuthenticationHeader tests that the authentication header is set correctly.
func TestClient_AuthenticationHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		apiKey      string
		wantAPIKey  string
		wantHeaders map[string]string
	}{
		{
			name:       "API key is set in header",
			apiKey:     "my-secret-api-key",
			wantAPIKey: "my-secret-api-key",
			wantHeaders: map[string]string{
				"X-API-Key":    "my-secret-api-key",
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
		{
			name:       "empty API key",
			apiKey:     "",
			wantAPIKey: "",
			wantHeaders: map[string]string{
				"X-API-Key":    "",
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
		{
			name:       "API key with special characters",
			apiKey:     "key-with-special_chars.123",
			wantAPIKey: "key-with-special_chars.123",
			wantHeaders: map[string]string{
				"X-API-Key": "key-with-special_chars.123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var receivedHeaders http.Header
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedHeaders = r.Header.Clone()
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := NewClient(config.Config{
				BaseURL:    server.URL,
				APIKey:     tt.apiKey,
				UserAgent:  "test-agent",
				HTTPClient: server.Client(),
			})

			err := client.Get(context.Background(), "/test", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for header, wantValue := range tt.wantHeaders {
				gotValue := receivedHeaders.Get(header)
				if gotValue != wantValue {
					t.Errorf("header %q = %q, want %q", header, gotValue, wantValue)
				}
			}
		})
	}
}

// TestClient_UserAgentHeader tests that the User-Agent header is set correctly.
func TestClient_UserAgentHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		userAgent     string
		wantUserAgent string
	}{
		{
			name:          "custom user agent",
			userAgent:     "my-app/2.0",
			wantUserAgent: "my-app/2.0",
		},
		{
			name:          "default user agent",
			userAgent:     config.DefaultUserAgent,
			wantUserAgent: config.DefaultUserAgent,
		},
		{
			name:          "empty user agent",
			userAgent:     "",
			wantUserAgent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var receivedUserAgent string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedUserAgent = r.Header.Get("User-Agent")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := NewClient(config.Config{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				UserAgent:  tt.userAgent,
				HTTPClient: server.Client(),
			})

			err := client.Get(context.Background(), "/test", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if receivedUserAgent != tt.wantUserAgent {
				t.Errorf("User-Agent = %q, want %q", receivedUserAgent, tt.wantUserAgent)
			}
		})
	}
}

// TestClient_JSONParsingError tests error handling for invalid JSON responses.
func TestClient_JSONParsingError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		serverHandler http.HandlerFunc
		wantErrMsg    string
	}{
		{
			name: "invalid JSON response",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("invalid json {"))
			},
			wantErrMsg: "failed to decode response",
		},
		{
			name: "HTML response instead of JSON",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte("<html><body>Error</body></html>"))
			},
			wantErrMsg: "failed to decode response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			client := NewClient(config.Config{
				BaseURL:    server.URL,
				APIKey:     "test-key",
				UserAgent:  "test-agent",
				HTTPClient: server.Client(),
			})

			type response struct {
				ID int `json:"id"`
			}
			var result response

			err := client.Get(context.Background(), "/test", &result)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !containsString(err.Error(), tt.wantErrMsg) {
				t.Errorf("error message = %q, want to contain %q", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

// TestClient_RequestBodyMarshalError tests error handling for invalid request body.
func TestClient_RequestBodyMarshalError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(config.Config{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		UserAgent:  "test-agent",
		HTTPClient: server.Client(),
	})

	// Create an unmarshalable value (channel cannot be marshaled to JSON)
	unmarshalable := make(chan int)

	err := client.Post(context.Background(), "/test", unmarshalable, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !containsString(err.Error(), "failed to marshal request body") {
		t.Errorf("error message = %q, want to contain 'failed to marshal request body'", err.Error())
	}
}

// TestClient_SetBaseURL tests the SetBaseURL method.
func TestClient_SetBaseURL(t *testing.T) {
	t.Parallel()

	client := NewClient(config.Config{
		BaseURL:    "https://original.example.com",
		APIKey:     "test-key",
		UserAgent:  "test-agent",
		HTTPClient: &http.Client{},
	})

	if client.baseURL != "https://original.example.com" {
		t.Errorf("initial baseURL = %q, want %q", client.baseURL, "https://original.example.com")
	}

	client.SetBaseURL("https://new.example.com")

	if client.baseURL != "https://new.example.com" {
		t.Errorf("updated baseURL = %q, want %q", client.baseURL, "https://new.example.com")
	}
}

// TestClient_ContextCancellation tests that requests respect context cancellation.
func TestClient_ContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(config.Config{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		UserAgent:  "test-agent",
		HTTPClient: server.Client(),
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := client.Get(ctx, "/test", nil)
	if err == nil {
		t.Fatal("expected error due to cancelled context, got nil")
	}
}

// TestClient_HTTPErrorWithRequestID tests that request ID is included in API errors.
func TestClient_HTTPErrorWithRequestID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req-12345")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	client := NewClient(config.Config{
		BaseURL:    server.URL,
		APIKey:     "test-key",
		UserAgent:  "test-agent",
		HTTPClient: server.Client(),
	})

	err := client.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *pkgerrors.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", apiErr.StatusCode, http.StatusInternalServerError)
	}

	if apiErr.RequestID != "req-12345" {
		t.Errorf("request ID = %q, want %q", apiErr.RequestID, "req-12345")
	}

	if apiErr.Message != "internal error" {
		t.Errorf("message = %q, want %q", apiErr.Message, "internal error")
	}
}

// containsString checks if s contains substr.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
