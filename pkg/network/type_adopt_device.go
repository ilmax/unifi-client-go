package network

import (
	"encoding/json"
	"time"
)

type AdoptDevice struct {
	ID                string               `json:"id"`
	MacAddress        string               `json:"macAddress"`
	IPAddress         string               `json:"ipAddress"`
	Name              string               `json:"name"`
	Model             string               `json:"model"`
	Supported         bool                 `json:"supported"`
	State             string               `json:"state"`
	FirmwareVersion   string               `json:"firmwareVersion"`
	FirmwareUpdatable bool                 `json:"firmwareUpdatable"`
	AdoptedAt         time.Time            `json:"adoptedAt"`
	ProvisionedAt     time.Time            `json:"provisionedAt"`
	ConfigurationID   string               `json:"configurationId"`
	Uplink            AdoptDeviceUplink    `json:"uplink"`
	Features          AdoptDeviceFeatures  `json:"features"`
	Interfaces        AdoptDeviceInterface `json:"interfaces"`
}

type AdoptDeviceUplink struct {
	DeviceID string `json:"deviceId"`
}

type AdoptDeviceFeatures struct {
	Switching   json.RawMessage `json:"switching"`
	AccessPoint json.RawMessage `json:"accessPoint"`
}

type AdoptDeviceInterface struct {
	Ports  []AdoptDeviceInterfacePort   `json:"ports"`
	Radios []AdoptDeviceInterfaceRadios `json:"radios"`
}

type AdoptDeviceInterfacePort struct {
	Idx          int                         `json:"idx"`
	State        string                      `json:"state"`
	Connector    string                      `json:"connector"`
	MaxSpeedMbps int                         `json:"maxSpeedMbps"`
	SpeedMbps    int                         `json:"speedMbps"`
	Poe          AdoptDeviceInterfacePortPoE `json:"poe"`
}

type AdoptDeviceInterfacePortPoE struct {
	Standard string `json:"standard"`
	Type     int    `json:"type"`
	Enabled  bool   `json:"enabled"`
	State    string `json:"state"`
}

type AdoptDeviceInterfaceRadios struct {
	WlanStandard    string `json:"wlanStandard"`
	FrequencyGHz    int    `json:"frequencyGHz"`
	ChannelWidthMHz int    `json:"channelWidthMHz"`
	Channel         int    `json:"channel"`
}
