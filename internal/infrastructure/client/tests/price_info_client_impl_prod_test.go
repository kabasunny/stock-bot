// internal/infrastructure/client/tests/price_info_client_impl_prod_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPriceInfoClientImpl_Production は、本番環境に対してテストを実行します。
//
// !!! 注意 !!!
// このテストを実行する直前に、必ず【手動で電話認証】を済ませてください。
// 電話認証が行われていない場合、ログインに失敗し、テストは実行されません。
// また、.envファイルが本番環境用に設定されている必要があります。
func TestPriceInfoClientImpl_Production(t *testing.T) {
	// .env から本番用の設定を読み込むために CreateTestClient を使用
	c := client.CreateTestClient(t)

	// --- 1. 本番環境へのログイン ---
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)

	// ログイン失敗時はテストを中断
	require.NoError(t, err, "本番環境へのログインに失敗しました。電話認証が完了しているか、.envの本番用設定（ID/PW/URL）が正しいか確認してください。")
	require.NotNil(t, session)
	require.Equal(t, "0", session.ResultCode, "ログインAPIからエラーが返されました。ResultCode: %s", session.ResultCode)

	t.Log("本番環境へのログインに成功しました。株価照会テストを開始します。")

	// --- 2. テストの実行 ---
	t.Run("本番系 (POST): トヨタの株価情報取得", func(t *testing.T) {
		req := request.ReqGetPriceInfo{
			TargetIssueCode: "7203", // トヨタ自動車
		}
		res, err := c.GetPriceInfo(context.Background(), session, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		// 本番なのでデータが返ってくることを期待
		assert.Greater(t, len(res.CLMMfdsMarketPrice), 0, "本番環境にもかかわらず、株価情報が返されませんでした。")
		if len(res.CLMMfdsMarketPrice) > 0 {
			assert.Equal(t, "7203", res.CLMMfdsMarketPrice[0].IssueCode)
			t.Logf("取得成功 (株価情報): %+v", res.CLMMfdsMarketPrice[0])
		}
	})

	t.Run("本番系 (POST): トヨタの株価履歴情報取得", func(t *testing.T) {
		req := request.ReqGetPriceInfoHistory{
			IssueCode: "7203", // トヨタ自動車
		}
		res, err := c.GetPriceInfoHistory(context.Background(), session, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		// 本番なのでデータが返ってくることを期待
		assert.Equal(t, "7203", res.IssueCode)
		assert.Greater(t, len(res.CLMMfdsGetMarketPriceHistory), 0, "本番環境にもかかわらず、株価履歴が返されませんでした。")
		if len(res.CLMMfdsGetMarketPriceHistory) > 0 {
			t.Logf("取得成功 (株価履歴): %d件のデータを取得", len(res.CLMMfdsGetMarketPriceHistory))
		}
	})
}

// go test -v ./internal/infrastructure/client/tests/price_info_client_impl_prod_test.go -run TestPriceInfoClientImpl_Production

// TestPriceInfo_Sequence_LoginWaitGetPrice は、無通信タイムアウトを確認するテストです。
// ログイン後、30分待機してから株価照会APIを呼び出します。
// このテストは完了まで30分以上かかります。
func TestPriceInfo_Sequence_LoginWaitGetPrice(t *testing.T) {
	t.Log("【シーケンステスト開始】ログイン → 30分待機 → 株価照会")

	// 1. ログイン
	c := client.CreateTestClient(t)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err, "シーケンステスト中のログインに失敗しました")
	require.NotNil(t, session)
	require.Equal(t, "0", session.ResultCode, "ログインAPIからエラーが返されました")
	t.Log("ログイン成功。")

	// 2. 30分間待機
	const waitMinutes = 30
	t.Logf("%d分間待機します...", waitMinutes)
	time.Sleep(waitMinutes * time.Minute)
	t.Log("待機完了。")

	// 3. 株価照会
	t.Log("セッションが有効か確認するため、株価照会APIを呼び出します。")
	req := request.ReqGetPriceInfo{
		TargetIssueCode: "7203", // トヨタ自動車
	}
	res, err := c.GetPriceInfo(context.Background(), session, req)
	if err != nil {
		t.Logf("株価照会APIの呼び出しでエラーが発生しました: %v", err)
	}
	if res != nil {
		t.Logf("株価照会APIの応答: %+v", res)
		if len(res.CLMMfdsMarketPrice) > 0 {
			t.Logf("取得成功 (株価情報): %+v", res.CLMMfdsMarketPrice[0])
		} else {
			t.Log("株価情報は返されませんでした。")
		}
	}
}

// go test -v ./internal/infrastructure/client/tests/price_info_client_impl_prod_test.go -run TestPriceInfo_Sequence_LoginWaitGetPrice -timeout 35m
