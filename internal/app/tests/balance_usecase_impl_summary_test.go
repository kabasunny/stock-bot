package tests

import (
	"context"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceUseCaseImpl_GetSummary_RealAPI(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// UseCaseImplのインスタンスを作成
	uc := app.NewBalanceUseCase(c)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: GetSummaryがエラーなくサマリーを返す", func(t *testing.T) {
		// APIを呼び出す
		summary, err := uc.GetSummary(context.Background())

		// 結果を検証する
		assert.NoError(t, err, "GetSummary should not return an error")
		assert.NotNil(t, summary, "Summary should not be nil")

		// 返された値の基本的なチェック (口座状況に依存しない)
		assert.Greater(t, summary.MarginRate, 0.0, "MarginRate should be greater than 0")
		assert.False(t, summary.UpdatedAt.IsZero(), "UpdatedAt should be set")
	})
}
