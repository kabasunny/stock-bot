# Stock Bot

このプロジェクトは、立花証券 e-支店 API を利用して株式の自動取引を行うためのGo言語製ボットです。

## 概要

指定されたロジックに基づき、株式の売買注文を自動的に実行します。APIとの通信、注文管理、ポジション管理などの機能を含みます。

## アーキテクチャ

このプロジェクトは、関心事の分離を目的としたレイヤードアーキテクチャ（クリーンアーキテクチャに類似）を採用しています。

-   `cmd/`: アプリケーションのエントリーポイントです。`main.go` が含まれ、依存関係の注入やアプリケーションの起動を行います。
-   `domain/`: コアとなるビジネスロジックとエンティティ（ビジネスオブジェクト）を定義します。
    -   `model/`: `Order` (注文), `Position` (建玉), `Stock` (株マスター) といった、システムの核となるデータ構造を定義します。
    -   `repository/`: データの永続化（取得・保存）に関するインターフェースを定義します。実際のデータベースやAPIとの通信方法は隠蔽されます。
-   `internal/`: このプロジェクト内部でのみ使用されるコードを配置します。
    -   `app/`: ユースケース層です。`domain` のエンティティやロジックを使い、具体的なアプリケーションの機能（例：「残高を確認する」）を実現します。
    -   `config/`: 設定ファイル (`.env`) の読み込みや管理を行います。
    -   `infrastructure/`: データベース、外部APIクライアントなど、外部システムとの接続に関する具体的な実装を配置します。
        -   `client/`: 立花証券APIとの通信を行うクライアントの実装です。
        -   `repository/`: `domain/repository` で定義されたインターフェースを実装し、データベースへのアクセスなどを具体的に行います。
    -   `interface/`: ユーザーや外部システムとの接点となる部分です。
        -   `web/`: HTTPサーバーのハンドラなどを定義します。
        -   `agent/`: 取引戦略（エージェント）の具体的な実装を配置します。
    -   `logger/`: アプリケーション全体で利用するロガーの設定を行います。

## セットアップ方法

### 1. 前提条件

-   Go 1.21 以上

### 2. 設定

プロジェクトのルートディレクトリに `.env` という名前のファイルを作成し、以下の内容を参考に設定を記述してください。

```.env
# .env.example

# Tachibana Securities API Settings
TACHIBANA_BASE_URL="https://demo-kabuka.e-shiten.jp/e_api_v4r6/" # デモ環境 or 本番環境
TACHIBANA_USER_ID="YOUR_USER_ID"
TACHIBANA_PASSWORD="YOUR_PASSWORD"

# Database Settings (Optional)
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="user"
DB_PASSWORD="password"
DB_NAME="stockbot_db"

# Log Level (debug, info, warn, error)
LOG_LEVEL="debug"

# HTTP Server Port
HTTP_PORT="8080"
```

### 3. 依存関係のインストール

```sh
go mod tidy
```

## 実行

### アプリケーションの起動

```sh
go run ./cmd/myapp/main.go
```

### テストの実行

```sh
go test -v ./...
```