// internal/infrastructure/client/tests/balance_client_impl_test.go
package tests

import (
	"context"
	"fmt"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"

	request_balance "stock-bot/internal/infrastructure/client/dto/balance/request"

	"github.com/stretchr/testify/assert"
)

func TestBalanceClientImpl_GetGenbutuKabuList(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 現物保有銘柄一覧が取得できること", func(t *testing.T) {
		// 全銘柄取得の場合は空文字列

		// GetGenbutuKabuList メソッドを実行
		res, err := c.GetGenbutuKabuList(context.Background())

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
			//assert.NotEmpty(t, res.AGenbutuKabuList) // 銘柄リストが空でないことを確認（実際に保有している場合）
			// 他のフィールドも必要に応じてチェック
		}
	})
}

func TestBalanceClientImpl_GetShinyouTategyokuList(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (共通の処理なので、BeforeEach のようなものがあればそちらに移動しても良い)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 信用建玉一覧が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (銘柄コードは指定しない)
		// 全建玉取得の場合は空文字列

		// GetShinyouTategyokuList メソッドを実行
		res, err := c.GetShinyouTategyokuList(context.Background())

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
			//assert.NotEmpty(t, res.SinyouTategyokuList) // 建玉リストが空でないことを確認 (実際に建玉がある場合)
			// 他のフィールド(売建代金合計、買建代金合計など)も必要に応じてチェック
			fmt.Println(res) //追加
		}
	})
}

func TestBalanceClientImpl_GetZanKaiKanougaku(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (共通の処理なので、BeforeEach のようなものがあればそちらに移動しても良い)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 買余力情報が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (銘柄コード、市場は未使用なので空文字列でOK)
		req := request_balance.ReqZanKaiKanougaku{
			IssueCode: "",
			SizyouC:   "",
		}

		// GetZanKaiKanougaku メソッドを実行
		res, err := c.GetZanKaiKanougaku(context.Background(), req)

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode) // 成功コードの確認
			// 他のフィールド(sSummaryGenkabuKaituke など)も必要に応じてチェック
			fmt.Println(res) //追加
		}
	})
}
func TestBalanceClientImpl_GetZanKaiKanougakuSuii(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (共通の処理なので、BeforeEach のようなものがあればそちらに移動しても良い)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 可能額推移情報が取得できること", func(t *testing.T) {
		// リクエストデータを作成
		req := request_balance.ReqZanKaiKanougakuSuii{}

		// GetZanKaiKanougakuSuii メソッドを実行
		res, err := c.GetZanKaiKanougakuSuii(context.Background(), req)

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

	})
}

func TestBalanceClientImpl_GetZanKaiSummary(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (共通の処理なので、BeforeEach のようなものがあればそちらに移動しても良い)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 可能額サマリーが取得できること", func(t *testing.T) {
		// リクエストデータを作成
		// GetZanKaiSummary メソッドを実行
		res, err := c.GetZanKaiSummary(context.Background())

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
			// 他のフィールドも必要に応じてチェック
			fmt.Println(res) //追加
		}
	})
}

func TestBalanceClientImpl_GetZanKaiGenbutuKaitukeSyousai(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (共通の処理)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 指定営業日の現物株式買付可能額詳細が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (例: 第4営業日を指定)
		tradingDay := 3 // 第4営業日

		// GetZanKaiGenbutuKaitukeSyousai メソッドを実行
		res, err := c.GetZanKaiGenbutuKaitukeSyousai(context.Background(), tradingDay)

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.ResultCode) // 成功コードの確認
			assert.NotEmpty(t, res.Hituke)       // 指定日 (YYYYMMDD) が返ってきていること
			// 他のフィールドも必要に応じてチェック (例: sGenbutuKaitukeKanougaku など)
			fmt.Println(res) //追加
		}
	})
}

func TestBalanceClientImpl_GetZanKaiSinyouSinkidateSyousai(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (共通の処理)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 指定営業日の信用新規建て可能額詳細が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (例: 第1営業日を指定)
		tradingDay := 0 // 第1営業日

		// GetZanKaiSinyouSinkidateSyousai メソッドを実行
		res, err := c.GetZanKaiSinyouSinkidateSyousai(context.Background(), tradingDay)

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode) // 成功コードの確認
			assert.NotEmpty(t, res.SHituke)       // 指定日 (YYYYMMDD) が返ってきていること
			// 他のフィールドも必要に応じてチェック (例: sSinyouSinkidateKanougaku など)
			fmt.Println(res) //追加
		}
	})
}

func TestBalanceClientImpl_GetZanRealHosyoukinRitu(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン (共通の処理)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: リアルタイム保証金率情報が取得できること", func(t *testing.T) {
		// リクエストデータを作成 (パラメータは不要)
		req := request_balance.ReqZanRealHosyoukinRitu{}

		// GetZanRealHosyoukinRitu メソッドを実行
		res, err := c.GetZanRealHosyoukinRitu(context.Background(), req)

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode) // 成功コードの確認
			// 他のフィールドも必要に応じてチェック (例: sItakuHosyoukinRitu など)
			fmt.Println(res)
		}
	})
}

func TestBalanceClientImpl_GetZanShinkiKanoIjiritu(t *testing.T) {
	c := client.CreateTestClient(t)

	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 信用新規建て可能維持率情報が取得できること", func(t *testing.T) {
		req := request_balance.ReqZanShinkiKanoIjiritu{}

		res, err := c.GetZanShinkiKanoIjiritu(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode)
			// 他のフィールドも必要に応じてチェック
			fmt.Println(res)
		}
	})
}

func TestBalanceClientImpl_GetZanUriKanousuu(t *testing.T) {
	c := client.CreateTestClient(t)

	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: 売却可能数量情報が取得できること", func(t *testing.T) {
		req := request_balance.ReqZanUriKanousuu{
			IssueCode: "8411", // 例としてみずほFGの銘柄コードを指定
		}

		res, err := c.GetZanUriKanousuu(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)

		if res != nil {
			assert.Equal(t, "0", res.SResultCode)
			// 他のフィールドも必要に応じてチェック
			fmt.Println(res)
		}
	})
}

// go test -v ./internal/infrastructure/client/tests/balance_client_impl_test.go
