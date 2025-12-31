package tests

import (
	"net/http/cookiejar"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewSession は NewSession() の基本動作をテストします
func TestNewSession(t *testing.T) {
	session := client.NewSession()

	// 基本的な初期化の確認
	assert.NotNil(t, session, "Session should not be nil")
	assert.NotNil(t, session.CookieJar, "CookieJar should be initialized")

	// CookieJarの型確認
	_, ok := session.CookieJar.(*cookiejar.Jar)
	assert.True(t, ok, "CookieJar should be of type *cookiejar.Jar")

	// 初期値の確認
	assert.Empty(t, session.ResultCode, "ResultCode should be empty initially")
	assert.Empty(t, session.ResultText, "ResultText should be empty initially")
	assert.Empty(t, session.SecondPassword, "SecondPassword should be empty initially")
	assert.Empty(t, session.RequestURL, "RequestURL should be empty initially")
	assert.Empty(t, session.MasterURL, "MasterURL should be empty initially")
	assert.Empty(t, session.PriceURL, "PriceURL should be empty initially")
	assert.Empty(t, session.EventURL, "EventURL should be empty initially")
}

// TestSession_GetPNo は GetPNo() のアトミック操作をテストします
func TestSession_GetPNo(t *testing.T) {
	session := client.NewSession()

	// 初回呼び出し（0から1にインクリメント）
	pNo1 := session.GetPNo()
	assert.Equal(t, int32(1), pNo1, "First GetPNo() should return 1")

	// 2回目呼び出し（1から2にインクリメント）
	pNo2 := session.GetPNo()
	assert.Equal(t, int32(2), pNo2, "Second GetPNo() should return 2")

	// 3回目呼び出し（2から3にインクリメント）
	pNo3 := session.GetPNo()
	assert.Equal(t, int32(3), pNo3, "Third GetPNo() should return 3")

	// 連続性の確認
	assert.Equal(t, pNo1+1, pNo2, "PNo should increment by 1")
	assert.Equal(t, pNo2+1, pNo3, "PNo should increment by 1")
}

// TestSession_GetPNo_Concurrent は GetPNo() の並行安全性をテストします
func TestSession_GetPNo_Concurrent(t *testing.T) {
	session := client.NewSession()
	const goroutines = 100
	const callsPerGoroutine = 10

	results := make(chan int32, goroutines*callsPerGoroutine)

	// 複数のgoroutineで同時にGetPNo()を呼び出し
	for i := 0; i < goroutines; i++ {
		go func() {
			for j := 0; j < callsPerGoroutine; j++ {
				results <- session.GetPNo()
			}
		}()
	}

	// 結果を収集
	pNos := make([]int32, 0, goroutines*callsPerGoroutine)
	for i := 0; i < goroutines*callsPerGoroutine; i++ {
		pNos = append(pNos, <-results)
	}

	// 重複がないことを確認
	seen := make(map[int32]bool)
	for _, pNo := range pNos {
		assert.False(t, seen[pNo], "PNo %d should be unique", pNo)
		seen[pNo] = true
	}

	// 期待される範囲内であることを確認
	for _, pNo := range pNos {
		assert.True(t, pNo >= 1 && pNo <= int32(goroutines*callsPerGoroutine),
			"PNo %d should be in range [1, %d]", pNo, goroutines*callsPerGoroutine)
	}
}

// TestSession_SetLoginResponse は SetLoginResponse() の動作をテストします
func TestSession_SetLoginResponse(t *testing.T) {
	session := client.NewSession()

	// テスト用のResLoginを作成
	resLogin := &response.ResLogin{
		ResultCode:         "0",
		ResultText:         "Success",
		RequestURL:         "https://example.com/request",
		MasterURL:          "https://example.com/master",
		PriceURL:           "https://example.com/price",
		SUrlEventWebSocket: "wss://example.com/event",
	}

	// SetLoginResponseを実行
	session.SetLoginResponse(resLogin)

	// 設定された値を確認
	assert.Equal(t, "0", session.ResultCode, "ResultCode should be set correctly")
	assert.Equal(t, "Success", session.ResultText, "ResultText should be set correctly")
	assert.Equal(t, "https://example.com/request", session.RequestURL, "RequestURL should be set correctly")
	assert.Equal(t, "https://example.com/master", session.MasterURL, "MasterURL should be set correctly")
	assert.Equal(t, "https://example.com/price", session.PriceURL, "PriceURL should be set correctly")
	assert.Equal(t, "wss://example.com/event", session.EventURL, "EventURL should be set correctly")
}

// TestSession_SetLoginResponse_NilInput は SetLoginResponse() にnilを渡した場合の動作をテストします
func TestSession_SetLoginResponse_NilInput(t *testing.T) {
	session := client.NewSession()

	// nilを渡してもパニックしないことを確認
	require.NotPanics(t, func() {
		session.SetLoginResponse(nil)
	}, "SetLoginResponse should not panic with nil input")
}

// TestSession_SetLoginResponse_EmptyValues は SetLoginResponse() に空の値を渡した場合の動作をテストします
func TestSession_SetLoginResponse_EmptyValues(t *testing.T) {
	session := client.NewSession()

	// 空の値を持つResLoginを作成
	resLogin := &response.ResLogin{
		ResultCode:         "",
		ResultText:         "",
		RequestURL:         "",
		MasterURL:          "",
		PriceURL:           "",
		SUrlEventWebSocket: "",
	}

	// SetLoginResponseを実行
	session.SetLoginResponse(resLogin)

	// 空の値が正しく設定されることを確認
	assert.Empty(t, session.ResultCode, "Empty ResultCode should be set")
	assert.Empty(t, session.ResultText, "Empty ResultText should be set")
	assert.Empty(t, session.RequestURL, "Empty RequestURL should be set")
	assert.Empty(t, session.MasterURL, "Empty MasterURL should be set")
	assert.Empty(t, session.PriceURL, "Empty PriceURL should be set")
	assert.Empty(t, session.EventURL, "Empty EventURL should be set")
}
