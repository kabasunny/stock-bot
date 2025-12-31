// internal/infrastructure/client/tests/order_client_impl_cancelorder_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderClientImpl_CancelOrder(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	t.Run("正常系: 注文取消が成功すること", func(t *testing.T) {
		// 事前準備: 取り消し可能な注文を発注しておく (NewOrder を利用)
		orderParams := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "3632", // 例: グリー
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "3",    // 買
			Condition:                "0",    // 指定なし
			OrderPrice:               "*",    // 指定なし (逆指値の場合)
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "0",    // 現物
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "1",    // 逆指値
			GyakusasiZyouken:         "565",  // 逆指値条件 (460円以上)
			GyakusasiPrice:           "530",  // 逆指値値段 (455円)
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}
		newOrderRes, err := c.NewOrder(context.Background(), session, orderParams)
		assert.NoError(t, err)
		assert.NotNil(t, newOrderRes)
		if newOrderRes == nil {
			t.Fatal("newOrderRes is nil")
		}

		// CancelOrder リクエストを作成
		cancelParams := client.CancelOrderParams{
			OrderNumber: newOrderRes.OrderNumber, // 発注した注文の番号
			EigyouDay:   newOrderRes.EigyouDay,   // 発注した注文の営業日
		}

		// CancelOrder 実行
		res, err := c.CancelOrder(context.Background(), session, cancelParams)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}
	})

	t.Run("異常系: ログインしていない状態で注文取消が失敗すること", func(t *testing.T) {
		// 意図的にnilセッションを渡してエラーを確認
		var invalidSession *client.Session = nil

		// CancelOrder リクエストを作成 (ダミーの値)
		cancelParams := client.CancelOrderParams{
			OrderNumber: "dummy_order_number",
			EigyouDay:   "20230101", // ダミーの値
		}

		// CancelOrder 実行
		_, err = c.CancelOrder(context.Background(), invalidSession, cancelParams)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session is nil") // エラーメッセージを検証
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_cancelorder_test.go
