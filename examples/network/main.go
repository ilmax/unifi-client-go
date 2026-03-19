//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	"github.com/ilmax/unifi-client-go/pkg/network"
	"github.com/ilmax/unifi-client-go/unifi"
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
	// Generated methods are added by the release workflow in pkg/network/openapi.gen.go.
	_ = client
	_ = ctx
}
