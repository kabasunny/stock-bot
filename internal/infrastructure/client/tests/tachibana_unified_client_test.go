package tests

import (
	"context"
	"log/slog"
	"stock-bot/internal/infrastructure/client"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTachibanaUnifiedClient_NewClient は、UnifiedClientの作成をテストします
func TestTachibanaUnifiedClient_NewClient(t *testing.T) {
	// テスト用クライアントを作成
	tachibanaClient := client.CreateTestClient(t)
	eventClient := client.NewEventClient(slog.Default())

	// UnifiedClientを作成
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		tachibanaClient.GetUserIDForTest(),
		tachibanaClient.GetPasswordForTest(),
		"dummy_second_password",
		slog.Default(),
	)

	// UnifiedClientが正しく作成されることを確認
	assert.NotNil(t, unifiedClient, "UnifiedClient should not be nil")

	// 初期状態では認証されていないことを確認
	assert.False(t, unifiedClient.IsAuthenticated(), "Should not be authenticated initially")
}

// TestTachibanaUnifiedClient_GetSession は、GetSession()の自動認証機能をテストします
func TestTachibanaUnifiedClient_GetSession(t *testing.T) {
	// テスト用クライアントを作成
	tachibanaClient := client.CreateTestClient(t)
	eventClient := client.NewEventClient(slog.Default())

	// UnifiedClientを作成
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		tachibanaClient.GetUserIDForTest(),
		tachibanaClient.GetPasswordForTest(),
		"dummy_second_password",
		slog.Default(),
	)

	// GetSession()を呼び出し（自動認証が実行される）
	session, err := unifiedClient.GetSession(context.Background())

	// 認証成功を確認
	require.NoError(t, err, "GetSession should not produce an error")
	require.NotNil(t, session, "Session should not be nil")
	assert.Equal(t, "0", session.ResultCode, "Session ResultCode should be 0")

	// 認証状態の確認
	assert.True(t, unifiedClient.IsAuthenticated(), "Should be authenticated after GetSession")

	// セッション内容の確認
	assert.NotEmpty(t, session.RequestURL, "RequestURL should not be empty")
	assert.NotEmpty(t, session.SecondPassword, "SecondPassword should be set")
	assert.Equal(t, "dummy_second_password", session.SecondPassword, "SecondPassword should match")

	t.Logf("GetSession successful - RequestURL: %s", session.RequestURL)
}

// TestTachibanaUnifiedClient_EnsureAuthenticated は、EnsureAuthenticated()の動作をテストします
func TestTachibanaUnifiedClient_EnsureAuthenticated(t *testing.T) {
	// テスト用クライアントを作成
	tachibanaClient := client.CreateTestClient(t)
	eventClient := client.NewEventClient(slog.Default())

	// UnifiedClientを作成
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		tachibanaClient.GetUserIDForTest(),
		tachibanaClient.GetPasswordForTest(),
		"dummy_second_password",
		slog.Default(),
	)

	// 初回認証
	err := unifiedClient.EnsureAuthenticated(context.Background())
	require.NoError(t, err, "First EnsureAuthenticated should not produce an error")
	assert.True(t, unifiedClient.IsAuthenticated(), "Should be authenticated after EnsureAuthenticated")

	// 2回目の呼び出し（既に認証済みなので再認証は行われない）
	err = unifiedClient.EnsureAuthenticated(context.Background())
	require.NoError(t, err, "Second EnsureAuthenticated should not produce an error")
	assert.True(t, unifiedClient.IsAuthenticated(), "Should still be authenticated")
}

// TestTachibanaUnifiedClient_MultipleGetSession は、複数回のGetSession呼び出しをテストします
func TestTachibanaUnifiedClient_MultipleGetSession(t *testing.T) {
	// テスト用クライアントを作成
	tachibanaClient := client.CreateTestClient(t)
	eventClient := client.NewEventClient(slog.Default())

	// UnifiedClientを作成
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		tachibanaClient.GetUserIDForTest(),
		tachibanaClient.GetPasswordForTest(),
		"dummy_second_password",
		slog.Default(),
	)

	// 1回目のGetSession
	session1, err1 := unifiedClient.GetSession(context.Background())
	require.NoError(t, err1, "First GetSession should not produce an error")
	require.NotNil(t, session1, "First session should not be nil")

	// 2回目のGetSession（同じセッションが返されるはず）
	session2, err2 := unifiedClient.GetSession(context.Background())
	require.NoError(t, err2, "Second GetSession should not produce an error")
	require.NotNil(t, session2, "Second session should not be nil")

	// 同じセッションオブジェクトが返されることを確認
	assert.Same(t, session1, session2, "Should return the same session object")
	assert.Equal(t, session1.RequestURL, session2.RequestURL, "RequestURL should be the same")
}

