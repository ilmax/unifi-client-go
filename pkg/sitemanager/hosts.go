package sitemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/murasame29/unifi-client-go/pkg/errors"
)

// Host represents a UniFi console or controller.
type Host struct {
	ID                        string          `json:"id"`
	HardwareID                string          `json:"hardwareId"`
	Type                      string          `json:"type"`
	IPAddress                 string          `json:"ipAddress"`
	Owner                     bool            `json:"owner"`
	IsBlocked                 bool            `json:"isBlocked"`
	RegistrationTime          time.Time       `json:"registrationTime"`
	LastConnectionStateChange time.Time       `json:"lastConnectionStateChange"`
	LatestBackupTime          string          `json:"latestBackupTime"`
	UserData                  json.RawMessage `json:"userData"`
	ReportedState             json.RawMessage `json:"reportedState"`
}

// ListHostsParams contains query parameters for ListHosts API.
type ListHostsParams struct {
	PageSize  string
	NextToken string
}

// ToQuery converts the params to URL query string.
func (p *ListHostsParams) ToQuery() string {
	if p == nil {
		return ""
	}
	params := url.Values{}
	if p.PageSize != "" {
		params.Set("pageSize", p.PageSize)
	}
	if p.NextToken != "" {
		params.Set("nextToken", p.NextToken)
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}

// ListHostsResponse represents the API response for listing hosts.
type ListHostsResponse struct {
	Data      []Host `json:"data"`
	NextToken string `json:"nextToken,omitempty"`
	SuccessResponse
}

// GetHostByIDResponse represents the API response for getting a host by ID.
type GetHostByIDResponse struct {
	Data Host `json:"data"`
	SuccessResponse
}

// ListHosts retrieves hosts with optional pagination parameters.
func (s *SiteManager) ListHosts(params *ListHostsParams) ([]Host, error) {
	return s.ListHostsWithContext(context.Background(), params)
}

// ListHostsWithContext retrieves hosts with optional pagination parameters.
func (s *SiteManager) ListHostsWithContext(ctx context.Context, params *ListHostsParams) ([]Host, error) {
	path := "/v1/hosts" + params.ToQuery()

	var resp ListHostsResponse
	if err := s.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetHostByID retrieves a single host by ID.
func (s *SiteManager) GetHostByID(hostID string) (*Host, error) {
	return s.GetHostByIDWithContext(context.Background(), hostID)
}

// GetHostByIDWithContext retrieves a single host by ID.
func (s *SiteManager) GetHostByIDWithContext(ctx context.Context, hostID string) (*Host, error) {
	if strings.TrimSpace(hostID) == "" {
		return nil, errors.ErrEmptyHostID
	}

	var resp GetHostByIDResponse
	if err := s.client.Get(ctx, fmt.Sprintf("/v1/hosts/%s", hostID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
