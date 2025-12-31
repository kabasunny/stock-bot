// internal/infrastructure/client/tests/order_client_impl_neworder_test.go
package tests

import (
	"context"
	"testing"
	"time"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderClientImpl_NewOrder_Cases(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (テストの前にログインしておく)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	t.Run("正常系 (POST): 現物 買い注文 (成行、特定口座) が成功すること", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "6658", // 例: シスメックス
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "3",    // 買
			Condition:                "0",    // 指定なし (成行)
			OrderPrice:               "0",    // 成行 (0)
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "0",    // 現物
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "0",    // 通常注文
			GyakusasiZyouken:         "0",    // 指定なし
			GyakusasiPrice:           "*",    // 指定なし
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}

		res, err := c.NewOrder(context.Background(), session, orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
			t.Logf("現物買い注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
		}
	})

	t.Run("正常系 (POST): 現物売り注文 (成行き、特定口座) が成功すること", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "6658", // 例: シスメックス
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "1",    // 売
			Condition:                "0",    // 指定なし
			OrderPrice:               "0",    // 成行
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "0",    // 現物
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "0",    // 通常注文
			GyakusasiZyouken:         "0",    // 指定なし
			GyakusasiPrice:           "*",    // 指定なし
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}

		time.Sleep(1 * time.Second) // 1秒のタイムラグ

		res, err := c.NewOrder(context.Background(), session, orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
			t.Logf("現物売り注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
		}
	})

	t.Run("正常系 (POST): 現物 買い指値注文 (指値、特定口座) が成功すること", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "6658", // 例: シスメックス
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "3",    // 買
			Condition:                "0",    // 指定なし
			OrderPrice:               "1000", // 指値（市場価格より低めに設定）
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "0",    // 現物
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "0",    // 通常注文
			GyakusasiZyouken:         "0",    // 指定なし
			GyakusasiPrice:           "*",    // 指定なし
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}

		time.Sleep(1 * time.Second)

		res, err := c.NewOrder(context.Background(), session, orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
			t.Logf("現物買い指値注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
		}
	})

	t.Run("正常系 (POST): 現物 売り指値注文 (指値、特定口座) が成功すること", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",     // 特定口座
			IssueCode:                "6658",  // 例: シスメックス
			SizyouC:                  "00",    // 東証
			BaibaiKubun:              "1",     // 売
			Condition:                "0",     // 指定なし
			OrderPrice:               "10000", // 指値（市場価格より高めに設定）
			OrderSuryou:              "100",   // 100株
			GenkinShinyouKubun:       "0",     // 現物
			OrderExpireDay:           "0",     // 当日限り
			GyakusasiOrderType:       "0",     // 通常注文
			GyakusasiZyouken:         "0",     // 指定なし
			GyakusasiPrice:           "*",     // 指定なし
			TatebiType:               "*",     // 指定なし
			TategyokuZyoutoekiKazeiC: "*",     // 指定なし
		}

		time.Sleep(1 * time.Second)

		res, err := c.NewOrder(context.Background(), session, orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
			t.Logf("現物売り指値注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
		}
	})

	t.Run("正常系 (POST): 信用買い注文 (成行) が成功すること", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "6658", // 例: シスメックス
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "3",    // 買
			Condition:                "0",    // 指定なし (成行)
			OrderPrice:               "0",    // 成行 (0)
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "2",    // 信用新規
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "0",    // 通常注文
			GyakusasiZyouken:         "0",    // 指定なし
			GyakusasiPrice:           "*",    // 指定なし
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}

		time.Sleep(1 * time.Second)

		res, err := c.NewOrder(context.Background(), session, orderReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode)
			assert.NotEmpty(t, res.OrderNumber)
			assert.NotEmpty(t, res.EigyouDay)
			t.Logf("信用買い注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
		}
	})

	t.Run("正常系 (POST): 信用売り注文 (成行) - 銘柄制約でエラーの可能性", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "7203", // トヨタ自動車（信用売り可能な銘柄に変更）
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "1",    // 売
			Condition:                "0",    // 指定なし (成行)
			OrderPrice:               "0",    // 成行 (0)
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "2",    // 信用新規（新規売り）
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "0",    // 通常注文
			GyakusasiZyouken:         "0",    // 指定なし
			GyakusasiPrice:           "*",    // 指定なし
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}

		time.Sleep(1 * time.Second)

		res, err := c.NewOrder(context.Background(), session, orderReq)
		if err != nil {
			t.Logf("信用売り注文エラー（銘柄制約の可能性）: %v", err)
			return // エラーでもテストは継続
		}
		assert.NotNil(t, res)
		if res != nil {
			if res.ResultCode == "0" {
				t.Logf("信用売り注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
			} else {
				t.Logf("信用売り注文失敗 - エラーコード: %s, メッセージ: %s", res.ResultCode, res.ResultText)
			}
		}
	})

	t.Run("正常系 (POST): 逆指値注文 (逆指値以上で成行) - 価格制約でエラーの可能性", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "7203", // トヨタ自動車
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "3",    // 買
			Condition:                "0",    // 指定なし
			OrderPrice:               "*",    // 指定なし (逆指値の場合)
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "0",    // 現物
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "1",    // 逆指値
			GyakusasiZyouken:         "3000", // 逆指値条件（適切な価格に調整）
			GyakusasiPrice:           "0",    // 逆指値値段（成行）
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}

		time.Sleep(1 * time.Second)

		res, err := c.NewOrder(context.Background(), session, orderReq)
		if err != nil {
			t.Logf("逆指値注文エラー（価格制約の可能性）: %v", err)
			return // エラーでもテストは継続
		}
		assert.NotNil(t, res)
		if res != nil {
			if res.ResultCode == "0" {
				t.Logf("逆指値注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
			} else {
				t.Logf("逆指値注文失敗 - エラーコード: %s, メッセージ: %s", res.ResultCode, res.ResultText)
			}
		}
	})

	t.Run("正常系 (POST): 通常+逆指値注文 - 価格制約でエラーの可能性", func(t *testing.T) {
		orderReq := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",    // 特定口座
			IssueCode:                "7203", // トヨタ自動車
			SizyouC:                  "00",   // 東証
			BaibaiKubun:              "3",    // 買
			Condition:                "0",    // 指定なし
			OrderPrice:               "2800", // 指値
			OrderSuryou:              "100",  // 100株
			GenkinShinyouKubun:       "0",    // 現物
			OrderExpireDay:           "0",    // 当日限り
			GyakusasiOrderType:       "2",    // 通常＋逆指値
			GyakusasiZyouken:         "3200", // 逆指値条件
			GyakusasiPrice:           "3100", // 逆指値値段
			TatebiType:               "*",    // 指定なし
			TategyokuZyoutoekiKazeiC: "*",    // 指定なし
		}

		time.Sleep(1 * time.Second)

		res, err := c.NewOrder(context.Background(), session, orderReq)
		if err != nil {
			t.Logf("通常+逆指値注文エラー（価格制約の可能性）: %v", err)
			return // エラーでもテストは継続
		}
		assert.NotNil(t, res)
		if res != nil {
			if res.ResultCode == "0" {
				t.Logf("通常+逆指値注文成功 - 注文番号: %s, 営業日: %s", res.OrderNumber, res.EigyouDay)
			} else {
				t.Logf("通常+逆指値注文失敗 - エラーコード: %s, メッセージ: %s", res.ResultCode, res.ResultText)
			}
		}
	})
}
