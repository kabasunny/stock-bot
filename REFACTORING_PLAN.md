# エージェント/クライアント分離 リファクタリング計画

## 1. 目的

`agent` コンポーネントと `client` コンポーネントを分離し、コードの保守性とテスト容易性を向上させる。

## 2. 背景と問題点

現状のコードベースでは、以下の問題点が確認されています。

1.  **密結合**: `TradeService` の実装である `GoaTradeService` がエージェントのロジック (`internal/agent`) 内に存在し、インフラ層と密結合している。
2.  **責務の混在**: エージェントが、データベースへの永続化を行う `ExecutionUseCase` を直接呼び出しており、責務が適切に分離されていない。

## 3. リファクタリング手順

以下の手順でリファクタリングを実施します。

### ステップ1: `TradeService` の再配置

1.  **`GoaTradeService` の移動**:
    -   `internal/agent/goa_trade_service.go` を `internal/tradeservice/goa_trade_service.go` に移動する。
2.  **`TradeService` インターフェースの移動**:
    -   `internal/agent/trade_service.go` 内のインターフェース定義を `domain/service/trade_service.go` に移動する。
3.  **インポートパスの修正**:
    -   移動に伴い、関連ファイル (`agent.go`, `main.go` 等) のインポートパスを更新する。

### ステップ2: `ExecutionUseCase` の依存関係逆転

1.  **`ExecutionEventHandler` インターフェースの作成**:
    -   `internal/agent/` パッケージに、約定イベントを処理するための新しいインターフェース `ExecutionEventHandler` を定義する。
2.  **`Agent` の修正**:
    -   `internal/agent/agent.go` を修正し、`ExecutionUseCase` の直接呼び出しを `ExecutionEventHandler` インターフェースの呼び出しに置き換える。
3.  **アダプターの作成**:
    -   `internal/app/` パッケージに、`ExecutionEventHandler` インターフェースを実装するアダプター (`execution_event_adapter.go`) を作成する。このアダプターが内部で `ExecutionUseCase` を呼び出す。

### ステップ3: DIコンテナの更新

1.  **`main.go` の更新**:
    -   `cmd/myapp/main.go` を修正し、新しいパッケージ構成とインターフェースに基づいて、依存性の注入（Dependency Injection）を正しく設定し直す。

### ステップ4: 検証

1.  **ビルドとテスト**:
    -   `go build ./...` と `go test ./...` を実行し、すべてのコンパイルエラーとテスト失敗が解消されていることを確認する。

## 4. 状況更新 (Status Update)

リファクタリングの途中で、ビルド時に `package stock-bot/... is not in std` というエラーが繰り返し発生し、Goモジュールパスの解決に問題があることが判明しました。これは、`go.mod`の`module stock-bot`という宣言がGoの標準ライブラリパスと競合している可能性が原因です。

この問題を解決するため、`go.mod`のモジュール名をより一意な`github.com/kabas/stock-bot`に変更しました。しかし、`go mod tidy`はソースコード内のインポートパスを自動的に変更しないため、現在、プロジェクト内のすべてのGoファイルのインポートパスを、新しいモジュール名に合わせて一括で置換する作業を進めています。

この一括置換が完了した後、再度ビルドとテストを実行し、リファクタリングが成功したことを検証します。