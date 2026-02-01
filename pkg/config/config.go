// Package config provides configuration for the UniFi SDK.
package config

import (
	"net/http"
	"strings"
	"time"
)

const (
	DefaultTimeout   = 30 * time.Second
	DefaultUserAgent = "unifi-go-sdk/1.0"
)

// Config contains the configuration for the UniFi SDK.
type Config struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
	Timeout    time.Duration
}

// ConfigOption is a function that configures the Config.
type ConfigOption func(*Config)

// New creates a new Config with default values.
func New() Config {
	return Config{
		UserAgent: DefaultUserAgent,
		Timeout:   DefaultTimeout,
	}
}

// Init initializes the config with the provided options.
func (c *Config) Init(opts []ConfigOption) error {
	for _, opt := range opts {
		opt(c)
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Timeout: c.Timeout,
		}
	}

	return nil
}

// ConfigAPIKey sets the API key.
func ConfigAPIKey(apiKey string) ConfigOption {
	return func(c *Config) {
		c.APIKey = strings.TrimSpace(apiKey)
	}
}

// ConfigBaseURL sets the base URL.
func ConfigBaseURL(baseURL string) ConfigOption {
	return func(c *Config) {
		c.BaseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// ConfigHTTPClient sets a custom HTTP client.
func ConfigHTTPClient(client *http.Client) ConfigOption {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// ConfigTimeout sets the HTTP timeout.
func ConfigTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// ConfigUserAgent sets the User-Agent header.
func ConfigUserAgent(userAgent string) ConfigOption {
	return func(c *Config) {
		c.UserAgent = userAgent
	}
}
