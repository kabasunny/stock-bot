# 株式取引システム 現在のアーキテクチャ

## 概要

本システムは立花証券APIを使用した自動株式取引システムです。Clean Architectureの原則に基づき、レイヤー分離とインターフェース駆動設計を採用しています。

## ディレクトリ構造

```
stock-bot/
├── cmd/                          # エントリーポイント
│   ├── myapp/                    # メインアプリケーション
│   ├── backtester/               # バックテスト機能
│   └── test-session/             # セッション管理テスト
├── domain/                       # ドメイン層
│   ├── model/                    # ドメインモデル
│   └── service/                  # ドメインサービス
├── internal/                     # 内部実装
│   ├── handler/                  # プレゼンテーション層
│   │   └── web/                  # HTTP APIハンドラー
│   ├── tradeservice/             # アプリケーション層
│   ├── infrastructure/           # インフラストラクチャ層
│   │   └── client/               # 外部API クライアント
│   └── eventprocessing/          # イベント処理
├── gen/                          # Goa生成コード
├── migrations/                   # データベースマイグレーション
├── signals/                      # 取引シグナル
└── data/                         # データファイル
```

## レイヤー構成

### 1. ドメイン層 (Domain Layer)

**場所**: `domain/`

**責務**: ビジネスロジックの中核、エンティティとドメインサービス

**主要コンポーネント**:
```go
// domain/model/
type Order struct {
    OrderID             string
    Symbol              string
    TradeType           TradeType
    OrderType           OrderType
    Quantity            int
    Price               float64
    TriggerPrice        float64
    OrderStatus         OrderStatus
    PositionAccountType PositionAccountType
    CreatedAt           time.Time
    UpdatedAt           time.Time
}

type Position struct {
    Symbol              string
    Quantity            int
    AveragePrice        float64
    CurrentPrice        float64
    UnrealizedPL        float64
    PositionType        PositionType
    PositionAccountType PositionAccountType
}

type Session struct {
    SessionID string
    UserID    string
    LoginTime time.Time
    ExpiresAt time.Time
}
```

**インターフェース**:
```go
// domain/service/
type TradeService interface {
    GetSession() *model.Session
    GetPositions(ctx context.Context) ([]*model.Position, error)
    GetOrders(ctx context.Context) ([]*model.Order, error)
    GetBalance(ctx context.Context) (*Balance, error)
    PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*model.Order, error)
    CancelOrder(ctx context.Context, orderID string) error
    CorrectOrder(ctx context.Context, orderID string, newPrice *float64, newQuantity *int) (*model.Order, error)
    GetPriceHistory(ctx context.Context, symbol string, days int) ([]*HistoricalPrice, error)
    HealthCheck(ctx context.Context) (*HealthStatus, error)
}
```

### 2. アプリケーション層 (Application Layer)

**場所**: `internal/tradeservice/`

**責務**: ユースケースの実装、ドメインサービスの調整

**主要コンポーネント**:
- `GoaTradeService`: TradeServiceインターフェースの実装
- `SessionRecoveryService`: セッション回復機能
- `ConversionService`: データ変換機能

**実装例**:
```go
type GoaTradeService struct {
    authClient       client.AuthClient
    balanceClient    client.BalanceClient
    orderClient      client.OrderClient
    priceInfoClient  client.PriceInfoClient
    masterDataClient client.MasterDataClient
    eventClient      client.EventClient
    session          *client.Session
    logger           *slog.Logger
}

func (s *GoaTradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
    // バリデーション
    if err := s.validateOrderRequest(req); err != nil {
        return nil, err
    }
    
    // API呼び出し
    apiReq := s.convertToAPIRequest(req)
    apiResp, err := s.orderClient.NewOrder(ctx, apiReq)
    if err != nil {
        return nil, err
    }
    
    // ドメインモデルに変換
    order := s.convertToOrder(apiResp)
    
    // データベースに保存
    if err := s.saveOrder(ctx, order); err != nil {
        s.logger.Warn("Failed to save order to database", "error", err)
    }
    
    return order, nil
}
```

### 3. インフラストラクチャ層 (Infrastructure Layer)

**場所**: `internal/infrastructure/`

**責務**: 外部システムとの通信、データ永続化

#### 3.1 API クライアント (`internal/infrastructure/client/`)

