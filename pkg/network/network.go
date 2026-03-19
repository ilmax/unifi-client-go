// Package network provides a client for the UniFi Network API.
package network

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	sdkerrors "github.com/ilmax/unifi-client-go/pkg/errors"
)

const (
	DefaultTimeout = 30 * time.Second
	apiKeyHeader   = "X-API-Key"
)

// Network wraps the generated UniFi Network clients with SDK-level configuration.
type Network struct {
	*Client
	*ClientWithResponses
}

// Config contains configuration for the Network client.
type Config struct {
	// BaseURL is the URL of the UniFi controller (e.g., "https://192.168.1.1:8443").
	BaseURL string
	// APIKey is the UniFi Network API key used for authentication.
	APIKey string
	// Timeout is the HTTP timeout.
	Timeout time.Duration
	// InsecureSkipVerify skips TLS certificate verification (useful for self-signed certificates).
	InsecureSkipVerify bool
	// HTTPClient overrides the default client. Intended mainly for tests and advanced callers.
	HTTPClient *http.Client
}

// New creates a new Network client backed by the generated OpenAPI client.
func New(cfg Config) (*Network, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, fmt.Errorf("baseURL is required")
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, sdkerrors.ErrEmptyAPIKey
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = newHTTPClient(cfg)
	}

	opts := []ClientOption{
		WithHTTPClient(httpClient),
		WithRequestEditorFn(apiKeyRequestEditor(cfg.APIKey)),
	}

	client, err := NewClient(strings.TrimSuffix(cfg.BaseURL, "/"), opts...)
	if err != nil {
		return nil, err
	}

	return &Network{
		Client:              client,
		ClientWithResponses: &ClientWithResponses{ClientInterface: client},
	}, nil
}

func newHTTPClient(cfg Config) *http.Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		},
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

func apiKeyRequestEditor(apiKey string) RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set(apiKeyHeader, apiKey)
		return nil
	}
}
