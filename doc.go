// Package unifi provides a Go SDK for the UniFi Site Manager API.
//
// This SDK provides a simple and intuitive interface for managing UniFi infrastructure,
// including hosts, sites, devices, ISP metrics, and SD-WAN configurations.
//
// Basic usage:
//
//	client, err := unifi.NewClient("your-api-key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	hosts, err := client.GetHosts(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// For pagination support, use the Paginator:
//
//	iter := client.Paginator.GetHosts(ctx)
//	for iter.Next() {
//	    host := iter.Value()
//	    // process host
//	}
//	if err := iter.Err(); err != nil {
//	    log.Fatal(err)
//	}
package unifi