**TachibanaUnifiedClient**: 統合クライアント
```go
type TachibanaUnifiedClient struct {
    authClient       AuthClient
    balanceClient    BalanceClient
    orderClient      OrderClient
    priceInfoClient  PriceInfoClient
    masterDataClient MasterDataClient
    eventClient      EventClient
    session          *Session
    logger           *slog.Logger
}
```

**個別クライアント**:
- `AuthClient`: 認証・ログイン
- `OrderClient`: 注文管理
- `BalanceClient`: 残高照会
- `MasterDataClient`: マスターデータ
- `PriceInfoClient`: 価格情報
- `EventClient`: WebSocketイベント

#### 3.2 セッション管理

**Session**: セッション情報管理
```go
type Session struct {
    ResultCode     string
    ResultText     string
    SecondPassword string
    RequestURL     string
    MasterURL      string
    PriceURL       string
    EventURL       string
    CookieJar      http.CookieJar
    pNo            atomic.Int32  // アトミックカウンタ
}
```

**SessionManager**: セッション戦略
- `TimeBasedSessionManager`: 時間ベース管理
- `DateBasedSessionManager`: 日付ベース管理

### 4. プレゼンテーション層 (Presentation Layer)

**場所**: `internal/handler/web/`

**責務**: HTTP APIエンドポイント、リクエスト/レスポンス処理

**主要エンドポイント**:
```go
// HTTP API エンドポイント
GET    /trade/session           # セッション情報取得
GET    /trade/positions         # ポジション一覧
GET    /trade/orders            # 注文一覧
GET    /trade/balance           # 残高情報
POST   /trade/orders            # 注文発行
DELETE /trade/orders/{orderID}  # 注文キャンセル
PUT    /trade/orders/{orderID}  # 注文訂正
GET    /trade/price-history/{symbol} # 価格履歴
GET    /trade/health            # ヘルスチェック
```

**実装例**:
```go
func (s *TradeService) PlaceOrder(ctx context.Context, p *trade.PlaceOrderPayload) (*trade.PlaceOrderResult, error) {
    s.logger.Info("TradeService.PlaceOrder called", 
        "symbol", p.Symbol, 
        "trade_type", p.TradeType, 
        "quantity", p.Quantity)

    // ペイロードをサービスリクエストに変換
    req := &service.PlaceOrderRequest{
        Symbol:              p.Symbol,
        TradeType:           convertTradeTypeFromAPI(p.TradeType),
        OrderType:           convertOrderTypeFromAPI(p.OrderType),
        Quantity:            p.Quantity,
        Price:               p.Price,
        TriggerPrice:        p.TriggerPrice,
        PositionAccountType: convertPositionAccountTypeFromAPI(p.PositionAccountType),
    }

    // サービス層を呼び出し
    order, err := s.tradeService.PlaceOrder(ctx, req)
    if err != nil {
        return nil, err
    }

    // レスポンスに変換
    return &trade.PlaceOrderResult{
        OrderID:             order.OrderID,
        Symbol:              order.Symbol,
        TradeType:           convertTradeType(order.TradeType),
        OrderType:           convertOrderType(order.OrderType),
        Quantity:            order.Quantity,
        Price:               order.Price,
        TriggerPrice:        order.TriggerPrice,
        OrderStatus:         convertOrderStatus(order.OrderStatus),
        PositionAccountType: convertPositionAccountType(order.PositionAccountType),
        CreatedAt:           order.CreatedAt.Format(time.RFC3339),
    }, nil
}
```

### 5. イベント処理 (Event Processing)

**場所**: `internal/eventprocessing/`

**責務**: WebSocketイベントの処理、リアルタイム更新

**主要コンポーネント**:
```go
// WebSocketEventService: イベント監視
type WebSocketEventService struct {
    eventClient     client.EventClient
    eventDispatcher service.EventDispatcher
    logger          *slog.Logger
}

// EventDispatcher: イベント振り分け
type EventDispatcher interface {
    RegisterHandler(eventType string, handler EventHandler)
    DispatchEvent(ctx context.Context, eventType string, data map[string]string) error
}

// 各種イベントハンドラー
type ExecutionEventHandler struct{}  // 約定通知
type PriceEventHandler struct{}      // 価格更新
type StatusEventHandler struct{}     // ステータス更新
```

## データベース設計

### テーブル構成

