// internal/infrastructure/client/tests/order_client_impl_cancelorderall_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderClientImpl_CancelOrderAll(t *testing.T) {
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

	t.Run("正常系: 全注文取消が成功すること", func(t *testing.T) {
		// 事前準備: 取り消し可能な注文を複数発注しておく (NewOrder を利用)
		// 約定しないような低い指値で発注する
		orderParams1 := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",
			IssueCode:                "3632", // 例: グリー
			SizyouC:                  "00",
			BaibaiKubun:              "3",
			Condition:                "2",   // 指値
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
		_, err := c.NewOrder(context.Background(), session, orderParams1)
		require.NoError(t, err)

		time.Sleep(1 * time.Second) // 連続注文のための待機

		orderParams2 := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",
			IssueCode:                "6658", // 例: シスメックス
			SizyouC:                  "00",
			BaibaiKubun:              "3",
			Condition:                "2",   // 指値
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
		_, err = c.NewOrder(context.Background(), session, orderParams2)
		require.NoError(t, err)

		time.Sleep(1 * time.Second) // 注文反映のための待機

		// CancelOrderAll リクエストを作成
		cancelAllParams := client.CancelOrderAllParams{}

		// CancelOrderAll 実行
		res, err := c.CancelOrderAll(context.Background(), session, cancelAllParams)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}

		// TODO: 実際に注文が全て取り消されたことを確認する方法を検討する
		//       GetOrderList で確認?
	})

	t.Run("異常系: ログインしていない状態で全注文取消が失敗すること", func(t *testing.T) {
		// 意図的にnilセッションを渡してエラーを確認
		var invalidSession *client.Session = nil

		// CancelOrderAll リクエストを作成
		cancelAllParams := client.CancelOrderAllParams{}

		// CancelOrderAll 実行
		_, err := c.CancelOrderAll(context.Background(), invalidSession, cancelAllParams)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session is nil")
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_cancelorderall_test.go
