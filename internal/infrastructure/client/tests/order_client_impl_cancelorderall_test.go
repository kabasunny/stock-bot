// internal/infrastructure/client/tests/order_client_impl_cancelorderall_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_CancelOrderAll(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 全注文取消が成功すること", func(t *testing.T) {
		// 事前準備: 取り消し可能な注文を複数発注しておく (NewOrder を利用)
		// ※ テスト実行ごとに一意の注文番号が採番されるため、
		//    ここではダミーの値を使用し、後で実際の発注結果で置き換える

		// 例: 現物買い注文を2件発注
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
			GyakusasiZyouken:         "550",                  // 逆指値条件 (460円以上)
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
			GyakusasiZyouken:         "590",                  // 逆指値条件 (460円以上)
			GyakusasiPrice:           "455",                  // 逆指値値段 (455円)
			TatebiType:               "*",                    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",                    // 指定なし
			SecondPassword:           c.GetPasswordForTest(), // 第二パスワード (発注パスワード)
		}
		_, err = c.NewOrder(context.Background(), orderReq2)
		assert.NoError(t, err)

		// CancelOrderAll リクエストを作成
		cancelAllReq := request.ReqCancelOrderAll{
			SecondPassword: c.GetPasswordForTest(),
		}

		// CancelOrderAll 実行
		res, err := c.CancelOrderAll(context.Background(), cancelAllReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}

		// TODO: 実際に注文が全て取り消されたことを確認する方法を検討する
		//       GetOrderList で確認?
	})
	// 異常系: ログインしていない
	// 異常系: 第二パスワード間違い
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_cancelorderall_test.go
