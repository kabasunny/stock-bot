# Stock Trading Bot

このプロジェクトは、証券会社APIを利用して株式の自動取引を行うための取引ボットです。

## 概要
柔軟な取引戦略を実装できる「エージェント」を中心に据えた、拡張性の高い自動取引システムです。証券会社APIとの通信、シグナルに基づいた注文執行、状態管理などの機能を提供します。

## アーキテクチャ (エージェント中心モデル)
このシステムは、責務が明確に分離された3つのコンポーネントで構成されます。

1.  **Go製 APIラッパー**
    -   **役割**: 証券会社APIとの通信をすべて担当し、その複雑さを抽象化します。
    -   **機能**: Goaフレームワークを用いて構築されたHTTP APIを公開し、外部からのリクエストに応じて注文執行やデータ取得を行います。
    -   **ディレクトリ**: `internal/infrastructure/client`, `internal/handler/web` など

2.  **Python製 シグナル生成サービス**
    -   **役割**: 高度な計算や機械学習モデルを用いて、売買シグナル（BUY/SELL/HOLD）を生成します。
    -   **機能**: 外部から市場データを受け取り、分析結果のシグナルを返すHTTPサーバーとして動作します。（このリポジトリ外で管理）

3.  **エージェント**
    -   **役割**: システム全体の「頭脳」として、主体的に意思決定を行う中央司令塔です。
    -   **機能**: 取引戦略のメインループとして、Go APIラッパー（データ取得/注文）とPythonサービス（シグナル問合せ）を呼び出し、全体のワークフローを指揮します。
    -   **ディレクトリ**: `internal/agent`

### 実装計画

-   **短期計画**: まず、Go APIラッパーとエージェントを単一のGoアプリケーションとして実装します。コンポーネント間の通信もGoaで定義したHTTP API (`localhost`経由) で行い、将来のサービス分割を見据えた疎結合な設計を維持します。
-   **長期目標**: システム安定稼働後、エージェントを**Rust製のマイクロサービス**として再実装し、パフォーマンスと安全性を極限まで高めることを目指します。

## ディレクトリ構造

このプロジェクトは、クリーンアーキテクチャと思想を参考に、責務に基づいたディレクトリ分割を行っています。

-   **/cmd**: アプリケーションのエントリーポイント（起動スクリプト）が配置されます。
    -   `main.go` は、依存関係の注入（DI）やサーバーの起動など、アプリケーション全体の初期化処理を担当します。

-   **/design**: GoaフレームワークのAPI設計ファイル (`design.go`) が配置されます。
    -   APIのエンドポイント、リクエスト/レスポンスのデータ構造などを定義します。`goa gen` コマンドは、この設計に基づいて `/gen` ディレクトリにコードを自動生成します。

-   **/domain**: システムの核となるビジネスルールとデータ構造（ドメインモデル）を定義します。
    -   **/model**: `Order` や `Position` といった、ビジネス上最も重要な概念を表現する構造体を定義します。
    -   **/repository**: データベースなどへの永続化処理のインターフェース（`OrderRepository`など）を定義します。具体的な実装は `infrastructure` 層が担当します。

-   **/gen**: Goaによって自動生成されたコードが格納されます。**このディレクトリ以下のファイルは直接編集しないでください。**

-   **/internal**: このプロジェクト内部でのみ使用されるGoパッケージを配置します。
    -   **/agent**: 取引戦略を実行する「エージェント」のロジックを実装します。システムの「頭脳」にあたる部分です。
    -   **/app**: ユースケース層です。アプリケーション固有のビジネスロジック（例: 「注文を実行する」）を実装します。ハンドラ層からの指示を受け、ドメインモデルを使って処理を行います。
    -   **/handler**: 外部からの入力を受け付ける層です。
        -   **/web**: Goaサービスの実装など、HTTPリクエストを処理するハンドラを配置します。リクエストを解釈し、対応するユースケースを呼び出します。
    -   **/infrastructure**: データベース、外部API（証券会社など）との通信といった、技術的な詳細を実装する層です。
        -   **/client**: 証券会社APIを呼び出すクライアントの実装などを配置します。
        -   **/repository**: `domain/repository` で定義されたインターフェースの具体的な実装（例: GORMを使ったDB操作）を配置します。

-   **/planning**: `SYSTEM_DESIGN_MEMO.md` などの設計関連ドキュメントを格納します。

## セットアップ方法

### 1. 前提条件

- Go 1.21 以上
- `goa` v3 コマンド

### 2. 設定

プロジェクトのルートディレクトリに `.env` という名前のファイルを作成し、以下の内容を参考に設定を記述してください。

```.env
# .env.example

# Tachibana Securities API Settings
TACHIBANA_BASE_URL="https://demo-kabuka.e-shiten.jp/e_api_v4r6/"
TACHIBANA_USER_ID="YOUR_USER_ID"
TACHIBANA_PASSWORD="YOUR_PASSWORD"

# WebSocket Event Settings
EVENT_RID=""
EVENT_BOARD_NO=""
EVENT_NO=""
EVENT_EVT_CMD=""

# Database Settings (Optional)
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="user"
DB_PASSWORD="password"
DB_NAME="stockbot_db"

# Log Level (debug, info, warn, error)
LOG_LEVEL="debug"

# HTTP Server Port (for Goa API)
HTTP_PORT="8080"
```

### 3. 依存関係のインストール

```sh
go mod tidy
```

## データベース (Database)

### マイグレーション

このプロジェクトは `golang-migrate/migrate` を使用してデータベースのスキーマを管理します。
アプリケーションを起動する前に、データベースのスキーマをセットアップまたは更新する必要があります。

1.  **`migrate` CLIのインストール (初回のみ)**
    ```sh
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

2.  **マイグレーションの適用**
    以下のコマンドを実行して、最新のスキーマをデータベースに適用します。
    コマンド内の `postgres://user:password@host:port/dbname` の部分は、ご自身の `.env` ファイルの内容に合わせて書き換えてください。

    ```sh
    # 例: migrate -database "postgres://user:password@localhost:5432/stockbot_db?sslmode=disable" -path ./migrations up
    migrate -database "YOUR_DATABASE_CONNECTION_STRING" -path ./migrations up
    ```

    マイグレーションを1つ前のバージョンに戻す場合は `down 1` を使用します。
    ```sh
    migrate -database "YOUR_DATABASE_CONNECTION_STRING" -path ./migrations down 1
    ```

## 実行

### アプリケーションの起動

```sh
go run ./cmd/myapp/main.go
```

### APIのテスト

サーバー起動後、以下のコマンドで注文APIをテストできます。

```sh
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"symbol": "7203", "trade_type": "BUY", "order_type": "MARKET", "quantity": 100}' \
  http://localhost:8080/order
```

### テストの実行

```sh
go test -v ./...
```