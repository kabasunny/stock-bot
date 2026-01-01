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
	fmt.Println("=== 立花証券クライアント直接STOP_LIMIT相当テスト ===")

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

	// 4. 逆指値指値注文の詳細テスト
	testCases := []struct {
		name   string
		params client.NewOrderParams
	}{
		{
			name: "逆指値指値注文 - GyakusasiOrderType=2（通常+逆指値）",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "3632", // 成功した銘柄
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし
				OrderPrice:               "450",  // 通常の指値価格
				OrderSuryou:              "100",  // 100株
				GenkinShinyouKubun:       "0",    // 現物
				OrderExpireDay:           "0",    // 当日限り
				GyakusasiOrderType:       "2",    // 通常＋逆指値
				GyakusasiZyouken:         "460",  // 逆指値条件（460円以上で発動）
				GyakusasiPrice:           "455",  // 逆指値時の価格
				TatebiType:               "*",    // 指定なし
				TategyokuZyoutoekiKazeiC: "*",    // 指定なし
			},
		},
		{
			name: "逆指値指値注文 - GyakusasiOrderType=1（逆指値のみ）+ OrderPrice指定",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "3632", // 成功した銘柄
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし
				OrderPrice:               "455",  // 指値価格を指定
				OrderSuryou:              "100",  // 100株
				GenkinShinyouKubun:       "0",    // 現物
				OrderExpireDay:           "0",    // 当日限り
				GyakusasiOrderType:       "1",    // 逆指値のみ
				GyakusasiZyouken:         "460",  // 逆指値条件（460円以上で発動）
				GyakusasiPrice:           "455",  // 逆指値時の価格
				TatebiType:               "*",    // 指定なし
				TategyokuZyoutoekiKazeiC: "*",    // 指定なし
			},
		},
	}

	// 5. 各テストケースを実行
	for i, tc := range testCases {
		fmt.Printf("\n%d. %s\n", i+1, tc.name)
		fmt.Printf("   パラメータ:\n")
		fmt.Printf("     OrderPrice: %s\n", tc.params.OrderPrice)
		fmt.Printf("     GyakusasiOrderType: %s\n", tc.params.GyakusasiOrderType)
		fmt.Printf("     GyakusasiZyouken: %s\n", tc.params.GyakusasiZyouken)
		fmt.Printf("     GyakusasiPrice: %s\n", tc.params.GyakusasiPrice)

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
