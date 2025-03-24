// internal/infrastructure/client/tests/order_client_impl_correctorder_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_CorrectOrder(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	// 以降に、t.Run を使って、各ケースのテストを記述していく

	// 正常系: 指値訂正
	t.Run("正常系: 指値訂正が成功すること", func(t *testing.T) {
		// 事前準備: 訂正可能な注文を発注しておく (NewOrder を利用)

		correctReq := request.ReqCorrectOrder{
			OrderNumber:      "24000028", // 発注した注文の番号 // あらかじめ調べる必要がある
			EigyouDay:        "20250324", // 発注した注文の営業日
			Condition:        "*",        // 変更なし
			OrderPrice:       "600",      // 新しい指値
			OrderSuryou:      "*",        // 変更なし
			OrderExpireDay:   "*",        // 変更なし
			GyakusasiZyouken: "*",        // 変更なし
			GyakusasiPrice:   "*",        // 変更なし
			SecondPassword:   c.GetPasswordForTest(),
		}

		res, err := c.CorrectOrder(context.Background(), correctReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		// ... レスポンスの検証 ...
	})

	t.Run("異常系: 存在しない注文番号で訂正が失敗すること", func(t *testing.T) {
		correctReq := request.ReqCorrectOrder{
			OrderNumber: "invalid_order_number", // 存在しない注文番号
			// ... 他のパラメータ ...
		}

		_, err := c.CorrectOrder(context.Background(), correctReq)
		assert.Error(t, err)
	})

	// ... 他の異常系テストケース (ログインしていない、第二パスワード間違いなど) ...
	// 異常系: ログインしていない
	// 異常系: 存在しない注文番号
	// 異常系: 第二パスワード間違い
	// 異常系: 増株
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_correctorder_test.go
