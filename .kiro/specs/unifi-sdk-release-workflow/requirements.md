# UniFi SDK Release Workflow - 要件定義

## 概要
UniFi APIのバージョンに対応したGo SDKを自動生成・リリースするワークフローを構築する。

## 背景
- 既存の `typegen` ツールでUniFi APIドキュメントからGo型を生成可能
- APIバージョンごとにSDKをリリースし、Goモジュールとして利用可能にしたい
- newrelic-client-goのファイル構成を参考にする

## ユーザーストーリー

### US-1: APIバージョン指定リリース
**As a** SDK開発者
**I want to** UniFi APIバージョンを指定してSDKをリリースできる
**So that** 特定のAPIバージョンに対応したSDKを提供できる

#### 受け入れ条件
- [ ] 1.1 workflow_dispatchでAPIバージョン（例: `v9.1.120`）を入力できる
- [ ] 1.2 入力されたバージョンがGoモジュールのバージョンタグとして使用される
- [ ] 1.3 リリースブランチ（例: `release/v9.1.120`）が自動作成される

### US-2: 型生成
**As a** SDK開発者
**I want to** UniFi APIドキュメントから自動的にGo型を生成できる
**So that** 手動での型定義作業を削減できる

#### 受け入れ条件
- [ ] 2.1 指定されたAPIバージョンのドキュメントURLからエンドポイントを発見できる
- [ ] 2.2 各エンドポイントのRequest/Response型が生成される
- [ ] 2.3 共通型（Voucher, Client, Site等）は `pkg/common/` に配置される
- [ ] 2.4 ドメイン別に `pkg/{category}/types.go` に型が配置される

### US-3: Clientメソッド生成
**As a** SDK利用者
**I want to** 生成されたClientメソッドでAPIを呼び出せる
**So that** 型安全にUniFi APIを利用できる

#### 受け入れ条件
- [ ] 3.1 各エンドポイントに対応するClientメソッドが生成される
- [ ] 3.2 メソッド名はAPIエンドポイント名をそのまま使用する
- [ ] 3.3 Path Parametersは関数引数として受け取る
- [ ] 3.4 Request Bodyがある場合は構造体として受け取る
- [ ] 3.5 Response型を返却する
- [ ] 3.6 エラーハンドリングが適切に行われる

### US-4: メインClient構造体
**As a** SDK利用者
**I want to** 単一のClientインスタンスから全APIにアクセスできる
**So that** 簡潔にSDKを利用できる

#### 受け入れ条件
- [ ] 4.1 `unifi/unifi.go` にメインClient構造体が定義される
- [ ] 4.2 各ドメインのサブClientが埋め込まれる（Clients, Sites, Vouchers等）
- [ ] 4.3 `unifi.New(config)` でClientを初期化できる
- [ ] 4.4 認証情報（API Key等）を設定できる

### US-5: GitHub Actionsワークフロー
**As a** SDK開発者
**I want to** GitHub Actionsでリリースを自動化できる
**So that** 手動作業を最小化できる

#### 受け入れ条件
- [ ] 5.1 workflow_dispatchトリガーで実行できる
- [ ] 5.2 APIバージョンを入力パラメータとして受け取る
- [ ] 5.3 型生成・Clientメソッド生成が実行される
- [ ] 5.4 go.modのバージョンが更新される
- [ ] 5.5 変更がコミット・プッシュされる
- [ ] 5.6 GitHubリリースが作成される
- [ ] 5.7 バージョンタグが付与される

### US-6: Goモジュールとしての利用
**As a** SDK利用者
**I want to** `go get` でSDKをインストールできる
**So that** 標準的なGoワークフローで利用できる

#### 受け入れ条件
- [ ] 6.1 `go get github.com/murasame29/unifi-go-sdk@v9.1.120` でインストールできる
- [ ] 6.2 標準ライブラリのみに依存する（外部依存なし）
- [ ] 6.3 適切なgo.modが生成される

## 非機能要件

### NFR-1: 依存関係
- 標準ライブラリのみ使用（`net/http`, `encoding/json`, `context`等）
- 外部パッケージへの依存は禁止

### NFR-2: ファイル構成
newrelic-client-goを参考にした構成:
```
unifi-go-sdk/
├── unifi/
│   └── unifi.go           # メインClient
├── pkg/
│   ├── common/            # 共通型
│   ├── clients/           # Clients API
│   ├── sites/             # Sites API
│   ├── vouchers/          # Vouchers API
│   ├── devices/           # Devices API
│   └── config/            # 設定
├── internal/
│   ├── http/              # HTTPクライアント実装
│   └── typegen/           # 型生成ツール（既存）
└── cmd/
    └── typegen/           # CLI（既存）
```

### NFR-3: エラーハンドリング
- `fmt.Errorf("%w", err)` でエラーをラップ
- カスタムエラー型を `pkg/errors/` に定義

### NFR-4: Context対応
- 全てのAPIメソッドは第一引数に `ctx context.Context` を受け取る

## 制約事項
- Cloud Connector APIは除外する
- ドキュメントページ（Overview等）は除外する
