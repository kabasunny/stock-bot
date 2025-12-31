package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	order_request "stock-bot/internal/infrastructure/client/dto/order/request"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrderClient_GetOrderList は注文一覧取得をテストします
func TestOrderClient_GetOrderList(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// 注文一覧取得
	orderList, err := c.GetOrderList(context.Background(), session, order_request.ReqOrderList{})
	require.NoError(t, err)
	require.NotNil(t, orderList)
	assert.Equal(t, "0", orderList.ResultCode)

	t.Logf("注文一覧取得成功 - 件数: %d", len(orderList.OrderList))

	// 注文がある場合、最初の注文の詳細を確認
	if len(orderList.OrderList) > 0 {
		firstOrder := orderList.OrderList[0]
		assert.NotEmpty(t, firstOrder.OrderOrderNumber, "注文番号が設定されているべき")
		assert.NotEmpty(t, firstOrder.OrderIssueCode, "銘柄コードが設定されているべき")
		t.Logf("最初の注文 - 注文番号: %s, 銘柄: %s, 数量: %s",
			firstOrder.OrderOrderNumber, firstOrder.OrderIssueCode, firstOrder.OrderOrderSuryou)
	}
}

// TestOrderClient_NewOrder_Market は成行注文をテストします
func TestOrderClient_NewOrder_Market(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// 成行買い注文のパラメータ
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

	// 新規注文発行
	newOrderRes, err := c.NewOrder(context.Background(), session, orderParams)
	require.NoError(t, err)
	require.NotNil(t, newOrderRes)
	assert.Equal(t, "0", newOrderRes.ResultCode)

	t.Logf("成行買い注文成功 - 注文番号: %s, 営業日: %s", newOrderRes.OrderNumber, newOrderRes.EigyouDay)

	// 注文番号が発行されていることを確認
	assert.NotEmpty(t, newOrderRes.OrderNumber, "注文番号が発行されているべき")
	assert.NotEmpty(t, newOrderRes.EigyouDay, "営業日が設定されているべき")
}

// TestOrderClient_NewOrder_Limit は指値注文をテストします
func TestOrderClient_NewOrder_Limit(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// 指値買い注文のパラメータ（現在価格より低い価格で指値）
	orderParams := client.NewOrderParams{
		ZyoutoekiKazeiC:          "1",    // 特定口座
		IssueCode:                "7203", // トヨタ
		SizyouC:                  "00",   // 東証
		BaibaiKubun:              "3",    // 買
		Condition:                "0",    // 指定なし
		OrderPrice:               "2000", // 指値価格（低めに設定）
		OrderSuryou:              "100",  // 100株
		GenkinShinyouKubun:       "0",    // 現物
		OrderExpireDay:           "0",    // 当日限り
		GyakusasiOrderType:       "0",    // 通常注文
		GyakusasiZyouken:         "0",    // 指定なし
		GyakusasiPrice:           "*",    // 指定なし
		TatebiType:               "*",    // 指定なし
		TategyokuZyoutoekiKazeiC: "*",    // 指定なし
	}

	// 新規注文発行
	newOrderRes, err := c.NewOrder(context.Background(), session, orderParams)
	require.NoError(t, err)
	require.NotNil(t, newOrderRes)
	assert.Equal(t, "0", newOrderRes.ResultCode)

	t.Logf("指値買い注文成功 - 注文番号: %s, 価格: %s", newOrderRes.OrderNumber, orderParams.OrderPrice)
}

// TestOrderClient_NewOrder_InvalidSymbol は不正銘柄コードをテストします
func TestOrderClient_NewOrder_InvalidSymbol(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// 不正な銘柄コードでの注文パラメータ
	orderParams := client.NewOrderParams{
		ZyoutoekiKazeiC:          "1",     // 特定口座
		IssueCode:                "99999", // 存在しない銘柄コード
		SizyouC:                  "00",    // 東証
		BaibaiKubun:              "3",     // 買
		Condition:                "0",     // 指定なし
		OrderPrice:               "0",     // 成行
		OrderSuryou:              "100",   // 100株
		GenkinShinyouKubun:       "0",     // 現物
		OrderExpireDay:           "0",     // 当日限り
		GyakusasiOrderType:       "0",     // 通常注文
		GyakusasiZyouken:         "0",     // 指定なし
		GyakusasiPrice:           "*",     // 指定なし
		TatebiType:               "*",     // 指定なし
		TategyokuZyoutoekiKazeiC: "*",     // 指定なし
	}

	// 新規注文発行（エラーが期待される）
	newOrderRes, err := c.NewOrder(context.Background(), session, orderParams)

	// エラーまたは失敗レスポンスが返されることを確認
	if err != nil {
		t.Logf("Expected error with invalid symbol: %v", err)
	} else {
		require.NotNil(t, newOrderRes)
		// エラーコードが返されることを確認
		assert.NotEqual(t, "0", newOrderRes.ResultCode, "Invalid symbol should return error code")
		t.Logf("Invalid symbol error - ResultCode: %s, ResultText: %s",
			newOrderRes.ResultCode, newOrderRes.ResultText)
	}
}

// TestOrderClient_CancelOrder は注文キャンセルをテストします
func TestOrderClient_CancelOrder(t *testing.T) {
	// テストクライアントを作成
	c := client.CreateTestClient(t)

	// ログイン
	session, err := c.LoginWithPost(context.Background(), request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	})
	require.NoError(t, err)
	require.NotNil(t, session)

	// まず注文を発行
	orderParams := client.NewOrderParams{
		ZyoutoekiKazeiC:          "1",    // 特定口座
		IssueCode:                "7203", // トヨタ
		SizyouC:                  "00",   // 東証
		BaibaiKubun:              "3",    // 買
		Condition:                "0",    // 指定なし
		OrderPrice:               "1000", // 低い指値（約定しにくい）
		OrderSuryou:              "100",  // 100株
		GenkinShinyouKubun:       "0",    // 現物
		OrderExpireDay:           "0",    // 当日限り
		GyakusasiOrderType:       "0",    // 通常注文
		GyakusasiZyouken:         "0",    // 指定なし
		GyakusasiPrice:           "*",    // 指定なし
		TatebiType:               "*",    // 指定なし
		TategyokuZyoutoekiKazeiC: "*",    // 指定なし
	}

	newOrderRes, err := c.NewOrder(context.Background(), session, orderParams)
	require.NoError(t, err)
	require.NotNil(t, newOrderRes)
	assert.Equal(t, "0", newOrderRes.ResultCode)

	t.Logf("キャンセル用注文発行成功 - 注文番号: %s", newOrderRes.OrderNumber)

	// 少し待機（注文が登録されるまで）
	time.Sleep(1 * time.Second)

	// 注文キャンセル
	cancelParams := client.CancelOrderParams{
		OrderNumber: newOrderRes.OrderNumber,
		EigyouDay:   newOrderRes.EigyouDay,
	}

	cancelRes, err := c.CancelOrder(context.Background(), session, cancelParams)
	require.NoError(t, err)
	require.NotNil(t, cancelRes)
	assert.Equal(t, "0", cancelRes.ResultCode)

	t.Logf("注文キャンセル成功 - 注文番号: %s", cancelParams.OrderNumber)
}
