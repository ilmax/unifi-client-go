package sitemanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/murasame29/unifi-go-sdk/pkg/errors"
)

// SDWANConfig represents an SD-WAN configuration.
type SDWANConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// SDWANStatus represents the status of an SD-WAN configuration.
type SDWANStatus struct {
	ConfigID    string      `json:"configId"`
	Status      string      `json:"status"`
	LastUpdated time.Time   `json:"lastUpdated"`
	Peers       []SDWANPeer `json:"peers"`
}

// SDWANPeer represents a peer in an SD-WAN configuration.
type SDWANPeer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Latency int    `json:"latency"`
}

// SDWANConfigsResponse represents the API response for listing SD-WAN configurations.
type SDWANConfigsResponse struct {
	Data []SDWANConfig `json:"data"`
}

// SDWANStatusResponse represents the API response for SD-WAN status.
type SDWANStatusResponse struct {
	Data SDWANStatus `json:"data"`
}

// GetSDWANConfigs retrieves all SD-WAN configurations.
func (s *SiteManager) GetSDWANConfigs() ([]SDWANConfig, error) {
	return s.GetSDWANConfigsWithContext(context.Background())
}

// GetSDWANConfigsWithContext retrieves all SD-WAN configurations.
func (s *SiteManager) GetSDWANConfigsWithContext(ctx context.Context) ([]SDWANConfig, error) {
	var resp SDWANConfigsResponse
	if err := s.client.Get(ctx, "/v1/sdwan/configs", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetSDWANStatus retrieves the status of an SD-WAN configuration.
func (s *SiteManager) GetSDWANStatus(configID string) (*SDWANStatus, error) {
	return s.GetSDWANStatusWithContext(context.Background(), configID)
}

// GetSDWANStatusWithContext retrieves the status of an SD-WAN configuration.
func (s *SiteManager) GetSDWANStatusWithContext(ctx context.Context, configID string) (*SDWANStatus, error) {
	if strings.TrimSpace(configID) == "" {
		return nil, errors.ErrEmptyConfigID
	}

	var resp SDWANStatusResponse
	if err := s.client.Get(ctx, fmt.Sprintf("/v1/sdwan/configs/%s/status", configID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
