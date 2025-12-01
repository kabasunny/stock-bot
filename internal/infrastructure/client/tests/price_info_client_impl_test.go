// internal/infrastructure/client/tests/price_info_client_impl_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriceInfoClientImpl_GetPriceInfo(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系 (POST): 株価情報取得が成功すること", func(t *testing.T) {
		// リクエストパラメータの設定
		req := request.ReqGetPriceInfo{
			TargetIssueCode: "6501", // 日立
		}

		// API呼び出し
		res, err := c.GetPriceInfo(context.Background(), req)
		if err != nil {
			t.Fatalf("API呼び出しエラー: %v", err)
		}

		// レスポンスの検証
		assert.NotNil(t, res)
		assert.Equal(t, "CLMMfdsGetMarketPrice", res.CLMID)

		if len(res.CLMMfdsMarketPrice) > 0 {
			assert.Equal(t, "6501", res.CLMMfdsMarketPrice[0].IssueCode)
		} else {
			t.Log("株価情報が存在しません")
		}
	})
}

func TestPriceInfoClientImpl_GetPriceInfoHistory(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系 (POST): 株価履歴情報取得が成功すること", func(t *testing.T) {
		// リクエストパラメータの設定
		req := request.ReqGetPriceInfoHistory{
			IssueCode: "6501", // 日立
		}

		// API呼び出し
		res, err := c.GetPriceInfoHistory(context.Background(), req)
		if err != nil {
			t.Fatalf("API呼び出しエラー: %v", err)
		}

		// レスポンスの検証
		assert.NotNil(t, res)
		assert.Equal(t, "CLMMfdsGetMarketPriceHistory", res.CLMID)

		if len(res.CLMMfdsGetMarketPriceHistory) > 0 {
			assert.Equal(t, "6501", res.IssueCode) // レスポンスのルートレベルにもIssueCodeがある
		} else {
			t.Log("株価履歴情報が存在しません")
		}
	})
}

// go test -v internal/infrastructure/client/tests/price_info_client_impl_test.go
