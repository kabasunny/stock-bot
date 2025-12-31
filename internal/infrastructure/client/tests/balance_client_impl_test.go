package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBalanceClient_GetZanKaiSummary は残高サマリー取得をテストします
func TestBalanceClient_GetZanKaiSummary(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// まずログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err, "Login should not produce an error")
	require.NotNil(t, session, "Session should not be nil")
	t.Logf("Login successful for balance test")

	// 残高サマリー取得
	summary, err := c.GetZanKaiSummary(context.Background(), session)

	// 結果の確認
	require.NoError(t, err, "GetZanKaiSummary should not produce an error")
	require.NotNil(t, summary, "Summary should not be nil")

	// 基本フィールドの存在確認
	assert.NotEmpty(t, summary.GenbutuKabuKaituke, "GenbutuKabuKaituke should not be empty")
	assert.NotEmpty(t, summary.SinyouSinkidate, "SinyouSinkidate should not be empty")
	assert.NotEmpty(t, summary.HosyouKinritu, "HosyouKinritu should not be empty")
	assert.NotEmpty(t, summary.Syukkin, "Syukkin should not be empty")
	assert.NotEmpty(t, summary.OisyouHasseiFlg, "OisyouHasseiFlg should not be empty")

	t.Logf("Balance Summary - GenbutuKaituke: %s, SinyouSinkidate: %s",
		summary.GenbutuKabuKaituke, summary.SinyouSinkidate)
}

// TestBalanceClient_WithoutLogin は未ログイン状態でのエラーテストです
func TestBalanceClient_WithoutLogin(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ダミーセッションを作成（実際にはログインしていない）
	dummySession := client.NewSession()
	dummySession.RequestURL = "https://example.com/dummy"

	// 残高サマリー取得試行
	summary, err := c.GetZanKaiSummary(context.Background(), dummySession)

	// エラーが発生することを確認
	require.Error(t, err, "GetZanKaiSummary with invalid session should produce an error")
	require.Nil(t, summary, "Summary should be nil when error occurs")

	t.Logf("Expected error with invalid session: %v", err)
}
