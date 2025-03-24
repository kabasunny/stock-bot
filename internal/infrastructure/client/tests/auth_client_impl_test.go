// internal/infrastructure/client/tests/auth_client_impl_test.go
package tests

import (
	"context"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"

	"github.com/stretchr/testify/assert"
)

func TestAuthClientImpl_LoginLogout(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t) // client パッケージの CreateTestClient を使用

	t.Run("正常系: 正しいIDとパスワードでログインできること", func(t *testing.T) {
		// Login リクエストを作成
		loginReq := request_auth.ReqLogin{
			UserId:   c.GetUserIDForTest(),   // ヘルパー関数を使用
			Password: c.GetPasswordForTest(), // ヘルパー関数を使用
		}

		// Login 実行
		res, err := c.Login(context.Background(), loginReq)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "0", res.ResultCode) // 成功コードの確認

		// ログイン状態を確認 (loggined フラグ、loginInfo など)
		assert.True(t, c.GetLogginedForTest())                 //test_helper.go
		assert.NotEmpty(t, c.GetLoginInfoForTest().RequestURL) // URLが空でないことを確認 //test_helper.go

		// Logout リクエストを作成
		logoutReq := request_auth.ReqLogout{}

		// Logout 実行
		logoutRes, err := c.Logout(context.Background(), logoutReq)
		assert.NoError(t, err)
		assert.NotNil(t, logoutRes)
		assert.Equal(t, "0", logoutRes.ResultCode)

		// ログアウト状態を確認
		assert.False(t, c.GetLogginedForTest()) //test_helper.go
		assert.Nil(t, c.GetLoginInfoForTest())  // LoginInfo が nil になっていることを確認 //test_helper.go
	})

	t.Run("異常系: 不正なIDとパスワードでログインできないこと", func(t *testing.T) {
		// 異常な値をセット
		originalUserID := c.GetUserIDForTest()
		originalPassword := c.GetPasswordForTest()
		c.SetUserIDForTest("invalid_user")
		c.SetPasswordForTest("invalid_password")

		// Login リクエストを作成 (不正なID/パスワード)
		loginReq := request_auth.ReqLogin{
			UserId:   c.GetUserIDForTest(),   // 設定された不正な値を使用
			Password: c.GetPasswordForTest(), // 設定された不正な値を使用
		}

		// Login 実行
		res, err := c.Login(context.Background(), loginReq)
		assert.Error(t, err) // エラーが発生することを確認
		assert.Nil(t, res)   // レスポンスがnilであること

		// ログイン状態が false であることを確認
		assert.False(t, c.GetLogginedForTest())
		defer func() {
			c.SetUserIDForTest(originalUserID)
			c.SetPasswordForTest(originalPassword)
		}()
	})
	t.Run("異常系: APIサーバーがエラーを返す場合", func(t *testing.T) {
		// オリジナルの baseURL を保持
		originalBaseURL := c.GetBaseURLForTest()

		// baseURL を無効な URL に変更
		c.SetBaseURLForTest("https://invalid.example.com/") // 存在しないURL

		// Login リクエストを作成 (正しいID/パスワードを使用)
		loginReq := request_auth.ReqLogin{
			UserId:   c.GetUserIDForTest(),
			Password: c.GetPasswordForTest(),
		}

		// Login 実行 (エラーが発生するはず)
		res, err := c.Login(context.Background(), loginReq)
		assert.Error(t, err)
		assert.Nil(t, res)

		// ログイン状態が false であることを確認
		assert.False(t, c.GetLogginedForTest())

		// baseURL を元に戻す (defer を使うと、このテストケースが終了する時に確実に実行される)
		defer func() {
			c.SetBaseURLForTest(originalBaseURL)
		}()
	})

}

// go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go
