// internal/infrastructure/client/tests/auth_client_impl_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"testing"
	"time" // timeパッケージをインポート

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestAuthClientImpl_LoginLogout(t *testing.T) {
// 	c := client.CreateTestClient(t)

// 	tests := []struct {
// 		name           string
// 		setup          func()
// 		loginReq       request_auth.ReqLogin
// 		expectLoginErr bool
// 		expectLoginRes bool
// 		expectLoggedIn bool
// 		expectLogout   bool
// 	}{
// 		{
// 			name: "正常系: 正しいIDとパスワードでログインできること",
// 			setup: func() {
// 				// 正常な認証情報を設定
// 				c.SetUserIDForTest(c.GetUserIDForTest())
// 				c.SetPasswordForTest(c.GetPasswordForTest())
// 				c.SetBaseURLForTest(c.GetBaseURLForTest())
// 			},
// 			loginReq: request_auth.ReqLogin{
// 				UserId:   c.GetUserIDForTest(),
// 				Password: c.GetPasswordForTest(),
// 			},
// 			expectLoginErr: false,
// 			expectLoginRes: true,
// 			expectLoggedIn: true,
// 			expectLogout:   true,
// 		},
// 		{
// 			name: "異常系: 不正なIDとパスワードでログインできないこと",
// 			setup: func() {
// 				c.SetUserIDForTest("invalid_user")
// 				c.SetPasswordForTest("invalid_password")
// 				c.SetBaseURLForTest(c.GetBaseURLForTest())
// 			},
// 			loginReq: request_auth.ReqLogin{
// 				UserId:   "invalid_user",
// 				Password: "invalid_password",
// 			},
// 			expectLoginErr: true,
// 			expectLoginRes: false,
// 			expectLoggedIn: false,
// 			expectLogout:   false,
// 		},
// 		{
// 			name: "異常系: APIサーバーがエラーを返す場合",
// 			setup: func() {
// 				c.SetUserIDForTest(c.GetUserIDForTest())
// 				c.SetPasswordForTest(c.GetPasswordForTest())
// 				c.SetBaseURLForTest("https://invalid.example.com/")
// 			},
// 			loginReq: request_auth.ReqLogin{
// 				UserId:   c.GetUserIDForTest(),
// 				Password: c.GetPasswordForTest(),
// 			},
// 			expectLoginErr: true,
// 			expectLoginRes: false,
// 			expectLoggedIn: false,
// 			expectLogout:   false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()

// 			res, err := c.Login(context.Background(), tt.loginReq)

// 			if tt.expectLoginErr {
// 				assert.Error(t, err)
// 				assert.Nil(t, res)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, res)
// 				assert.Equal(t, "0", res.ResultCode)
// 				t.Logf("Login成功後の p_no: %d", c.GetPNoForTest()) // ★追加
// 			}

// 			assert.Equal(t, tt.expectLoggedIn, c.GetLogginedForTest())

// 			if tt.expectLoggedIn {
// 				assert.NotEmpty(t, c.GetLoginInfoForTest().RequestURL)
// 			} else {
// 				assert.Nil(t, c.GetLoginInfoForTest())
// 			}

// 			if tt.expectLogout {
// 				t.Logf("Logout前の p_no: %d", c.GetPNoForTest()) // ★追加
// 				logoutRes, err := c.Logout(context.Background(), request_auth.ReqLogout{})
// 				t.Logf("Logout実行後の p_no: %d", c.GetPNoForTest()) // ★追加
// 				assert.NoError(t, err)
// 				assert.NotNil(t, logoutRes)
// 				assert.Equal(t, "0", logoutRes.ResultCode)
// 				assert.False(t, c.GetLogginedForTest())
// 				assert.Nil(t, c.GetLoginInfoForTest())
// 			}
// 		})
// 	}
// }

// func TestAuthClientImpl_LoginWithPost(t *testing.T) {
// 	c := client.CreateTestClient(t)

// 	tests := []struct {
// 		name           string
// 		setup          func()
// 		loginReq       request_auth.ReqLogin
// 		expectLoginErr bool
// 		expectLoginRes bool
// 		expectLoggedIn bool
// 		expectLogout   bool
// 	}{
// 		{
// 			name: "正常系 (POST): 正しいIDとパスワードでログインできること",
// 			setup: func() {
// 				// 正常な認証情報を設定
// 				c.SetUserIDForTest(c.GetUserIDForTest())
// 				c.SetPasswordForTest(c.GetPasswordForTest())
// 				c.SetBaseURLForTest(c.GetBaseURLForTest())
// 			},
// 			loginReq: request_auth.ReqLogin{
// 				UserId:   c.GetUserIDForTest(),
// 				Password: c.GetPasswordForTest(),
// 			},
// 			expectLoginErr: false,
// 			expectLoginRes: true,
// 			expectLoggedIn: true,
// 			expectLogout:   true,
// 		},
// 		{
// 			name: "異常系 (POST): 不正なIDとパスワードでログインできないこと",
// 			setup: func() {
// 				c.SetUserIDForTest("invalid_user")
// 				c.SetPasswordForTest("invalid_password")
// 				c.SetBaseURLForTest(c.GetBaseURLForTest())
// 			},
// 			loginReq: request_auth.ReqLogin{
// 				UserId:   "invalid_user",
// 				Password: "invalid_password",
// 			},
// 			expectLoginErr: true,
// 			expectLoginRes: false,
// 			expectLoggedIn: false,
// 			expectLogout:   false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.setup()

// 			res, err := c.LoginWithPost(context.Background(), tt.loginReq)

// 			if tt.expectLoginErr {
// 				assert.Error(t, err)
// 				assert.Nil(t, res)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, res)
// 				assert.Equal(t, "0", res.ResultCode)
// 				t.Logf("LoginWithPost成功後の p_no: %d", c.GetPNoForTest()) // ★追加
// 			}

// 			assert.Equal(t, tt.expectLoggedIn, c.GetLogginedForTest())

// 			if tt.expectLoggedIn {
// 				assert.NotEmpty(t, c.GetLoginInfoForTest().RequestURL)
// 			} else {
// 				assert.Nil(t, c.GetLoginInfoForTest())
// 			}

// 			if tt.expectLogout {
// 				t.Logf("LogoutWithPost前の p_no: %d", c.GetPNoForTest()) // ★追加
// 				logoutRes, err := c.LogoutWithPost(context.Background(), request_auth.ReqLogout{})
// 				t.Logf("LogoutWithPost実行後の p_no: %d", c.GetPNoForTest()) // ★追加
// 				assert.NoError(t, err)
// 				assert.NotNil(t, logoutRes)
// 				assert.Equal(t, "0", logoutRes.ResultCode)
// 				assert.False(t, c.GetLogginedForTest())
// 				assert.Nil(t, c.GetLoginInfoForTest())
// 			}
// 		})
// 	}
// }

// go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go

// TestAuthClientImpl_SimpleLoginLogout は、単純なログイン・ログアウトのシーケンスをテストします。
// func TestAuthClientImpl_SimpleLoginLogout(t *testing.T) {
// 	// 1. テストクライアントの作成
// 	c := client.CreateTestClient(t)

// 	// 2. ログイン
// 	loginReq := request_auth.ReqLogin{
// 		UserId:   c.GetUserIDForTest(),
// 		Password: c.GetPasswordForTest(),
// 	}
// 	loginRes, err := c.LoginWithPost(context.Background(), loginReq)

// 	// ログインの成功を確認
// 	require.NoError(t, err, "Login should not produce an error")
// 	require.NotNil(t, loginRes, "Login response should not be nil")
// 	require.Equal(t, "0", loginRes.ResultCode, "Login result code should be 0")
// 	require.True(t, c.GetLogginedForTest(), "Client should be in a logged-in state")
// 	t.Logf("ログイン成功後の p_no: %d", c.GetPNoForTest())

// 	// 3. ログアウト
// 	t.Logf("ログアウト前の p_no: %d", c.GetPNoForTest())
// 	logoutRes, err := c.LogoutWithPost(context.Background(), request_auth.ReqLogout{})

// 	// ログアウトの成功を確認
// 	require.NoError(t, err, "Logout should not produce an error")
// 	require.NotNil(t, logoutRes, "Logout response should not be nil")
// 	require.Equal(t, "0", logoutRes.ResultCode, "Logout result code should be 0")
// 	require.False(t, c.GetLogginedForTest(), "Client should be in a logged-out state")
// 	t.Logf("ログアウト実行後の p_no: %d", c.GetPNoForTest())
// }

// TestAuthClientImpl_LoginOnly は、ログインのみを行うテストです。
func TestAuthClientImpl_LoginOnly(t *testing.T) {
	// 1. テストクライアントの作成
	c := client.CreateTestClient(t)

	// 2. ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	loginRes, err := c.LoginWithPost(context.Background(), loginReq)

	// ログインの成功を確認
	require.NoError(t, err, "Login should not produce an error")
	require.NotNil(t, loginRes, "Login response should not be nil")
	require.Equal(t, "0", loginRes.ResultCode, "Login result code should be 0")
	require.True(t, c.GetLogginedForTest(), "Client should be in a logged-in state")
	t.Logf("TestAuthClientImpl_LoginOnly - ログイン成功後の p_no: %d", c.GetPNoForTest())
}

// TestAuthClientImpl_LogoutOnly は、ログイン後にログアウトのみを行うテストです。
func TestAuthClientImpl_LogoutOnly(t *testing.T) {
	// 1. テストクライアントの作成
	c := client.CreateTestClient(t)

	// 2. ログイン (ログアウトテストを行うためにまずログイン状態にする)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	loginRes, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err, "Login should not produce an error before logout test")
	require.NotNil(t, loginRes, "Login response should not be nil before logout test")
	require.Equal(t, "0", loginRes.ResultCode, "Login result code should be 0 before logout test")
	require.True(t, c.GetLogginedForTest(), "Client should be in a logged-in state before logout test")
	t.Logf("TestAuthClientImpl_LogoutOnly - ログイン成功後の p_no (ログアウト前): %d", c.GetPNoForTest())

	// 3. ログアウト
	t.Logf("TestAuthClientImpl_LogoutOnly - ログアウト前の p_no: %d", c.GetPNoForTest())
	logoutRes, err := c.LogoutWithPost(context.Background(), request_auth.ReqLogout{})

	// ログアウトの成功を確認
	require.NoError(t, err, "Logout should not produce an error")
	require.NotNil(t, logoutRes, "Logout response should not be nil")
	require.Equal(t, "0", logoutRes.ResultCode, "Logout result code should be 0")
	require.False(t, c.GetLogginedForTest(), "Client should be in a logged-out state")
	t.Logf("TestAuthClientImpl_LogoutOnly - ログアウト実行後の p_no: %d", c.GetPNoForTest())
}

// TestAuthClientImpl_Sequence_LoginWaitLogout は、ログイン、5分待機、ログアウトの一連のシーケンスをテストします。
// このテストは実行に5分以上かかります。
func TestAuthClientImpl_Sequence_LoginWaitLogout(t *testing.T) {
	t.Log("【シーケンステスト開始】ログイン → 5分待機 → ログアウト")

	// 1. ログイン
	c := client.CreateTestClient(t)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	loginRes, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err, "シーケンステスト中のログインに失敗しました")
	require.NotNil(t, loginRes)
	require.Equal(t, "0", loginRes.ResultCode, "ログインAPIからエラーが返されました")
	t.Logf("ログイン成功。p_no: %d", c.GetPNoForTest())

	// 2. 5分間待機
	const waitMinutes = 5
	t.Logf("%d分間待機します...", waitMinutes)
	time.Sleep(waitMinutes * time.Minute)
	t.Log("待機完了。")

	// 3. ログアウト
	// ここでは、新しくログインせず、既存のセッションを使ってログアウトのみを実行します。
	t.Logf("ログアウト前の p_no: %d", c.GetPNoForTest())
	logoutRes, err := c.LogoutWithPost(context.Background(), request_auth.ReqLogout{})
	require.NoError(t, err, "ログアウトAPIの呼び出し自体に失敗しました")
	require.NotNil(t, logoutRes)
	t.Logf("ログアウトAPIの応答: ResultCode=%s, ResultText=%s", logoutRes.ResultCode, logoutRes.ResultText)
	t.Logf("ログアウト実行後の p_no: %d", c.GetPNoForTest())

	// このテストでは、ログアウトAPIがどのような結果を返すか（成功するのか、セッション切れエラーか）を観察します。
}
