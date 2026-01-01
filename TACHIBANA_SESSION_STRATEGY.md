# 立花証券セッション管理戦略

## 立花証券の制約事項

### **1. セッション排他性**
- 新しいログインで以前のセッションが無効化
- 同時に複数のセッションを維持不可
- セッション競合によるAPI呼び出し失敗

### **2. 本番環境の制約**
- 電話による二要素認証が必要
- ログイン頻度を最小限に抑制する必要
- 自動再ログインは実質的に不可能

## 現在の実装の問題点

### **問題1: 自動再認証ロジック**
```go
// 現在の実装（問題あり）
func (c *TachibanaUnifiedClient) EnsureAuthenticated(ctx context.Context) error {
    // 8時間経過で自動再ログイン → 本番では不可能
    if c.session != nil && time.Since(c.lastLoginTime) < 8*time.Hour {
        return nil
    }
    
    // 自動ログイン実行 → 電話認証で停止
    session, err := c.authClient.LoginWithPost(ctx, loginReq)
    // ...
}
```

### **問題2: セッション競合**
```go
// 複数のAPIリクエストが同時実行される場合
goroutine1: GetBalance() → EnsureAuthenticated() → Login()
goroutine2: GetOrders()  → EnsureAuthenticated() → Login() // 前のセッション無効化
goroutine1: API呼び出し → セッション無効エラー
```

### **問題3: セッション永続化なし**
- サーバー再起動でセッション消失
- 手動ログインが必要

## 改善されたセッション管理戦略

### **戦略1: セッション永続化 + 手動ログイン**

```go
// internal/session/tachibana_session_manager.go
package session

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "stock-bot/internal/infrastructure/client"
    "sync"
    "time"
)

type TachibanaSessionManager struct {
    session      *client.Session
    sessionMutex sync.RWMutex
    logger       *slog.Logger
    
    // セッション永続化
    sessionStore SessionStore
    
    // セッション状態管理
    isAuthenticated bool
    lastUsedAt      time.Time
    
    // ログイン制御
    loginInProgress bool
    loginMutex      sync.Mutex
}

type SessionStore interface {
    SaveSession(session *client.Session) error
    LoadSession() (*client.Session, error)
    DeleteSession() error
}

func NewTachibanaSessionManager(sessionStore SessionStore, logger *slog.Logger) *TachibanaSessionManager {
    return &TachibanaSessionManager{
        sessionStore: sessionStore,
        logger:       logger,
    }
}

// Initialize はサーバー起動時にセッションを復元する
func (tsm *TachibanaSessionManager) Initialize(ctx context.Context) error {
    tsm.sessionMutex.Lock()
    defer tsm.sessionMutex.Unlock()
    
    // 永続化されたセッションを復元
    session, err := tsm.sessionStore.LoadSession()
    if err != nil {
        tsm.logger.Warn("No saved session found", "error", err)
        return nil // エラーではない
    }
    
    tsm.session = session
    tsm.isAuthenticated = true
    tsm.lastUsedAt = time.Now()
    
    tsm.logger.Info("Session restored from storage")
    return nil
}

// GetSession はセッションを取得する（自動ログインなし）
func (tsm *TachibanaSessionManager) GetSession(ctx context.Context) (*client.Session, error) {
    tsm.sessionMutex.RLock()
    defer tsm.sessionMutex.RUnlock()
    
    if !tsm.isAuthenticated || tsm.session == nil {
        return nil, fmt.Errorf("not authenticated - manual login required")
    }
    
    // 最終使用時刻を更新
    tsm.lastUsedAt = time.Now()
    
    return tsm.session, nil
}

// ManualLogin は手動ログインを実行する（管理者用）
func (tsm *TachibanaSessionManager) ManualLogin(ctx context.Context, authClient client.AuthClient, userID, password string) error {
    tsm.loginMutex.Lock()
    defer tsm.loginMutex.Unlock()
    
    if tsm.loginInProgress {
        return fmt.Errorf("login already in progress")
    }
    
    tsm.loginInProgress = true
    defer func() { tsm.loginInProgress = false }()
    
    tsm.logger.Info("Starting manual login to Tachibana API")
    
    // ログイン実行
    loginReq := request.ReqLogin{
        UserId:   userID,
        Password: password,
    }
    
    session, err := authClient.LoginWithPost(ctx, loginReq)
    if err != nil {
        return fmt.Errorf("login failed: %w", err)
    }
    
    // セッション保存
    tsm.sessionMutex.Lock()
    tsm.session = session
    tsm.isAuthenticated = true
    tsm.lastUsedAt = time.Now()
    tsm.sessionMutex.Unlock()
    
    // セッション永続化
    if err := tsm.sessionStore.SaveSession(session); err != nil {
        tsm.logger.Warn("Failed to save session", "error", err)
    }
    
    tsm.logger.Info("Manual login successful")
    return nil
}

// Logout はセッションを終了する
func (tsm *TachibanaSessionManager) Logout(ctx context.Context, authClient client.AuthClient) error {
    tsm.sessionMutex.Lock()
    defer tsm.sessionMutex.Unlock()
    
    if tsm.session == nil {
        return nil
    }
    
    // ログアウト実行
    logoutReq := request.ReqLogout{}
    _, err := authClient.LogoutWithPost(ctx, tsm.session, logoutReq)
    if err != nil {
        tsm.logger.Warn("Logout request failed", "error", err)
    }
    
    // セッション削除
    tsm.session = nil
    tsm.isAuthenticated = false
    tsm.sessionStore.DeleteSession()
    
    tsm.logger.Info("Logout completed")
    return err
}

// GetSessionStatus はセッション状態を取得する
func (tsm *TachibanaSessionManager) GetSessionStatus() SessionStatus {
    tsm.sessionMutex.RLock()
    defer tsm.sessionMutex.RUnlock()
    
    return SessionStatus{
        IsAuthenticated: tsm.isAuthenticated,
        LastUsedAt:      tsm.lastUsedAt,
        LoginInProgress: tsm.loginInProgress,
    }
}

type SessionStatus struct {
    IsAuthenticated bool      `json:"is_authenticated"`
    LastUsedAt      time.Time `json:"last_used_at"`
    LoginInProgress bool      `json:"login_in_progress"`
}
```

