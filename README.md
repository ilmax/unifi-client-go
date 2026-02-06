# UniFi Go SDK

Go SDK for the UniFi API. It provides auto-generated types and client methods based on the UniFi API documentation.

## Features

- **Auto-generated**: Types and client methods are generated from the UniFi API documentation
- **Type-safe**: Strongly typed request and response models for all APIs
- **Standard library only**: No external dependencies (uses only `net/http`, `encoding/json`, `context`, etc.)
- **Context-first**: All API methods accept `ctx context.Context` as the first argument

## Installation

```bash
go get github.com/ilmax/unifi-client-go@v9.1.120
```

Versions correspond to UniFi API versions. See [Releases](https://github.com/ilmax/unifi-client-go/releases) for available versions.

## Usage

### Site Manager API (Cloud API)

The Site Manager API manages devices through UniFi cloud services. Authentication requires an API key.

```go
package main

import (
    "context"
    "log"

    "github.com/ilmax/unifi-client-go/unifi"
)

func main() {
    // Initialize client
    client, err := unifi.New(
        unifi.ConfigAPIKey("your-api-key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Use Site Manager API
    // client.SiteManager.XXX(ctx, ...)
    _ = client
    _ = ctx
}
```

### Network API (Local Controller)

The Network API communicates directly with a local UniFi controller (UDM, Cloud Key, etc.).

```go
package main

import (
    "context"
    "log"

    "github.com/ilmax/unifi-client-go/unifi"
    "github.com/ilmax/unifi-client-go/pkg/network"
)

func main() {
    // Initialize network client
    client, err := unifi.NewNetwork(network.Config{
        BaseURL:            "https://192.168.1.1:8443",
        Site:               "default",
        InsecureSkipVerify: true, // For self-signed certificates
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Use Network API
    // client.XXX(ctx, ...)
    _ = client
    _ = ctx
}
```

## Directory Structure

```
unifi-go-sdk/
├── .github/
│   └── workflows/
│       └── release.yml          # Release workflow
├── unifi/
│   └── unifi.go                 # Main client
├── pkg/
│   ├── config/                  # Configuration
│   ├── errors/                  # Custom errors
│   ├── network/                 # Network API (generated)
│   │   ├── types.go             # Type definitions
│   │   └── network.go           # Client methods
│   └── sitemanager/             # Site Manager API
├── internal/
│   ├── http/                    # HTTP client
│   ├── typegen/                 # Type generator
│   └── clientgen/               # Client method generator
├── cmd/
│   ├── typegen/                 # Type generator CLI
│   └── clientgen/               # Client method generator CLI
├── go.mod
└── README.md
```

## About Auto-Generation

Types and client methods in this SDK are auto-generated from UniFi API documentation.

### Generated Types

- **Request types**: API request body structures
- **Response types**: API response structures
- **Shared types**: Common entities such as Voucher, Client, Site, etc.

### Generated Client Methods

Methods are generated for each API endpoint:

| HTTP Method | Endpoint Name | Generated Method |
|-------------|---------------|------------------|
| GET (list)  | List Clients  | ListClients      |
| GET (single)| Get Client    | GetClient        |
| POST        | Create Voucher| CreateVoucher    |
| PUT         | Update Device | UpdateDevice     |
| DELETE      | Delete Voucher| DeleteVoucher    |

## Release Workflow

To release an SDK for a new API version, use GitHub Actions `workflow_dispatch`.

### Manual Release Steps

1. Open the "Actions" tab in the GitHub repository
2. Select the "Release SDK" workflow
3. Click "Run workflow"
4. Enter the UniFi API version (e.g. `v9.1.120`) in `api_version`
5. Click "Run workflow" to execute

### Workflow Steps

1. Create a release branch (`release/vX.Y.Z`)
2. Generate types from the UniFi API documentation (`typegen`)
3. Generate client methods (`clientgen`)
4. Format, build, and test code
5. Commit and push changes
6. Create a version tag
7. Create a GitHub release

## Error Handling

The SDK provides custom error types in `pkg/errors`:

```go
import "github.com/ilmax/unifi-client-go/pkg/errors"

// Check API errors
if errors.Is(err, errors.ErrUnauthorized) {
    // Authentication error
}

if errors.Is(err, errors.ErrNotFound) {
    // Resource not found
}

// Get details of APIError
var apiErr *errors.APIError
if errors.As(err, &apiErr) {
    log.Printf("Status: %d, Message: %s", apiErr.StatusCode, apiErr.Message)
}
```

## Configuration Options

### Site Manager API

```go
client, err := unifi.New(
    unifi.ConfigAPIKey("your-api-key"),           // API key (required)
    unifi.ConfigBaseURL("https://api.ui.com"),    // Base URL (optional)
    unifi.ConfigUserAgent("my-app/1.0"),          // User-Agent (optional)
)
```

### Network API

```go
client, err := unifi.NewNetwork(network.Config{
    BaseURL:            "https://192.168.1.1:8443", // Controller URL (required)
    Site:               "default",                   // Site name (default: "default")
    Timeout:            30 * time.Second,            // Timeout (default: 30s)
    InsecureSkipVerify: true,                        // Skip TLS verification for self-signed certs
})
```

## Development

### Run the type generator

```bash
go run ./cmd/typegen/main.go \
    -discover "https://developer.ui.com/network/v9.1.120" \
    -output-dir ./pkg \
    -package network \
    -workers 4
```

### Run the client generator

```bash
go run ./cmd/clientgen/main.go \
    -input ./pkg \
    -output ./pkg
```

## License

MIT
