// Package main demonstrates basic usage of the UniFi SDK.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/murasame29/unifi-client-go/pkg/sitemanager"
	"github.com/murasame29/unifi-client-go/unifi"
)

func main() {
	apiKey := os.Getenv("UNIFI_API_KEY")
	if apiKey == "" {
		log.Fatal("UNIFI_API_KEY environment variable is required")
	}

	client, err := unifi.New(unifi.ConfigAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get hosts
	fmt.Println("=== Hosts ===")
	hosts, err := client.SiteManager.ListHosts(nil)
	if err != nil {
		log.Fatalf("Failed to get hosts: %v", err)
	}
	for _, host := range hosts {
		fmt.Printf("Host: %s (ID: %s, Type: %s)\n", host.ID, host.ID, host.Type)
	}

	// Get sites
	fmt.Println("\n=== Sites ===")
	sites, err := client.SiteManager.ListSites(nil)
	if err != nil {
		log.Fatalf("Failed to get sites: %v", err)
	}
	for _, site := range sites {
		fmt.Printf("Site: %s (ID: %s)\n", site.Meta.Name, site.SiteID)
	}

	// Get devices
	fmt.Println("\n=== Devices ===")
	hostDevices, err := client.SiteManager.ListDevices(nil)
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}
	for _, hd := range hostDevices {
		fmt.Printf("Host: %s\n", hd.HostName)
		for _, device := range hd.Devices {
			fmt.Printf("  Device: %s (Model: %s, Status: %s)\n", device.Name, device.Model, device.Status)
		}
	}

	// Get ISP metrics (5 minute interval)
	fmt.Println("\n=== ISP Metrics (5m) ===")
	metrics, err := client.SiteManager.GetISPMetrics(sitemanager.ISPMetricsInterval5m, nil)
	if err != nil {
		log.Fatalf("Failed to get ISP metrics: %v", err)
	}
	for _, data := range metrics {
		fmt.Printf("Host: %s, Site: %s\n", data.HostID, data.SiteID)
		for _, period := range data.Periods {
			wan := period.Data.WAN
			fmt.Printf("  Time: %s, AvgLatency: %dms, Download: %dkbps, Upload: %dkbps\n",
				period.MetricTime.Format("2006-01-02 15:04:05"), wan.AvgLatency, wan.DownloadKbps, wan.UploadKbps)
		}
	}
}
