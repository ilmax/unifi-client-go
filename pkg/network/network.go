// Package network provides a client for the UniFi Network API.
// This API is used to interact with local UniFi controllers (UDM, Cloud Key, etc.)
package network

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)

// Network is used to interact with the UniFi Network API.
type Network struct {
	httpClient *http.Client
	baseURL    string
	site       string
	csrfToken  string
	isUDM      bool
}

// Config contains configuration for the Network client.
type Config struct {
	// BaseURL is the URL of the UniFi controller (e.g., "https://192.168.1.1:8443")
	BaseURL string
	// Site is the site name (default: "default")
	Site string
	// Timeout is the HTTP timeout
	Timeout time.Duration
	// InsecureSkipVerify skips TLS certificate verification (useful for self-signed certs)
	InsecureSkipVerify bool
}

// New creates a new Network client.
func New(cfg Config) (*Network, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("baseURL is required")
	}

	if cfg.Site == "" {
		cfg.Site = "default"
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		},
	}

	return &Network{
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Jar:       jar,
			Transport: transport,
		},
		baseURL: strings.TrimSuffix(cfg.BaseURL, "/"),
		site:    cfg.Site,
	}, nil
}
