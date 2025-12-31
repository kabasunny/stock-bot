// internal/infrastructure/client/tests/auth_client_impl_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"testing"

	"time" // timeパッケージをインポート

	"github.com/stretchr/testify/assert"
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
	session, err := c.LoginWithPost(context.Background(), loginReq)

	// ログインの成功を確認
	require.NoError(t, err, "Login should not produce an error")
	require.NotNil(t, session, "Session should not be nil")
	require.Equal(t, "0", session.ResultCode, "Login result code should be 0")

	// セッション内容の詳細確認
	assert.NotEmpty(t, session.RequestURL, "RequestURL should not be empty")
	assert.NotEmpty(t, session.MasterURL, "MasterURL should not be empty")
	assert.NotEmpty(t, session.PriceURL, "PriceURL should not be empty")
	assert.NotEmpty(t, session.EventURL, "EventURL should not be empty")
	assert.NotNil(t, session.CookieJar, "CookieJar should not be nil")

	// p_noの初期値確認
	pNo := session.GetPNo()
	assert.Greater(t, pNo, int32(0), "PNo should be greater than 0")

	t.Logf("TestAuthClientImpl_LoginOnly - ログイン成功")
	t.Logf("Session details - RequestURL: %s", session.RequestURL)
	t.Logf("Session details - PNo: %d", pNo)
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
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err, "Login should not produce an error before logout test")
	require.NotNil(t, session, "Session should not be nil before logout test")
	t.Logf("TestAuthClientImpl_LogoutOnly - ログイン成功")

	// ログイン直後のp_no確認
	pNoBeforeLogout := session.GetPNo()
	t.Logf("PNo before logout: %d", pNoBeforeLogout)

	// 3. ログアウト
	logoutRes, err := c.LogoutWithPost(context.Background(), session, request_auth.ReqLogout{})

	// ログアウトの成功を確認
	require.NoError(t, err, "Logout should not produce an error")
	require.NotNil(t, logoutRes, "Logout response should not be nil")
	require.Equal(t, "0", logoutRes.ResultCode, "Logout result code should be 0")

	// ログアウト後のp_no確認（インクリメントされているはず）
	pNoAfterLogout := session.GetPNo()
	assert.Greater(t, pNoAfterLogout, pNoBeforeLogout, "PNo should increment after logout")

	t.Logf("TestAuthClientImpl_LogoutOnly - ログアウト成功")
	t.Logf("PNo after logout: %d", pNoAfterLogout)
}

func TestAuthClientImpl_MultipleSessions(t *testing.T) {
	c := client.CreateTestClient(t)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}

	// 1. ログイン1 -> セッション1 を取得
	t.Log("Executing Login 1...")
	session1, err1 := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err1)
	require.NotNil(t, session1)
	require.Equal(t, "0", session1.ResultCode, "Login 1 failed. ResultCode: %s", session1.ResultCode)
	t.Logf("Login 1 successful. Session1 ResultCode: %s", session1.ResultCode)

	// 2. すぐにログイン2 -> セッション2 を取得
	t.Log("Executing Login 2...")
	session2, err2 := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err2)
	require.NotNil(t, session2)
	require.Equal(t, "0", session2.ResultCode, "Login 2 failed. ResultCode: %s", session2.ResultCode)
	t.Logf("Login 2 successful. Session2 ResultCode: %s", session2.ResultCode)

	// セッションオブジェクト自体が異なることを確認（念のため）
	assert.NotSame(t, session1, session2, "Sessions should be different objects")
	// Session内のpNoが異なることを確認 (atomic.Int32はポインタ比較では同じになる可能性があるため、具体的な値を確認)
	// ただし、pNoは次のリクエストでインクリメントされるため、直後の比較は意味がない。
	// APIから返される各Sessionは独立していることを期待するため、オブジェクト自体の比較で十分。

	// 3. ログアウト1 (セッション1を使用)
	t.Log("Executing Logout 1 with session 1...")
	logoutRes1, err3 := c.LogoutWithPost(context.Background(), session1, request_auth.ReqLogout{})
	require.NoError(t, err3)
	require.NotNil(t, logoutRes1)
	// ログイン2によってセッション1は無効化されているため、ResultCodeは"0"ではないことを期待する
	assert.NotEqual(t, "0", logoutRes1.ResultCode, "Logout 1 with session 1 should have failed, but it succeeded.")
	t.Logf("Logout 1 with invalidated session failed as expected. ResultCode: %s, ResultText: %s", logoutRes1.ResultCode, logoutRes1.ResultText)

	// 4. ログアウト2 (セッション2を使用)
	t.Log("Executing Logout 2 with session 2...")
	logoutRes2, err4 := c.LogoutWithPost(context.Background(), session2, request_auth.ReqLogout{})
	require.NoError(t, err4)
	require.NotNil(t, logoutRes2)
	assert.Equal(t, "0", logoutRes2.ResultCode, "Logout 2 failed with session 2. ResultCode: %s, ResultText: %s", logoutRes2.ResultCode, logoutRes2.ResultText)
	t.Logf("Logout 2 Result (Session2): Code=%s, Text=%s", logoutRes2.ResultCode, logoutRes2.ResultText)
}

// go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_MultipleSessions