```sql
-- 注文テーブル
CREATE TABLE orders (
    order_id VARCHAR(255) PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    trade_type VARCHAR(10) NOT NULL,
    order_type VARCHAR(20) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2),
    trigger_price DECIMAL(10,2),
    order_status VARCHAR(20) NOT NULL,
    position_account_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ポジションテーブル
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    quantity INTEGER NOT NULL,
    average_price DECIMAL(10,2) NOT NULL,
    current_price DECIMAL(10,2),
    unrealized_pl DECIMAL(12,2),
    position_type VARCHAR(10) NOT NULL,
    position_account_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 約定テーブル
CREATE TABLE executions (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- マスターデータテーブル
CREATE TABLE master_stocks (
    symbol VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    market VARCHAR(50) NOT NULL,
    trading_unit INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 設定管理

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

### Docker Compose設定
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=stock_trading
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: stock_trading
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

volumes:
  postgres_data:
```

## API仕様

### 立花証券API連携

**認証フロー**:
1. ログインAPI呼び出し
2. セッション情報取得
3. 各種APIでセッション使用
4. 8時間後自動再認証

**主要API**:
- `POST /login`: ログイン
- `POST /logout`: ログアウト
- `POST /order`: 注文発行
- `DELETE /order/{orderID}`: 注文キャンセル
- `GET /balance`: 残高照会
- `GET /positions`: ポジション照会
- `WebSocket /events`: リアルタイムイベント

### 内部API仕様

**リクエスト例**:
```json
POST /trade/orders
{
  "symbol": "1301",
  "trade_type": "BUY",
  "order_type": "LIMIT",
  "quantity": 100,
  "price": 1500.0,
  "position_account_type": "CASH"
}
```

**レスポンス例**:
```json
{
  "order_id": "ORD-20240101-001",
  "symbol": "1301",
  "trade_type": "BUY",
  "order_type": "LIMIT",
  "quantity": 100,
  "price": 1500.0,
  "order_status": "NEW",
  "position_account_type": "CASH",
  "created_at": "2024-01-01T09:00:00Z"
}
```

## エラーハンドリング

### エラー分類
1. **ビジネスエラー**: 注文条件不正、残高不足等
2. **システムエラー**: API通信エラー、DB接続エラー等
3. **バリデーションエラー**: 入力値不正等

### エラーレスポンス形式
```json
{
  "name": "validation_error",
  "message": "Invalid order parameters",
  "details": {
    "field": "quantity",
    "reason": "must be positive"
  }
}
```

## ログ設計

### ログレベル
- `DEBUG`: 詳細なデバッグ情報
- `INFO`: 一般的な情報（API呼び出し等）
- `WARN`: 警告（リトライ等）
- `ERROR`: エラー（処理失敗等）

### ログ形式
```json
{
  "time": "2024-01-01T09:00:00Z",
  "level": "INFO",
  "msg": "Order placed successfully",
  "order_id": "ORD-20240101-001",
  "symbol": "1301",
  "quantity": 100
}
```

## セキュリティ

### 認証・認可
- セッションベース認証
- HTTPS通信必須
- 認証情報の環境変数管理

### データ保護
- パスワードの暗号化保存
- API通信の暗号化
- ログでの機密情報マスキング

## パフォーマンス

### 現在の性能指標
- **スループット**: 758 req/sec
- **レスポンス時間**: 平均1.2ms
- **メモリ使用量**: 1.25 bytes/request
- **並行処理**: 100並行セッション対応

### 最適化ポイント
- データベース接続プール
- HTTPクライアント再利用
- メモリプール使用
- 非同期処理活用

## 監視・運用

### ヘルスチェック
```go
func (s *GoaTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
    status := &service.HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Services:  make(map[string]string),
    }
    
    // データベース接続確認
    if err := s.checkDatabase(ctx); err != nil {
        status.Status = "unhealthy"
        status.Services["database"] = "error"
    } else {
        status.Services["database"] = "healthy"
    }
    
    // API接続確認
    if err := s.checkAPIConnection(ctx); err != nil {
        status.Status = "unhealthy"
        status.Services["api"] = "error"
    } else {
        status.Services["api"] = "healthy"
    }
    
    return status, nil
}
```

### メトリクス
- リクエスト数・レスポンス時間
- エラー率
- データベース接続数
- メモリ・CPU使用率

この現在のアーキテクチャにより、拡張可能で保守性の高い株式取引システムが実現されています。