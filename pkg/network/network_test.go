package network

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sdkerrors "github.com/ilmax/unifi-client-go/pkg/errors"
)

func TestNewRequiresBaseURL(t *testing.T) {
	t.Parallel()

	_, err := New(Config{APIKey: "test-key"})
	if err == nil {
		t.Fatal("expected error for missing baseURL")
	}
}

func TestNewRequiresAPIKey(t *testing.T) {
	t.Parallel()

	_, err := New(Config{BaseURL: "https://example.com"})
	if !errors.Is(err, sdkerrors.ErrEmptyAPIKey) {
		t.Fatalf("expected ErrEmptyAPIKey, got %v", err)
	}
}

func TestNewInjectsAPIKeyHeader(t *testing.T) {
	t.Parallel()

	var gotHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get(apiKeyHeader)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := New(Config{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	resp, err := client.GetInfo(context.Background())
	if err != nil {
		t.Fatalf("GetInfo() error = %v", err)
	}
	defer resp.Body.Close()

	if gotHeader != "test-api-key" {
		t.Fatalf("expected %s header to be set, got %q", apiKeyHeader, gotHeader)
	}
}

func TestNewUsesDefaultTimeout(t *testing.T) {
	t.Parallel()

	client := newHTTPClient(Config{})
	if client.Timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultTimeout, client.Timeout)
	}
}

func TestNewUsesConfiguredTimeout(t *testing.T) {
	t.Parallel()

	want := 5 * time.Second
	client := newHTTPClient(Config{Timeout: want})
	if client.Timeout != want {
		t.Fatalf("expected timeout %v, got %v", want, client.Timeout)
	}
}
