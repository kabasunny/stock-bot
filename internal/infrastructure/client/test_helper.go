// internal/infrastructure/client/test_helper.go
package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"stock-bot/internal/config"
)

const DummyNewOrderResponse = `{
    "ResultCode": "0",
    "OrderNumber": "123456789",
    "EigyouDay": "2025/12/09"
}`

// CreateTestClientWithServer はテスト用のTachibanaClientとテストサーバーを作成します
func CreateTestClientWithServer(t *testing.T, handler http.HandlerFunc) (*TachibanaClientImpl, *httptest.Server) {
	t.Helper()

	// テストサーバーを起動
	server := httptest.NewServer(handler)

	// .env ファイルのパスを修正
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get caller information")
	}
	envPath := filepath.Join(filepath.Dir(filename), "../../../.env")

	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	// baseURLをテストサーバーのURLに上書き
	cfg.TachibanaBaseURL = server.URL

	tachibanaClient := NewTachibanaClient(cfg)

	return tachibanaClient, server
}

// CreateTestClient はテスト用の TachibanaClient インスタンスを作成
func CreateTestClient(t *testing.T) *TachibanaClientImpl {
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
		t.Log("APIの　デモ　環境に接続")
	} else {
		t.Log("APIの　本番　環境に接続")
	}

	// TachibanaClient インスタンスの作成
	tachibanaClient := NewTachibanaClient(cfg)

	return tachibanaClient
}

// GetLogginedForTest はテスト用にlogginedを取得
func (tc *TachibanaClientImpl) GetLogginedForTest() bool { // レシーバを変更
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loggined
}

// GetLoginInfoForTest はテスト用に loginInfo を取得 (テストヘルパー)
func (tc *TachibanaClientImpl) GetLoginInfoForTest() *LoginInfo { // レシーバを変更
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loginInfo
}

// GetUserIDForTest はテスト用に userID を取得します。
func (tc *TachibanaClientImpl) GetUserIDForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.sUserId
}

// GetPasswordForTest はテスト用に password を取得します。
func (tc *TachibanaClientImpl) GetPasswordForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.sPassword
}

// GetBaseURLForTest はテスト用に baseURL を取得します。
func (tc *TachibanaClientImpl) GetBaseURLForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.baseURL.String() // 文字列で返す
}

// SetBaseURLForTest はテスト用に baseURL を設定します。
func (tc *TachibanaClientImpl) SetBaseURLForTest(baseURL string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	parsedURL, _ := url.Parse(baseURL) // 文字列から *url.URL に変換
	tc.baseURL = parsedURL
}

// SetUserIDForTest はテスト用に userID を設定します。
func (tc *TachibanaClientImpl) SetUserIDForTest(userID string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.sUserId = userID
}

// SetPasswordForTest はテスト用に password を設定します。
func (tc *TachibanaClientImpl) SetPasswordForTest(password string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.sPassword = password
}

// GetBaseURLForTest はテスト用に baseURL を取得します。
func (tc *TachibanaClientImpl) GetLastRequestURLForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.loginInfo.RequestURL // 文字列で返す
}

// GetBaseURLForTest はテスト用に baseURL を取得します。
func (tc *TachibanaClientImpl) FormatSDDateForTest() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return formatSDDate(time.Now())

}

// SetLoginStateForTest は、テスト用に手動でログイン状態を設定します。
func (tc *TachibanaClientImpl) SetLoginStateForTest(t *testing.T, sessionCookie, priceURL string, pNo int64) {
	t.Helper()
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// 1. ログイン状態をtrueに設定
	tc.loggined = true

	// 2. p_no を設定
	tc.p_no = pNo

	// 3. LoginInfo を設定 (PriceURLのみで初期化)
	tc.loginInfo = &LoginInfo{
		PriceURL: priceURL,
		// 他のURLも必要であれば、引数に追加する
	}

	// 4. Cookieを設定
	header := http.Header{}
	header.Add("Set-Cookie", sessionCookie)
	cookies := (&http.Response{Header: header}).Cookies()
	tc.httpClient.Jar.SetCookies(tc.baseURL, cookies)

	t.Log("手動でクライアントのログイン状態を設定しました。")
}

// GetPNoForTest はテスト用に p_no を取得します。
func (tc *TachibanaClientImpl) GetPNoForTest() int64 {
	tc.p_NoMu.Lock()
	defer tc.p_NoMu.Unlock()
	return tc.p_no
}