// TestTachibanaUnifiedClient_Logout は、Logout()の動作をテストします
func TestTachibanaUnifiedClient_Logout(t *testing.T) {
	// テスト用クライアントを作成
	tachibanaClient := client.CreateTestClient(t)
	eventClient := client.NewEventClient(slog.Default())

	// UnifiedClientを作成
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		tachibanaClient.GetUserIDForTest(),
		tachibanaClient.GetPasswordForTest(),
		"dummy_second_password",
		slog.Default(),
	)

	// まずログイン
	session, err := unifiedClient.GetSession(context.Background())
	require.NoError(t, err, "Login should not produce an error")
	require.NotNil(t, session, "Session should not be nil")
	assert.True(t, unifiedClient.IsAuthenticated(), "Should be authenticated after login")

	// ログアウト
	err = unifiedClient.Logout(context.Background())
	require.NoError(t, err, "Logout should not produce an error")

	// ログアウト後の状態確認
	assert.False(t, unifiedClient.IsAuthenticated(), "Should not be authenticated after logout")

	// ログアウト後のGetSession（再認証が実行される）
	newSession, err := unifiedClient.GetSession(context.Background())
	require.NoError(t, err, "GetSession after logout should not produce an error")
	require.NotNil(t, newSession, "New session should not be nil")
	assert.True(t, unifiedClient.IsAuthenticated(), "Should be authenticated again after GetSession")

	// 新しいセッションは異なるオブジェクトであることを確認
	assert.NotSame(t, session, newSession, "New session should be a different object")
}

// TestTachibanaUnifiedClient_InvalidCredentials は、不正な認証情報でのテストです
func TestTachibanaUnifiedClient_InvalidCredentials(t *testing.T) {
	// テスト用クライアントを作成
	tachibanaClient := client.CreateTestClient(t)
	eventClient := client.NewEventClient(slog.Default())

	// 不正な認証情報でUnifiedClientを作成
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		"invalid_user",
		"invalid_password",
		"dummy_second_password",
		slog.Default(),
	)

	// GetSession()を呼び出し（認証失敗が期待される）
	session, err := unifiedClient.GetSession(context.Background())

	// 認証失敗を確認
	require.Error(t, err, "GetSession with invalid credentials should produce an error")
	require.Nil(t, session, "Session should be nil when authentication fails")
	assert.False(t, unifiedClient.IsAuthenticated(), "Should not be authenticated with invalid credentials")

	// エラーメッセージの確認
	assert.Contains(t, err.Error(), "failed to authenticate", "Error message should contain authentication failure")

	t.Logf("Expected authentication failure: %v", err)
}

// TestTachibanaUnifiedClient_LogoutWithoutLogin は、ログインせずにログアウトを試行するテストです
func TestTachibanaUnifiedClient_LogoutWithoutLogin(t *testing.T) {
	// テスト用クライアントを作成
	tachibanaClient := client.CreateTestClient(t)
	eventClient := client.NewEventClient(slog.Default())

	// UnifiedClientを作成
	unifiedClient := client.NewTachibanaUnifiedClient(
		tachibanaClient, // AuthClient
		tachibanaClient, // BalanceClient
		tachibanaClient, // OrderClient
		tachibanaClient, // PriceInfoClient
		tachibanaClient, // MasterDataClient
		eventClient,     // EventClient
		tachibanaClient.GetUserIDForTest(),
		tachibanaClient.GetPasswordForTest(),
		"dummy_second_password",
		slog.Default(),
	)

	// ログインせずにログアウト試行
	err := unifiedClient.Logout(context.Background())

	// ログアウトはエラーにならない（既にログアウト済みとして扱われる）
	require.NoError(t, err, "Logout without login should not produce an error")
	assert.False(t, unifiedClient.IsAuthenticated(), "Should not be authenticated")
}
