# UniFi SDK Release Workflow - 設計書

## アーキテクチャ概要

```
┌─────────────────────────────────────────────────────────────┐
│                    GitHub Actions                            │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  workflow_dispatch (api_version: v9.1.120)          │    │
│  └─────────────────────────────────────────────────────┘    │
│                           │                                  │
│                           ▼                                  │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  1. Create release branch                           │    │
│  │  2. Run typegen (generate types)                    │    │
│  │  3. Run clientgen (generate client methods)         │    │
│  │  4. Update go.mod                                   │    │
│  │  5. Commit & Push                                   │    │
│  │  6. Create Release & Tag                            │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## コンポーネント設計

### 1. 型生成ツール (typegen) - 既存拡張

#### 1.1 APISchema拡張
```go
// internal/typegen/types.go
type APISchema struct {
    Endpoint    string
    Method      string      // GET, POST, PUT, DELETE
    Path        string      // /api/v1/sites/{siteId}/clients
    Description string
    PathParams  []Property
    Request     *SchemaObject
    Response    *SchemaObject
    Category    string      // clients, sites, vouchers, devices
}
```

#### 1.2 出力先の変更
- 現在: `tmp/` に全ファイル出力
- 変更後: `pkg/{category}/types.go` に出力

### 2. Clientメソッド生成ツール (clientgen) - 新規

#### 2.1 生成されるコード例
```go
// pkg/clients/clients.go
package clients

import (
    "context"
    "fmt"
    "net/http"
)

// Clients provides access to the Clients API.
type Clients struct {
    client *http.Client
    baseURL string
    apiKey string
}

// New creates a new Clients instance.
func New(client *http.Client, baseURL, apiKey string) *Clients {
    return &Clients{
        client:  client,
        baseURL: baseURL,
        apiKey:  apiKey,
    }
}

// ListClients retrieves all clients for a site.
// GET /api/v1/sites/{siteId}/clients
func (c *Clients) ListClients(ctx context.Context, siteId string) (*ListClientsResponse, error) {
    path := fmt.Sprintf("/api/v1/sites/%s/clients", siteId)
    // ... HTTP request implementation
}

// GetClient retrieves a specific client.
// GET /api/v1/sites/{siteId}/clients/{clientId}
func (c *Clients) GetClient(ctx context.Context, siteId, clientId string) (*Client, error) {
    path := fmt.Sprintf("/api/v1/sites/%s/clients/%s", siteId, clientId)
    // ... HTTP request implementation
}
```

#### 2.2 メソッド名生成ルール
| HTTP Method | Endpoint Name | Generated Method |
|-------------|---------------|------------------|
| GET (list)  | List Clients  | ListClients      |
| GET (single)| Get Client    | GetClient        |
| POST        | Create Voucher| CreateVoucher    |
| PUT         | Update Device | UpdateDevice     |
| DELETE      | Delete Voucher| DeleteVoucher    |

### 3. メインClient (unifi/unifi.go)

```go
// unifi/unifi.go
package unifi

import (
    "net/http"
    "time"

    "github.com/murasame29/unifi-go-sdk/pkg/clients"
    "github.com/murasame29/unifi-go-sdk/pkg/config"
    "github.com/murasame29/unifi-go-sdk/pkg/devices"
    "github.com/murasame29/unifi-go-sdk/pkg/sites"
    "github.com/murasame29/unifi-go-sdk/pkg/vouchers"
)

// UniFi is a collection of UniFi APIs.
type UniFi struct {
    Clients  *clients.Clients
    Sites    *sites.Sites
    Vouchers *vouchers.Vouchers
    Devices  *devices.Devices

    config *config.Config
}

// New returns a new UniFi client.
func New(opts ...ConfigOption) (*UniFi, error) {
    cfg := config.New()
    if err := cfg.Apply(opts...); err != nil {
        return nil, err
    }

    httpClient := &http.Client{
        Timeout: cfg.Timeout,
    }

    return &UniFi{
        Clients:  clients.New(httpClient, cfg.BaseURL, cfg.APIKey),
        Sites:    sites.New(httpClient, cfg.BaseURL, cfg.APIKey),
        Vouchers: vouchers.New(httpClient, cfg.BaseURL, cfg.APIKey),
        Devices:  devices.New(httpClient, cfg.BaseURL, cfg.APIKey),
        config:   cfg,
    }, nil
}

// ConfigOption configures the UniFi client.
type ConfigOption func(*config.Config) error

// ConfigAPIKey sets the API key.
func ConfigAPIKey(apiKey string) ConfigOption {
    return func(c *config.Config) error {
        c.APIKey = apiKey
        return nil
    }
}

// ConfigBaseURL sets the base URL.
func ConfigBaseURL(baseURL string) ConfigOption {
    return func(c *config.Config) error {
        c.BaseURL = baseURL
        return nil
    }
}

// ConfigTimeout sets the HTTP timeout.
func ConfigTimeout(timeout time.Duration) ConfigOption {
    return func(c *config.Config) error {
        c.Timeout = timeout
        return nil
    }
}
```

### 4. 設定 (pkg/config/config.go)

```go
// pkg/config/config.go
package config

import "time"

// Config holds the configuration for the UniFi client.
type Config struct {
    APIKey  string
    BaseURL string
    Timeout time.Duration
}

