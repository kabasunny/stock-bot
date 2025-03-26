// internal/infrastructure/client/tests/order_client_impl_getorderlistdetail_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_GetOrderListDetail(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 注文詳細取得が成功すること", func(t *testing.T) {
		// 事前準備: 注文を発注しておく (NewOrder を利用)
		orderReq := request.ReqNewOrder{
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
			GyakusasiPrice:           "550",                  // 逆指値値段 (455円)
			TatebiType:               "*",                    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",                    // 指定なし
			SecondPassword:           c.GetPasswordForTest(), // 第二パスワード (発注パスワード)
		}
		newOrderRes, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, newOrderRes)

		// GetOrderListDetail リクエストを作成
		detailReq := request.ReqOrderListDetail{
			OrderNumber: newOrderRes.OrderNumber, // 発注した注文の番号
			EigyouDay:   newOrderRes.EigyouDay,   // 発注した注文の営業日
		}

		// GetOrderListDetail 実行
		res, err := c.GetOrderListDetail(context.Background(), detailReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)                      // 成功コードの確認
			assert.Equal(t, newOrderRes.OrderNumber, res.OrderNumber) // 注文番号が一致することを確認
			assert.Equal(t, newOrderRes.EigyouDay, res.EigyouDay)     // 営業日が一致することを確認
			assert.NotEmpty(t, res.IssueCode)                         // 銘柄コードが空でないことを確認 (追加)
			// 他のフィールド (約定情報など) も必要に応じて検証
		}
	})

	// 以降に、異常系のテストケース (ログインしていない、存在しない注文番号) を追加
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_getorderlistdetail_test.go
