package app

import (
	"context"
	"fmt"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/internal/infrastructure/client"
	// "stock-bot/internal/infrastructure/client/dto/order/request" // Removed as request.ReqNewOrder is no longer directly used
)

// OrderUseCaseの実装
type OrderUseCaseImpl struct {
	orderClient client.OrderClient
	orderRepo   repository.OrderRepository
	// secondPassword string // Removed
}

// NewOrderUseCaseImpl はOrderUseCaseImplの新しいインスタンスを生成します
func NewOrderUseCaseImpl(orderClient client.OrderClient, orderRepo repository.OrderRepository) OrderUseCase {
	return &OrderUseCaseImpl{
		orderClient: orderClient,
		orderRepo:   orderRepo,
		// secondPassword: secondPassword, // Removed
	}
}

// ExecuteOrder は注文を実行します
func (uc *OrderUseCaseImpl) ExecuteOrder(ctx context.Context, session *client.Session, params OrderParams) (*model.Order, error) {
	// 1. 外部APIへのリクエストDTOに変換
	// TradeType のマッピング
	var baibaiKubun string
	switch params.TradeType {
	case model.TradeTypeBuy:
		baibaiKubun = "3" // 買
	case model.TradeTypeSell:
		baibaiKubun = "1" // 売
	default:
		return nil, fmt.Errorf("invalid trade type: %s", params.TradeType)
	}

	// OrderType のマッピング
	var orderPrice string
	var condition string
	switch params.OrderType {
	case model.OrderTypeMarket:
		orderPrice = "0" // 成行
		condition = "0"  // 指定なし
	case model.OrderTypeLimit:
		orderPrice = fmt.Sprintf("%.1f", params.Price) // 指値価格
		condition = "0"                                // 指定なし
	default:
		return nil, fmt.Errorf("invalid order type: %s", params.OrderType)
	}

	// GenkinShinyouKubun のマッピング
	var genkinShinyouKubun string
	switch params.PositionAccountType {
	case model.PositionAccountTypeCash:
		genkinShinyouKubun = "0" // 現物 (Cash)
	case model.PositionAccountTypeMarginNew:
		genkinShinyouKubun = "2" // 信用新規 (Margin New)
	case model.PositionAccountTypeMarginRepay:
		genkinShinyouKubun = "4" // 信用返済 (Margin Repay)
	default:
		// デフォルトは現物とするが、不明な場合は上位層でログ出力済みの前提
		genkinShinyouKubun = "0"
	}

	req := client.NewOrderParams{ // Changed to client.NewOrderParams
		// SecondPassword:           uc.secondPassword, // Removed
		ZyoutoekiKazeiC:          "1", // 特定口座
		IssueCode:                params.Symbol,
		SizyouC:                  "00", // 東証
		BaibaiKubun:              baibaiKubun,
		Condition:                condition,
		OrderPrice:               orderPrice,
		OrderSuryou:              fmt.Sprintf("%d", params.Quantity),
		GenkinShinyouKubun:       genkinShinyouKubun, // ここでマッピングした値を使用
		OrderExpireDay:           "0",                // 当日
		GyakusasiOrderType:       "0",                // 通常注文
		GyakusasiZyouken:         "0",                // 指定なし
		GyakusasiPrice:           "*",                // 指定なし
		TatebiType:               "*",                // 指定なし
		TategyokuZyoutoekiKazeiC: "*",                // 指定なし
	}

	// 2. 外部API（証券会社）を呼び出す
	res, err := uc.orderClient.NewOrder(ctx, session, req) // No change in call, but req type changed
	if err != nil {
		return nil, fmt.Errorf("failed to execute order via client: %w", err)
	}
	if res.ResultCode != "0" {
		return nil, fmt.Errorf("order failed with result code %s: %s", res.ResultCode, res.ResultText)
	}

	// 3. 結果をドメインモデルに変換
	newOrder := &model.Order{
		OrderID:             res.OrderNumber,
		Symbol:              params.Symbol,
		TradeType:           params.TradeType,
		OrderType:           params.OrderType,
		Quantity:            int(params.Quantity),
		Price:               params.Price,
		OrderStatus:         model.OrderStatusNew,
		IsMargin:            (params.PositionAccountType == model.PositionAccountTypeMarginNew || params.PositionAccountType == model.PositionAccountTypeMarginRepay), // PositionAccountTypeからIsMarginを設定
		PositionAccountType: params.PositionAccountType,                                      // ここでPositionAccountTypeを設定
	}

	// 4. リポジトリで永続化
	if err := uc.orderRepo.Save(ctx, newOrder); err != nil {
		// APIは成功したがDB保存に失敗した場合。
		// ここではエラーを返すだけだが、実際にはリトライや補正処理を検討するべき。
		return nil, fmt.Errorf("failed to save order to repository: %w", err)
	}

	return newOrder, nil
}
