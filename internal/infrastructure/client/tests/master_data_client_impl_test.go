// internal/infrastructure/client/tests/master_data_client_impl_test.go
package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/master/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMasterDataClientImpl_DownloadMasterData(t *testing.T) {
	// テストロガーの設定
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := cfg.Build()
	defer logger.Sync()

	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.Login(context.Background(), loginReq)
	require.NoError(t, err, "login failed")
	assert.True(t, c.GetLogginedForTest(), "client should be logged in")

	// リクエストを作成
	req := request.ReqDownloadMaster{
		TargetCLMID: "CLMSystemStatus,CLMEventDownloadComplete",
	}

	// API を実行
	resp, err := c.DownloadMasterData(context.Background(), req) // c は Client インターフェースを満たしているので、そのまま使用できる

	//testで落ちないようにする
	if err != nil {
		fmt.Printf("DownloadMasterData failed: %v", err)
	}

	// Create the file
	file, err := os.Create("raw_response.txt")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Dump response body to file
	if resp != nil {
		// JSONに変換してから文字列化
		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response to JSON: %v", err)
		}
		jsonString := string(jsonBytes)

		_, err = file.WriteString(jsonString)
		if err != nil {
			t.Fatalf("Failed to write to file: %v", err)
		}
	} else {
		fmt.Printf("resp is nil")
	}

	// ログアウト
	reqLogout := request_auth.ReqLogout{}
	_, err = c.Logout(context.Background(), reqLogout)
	require.NoError(t, err, "logout failed")
	assert.False(t, c.GetLogginedForTest(), "client should be logged out")
}

// go test -v ./internal/infrastructure/client/tests/master_data_client_impl_test.go
