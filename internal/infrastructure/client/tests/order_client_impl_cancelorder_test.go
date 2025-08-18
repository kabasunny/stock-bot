// internal/infrastructure/client/tests/order_client_impl_cancelorder_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_CancelOrder(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 注文取消が成功すること", func(t *testing.T) {
		// 事前準備: 取り消し可能な注文を発注しておく (NewOrder を利用)
		// ※ テスト実行ごとに一意の注文番号が採番されるため、
		//    ここではダミーの値を使用し、後で実際の発注結果で置き換える
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
			GyakusasiZyouken:         "565",                  // 逆指値条件 (460円以上)
			GyakusasiPrice:           "530",                  // 逆指値値段 (455円)
			TatebiType:               "*",                    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",                    // 指定なし
			SecondPassword:           c.GetPasswordForTest(), // 第二パスワード (発注パスワード)
		}
		newOrderRes, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, newOrderRes)

		// CancelOrder リクエストを作成
		cancelReq := request.ReqCancelOrder{
			OrderNumber:    newOrderRes.OrderNumber, // 発注した注文の番号
			EigyouDay:      newOrderRes.EigyouDay,   // 発注した注文の営業日
			SecondPassword: c.GetPasswordForTest(),  // 第二パスワード
		}

		// CancelOrder 実行
		res, err := c.CancelOrder(context.Background(), cancelReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
			// 他にも、必要に応じてレスポンスの内容を検証 (ResultText など)
		}
	})

	t.Run("異常系: ログインしていない状態で注文取消が失敗すること", func(t *testing.T) {
		// ログアウト
		logoutReq := request_auth.ReqLogout{}
		_, err := c.Logout(context.Background(), logoutReq)
		assert.NoError(t, err)

		// CancelOrder リクエストを作成 (ダミーの値)
		cancelReq := request.ReqCancelOrder{
			OrderNumber:    "dummy_order_number",
			EigyouDay:      "20230101", // ダミーの値
			SecondPassword: c.GetPasswordForTest(),
		}

		// CancelOrder 実行
		_, err = c.CancelOrder(context.Background(), cancelReq)
		assert.Error(t, err)
		assert.Equal(t, "not logged in", err.Error()) // エラーメッセージを検証

		// ログイン (後処理)
		loginReq := request_auth.ReqLogin{
			UserId:   c.GetUserIDForTest(),
			Password: c.GetPasswordForTest(),
		}
		_, err = c.Login(context.Background(), loginReq)
		assert.NoError(t, err)
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_cancelorder_test.go
