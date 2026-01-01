# 複数証券会社対応アーキテクチャ

## 拡張されたアーキテクチャ概要

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Strategy Agent Layer                              │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌───────────┐ │
│  │ LightweightAgent│ │ LightweightAgent│ │ LightweightAgent│ │    ...    │ │
│  │ + Strategy      │ │ + Strategy      │ │ + Strategy      │ │           │ │
│  │ + BrokerConfig  │ │ + BrokerConfig  │ │ + BrokerConfig  │ │           │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └───────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │ HTTP API Calls
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Goa Service Layer                                │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                     Multi-Broker API Gateway                            │ │
│  │                                                                         │ │
│  │  POST /brokers/{broker}/orders    - 指定証券会社への注文発行             │ │
│  │  GET  /brokers/{broker}/balance   - 指定証券会社の残高取得               │ │
│  │  GET  /brokers/{broker}/positions - 指定証券会社のポジション取得         │ │
│  │  GET  /brokers                    - 利用可能証券会社一覧                 │ │
│  │                                                                         │ │
│  │  # 既存API（デフォルト証券会社）                                         │ │
│  │  POST /trade/orders               - デフォルト証券会社への注文           │ │
│  │  GET  /trade/balance              - デフォルト証券会社の残高             │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Broker Service Factory                                │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                      TradeServiceFactory                                │ │
│  │                                                                         │ │
│  │  func GetTradeService(brokerType string) TradeService {                 │ │
│  │    switch brokerType {                                                  │ │
│  │      case "tachibana":                                                  │ │
│  │        return NewTachibanaTradeService(...)                             │ │
│  │      case "sbi":                                                        │ │
│  │        return NewSBITradeService(...)                                   │ │
│  │      case "rakuten":                                                    │ │
│  │        return NewRakutenTradeService(...)                               │ │
│  │      case "matsui":                                                     │ │
│  │        return NewMatsuiTradeService(...)                                │ │
│  │      default:                                                           │ │
│  │        return NewTachibanaTradeService(...) // デフォルト               │ │
│  │    }                                                                    │ │
│  │  }                                                                      │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        TradeService Interface                               │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │  type TradeService interface {                                          │ │
│  │    GetSession() *model.Session                                          │ │
│  │    GetBalance(ctx) (*service.Balance, error)                           │ │
│  │    GetPositions(ctx) ([]*model.Position, error)                        │ │
│  │    PlaceOrder(ctx, *service.PlaceOrderRequest) (*model.Order, error)   │ │
│  │    CancelOrder(ctx, orderID string) error                              │ │
│  │    HealthCheck(ctx) (*service.HealthStatus, error)                     │ │
│  │    // ... その他のメソッド                                               │ │
│  │  }                                                                      │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Multiple Broker Implementations                         │
│                                                                             │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌───────────┐ │
│  │TachibanaTradeService│ │  SBITradeService │ │RakutenTradeService│ │   ...   │ │
│  │                 │ │                 │ │                 │ │           │ │
│  │ • TachibanaAPI  │ │ • SBI API       │ │ • Rakuten API   │ │           │ │
│  │ • 立花固有ロジック│ │ • SBI固有ロジック│ │ • 楽天固有ロジック│ │           │ │
│  │ • 認証方式      │ │ • 認証方式      │ │ • 認証方式      │ │           │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └───────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Broker-Specific APIs                                │
│                                                                             │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌───────────┐ │
│  │  Tachibana API  │ │    SBI API      │ │   Rakuten API   │ │    ...    │ │
│  │                 │ │                 │ │                 │ │           │ │
│  │ • REST API      │ │ • REST API      │ │ • REST API      │ │           │ │
│  │ • WebSocket     │ │ • WebSocket     │ │ • WebSocket     │ │           │ │
│  │ • 独自認証      │ │ • 独自認証      │ │ • 独自認証      │ │           │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └───────────┘ │
└─────────────────────────────────────────────────────────────────────────────┐
```

## 実装例

### 1. TradeServiceFactory の実装

```go
// internal/tradeservice/factory.go
package tradeservice

