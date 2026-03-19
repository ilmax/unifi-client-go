// Package unifi provides constructors for UniFi API clients.
package unifi

import "github.com/ilmax/unifi-client-go/pkg/network"

// NewNetwork creates a new Network API client for local UniFi controllers.
func NewNetwork(cfg network.Config) (*network.Network, error) {
	return network.New(cfg)
}
