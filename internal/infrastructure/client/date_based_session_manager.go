package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	"sync"
	"time"
)

// SessionData はセッションファイルに保存するデータ構造
type SessionData struct {
	Session    *Session  `json:"session"`
	Date       string    `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

// DateBasedSessionManager は営業日ベースのセッション管理を行う
type DateBasedSessionManager struct {
	authClient     AuthClient
	userID         string
	password       string
	secondPassword string
	sessionDir     string
	logger         *slog.Logger

	// セッション状態
	session         *Session
	currentDate     string
	sessionDate     string
	isAuthenticated bool
	mutex           sync.RWMutex
}

// NewDateBasedSessionManager は新しい日付ベースセッション管理を作成
func NewDateBasedSessionManager(
	authClient AuthClient,
	userID, password, secondPassword string,
	sessionDir string,
	logger *slog.Logger,
) *DateBasedSessionManager {
	return &DateBasedSessionManager{
		authClient:     authClient,
		userID:         userID,
		password:       password,
		secondPassword: secondPassword,
		sessionDir:     sessionDir,
		logger:         logger,
		currentDate:    getCurrentBusinessDate(),
	}
}

// getCurrentBusinessDate は現在の営業日を取得（土日は前の金曜日）
func getCurrentBusinessDate() string {
	now := time.Now()

	// 土日の場合は前の金曜日を返す
	for now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		now = now.AddDate(0, 0, -1)
	}

	return now.Format("2006-01-02")
}

// EnsureAuthenticated はセッションが有効であることを確認し、必要に応じて再認証を行う
func (dsm *DateBasedSessionManager) EnsureAuthenticated(ctx context.Context) error {
	dsm.mutex.Lock()
	defer dsm.mutex.Unlock()

	// 営業日の変更をチェック
	if err := dsm.checkDateChange(ctx); err != nil {
		return fmt.Errorf("failed to handle date change: %w", err)
	}

	// 既にセッションが有効な場合
	if dsm.isAuthenticated && dsm.session != nil {
		return nil
	}

	// 当日のセッション復元を試行
	if err := dsm.loadTodaysSession(); err == nil {
		dsm.logger.Info("session restored from file", "date", dsm.currentDate)
		return nil
	}

	// 新しいログインが必要
	return dsm.performLogin(ctx)
}

// checkDateChange は営業日の変更をチェック
func (dsm *DateBasedSessionManager) checkDateChange(ctx context.Context) error {
	newDate := getCurrentBusinessDate()

	if newDate == dsm.currentDate {
		return nil // 日付変更なし
	}

	dsm.logger.Info("business date changed",
		"old_date", dsm.currentDate,
		"new_date", newDate)

	// 営業日が変わった場合、セッションを無効化
	dsm.invalidateSession()
	dsm.currentDate = newDate

	return nil
}

// loadTodaysSession は当日のセッションファイルから復元
func (dsm *DateBasedSessionManager) loadTodaysSession() error {
	filePath := dsm.getSessionFilePath(dsm.currentDate)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("no session file for date %s: %w", dsm.currentDate, err)
	}

	var sessionData SessionData
	if err := json.Unmarshal(data, &sessionData); err != nil {
		return fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// 日付の整合性チェック
	if sessionData.Date != dsm.currentDate {
		return fmt.Errorf("session date mismatch: expected %s, got %s",
			dsm.currentDate, sessionData.Date)
	}

	dsm.session = sessionData.Session
	dsm.isAuthenticated = true
	dsm.sessionDate = sessionData.Date

	// 最終使用時刻を更新
	dsm.updateLastUsedTime()

	return nil
}

// performLogin は新しいログインを実行
func (dsm *DateBasedSessionManager) performLogin(ctx context.Context) error {
	dsm.logger.Info("performing new login", "date", dsm.currentDate)

	// ログインリクエストを作成
	loginReq := request.ReqLogin{
		UserId:   dsm.userID,
		Password: dsm.password,
	}

	// 認証実行
	session, err := dsm.authClient.LoginWithPost(ctx, loginReq)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// セッション情報を設定
	session.SecondPassword = dsm.secondPassword

	dsm.session = session
	dsm.isAuthenticated = true
	dsm.sessionDate = dsm.currentDate

	// セッションファイルに保存
	if err := dsm.saveSession(); err != nil {
		dsm.logger.Warn("failed to save session to file", "error", err)
		// 保存エラーでもログインは成功とする
	}

	dsm.logger.Info("login successful", "date", dsm.currentDate)
	return nil
}

// saveSession はセッションをファイルに保存
func (dsm *DateBasedSessionManager) saveSession() error {
	if dsm.session == nil {
		return fmt.Errorf("no session to save")
	}

	// セッションディレクトリを作成
	if err := os.MkdirAll(dsm.sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	sessionData := SessionData{
		Session:    dsm.session,
		Date:       dsm.currentDate,
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}

	data, err := json.MarshalIndent(sessionData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	filePath := dsm.getSessionFilePath(dsm.currentDate)
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	dsm.logger.Info("session saved to file", "path", filePath)
	return nil
}

// updateLastUsedTime は最終使用時刻を更新
func (dsm *DateBasedSessionManager) updateLastUsedTime() {
	filePath := dsm.getSessionFilePath(dsm.currentDate)

	// 既存ファイルを読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		dsm.logger.Warn("failed to read session file for update", "error", err)
		return
	}

	var sessionData SessionData
	if err := json.Unmarshal(data, &sessionData); err != nil {
		dsm.logger.Warn("failed to unmarshal session for update", "error", err)
		return
	}

	// 最終使用時刻を更新
	sessionData.LastUsedAt = time.Now()

	// ファイルに書き戻し
	updatedData, err := json.MarshalIndent(sessionData, "", "  ")
	if err != nil {
		dsm.logger.Warn("failed to marshal updated session", "error", err)
		return
	}

	if err := os.WriteFile(filePath, updatedData, 0600); err != nil {
		dsm.logger.Warn("failed to update session file", "error", err)
	}
}

// getSessionFilePath は指定日付のセッションファイルパスを取得
func (dsm *DateBasedSessionManager) getSessionFilePath(date string) string {
	filename := fmt.Sprintf("tachibana_session_%s.json", date)
	return filepath.Join(dsm.sessionDir, filename)
}

// GetSession は現在のセッションを取得
func (dsm *DateBasedSessionManager) GetSession(ctx context.Context) (*Session, error) {
	if err := dsm.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	dsm.mutex.RLock()
	defer dsm.mutex.RUnlock()

	// 使用時刻を更新（非同期）
	go dsm.updateLastUsedTime()

	return dsm.session, nil
}

// invalidateSession はセッションを無効化
func (dsm *DateBasedSessionManager) invalidateSession() {
	dsm.session = nil
	dsm.isAuthenticated = false
	dsm.sessionDate = ""
}

// IsAuthenticated はセッションが有効かどうかを確認
func (dsm *DateBasedSessionManager) IsAuthenticated() bool {
	dsm.mutex.RLock()
	defer dsm.mutex.RUnlock()
	return dsm.isAuthenticated && dsm.session != nil
}

// Logout はセッションを終了
func (dsm *DateBasedSessionManager) Logout(ctx context.Context) error {
	dsm.mutex.Lock()
	defer dsm.mutex.Unlock()

	if dsm.session == nil {
		return nil // 既にログアウト済み
	}

	// ログアウトリクエストを作成
	logoutReq := request.ReqLogout{}

	// ログアウト実行
	_, err := dsm.authClient.LogoutWithPost(ctx, dsm.session, logoutReq)
	if err != nil {
		dsm.logger.Warn("logout request failed", "error", err)
	}

	// セッションファイルを削除
	filePath := dsm.getSessionFilePath(dsm.currentDate)
	if err := os.Remove(filePath); err != nil {
		dsm.logger.Warn("failed to remove session file", "error", err)
	}

	dsm.invalidateSession()
	dsm.logger.Info("logout completed")

	return err
}
