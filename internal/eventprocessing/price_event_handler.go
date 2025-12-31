package eventprocessing

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/internal/state"
	"strings"
)

// PriceEventHandlerImpl は価格データイベントハンドラーの実装
type PriceEventHandlerImpl struct {
	state          *state.State
	gyouNoToSymbol map[string]string // 行番号から銘柄コードへのマッピング
	logger         *slog.Logger
}

// NewPriceEventHandler は新しい価格イベントハンドラーを作成する
func NewPriceEventHandler(agentState *state.State, gyouNoToSymbol map[string]string, logger *slog.Logger) *PriceEventHandlerImpl {
	return &PriceEventHandlerImpl{
		state:          agentState,
		gyouNoToSymbol: gyouNoToSymbol,
		logger:         logger,
	}
}

// HandleEvent はイベントハンドラーインターフェースの実装
func (h *PriceEventHandlerImpl) HandleEvent(ctx context.Context, eventType string, data map[string]string) error {
	if eventType != "FD" {
		return fmt.Errorf("unsupported event type for price handler: %s", eventType)
	}

	return h.processPriceData(ctx, data)
}

// HandlePriceUpdate は価格更新を処理し、状態を更新する
func (h *PriceEventHandlerImpl) HandlePriceUpdate(ctx context.Context, symbol string, price float64) error {
	h.state.UpdatePrice(symbol, price)
	h.logger.Debug("updated price from event", "symbol", symbol, "price", price)
	return nil
}

// processPriceData は価格情報（時価配信）イベントを処理する
func (h *PriceEventHandlerImpl) processPriceData(ctx context.Context, data map[string]string) error {
	// FDイベントのデータは p_行番号_項目名 という形式で来る
	// 例: p_1_DPP -> 行番号1の銘柄の現在値
	// 行番号を特定し、gyouNoToSymbolマップから銘柄コードを取得する

	for key, value := range data {
		// p_N_DPP の形式を想定 (Nは行番号)
		if strings.HasPrefix(key, "p_") && strings.HasSuffix(key, "_DPP") {
			parts := strings.Split(key, "_")
			if len(parts) != 3 {
				continue // 予期しない形式
			}
			gyouNo := parts[1] // 行番号文字列

			symbol, ok := h.gyouNoToSymbol[gyouNo]
			if !ok {
				h.logger.Warn("unknown gyouNo in price data", "gyouNo", gyouNo, "key", key)
				continue
			}

			price, err := parseFloat(value)
			if err != nil {
				h.logger.Error("failed to parse price from FD event", "symbol", symbol, "key", key, "value", value, "error", err)
				continue
			}

			if err := h.HandlePriceUpdate(ctx, symbol, price); err != nil {
				h.logger.Error("failed to handle price update", "symbol", symbol, "price", price, "error", err)
				continue
			}
		}
	}

	return nil
}
