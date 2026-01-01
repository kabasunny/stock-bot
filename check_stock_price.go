package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"stock-bot/internal/config"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
)

func main() {
	fmt.Println("=== 銘柄価格確認 ===")

	// 1. 設定ファイルの読み込み
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get caller information")
	}
	envPath := filepath.Join(filepath.Dir(filename), ".env")

	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. クライアント作成
	tachibanaClient := client.NewTachibanaClient(cfg)
	ctx := context.Background()

	// 3. ログイン
	loginReq := request.ReqLogin{
		UserId:   cfg.TachibanaUserID,
		Password: cfg.TachibanaPassword,
	}

	session, err := tachibanaClient.LoginWithPost(ctx, loginReq)
	if err != nil {
		fmt.Printf("ログインエラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("ログイン成功\n")

	// 4. 銘柄価格を取得
	symbols := []string{"6658", "7203", "9984"} // シスメックス、トヨタ、ソフトバンクG

	for _, symbol := range symbols {
		fmt.Printf("\n銘柄コード: %s\n", symbol)

		// 価格情報を取得（実装されている場合）
		// ここでは簡単な注文テストで価格範囲を確認
		testParams := client.NewOrderParams{
			ZyoutoekiKazeiC:          "1",
			IssueCode:                symbol,
			SizyouC:                  "00",
			BaibaiKubun:              "3", // 買
			Condition:                "0",
			OrderPrice:               "1", // 最小価格でテスト
			OrderSuryou:              "100",
			GenkinShinyouKubun:       "0",
			OrderExpireDay:           "0",
			GyakusasiOrderType:       "0",
			GyakusasiZyouken:         "0",
			GyakusasiPrice:           "*",
			TatebiType:               "*",
			TategyokuZyoutoekiKazeiC: "*",
		}

		res, err := tachibanaClient.NewOrder(ctx, session, testParams)
		if err != nil {
			fmt.Printf("  エラー（価格1円）: %v\n", err)
		} else if res != nil && res.ResultCode != "0" {
			fmt.Printf("  エラー（価格1円）: %s - %s\n", res.ResultCode, res.ResultText)
		}

		// より高い価格でテスト
		testParams.OrderPrice = "10000"
		res, err = tachibanaClient.NewOrder(ctx, session, testParams)
		if err != nil {
			fmt.Printf("  エラー（価格10000円）: %v\n", err)
		} else if res != nil && res.ResultCode != "0" {
			fmt.Printf("  エラー（価格10000円）: %s - %s\n", res.ResultCode, res.ResultText)
		}
	}

	fmt.Println("\n=== 価格確認完了 ===")
}