import (
    "fmt"
    "log/slog"
    "stock-bot/domain/service"
    "stock-bot/internal/config"
)

type BrokerType string

const (
    BrokerTachibana BrokerType = "tachibana"
    BrokerSBI       BrokerType = "sbi"
    BrokerRakuten   BrokerType = "rakuten"
    BrokerMatsui    BrokerType = "matsui"
)

type TradeServiceFactory struct {
    config *config.Config
    logger *slog.Logger
}

func NewTradeServiceFactory(config *config.Config, logger *slog.Logger) *TradeServiceFactory {
    return &TradeServiceFactory{
        config: config,
        logger: logger,
    }
}

func (f *TradeServiceFactory) CreateTradeService(brokerType BrokerType) (service.TradeService, error) {
    switch brokerType {
    case BrokerTachibana:
        return f.createTachibanaService()
    case BrokerSBI:
        return f.createSBIService()
    case BrokerRakuten:
        return f.createRakutenService()
    case BrokerMatsui:
        return f.createMatsuiService()
    default:
        return nil, fmt.Errorf("unsupported broker type: %s", brokerType)
    }
}

func (f *TradeServiceFactory) createTachibanaService() (service.TradeService, error) {
    // 既存のTachibanaTradeService実装
    return NewGoaTradeService(/* 立花用の設定 */), nil
}

func (f *TradeServiceFactory) createSBIService() (service.TradeService, error) {
    // SBI証券用の実装
    return NewSBITradeService(/* SBI用の設定 */), nil
}

func (f *TradeServiceFactory) createRakutenService() (service.TradeService, error) {
    // 楽天証券用の実装
    return NewRakutenTradeService(/* 楽天用の設定 */), nil
}

func (f *TradeServiceFactory) createMatsuiService() (service.TradeService, error) {
    // 松井証券用の実装
    return NewMatsuiTradeService(/* 松井用の設定 */), nil
}
```

### 2. SBI証券用TradeService実装例

```go
// internal/tradeservice/sbi_trade_service.go
package tradeservice

import (
    "context"
    "stock-bot/domain/model"
    "stock-bot/domain/service"
    "stock-bot/internal/infrastructure/client/sbi"
)

type SBITradeService struct {
    sbiClient *sbi.SBIClient
    logger    *slog.Logger
}

func NewSBITradeService(sbiClient *sbi.SBIClient, logger *slog.Logger) *SBITradeService {
    return &SBITradeService{
        sbiClient: sbiClient,
        logger:    logger,
    }
}

func (s *SBITradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
    // SBI証券APIを使用した残高取得
    balance, err := s.sbiClient.GetBalance(ctx)
    if err != nil {
        return nil, err
    }
    
    return &service.Balance{
        Cash:        balance.Cash,
        BuyingPower: balance.BuyingPower,
    }, nil
}

func (s *SBITradeService) PlaceOrder(ctx context.Context, req *service.PlaceOrderRequest) (*model.Order, error) {
    // SBI証券APIを使用した注文発行
    sbiOrder := &sbi.OrderRequest{
        Symbol:    req.Symbol,
        Side:      convertTradeType(req.TradeType),
        OrderType: convertOrderType(req.OrderType),
        Quantity:  req.Quantity,
        Price:     req.Price,
    }
    
    response, err := s.sbiClient.PlaceOrder(ctx, sbiOrder)
    if err != nil {
        return nil, err
    }
    
    return &model.Order{
        OrderID:     response.OrderID,
        Symbol:      req.Symbol,
        TradeType:   req.TradeType,
        OrderType:   req.OrderType,
        Quantity:    req.Quantity,
        Price:       req.Price,
        OrderStatus: model.OrderStatusNew,
    }, nil
}

