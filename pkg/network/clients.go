package network

type ExecuteClientAction string

const (
	ExecuteClientActionAuthorizeGuestAccess   ExecuteClientAction = "AUTHORIZE_GUEST_ACCESS"
	ExecuteClientActionUnAuthorizeGuestAccess ExecuteClientAction = "UNAUTHORIZE_GUEST_ACCESS"
)

type ExecuteClientActionRequest struct {
	ClientID string `json:"clientId"`
	SiteID   string `json:"siteId"`

	Action               ExecuteClientAction `json:"action"`
	TimeLimitMinutes     int                 `json:"timeLimitMinutes,omitempty"`
	DataUsageLimitMBytes int                 `json:"dataUsageLimitMBytes,omitempty"`
	RxRateLimitKbps      int                 `json:"rxRateLimitKbps,omitempty"`
	TxRateLimitKbps      int                 `json:"txRateLimitKbps,omitempty"`
}

type ExecuteClientActionResponse struct {
	Action ExecuteClientAction `josn:"action"`

	RevokedAuthorization ClientActionAuthorization `json:"revokedAuthorization,omitempty"`
	GrantedAuthorization ClientActionAuthorization `json:"grantedAuthorization,omitempty"`
}

// POST /v1/sites/{siteId}/clients/{clientId}/actions

type ConnectedClientsRequest struct {
	SiteID string `json:"siteId"`

	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

type ConnectedClientsResponse struct {
	Offset     int               `json:"offset"`
	Limit      int               `json:"limit"`
	Count      int               `json:"count"`
	TotalCount int               `json:"totalCount"`
	Data       []ConnectedClient `json:"data"`
}

// GET /v1/sites/{siteId}/clients

type ConnectedClientDetailsRequest struct {
	ClientID string `json:"clientId"`
	SiteID   string `json:"siteId"`
}

type ConnectedClientDetailsResponse struct {
	ConnectedClient
}

// GET /v1/sites/{siteId}/clients/{clientId}
