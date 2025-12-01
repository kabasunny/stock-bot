// internal/infrastructure/client/tests/master_data_client_impl_test.go
package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/master/request"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestMasterDataClientImpl_DownloadMasterDataWithPost(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	_, err := c.LoginWithPost(context.Background(), loginReq)
	assert.NoError(t, err)

	// 正常系: システムステータスのダウンロードが成功すること
	// t.Run("正常系: 1.システムステータスのダウンロードが成功すること", func(t *testing.T) {
	// 	// マスタ情報ダウンロードのリクエストデータを作成 (システムステータスのみ)
	// 	downloadReq := request.ReqDownloadMaster{
	// 		TargetCLMID: "CLMSystemStatus,CLMEventDownloadComplete", // システムステータスのみ
	// 	}

	// 	// DownloadMasterData メソッドを実行
	// 	res, err := c.DownloadMasterData(context.Background(), downloadReq)

	// 	// レスポンスとエラーをチェック
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, res)

	// 	// SystemStatus が取得できているかチェック
	// 	if res != nil && res.SystemStatus.CLMID == "CLMSystemStatus" {
	// 		assert.Equal(t, "CLMSystemStatus", res.SystemStatus.CLMID)

	// 		// 構造体を JSON 形式に変換して出力
	// 		jsonData, err := json.MarshalIndent(res.SystemStatus, "", "  ") // インデント付きで出力
	// 		if err != nil {
	// 			t.Errorf("JSON 変換エラー: %v", err)
	// 			return
	// 		}
	// 		fmt.Println("SystemStatus Data (JSON):")
	// 		fmt.Println(string(jsonData))

	// 	} else {
	// 		t.Errorf("SystemStatus が取得できていません。レスポンス: %v", res)
	// 	}
	// })

	// 正常系: DateInfo のダウンロードが成功すること
	// t.Run("正常系: 2.DateInfo のダウンロードが成功すること", func(t *testing.T) {
	// 	// マスタ情報ダウンロードのリクエストデータを作成 (DateInfoのみ)
	// 	downloadReq := request.ReqDownloadMaster{
	// 		TargetCLMID: "CLMDateZyouhou,CLMEventDownloadComplete", // DateInfo のみ
	// 	}

	// 	// DownloadMasterData メソッドを実行
	// 	res, err := c.DownloadMasterData(context.Background(), downloadReq)

	// 	// レスポンスとエラーをチェック
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, res)

	// 	// DateInfo が取得できているかチェック
	// 	if res != nil && res.DateInfo.CLMID == "CLMDateZyouhou" {
	// 		assert.Equal(t, "CLMDateZyouhou", res.DateInfo.CLMID)

	// 		// 構造体を JSON 形式に変換して出力
	// 		jsonData, err := json.MarshalIndent(res.DateInfo, "", "  ") // インデント付きで出力
	// 		if err != nil {
	// 			t.Errorf("JSON 変換エラー: %v", err)
	// 			return
	// 		}
	// 		fmt.Println("DateInfo Data (JSON):")
	// 		fmt.Println(string(jsonData))

	// 	} else {
	// 		t.Errorf("DateInfo が取得できていません。レスポンス: %v", res)
	// 	}
	// })

	// 正常系: 呼値 (TickRule) のダウンロードが成功すること
	// t.Run("正常系: 3.呼値 (TickRule) のダウンロードが成功すること", func(t *testing.T) {
	// 	// マスタ情報ダウンロードのリクエストデータを作成 (TickRuleのみ)
	// 	downloadReq := request.ReqDownloadMaster{
	// 		TargetCLMID: "CLMYobine,CLMEventDownloadComplete", // TickRule のみ
	// 	}

	// 	// DownloadMasterData メソッドを実行
	// 	res, err := c.DownloadMasterData(context.Background(), downloadReq)

	// 	// レスポンスとエラーをチェック
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, res)

	// 	// TickRule が取得できているかチェック
	// 	if res != nil && len(res.TickRule) > 0 {
	// 		found := false
	// 		for _, tickRule := range res.TickRule {
	// 			if tickRule.CLMID == "CLMYobine" && tickRule.TickUnitNumber == "101" {
	// 				found = true

	// 				// 構造体を JSON 形式に変換して出力
	// 				jsonData, err := json.MarshalIndent(tickRule, "", "  ") // インデント付きで出力
	// 				if err != nil {
	// 					t.Errorf("JSON 変換エラー: %v", err)
	// 					return
	// 				}
	// 				fmt.Println("TickRule Data (JSON):")
	// 				fmt.Println(string(jsonData))

	// 				// 必要に応じて、他のフィールドもチェック
	// 				// assert.Equal(t, "20140101", tickRule.ApplicableDate) // 例: 適用日が 20140101 であることを確認
	// 				break // "101" が見つかったらループを抜ける
	// 			}
	// 		}

	// 		if !found {
	// 			t.Errorf("sYobineTaniNumber が 101 の TickRule が見つかりませんでした。レスポンス: %v", res)
	// 		}

	// 	} else {
	// 		t.Errorf("TickRule が取得できていません。レスポンス: %v", res)
	// 	}
	// })

	// 正常系: 全マスタ情報のダウンロードが成功すること (TargetCLMID 未指定)
	t.Run("正常系: 4.全マスタ情報のダウンロードが成功すること (TargetCLMID 未指定)", func(t *testing.T) {
		// マスタ情報ダウンロードのリクエストデータを作成 (TargetCLMID 未指定)
		downloadReq := request.ReqDownloadMaster{}

		// DownloadMasterData メソッドを実行
		res, err := c.DownloadMasterData(context.Background(), downloadReq)

		// レスポンスとエラーをチェック
		assert.NoError(t, err)
		assert.NotNil(t, res)

		// 構造体の内容を表示 (存在する場合のみ)
		if res != nil {
			fmt.Println("--- 全マスタ情報 ---")

			if res.SystemStatus.CLMID != "" {
				fmt.Println("SystemStatus Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.SystemStatus, "", "  ")
				fmt.Println(string(jsonData))
			}

			if len(res.DateInfo) > 0 {
				fmt.Println("DateInfo Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.DateInfo[0], "", "  ")
				fmt.Println(string(jsonData))
			}

			if len(res.TickRule) > 0 {
				fmt.Println("TickRule Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.TickRule[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.OperationStatus) > 0 {
				fmt.Println("OperationStatus Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.OperationStatus[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.OperationStatusStock) > 0 {
				fmt.Println("OperationStatusStock Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.OperationStatusStock[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.OperationStatusDerivative) > 0 {
				fmt.Println("OperationStatusDerivative Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.OperationStatusDerivative[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.StockMaster) > 0 {
				fmt.Println("StockMaster Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.StockMaster[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.StockMarketMaster) > 0 {
				fmt.Println("StockMarketMaster Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.StockMarketMaster[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.StockIssueRegulation) > 0 {
				fmt.Println("StockIssueRegulation Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.StockIssueRegulation[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.FutureMaster) > 0 {
				fmt.Println("FutureMaster Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.FutureMaster[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.OptionMaster) > 0 {
				fmt.Println("OptionMaster Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.OptionMaster[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.FutureOptionRegulation) > 0 {
				fmt.Println("FutureOptionRegulation Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.FutureOptionRegulation[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.MarginRate) > 0 {
				fmt.Println("MarginRate Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.MarginRate[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.MarginMaster) > 0 {
				fmt.Println("MarginMaster Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.MarginMaster[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}

			if len(res.ErrorReason) > 0 {
				fmt.Println("ErrorReason Data (JSON):")
				jsonData, _ := json.MarshalIndent(res.ErrorReason[0], "", "  ") // 最初の要素のみ
				fmt.Println(string(jsonData))
			}
		} else {
			t.Errorf("マスタ情報が取得できていません。レスポンス: %v", res)
		}
	})

}

// go test -v ./internal/infrastructure/client/tests/master_data_client_impl_download_masterdata_test.go
