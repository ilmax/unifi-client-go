//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/murasame29/unifi-client-go/pkg/network"
	"github.com/murasame29/unifi-client-go/unifi"
)

func main() {
	// Create Network API client for local UniFi controller
	client, err := unifi.NewNetwork(network.Config{
		BaseURL:            os.Getenv("UNIFI_CONTROLLER_URL"), // e.g., "https://192.168.1.1"
		Site:               "default",
		InsecureSkipVerify: true, // For self-signed certificates
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Login with username and password
	username := os.Getenv("UNIFI_USERNAME")
	password := os.Getenv("UNIFI_PASSWORD")
	if err := client.Login(ctx, username, password); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer client.Logout(ctx)

	// List all sites
	sites, err := client.ListSites(ctx)
	if err != nil {
		log.Fatalf("Failed to list sites: %v", err)
	}
	fmt.Printf("Found %d sites\n", len(sites))
	for _, site := range sites {
		fmt.Printf("  - %s (%s)\n", site.Desc, site.Name)
	}

	// Get site health
	health, err := client.GetSiteHealth(ctx)
	if err != nil {
		log.Fatalf("Failed to get site health: %v", err)
	}
	for _, h := range health {
		fmt.Printf("Subsystem: %s, Status: %s\n", h.Subsystem, h.Status)
	}

	// List all devices
	devices, err := client.ListDevices(ctx)
	if err != nil {
		log.Fatalf("Failed to list devices: %v", err)
	}
	fmt.Printf("\nFound %d devices\n", len(devices))
	for _, device := range devices {
		fmt.Printf("  - %s (%s) - %s\n", device.Name, device.Model, device.IP)
	}

	// List all clients
	clients, err := client.ListClients(ctx)
	if err != nil {
		log.Fatalf("Failed to list clients: %v", err)
	}
	fmt.Printf("\nFound %d connected clients\n", len(clients))

	// List WLANs
	wlans, err := client.ListWLANs(ctx)
	if err != nil {
		log.Fatalf("Failed to list WLANs: %v", err)
	}
	fmt.Printf("\nFound %d WLANs\n", len(wlans))
	for _, wlan := range wlans {
		status := "disabled"
		if wlan.Enabled {
			status = "enabled"
		}
		fmt.Printf("  - %s (%s)\n", wlan.Name, status)
	}
}
