# Stock Trading Bot - 現在のアーキテクチャ図

## 全体アーキテクチャ概要

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Strategy Agent Layer                              │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌───────────┐ │
│  │ LightweightAgent│ │ LightweightAgent│ │ LightweightAgent│ │    ...    │ │
│  │ + SimpleStrategy│ │ + SwingStrategy │ │ + DayStrategy   │ │           │ │
│  │                 │ │                 │ │                 │ │           │ │
│  │ • 戦略実行       │ │ • 戦略実行       │ │ • 戦略実行       │ │           │ │
│  │ • リスク管理     │ │ • リスク管理     │ │ • リスク管理     │ │           │ │
│  │ • HTTP Client   │ │ • HTTP Client   │ │ • HTTP Client   │ │           │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └───────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │ HTTP API Calls
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Goa Service Layer                                │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                        HTTP API Gateway                                 │ │
│  │                                                                         │ │
│  │  GET  /trade/health           - ヘルスチェック                          │ │
│  │  GET  /trade/session          - セッション情報取得                      │ │
│  │  GET  /trade/balance          - 残高情報取得                            │ │
│  │  GET  /trade/positions        - ポジション一覧取得                      │ │
│  │  GET  /trade/orders           - 注文一覧取得                            │ │
│  │  POST /trade/orders           - 注文発行                                │ │
│  │  PUT  /trade/orders/{id}      - 注文訂正                                │ │
│  │  DELETE /trade/orders/{id}    - 注文キャンセル                          │ │
│  │  DELETE /trade/orders         - 全注文キャンセル                        │ │
│  │  GET  /trade/price-history/{symbol} - 価格履歴取得                     │ │
│  │  GET  /trade/symbols/{symbol}/validate - 銘柄妥当性チェック            │ │
│  │  GET  /trade/orders/history   - 注文履歴取得                            │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Domain Service Layer                               │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                         TradeService                                    │ │
│  │  • GetSession()     - セッション管理                                    │ │
│  │  • GetBalance()     - 残高取得                                          │ │
│  │  • GetPositions()   - ポジション管理                                    │ │
│  │  • GetOrders()      - 注文管理                                          │ │
│  │  • PlaceOrder()     - 注文実行                                          │ │
│  │  • CancelOrder()    - 注文キャンセル                                    │ │
│  │  • HealthCheck()    - システム状態監視                                  │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Infrastructure Layer                                   │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                    TachibanaUnifiedClient                               │ │
│  │  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌───────┐ │ │
│  │  │   AuthClient    │ │  BalanceClient  │ │   OrderClient   │ │  ...  │ │ │
│  │  │   (認証I/F)      │ │  (残高I/F)      │ │   (注文I/F)      │ │       │ │ │
│  │  └─────────────────┘ └─────────────────┘ └─────────────────┘ └───────┘ │ │
│  │                                                                         │ │
│  │  • セッション管理（8時間自動更新）                                       │ │
│  │  • 自動再認証機能                                                       │ │
│  │  • 統一されたAPI呼び出しインターフェース                                 │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Tachibana Securities API                            │
│  • 認証API (LOGIN/LOGOUT)                                                   │
│  • 注文API (ORDER/CANCEL/CORRECT)                                           │
│  • 残高API (BALANCE/MARGIN)                                                 │
│  • ポジションAPI (POSITIONS)                                                │
│  • 価格API (PRICE/HISTORY)                                                  │
│  • マスターデータAPI (MASTER)                                               │
│  • WebSocketイベントAPI (EVENTS)                                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 戦略エージェント詳細

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          LightweightAgent                                   │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                        Strategy Interface                               │ │
│  │  • Name() string                                                        │ │
│  │  • Evaluate(MarketData) (*StrategySignal, error)                       │ │
│  │  • GetExecutionInterval() time.Duration                                 │ │
│  │  • GetRiskLimits() *RiskLimits                                          │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐               │
│  │ SimpleStrategy  │ │ SwingStrategy   │ │ DayTradingStrategy│               │
│  │                 │ │                 │ │                 │               │
│  │ • 30秒間隔       │ │ • 1時間間隔      │ │ • 5分間隔        │               │
│  │ • 基本ロジック   │ │ • 9:00-11:30   │ │ • 9:00-14:30   │               │
│  │ • リスク管理     │ │ • スイング戦略   │ │ • デイトレード   │               │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘               │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                        Execution Loop                                   │ │
│  │  1. ヘルスチェック (GET /trade/health)                                   │ │
│  │  2. 市場データ収集:                                                      │ │
│  │     • 残高取得 (GET /trade/balance)                                      │ │
│  │     • ポジション取得 (GET /trade/positions)                              │ │
│  │     • 注文一覧取得 (GET /trade/orders)                                   │ │
│  │  3. 戦略評価 (Strategy.Evaluate())                                      │ │
│  │  4. 注文実行 (POST /trade/orders)                                       │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

