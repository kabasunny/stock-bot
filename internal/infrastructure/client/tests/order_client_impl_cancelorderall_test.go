// internal/infrastructure/client/tests/order_client_impl_cancelorderall_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"

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
		orderParams1 := client.NewOrderParams{
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
			GyakusasiZyouken:         "590",  // 逆指値条件 (460円以上)
			GyakusasiPrice:           "520",  // 逆指値値段 (455円)
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}
		_, err := c.NewOrder(context.Background(), orderParams1)
		assert.NoError(t, err)

		orderParams2 := client.NewOrderParams{
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
			GyakusasiZyouken:         "590",  // 逆指値条件 (460円以上)
			GyakusasiPrice:           "530",  // 逆指値値段 (455円)
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}
		_, err = c.NewOrder(context.Background(), orderParams2)
		assert.NoError(t, err)

		// CancelOrderAll リクエストを作成
		cancelAllParams := client.CancelOrderAllParams{}

		// CancelOrderAll 実行
		res, err := c.CancelOrderAll(context.Background(), cancelAllParams)
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

// TestOrderClientImpl_CancelOrderAllWithPost は CancelOrderAllWithPost メソッドのテストケース
func TestOrderClientImpl_CancelOrderAllWithPost(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.LoginWithPost(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系 (POST): 全注文取消が成功すること", func(t *testing.T) {
		// 事前準備: 取り消し可能な注文を複数発注しておく (NewOrderWithPost を利用)
		// 約定しないような低い指値で発注する
		orderParams1 := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",
			IssueCode:                "3632", // 例: グリー
			SizyouC:                  "00",
			BaibaiKubun:              "3",
			Condition:                "0",
			OrderPrice:               "100", // 約定しないような低い指値
			OrderSuryou:              "100",
			GenkinShinyouKubun:       "0",
			OrderExpireDay:           "0",
			GyakusasiOrderType:       "0",
			GyakusasiZyouken:         "0",
			GyakusasiPrice:           "*",
			TatebiType:               "*",
			TategyokuZyoutoekiKazeiC: "*",
		}
		_, err := c.NewOrder(context.Background(), orderParams1)
		assert.NoError(t, err)

		orderParams2 := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",
			IssueCode:                "6658", // 例: シスメックス
			SizyouC:                  "00",
			BaibaiKubun:              "3",
			Condition:                "0",
			OrderPrice:               "100", // 約定しないような低い指値
			OrderSuryou:              "100",
			GenkinShinyouKubun:       "0",
			OrderExpireDay:           "0",
			GyakusasiOrderType:       "0",
			GyakusasiZyouken:         "0",
			GyakusasiPrice:           "*",
			TatebiType:               "*",
			TategyokuZyoutoekiKazeiC: "*",
		}
		_, err = c.NewOrder(context.Background(), orderParams2)
		assert.NoError(t, err)

		// CancelOrderAllWithPost リクエストを作成
		cancelAllParams := client.CancelOrderAllParams{}

		// CancelOrderAllWithPost 実行
		res, err := c.CancelOrderAll(context.Background(), cancelAllParams)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}

		// TODO: 実際に注文が全て取り消されたことを確認する方法を検討する
		//       GetOrderList で確認?
	})

	t.Run("異常系 (POST): ログインしていない状態で全注文取消が失敗すること", func(t *testing.T) {
		// ログアウト
		logoutReq := request_auth.ReqLogout{}
		_, err := c.LogoutWithPost(context.Background(), logoutReq)
		assert.NoError(t, err)

		// CancelOrderAllWithPost リクエストを作成
		cancelAllParams := client.CancelOrderAllParams{}

		// CancelOrderAllWithPost 実行
		_, err = c.CancelOrderAll(context.Background(), cancelAllParams)
		assert.Error(t, err)
		assert.Equal(t, "not logged in", err.Error()) // エラーメッセージを検証

		// ログイン (後処理)
		loginReq := request_auth.ReqLogin{
			UserId:   c.GetUserIDForTest(),
			Password: c.GetPasswordForTest(),
		}
		_, err = c.LoginWithPost(context.Background(), loginReq)
		assert.NoError(t, err)
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_cancelorderall_test.go
