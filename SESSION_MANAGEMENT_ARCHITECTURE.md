# マルチユーザー・セッション管理アーキテクチャ

## 現在の課題と解決策

### **現在の問題**
- 単一セッション（グローバル変数的な管理）
- 複数ユーザー・エージェントの同時利用不可
- セッション情報の永続化なし

### **解決すべき要件**
1. **マルチユーザー対応**: 複数のユーザーが同時にログイン
2. **マルチエージェント対応**: 1ユーザーが複数エージェントを実行
3. **セッション永続化**: サーバー再起動後もセッション維持
4. **セッション有効期限管理**: 自動ログアウト・再認証
5. **セキュリティ**: セッションハイジャック対策

## 拡張されたアーキテクチャ

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Client Layer                                      │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌───────────┐ │
│  │ Agent A         │ │ Agent B         │ │ Web Client      │ │    ...    │ │
│  │ User: alice     │ │ User: bob       │ │ User: charlie   │ │           │ │
│  │ Token: abc123   │ │ Token: def456   │ │ Token: ghi789   │ │           │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └───────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │ HTTP + Authorization Header
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Goa Service Layer                                    │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                    Authentication Middleware                            │ │
│  │  • JWT Token Validation                                                 │ │
│  │  • Session Lookup                                                       │ │
│  │  • User Context Injection                                               │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                        API Endpoints                                    │ │
│  │                                                                         │ │
│  │  POST /auth/login                 - ログイン（トークン発行）             │ │
│  │  POST /auth/logout                - ログアウト                          │ │
│  │  POST /auth/refresh               - トークン更新                        │ │
│  │  GET  /auth/sessions              - セッション一覧                      │ │
│  │                                                                         │ │
│  │  GET  /trade/balance              - 残高取得（認証必須）                │ │
│  │  POST /trade/orders               - 注文発行（認証必須）                │ │
│  │  GET  /users/me                   - 現在のユーザー情報                  │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Session Management Layer                              │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                        SessionManager                                   │ │
│  │                                                                         │ │
│  │  • CreateSession(userID, brokerType) (*Session, error)                 │ │
│  │  • GetSession(sessionID) (*Session, error)                             │ │
│  │  • RefreshSession(sessionID) error                                      │ │
│  │  • DeleteSession(sessionID) error                                       │ │
│  │  • CleanupExpiredSessions()                                             │ │
│  │  • GetUserSessions(userID) ([]*Session, error)                         │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Session Storage Layer                                │
│                                                                             │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌───────────┐ │
│  │   Redis Cache   │ │   Database      │ │   Memory Store  │ │    ...    │ │
│  │                 │ │                 │ │                 │ │           │ │
│  │ • 高速アクセス   │ │ • 永続化        │ │ • 開発・テスト用 │ │           │ │
│  │ • TTL管理       │ │ • 履歴管理      │ │ • 軽量          │ │           │ │
│  │ • 分散対応      │ │ • 監査ログ      │ │ • 設定不要      │ │           │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └───────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Broker Client Management                               │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │                    BrokerClientManager                                  │ │
│  │                                                                         │ │
│  │  Map[SessionID]*BrokerClient {                                          │ │
│  │    "session_alice_tachibana": TachibanaClient{...},                     │ │
│  │    "session_bob_sbi":         SBIClient{...},                           │ │
│  │    "session_charlie_rakuten": RakutenClient{...},                       │ │
│  │  }                                                                      │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

## データモデル設計

### 1. Session Entity

