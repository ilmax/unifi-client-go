package network

import "time"

type ListAdoptedDevicesRequest struct {
	SiteID string `json:"siteId"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

type ListAdoptedDevicesResponse struct {
	Offset     int                     `json:"offset"`
	Limit      int                     `json:"limit"`
	Count      int                     `json:"count"`
	TotalCount int                     `json:"totalCount"`
	Data       []AdoptedDeviceOverview `json:"data"`
}

type AdoptedDeviceOverview struct {
	ID                string   `json:"id"`
	MacAddress        string   `json:"macAddress"`
	IPAddress         string   `json:"ipAddress"`
	Name              string   `json:"name"`
	Model             string   `json:"model"`
	State             string   `json:"state"`
	Supported         bool     `json:"supported"`
	FirmwareVersion   string   `json:"firmwareVersion"`
	FirmwareUpdatable bool     `json:"firmwareUpdatable"`
	Features          []string `json:"features"`
	Interfaces        []string `json:"interfaces"`
}

// GET /v1/sites/{siteId}/devices

type AdoptDeviceRequest struct {
	SiteID            string `json:"siteId"`
	MacAddress        string `json:"macAddress"`
	IgnoreDeviceLimit bool   `json:"ignoreDeviceLimit"`
}

type AdoptDeviceResponse struct {
	AdoptDevice
}

// POST /v1/sites/{siteId}/devices

type ExecutePortActionRequest struct {
	PortIDx  int    `json:"portIdx"`
	SiteID   string `json:"siteId"`
	DeviceID string `json:"deviceId"`
	Action   string `json:"action"`
}

// POST /v1/sites/{siteId}/devices/{deviceId}/interfaces/ports/{portIdx}/actions

type ExecuteAdoptDeviceActionRequest struct {
	SiteID   string `json:"siteId"`
	DeviceID string `json:"deviceId"`
	Action   string `json:"action"`
}

// POST /v1/sites/{siteId}/devices/{deviceId}/actions

type AdoptDeviceDetailRequest struct {
	SiteID   string `json:"siteId"`
	DeviceID string `json:"deviceId"`
}

type AdoptDeviceDetailResponse struct {
	AdoptDevice
}

// GET /v1/sites/{siteId}/devices/{deviceId}

type LatestAdoptedDeviceStatisticsRequest struct {
	SiteID   string `json:"siteId"`
	DeviceID string `json:"deviceId"`
}

type LatestAdoptedDeviceStatisticsResponse struct {
	UptimeSec            int                                    `json:"uptimeSec"`
	LastHeartbeatAt      time.Time                              `json:"lastHeartbeatAt"`
	NextHeartbeatAt      time.Time                              `json:"nextHeartbeatAt"`
	LoadAverage1Min      int                                    `json:"loadAverage1Min"`
	LoadAverage5Min      int                                    `json:"loadAverage5Min"`
	LoadAverage15Min     int                                    `json:"loadAverage15Min"`
	CPUUtilizationPct    int                                    `json:"cpuUtilizationPct"`
	MemoryUtilizationPct int                                    `json:"memoryUtilizationPct"`
	Uplink               LatestAdoptedDeviceStatisticsUplink    `json:"uplink"`
	Interfaces           LatestAdoptedDeviceStatisticsInterface `json:"interfaces"`
}

type LatestAdoptedDeviceStatisticsUplink struct {
	TxRateBps int `json:"txRateBps"`
	RxRateBps int `json:"rxRateBps"`
}

type LatestAdoptedDeviceStatisticsInterface struct {
	Radios []LatestAdoptedDeviceStatisticsInterfaceRadio `json:"radios"`
}

type LatestAdoptedDeviceStatisticsInterfaceRadio struct {
	FrequencyGHz int `json:"frequencyGHz"`
	TxRetriesPct int `json:"txRetriesPct"`
}

// GET /v1/sites/{siteId}/devices/{deviceId}/statistics/latest

type DevicesPendingAdoptionRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

type DevicesPendingAdoptionResponse struct {
	Offset     int                          `json:"offset"`
	Limit      int                          `json:"limit"`
	Count      int                          `json:"count"`
	TotalCount int                          `json:"totalCount"`
	Data       []DevicesPendingAdoptionData `json:"data"`
}

type DevicesPendingAdoptionData struct {
	MacAddress        string   `json:"macAddress"`
	IPAddress         string   `json:"ipAddress"`
	Model             string   `json:"model"`
	State             string   `json:"state"`
	Supported         bool     `json:"supported"`
	FirmwareVersion   string   `json:"firmwareVersion"`
	FirmwareUpdatable bool     `json:"firmwareUpdatable"`
	Features          []string `json:"features"`
}

// GET /v1/pending-devices
