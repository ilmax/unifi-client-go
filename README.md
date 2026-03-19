# UniFi Go SDK

Go SDK for the UniFi Network API. It provides generated types and client methods based on the published UniFi Network OpenAPI spec.

## Features

- **Auto-generated**: Types and client methods are generated from the UniFi Network OpenAPI spec
- **Type-safe**: Strongly typed request and response models for UniFi Network APIs
- **oapi-codegen**: Client code is generated via `oapi-codegen`
- **Context-first**: All API methods accept `ctx context.Context` as the first argument

## Installation

```bash
go get github.com/ilmax/unifi-client-go@v0.1.0
```

SDK versions are auto-incremented and independent from UniFi API versions. See [Releases](https://github.com/ilmax/unifi-client-go/releases) for available versions.

## Usage

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
│   ├── errors/                  # Custom errors
│   ├── network/                 # Network API (generated)
│   │   ├── openapi.gen.go       # Generated types + client
│   │   └── network.go           # Client helpers/config
├── openapi/
│   └── unifi-network/           # Committed UniFi Network specs by version
├── go.mod
└── README.md
```

## About Auto-Generation

Types and client methods in this SDK are auto-generated from the UniFi Network OpenAPI spec:

`https://raw.githubusercontent.com/beezly/unifi-apis/main/unifi-network/{version}.json`

## Release Workflow

To release an SDK for a new API version, use GitHub Actions `workflow_dispatch`.

### Manual Release Steps

1. Open the "Actions" tab in the GitHub repository
2. Select the "Release SDK" workflow
3. Click "Run workflow"
4. Enter the UniFi API version (e.g. `v9.1.120`) in `api_version`
5. Click "Run workflow" to execute

### Workflow Steps

1. Determine the next free SDK version by starting from [VERSION](VERSION) and incrementing until both tag and branch names are available
2. Create a release branch named `release/sdk-vX.Y.Z-network-vA.B.C`
3. Download the OpenAPI spec into `openapi/unifi-network/<version>.json`
4. If `unifi-network-version` contains a valid version, compute OpenAPI diffs via `oasdiff` and include them in the release notes
5. Generate `pkg/network/openapi.gen.go` from the committed spec via `oapi-codegen`
6. Format, build, and test code
7. Update `unifi-network-version` to the released API version
8. Commit and push the spec, generated code, and metadata to both the release branch and `main`
9. Create a version tag (SDK version)
10. Create a GitHub release

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
    "https://raw.githubusercontent.com/beezly/unifi-apis/main/unifi-network/10.0.162.json" \
    -o ./openapi/unifi-network/10.0.162.json

go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
    -generate types,client \
    -package network \
    -o ./pkg/network/openapi.gen.go \
    ./openapi/unifi-network/10.0.162.json
```

## License

MIT