### **戦略2: ファイルベースセッション永続化**

```go
// internal/session/file_session_store.go
package session

import (
    "encoding/json"
    "os"
    "stock-bot/internal/infrastructure/client"
)

type FileSessionStore struct {
    filePath string
}

func NewFileSessionStore(filePath string) *FileSessionStore {
    return &FileSessionStore{
        filePath: filePath,
    }
}

func (fss *FileSessionStore) SaveSession(session *client.Session) error {
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }
    
    return os.WriteFile(fss.filePath, data, 0600) // 読み書き権限を制限
}

func (fss *FileSessionStore) LoadSession() (*client.Session, error) {
    data, err := os.ReadFile(fss.filePath)
    if err != nil {
        return nil, err
    }
    
    var session client.Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, err
    }
    
    return &session, nil
}

func (fss *FileSessionStore) DeleteSession() error {
    return os.Remove(fss.filePath)
}
```

### **戦略3: 管理用API追加**

```go
// design/design.go に追加
var _ = Service("admin", func() {
    Description("Administrative functions for session management")

    // POST /admin/login
    Method("login", func() {
        Description("Manual login to broker (admin only)")
        Payload(func() {
            Attribute("user_id", String, "証券会社ユーザーID")
            Attribute("password", String, "証券会社パスワード")
            Attribute("second_password", String, "第二パスワード")
            Required("user_id", "password")
        })
        Result(func() {
            Attribute("success", Boolean, "ログイン成功")
            Attribute("message", String, "メッセージ")
            Required("success", "message")
        })
        HTTP(func() {
            POST("/admin/login")
            Response(StatusOK)
        })
    })

    // POST /admin/logout
    Method("logout", func() {
        Description("Logout from broker")
        Payload(Empty)
        Result(func() {
            Attribute("success", Boolean, "ログアウト成功")
            Attribute("message", String, "メッセージ")
            Required("success", "message")
        })
        HTTP(func() {
            POST("/admin/logout")
            Response(StatusOK)
        })
    })

    // GET /admin/session-status
    Method("session_status", func() {
        Description("Get current session status")
        Payload(Empty)
        Result(func() {
            Attribute("is_authenticated", Boolean, "認証状態")
            Attribute("last_used_at", String, "最終使用時刻")
            Attribute("login_in_progress", Boolean, "ログイン処理中")
            Required("is_authenticated")
        })
        HTTP(func() {
            GET("/admin/session-status")
            Response(StatusOK)
        })
    })
})
```

### **戦略4: 改善されたTachibanaUnifiedClient**

