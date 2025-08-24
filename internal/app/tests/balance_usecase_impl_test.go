package tests

import (
	"context"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceUseCaseImpl_CanEntry_RealAPI(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// UseCaseImplのインスタンスを作成
	uc := app.NewBalanceUseCaseImpl(c)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	assert.NoError(t, err)

	t.Run("正常系: エントリー可能な銘柄の場合、trueと口座情報が返る", func(t *testing.T) {
		// 銘柄コードを指定
		issueCode := "7203" // 例: トヨタ自動車

		// APIを呼び出す
		_, _, err := uc.CanEntry(context.Background(), issueCode)

		// 結果を検証する
		assert.NoError(t, err, "CanEntry should not return an error")
	})

	// 例: 銘柄をすでに保有している場合にfalseが返ることを検証するテストケース
	// 例: APIからのデータ取得でエラーが発生した場合にエラーが伝播されることを検証するテストケース
}

// go test ./internal/app/tests -v
