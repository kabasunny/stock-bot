package main

import (
	"context"
	"fmt"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/auth/request"
	"time"
)

func main() {
	fmt.Println("=== 全注文タイプ立花クライアント単体テスト ===")

	// 1. クライアント作成とログイン
	c := client.CreateTestClient(nil)
	ctx := context.Background()

	loginReq := request.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}

	session, err := c.LoginWithPost(ctx, loginReq)
	if err != nil {
		fmt.Printf("ログインエラー: %v\n", err)
		return
	}
	fmt.Printf("ログイン成功\n")

	// 2. テストケース定義
	testCases := []struct {
		name   string
		params client.NewOrderParams
	}{
		{
			name: "成行注文（現物買い）",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "6658", // シスメックス
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし (成行)
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
			name: "成行注文（信用新規買い）",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "6658", // シスメックス
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし (成行)
				OrderPrice:               "0",    // 成行
				OrderSuryou:              "100",  // 100株
				GenkinShinyouKubun:       "2",    // 信用新規（制度信用6ヶ月）
				OrderExpireDay:           "0",    // 当日限り
				GyakusasiOrderType:       "0",    // 通常注文
				GyakusasiZyouken:         "0",    // 指定なし
				GyakusasiPrice:           "*",    // 指定なし
				TatebiType:               "*",    // 指定なし
				TategyokuZyoutoekiKazeiC: "*",    // 指定なし
			},
		},
		{
			name: "指値注文（現物買い）",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "6658", // シスメックス
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "3",    // 買
				Condition:                "0",    // 指定なし
				OrderPrice:               "8000", // 指値（低めに設定）
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
			name: "逆指値注文（現物売り）",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",    // 特定口座
				IssueCode:                "6658", // シスメックス
				SizyouC:                  "00",   // 東証
				BaibaiKubun:              "1",    // 売
				Condition:                "0",    // 指定なし
				OrderPrice:               "*",    // 指定なし (逆指値の場合)
				OrderSuryou:              "100",  // 100株
				GenkinShinyouKubun:       "0",    // 現物
				OrderExpireDay:           "0",    // 当日限り
				GyakusasiOrderType:       "1",    // 逆指値
				GyakusasiZyouken:         "8000", // 逆指値条件（8000円以下で売り）
				GyakusasiPrice:           "0",    // 逆指値値段（成行）
				TatebiType:               "*",    // 指定なし
				TategyokuZyoutoekiKazeiC: "*",    // 指定なし
			},
		},
		{
			name: "逆指値指値注文（現物買い）",
			params: client.NewOrderParams{
				ZyoutoekiKazeiC:          "1",     // 特定口座
				IssueCode:                "6658",  // シスメックス
				SizyouC:                  "00",    // 東証
				BaibaiKubun:              "3",     // 買
				Condition:                "0",     // 指定なし
				OrderPrice:               "12000", // 指値
				OrderSuryou:              "100",   // 100株
				GenkinShinyouKubun:       "0",     // 現物
				OrderExpireDay:           "0",     // 当日限り
				GyakusasiOrderType:       "2",     // 逆指値指値
				GyakusasiZyouken:         "11000", // 逆指値条件（11000円以上で発動）
				GyakusasiPrice:           "12000", // 逆指値値段（指値）
				TatebiType:               "*",     // 指定なし
				TategyokuZyoutoekiKazeiC: "*",     // 指定なし
			},
		},
	}

	// 3. 各テストケースを実行
	for i, tc := range testCases {
		fmt.Printf("\n%d. %s\n", i+1, tc.name)
		fmt.Printf("   パラメータ: %+v\n", tc.params)

		res, err := c.NewOrder(ctx, session, tc.params)
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
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n=== テスト完了 ===")
}
