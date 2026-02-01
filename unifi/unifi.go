// Package unifi provides a unified client for all UniFi APIs.
package unifi

import (
	"github.com/murasame29/unifi-client-go/pkg/config"
	"github.com/murasame29/unifi-client-go/pkg/errors"
	"github.com/murasame29/unifi-client-go/pkg/network"
	"github.com/murasame29/unifi-client-go/pkg/sitemanager"
)

// UniFi is a collection of UniFi APIs.
// Use New() for Site Manager API (cloud API with API key).
// Use NewNetwork() for Network API (local controller with username/password).
type UniFi struct {
	SiteManager *sitemanager.SiteManager

	config config.Config
}

// New returns a collection of UniFi APIs for Site Manager (cloud API).
// Requires an API key for authentication.
func New(opts ...ConfigOption) (*UniFi, error) {
	cfg := config.New()

	if err := cfg.Init(opts); err != nil {
		return nil, err
	}

	if cfg.APIKey == "" {
		return nil, errors.ErrEmptyAPIKey
	}

	u := &UniFi{
		config:      cfg,
		SiteManager: sitemanager.New(cfg),
	}

	return u, nil
}

// NewNetwork creates a new Network API client for local UniFi controllers.
// Use this for UDM, Cloud Key, or software-based controllers.
func NewNetwork(cfg network.Config) (*network.Network, error) {
	return network.New(cfg)
}

// ConfigOption configures the UniFi client.
type ConfigOption = config.ConfigOption

// ConfigAPIKey sets the API key for authentication.
func ConfigAPIKey(apiKey string) ConfigOption {
	return config.ConfigAPIKey(apiKey)
}

// ConfigBaseURL sets the base URL for API requests.
func ConfigBaseURL(baseURL string) ConfigOption {
	return config.ConfigBaseURL(baseURL)
}

// ConfigTimeout sets the HTTP timeout.
func ConfigTimeout(timeout config.ConfigOption) ConfigOption {
	return timeout
}

// ConfigUserAgent sets the User-Agent header.
func ConfigUserAgent(userAgent string) ConfigOption {
	return config.ConfigUserAgent(userAgent)
}
