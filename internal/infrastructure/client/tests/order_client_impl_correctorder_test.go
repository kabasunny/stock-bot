// internal/infrastructure/client/tests/order_client_impl_correctorder_test.go
package tests

import (
	"context"
	"testing"
	"time"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrderClientImpl_CorrectOrder は CorrectOrder メソッドのテストケース
func TestOrderClientImpl_CorrectOrder(t *testing.T) {
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

	// 正常系: 指値訂正
	t.Run("正常系: 指値訂正が成功すること", func(t *testing.T) {
		// 事前準備: 訂正可能な注文を発注しておく (NewOrder を利用)
		// デモ環境では注文が即座に約定しないような、市場価格から離れた指値で注文する必要がある。
		orderParams := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "6658", // 例: シスメックス
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "3",    // 買
			Condition:                "2",    // 指値
			OrderPrice:               "100",  // 市場価格から大きく離れた指値
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "0",    // 現物
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "0",    // 通常注文
			GyakusasiZyouken:         "0",    // 指定なし
			GyakusasiPrice:           "*",    // 指定なし
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}
		resNewOrder, err := c.NewOrder(context.Background(), session, orderParams)
		require.NoError(t, err)
		require.NotNil(t, resNewOrder)
		require.Equal(t, "0", resNewOrder.ResultCode, "事前準備の新規注文に失敗しました")

		time.Sleep(1 * time.Second) // 注文がシステムに反映されるのを待つ

		correctParams := client.CorrectOrderParams{
			OrderNumber:      resNewOrder.OrderNumber, // 発注した注文の番号
			EigyouDay:        resNewOrder.EigyouDay,   // 発注した注文の営業日
			Condition:        "*",                     // 変更なし
			OrderPrice:       "101",                   // "100" -> "101" に訂正
			OrderSuryou:      "*",                     // 変更なし
			OrderExpireDay:   "*",                     // 変更なし
			GyakusasiZyouken: "*",                     // 変更なし
			GyakusasiPrice:   "*",                     // 変更なし
		}

		res, err := c.CorrectOrder(context.Background(), session, correctParams)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
		}
	})

	t.Run("異常系: 存在しない注文番号で訂正が失敗すること", func(t *testing.T) {
		badCorrectParams := client.CorrectOrderParams{
			OrderNumber: "invalid_order_number", // 存在しない注文番号
			EigyouDay:   "20251220",             // 営業日を適切に設定 (例)
		}

		res, err := c.CorrectOrder(context.Background(), session, badCorrectParams)
		assert.NoError(t, err) // API通信自体は成功するはず
		assert.NotNil(t, res)
		if res != nil {
			assert.NotEqual(t, "0", res.ResultCode) // エラーコードが返ることを期待
		}
	})

	t.Run("異常系: ログインしていない状態で訂正が失敗すること", func(t *testing.T) {
		// 意図的にnilセッションを渡してエラーを確認
		var invalidSession *client.Session = nil

		badCorrectParams := client.CorrectOrderParams{
			OrderNumber: "dummy_order_number",
			EigyouDay:   "20251220",
		}

		_, err := c.CorrectOrder(context.Background(), invalidSession, badCorrectParams)
		assert.Error(t, err) // クライアント内部でエラーになることを期待
		assert.Contains(t, err.Error(), "session is nil")
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_correctorder_test.go
