package sitemanager

import (
	"context"
	"encoding/json"
	"net/url"
)

// Site represents a UniFi Network site.
type Site struct {
	SiteID     string          `json:"siteId"`
	HostID     string          `json:"hostId"`
	Meta       SiteMeta        `json:"meta"`
	Statistics json.RawMessage `json:"statistics"`
	Permission string          `json:"permission"`
	IsOwner    bool            `json:"isOwner"`
}

// SiteMeta contains metadata for a UniFi site.
type SiteMeta struct {
	Desc       string `json:"desc"`
	GatewayMAC string `json:"gatewayMac"`
	Name       string `json:"name"`
	Timezone   string `json:"timezone"`
}

// ListSitesParams contains query parameters for ListSites API.
type ListSitesParams struct {
	PageSize  string
	NextToken string
}

// ToQuery converts the params to URL query string.
func (p *ListSitesParams) ToQuery() string {
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

// ListSitesResponse represents the API response for listing sites.
type ListSitesResponse struct {
	Data      []Site `json:"data"`
	NextToken string `json:"nextToken,omitempty"`
	SuccessResponse
}

// ListSites retrieves sites with optional pagination parameters.
func (s *SiteManager) ListSites(params *ListSitesParams) ([]Site, error) {
	return s.ListSitesWithContext(context.Background(), params)
}

// ListSitesWithContext retrieves sites with optional pagination parameters.
func (s *SiteManager) ListSitesWithContext(ctx context.Context, params *ListSitesParams) ([]Site, error) {
	path := "/v1/sites" + params.ToQuery()

	var resp ListSitesResponse
	if err := s.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
