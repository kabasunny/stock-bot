package tests

import (
	"context"
	"fmt"
	"os"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBalanceUseCaseImpl_CanEntry_RealAPI(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ロガーの設定
	logger := zap.NewNop()

	// UseCaseImplのインスタンスを作成
	uc := app.NewBalanceUseCaseImpl(c, logger)

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
		isHolding, totalAssets, err := uc.CanEntry(context.Background(), issueCode)

		// 結果を検証する
		assert.NoError(t, err, "CanEntry should not return an error")

		//口座の状態によって結果が変わるので、環境変数を元にテストの評価をする
		expectedIsHoldingStr := os.Getenv("EXPECTED_ISHOLDING")
		expectedTotalAssetsStr := os.Getenv("EXPECTED_TOTAL_ASSETS")

		expectedIsHolding, err := strconv.ParseBool(expectedIsHoldingStr)
		if err != nil {
			t.Fatalf("Invalid value for EXPECTED_ISHOLDING: %v", err)
		}

		expectedTotalAssets, err := strconv.ParseFloat(expectedTotalAssetsStr, 64)
		if err != nil {
			t.Fatalf("Invalid value for EXPECTED_TOTAL_ASSETS: %v", err)
		}

		assert.Equal(t, expectedIsHolding, isHolding, "IsHolding value is not as expected")
		assert.Equal(t, expectedTotalAssets, totalAssets, "TotalAssets value is not as expected")

		fmt.Printf("Is holding: %v\n", isHolding)     // 実際の結果を出力して確認
		fmt.Printf("Total assets: %v\n", totalAssets) // 実際の結果を出力して確認
		// 資金の検証. 口座の状態によって結果が変わるので、検証ロジックは調整が必要
		// assert.Greater(t, totalAssets, float64(0))
	})

	// 例: 銘柄をすでに保有している場合にfalseが返ることを検証するテストケース
	// 例: APIからのデータ取得でエラーが発生した場合にエラーが伝播されることを検証するテストケース
}

// go test ./internal/app/tests -v
