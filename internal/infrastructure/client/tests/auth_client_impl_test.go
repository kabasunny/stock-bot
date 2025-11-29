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
	c := client.CreateTestClient(t)

	tests := []struct {
		name           string
		setup          func()
		loginReq       request_auth.ReqLogin
		expectLoginErr bool
		expectLoginRes bool
		expectLoggedIn bool
		expectLogout   bool
	}{
		{
			name: "正常系: 正しいIDとパスワードでログインできること",
			setup: func() {
				// 正常な認証情報を設定
				c.SetUserIDForTest(c.GetUserIDForTest())
				c.SetPasswordForTest(c.GetPasswordForTest())
				c.SetBaseURLForTest(c.GetBaseURLForTest())
			},
			loginReq: request_auth.ReqLogin{
				UserId:   c.GetUserIDForTest(),
				Password: c.GetPasswordForTest(),
			},
			expectLoginErr: false,
			expectLoginRes: true,
			expectLoggedIn: true,
			expectLogout:   true,
		},
		{
			name: "異常系: 不正なIDとパスワードでログインできないこと",
			setup: func() {
				c.SetUserIDForTest("invalid_user")
				c.SetPasswordForTest("invalid_password")
				c.SetBaseURLForTest(c.GetBaseURLForTest())
			},
			loginReq: request_auth.ReqLogin{
				UserId:   "invalid_user",
				Password: "invalid_password",
			},
			expectLoginErr: true,
			expectLoginRes: false,
			expectLoggedIn: false,
			expectLogout:   false,
		},
		{
			name: "異常系: APIサーバーがエラーを返す場合",
			setup: func() {
				c.SetUserIDForTest(c.GetUserIDForTest())
				c.SetPasswordForTest(c.GetPasswordForTest())
				c.SetBaseURLForTest("https://invalid.example.com/")
			},
			loginReq: request_auth.ReqLogin{
				UserId:   c.GetUserIDForTest(),
				Password: c.GetPasswordForTest(),
			},
			expectLoginErr: true,
			expectLoginRes: false,
			expectLoggedIn: false,
			expectLogout:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			res, err := c.Login(context.Background(), tt.loginReq)

			if tt.expectLoginErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "0", res.ResultCode)
			}

			assert.Equal(t, tt.expectLoggedIn, c.GetLogginedForTest())

			if tt.expectLoggedIn {
				assert.NotEmpty(t, c.GetLoginInfoForTest().RequestURL)
			} else {
				assert.Nil(t, c.GetLoginInfoForTest())
			}

			if tt.expectLogout {
				logoutRes, err := c.Logout(context.Background(), request_auth.ReqLogout{})
				assert.NoError(t, err)
				assert.NotNil(t, logoutRes)
				assert.Equal(t, "0", logoutRes.ResultCode)
				assert.False(t, c.GetLogginedForTest())
				assert.Nil(t, c.GetLoginInfoForTest())
			}
		})
	}
}

// go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go

func TestAuthClientImpl_LoginWithPost(t *testing.T) {
	c := client.CreateTestClient(t)

	tests := []struct {
		name           string
		setup          func()
		loginReq       request_auth.ReqLogin
		expectLoginErr bool
		expectLoginRes bool
		expectLoggedIn bool
		expectLogout   bool
	}{
		{
			name: "正常系 (POST): 正しいIDとパスワードでログインできること",
			setup: func() {
				// 正常な認証情報を設定
				c.SetUserIDForTest(c.GetUserIDForTest())
				c.SetPasswordForTest(c.GetPasswordForTest())
				c.SetBaseURLForTest(c.GetBaseURLForTest())
			},
			loginReq: request_auth.ReqLogin{
				UserId:   c.GetUserIDForTest(),
				Password: c.GetPasswordForTest(),
			},
			expectLoginErr: false,
			expectLoginRes: true,
			expectLoggedIn: true,
			expectLogout:   true,
		},
		{
			name: "異常系 (POST): 不正なIDとパスワードでログインできないこと",
			setup: func() {
				c.SetUserIDForTest("invalid_user")
				c.SetPasswordForTest("invalid_password")
				c.SetBaseURLForTest(c.GetBaseURLForTest())
			},
			loginReq: request_auth.ReqLogin{
				UserId:   "invalid_user",
				Password: "invalid_password",
			},
			expectLoginErr: true,
			expectLoginRes: false,
			expectLoggedIn: false,
			expectLogout:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			res, err := c.LoginWithPost(context.Background(), tt.loginReq)

			if tt.expectLoginErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "0", res.ResultCode)
			}

			assert.Equal(t, tt.expectLoggedIn, c.GetLogginedForTest())

			if tt.expectLoggedIn {
				assert.NotEmpty(t, c.GetLoginInfoForTest().RequestURL)
			} else {
				assert.Nil(t, c.GetLoginInfoForTest())
			}

			if tt.expectLogout {
				logoutRes, err := c.LogoutWithPost(context.Background(), request_auth.ReqLogout{})
				assert.NoError(t, err)
				assert.NotNil(t, logoutRes)
				assert.Equal(t, "0", logoutRes.ResultCode)
				assert.False(t, c.GetLogginedForTest())
				assert.Nil(t, c.GetLoginInfoForTest())
			}
		})
	}
}