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
	"time"
)

func main() {
	fmt.Println("=== 逆指値指値注文テスト ===")

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

	// 4. 段階的テスト（成行→指値→逆指値→逆指値指値）
	testCases := []struct {
		name   string
		params client.NewOrderParams
	}{
		{
			name: "1. 成行注文（ベースライン確認）",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "6658", // シスメックス
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし
				OrderPrice:               "0",    // 成行
				OrderSuryou:              "100",  // 100株
				GenkinShinyouKubun:       "0",    // 現物
				OrderExpireDay:           "0",    // 当日限り
				GyakusasiOrderType:       "0",    // 通常注文
				GyakusasiZyouken:         "0",    // 指定なし
				GyakusasiPrice:           "*",    // 指定なし
				TatebiType:               "*",    // 指定なし
				TategyokuZyoutoekiKazeiC: "*",    // 指定なし
			},
		},
		{
			name: "2. 逆指値注文（成行）- 公式例準拠",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "3632", // 公式例の銘柄
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし
				OrderPrice:               "*",    // 指定なし（逆指値の場合）
				OrderSuryou:              "100",  // 100株
				GenkinShinyouKubun:       "0",    // 現物
				OrderExpireDay:           "0",    // 当日限り
				GyakusasiOrderType:       "1",    // 逆指値
				GyakusasiZyouken:         "460",  // 公式例の条件
				GyakusasiPrice:           "455",  // 公式例の価格
				TatebiType:               "*",    // 指定なし
				TategyokuZyoutoekiKazeiC: "*",    // 指定なし
			},
		},
		{
			name: "3. 通常＋逆指値注文（指値＋逆指値）- 公式例準拠",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "3668", // 公式例の銘柄
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし
				OrderPrice:               "970",  // 公式例の通常指値
				OrderSuryou:              "100",  // 100株
				GenkinShinyouKubun:       "0",    // 現物
				OrderExpireDay:           "0",    // 当日限り
				GyakusasiOrderType:       "2",    // 通常＋逆指値
				GyakusasiZyouken:         "974",  // 公式例の条件
				GyakusasiPrice:           "972",  // 公式例の価格
				TatebiType:               "*",    // 指定なし
				TategyokuZyoutoekiKazeiC: "*",    // 指定なし
			},
		},
	}

	// 5. 各テストケースを実行
	for i, tc := range testCases {
		fmt.Printf("\n%d. %s\n", i+1, tc.name)
		fmt.Printf("   パラメータ:\n")
		fmt.Printf("     銘柄コード: %s\n", tc.params.IssueCode)
		fmt.Printf("     売買区分: %s (3=買)\n", tc.params.BaibaiKubun)
		fmt.Printf("     注文価格: %s\n", tc.params.OrderPrice)
		fmt.Printf("     逆指値種別: %s\n", tc.params.GyakusasiOrderType)
		fmt.Printf("     逆指値条件: %s円\n", tc.params.GyakusasiZyouken)
		fmt.Printf("     逆指値価格: %s円\n", tc.params.GyakusasiPrice)

		res, err := tachibanaClient.NewOrder(ctx, session, tc.params)
		if err != nil {
			fmt.Printf("   ❌ エラー: %v\n", err)
		} else if res != nil {
			if res.ResultCode == "0" {
				fmt.Printf("   ✅ 成功 - 注文番号: %s, 営業日: %s\n", res.OrderNumber, res.EigyouDay)
			} else {
				fmt.Printf("   ❌ 失敗 - エラーコード: %s, メッセージ: %s\n", res.ResultCode, res.ResultText)
			}
		}

		// API制限を考慮して待機
		time.Sleep(2 * time.Second)
	}

	fmt.Println("\n=== テスト完了 ===")
}
