// internal/infrastructure/client/test_helper.go
package client

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"stock-bot/internal/config"

	"go.uber.org/zap/zaptest"
)

// CreateTestClient はテスト用の TachibanaClient インスタンスを作成します。
func CreateTestClient(t *testing.T) *TachibanaClient {
	t.Helper()

	// .env ファイルのパスを修正
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get caller information")
	}
	// test_helper.go から見た .env の相対パス (プロジェクトルート)
	envPath := filepath.Join(filepath.Dir(filename), "../../../.env") // パスを修正

	// 設定ファイルの読み込み
	cfg, err := config.LoadConfig(envPath) // 絶対パスまたは相対パスを指定
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// デモ環境かどうかのチェックと表示
	if strings.Contains(cfg.TachibanaBaseURL, "demo") {
		fmt.Println("APIの　デモ　環境に接続")
	} else {
		fmt.Println("APIの　本番　環境に接続")
	}

	// ロガーの作成 (テスト用)
	logger := zaptest.NewLogger(t) // テストログを出力

	// TachibanaClient インスタンスの作成
	tachibanaClient := NewTachibanaClient(cfg, logger)

	return tachibanaClient
}

// GetLogginedForTest はテスト用にlogginedを取得
func (tc *TachibanaClient) GetLogginedForTest() bool { // レシーバを変更
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loggined
}

// GetLoginInfoForTest はテスト用に loginInfo を取得 (テストヘルパー)
func (tc *TachibanaClient) GetLoginInfoForTest() *LoginInfo { // レシーバを変更
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loginInfo
}

// GetUserIDForTest はテスト用に userID を取得します。
func (tc *TachibanaClient) GetUserIDForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.sUserId
}

// GetPasswordForTest はテスト用に password を取得します。
func (tc *TachibanaClient) GetPasswordForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.sPassword
}

// GetBaseURLForTest はテスト用に baseURL を取得します。
func (tc *TachibanaClient) GetBaseURLForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.baseURL.String() // 文字列で返す
}

// SetBaseURLForTest はテスト用に baseURL を設定します。
func (tc *TachibanaClient) SetBaseURLForTest(baseURL string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	parsedURL, _ := url.Parse(baseURL) // 文字列から *url.URL に変換
	tc.baseURL = parsedURL
}

// SetUserIDForTest はテスト用に userID を設定します。
func (tc *TachibanaClient) SetUserIDForTest(userID string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.sUserId = userID
}

// SetPasswordForTest はテスト用に password を設定します。
func (tc *TachibanaClient) SetPasswordForTest(password string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.sPassword = password
}

// GetBaseURLForTest はテスト用に baseURL を取得します。
func (tc *TachibanaClient) GetLastRequestURLForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loginInfo.RequestURL // 文字列で返す
}

// GetBaseURLForTest はテスト用に baseURL を取得します。
func (tc *TachibanaClient) FormatSDDateForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return formatSDDate(time.Now())

}