```go
// domain/model/session.go
package model

import (
    "time"
)

type Session struct {
    SessionID    string    `json:"session_id" db:"session_id"`
    UserID       string    `json:"user_id" db:"user_id"`
    BrokerType   string    `json:"broker_type" db:"broker_type"`
    BrokerUserID string    `json:"broker_user_id" db:"broker_user_id"`
    
    // JWT関連
    AccessToken  string    `json:"access_token" db:"access_token"`
    RefreshToken string    `json:"refresh_token" db:"refresh_token"`
    
    // セッション管理
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
    LastUsedAt   time.Time `json:"last_used_at" db:"last_used_at"`
    
    // 証券会社セッション情報
    BrokerSession *BrokerSession `json:"broker_session"`
    
    // メタデータ
    ClientIP     string    `json:"client_ip" db:"client_ip"`
    UserAgent    string    `json:"user_agent" db:"user_agent"`
    IsActive     bool      `json:"is_active" db:"is_active"`
}

type BrokerSession struct {
    BrokerSessionID string            `json:"broker_session_id"`
    BrokerToken     string            `json:"broker_token"`
    BrokerExpiresAt time.Time         `json:"broker_expires_at"`
    BrokerData      map[string]string `json:"broker_data"`
}

type User struct {
    UserID       string    `json:"user_id" db:"user_id"`
    Username     string    `json:"username" db:"username"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    IsActive     bool      `json:"is_active" db:"is_active"`
}
```

### 2. SessionManager Implementation

```go
// internal/session/session_manager.go
package session

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "time"
    "stock-bot/domain/model"
)

type SessionManager interface {
    CreateSession(ctx context.Context, req *CreateSessionRequest) (*model.Session, error)
    GetSession(ctx context.Context, sessionID string) (*model.Session, error)
    RefreshSession(ctx context.Context, sessionID string) (*model.Session, error)
    DeleteSession(ctx context.Context, sessionID string) error
    GetUserSessions(ctx context.Context, userID string) ([]*model.Session, error)
    CleanupExpiredSessions(ctx context.Context) error
}

type CreateSessionRequest struct {
    UserID       string
    BrokerType   string
    BrokerUserID string
    BrokerPassword string
    ClientIP     string
    UserAgent    string
}

type SessionManagerImpl struct {
    sessionStore   SessionStore
    brokerFactory  BrokerClientFactory
    jwtService     JWTService
    logger         *slog.Logger
}

func (sm *SessionManagerImpl) CreateSession(ctx context.Context, req *CreateSessionRequest) (*model.Session, error) {
    // 1. セッションID生成
    sessionID, err := generateSessionID()
    if err != nil {
        return nil, err
    }
    
    // 2. 証券会社クライアント作成・認証
    brokerClient, err := sm.brokerFactory.CreateClient(req.BrokerType)
    if err != nil {
        return nil, err
    }
    
    brokerSession, err := brokerClient.Login(ctx, req.BrokerUserID, req.BrokerPassword)
    if err != nil {
        return nil, fmt.Errorf("broker login failed: %w", err)
    }
    
    // 3. JWT トークン生成
    accessToken, refreshToken, err := sm.jwtService.GenerateTokens(req.UserID, sessionID)
    if err != nil {
        return nil, err
    }
    
    // 4. セッション作成
    session := &model.Session{
        SessionID:     sessionID,
        UserID:        req.UserID,
        BrokerType:    req.BrokerType,
        BrokerUserID:  req.BrokerUserID,
        AccessToken:   accessToken,
        RefreshToken:  refreshToken,
        CreatedAt:     time.Now(),
        ExpiresAt:     time.Now().Add(8 * time.Hour), // 8時間有効
        LastUsedAt:    time.Now(),
        BrokerSession: brokerSession,
        ClientIP:      req.ClientIP,
        UserAgent:     req.UserAgent,
        IsActive:      true,
    }
    
    // 5. セッション保存
    if err := sm.sessionStore.Save(ctx, session); err != nil {
        return nil, err
    }
    
    sm.logger.Info("Session created", 
        "session_id", sessionID,
        "user_id", req.UserID,
        "broker_type", req.BrokerType)
    
    return session, nil
}

func generateSessionID() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}
```

### 3. Session Store Interface & Implementations

```go
// internal/session/session_store.go
package session

import (
    "context"
    "stock-bot/domain/model"
)

type SessionStore interface {
    Save(ctx context.Context, session *model.Session) error
    Get(ctx context.Context, sessionID string) (*model.Session, error)
    Update(ctx context.Context, session *model.Session) error
    Delete(ctx context.Context, sessionID string) error
    GetByUserID(ctx context.Context, userID string) ([]*model.Session, error)
    DeleteExpired(ctx context.Context) error
}

// Redis Implementation
type RedisSessionStore struct {
    client redis.Client
    ttl    time.Duration
}

