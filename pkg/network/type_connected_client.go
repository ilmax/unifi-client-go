package network

type ConnectedClientType string

const (
	ConnectedClientTypeWired    ConnectedClientType = "WIRED"
	ConnectedClientTypeWireless ConnectedClientType = "WIRELESS"
	ConnectedClientTypeVPN      ConnectedClientType = "VPN"
	ConnectedClientTypeTeleport ConnectedClientType = "TELEPORT"
)

type ConnectedClient struct {
	ConnectedType  ConnectedClientType   `json:"type"`
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	ConnectedAt    string                `json:"connectedAt"`
	IpAddress      string                `json:"ipAddress"`
	Access         ConnectedClientAccess `json:"access"`
	MacAddress     string                `json:"macAddress,omitempty"`
	UplinkDeviceID string                `json:"uplinkDeviceId,omitempty"`
}

type ConnectedClientAccessType string

const (
	ConnectedClientAccessDefault ConnectedClientAccessType = "DEFAULT"
	ConnectedClientAccessGuest   ConnectedClientAccessType = "GUEST"
)

type ConnectedClientAccess struct {
	ConnectedClientAccessType ConnectedClientAccessType `json:"type"`
	Authorized                bool                      `json:"authorized,omitempty"`
	Athorization              ClientActionAuthorization `json:"authorization,omitempty"`
}

type ClientActionAuthorization struct {
	AuthorizedAt         string                         `json:"authorizedAt"`
	AuthorizationMethod  string                         `json:"authorizationMethod"`
	ExpiresAt            string                         `json:"expiresAt"`
	DataUsageLimitMBytes int                            `json:"dataUsageLimitMBytes,omitempty"`
	RxRateLimitKbps      int                            `json:"rxRateLimitKbps,omitempty"`
	TxRateLimitKbps      int                            `json:"txRateLimitKbps,omitempty"`
	Usage                ClientActionAuthorizationUsage `json:"usage"`
}

type ClientActionAuthorizationUsage struct {
	DurationSec int `json:"durationSec"`
	RxBytes     int `json:"rxBytes"`
	TxBytes     int `json:"txBytes"`
	Bytes       int `json:"Bytes"`
}
