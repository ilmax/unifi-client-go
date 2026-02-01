package network

type NetworkManagementType string

const (
	NetworkManagementTypeUnmanagement NetworkManagementType = "UNMANAGED"
	NetworkManagementTypeGateway      NetworkManagementType = "GATEWAY"
	NetworkManagementTypeSwitch       NetworkManagementType = "SWITCH"
)

type NetworkDetail struct {
	Management            NetworkManagementType     `json:"management"`
	ID                    string                    `json:"id"`
	Name                  string                    `json:"name"`
	Enabled               bool                      `json:"enabled"`
	VlanID                int                       `json:"vlanId"`
	Metadata              NetworkDetailMetadata     `json:"metadata"`
	DHCPGuarding          NetworkDetailDHCPGuarding `json:"dhcpGuarding"`
	IsolationEnabled      bool                      `json:"isolationEnabled"`
	CellularBackupEnabled bool                      `json:"cellularBackupEnabled"`
	ZoneID                string                    `json:"zoneId"`
	DeviceID              string                    `json:"deviceId"`
	InternetAccessEnabled bool                      `json:"internetAccessEnabled"`
	MDNSForwardingEnabled bool                      `json:"mdnsForwardingEnabled"`
	IPv4Configuration     IPv4Configuration         `json:"ipv4Configuration"`
	IPv6Configuration     IPv6Configuration         `json:"ipv6Configuration"`
}

type NetworkMetadataOrigin string

const (
	NetworkMetadataOriginUserDefined   NetworkMetadataOrigin = "USER_DEFINED"
	NetworkMetadataOriginSystemDefined NetworkMetadataOrigin = "SYSTEM_DEFINED"
	NetworkMetadataOriginOrchestrated  NetworkMetadataOrigin = "ORCHESTRATED"
)

type NetworkDetailMetadata struct {
	Origin NetworkMetadataOrigin `json:"origin"`
}

type NetworkDetailDHCPGuarding struct {
	TrustedDHCPServerIpAddresses []string `json:"trustedDhcpServerIpAddresses"`
}

type IPv4Configuration struct {
	AutoScaleEnabled                  bool                              `json:"autoScaleEnabled"`
	HostIpAddress                     string                            `json:"hostIpAddress"`
	PrefixLength                      int                               `json:"prefixLength"`
	AdditionalHostIpSubnets           []string                          `json:"additionalHostIpSubnets"`
	DHCPConfiguration                 IPv4DHCPConfiguration             `json:"dhcpConfiguration"`
	NATOutboundIpAddressConfiguration NATOutboundIpAddressConfiguration `json:"natOutboundIpAddressConfiguration"`
}

type DHCPConfigurationMode string

const (
	DHCPConfigurationModeServer DHCPConfigurationMode = "SERVER"
	DHCPConfigurationModeRELAY  DHCPConfigurationMode = "RELAY"
)

type IPv4DHCPConfiguration struct {
	Mode                         DHCPConfigurationMode `json:"mode"`
	IpAddressRange               IpAddressRange        `json:"ipAddressRange"`
	GatewayIpAddressOverride     string                `json:"gatewayIpAddressOverride"`
	DNSServerIpAddressesOverride []string              `json:"dnsServerIpAddressesOverride"`
	LeaseTimeSeconds             int                   `json:"leaseTimeSeconds"`
	DomainName                   string                `json:"domainName"`
	PingConflictDetectionEnabled bool                  `json:"pingConflictDetectionEnabled"`
	PXEConfiguration             PXEConfiguration      `json:"pxeConfiguration"`
	NTPServerIpAddresses         []string              `json:"ntpServerIpAddresses"`
	Option43Value                string                `json:"option43Value"`
	TFTPServerAddress            string                `json:"tftpServerAddress"`
	TimeOffsetSeconds            int                   `json:"timeOffsetSeconds"`
	WPADURL                      string                `json:"wpadUrl"`
	WINSServerIpAddresses        []string              `json:"winsServerIpAddresses"`
}

type IpAddressRange struct {
	Start string `json:"start"`
	Stop  string `json:"stop"`
}

type PXEConfiguration struct {
	ServerIpAddress string `json:"serverIpAddress"`
	Filename        string `json:"filename"`
}

type NATOutboundIpAddressConfigurationType string

const (
	NATOutboundIpAddressConfigurationTypeAuto   NATOutboundIpAddressConfigurationType = "AUTO"
	NATOutboundIpAddressConfigurationTypeStatic NATOutboundIpAddressConfigurationType = "STATIC"
)

type NATOutboundIpAddressConfiguration struct {
	NATOutboundIpAddressConfigurationType NATOutboundIpAddressConfigurationType `json:"type"`
	WANInterfaceId                        string                                `json:"wanInterfaceId"`
	IpAddressSelectionMode                string                                `json:"ipAddressSelectionMode"`
	IpAddressSelectors                    []IpAddressSelector                   `json:"ipAddressSelectors"`
}

type IpAddressSelectorType string

const (
	IpAddressSelectorTypeIpAddress    IpAddressSelectorType = "IP_ADDRESS"
	IpAddressSelectorTypeAddressRange IpAddressSelectorType = "IP_ADDRESS_RANGE"
)

type IpAddressSelector struct {
	IpAddressSelectorType IpAddressSelectorType `json:"type"`
	IpAddressRange
	Value string `json:"value"`
}

type IPv6InterfaceType string

const (
	IPv6InterfaceTypePrefixDelegation IPv6InterfaceType = "PREFIX_DELEGATION"
	IPv6InterfaceTypeStatic           IPv6InterfaceType = "STATIC"
)

type IPv6Configuration struct {
	InterfaceType                  IPv6InterfaceType       `json:"interfaceType"`
	ClientAddressAssignment        ClientAddressAssignment `json:"clientAddressAssignment"`
	RouterAdvertisement            RouterAdvertisement     `json:"routerAdvertisement"`
	DNSServerIpAddressesOverride   []string                `json:"dnsServerIpAddressesOverride"`
	AdditionalHostIpSubnets        []string                `json:"additionalHostIpSubnets"`
	PrefixDelegationWanInterfaceID string                  `json:"prefixDelegationWanInterfaceId"`
	HostIpAddress                  string                  `json:"hostIpAddress"`
	PrefixLength                   int                     `json:"prefixLength"`
}

type ClientAddressAssignment struct {
	DHCPConfiguration IPv6DHCPConfiguration `json:"dhcpConfiguration"`
	SlaacEnabled      bool                  `json:"slaacEnabled"`
}

type IPv6DHCPConfiguration struct {
	IPAddressSuffixRange IpAddressRange `json:"ipAddressSuffixRange"`
	LeaseTimeSeconds     int            `json:"leaseTimeSeconds"`
}

type RouterAdvertisement struct {
	Priority string `json:"priority"`
}

// GET /v1/sites/{siteId}/networks/{networkId}