// その他のメソッド実装...
```

### 3. 拡張されたGoa API設計

```go
// design/design.go に追加
var _ = Service("broker", func() {
    Description("Multi-broker trading service")

    // GET /brokers
    Method("list", func() {
        Description("List available brokers")
        Payload(Empty)
        Result(func() {
            Attribute("brokers", ArrayOf(BrokerInfo), "利用可能な証券会社一覧")
            Required("brokers")
        })
        HTTP(func() {
            GET("/brokers")
            Response(StatusOK)
        })
    })

    // POST /brokers/{broker}/orders
    Method("place_order", func() {
        Description("Place order with specific broker")
        Payload(func() {
            Attribute("broker", String, "証券会社ID", func() {
                Enum("tachibana", "sbi", "rakuten", "matsui")
            })
            Attribute("symbol", String, "銘柄コード")
            Attribute("trade_type", String, "売買区分", func() {
                Enum("BUY", "SELL")
            })
            Attribute("order_type", String, "注文種別", func() {
                Enum("MARKET", "LIMIT", "STOP")
            })
            Attribute("quantity", UInt, "数量")
            Attribute("price", Float64, "価格")
            Required("broker", "symbol", "trade_type", "order_type", "quantity")
        })
        Result(TradeOrderResult)
        HTTP(func() {
            POST("/brokers/{broker}/orders")
            Response(StatusCreated)
        })
    })

    // GET /brokers/{broker}/balance
    Method("get_balance", func() {
        Description("Get balance from specific broker")
        Payload(func() {
            Attribute("broker", String, "証券会社ID")
            Required("broker")
        })
        Result(TradeBalanceResult)
        HTTP(func() {
            GET("/brokers/{broker}/balance")
            Response(StatusOK)
        })
    })
})

var BrokerInfo = Type("BrokerInfo", func() {
    Description("Broker information")
    Attribute("id", String, "証券会社ID")
    Attribute("name", String, "証券会社名")
    Attribute("status", String, "接続状態", func() {
        Enum("connected", "disconnected", "error")
    })
    Attribute("features", ArrayOf(String), "対応機能")
    Required("id", "name", "status")
})
```

### 4. エージェント側での証券会社指定

```go
// cmd/lightweight-agent/main.go に追加
func main() {
    baseURL := flag.String("base-url", "http://localhost:8080", "Base URL of the Goa service")
    broker := flag.String("broker", "tachibana", "Broker to use (tachibana, sbi, rakuten, matsui)")
    strategyType := flag.String("strategy", "simple", "Strategy type")
    // ... その他のフラグ

    // 証券会社指定でエージェント作成
    lightAgent := agent.NewLightweightAgentWithBroker(*baseURL, *broker, strategy, logger)
}
```

```go
// internal/agent/lightweight_agent.go に追加
type LightweightAgent struct {
    httpClient *http.Client
    baseURL    string
    broker     string  // 追加: 証券会社指定
    strategy   Strategy
    logger     *slog.Logger
}

func (a *LightweightAgent) placeOrder(ctx context.Context, req *PlaceOrderRequest) (*Order, error) {
    // 証券会社指定のエンドポイントを使用
    endpoint := fmt.Sprintf("%s/brokers/%s/orders", a.baseURL, a.broker)
    
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    resp, err := a.httpClient.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
    // ... 残りの実装
}
```

## 利用例

### 複数証券会社での並列実行
```bash
# 立花証券でスイング戦略
./lightweight-agent.exe --broker=tachibana --strategy=swing --symbol=7203

# SBI証券でデイトレード戦略  
./lightweight-agent.exe --broker=sbi --strategy=day --symbol=6758

# 楽天証券で基本戦略
./lightweight-agent.exe --broker=rakuten --strategy=simple --symbol=9984
```

### 証券会社別の残高確認
```bash
curl http://localhost:8080/brokers/tachibana/balance
curl http://localhost:8080/brokers/sbi/balance
curl http://localhost:8080/brokers/rakuten/balance
```

## 実装の利点

### 1. **拡張性**
- 新しい証券会社の追加が容易
- 既存コードへの影響を最小限に抑制

### 2. **柔軟性**
- 証券会社ごとに異なる戦略を実行可能
- リスク分散（複数口座での取引）

### 3. **保守性**
- 各証券会社の実装が独立
- インターフェースによる統一API

### 4. **テスト容易性**
- 証券会社ごとのモック実装
- 統合テストの実行

この設計により、将来的に複数の証券会社に対応でき、各証券会社の特性に応じた最適化も可能になります。