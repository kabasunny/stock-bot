# 個人用軽量版システム提案

## 概要

個人投資家向けに、現在のシステムから過剰な部分を削減し、クラウド運用コストを大幅に削減する軽量版を提案します。

## コスト削減ポイント

### 1. データベース簡略化

#### 現在版
```yaml
services:
  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=stockuser
      - POSTGRES_PASSWORD=stockpassword
    volumes:
      - postgres-data:/var/lib/postgresql/data
```
**月額コスト**: $20-50

#### 軽量版
```go
// SQLite使用
import "database/sql"
import _ "github.com/mattn/go-sqlite3"

db, err := sql.Open("sqlite3", "./data/trading.db")
```
**月額コスト**: $0 (ローカルファイル)

### 2. セッション管理簡略化

#### 現在版
- ファクトリーパターン
- 2種類のマネージャー
- 複雑な設定管理

#### 軽量版
```go
type SimpleSession struct {
    SessionID string    `json:"session_id"`
    UserID    string    `json:"user_id"`
    LoginTime time.Time `json:"login_time"`
    FilePath  string    `json:"-"`
}

func (s *SimpleSession) Save() error {
    data, _ := json.Marshal(s)
    return os.WriteFile(s.FilePath, data, 0600)
}

func LoadSession(filePath string) (*SimpleSession, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    var session SimpleSession
    json.Unmarshal(data, &session)
    session.FilePath = filePath
    return &session, nil
}
```

### 3. ストレージ簡略化

#### 現在版: 複数テーブル + インデックス
```sql
CREATE TABLE orders (...);
CREATE TABLE executions (...);
CREATE TABLE positions (...);
CREATE TABLE stock_masters (...);
-- + 多数のインデックス
```

#### 軽量版: 単一JSONファイル
```go
type TradingData struct {
    Orders     []Order     `json:"orders"`
    Executions []Execution `json:"executions"`
    Positions  []Position  `json:"positions"`
    UpdatedAt  time.Time   `json:"updated_at"`
}

func (td *TradingData) Save(filePath string) error {
    td.UpdatedAt = time.Now()
    data, _ := json.MarshalIndent(td, "", "  ")
    return os.WriteFile(filePath, data, 0644)
}
```

### 4. API簡略化

#### 現在版: Goa Framework
- 複雑なコード生成
- 多数のファイル
- OpenAPI仕様

#### 軽量版: 標準net/http
```go
func main() {
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/orders", handleOrders)
    http.HandleFunc("/positions", handlePositions)
    http.HandleFunc("/balance", handleBalance)
    
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        // 注文一覧取得
    case "POST":
        // 注文発行
    case "DELETE":
        // 注文キャンセル
    }
}
```

### 5. ログ簡略化

#### 現在版: 構造化ログ + 外部保存
```go
slog.Info("order placed", 
    slog.String("order_id", orderID),
    slog.String("symbol", symbol),
    slog.Float64("price", price))
```

#### 軽量版: 標準ログ
```go
log.Printf("Order placed: %s %s @%.2f", orderID, symbol, price)
```

## 軽量版アーキテクチャ

```
┌─────────────────────────────────────────┐
│              Trading Bot                │
│            (Personal Agent)             │
└─────────────────┬───────────────────────┘
                  │ HTTP REST API
┌─────────────────┴───────────────────────┐
│           Lightweight Server            │
│  ┌─────────────┐  ┌─────────────────┐   │
│  │   HTTP      │  │   File Storage  │   │
│  │  Handlers   │  │   (JSON/SQLite) │   │
│  └─────────────┘  └─────────────────┘   │
│  ┌─────────────────────────────────────┐ │
│  │      Tachibana API Client           │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
                  │
┌─────────────────┴───────────────────────┐
│          Tachibana Securities           │
│             REST API                    │
└─────────────────────────────────────────┘
```

## 実装サイズ比較

| 項目 | 現在版 | 軽量版 | 削減率 |
|------|--------|--------|--------|
| **ファイル数** | ~200 | ~20 | 90% |
| **コード行数** | ~15,000 | ~2,000 | 87% |
| **依存関係** | 30+ | 5-10 | 70% |
| **バイナリサイズ** | ~50MB | ~10MB | 80% |
| **メモリ使用量** | ~100MB | ~20MB | 80% |

