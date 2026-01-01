package client

import (
	"context"
	"log/slog"
	"time"
)

// SessionManagerType はセッション管理の種類を定義
type SessionManagerType string

const (
	SessionManagerTypeTime SessionManagerType = "time"
	SessionManagerTypeDate SessionManagerType = "date"
)

// SessionManagerConfig はセッション管理の設定
type SessionManagerConfig struct {
	Type           SessionManagerType `json:"type"`            // "time" または "date"
	SessionTimeout time.Duration      `json:"session_timeout"` // 時間ベース管理の場合のタイムアウト（デフォルト: 8時間）
	SessionDir     string             `json:"session_dir"`     // 日付ベース管理の場合のセッションディレクトリ
}

// SessionManager はセッション管理の共通インターフェース
type SessionManager interface {
	EnsureAuthenticated(ctx context.Context) error
	GetSession(ctx context.Context) (*Session, error)
	IsAuthenticated() bool
	Logout(ctx context.Context) error
}

// SessionManagerFactory はセッション管理を作成するファクトリー
type SessionManagerFactory struct {
	config SessionManagerConfig
	logger *slog.Logger
}

// NewSessionManagerFactory は新しいセッション管理ファクトリーを作成
func NewSessionManagerFactory(config SessionManagerConfig, logger *slog.Logger) *SessionManagerFactory {
	// デフォルト値の設定
	if config.SessionTimeout == 0 {
		config.SessionTimeout = 8 * time.Hour
	}
	if config.SessionDir == "" {
		config.SessionDir = "./data/sessions"
	}

	return &SessionManagerFactory{
		config: config,
		logger: logger,
	}
}

// CreateSessionManager は設定に基づいてセッション管理を作成
func (f *SessionManagerFactory) CreateSessionManager(
	authClient AuthClient,
	userID, password, secondPassword string,
) SessionManager {
	switch f.config.Type {
	case SessionManagerTypeDate:
		f.logger.Info("creating date-based session manager")
		return NewDateBasedSessionManager(
			authClient,
			userID,
			password,
			secondPassword,
			f.config.SessionDir,
			f.logger,
		)
	case SessionManagerTypeTime:
		f.logger.Info("creating time-based session manager",
			"timeout", f.config.SessionTimeout)
		return NewTimeBasedSessionManager(
			authClient,
			userID,
			password,
			secondPassword,
			f.config.SessionTimeout,
			f.logger,
		)
	default:
		f.logger.Warn("unknown session manager type, using time-based",
			"type", f.config.Type)
		return NewTimeBasedSessionManager(
			authClient,
			userID,
			password,
			secondPassword,
			f.config.SessionTimeout,
			f.logger,
		)
	}
}

// GetDefaultConfig はデフォルトのセッション管理設定を取得
func GetDefaultConfig() SessionManagerConfig {
	return SessionManagerConfig{
		Type:           SessionManagerTypeDate, // デフォルトは日付ベース
		SessionTimeout: 8 * time.Hour,
		SessionDir:     "./data/sessions",
	}
}

// GetTimeBasedConfig は時間ベースのセッション管理設定を取得
func GetTimeBasedConfig(timeout time.Duration) SessionManagerConfig {
	return SessionManagerConfig{
		Type:           SessionManagerTypeTime,
		SessionTimeout: timeout,
		SessionDir:     "./data/sessions",
	}
}

// GetDateBasedConfig は日付ベースのセッション管理設定を取得
func GetDateBasedConfig(sessionDir string) SessionManagerConfig {
	return SessionManagerConfig{
		Type:           SessionManagerTypeDate,
		SessionTimeout: 8 * time.Hour, // 使用されないが設定
		SessionDir:     sessionDir,
	}
}
