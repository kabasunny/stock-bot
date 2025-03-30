// internal/app/account_usecase.go
package app

import (
	"context"
	"fmt"
	"stock-bot/internal/infrastructure/client"
	"strconv"
)

// BalanceInformation は、口座情報を表す
type BalanceUseCaseImpl struct {
	client *client.TachibanaClient
}

func (uc *BalanceUseCaseImpl) CanEntry(ctx context.Context, issueCode string) (bool, float64, error) {
	// 1. API から必要な情報を取得
	genbutuKabuList, err := uc.client.GetGenbutuKabuList(ctx)
	if err != nil {
		return false, 0, fmt.Errorf("GetGenbutuKabuList failed: %w", err)
	}

	shinyouTategyokuList, err := uc.client.GetShinyouTategyokuList(ctx)
	if err != nil {
		return false, 0, fmt.Errorf("GetShinyouTategyokuList failed: %w", err)
	}

	zanKaiSummary, err := uc.client.GetZanKaiSummary(ctx)
	if err != nil {
		return false, 0, fmt.Errorf("GetZanKaiGenbutuKaitukeSyousai failed: %w", err)
	}
	// 2. 取得した情報を基に口座情報を計算
	//  genbutuKabuList 及び shinyouTategyokuList に issueCode が含まれなければ、エントリ可
	// isHolding を初期化
	isHolding := false

	// genbutuKabuList に issueCode が含まれているかチェック
	for _, item := range genbutuKabuList.GenbutuKabuList {
		if item.UriOrderIssueCode == issueCode {
			isHolding = true
			break
		}
	}

	// shinyouTategyokuList に issueCode が含まれているかチェック
	if !isHolding { // すでに含まれている場合はチェック不要
		for _, item := range shinyouTategyokuList.SinyouTategyokuList {
			if item.OrderIssueCode == issueCode {
				isHolding = true
				break
			}
		}
	}
	// 総資金の計算 総資金 = 預り金 + 信用評価損益合計 + 現物評価額合計
	// stringをfloat64に変換
	zanKaiSummaryValue, err1 := strconv.ParseFloat(zanKaiSummary.GenbutuKabuKaituke, 64)
	if err1 != nil {
		fmt.Println("変換エラー:", err1)
		return false, 0, err1
	}

	shinyouTategyokuListValue, err2 := strconv.ParseFloat(shinyouTategyokuList.TotalHyoukaSonekiGoukei, 64)
	if err2 != nil {
		fmt.Println("変換エラー:", err2)
		return false, 0, err1
	}

	genbutuKabuListValue, err3 := strconv.ParseFloat(genbutuKabuList.TotalGaisanHyoukagakuGoukei, 64)
	if err3 != nil {
		fmt.Println("変換エラー:", err3)
		return false, 0, err1
	}

	// 加算処理
	totalAssets := zanKaiSummaryValue + shinyouTategyokuListValue + genbutuKabuListValue
	fmt.Println("総資産:", totalAssets)

	return !isHolding, totalAssets, nil
}
