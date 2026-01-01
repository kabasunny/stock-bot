package client

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	"sync"
	"time"
)

// TimeBasedSessionManager は時間ベースのセッション管理を行う
type TimeBasedSessionManager struct {
	authClient     AuthClient
	userID         string
	password       string
	secondPassword string
	sessionTimeout time.Duration
	logger         *slog.Logger

	// セッション状態
	session       *Session
	lastLoginTime time.Time
	mutex         sync.RWMutex
}

// NewTimeBasedSessionManager は新しい時間ベースセッション管理を作成
func NewTimeBasedSessionManager(
	authClient AuthClient,
	userID, password, secondPassword string,
	sessionTimeout time.Duration,
	logger *slog.Logger,
) *TimeBasedSessionManager {
	return &TimeBasedSessionManager{
		authClient:     authClient,
		userID:         userID,
		password:       password,
		secondPassword: secondPassword,
		sessionTimeout: sessionTimeout,
		logger:         logger,
	}
}

// EnsureAuthenticated はセッションが有効であることを確認し、必要に応じて再認証を行う
func (tsm *TimeBasedSessionManager) EnsureAuthenticated(ctx context.Context) error {
	tsm.mutex.Lock()
	defer tsm.mutex.Unlock()

	// セッションが存在し、まだ有効な場合はそのまま使用
	if tsm.session != nil && time.Since(tsm.lastLoginTime) < tsm.sessionTimeout {
		return nil
	}

	tsm.logger.Info("performing authentication to Tachibana API",
		"timeout", tsm.sessionTimeout)

	// ログインリクエストを作成
	loginReq := request.ReqLogin{
		UserId:   tsm.userID,
		Password: tsm.password,
	}

	// 認証実行
	session, err := tsm.authClient.LoginWithPost(ctx, loginReq)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// セッション情報を設定
	session.SecondPassword = tsm.secondPassword

	tsm.session = session
	tsm.lastLoginTime = time.Now()

	tsm.logger.Info("authentication successful",
		"expires_at", tsm.lastLoginTime.Add(tsm.sessionTimeout))
	return nil
}

// GetSession は現在のセッションを取得
func (tsm *TimeBasedSessionManager) GetSession(ctx context.Context) (*Session, error) {
	if err := tsm.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	tsm.mutex.RLock()
	defer tsm.mutex.RUnlock()
	return tsm.session, nil
}

// IsAuthenticated はセッションが有効かどうかを確認
func (tsm *TimeBasedSessionManager) IsAuthenticated() bool {
	tsm.mutex.RLock()
	defer tsm.mutex.RUnlock()
	return tsm.session != nil && time.Since(tsm.lastLoginTime) < tsm.sessionTimeout
}

// Logout はセッションを終了
func (tsm *TimeBasedSessionManager) Logout(ctx context.Context) error {
	tsm.mutex.Lock()
	defer tsm.mutex.Unlock()

	if tsm.session == nil {
		return nil // 既にログアウト済み
	}

	// ログアウトリクエストを作成
	logoutReq := request.ReqLogout{}

	// ログアウト実行
	_, err := tsm.authClient.LogoutWithPost(ctx, tsm.session, logoutReq)
	if err != nil {
		tsm.logger.Warn("logout request failed", "error", err)
		// ログアウトエラーでもセッションはクリアする
	}

	tsm.session = nil
	tsm.lastLoginTime = time.Time{}

	tsm.logger.Info("logout completed")
	return err
}

// GetSessionInfo はセッション情報を取得（デバッグ用）
func (tsm *TimeBasedSessionManager) GetSessionInfo() map[string]interface{} {
	tsm.mutex.RLock()
	defer tsm.mutex.RUnlock()

	info := map[string]interface{}{
		"type":            "time-based",
		"timeout":         tsm.sessionTimeout,
		"has_session":     tsm.session != nil,
		"last_login_time": tsm.lastLoginTime,
	}

	if tsm.session != nil {
		remainingTime := tsm.sessionTimeout - time.Since(tsm.lastLoginTime)
		info["remaining_time"] = remainingTime
		info["expires_at"] = tsm.lastLoginTime.Add(tsm.sessionTimeout)
		info["is_expired"] = remainingTime <= 0
	}

	return info
}
