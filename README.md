# UniFi Go SDK

Go SDK for the UniFi API. It provides auto-generated types and client methods based on the UniFi Network OpenAPI spec.

## Features

- **Auto-generated**: Types and client methods are generated from the UniFi Network OpenAPI spec
- **Type-safe**: Strongly typed request and response models for all APIs
- **oapi-codegen**: Client code is generated via `oapi-codegen`
- **Context-first**: All API methods accept `ctx context.Context` as the first argument

## Installation

```bash
go get github.com/ilmax/unifi-client-go@v0.1.0
```

SDK versions are auto-incremented and independent from UniFi API versions. See [Releases](https://github.com/ilmax/unifi-client-go/releases) for available versions.

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
│   │   ├── openapi.gen.go       # Generated types + client
│   │   └── network.go           # Client helpers/config
│   └── sitemanager/             # Site Manager API
├── internal/
│   └── http/                    # HTTP client
├── go.mod
└── README.md
```

## About Auto-Generation

Types and client methods in this SDK are auto-generated from the UniFi Network OpenAPI spec:

`https://raw.githubusercontent.com/beezly/unifi-apis/main/unifi-network/{version}.json`

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

1. Determine the next free SDK version by starting from [VERSION](VERSION) and incrementing until both tag and release branch names are available
2. Download the OpenAPI spec for the requested API version
3. Generate types and client via `oapi-codegen`
4. If `unifi-network-version` contains a valid version, compute OpenAPI diffs via `oasdiff` and include them in the release notes
5. Format, build, and test code
6. Update `unifi-network-version` to the released API version
7. Commit and push changes to both the release branch and `main`
8. Create a version tag (SDK version)
9. Create a GitHub release

Release notes include an OpenAPI diff and a breaking-changes section when `unifi-network-version` is present and valid.
If the file is missing or invalid, the workflow falls back to the latest GitHub release and parses only the first 10 body lines for `UniFi SDK for API version {version}` (`v` prefix optional).

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

### Generate the Network client from OpenAPI

```bash
curl -fsSL \
    "https://raw.githubusercontent.com/beezly/unifi-apis/main/unifi-network/v9.1.120.json" \
    -o /tmp/unifi-network.json

go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
    -generate types,client \
    -package network \
    -o ./pkg/network/openapi.gen.go \
    /tmp/unifi-network.json
```

## License

MIT
