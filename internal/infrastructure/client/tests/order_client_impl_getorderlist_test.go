// internal/infrastructure/client/tests/order_client_impl_getorderlist_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_GetOrderList(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 注文一覧取得が成功すること", func(t *testing.T) {
		// 事前準備: 注文をいくつか発注しておく (NewOrder を利用)
		// 例: 現物買い注文を2件発注 (銘柄コード、数量などは適当に変更してください)
		orderReq1 := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "3632",                 // 例: グリー
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 買
			Condition:                "0",                    // 指定なし
			OrderPrice:               "*",                    // 指定なし (逆指値の場合)
			OrderSuryou:              "100",                  // 100株
			GenkinShinyouKubun:       "0",                    // 現物
			OrderExpireDay:           "0",                    // 当日限り
			GyakusasiOrderType:       "1",                    // 逆指値
			GyakusasiZyouken:         "570",                  // 逆指値条件 (460円以上)
			GyakusasiPrice:           "455",                  // 逆指値値段 (455円)
			TatebiType:               "*",                    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",                    // 指定なし
			SecondPassword:           c.GetPasswordForTest(), // 第二パスワード (発注パスワード)
		}
		_, err := c.NewOrder(context.Background(), orderReq1)
		assert.NoError(t, err)

		orderReq2 := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "3632",                 // 例: グリー
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 買
			Condition:                "0",                    // 指定なし
			OrderPrice:               "*",                    // 指定なし (逆指値の場合)
			OrderSuryou:              "100",                  // 100株
			GenkinShinyouKubun:       "0",                    // 現物
			OrderExpireDay:           "0",                    // 当日限り
			GyakusasiOrderType:       "1",                    // 逆指値
			GyakusasiZyouken:         "580",                  // 逆指値条件 (460円以上)
			GyakusasiPrice:           "455",                  // 逆指値値段 (455円)
			TatebiType:               "*",                    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",                    // 指定なし
			SecondPassword:           c.GetPasswordForTest(), // 第二パスワード (発注パスワード)
		}
		_, err = c.NewOrder(context.Background(), orderReq2)
		assert.NoError(t, err)

		// GetOrderList リクエストを作成 (すべての注文を取得)
		orderListReq := request.ReqOrderList{}

		// GetOrderList 実行
		res, err := c.GetOrderList(context.Background(), orderListReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
			assert.NotEmpty(t, res.OrderList)    // 注文リストが空でないことを確認
			// 必要であれば、取得した注文一覧の内容を検証
			// (例: 注文数、各注文の情報が正しいか)
		}
	})
	// 異常系: ログインしていない
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_getorderlist_test.go

func TestOrderClientImpl_GetOrderListWithPost(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.LoginWithPost(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系 (POST): 注文一覧取得が成功すること", func(t *testing.T) {
		// GetOrderList リクエストを作成 (すべての注文を取得)
		orderListReq := request.ReqOrderList{}

		// GetOrderListWithPost 実行
		res, err := c.GetOrderListWithPost(context.Background(), orderListReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}
	})
}