// New creates a new Config with default values.
func New() *Config {
    return &Config{
        Timeout: 30 * time.Second,
    }
}

// Apply applies the given options to the config.
func (c *Config) Apply(opts ...func(*Config) error) error {
    for _, opt := range opts {
        if err := opt(c); err != nil {
            return err
        }
    }
    return nil
}
```

### 5. HTTPクライアント (internal/http/client.go)

```go
// internal/http/client.go
package http

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// Client wraps the standard http.Client with UniFi-specific functionality.
type Client struct {
    httpClient *http.Client
    baseURL    string
    apiKey     string
}

// New creates a new HTTP client.
func New(httpClient *http.Client, baseURL, apiKey string) *Client {
    return &Client{
        httpClient: httpClient,
        baseURL:    baseURL,
        apiKey:     apiKey,
    }
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, result any) error {
    return c.do(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body, result any) error {
    return c.do(ctx, http.MethodPost, path, body, result)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body, result any) error {
    return c.do(ctx, http.MethodPut, path, body, result)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) error {
    return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
    url := c.baseURL + path

    var bodyReader io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("failed to marshal request body: %w", err)
        }
        bodyReader = bytes.NewReader(jsonBody)
    }

    req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("API error: status %d", resp.StatusCode)
    }

    if result != nil {
        if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
            return fmt.Errorf("failed to decode response: %w", err)
        }
    }

    return nil
}
```

### 6. GitHub Actions Workflow

```yaml
# .github/workflows/release.yml
name: Release SDK

on:
  workflow_dispatch:
    inputs:
      api_version:
        description: 'UniFi API Version (e.g., v9.1.120)'
        required: true
        type: string

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Setup Chrome
        uses: browser-actions/setup-chrome@v1

      - name: Create release branch
        run: |
          git checkout -b release/${{ inputs.api_version }}

      - name: Generate types
        run: |
          go run ./cmd/typegen/main.go \
            -discover "https://developer.ui.com/network/${{ inputs.api_version }}" \
            -output-dir ./pkg \
            -package network \
            -workers 4

      - name: Generate client methods
        run: |
          go run ./cmd/clientgen/main.go \
            -input ./pkg \
            -output ./pkg

      - name: Update go.mod version
        run: |
          # Update module version if needed

      - name: Build and test
        run: |
          go build ./...
          go test ./...

      - name: Commit changes
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add .
          git commit -m "Release ${{ inputs.api_version }}"
          git push origin release/${{ inputs.api_version }}

      - name: Create tag
        run: |
          git tag ${{ inputs.api_version }}
          git push origin ${{ inputs.api_version }}

      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          tag_name: ${{ inputs.api_version }}
          name: Release ${{ inputs.api_version }}
          body: |
            UniFi SDK for API version ${{ inputs.api_version }}
            
            ## Installation
            ```bash
            go get github.com/murasame29/unifi-go-sdk@${{ inputs.api_version }}
            ```
          draft: false
          prerelease: false
```

## ファイル構成（最終形）

```
unifi-go-sdk/
├── .github/
│   └── workflows/
│       └── release.yml          # リリースワークフロー
├── unifi/
│   └── unifi.go                 # メインClient
├── pkg/
│   ├── common/
│   │   └── types.go             # 共通型（Voucher, Client, Site等）
│   ├── config/
│   │   └── config.go            # 設定
│   ├── errors/
│   │   └── errors.go            # カスタムエラー
│   ├── clients/
│   │   ├── types.go             # 型定義（生成）
│   │   └── clients.go           # Clientメソッド（生成）
│   ├── sites/
│   │   ├── types.go
│   │   └── sites.go
│   ├── vouchers/
│   │   ├── types.go
│   │   └── vouchers.go
│   └── devices/
│       ├── types.go
│       └── devices.go
├── internal/
│   ├── http/
│   │   └── client.go            # HTTPクライアント
│   └── typegen/                 # 型生成ツール（既存）
├── cmd/
│   ├── typegen/                 # 型生成CLI（既存）
│   └── clientgen/               # Clientメソッド生成CLI（新規）
├── go.mod
└── README.md
```

## 実装フェーズ

### Phase 1: 基盤整備
1. `pkg/config/` - 設定パッケージ
2. `pkg/errors/` - エラーパッケージ
3. `internal/http/` - HTTPクライアント

### Phase 2: 型生成拡張
1. typegen出力先を `pkg/{category}/types.go` に変更
2. カテゴリ分類ロジックの実装

### Phase 3: Clientメソッド生成
1. `cmd/clientgen/` - CLI作成
2. `internal/clientgen/` - 生成ロジック

### Phase 4: メインClient
1. `unifi/unifi.go` - メインClient実装

### Phase 5: GitHub Actions
1. `.github/workflows/release.yml` - ワークフロー作成

## 正当性プロパティ

### P-1: 型生成の正当性
- 生成された型がGoの構文として有効である
- 全てのフィールドに適切なJSONタグが付与される

### P-2: Clientメソッドの正当性
- 全てのメソッドが `ctx context.Context` を第一引数に持つ
- Path Parametersが正しく関数引数にマッピングされる
- HTTPメソッドが正しく設定される

### P-3: エラーハンドリング
- 全てのエラーが適切にラップされる
- HTTPステータスコード >= 400 でエラーが返される