## 機能比較

| 機能 | 現在版 | 軽量版 | 個人用十分性 |
|------|--------|--------|-------------|
| **ログイン** | ✅ | ✅ | ✅ |
| **注文発行** | ✅ | ✅ | ✅ |
| **ポジション取得** | ✅ | ✅ | ✅ |
| **残高取得** | ✅ | ✅ | ✅ |
| **価格取得** | ✅ | ✅ | ✅ |
| **WebSocket** | ✅ | ❌ | ⚠️ (ポーリングで代替) |
| **複雑なエラーハンドリング** | ✅ | ❌ | ⚠️ (基本的なもので十分) |
| **マスターデータ同期** | ✅ | ❌ | ⚠️ (手動更新で十分) |
| **包括的テスト** | ✅ | ❌ | ⚠️ (基本テストで十分) |

## クラウド運用コスト

### AWS/GCP での月額コスト比較

#### 現在版
```
- EC2 t3.small (2GB RAM): $15-20
- RDS PostgreSQL: $20-30
- CloudWatch Logs: $5-10
- Load Balancer: $15-20
- Storage (20GB): $2-3
合計: $57-83/月
```

#### 軽量版
```
- EC2 t3.nano (512MB RAM): $4-6
- Storage (1GB): $0.1
- 基本監視: $0
合計: $4-6/月
```

**削減額: $50-75/月 (約90%削減)**

## 軽量版実装例

### main.go (軽量版)
```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
)

type TradingServer struct {
    tachibanaClient *TachibanaClient
    storage         *FileStorage
    session         *SimpleSession
}

type FileStorage struct {
    filePath string
    data     *TradingData
}

type TradingData struct {
    Orders     []Order     `json:"orders"`
    Positions  []Position  `json:"positions"`
    UpdatedAt  time.Time   `json:"updated_at"`
}

func main() {
    server := &TradingServer{
        tachibanaClient: NewTachibanaClient(),
        storage:         NewFileStorage("./data/trading.json"),
    }
    
    http.HandleFunc("/login", server.handleLogin)
    http.HandleFunc("/orders", server.handleOrders)
    http.HandleFunc("/positions", server.handlePositions)
    http.HandleFunc("/balance", server.handleBalance)
    
    log.Println("Lightweight trading server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *TradingServer) handleLogin(w http.ResponseWriter, r *http.Request) {
    session, err := s.tachibanaClient.Login()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    s.session = session
    session.Save("./data/session.json")
    
    json.NewEncoder(w).Encode(session)
}

func (s *TradingServer) handleOrders(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        orders := s.storage.GetOrders()
        json.NewEncoder(w).Encode(orders)
    case "POST":
        var order Order
        json.NewDecoder(r.Body).Decode(&order)
        
        result, err := s.tachibanaClient.PlaceOrder(s.session, &order)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        s.storage.AddOrder(&order)
        json.NewEncoder(w).Encode(result)
    }
}
```

### Dockerfile (軽量版)
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o trading-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/trading-server .
EXPOSE 8080
CMD ["./trading-server"]
```

## 移行戦略

### Phase 1: 機能削減
1. PostgreSQL → SQLite移行
2. Goa → 標準HTTP移行
3. 複雑なセッション管理削除

### Phase 2: 軽量化
1. WebSocket → ポーリング変更
2. 構造化ログ削除
3. テストスイート簡略化

### Phase 3: 最適化
1. バイナリサイズ最適化
2. メモリ使用量削減
3. 起動時間短縮

## 推奨事項

### 個人使用なら軽量版推奨
- **コスト**: 90%削減 ($60 → $6/月)
- **保守性**: シンプルで理解しやすい
- **機能**: 個人取引には十分

### 現在版を維持すべき場合
- 複数ユーザー対応予定
- 高頻度取引 (HFT)
- 商用利用予定
- 高可用性要求

## 結論

個人投資家の用途であれば、**軽量版で十分**かつ**大幅なコスト削減**が可能です。
現在のシステムは素晴らしい設計ですが、個人使用には確実に過剰仕様です。

軽量版実装をご希望でしたら、段階的な移行プランを提案できます。