func (r *RedisSessionStore) Save(ctx context.Context, session *model.Session) error {
    key := fmt.Sprintf("session:%s", session.SessionID)
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }
    
    return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *RedisSessionStore) Get(ctx context.Context, sessionID string) (*model.Session, error) {
    key := fmt.Sprintf("session:%s", sessionID)
    data, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }
    
    var session model.Session
    if err := json.Unmarshal([]byte(data), &session); err != nil {
        return nil, err
    }
    
    return &session, nil
}

// Database Implementation
type DatabaseSessionStore struct {
    db *sql.DB
}

func (d *DatabaseSessionStore) Save(ctx context.Context, session *model.Session) error {
    query := `
        INSERT INTO sessions (
            session_id, user_id, broker_type, broker_user_id,
            access_token, refresh_token, created_at, expires_at,
            last_used_at, broker_session, client_ip, user_agent, is_active
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        ON CONFLICT (session_id) DO UPDATE SET
            last_used_at = $9, broker_session = $10, is_active = $13
    `
    
    brokerSessionJSON, _ := json.Marshal(session.BrokerSession)
    
    _, err := d.db.ExecContext(ctx, query,
        session.SessionID, session.UserID, session.BrokerType, session.BrokerUserID,
        session.AccessToken, session.RefreshToken, session.CreatedAt, session.ExpiresAt,
        session.LastUsedAt, brokerSessionJSON, session.ClientIP, session.UserAgent, session.IsActive)
    
    return err
}

// Memory Implementation (for development)
type MemorySessionStore struct {
    sessions map[string]*model.Session
    mutex    sync.RWMutex
}

func (m *MemorySessionStore) Save(ctx context.Context, session *model.Session) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.sessions[session.SessionID] = session
    return nil
}
```

### 4. Authentication Middleware

```go
// internal/middleware/auth_middleware.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "stock-bot/internal/session"
)