// TestAuthClientImpl_Sequence_LoginWaitLogoutLogin は、ログイン、5分待機、ログアウト、再ログインの一連のシーケンスをテストします。
// このテストは実行に5分以上かかります。
func TestAuthClientImpl_Sequence_LoginWaitLogoutLogin(t *testing.T) {
	t.Log("【シーケンステスト開始】ログイン → 5分待機 → ログアウト → 再ログイン")

	// 1. ログイン
	c := client.CreateTestClient(t)
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err, "シーケンステスト中の初回ログインに失敗しました")
	require.NotNil(t, session, "初回ログインのSessionがnilです")
	require.Equal(t, "0", session.ResultCode, "初回ログインAPIからエラーが返されました")
	t.Logf("ログイン成功。")

	// 2. 5分間待機
	const waitMinutes = 5
	t.Logf("%d分間待機します...", waitMinutes)
	time.Sleep(waitMinutes * time.Minute)
	t.Log("待機完了。")

	// 3. ログアウト
	logoutRes, err := c.LogoutWithPost(context.Background(), session, request_auth.ReqLogout{})
	require.NoError(t, err, "ログアウトAPIの呼び出し自体に失敗しました")
	require.NotNil(t, logoutRes, "ログアウトのレスポンスがnilです")
	t.Logf("ログアウトAPIの応答: ResultCode=%s, ResultText=%s", logoutRes.ResultCode, logoutRes.ResultText)
	// ログアウト自体は成功することを期待
	assert.Equal(t, "0", logoutRes.ResultCode, "5分後のログアウトが失敗しました")

	// 4. 再ログイン
	t.Log("再ログインを試みます...")
	reloginSession, reloginErr := c.LoginWithPost(context.Background(), loginReq)

	// 電話認証の有効期間(3分)が切れているため、再ログインは失敗することを期待する
	// APIからエラーが返されるため、reloginErrがnilではないことを検証する
	require.Error(t, reloginErr, "電話認証が切れているため、再ログインはエラーを返すはずでした")
	// TEST-002の結果から、エラーメッセージに "10089" が含まれることを期待
	assert.Contains(t, reloginErr.Error(), "10089", "エラーメッセージに result code 10089 が含まれていません")
	// reloginSession はエラー時にはnilである可能性が高いため、nilであることを確認
	assert.Nil(t, reloginSession, "再ログイン失敗時、セッションはnilであるはずです")

	t.Logf("期待通り、再ログインに失敗しました。エラー: %v", reloginErr)
}

// go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_Sequence_LoginWaitLogoutLogin

// TestAuthClientImpl_InvalidCredentials は、不正な認証情報でのログインテストです。
func TestAuthClientImpl_InvalidCredentials(t *testing.T) {
	c := client.CreateTestClient(t)

	// 不正な認証情報でログイン試行
	loginReq := request_auth.ReqLogin{
		UserId:   "invalid_user",
		Password: "invalid_password",
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)

	// ログインの失敗を確認
	require.Error(t, err, "Login with invalid credentials should produce an error")
	require.Nil(t, session, "Session should be nil when login fails")

	// エラーメッセージの確認
	assert.Contains(t, err.Error(), "login failed", "Error message should contain 'login failed'")

	t.Logf("Expected login failure with invalid credentials: %v", err)
}

// TestAuthClientImpl_EmptyCredentials は、空の認証情報でのログインテストです。
func TestAuthClientImpl_EmptyCredentials(t *testing.T) {
	c := client.CreateTestClient(t)

	// 空の認証情報でログイン試行
	loginReq := request_auth.ReqLogin{
		UserId:   "",
		Password: "",
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)

	// ログインの失敗を確認
	require.Error(t, err, "Login with empty credentials should produce an error")
	require.Nil(t, session, "Session should be nil when login fails")

	t.Logf("Expected login failure with empty credentials: %v", err)
}

// TestAuthClientImpl_LogoutWithoutLogin は、ログインせずにログアウトを試行するテストです。
func TestAuthClientImpl_LogoutWithoutLogin(t *testing.T) {
	c := client.CreateTestClient(t)

	// ダミーセッションを作成（実際にはログインしていない）
	dummySession := client.NewSession()
	dummySession.RequestURL = "https://example.com/dummy"

	// ログアウト試行
	logoutRes, err := c.LogoutWithPost(context.Background(), dummySession, request_auth.ReqLogout{})

	// ログアウトAPIの呼び出し自体は成功するが、ResultCodeでエラーが返される
	require.NoError(t, err, "Logout API call should not produce an error")
	require.NotNil(t, logoutRes, "Logout response should not be nil")

	// 無効なセッションなのでResultCodeは"0"以外になるはず
	assert.NotEqual(t, "0", logoutRes.ResultCode, "Logout with invalid session should fail")

	t.Logf("Expected logout failure with invalid session: ResultCode=%s, ResultText=%s",
		logoutRes.ResultCode, logoutRes.ResultText)
}

// TestAuthClientImpl_LogoutWithNilSession は、nilセッションでのログアウトテストです。
func TestAuthClientImpl_LogoutWithNilSession(t *testing.T) {
	c := client.CreateTestClient(t)

	// nilセッションでログアウト試行
	logoutRes, err := c.LogoutWithPost(context.Background(), nil, request_auth.ReqLogout{})

	// nilセッションの場合はエラーが返されるはず
	require.Error(t, err, "Logout with nil session should produce an error")
	require.Nil(t, logoutRes, "Logout response should be nil when session is nil")

	// エラーメッセージの確認
	assert.Contains(t, err.Error(), "session is nil", "Error message should contain 'session is nil'")

	t.Logf("Expected logout failure with nil session: %v", err)
}
