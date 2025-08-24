// internal/app/account_usecase.go
package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"stock-bot/domain/model"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/client/dto/balance/response"
	_ "stock-bot/internal/logger"

	"golang.org/x/sync/errgroup"
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
	zanKaiSummaryValue, err1 := parseFloat(zanKaiSummary.GenbutuKabuKaituke)
	if err1 != nil {
		slog.Error("failed to parse float for zanKaiSummaryValue", slog.Any("error", err1), slog.String("value", zanKaiSummary.GenbutuKabuKaituke))
		return false, 0, err1
	}

	shinyouTategyokuListValue, err2 := parseFloat(shinyouTategyokuList.TotalHyoukaSonekiGoukei)
	if err2 != nil {
		slog.Error("failed to parse float for shinyouTategyokuListValue", slog.Any("error", err2), slog.String("value", shinyouTategyokuList.TotalHyoukaSonekiGoukei))
		return false, 0, err1
	}

	genbutuKabuListValue, err3 := parseFloat(genbutuKabuList.TotalGaisanHyoukagakuGoukei)
	if err3 != nil {
		slog.Error("failed to parse float for genbutuKabuListValue", slog.Any("error", err3), slog.String("value", genbutuKabuList.TotalGaisanHyoukagakuGoukei))
		return false, 0, err1
	}

	// 加算処理
	totalAssets := zanKaiSummaryValue + shinyouTategyokuListValue + genbutuKabuListValue
	slog.Info("calculated total assets", slog.Float64("totalAssets", totalAssets))

	return !isHolding, totalAssets, nil
}

func (uc *BalanceUseCaseImpl) GetSummary(ctx context.Context) (*model.BalanceSummary, error) {
	var (
		zanKaiSummary      *response.ResZanKaiSummary
		shinyouTategyokuList *response.ResShinyouTategyokuList
		genbutuKabuList    *response.ResGenbutuKabuList
	)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var err error
		zanKaiSummary, err = uc.client.GetZanKaiSummary(ctx)
		return err
	})

	eg.Go(func() error {
		var err error
		shinyouTategyokuList, err = uc.client.GetShinyouTategyokuList(ctx)
		return err
	})

	eg.Go(func() error {
		var err error
		genbutuKabuList, err = uc.client.GetGenbutuKabuList(ctx)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get balance summary data: %w", err)
	}

	// Parse values
	// Note: TotalAssets calculation is based on the existing logic in CanEntry.
	assetFromSummary, err := parseFloat(zanKaiSummary.GenbutuKabuKaituke)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GenbutuKabuKaituke for total assets: %w", err)
	}
	marginPL, err := parseFloat(shinyouTategyokuList.TotalHyoukaSonekiGoukei)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TotalHyoukaSonekiGoukei: %w", err)
	}
	stockValue, err := parseFloat(genbutuKabuList.TotalGaisanHyoukagakuGoukei)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TotalGaisanHyoukagakuGoukei: %w", err)
	}
	cashBuyingPower, err := parseFloat(zanKaiSummary.GenbutuKabuKaituke)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GenbutuKabuKaituke for buying power: %w", err)
	}
	marginBuyingPower, err := parseFloat(zanKaiSummary.SinyouSinkidate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SinyouSinkidate: %w", err)
	}
	withdrawalPossibleAmount, err := parseFloat(zanKaiSummary.Syukkin)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Syukkin: %w", err)
	}
	marginRate, err := parseFloat(zanKaiSummary.HosyouKinritu)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HosyouKinritu: %w", err)
	}

	updatedAt, err := time.Parse("200601021504", zanKaiSummary.UpdateDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UpdateDate: %w", err)
	}

	// Create summary
	summary := &model.BalanceSummary{
		TotalAssets:              assetFromSummary + marginPL + stockValue,
		CashBuyingPower:          cashBuyingPower,
		MarginBuyingPower:        marginBuyingPower,
		WithdrawalPossibleAmount: withdrawalPossibleAmount,
		MarginRate:               marginRate,
		UpdatedAt:                updatedAt,
	}

	return summary, nil
}

// parseFloat converts a string to a float64, logging an error on failure.
func parseFloat(val string) (float64, error) {
	if val == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		slog.Error("failed to parse float", slog.Any("error", err), slog.String("value", val))
		return 0, err
	}
	return f, nil
}
