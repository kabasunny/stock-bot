// internal/infrastructure/client/tests/order_client_impl_correctorder_test.go
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

		// 例: 現物買い注文　トリガー前の逆指値注文の注文値段は訂正できません
		// 休憩中に指値か成行きを入れて検証する
		// 指値は即座に約定し、テストは失敗する
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "6658",                 // 例: シスメックス
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 買
			Condition:                "0",                    // 指定なし (成行)
			OrderPrice:               "610",                  // 指値　/　成行 (0)
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
		resNewOrder, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)

		time.Sleep(5 * time.Second) // 1秒のタイムラグ

		correctReq := request.ReqCorrectOrder{
			OrderNumber:      resNewOrder.OrderNumber, // 発注した注文の番号
			EigyouDay:        resNewOrder.EigyouDay,   // 発注した注文の営業日
			Condition:        "*",                     // 変更なし
			OrderPrice:       "0",                     // 成行きに変更
			OrderSuryou:      "*",                     // 変更なし
			OrderExpireDay:   "*",                     // 変更なし
			GyakusasiZyouken: "*",                     // 変更なし
			GyakusasiPrice:   "*",                     // 変更なし
			SecondPassword:   c.GetPasswordForTest(),
		}

		res, err := c.CorrectOrder(context.Background(), correctReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		// ... レスポンスの検証 ...
	})

	t.Run("異常系: 存在しない注文番号で訂正が失敗すること", func(t *testing.T) {
		badCorrectReq := request.ReqCorrectOrder{
			OrderNumber:    "invalid_order_number", // 存在しない注文番号
			EigyouDay:      "20250818",             // 営業日を適切に設定 (例)
			SecondPassword: c.GetPasswordForTest(), // 第二パスワードを設定
		}

		_, err := c.CorrectOrder(context.Background(), badCorrectReq)
		assert.Error(t, err)
	})

	// ... 他の異常系テストケース (ログインしていない、第二パスワード間違いなど) ...
	// 異常系: ログインしていない
	// 異常系: 存在しない注文番号
	// 異常系: 第二パスワード間違い
	// 異常系: 増株
}

// TestOrderClientImpl_CorrectOrderWithPost_Cases は CorrectOrderWithPost メソッドのテストケース
func TestOrderClientImpl_CorrectOrderWithPost_Cases(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく) - POST版
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.LoginWithPost(context.Background(), loginReq)
	assert.NoError(t, err)

	// 正常系: 指値訂正 (POST版) - GET版のテストロジックを移植
	t.Run("正常系 (POST): 指値訂正が成功すること", func(t *testing.T) {
		// 事前準備: 訂正可能な注文を発注しておく (NewOrderWithPost を利用)
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "6658",                 // 例: シスメックス
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 買
			Condition:                "0",                    // 指定なし (成行)
			OrderPrice:               "610",                  // 指値　/　成行 (0)
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
		resNewOrder, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)
		// resNewOrderがnilの場合、後続のテストは無意味なので早期リターン
		if !assert.NotNil(t, resNewOrder) {
			t.FailNow()
		}

		time.Sleep(5 * time.Second) // 1秒のタイムラグ

		correctReq := request.ReqCorrectOrder{
			OrderNumber:      resNewOrder.OrderNumber, // 発注した注文の番号
			EigyouDay:        resNewOrder.EigyouDay,   // 発注した注文の営業日
			Condition:        "*",                     // 変更なし
			OrderPrice:       "0",                     // 成行きに変更
			OrderSuryou:      "*",                     // 変更なし
			OrderExpireDay:   "*",                     // 変更なし
			GyakusasiZyouken: "*",                     // 変更なし
			GyakusasiPrice:   "*",                     // 変更なし
			SecondPassword:   c.GetPasswordForTest(),
		}

		res, err := c.CorrectOrder(context.Background(), correctReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
		}
	})

	// 異常系: 存在しない注文番号で訂正が失敗すること (POST版) - GET版のテストロジックを移植
	t.Run("異常系 (POST): 存在しない注文番号で訂正が失敗すること", func(t *testing.T) {
		badCorrectReq := request.ReqCorrectOrder{
			OrderNumber:    "invalid_order_number", // 存在しない注文番号
			EigyouDay:      "20250818",             // 営業日を適切に設定 (例)
			SecondPassword: c.GetPasswordForTest(), // 第二パスワードを設定
		}

		_, err := c.CorrectOrder(context.Background(), badCorrectReq)
		assert.Error(t, err)
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_correctorder_test.go
