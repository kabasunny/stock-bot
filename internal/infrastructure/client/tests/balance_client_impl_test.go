// internal/infrastructure/client/tests/balance_client_impl_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	request_balance "stock-bot/internal/infrastructure/client/dto/balance/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require" // 追加
)

func TestBalanceClientImpl_GetGenbutuKabuList(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	// POST版のログインメソッドを呼び出すように修正
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 現物保有銘柄一覧が取得できること", func(t *testing.T) {
		// GetGenbutuKabuList メソッドを実行
		res, err := c.GetGenbutuKabuList(context.Background(), session) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}
	})
}
func TestBalanceClientImpl_GetShinyouTategyokuList(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 信用建玉一覧が取得できること", func(t *testing.T) {
		// GetShinyouTategyokuList メソッドを実行
		res, err := c.GetShinyouTategyokuList(context.Background(), session) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}
	})
}
func TestBalanceClientImpl_GetZanKaiKanougaku(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 買余力情報が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (銘柄コード、市場は未使用なので空文字列でOK)
		req := request_balance.ReqZanKaiKanougaku{
			IssueCode: "",
			SizyouC:   "",
		}

		// GetZanKaiKanougaku メソッドを実行
		res, err := c.GetZanKaiKanougaku(context.Background(), session, req) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode) // 成功コードの確認
		}
	})
}
func TestBalanceClientImpl_GetZanKaiKanougakuSuii(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 可能額推移情報が取得できること", func(t *testing.T) {
		// リクエストデータを作成
		req := request_balance.ReqZanKaiKanougakuSuii{}

		// GetZanKaiKanougakuSuii メソッドを実行
		res, err := c.GetZanKaiKanougakuSuii(context.Background(), session, req) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode) // 成功コードの確認
		}
	})
}
func TestBalanceClientImpl_GetZanKaiSummary(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 可能額サマリーが取得できること", func(t *testing.T) {
		// GetZanKaiSummary メソッドを実行
		res, err := c.GetZanKaiSummary(context.Background(), session) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}
	})
}
func TestBalanceClientImpl_GetZanKaiGenbutuKaitukeSyousai(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 指定営業日の現物株式買付可能額詳細が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (例: 第4営業日を指定)
		tradingDay := 3 // 第4営業日

		// GetZanKaiGenbutuKaitukeSyousai メソッドを実行
		res, err := c.GetZanKaiGenbutuKaitukeSyousai(context.Background(), session, tradingDay) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
		}
	})
}
func TestBalanceClientImpl_GetZanKaiSinyouSinkidateSyousai(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 指定営業日の信用新規建て可能額詳細が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (例: 第1営業日を指定)
		tradingDay := 0 // 第1営業日

		// GetZanKaiSinyouSinkidateSyousai メソッドを実行
		res, err := c.GetZanKaiSinyouSinkidateSyousai(context.Background(), session, tradingDay) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode) // 成功コードの確認
		}
	})
}
func TestBalanceClientImpl_GetZanRealHosyoukinRitu(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// POST版でログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): リアルタイム保証金率情報が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (パラメータは不要)
		req := request_balance.ReqZanRealHosyoukinRitu{}

		// GetZanRealHosyoukinRitu メソッドを実行
		res, err := c.GetZanRealHosyoukinRitu(context.Background(), session, req) // session引数を追加

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode) // 成功コードの確認
		}
	})
}

func TestBalanceClientImpl_GetZanShinkiKanoIjiritu(t *testing.T) {
	c := client.CreateTestClient(t)

	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 信用新規建て可能維持率情報が取得できること", func(t *testing.T) {
		req := request_balance.ReqZanShinkiKanoIjiritu{}

		res, err := c.GetZanShinkiKanoIjiritu(context.Background(), session, req) // session引数を追加

		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode)
		}
	})
}
func TestBalanceClientImpl_GetZanUriKanousuu(t *testing.T) {
	c := client.CreateTestClient(t)

	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.Equal(t, "0", session.ResultCode) // ログイン成功コードの確認

	t.Run("正常系 (POST): 売却可能数量情報が取得できること", func(t *testing.T) {
		req := request_balance.ReqZanUriKanousuu{
			IssueCode: "8411", // 例としてみずほFGの銘柄コードを指定
		}

		res, err := c.GetZanUriKanousuu(context.Background(), session, req) // session引数を追加

		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode)
		}
	})
}

// go test -v ./internal/infrastructure/client/tests/balance_client_impl_test.go