## データフロー図

```
┌─────────────────┐    HTTP Request     ┌─────────────────┐
│ LightweightAgent│ ──────────────────► │   Goa Service   │
│                 │                     │                 │
│ • Strategy Logic│                     │ • HTTP Handlers │
│ • Risk Mgmt     │                     │ • Request/Resp  │
│ • HTTP Client   │                     │ • Validation    │
└─────────────────┘                     └─────────────────┘
         ▲                                        │
         │                                        ▼
         │ JSON Response              ┌─────────────────┐
         └────────────────────────────│  TradeService   │
                                      │                 │
                                      │ • Business Logic│
                                      │ • Domain Rules  │
                                      │ • Data Transform│
                                      └─────────────────┘
                                                │
                                                ▼
                                      ┌─────────────────┐
                                      │TachibanaUnified │
                                      │     Client      │
                                      │                 │
                                      │ • Session Mgmt  │
                                      │ • Auto Re-auth  │
                                      │ • API Calls     │
                                      └─────────────────┘
                                                │
                                                ▼
                                      ┌─────────────────┐
                                      │ Tachibana API   │
                                      │                 │
                                      │ • REST API      │
                                      │ • WebSocket     │
                                      │ • Authentication│
                                      └─────────────────┘
```

## 実行時の並列処理

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Runtime Architecture                              │
└─────────────────────────────────────────────────────────────────────────────┘

Process 1: goa-service.exe --no-tachibana --no-db
┌─────────────────┐
│   Goa Service   │ ← HTTP Server (Port 8080)
│                 │
│ • HTTP Handlers │
│ • TradeService  │
│ • Mock Data     │
└─────────────────┘

Process 2: lightweight-agent.exe --strategy=simple --symbol=7203
┌─────────────────┐
│ Simple Agent    │ ← HTTP Client → Goa Service
│                 │
│ • 30sec interval│
│ • Basic Logic   │
│ • Risk Limits   │
└─────────────────┘

Process 3: lightweight-agent.exe --strategy=swing --symbol=6758
┌─────────────────┐
│ Swing Agent     │ ← HTTP Client → Goa Service
│                 │
│ • 1hour interval│
│ • Time Limits   │
│ • Swing Logic   │
└─────────────────┘

Process 4: lightweight-agent.exe --strategy=day --symbol=9984
┌─────────────────┐
│ Day Agent       │ ← HTTP Client → Goa Service
│                 │
│ • 5min interval │
│ • Day Trading   │
│ • Time Limits   │
└─────────────────┘
```

## 主要な設計原則

### 1. **責務分離 (Separation of Concerns)**
- **Goa Service**: 取引API抽象化層
- **LightweightAgent**: 戦略実行エンジン
- **Strategy**: 戦略ロジック

### 2. **疎結合 (Loose Coupling)**
- HTTP APIによる通信
- インターフェースベースの設計
- 独立したプロセス実行

### 3. **拡張性 (Extensibility)**
- 新戦略の簡単な追加
- 複数エージェントの並列実行
- マイクロサービス対応

### 4. **テスト容易性 (Testability)**
- HTTP APIのモック化
- 各層の独立テスト
- 統合テストの実行

### 5. **運用性 (Operability)**
- ヘルスチェック機能
- 構造化ログ出力
- Graceful Shutdown

この構成により、戦略管理をエージェント側で行い、Goaサービスを純粋な取引APIラッパーとして機能させる、クリーンで拡張性の高いアーキテクチャを実現しています。