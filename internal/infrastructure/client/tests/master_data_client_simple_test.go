package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	master_request "stock-bot/internal/infrastructure/client/dto/master/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMasterDataClient_DownloadMasterData_SystemStatus はシステムステータス取得をテストします
func TestMasterDataClient_DownloadMasterData_SystemStatus(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// まずログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err, "Login should not produce an error")
	require.NotNil(t, session, "Session should not be nil")
	t.Logf("Login successful for master data test")

	// システムステータスのダウンロード
	downloadReq := master_request.ReqDownloadMaster{
		TargetCLMID: "CLMSystemStatus,CLMEventDownloadComplete",
	}

	res, err := c.DownloadMasterData(context.Background(), session, downloadReq)

	// 結果の確認
	require.NoError(t, err, "DownloadMasterData should not produce an error")
	require.NotNil(t, res, "Response should not be nil")

	// SystemStatusの確認
	assert.Equal(t, "CLMSystemStatus", res.SystemStatus.CLMID, "SystemStatus CLMID should match")
	assert.NotEmpty(t, res.SystemStatus.SystemStatus, "SystemStatus should not be empty")

	t.Logf("SystemStatus - CLMID: %s, SystemStatus: %s",
		res.SystemStatus.CLMID, res.SystemStatus.SystemStatus)
}

// TestMasterDataClient_WithoutLogin は未ログイン状態でのエラーテストです
func TestMasterDataClient_WithoutLogin(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ダミーセッションを作成（実際にはログインしていない）
	dummySession := client.NewSession()
	dummySession.MasterURL = "https://example.com/dummy"

	// マスターデータダウンロード試行
	downloadReq := master_request.ReqDownloadMaster{
		TargetCLMID: "CLMSystemStatus,CLMEventDownloadComplete",
	}

	res, err := c.DownloadMasterData(context.Background(), dummySession, downloadReq)

	// エラーが発生することを確認
	require.Error(t, err, "DownloadMasterData with invalid session should produce an error")
	require.Nil(t, res, "Response should be nil when error occurs")

	t.Logf("Expected error with invalid session: %v", err)
}