type AuthMiddleware struct {
    sessionManager session.SessionManager
    jwtService     JWTService
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Authorization ヘッダーからトークン取得
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            http.Error(w, "Bearer token required", http.StatusUnauthorized)
            return
        }
        
        // 2. JWT トークン検証
        claims, err := am.jwtService.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // 3. セッション取得
        session, err := am.sessionManager.GetSession(r.Context(), claims.SessionID)
        if err != nil {
            http.Error(w, "Session not found", http.StatusUnauthorized)
            return
        }
        
        // 4. セッション有効性チェック
        if !session.IsActive || session.ExpiresAt.Before(time.Now()) {
            http.Error(w, "Session expired", http.StatusUnauthorized)
            return
        }
        
        // 5. コンテキストにセッション情報を追加
        ctx := context.WithValue(r.Context(), "session", session)
        ctx = context.WithValue(ctx, "user_id", session.UserID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 5. 拡張されたGoa API設計

```go
// design/design.go に追加
var _ = Service("auth", func() {
    Description("Authentication and session management service")

    // POST /auth/login
    Method("login", func() {
        Description("Login and create session")
        Payload(func() {
            Attribute("username", String, "ユーザー名")
            Attribute("password", String, "パスワード")
            Attribute("broker_type", String, "証券会社", func() {
                Enum("tachibana", "sbi", "rakuten", "matsui")
            })
            Attribute("broker_user_id", String, "証券会社ユーザーID")
            Attribute("broker_password", String, "証券会社パスワード")
            Required("username", "password", "broker_type", "broker_user_id", "broker_password")
        })
        Result(func() {
            Attribute("access_token", String, "アクセストークン")
            Attribute("refresh_token", String, "リフレッシュトークン")
            Attribute("expires_in", UInt, "有効期限（秒）")
            Attribute("session_id", String, "セッションID")
            Required("access_token", "refresh_token", "expires_in", "session_id")
        })
        HTTP(func() {
            POST("/auth/login")
            Response(StatusOK)
        })
    })

    // POST /auth/logout
    Method("logout", func() {
        Description("Logout and destroy session")
        Payload(Empty)
        Result(Empty)
        HTTP(func() {
            POST("/auth/logout")
            Response(StatusNoContent)
        })
    })

    // POST /auth/refresh
    Method("refresh", func() {
        Description("Refresh access token")
        Payload(func() {
            Attribute("refresh_token", String, "リフレッシュトークン")
            Required("refresh_token")
        })
        Result(func() {
            Attribute("access_token", String, "新しいアクセストークン")
            Attribute("expires_in", UInt, "有効期限（秒）")
            Required("access_token", "expires_in")
        })
        HTTP(func() {
            POST("/auth/refresh")
            Response(StatusOK)
        })
    })

    // GET /auth/sessions
    Method("list_sessions", func() {
        Description("List user sessions")
        Payload(Empty)
        Result(func() {
            Attribute("sessions", ArrayOf(SessionInfo), "セッション一覧")
            Required("sessions")
        })
        HTTP(func() {
            GET("/auth/sessions")
            Response(StatusOK)
        })
    })
})

var SessionInfo = Type("SessionInfo", func() {
    Description("Session information")
    Attribute("session_id", String, "セッションID")
    Attribute("broker_type", String, "証券会社")
    Attribute("created_at", String, "作成日時")
    Attribute("last_used_at", String, "最終使用日時")
    Attribute("expires_at", String, "有効期限")
    Attribute("client_ip", String, "クライアントIP")
    Attribute("is_active", Boolean, "アクティブ状態")
    Required("session_id", "broker_type", "created_at", "is_active")
})
```

## エージェント側の実装

### 1. 認証機能付きLightweightAgent

```go
// internal/agent/authenticated_agent.go
package agent

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type AuthenticatedAgent struct {
    *LightweightAgent
    credentials *Credentials
    authToken   string
    refreshToken string
    tokenExpiry time.Time
}

type Credentials struct {
    Username       string `json:"username"`
    Password       string `json:"password"`
    BrokerType     string `json:"broker_type"`
    BrokerUserID   string `json:"broker_user_id"`
    BrokerPassword string `json:"broker_password"`
}

func NewAuthenticatedAgent(baseURL string, credentials *Credentials, strategy Strategy, logger *slog.Logger) *AuthenticatedAgent {
    return &AuthenticatedAgent{
        LightweightAgent: NewLightweightAgent(baseURL, strategy, logger),
        credentials:      credentials,
    }
}

func (a *AuthenticatedAgent) Run(ctx context.Context) error {
    // 1. 初回ログイン
    if err := a.login(ctx); err != nil {
        return fmt.Errorf("initial login failed: %w", err)
    }

    // 2. トークン更新用のgoroutine起動
    go a.tokenRefreshLoop(ctx)

    // 3. 通常のエージェント処理開始
    return a.LightweightAgent.Run(ctx)
}

func (a *AuthenticatedAgent) login(ctx context.Context) error {
    loginData := map[string]string{
        "username":        a.credentials.Username,
        "password":        a.credentials.Password,
        "broker_type":     a.credentials.BrokerType,
        "broker_user_id":  a.credentials.BrokerUserID,
        "broker_password": a.credentials.BrokerPassword,
    }

    jsonData, err := json.Marshal(loginData)
    if err != nil {
        return err
    }

    resp, err := a.httpClient.Post(
        a.baseURL+"/auth/login",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("login failed with status: %d", resp.StatusCode)
    }

    var loginResp struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    uint   `json:"expires_in"`
        SessionID    string `json:"session_id"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
        return err
    }

    a.authToken = loginResp.AccessToken
    a.refreshToken = loginResp.RefreshToken
    a.tokenExpiry = time.Now().Add(time.Duration(loginResp.ExpiresIn) * time.Second)

    a.logger.Info("Login successful", 
        "session_id", loginResp.SessionID,
        "expires_at", a.tokenExpiry)

    return nil
}

func (a *AuthenticatedAgent) tokenRefreshLoop(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Minute) // 30分ごとにチェック
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // トークンの有効期限が1時間以内の場合、更新
            if time.Until(a.tokenExpiry) < time.Hour {
                if err := a.refreshAccessToken(ctx); err != nil {
                    a.logger.Error("Token refresh failed", "error", err)
                }
            }
        }
    }
}

func (a *AuthenticatedAgent) makeAuthenticatedRequest(req *http.Request) (*http.Response, error) {
    // Authorization ヘッダーを追加
    req.Header.Set("Authorization", "Bearer "+a.authToken)
    return a.httpClient.Do(req)
}

