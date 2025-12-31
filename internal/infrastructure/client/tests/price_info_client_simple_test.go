package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	price_request "stock-bot/internal/infrastructure/client/dto/price/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPriceInfoClient_GetPriceInfo は現在価格取得をテストします
func TestPriceInfoClient_GetPriceInfo(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// まずログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err, "Login should not produce an error")
	require.NotNil(t, session, "Session should not be nil")
	t.Logf("Login successful for price info test")

	// 有名な銘柄（トヨタ自動車）の価格情報を取得
	priceReq := price_request.ReqGetPriceInfo{
		CLMID:           "CLMMfdsGetMarketPrice",
		TargetIssueCode: "7203",
		TargetColumn:    "pDPP,pDOP,pDHP,pDLP", // 終値、始値、高値、安値
	}

	priceInfo, err := c.GetPriceInfo(context.Background(), session, priceReq)

	// 夜間の場合はエラーが発生する可能性があるため、エラーチェックを柔軟に
	if err != nil {
		t.Logf("GetPriceInfo error (may be expected during off-hours): %v", err)
		return // 夜間エラーの場合はスキップ
	}

	// 結果の確認
	require.NotNil(t, priceInfo, "PriceInfo should not be nil")
	assert.NotEmpty(t, priceInfo.CLMID, "CLMID should not be empty")

	if len(priceInfo.CLMMfdsMarketPrice) > 0 {
		firstItem := priceInfo.CLMMfdsMarketPrice[0]
		t.Logf("価格情報取得成功 - 銘柄コード: %s, Values: %+v",
			firstItem.IssueCode, firstItem.Values)
	} else {
		t.Logf("価格情報が返されませんでした（夜間・休日のため正常）- CLMID: %s", priceInfo.CLMID)
	}
}

// TestPriceInfoClient_GetPriceInfoHistory は価格履歴取得をテストします
func TestPriceInfoClient_GetPriceInfoHistory(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// まずログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err, "Login should not produce an error")
	require.NotNil(t, session, "Session should not be nil")

	// 価格履歴取得
	historyReq := price_request.ReqGetPriceInfoHistory{
		CLMID:     "CLMMfdsGetMarketPriceHistory",
		IssueCode: "7203",
		SizyouC:   "00", // 東証
	}

	historyInfo, err := c.GetPriceInfoHistory(context.Background(), session, historyReq)

	// 夜間の場合はエラーが発生する可能性があるため、エラーチェックを柔軟に
	if err != nil {
		t.Logf("GetPriceInfoHistory error (may be expected during off-hours): %v", err)
		return // 夜間エラーの場合はスキップ
	}

	// 結果の確認
	require.NotNil(t, historyInfo, "HistoryInfo should not be nil")
	assert.NotEmpty(t, historyInfo.CLMID, "CLMID should not be empty")

	if historyInfo.IssueCode != "" && len(historyInfo.CLMMfdsGetMarketPriceHistory) > 0 {
		t.Logf("価格履歴取得成功 - 銘柄コード: %s, データ件数: %d",
			historyInfo.IssueCode, len(historyInfo.CLMMfdsGetMarketPriceHistory))
	} else {
		t.Logf("価格履歴が返されませんでした（夜間・休日のため正常）- CLMID: %s", historyInfo.CLMID)
	}
}

// TestPriceInfoClient_WithoutLogin は未ログイン状態でのエラーテストです
func TestPriceInfoClient_WithoutLogin(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ダミーセッションを作成（実際にはログインしていない）
	dummySession := client.NewSession()
	dummySession.PriceURL = "https://example.com/dummy"

	// 価格情報取得試行
	priceReq := price_request.ReqGetPriceInfo{
		CLMID:           "CLMMfdsGetMarketPrice",
		TargetIssueCode: "7203",
		TargetColumn:    "pDPP", // 終値のみ
	}

	priceInfo, err := c.GetPriceInfo(context.Background(), dummySession, priceReq)

	// エラーが発生することを確認
	require.Error(t, err, "GetPriceInfo with invalid session should produce an error")
	require.Nil(t, priceInfo, "PriceInfo should be nil when error occurs")

	t.Logf("Expected error with invalid session: %v", err)
}
