// internal/app/account_usecase.go
package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"stock-bot/internal/infrastructure/client"
	_ "stock-bot/internal/logger"
)

// BalanceUseCaseImpl は、口座情報を表す
type BalanceUseCaseImpl struct {
	client client.BalanceClient // BalanceClientインターフェースを使用
}

// NewBalanceUseCaseImpl は、BalanceUseCaseImplのコンストラクタ
func NewBalanceUseCaseImpl(client client.BalanceClient) BalanceUseCase {
	return &BalanceUseCaseImpl{
		client: client,
	}
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
		return false, 0, fmt.Errorf("GetZanKaiSummary failed: %w", err)
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
		slog.Error("failed to parse float for zanKaiSummaryValue", slog.Any("error", err1), slog.String("value", zanKaiSummary.GenbutuKabuKaituke))
		return false, 0, err1
	}

	shinyouTategyokuListValue, err2 := strconv.ParseFloat(shinyouTategyokuList.TotalHyoukaSonekiGoukei, 64)
	if err2 != nil {
		slog.Error("failed to parse float for shinyouTategyokuListValue", slog.Any("error", err2), slog.String("value", shinyouTategyokuList.TotalHyoukaSonekiGoukei))
		return false, 0, err1
	}

	genbutuKabuListValue, err3 := strconv.ParseFloat(genbutuKabuList.TotalGaisanHyoukagakuGoukei, 64)
	if err3 != nil {
		slog.Error("failed to parse float for genbutuKabuListValue", slog.Any("error", err3), slog.String("value", genbutuKabuList.TotalGaisanHyoukagakuGoukei))
		return false, 0, err1
	}

	// 加算処理
	totalAssets := zanKaiSummaryValue + shinyouTategyokuListValue + genbutuKabuListValue
	slog.Info("calculated total assets", slog.Float64("totalAssets", totalAssets))

	return !isHolding, totalAssets, nil
}