// 既存のHTTPリクエストメソッドをオーバーライド
func (a *AuthenticatedAgent) getBalance(ctx context.Context) (*Balance, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", a.baseURL+"/trade/balance", nil)
    if err != nil {
        return nil, err
    }

    resp, err := a.makeAuthenticatedRequest(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusUnauthorized {
        // トークンが無効な場合、再ログイン
        if err := a.login(ctx); err != nil {
            return nil, err
        }
        // リトライ
        return a.getBalance(ctx)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
    }

    var balance Balance
    if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
        return nil, err
    }

    return &balance, nil
}
```

### 2. エージェント起動例

```go
// cmd/authenticated-agent/main.go
package main

import (
    "context"
    "flag"
    "log/slog"
    "os"
    "stock-bot/internal/agent"
)

func main() {
    // コマンドラインフラグ
    baseURL := flag.String("base-url", "http://localhost:8080", "Base URL")
    username := flag.String("username", "", "Username")
    password := flag.String("password", "", "Password")
    brokerType := flag.String("broker", "tachibana", "Broker type")
    brokerUserID := flag.String("broker-user", "", "Broker user ID")
    brokerPassword := flag.String("broker-password", "", "Broker password")
    strategyType := flag.String("strategy", "simple", "Strategy type")
    targetSymbol := flag.String("symbol", "7203", "Target symbol")
    quantity := flag.Uint("quantity", 100, "Order quantity")
    flag.Parse()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    // 認証情報
    credentials := &agent.Credentials{
        Username:       *username,
        Password:       *password,
        BrokerType:     *brokerType,
        BrokerUserID:   *brokerUserID,
        BrokerPassword: *brokerPassword,
    }

    // 戦略作成
    var strategy agent.Strategy
    switch *strategyType {
    case "swing":
        strategy = agent.NewSwingStrategy("SwingStrategy", *targetSymbol, *quantity, logger)
    case "day":
        strategy = agent.NewDayTradingStrategy("DayTradingStrategy", *targetSymbol, *quantity, logger)
    default:
        strategy = agent.NewSimpleStrategy("SimpleStrategy", *targetSymbol, *quantity, logger)
    }

    // 認証付きエージェント作成
    authAgent := agent.NewAuthenticatedAgent(*baseURL, credentials, strategy, logger)

    // 実行
    ctx := context.Background()
    if err := authAgent.Run(ctx); err != nil {
        logger.Error("Agent execution failed", "error", err)
        os.Exit(1)
    }
}
```

## 利用例

### 1. 複数ユーザーでの並列実行
```bash
# Alice（立花証券）
./authenticated-agent.exe \
  --username=alice --password=alice123 \
  --broker=tachibana --broker-user=alice_tachibana --broker-password=tachibana123 \
  --strategy=swing --symbol=7203

# Bob（SBI証券）
./authenticated-agent.exe \
  --username=bob --password=bob123 \
  --broker=sbi --broker-user=bob_sbi --broker-password=sbi123 \
  --strategy=day --symbol=6758

# Charlie（楽天証券）
./authenticated-agent.exe \
  --username=charlie --password=charlie123 \
  --broker=rakuten --broker-user=charlie_rakuten --broker-password=rakuten123 \
  --strategy=simple --symbol=9984
```

### 2. セッション管理API
```bash
# ログイン
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "alice123",
    "broker_type": "tachibana",
    "broker_user_id": "alice_tachibana",
    "broker_password": "tachibana123"
  }'

# 認証付きAPI呼び出し
curl -X GET http://localhost:8080/trade/balance \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# セッション一覧
curl -X GET http://localhost:8080/auth/sessions \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# ログアウト
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## セキュリティ考慮事項

### 1. **JWT トークン**
- 短い有効期限（1-2時間）
- リフレッシュトークンによる自動更新
- 署名検証による改ざん防止

### 2. **セッション管理**
- セッションID の暗号学的安全性
- セッション有効期限の適切な管理
- 不正アクセス検知・ログ記録

### 3. **パスワード管理**
- bcrypt等による安全なハッシュ化
- 証券会社パスワードの暗号化保存
- 環境変数・設定ファイルでの管理

### 4. **通信セキュリティ**
- HTTPS必須（本番環境）
- CORS設定
- Rate Limiting

この設計により、複数ユーザー・エージェントが安全に同時利用できるセッション管理システムが実現できます。