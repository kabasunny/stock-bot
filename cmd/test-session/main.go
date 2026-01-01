package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"stock-bot/internal/config"
	"stock-bot/internal/infrastructure/client"
	"time"
)

func main() {
	// ロガーの設定
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// 設定の読み込み
	cfg, err := config.LoadConfig("agent_config.yaml")
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// セッションディレクトリの確保
	sessionDir := "./data/sessions"
	if err := client.EnsureSessionDirectory(sessionDir); err != nil {
		logger.Error("failed to create session directory", "error", err)
		os.Exit(1)
	}

	// 基本クライアントの初期化
	tachibanaClient := client.NewTachibanaClient(cfg)

	// 日付ベースセッション管理の初期化
	sessionManager := client.NewDateBasedSessionManager(
		tachibanaClient, // AuthClient
		cfg.TachibanaUserID,
		cfg.TachibanaPassword,
		cfg.TachibanaSecondPassword,
		sessionDir,
		logger,
	)

	ctx := context.Background()

	// セッション管理のテスト
	logger.Info("=== 日付ベースセッション管理テスト開始 ===")

	// 1. 初回認証テスト
	logger.Info("1. 初回認証テスト")
	session1, err := sessionManager.GetSession(ctx)
	if err != nil {
		logger.Error("初回認証に失敗", "error", err)
		os.Exit(1)
	}
	logger.Info("初回認証成功", "session_valid", session1 != nil)

	// 2. セッション復元テスト
	logger.Info("2. セッション復元テスト")
	session2, err := sessionManager.GetSession(ctx)
	if err != nil {
		logger.Error("セッション復元に失敗", "error", err)
		os.Exit(1)
	}
	logger.Info("セッション復元成功", "same_session", session1 == session2)

	// 3. 認証状態確認
	logger.Info("3. 認証状態確認")
	isAuth := sessionManager.IsAuthenticated()
	logger.Info("認証状態", "authenticated", isAuth)

	// 4. セッションファイル確認
	logger.Info("4. セッションファイル確認")
	currentDate := getCurrentBusinessDate()
	sessionFile := fmt.Sprintf("%s/tachibana_session_%s.json", sessionDir, currentDate)
	if _, err := os.Stat(sessionFile); err == nil {
		logger.Info("セッションファイル存在確認", "file", sessionFile, "exists", true)
	} else {
		logger.Warn("セッションファイルが見つかりません", "file", sessionFile, "error", err)
	}

	// 5. 統合クライアントテスト
	logger.Info("5. 統合クライアントテスト")
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		nil,             // EventClient (テストでは不要)
		cfg.TachibanaUserID,
		cfg.TachibanaPassword,
		cfg.TachibanaSecondPassword,
		logger,
	)

	session3, err := unifiedClient.GetSession(ctx)
	if err != nil {
		logger.Error("統合クライアント認証に失敗", "error", err)
		os.Exit(1)
	}
	logger.Info("統合クライアント認証成功", "session_valid", session3 != nil)

	// 6. 複数回呼び出しテスト
	logger.Info("6. 複数回呼び出しテスト")
	for i := 0; i < 3; i++ {
		session, err := unifiedClient.GetSession(ctx)
		if err != nil {
			logger.Error("複数回呼び出しテスト失敗", "iteration", i, "error", err)
			continue
		}
		logger.Info("複数回呼び出し成功", "iteration", i, "session_valid", session != nil)
		time.Sleep(1 * time.Second)
	}

	logger.Info("=== 日付ベースセッション管理テスト完了 ===")
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