```go
// internal/infrastructure/client/tachibana_unified_client_v2.go
package client

import (
    "context"
    "fmt"
    "stock-bot/internal/session"
)

type TachibanaUnifiedClientV2 struct {
    authClient       AuthClient
    balanceClient    BalanceClient
    orderClient      OrderClient
    priceClient      PriceInfoClient
    masterClient     MasterDataClient
    eventClient      EventClient
    
    sessionManager   *session.TachibanaSessionManager
    logger           *slog.Logger
}

func NewTachibanaUnifiedClientV2(
    authClient AuthClient,
    balanceClient BalanceClient,
    orderClient OrderClient,
    priceClient PriceInfoClient,
    masterClient MasterDataClient,
    eventClient EventClient,
    sessionManager *session.TachibanaSessionManager,
    logger *slog.Logger,
) *TachibanaUnifiedClientV2 {
    return &TachibanaUnifiedClientV2{
        authClient:     authClient,
        balanceClient:  balanceClient,
        orderClient:    orderClient,
        priceClient:    priceClient,
        masterClient:   masterClient,
        eventClient:    eventClient,
        sessionManager: sessionManager,
        logger:         logger,
    }
}

// GetZanKaiSummary は残高サマリーを取得する（自動ログインなし）
func (c *TachibanaUnifiedClientV2) GetZanKaiSummary(ctx context.Context) (*balance_response.ResZanKaiSummary, error) {
    session, err := c.sessionManager.GetSession(ctx)
    if err != nil {
        return nil, fmt.Errorf("session not available: %w", err)
    }
    
    result, err := c.balanceClient.GetZanKaiSummary(ctx, session)
    if err != nil {
        // セッションエラーの場合、詳細ログを出力
        c.logger.Error("API call failed - session may be invalid", 
            "error", err,
            "api", "GetZanKaiSummary")
        return nil, fmt.Errorf("API call failed: %w", err)
    }
    
    return result, nil
}

// ManualLogin は管理者による手動ログインを実行する
func (c *TachibanaUnifiedClientV2) ManualLogin(ctx context.Context, userID, password, secondPassword string) error {
    return c.sessionManager.ManualLogin(ctx, c.authClient, userID, password)
}

// Logout はセッションを終了する
func (c *TachibanaUnifiedClientV2) Logout(ctx context.Context) error {
    return c.sessionManager.Logout(ctx, c.authClient)
}

// GetSessionStatus はセッション状態を取得する
func (c *TachibanaUnifiedClientV2) GetSessionStatus() session.SessionStatus {
    return c.sessionManager.GetSessionStatus()
}
```

## 運用フロー

### **1. 初回セットアップ**
```bash
# 1. Goaサービス起動
./goa-service.exe

# 2. 管理者による手動ログイン（電話認証含む）
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "your_tachibana_id",
    "password": "your_password",
    "second_password": "your_second_password"
  }'

# 3. セッション状態確認
curl http://localhost:8080/admin/session-status

# 4. エージェント起動
./lightweight-agent.exe --strategy=simple --symbol=7203
```

### **2. 日常運用**
```bash
# セッション状態監視
curl http://localhost:8080/admin/session-status

# セッション無効時の対応
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -d '{"user_id":"...","password":"..."}'
```

### **3. セッション永続化**
```
# セッションファイル保存場所
./data/tachibana_session.json

# サーバー再起動時
1. セッションファイルから復元
2. 無効な場合は手動ログイン要求
3. 有効な場合はそのまま利用継続
```

## エラーハンドリング戦略

### **1. セッション無効エラー**
```go
func (s *GoaTradeService) GetBalance(ctx context.Context) (*service.Balance, error) {
    session, err := s.sessionManager.GetSession(ctx)
    if err != nil {
        return nil, &SessionRequiredError{
            Message: "Manual login required",
            Action:  "POST /admin/login",
        }
    }
    
    balance, err := s.balanceClient.GetZanKaiSummary(ctx, session)
    if err != nil {
        // セッション無効の可能性
        if isSessionError(err) {
            return nil, &SessionInvalidError{
                Message: "Session may be invalid - manual re-login required",
                Action:  "POST /admin/login",
            }
        }
        return nil, err
    }
    
    return convertBalance(balance), nil
}
```

### **2. エラーレスポンス**
```json
{
  "error": "session_required",
  "message": "Manual login required",
  "action": "POST /admin/login",
  "details": {
    "last_used_at": "2024-01-01T10:00:00Z",
    "is_authenticated": false
  }
}
```

## 監視・アラート

### **1. セッション監視**
```go
// セッション状態を定期監視
func (tsm *TachibanaSessionManager) StartMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            status := tsm.GetSessionStatus()
            
            // 長時間未使用の警告
            if time.Since(status.LastUsedAt) > 6*time.Hour {
                tsm.logger.Warn("Session unused for long time", 
                    "last_used", status.LastUsedAt)
            }
            
            // セッション無効の警告
            if !status.IsAuthenticated {
                tsm.logger.Error("Session not authenticated - manual login required")
            }
        }
    }()
}
```

### **2. ヘルスチェック拡張**
```go
func (s *GoaTradeService) HealthCheck(ctx context.Context) (*service.HealthStatus, error) {
    sessionStatus := s.sessionManager.GetSessionStatus()
    
    status := "healthy"
    if !sessionStatus.IsAuthenticated {
        status = "unhealthy"
    } else if time.Since(sessionStatus.LastUsedAt) > 6*time.Hour {
        status = "degraded"
    }
    
    return &service.HealthStatus{
        Status:             status,
        Timestamp:          time.Now(),
        SessionValid:       sessionStatus.IsAuthenticated,
        DatabaseConnected:  true,
        WebSocketConnected: true,
        SessionLastUsed:    sessionStatus.LastUsedAt,
    }, nil
}
```

この戦略により、立花証券の制約に対応した安全で実用的なセッション管理が実現できます。