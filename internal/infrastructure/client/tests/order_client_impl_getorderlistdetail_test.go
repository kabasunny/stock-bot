// internal/infrastructure/client/tests/order_client_impl_getorderlistdetail_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderClientImpl_GetOrderListDetail(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	t.Run("正常系 (POST): 注文詳細取得リクエストが成功すること", func(t *testing.T) {
		// GetOrderListDetail リクエストを作成
		// 注意: このテストはPOSTリクエストの仕組みを検証するものであり、
		// 実際に存在する注文を取得するものではありません。
		// そのため、ResultCodeが"0"で返ってくることは期待しません。
		detailReq := request.ReqOrderListDetail{
			OrderNumber: "1", // 仮の注文番号
		}

		// GetOrderListDetail 実行
		res, err := c.GetOrderListDetail(context.Background(), session, detailReq) // session引数を追加
		assert.NoError(t, err)                                                  // APIからのエラー応答（例: 注文なし）はerrではなく、resに含まれる
		assert.NotNil(t, res)
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_getorderlistdetail_test.go
