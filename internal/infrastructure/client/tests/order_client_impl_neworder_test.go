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

	// 以降に、t.Run を使って、各ケースのテストを記述していく

	t.Run("正常系: 現物 買い注文 (成行、特定口座) が成功すること", func(t *testing.T) {
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

	t.Run("正常系: 現物売り注文 (成行き、特定口座) が成功すること", func(t *testing.T) {
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

	t.Run("正常系: 現物 買い指値注文 (指値、特定口座) が成功すること", func(t *testing.T) {
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

		res, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
		}
	})

	t.Run("正常系: 現物 売り指値注文 (指値、特定口座) が成功すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "6658",                 // 例: シスメックス
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "1",                    // 売
			Condition:                "0",                    // 指定なし
			OrderPrice:               "700",                  // 指値
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

		time.Sleep(2 * time.Second) // 1秒のタイムラグ

		res, err := c.NewOrder(context.Background(), orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
		}
	})

	// 正常系: 信用 買建→売返済注文 (個別指定、成行、特定口座) で返済をする対象になる
	t.Run("正常系: 信用 新規買い注文 (成行、特定口座) が成功すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "3556",                 // 例: リネットジャパン
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 買
			Condition:                "0",                    // 指定なし (成行)
			OrderPrice:               "0",                    // 成行 (0)
			OrderSuryou:              "200",                  // 100株
			GenkinShinyouKubun:       "2",                    // 信用新規 (制度信用6ヶ月)
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

	t.Run("正常系: 信用 買建→売返済注文 (個別指定、成行、特定口座)", func(t *testing.T) {

		time.Sleep(3 * time.Second) // 買建の待ち

		tategyokuRes, err := c.GetShinyouTategyokuList(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, tategyokuRes)

		// 今回建てた銘柄の建玉番号を探す
		var tategyokuNo string
		for _, tg := range tategyokuRes.SinyouTategyokuList {
			if tg.OrderIssueCode == "3556" {
				tategyokuNo = tg.OrderTategyokuNumber
				break
			}
		}
		assert.NotEmpty(t, tategyokuNo, "建玉番号が取得できませんでした")

		// 3. 信用返済注文（売）
		hensaiOrderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1", // 特定口座
			IssueCode:                "3556",
			SizyouC:                  "00",
			BaibaiKubun:              "1", // 売 (返済)
			Condition:                "0",
			OrderPrice:               "0", // 成行
			OrderSuryou:              "100",
			GenkinShinyouKubun:       "4", // 返済 (制度信用6ヶ月)
			OrderExpireDay:           "0",
			GyakusasiOrderType:       "0",
			GyakusasiZyouken:         "0",
			GyakusasiPrice:           "*",
			TatebiType:               "1", // 個別指定
			TategyokuZyoutoekiKazeiC: "*",
			SecondPassword:           c.GetPasswordForTest(),
			CLMKabuHensaiData: []request.ReqHensaiData{
				{
					TategyokuNumber: tategyokuNo,
					TatebiZyuni:     "1",
					OrderSuryou:     "100",
				},
			},
		}

		resHensai, err := c.NewOrder(context.Background(), hensaiOrderReq)
		assert.NoError(t, err)
		assert.NotNil(t, resHensai)
		if resHensai != nil {
			assert.Equal(t, "0", resHensai.ResultCode)
			assert.NotEmpty(t, resHensai.OrderNumber)
			assert.NotEmpty(t, resHensai.EigyouDay)
		}
	})

	t.Run("正常系: 信用 新規売り注文 (成行、特定口座) が成功すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "3632",                 // 例: グリー
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "1",                    // 売
			Condition:                "0",                    // 指定なし (成行)
			OrderPrice:               "0",                    // 成行 (0)
			OrderSuryou:              "100",                  // 100株
			GenkinShinyouKubun:       "2",                    // 信用新規 (制度信用6ヶ月)
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

	t.Run("正常系: 信用 売りの買返済注文 (建日順、成行、特定口座) が成功すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "3632",                 // 例: アテクト
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 売 (信用返済)
			Condition:                "0",                    // 指定なし
			OrderPrice:               "0",                    // 指値
			OrderSuryou:              "100",                  // 100株 (数量は減らす必要あり)
			GenkinShinyouKubun:       "4",                    // 信用返済 (制度信用6ヶ月)
			OrderExpireDay:           "0",                    // 当日限り
			GyakusasiOrderType:       "0",                    // 通常注文
			GyakusasiZyouken:         "0",                    // 指定なし
			GyakusasiPrice:           "*",                    // 指定なし
			TatebiType:               "2",                    //建日順
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

	t.Run("正常系: 現物 買い注文 (逆指値以上で指値) が成功すること", func(t *testing.T) {
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
			GyakusasiZyouken:         "461",                  // 逆指値条件 (テスト時、リアル株価を入れる)
			GyakusasiPrice:           "461",                  // 逆指値値段 (テスト時、リアル株価を入れる)
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

	t.Run("正常系: 現物 買い注文 (通常+逆指値) が成功すること", func(t *testing.T) {
		orderReq := request.ReqNewOrder{
			ZyoutoekiKazeiC:          "1",                    // 特定口座
			IssueCode:                "3668",                 // 例: コロプラ
			SizyouC:                  "00",                   // 東証
			BaibaiKubun:              "3",                    // 買
			Condition:                "0",                    // 指定なし
			OrderPrice:               "490",                  // 指値 (970円)
			OrderSuryou:              "100",                  // 100株
			GenkinShinyouKubun:       "0",                    // 現物
			OrderExpireDay:           "0",                    // 当日限り
			GyakusasiOrderType:       "2",                    // 通常+逆指値
			GyakusasiZyouken:         "510",                  // 逆指値条件 (974円以上)
			GyakusasiPrice:           "500",                  // 逆指値値段 (972円)
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
	// 他のテストケース (信用返済、現引き/現渡し、逆指値など) は後で追加
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_neworder_test.go
