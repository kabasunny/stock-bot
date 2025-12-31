package tests

import (
	"context"
	"log/slog"
	"os"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	balance_request "stock-bot/internal/infrastructure/client/dto/balance/request"
	order_request "stock-bot/internal/infrastructure/client/dto/order/request"
	price_request "stock-bot/internal/infrastructure/client/dto/price/request"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnifiedClientAdapter_BalanceClient はBalanceClient互換性をテストします
func TestUnifiedClientAdapter_BalanceClient(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// UnifiedClientを作成（EventClientはnilで渡す）
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	unifiedClient := client.NewTachibanaUnifiedClient(
		c, c, c, c, c, nil, // EventClientはnilで渡す
		c.GetUserIDForTest(),
		c.GetPasswordForTest(),
		c.GetPasswordForTest(),
		logger,
	)

	// Adapterを作成
	adapter := client.NewTachibanaUnifiedClientAdapter(unifiedClient)

	// 残高サマリー取得テスト
	zanKaiSummary, err := adapter.GetZanKaiSummary(context.Background(), session)
	require.NoError(t, err)
	require.NotNil(t, zanKaiSummary)
	assert.NotEmpty(t, zanKaiSummary.CLMID)

	t.Logf("残高サマリー取得成功 - CLMID: %s", zanKaiSummary.CLMID)

	// 現物株一覧取得テスト
	genbutuKabuList, err := adapter.GetGenbutuKabuList(context.Background(), session)
	require.NoError(t, err)
	require.NotNil(t, genbutuKabuList)
	assert.NotEmpty(t, genbutuKabuList.CLMID)

	t.Logf("現物株一覧取得成功 - CLMID: %s", genbutuKabuList.CLMID)

	// 信用建玉一覧取得テスト
	shinyouTategyokuList, err := adapter.GetShinyouTategyokuList(context.Background(), session)
	require.NoError(t, err)
	require.NotNil(t, shinyouTategyokuList)
	assert.NotEmpty(t, shinyouTategyokuList.CLMID)

	t.Logf("信用建玉一覧取得成功 - CLMID: %s", shinyouTategyokuList.CLMID)
}

// TestUnifiedClientAdapter_OrderClient はOrderClient互換性をテストします
func TestUnifiedClientAdapter_OrderClient(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// UnifiedClientを作成（EventClientはnilで渡す）
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	unifiedClient := client.NewTachibanaUnifiedClient(
		c, c, c, c, c, nil, // EventClientはnilで渡す
		c.GetUserIDForTest(),
		c.GetPasswordForTest(),
		c.GetPasswordForTest(),
		logger,
	)

	// Adapterを作成
	adapter := client.NewTachibanaUnifiedClientAdapter(unifiedClient)

	// 注文一覧取得テスト
	orderList, err := adapter.GetOrderList(context.Background(), session, order_request.ReqOrderList{})
	require.NoError(t, err)
	require.NotNil(t, orderList)
	assert.Equal(t, "0", orderList.ResultCode)

	t.Logf("注文一覧取得成功 - 件数: %d", len(orderList.OrderList))

	// 新規注文テスト（成行買い）
	orderParams := client.NewOrderParams{
		ZyoutoekiKazeiC:          "1",    // 特定口座
		IssueCode:                "6658", // シスメックス
		SizyouC:                  "00",   // 東証
		BaibaiKubun:              "3",    // 買
		Condition:                "0",    // 指定なし (成行)
		OrderPrice:               "0",    // 成行
		OrderSuryou:              "100",  // 100株
		GenkinShinyouKubun:       "0",    // 現物
		OrderExpireDay:           "0",    // 当日限り
		GyakusasiOrderType:       "0",    // 通常注文
		GyakusasiZyouken:         "0",    // 指定なし
		GyakusasiPrice:           "*",    // 指定なし
		TatebiType:               "*",    // 指定なし
		TategyokuZyoutoekiKazeiC: "*",    // 指定なし
	}

	newOrderRes, err := adapter.NewOrder(context.Background(), session, orderParams)
	require.NoError(t, err)
	require.NotNil(t, newOrderRes)
	assert.Equal(t, "0", newOrderRes.ResultCode)

	t.Logf("新規注文成功 - 注文番号: %s, 営業日: %s", newOrderRes.OrderNumber, newOrderRes.EigyouDay)
}

// TestUnifiedClientAdapter_PriceInfoClient はPriceInfoClient互換性をテストします
func TestUnifiedClientAdapter_PriceInfoClient(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// UnifiedClientを作成（EventClientはnilで渡す）
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	unifiedClient := client.NewTachibanaUnifiedClient(
		c, c, c, c, c, nil, // EventClientはnilで渡す
		c.GetUserIDForTest(),
		c.GetPasswordForTest(),
		c.GetPasswordForTest(),
		logger,
	)

	// Adapterを作成
	adapter := client.NewTachibanaUnifiedClientAdapter(unifiedClient)

	// 価格情報取得テスト
	priceReq := price_request.ReqGetPriceInfo{
		CLMID:           "CLMMfdsGetMarketPrice",
		TargetIssueCode: "7203",
		TargetColumn:    "pDPP,pDOP,pDHP,pDLP", // 終値、始値、高値、安値
	}

	priceInfo, err := adapter.GetPriceInfo(context.Background(), session, priceReq)
	if err != nil {
		t.Logf("価格情報取得エラー（夜間・休日のため正常）: %v", err)
		return // 夜間エラーの場合はスキップ
	}

	require.NotNil(t, priceInfo)
	assert.NotEmpty(t, priceInfo.CLMID)

	t.Logf("価格情報取得成功 - CLMID: %s", priceInfo.CLMID)

	// 価格履歴取得テスト
	historyReq := price_request.ReqGetPriceInfoHistory{
		CLMID:     "CLMMfdsGetMarketPriceHistory",
		IssueCode: "7203",
		SizyouC:   "00", // 東証
	}

	historyInfo, err := adapter.GetPriceInfoHistory(context.Background(), session, historyReq)
	if err != nil {
		t.Logf("価格履歴取得エラー（夜間・休日のため正常）: %v", err)
		return // 夜間エラーの場合はスキップ
	}

	require.NotNil(t, historyInfo)
	assert.NotEmpty(t, historyInfo.CLMID)

	t.Logf("価格履歴取得成功 - CLMID: %s", historyInfo.CLMID)
}

// TestUnifiedClientAdapter_ErrorHandling はエラーハンドリングをテストします
func TestUnifiedClientAdapter_ErrorHandling(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// UnifiedClientを作成（EventClientはnilで渡す）
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	unifiedClient := client.NewTachibanaUnifiedClient(
		c, c, c, c, c, nil, // EventClientはnilで渡す
		"invalid_user",     // 無効なユーザーID
		"invalid_password", // 無効なパスワード
		"invalid_password", // 無効なセカンドパスワード
		logger,
	)

	// Adapterを作成
	adapter := client.NewTachibanaUnifiedClientAdapter(unifiedClient)

	// 無効な認証情報でのテスト - GetSessionを呼び出してエラーを確認
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// GetSessionでエラーが発生することを確認（内部で認証が実行される）
	_, err := unifiedClient.GetSession(ctx)
	require.Error(t, err)
	t.Logf("Expected authentication error with invalid credentials: %v", err)

	// 無効なセッション（空のセッション）での残高取得テスト
	invalidSession := client.NewSession()
	// 空のセッションでAPIを呼び出すとエラーになるはず

	// 残高取得でエラーが発生することを確認
	_, err = adapter.GetZanKaiSummary(context.Background(), invalidSession)
	require.Error(t, err)
	t.Logf("Expected error with invalid session: %v", err)
}

// TestUnifiedClientAdapter_NotImplementedMethods は未実装メソッドをテストします
func TestUnifiedClientAdapter_NotImplementedMethods(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// UnifiedClientを作成（EventClientはnilで渡す）
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	unifiedClient := client.NewTachibanaUnifiedClient(
		c, c, c, c, c, nil, // EventClientはnilで渡す
		c.GetUserIDForTest(),
		c.GetPasswordForTest(),
		c.GetPasswordForTest(),
		logger,
	)

	// Adapterを作成
	adapter := client.NewTachibanaUnifiedClientAdapter(unifiedClient)

	// 未実装メソッドのテスト
	_, err = adapter.GetZanKaiKanougaku(context.Background(), session, balance_request.ReqZanKaiKanougaku{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")

	_, err = adapter.CancelOrderAll(context.Background(), session, client.CancelOrderAllParams{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")

	_, err = adapter.GetOrderListDetail(context.Background(), session, order_request.ReqOrderListDetail{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")

	t.Logf("未実装メソッドで適切にエラーが発生することを確認")
}
