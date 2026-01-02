# 株式取引システム (Stock Trading System)

立花証券APIを使用した自動株式取引システムです。Clean Architectureに基づく設計で、高い拡張性と保守性を実現しています。

## 🎯 主要機能

- **自動取引**: 戦略に基づく自動売買
- **リアルタイム監視**: WebSocketによる価格・約定監視
- **注文管理**: 成行・指値・逆指値注文対応
- **ポートフォリオ管理**: ポジション・残高管理
- **バックテスト**: 過去データでの戦略検証
- **REST API**: HTTP APIによる外部連携

## 🏗️ アーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│                    株式取引システム                          │
├─────────────────────────────────────────────────────────────┤
│  Web API     │  Trading Bot  │  Backtester                 │
│  (HTTP/REST) │  (Agent)      │  (Analysis)                 │
├─────────────────────────────────────────────────────────────┤
│              Application Layer                              │
│  Trade Service │ Event Handler │ State Manager              │
├─────────────────────────────────────────────────────────────┤
│             Infrastructure Layer                            │
│  Tachibana API │ Database      │ WebSocket                  │
├─────────────────────────────────────────────────────────────┤
│               External APIs                                 │
│  立花証券 API   │ Market Data   │ Price Feed                 │
└─────────────────────────────────────────────────────────────┘
```

詳細は [システムアーキテクチャ概要](SYSTEM_ARCHITECTURE_OVERVIEW.md) を参照してください。

## 🚀 クイックスタート

### 前提条件

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15
- 立花証券API アクセス権限

### セットアップ

1. **リポジトリクローン**
```bash
git clone <repository-url>
cd stock-bot
```

2. **環境変数設定**
```bash
cp .env.example .env
# .envファイルを編集して認証情報を設定
```

3. **データベース起動**
```bash
docker-compose up -d postgres
```

4. **マイグレーション実行**
```bash
go run migrations/migrate.go
```

5. **アプリケーション起動**
```bash
go run cmd/myapp/main.go
```

### API使用例

```bash
# セッション確認
curl http://localhost:8080/trade/session

# 残高照会
curl http://localhost:8080/trade/balance

# 注文発行
curl -X POST http://localhost:8080/trade/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "1301",
    "trade_type": "BUY",
    "order_type": "LIMIT",
    "quantity": 100,
    "price": 1500.0,
    "position_account_type": "CASH"
  }'
```

## 📊 パフォーマンス

- **スループット**: 758 req/sec (同時接続)
- **注文処理速度**: 1,345 orders/sec
- **平均レスポンス時間**: 1.2ms
- **95パーセンタイル**: 4.9ms
- **並行処理**: 100並行セッション対応

## 🧪 テスト

包括的なテストスイートを提供しています：

```bash
# 全テスト実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# パフォーマンステスト
go test -v ./internal/handler/web/tests/performance_test.go

# E2Eテスト
go test -v ./internal/handler/web/tests/e2e_test.go
```

**テスト完了率**: 95% (114/120項目)

詳細は [テストプラン](TEST_PLAN.md) を参照してください。

## 📁 プロジェクト構造

```
stock-bot/
├── cmd/                    # エントリーポイント
│   ├── myapp/             # メインアプリケーション
│   ├── backtester/        # バックテスト機能
│   └── test-session/      # セッション管理テスト
├── domain/                # ドメイン層
│   ├── model/            # ドメインモデル
│   └── service/          # ドメインサービス
├── internal/              # 内部実装
│   ├── handler/web/      # HTTP APIハンドラー
│   ├── tradeservice/     # アプリケーション層
│   ├── infrastructure/   # インフラストラクチャ層
│   └── eventprocessing/  # イベント処理
├── gen/                   # Goa生成コード
├── migrations/            # データベースマイグレーション
└── data/                  # データファイル
```

## 🔧 設定

### 環境変数

```bash
# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_NAME=stock_trading
DB_USER=postgres
DB_PASSWORD=password

# 立花証券API設定
TACHIBANA_USER_ID=your_user_id
TACHIBANA_PASSWORD=your_password
TACHIBANA_SECOND_PASSWORD=your_second_password

# アプリケーション設定
HTTP_PORT=8080
LOG_LEVEL=info
SESSION_STRATEGY=time_based
```

### Docker Compose

```bash
# 全サービス起動
docker-compose up -d

# ログ確認
docker-compose logs -f app

# サービス停止
docker-compose down
```

## 📚 ドキュメント

- [システムアーキテクチャ概要](SYSTEM_ARCHITECTURE_OVERVIEW.md) - 全体設計と図解
- [現在のアーキテクチャ](CURRENT_ARCHITECTURE.md) - 詳細な技術仕様
- [テストプラン](TEST_PLAN.md) - テスト戦略と進捗
- [セッション管理戦略](SESSION_MANAGEMENT_ARCHITECTURE.md) - セッション管理設計
- [マルチブローカー対応](MULTI_BROKER_ARCHITECTURE.md) - 拡張アーキテクチャ

## 🔒 セキュリティ

- セッションベース認証
- HTTPS通信必須
- 認証情報の環境変数管理
- API通信の暗号化
- ログでの機密情報マスキング

## 📈 監視・運用

### ヘルスチェック
```bash
curl http://localhost:8080/trade/health
```

### ログ監視
```bash
# アプリケーションログ
tail -f logs/app.log

# エラーログ
tail -f logs/error.log
```

### メトリクス
- リクエスト数・レスポンス時間
- エラー率
- データベース接続数
- メモリ・CPU使用率

## 🛠️ 開発

### 開発環境セットアップ

```bash
# 依存関係インストール
go mod download

# 開発用データベース起動
docker-compose -f docker-compose.dev.yml up -d

# ホットリロード（Air使用）
air
```

### コード生成

```bash
# Goaコード生成
goa gen stock-bot/design

# モック生成
go generate ./...
```

### 品質チェック

```bash
# リント
golangci-lint run

# フォーマット
gofmt -s -w .

# 脆弱性チェック
gosec ./...
```

## 🤝 コントリビューション

1. フォークしてブランチを作成
2. 変更を実装
3. テストを追加・実行
4. プルリクエストを作成

## 📄 ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルを参照

## 🆘 サポート

- Issues: GitHub Issues
- ドキュメント: `/docs` ディレクトリ
- API仕様: OpenAPI仕様書

---

**注意**: 本システムは教育・研究目的で作成されています。実際の取引での使用は自己責任でお願いします。