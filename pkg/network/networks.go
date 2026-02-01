package network

type NetworkDetailsRequest struct {
	NetworkID string `json:"networkId"`
	SiteID    string `json:"siteId"`
}

type NetworkDetailsResponse struct {
	NetworkDetail
}

// GET /v1/sites/{siteId}/networks/{networkId}

type UpdateNetworkRequest struct {
	NetworkID string `json:"networkId"`
	SiteID    string `json:"siteId"`

	Management   NetworkManagementType `json:"management"`
	Name         string                `json:"name"`
	Enabled      bool                  `json:"enabled"`
	VlanID       int                   `json:"vlanId"`
	DHCPGuarding DHCPConfigurationMode `json:"dhcpGuarding"`

	IsolationEnabled      bool              `json:"isolationEnabled"`
	CellularBackupEnabled bool              `json:"cellularBackupEnabled"`
	ZoneID                string            `json:"zoneId"`
	DeviceID              string            `json:"deviceId"`
	InternetAccessEnabled bool              `json:"internetAccessEnabled"`
	MDNSForwardingEnabled bool              `json:"mdnsForwardingEnabled"`
	IPv4Configuration     IPv4Configuration `json:"ipv4Configuration"`
	IPv6Configuration     IPv6Configuration `json:"ipv6Configuration"`
}

type UpdateNetworkResponse struct {
	NetworkDetail
}

// PUT /v1/sites/{siteId}/networks/{networkId}

type DeleteNetworkResponse struct {
	NetworkID string `json:"networkId"`
	SiteID    string `json:"siteId"`
	Cascade   bool   `json:"cascade"`
	Force     bool   `json:"force"`
}

// DELETE /v1/sites/{siteId}/networks/{networkId}

type ListNetworksRequest struct {
	SiteID string `json:"siteId"`

	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

type ListNetworksResponse struct {
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
	Count      int             `json:"count"`
	TotalCount int             `json:"totalCount"`
	Data       []NetworkDetail `json:"data"`
}

// GET /v1/sites/{siteId}/networks

type CreateNetworkRequest struct {
	NetworkID string `json:"networkId"`
	SiteID    string `json:"siteId"`

	Management   NetworkManagementType `json:"management"`
	Name         string                `json:"name"`
	Enabled      bool                  `json:"enabled"`
	VlanID       int                   `json:"vlanId"`
	DHCPGuarding DHCPConfigurationMode `json:"dhcpGuarding"`

	IsolationEnabled      bool              `json:"isolationEnabled"`
	CellularBackupEnabled bool              `json:"cellularBackupEnabled"`
	ZoneID                string            `json:"zoneId"`
	DeviceID              string            `json:"deviceId"`
	InternetAccessEnabled bool              `json:"internetAccessEnabled"`
	MDNSForwardingEnabled bool              `json:"mdnsForwardingEnabled"`
	IPv4Configuration     IPv4Configuration `json:"ipv4Configuration"`
	IPv6Configuration     IPv6Configuration `json:"ipv6Configuration"`
}

type CreateNetworkResponse struct {
	NetworkDetail
}

// POST /v1/sites/{siteId}/networks

type NetworkReferencesRequest struct {
	NetworkID string `json:"networkId"`
	SiteID    string `json:"siteId"`
}

type NetworkReferencesResponse struct {
	ReferenceResources []ReferenceResource `json:"referenceResources"`
}

type ReferenceResource struct {
	ResourceType   string                       `json:"resourceType"`
	ReferenceCount int                          `json:"referenceCount"`
	References     []ReferenceResourceReference `json:"references"`
}

type ReferenceResourceReference struct {
	ReferenceID string `json:"referenceId"`
}

// GET /v1/sites/{siteId}/networks/{networkId}/references
