// internal/infrastructure/client/tests/order_client_impl_getorderlist_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/order/request"

	"github.com/stretchr/testify/assert"
)

func TestOrderClientImpl_GetOrderList(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系 (POST): 注文一覧取得が成功すること", func(t *testing.T) {
		// GetOrderList リクエストを作成 (すべての注文を取得)
		orderListReq := request.ReqOrderList{}

		// GetOrderList 実行
		res, err := c.GetOrderList(context.Background(), orderListReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}
	})
}

// go test -v ./internal/infrastructure/client/tests/order_client_impl_getorderlist_test.go
