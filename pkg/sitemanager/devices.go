package sitemanager

import (
	"context"
	"net/url"
	"time"
)

// Device represents a UniFi device.
type Device struct {
	ID              string      `json:"id"`
	MAC             string      `json:"mac"`
	Name            string      `json:"name"`
	Model           string      `json:"model"`
	Shortname       string      `json:"shortname"`
	IP              string      `json:"ip"`
	ProductLine     string      `json:"productLine"`
	Status          string      `json:"status"`
	Version         string      `json:"version"`
	FirmwareStatus  string      `json:"firmwareStatus"`
	UpdateAvailable string      `json:"updateAvailable"`
	IsConsole       bool        `json:"isConsole"`
	IsManaged       bool        `json:"isManaged"`
	StartupTime     *time.Time  `json:"startupTime"`
	AdoptionTime    *time.Time  `json:"adoptionTime"`
	Note            string      `json:"note"`
	UIDB            *DeviceUIDB `json:"uidb,omitempty"`
}

// DeviceUIDB contains UIDB information for a device.
type DeviceUIDB struct {
	GUID   string           `json:"guid"`
	IconID string           `json:"iconId"`
	ID     string           `json:"id"`
	Images DeviceUIDBImages `json:"images"`
}

// DeviceUIDBImages contains image identifiers for a device.
type DeviceUIDBImages struct {
	Default   string `json:"default"`
	NoPadding string `json:"nopadding"`
	Topology  string `json:"topology"`
}

// HostDevices represents devices grouped by host.
type HostDevices struct {
	HostID    string    `json:"hostId"`
	HostName  string    `json:"hostName"`
	Devices   []Device  `json:"devices"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ListDevicesParams contains query parameters for ListDevices API.
type ListDevicesParams struct {
	HostIDs   []string
	Time      string
	PageSize  string
	NextToken string
}

// ToQuery converts the params to URL query string.
func (p *ListDevicesParams) ToQuery() string {
	if p == nil {
		return ""
	}
	params := url.Values{}
	for _, hostID := range p.HostIDs {
		params.Add("hostIds", hostID)
	}
	if p.Time != "" {
		params.Set("time", p.Time)
	}
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

// ListDevicesResponse represents the API response for listing devices.
type ListDevicesResponse struct {
	Data      []HostDevices `json:"data"`
	NextToken string        `json:"nextToken,omitempty"`
	SuccessResponse
}

// ListDevices retrieves devices with optional pagination parameters.
func (s *SiteManager) ListDevices(params *ListDevicesParams) ([]HostDevices, error) {
	return s.ListDevicesWithContext(context.Background(), params)
}

// ListDevicesWithContext retrieves devices with optional pagination parameters.
func (s *SiteManager) ListDevicesWithContext(ctx context.Context, params *ListDevicesParams) ([]HostDevices, error) {
	path := "/v1/devices" + params.ToQuery()

	var resp ListDevicesResponse
	if err := s.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
