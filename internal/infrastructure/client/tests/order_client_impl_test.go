// internal/infrastructure/client/tests/order_client_impl_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_NewOrder(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 現物買い注文が成功すること", func(t *testing.T) {
		// 新規注文のリクエストデータを作成 (現物買いの例)
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "8411",                 // 例: みずほFG
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

		// NewOrder メソッドを実行
		res, err := c.NewOrder(context.Background(), orderReq)

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
			assert.NotEmpty(t, res.OrderNumber)  // 注文番号が返ってきていること
			assert.NotEmpty(t, res.EigyouDay)    // 営業日が返ってきていること
		}
	})

	t.Run("異常系: ログインしていない状態で注文が失敗すること", func(t *testing.T) {
		// ログアウト
		logoutReq := request_auth.ReqLogout{}
		_, err := c.Logout(context.Background(), logoutReq)
		assert.NoError(t, err)

		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC: "1",
			IssueCode:       "8411",
			// ... 他のパラメータは省略 ...
			SecondPassword: c.GetPasswordForTest(),
		}

		_, err = c.NewOrder(context.Background(), orderReq)
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

	t.Run("異常系: 不正な銘柄コードで注文が失敗すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC: "1",
			IssueCode:       "invalid_code", // 不正な銘柄コード
			// ... 他のパラメータは省略 ...
			SecondPassword: c.GetPasswordForTest(),
		}

		_, err := c.NewOrder(context.Background(), orderReq)
		assert.Error(t, err)
	})

	t.Run("異常系: 第二パスワードが間違っている場合に注文が失敗すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:    "1",
			IssueCode:          "8411",
			SizyouC:            "00",
			BaibaiKubun:        "3",
			Condition:          "0",
			OrderPrice:         "0",
			OrderSuryou:        "100",
			GenkinShinyouKubun: "0",
			OrderExpireDay:     "0",
			GyakusasiOrderType: "0",
			SecondPassword:     "wrong_password", // 間違ったパスワード
		}

		_, err := c.NewOrder(context.Background(), orderReq)
		assert.Error(t, err)
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_test.go
