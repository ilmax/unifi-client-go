package sitemanager

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ilmax/unifi-client-go/pkg/errors"
)

// ISPMetricsInterval represents the interval for ISP metrics data.
type ISPMetricsInterval string

const (
	ISPMetricsInterval5m ISPMetricsInterval = "5m"
	ISPMetricsInterval1h ISPMetricsInterval = "1h"
)

// IsValidISPMetricsInterval checks if the given interval is valid.
func IsValidISPMetricsInterval(interval ISPMetricsInterval) bool {
	return interval == ISPMetricsInterval5m || interval == ISPMetricsInterval1h
}

// ISPMetricsData represents ISP metrics data.
type ISPMetricsData struct {
	MetricType string            `json:"metricType"`
	Periods    []ISPMetricPeriod `json:"periods"`
	HostID     string            `json:"hostId"`
	SiteID     string            `json:"siteId"`
}

// ISPMetricPeriod represents a single period of ISP metrics.
type ISPMetricPeriod struct {
	Data       ISPMetricPeriodData `json:"data"`
	MetricTime time.Time           `json:"metricTime"`
	Version    string              `json:"version"`
}

// ISPMetricPeriodData contains the WAN metrics for a period.
type ISPMetricPeriodData struct {
	WAN ISPMetricWAN `json:"wan"`
}

// ISPMetricWAN contains WAN-specific ISP metrics.
type ISPMetricWAN struct {
	AvgLatency   int    `json:"avgLatency"`
	DownloadKbps int    `json:"download_kbps"`
	Downtime     int    `json:"downtime"`
	ISPAsn       string `json:"ispAsn"`
	ISPName      string `json:"ispName"`
	MaxLatency   int    `json:"maxLatency"`
	PacketLoss   int    `json:"packetLoss"`
	UploadKbps   int    `json:"upload_kbps"`
	Uptime       int    `json:"uptime"`
}

// GetISPMetricsParams contains query parameters for GetISPMetrics API.
type GetISPMetricsParams struct {
	BeginTimestamp string
	EndTimestamp   string
	Duration       string
}

// ToQuery converts the params to URL query string.
func (p *GetISPMetricsParams) ToQuery() string {
	if p == nil {
		return ""
	}
	params := url.Values{}
	if p.BeginTimestamp != "" {
		params.Set("beginTimestamp", p.BeginTimestamp)
	}
	if p.EndTimestamp != "" {
		params.Set("endTimestamp", p.EndTimestamp)
	}
	if p.Duration != "" {
		params.Set("duration", p.Duration)
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}

// GetISPMetricsResponse represents the API response for ISP metrics.
type GetISPMetricsResponse struct {
	Data []ISPMetricsData `json:"data"`
	SuccessResponse
}

// GetISPMetrics retrieves ISP performance metrics for the specified interval.
func (s *SiteManager) GetISPMetrics(interval ISPMetricsInterval, params *GetISPMetricsParams) ([]ISPMetricsData, error) {
	return s.GetISPMetricsWithContext(context.Background(), interval, params)
}

// GetISPMetricsWithContext retrieves ISP performance metrics for the specified interval.
func (s *SiteManager) GetISPMetricsWithContext(ctx context.Context, interval ISPMetricsInterval, params *GetISPMetricsParams) ([]ISPMetricsData, error) {
	if !IsValidISPMetricsInterval(interval) {
		return nil, errors.ErrInvalidInterval
	}

	path := fmt.Sprintf("/v1/isp-metrics/%s", interval) + params.ToQuery()

	var resp GetISPMetricsResponse
	if err := s.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
