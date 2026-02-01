# UniFi Go SDK

UniFi API用のGo SDKです。UniFi APIドキュメントから自動生成された型とClientメソッドを提供します。

## 特徴

- **自動生成**: UniFi APIドキュメントから型とClientメソッドを自動生成
- **型安全**: 全てのAPIリクエスト/レスポンスに型定義
- **標準ライブラリのみ**: 外部依存なし（`net/http`, `encoding/json`, `context`等のみ使用）
- **Context対応**: 全てのAPIメソッドは `ctx context.Context` を第一引数に受け取る

## インストール

```bash
go get github.com/murasame29/unifi-go-sdk@v9.1.120
```

バージョンはUniFi APIバージョンに対応しています。利用可能なバージョンは[Releases](https://github.com/murasame29/unifi-go-sdk/releases)を確認してください。

## 使用方法

### Site Manager API（クラウドAPI）

Site Manager APIはUniFiクラウドサービスを通じてデバイスを管理するためのAPIです。APIキーによる認証が必要です。

```go
package main

import (
    "context"
    "log"

    "github.com/murasame29/unifi-go-sdk/unifi"
)

func main() {
    // Clientの初期化
    client, err := unifi.New(
        unifi.ConfigAPIKey("your-api-key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Site Manager APIを使用
    // client.SiteManager.XXX(ctx, ...)
    _ = client
    _ = ctx
}
```

### Network API（ローカルコントローラー）

Network APIはローカルのUniFiコントローラー（UDM、Cloud Key等）と直接通信するためのAPIです。

```go
package main

import (
    "context"
    "log"

    "github.com/murasame29/unifi-go-sdk/unifi"
    "github.com/murasame29/unifi-go-sdk/pkg/network"
)

func main() {
    // Network Clientの初期化
    client, err := unifi.NewNetwork(network.Config{
        BaseURL:            "https://192.168.1.1:8443",
        Site:               "default",
        InsecureSkipVerify: true, // 自己署名証明書の場合
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Network APIを使用
    // client.XXX(ctx, ...)
    _ = client
    _ = ctx
}
```

## ディレクトリ構造

```
unifi-go-sdk/
├── .github/
│   └── workflows/
│       └── release.yml          # リリースワークフロー
├── unifi/
│   └── unifi.go                 # メインClient
├── pkg/
│   ├── config/                  # 設定
│   ├── errors/                  # カスタムエラー
│   ├── network/                 # Network API（自動生成）
│   │   ├── types.go             # 型定義
│   │   └── network.go           # Clientメソッド
│   └── sitemanager/             # Site Manager API
├── internal/
│   ├── http/                    # HTTPクライアント
│   ├── typegen/                 # 型生成ツール
│   └── clientgen/               # Clientメソッド生成ツール
├── cmd/
│   ├── typegen/                 # 型生成CLI
│   └── clientgen/               # Clientメソッド生成CLI
├── go.mod
└── README.md
```

## 自動生成について

このSDKの型とClientメソッドは、UniFi APIドキュメントから自動生成されています。

### 生成される型

- **Request型**: APIリクエストのボディ構造体
- **Response型**: APIレスポンスの構造体
- **共通型**: Voucher, Client, Site等の共通エンティティ

### 生成されるClientメソッド

各APIエンドポイントに対応するメソッドが生成されます：

| HTTPメソッド | エンドポイント名 | 生成されるメソッド |
|-------------|-----------------|-------------------|
| GET (list)  | List Clients    | ListClients       |
| GET (single)| Get Client      | GetClient         |
| POST        | Create Voucher  | CreateVoucher     |
| PUT         | Update Device   | UpdateDevice      |
| DELETE      | Delete Voucher  | DeleteVoucher     |

## リリースワークフロー

新しいAPIバージョンのSDKをリリースするには、GitHub Actionsの`workflow_dispatch`を使用します。

### 手動リリース手順

1. GitHubリポジトリの「Actions」タブを開く
2. 「Release SDK」ワークフローを選択
3. 「Run workflow」をクリック
4. `api_version`にUniFi APIバージョン（例: `v9.1.120`）を入力
5. 「Run workflow」をクリックして実行

### ワークフローの処理内容

1. リリースブランチ（`release/vX.Y.Z`）を作成
2. UniFi APIドキュメントから型を生成（`typegen`）
3. Clientメソッドを生成（`clientgen`）
4. コードをフォーマット・ビルド・テスト
5. 変更をコミット・プッシュ
6. バージョンタグを作成
7. GitHubリリースを作成

## エラーハンドリング

SDKは`pkg/errors`パッケージでカスタムエラー型を提供しています：

```go
import "github.com/murasame29/unifi-go-sdk/pkg/errors"

// APIエラーのチェック
if errors.Is(err, errors.ErrUnauthorized) {
    // 認証エラー
}

if errors.Is(err, errors.ErrNotFound) {
    // リソースが見つからない
}

// APIErrorの詳細を取得
var apiErr *errors.APIError
if errors.As(err, &apiErr) {
    log.Printf("Status: %d, Message: %s", apiErr.StatusCode, apiErr.Message)
}
```

## 設定オプション

### Site Manager API

```go
client, err := unifi.New(
    unifi.ConfigAPIKey("your-api-key"),           // APIキー（必須）
    unifi.ConfigBaseURL("https://api.ui.com"),    // ベースURL（オプション）
    unifi.ConfigUserAgent("my-app/1.0"),          // User-Agent（オプション）
)
```

### Network API

```go
client, err := unifi.NewNetwork(network.Config{
    BaseURL:            "https://192.168.1.1:8443", // コントローラーURL（必須）
    Site:               "default",                   // サイト名（デフォルト: "default"）
    Timeout:            30 * time.Second,            // タイムアウト（デフォルト: 30秒）
    InsecureSkipVerify: true,                        // TLS検証スキップ（自己署名証明書用）
})
```

## 開発

### 型生成ツールの実行

```bash
go run ./cmd/typegen/main.go \
    -discover "https://developer.ui.com/network/v9.1.120" \
    -output-dir ./pkg \
    -package network \
    -workers 4
```

### Clientメソッド生成ツールの実行

```bash
go run ./cmd/clientgen/main.go \
    -input ./pkg \
    -output ./pkg
```

## License

MIT
