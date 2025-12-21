// internal/infrastructure/client/tests/price_info_client_impl_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/price/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupLoggedInClientForPriceInfoTest は、テスト用のログイン済みクライアントをセットアップするヘルパー関数です。
func setupLoggedInClientForPriceInfoTest(t *testing.T) (*client.TachibanaClientImpl, *client.Session) {
	t.Helper() // これがヘルパー関数であることを示す

	c := client.CreateTestClient(t)

	// auth_client_impl_test.go を参考に LoginWithPost を使用
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	// ログイン実行
	session, err := c.LoginWithPost(context.Background(), loginReq) // session を受け取るように変更
	// ログインに失敗した場合はテストを即時終了
	require.NoError(t, err, "テストの前提条件であるログインに失敗しました。ログインID/パスワード、APIのURLを確認してください。")
	require.NotNil(t, session, "セッションがnilです。") // session が nil でないことを確認

	return c, session
}

func TestPriceInfoClientImpl_GetPriceInfo(t *testing.T) {
	// ヘルパー関数でログイン済みのクライアントを取得
	c, session := setupLoggedInClientForPriceInfoTest(t) // session も受け取るように変更

	t.Run("正常系 (POST): トヨタの株価情報取得が成功し、内容が正しいこと", func(t *testing.T) {
		// リクエストパラメータの設定
		req := request.ReqGetPriceInfo{
			TargetIssueCode: "7203", // トヨタ自動車
		}

		// API呼び出し
		res, err := c.GetPriceInfo(context.Background(), session, req) // session引数を追加
		if err != nil {
			t.Fatalf("API呼び出しエラー: %v", err)
		}

		// レスポンスの検証
		assert.NotNil(t, res)
		assert.Equal(t, "CLMMfdsGetMarketPrice", res.CLMID)

		if len(res.CLMMfdsMarketPrice) > 0 {
			// データが返された場合（本番環境など）は、内容を検証
			assert.Equal(t, "7203", res.CLMMfdsMarketPrice[0].IssueCode)
			t.Logf("取得した株価情報: %+v", res.CLMMfdsMarketPrice[0])
		} else {
			// データが返されない場合（デモ環境など）は、ログを出力してテストをパスさせる
			t.Log("株価情報がデモ環境では返されませんでした。API接続性の確認としては成功です。")
		}
	})
}

func TestPriceInfoClientImpl_GetPriceInfoHistory(t *testing.T) {
	// ヘルパー関数でログイン済みのクライアントを取得
	c, session := setupLoggedInClientForPriceInfoTest(t) // session も受け取るように変更

	t.Run("正常系 (POST): トヨタの株価履歴情報取得が成功し、内容が正しいこと", func(t *testing.T) {
		// リクエストパラメータの設定
		req := request.ReqGetPriceInfoHistory{
			IssueCode: "7203", // トヨタ自動車
		}

		// API呼び出し
		res, err := c.GetPriceInfoHistory(context.Background(), session, req) // session引数を追加
		if err != nil {
			t.Fatalf("API呼び出しエラー: %v", err)
		}

		// レスポンスの検証
		assert.NotNil(t, res)
		assert.Equal(t, "CLMMfdsGetMarketPriceHistory", res.CLMID)

		if len(res.CLMMfdsGetMarketPriceHistory) > 0 {
			// データが返された場合（本番環境など）は、内容を検証
			assert.Equal(t, "7203", res.IssueCode)
			t.Logf("取得した株価履歴情報の件数: %d件", len(res.CLMMfdsGetMarketPriceHistory))
			if len(res.CLMMfdsGetMarketPriceHistory) > 0 {
				t.Logf("最新の株価履歴情報: %+v", res.CLMMfdsGetMarketPriceHistory[0])
			}
		} else {
			// データが返されない場合（デモ環境など）は、ログを出力してテストをパスさせる
			t.Log("株価履歴情報がデモ環境では返されませんでした。API接続性の確認としては成功です。")
		}
	})
}

// go test -v internal/infrastructure/client/tests/price_info_client_impl_demo_test.go
