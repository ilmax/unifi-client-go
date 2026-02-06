// Package sitemanager provides a client for the UniFi Site Manager API.
package sitemanager

import (
	"github.com/ilmax/unifi-client-go/internal/http"
	"github.com/ilmax/unifi-client-go/pkg/config"
)

const DefaultBaseURL = "https://api.ui.com"

// SiteManager is used to interact with the UniFi Site Manager API.
type SiteManager struct {
	client http.Client
}

// New returns a new client for interacting with the Site Manager API.
func New(cfg config.Config) *SiteManager {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	return &SiteManager{
		client: http.NewClient(cfg),
	}
}
