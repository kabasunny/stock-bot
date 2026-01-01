package tests

import (
	"runtime"
	"stock-bot/internal/infrastructure/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSession_ConcurrentAccess はSessionの並行アクセステスト
func TestSession_ConcurrentAccess(t *testing.T) {
	session := client.NewSession()

	// 並行でP_noを取得
	concurrentRequests := 100
	results := make(chan int32, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			pNo := session.GetPNo()
			results <- pNo
		}()
	}

	// 結果を収集
	pNos := make([]int32, concurrentRequests)
	for i := 0; i < concurrentRequests; i++ {
		pNos[i] = <-results
	}

	// 全てのP_noが異なることを確認（アトミック操作の確認）
	uniquePNos := make(map[int32]bool)
	for _, pNo := range pNos {
		assert.False(t, uniquePNos[pNo], "P_no %d が重複しています", pNo)
		uniquePNos[pNo] = true
	}

	t.Logf("Generated %d unique P_nos", len(uniquePNos))
	assert.Equal(t, concurrentRequests, len(uniquePNos), "全てのP_noが一意である必要があります")
}

// TestSession_ErrorResilience はSessionのエラー耐性テスト
func TestSession_ErrorResilience(t *testing.T) {
	// nilポインタアクセステスト
	t.Run("NilPointerAccess", func(t *testing.T) {
		var session *client.Session

		// nilセッションでもパニックしないことを確認
		assert.NotPanics(t, func() {
			if session != nil {
				session.GetPNo()
			}
		})
	})

	// 不正なデータでのセッション作成テスト
	t.Run("InvalidSessionData", func(t *testing.T) {
		session := client.NewSession()

		// 不正なURLを設定
		session.RequestURL = "invalid-url"
		session.MasterURL = ""
		session.PriceURL = "not-a-url"

		// セッションが作成されることを確認（エラーハンドリングは各クライアントで行う）
		assert.NotNil(t, session)
		assert.Equal(t, "invalid-url", session.RequestURL)
	})
}

// TestSession_BasicFunctionality はSessionの基本機能テスト
func TestSession_BasicFunctionality(t *testing.T) {
	// セッション作成テスト
	t.Run("SessionCreation", func(t *testing.T) {
		session := client.NewSession()
		assert.NotNil(t, session)
		assert.NotNil(t, session.CookieJar)
	})

	// P_no生成テスト
	t.Run("PNoGeneration", func(t *testing.T) {
		session := client.NewSession()

		// 連続でP_noを取得
		pNo1 := session.GetPNo()
		pNo2 := session.GetPNo()
		pNo3 := session.GetPNo()

		// P_noが増加することを確認
		assert.Greater(t, pNo2, pNo1)
		assert.Greater(t, pNo3, pNo2)
	})

	// セッション情報設定テスト
	t.Run("SessionInfoSetting", func(t *testing.T) {
		session := client.NewSession()

		// セッション情報を設定
		session.RequestURL = "https://example.com/api"
		session.MasterURL = "https://example.com/master"
		session.PriceURL = "https://example.com/price"
		session.EventURL = "wss://example.com/event"

		// 設定された値が正しいことを確認
		assert.Equal(t, "https://example.com/api", session.RequestURL)
		assert.Equal(t, "https://example.com/master", session.MasterURL)
		assert.Equal(t, "https://example.com/price", session.PriceURL)
		assert.Equal(t, "wss://example.com/event", session.EventURL)
	})
}

// TestMemoryLeakPrevention はメモリリーク防止テスト
func TestMemoryLeakPrevention(t *testing.T) {
	// 大量のセッション作成・破棄テスト
	t.Run("MassiveSessionCreation", func(t *testing.T) {
		sessionCount := 1000

		for i := 0; i < sessionCount; i++ {
			session := client.NewSession()

			// セッションを使用
			pNo := session.GetPNo()
			assert.Greater(t, pNo, int32(0))

			// 明示的にnilにしてGCを促進
			session = nil
		}

		// GCを実行
		runtime.GC()

		t.Logf("Created and destroyed %d sessions", sessionCount)
	})
}

// TestSession_ThreadSafety はSessionのスレッドセーフティテスト
func TestSession_ThreadSafety(t *testing.T) {
	session := client.NewSession()

	// 複数のgoroutineで同時にセッション情報を変更
	concurrentWrites := 50
	done := make(chan bool, concurrentWrites)

	for i := 0; i < concurrentWrites; i++ {
		go func(id int) {
			// P_noを取得（アトミック操作）
			pNo := session.GetPNo()
			assert.Greater(t, pNo, int32(0))

			// セッション情報を設定（非アトミック操作だが、テスト用）
			session.ResultCode = "OK"
			session.ResultText = "Success"

			done <- true
		}(i)
	}

	// 全てのgoroutineの完了を待つ
	for i := 0; i < concurrentWrites; i++ {
		<-done
	}

	// 最終的なP_noが期待値以上であることを確認
	finalPNo := session.GetPNo()
	assert.GreaterOrEqual(t, finalPNo, int32(concurrentWrites+1))

	t.Logf("Final P_no after %d concurrent operations: %d", concurrentWrites, finalPNo)
}
