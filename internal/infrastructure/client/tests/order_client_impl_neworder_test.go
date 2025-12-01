// internal/infrastructure/client/tests/order_client_impl_neworder_test.go
package tests

import (
	"context"
	"testing"
	"time"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_NewOrder_Cases(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系 (POST): 現物 買い注文 (成行、特定口座) が成功すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "6658",                 // 例: シスメックス
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 買
			Condition:                "0",                    // 指定なし (成行)
			OrderPrice:               "0",                    // 成行 (0)
			OrderSuryou:              "100",                  // 100株
			GenkinShinyouKubun:       "0",                    // 現物
			OrderExpireDay:           "0",                    // 当日限り
			GyakusasiOrderType:       "0",                    // 通常注文
			GyakusasiZyouken:         "0",                    // 指定なし
			GyakusasiPrice:           "*",                    // 指定なし
			TatebiType:               "*",                    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",                    // 指定なし
			SecondPassword:           c.GetPasswordForTest(), // 第二パスワード (発注パスワード)
		}

		res, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
		}
	})

	t.Run("正常系 (POST): 現物売り注文 (成行き、特定口座) が成功すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "6658",                 // 例: シスメックス
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "1",                    // 売
			Condition:                "0",                    // 指定なし
			OrderPrice:               "0",                    // 成行
			OrderSuryou:              "100",                  // 100株
			GenkinShinyouKubun:       "0",                    // 現物
			OrderExpireDay:           "0",                    // 当日限り
			GyakusasiOrderType:       "0",                    // 通常注文
			GyakusasiZyouken:         "0",                    // 指定なし
			GyakusasiPrice:           "*",                    // 指定なし
			TatebiType:               "*",                    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",                    // 指定なし
			SecondPassword:           c.GetPasswordForTest(), // 第二パスワード (発注パスワード)
		}

		time.Sleep(1 * time.Second) // 1秒のタイムラグ

		res, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
		}
	})
	// 他のテストケース (信用返済、現引き/現渡し、逆指値など) は後で追加
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_neworder_test.go
