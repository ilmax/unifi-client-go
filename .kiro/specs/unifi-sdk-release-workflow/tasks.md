# UniFi SDK Release Workflow - タスクリスト

## Phase 1: 基盤整備

- [x] 1. pkg/config パッケージの作成
  - [x] 1.1 `pkg/config/config.go` - Config構造体とNew関数
  - [x] 1.2 ConfigOption パターンの実装

- [x] 2. pkg/errors パッケージの作成
  - [x] 2.1 `pkg/errors/errors.go` - カスタムエラー型定義
  - [x] 2.2 APIError, ValidationError等の実装

- [x] 3. internal/http パッケージの作成
  - [x] 3.1 `internal/http/client.go` - HTTPクライアント実装
  - [x] 3.2 Get, Post, Put, Delete メソッド
  - [x] 3.3 認証ヘッダー設定
  - [x] 3.4 エラーハンドリング

## Phase 2: 型生成拡張

- [x] 4. typegen出力先の変更
  - [x] 4.1 カテゴリ分類ロジックの実装（endpoint URLからカテゴリを抽出）
  - [x] 4.2 `pkg/{category}/types.go` への出力対応
  - [x] 4.3 共通型を `pkg/common/types.go` に出力

- [x] 5. APISchema拡張
  - [x] 5.1 Category フィールドの追加
  - [x] 5.2 スクレイピング時にカテゴリを設定

## Phase 3: Clientメソッド生成

- [x] 6. internal/clientgen パッケージの作成
  - [x] 6.1 `internal/clientgen/generator.go` - メイン生成ロジック
  - [x] 6.2 APISchemaからメソッドコードを生成
  - [x] 6.3 メソッド名生成ルール（ListXxx, GetXxx, CreateXxx等）
  - [x] 6.4 Path Parameter置換ロジック

- [x] 7. cmd/clientgen CLIの作成
  - [x] 7.1 `cmd/clientgen/main.go` - CLIエントリポイント
  - [x] 7.2 入力ディレクトリ指定オプション
  - [x] 7.3 出力ディレクトリ指定オプション

## Phase 4: メインClient

- [x] 8. unifi パッケージの作成
  - [x] 8.1 `unifi/unifi.go` - UniFi構造体
  - [x] 8.2 New関数とConfigOption
  - [x] 8.3 各サブClientの初期化

## Phase 5: GitHub Actions

- [x] 9. リリースワークフローの作成
  - [x] 9.1 `.github/workflows/release.yml` - ワークフロー定義
  - [x] 9.2 workflow_dispatch入力パラメータ
  - [x] 9.3 Chrome/Chromiumセットアップ
  - [x] 9.4 型生成ステップ
  - [x] 9.5 Clientメソッド生成ステップ
  - [x] 9.6 ビルド・テストステップ
  - [x] 9.7 コミット・プッシュステップ
  - [x] 9.8 タグ作成ステップ
  - [x] 9.9 GitHubリリース作成ステップ

## Phase 6: ドキュメント・仕上げ

- [x] 10. ドキュメント整備
  - [x] 10.1 README.md更新（使用方法、インストール方法）
  - [x] 10.2 使用例の追加

- [x] 11. テスト
  - [x] 11.1 pkg/config のユニットテスト
  - [x] 11.2 internal/http のユニットテスト
  - [x] 11.3 生成コードのビルド確認
