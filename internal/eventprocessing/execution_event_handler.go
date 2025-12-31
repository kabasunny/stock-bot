package eventprocessing

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/model"
	"stock-bot/internal/app"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

// ExecutionEventHandlerImpl は約定通知イベントハンドラーの実装
type ExecutionEventHandlerImpl struct {
	executionUseCase app.ExecutionUseCase
	logger           *slog.Logger
}

// NewExecutionEventHandler は新しい約定イベントハンドラーを作成する
func NewExecutionEventHandler(executionUseCase app.ExecutionUseCase, logger *slog.Logger) *ExecutionEventHandlerImpl {
	return &ExecutionEventHandlerImpl{
		executionUseCase: executionUseCase,
		logger:           logger,
	}
}

// HandleEvent はイベントハンドラーインターフェースの実装
func (h *ExecutionEventHandlerImpl) HandleEvent(ctx context.Context, eventType string, data map[string]string) error {
	if eventType != "EC" {
		return fmt.Errorf("unsupported event type for execution handler: %s", eventType)
	}

	execution, err := h.parseExecutionFromEventData(data)
	if err != nil {
		return fmt.Errorf("failed to parse execution from event data: %w", err)
	}

	return h.HandleExecution(ctx, execution)
}

// HandleExecution は約定通知を処理し、注文状態を更新する
func (h *ExecutionEventHandlerImpl) HandleExecution(ctx context.Context, execution *model.Execution) error {
	h.logger.Debug("processing execution event", "execution_id", execution.ExecutionID, "order_id", execution.OrderID)

	if err := h.executionUseCase.Execute(ctx, execution); err != nil {
		// "order with ID ... not found" エラーの場合は Warn レベルでログを出力
		if strings.Contains(err.Error(), "order with ID") && strings.Contains(err.Error(), "not found") {
			h.logger.Warn("execution event for non-existent order received", "error", err)
			return nil // エラーとして扱わない
		}

		h.logger.Error("failed to execute execution use case", "execution_id", execution.ExecutionID, "error", err)
		return errors.Wrapf(err, "failed to execute execution use case for execution_id %s", execution.ExecutionID)
	}

	h.logger.Debug("successfully processed execution event", "execution_id", execution.ExecutionID, "order_id", execution.OrderID)
	return nil
}

// parseExecutionFromEventData はイベントデータから約定情報を解析する
func (h *ExecutionEventHandlerImpl) parseExecutionFromEventData(data map[string]string) (*model.Execution, error) {
	execution := &model.Execution{}

	// ExecutionID は約定ごとにユニークなIDが必要だが、ECイベントには直接存在しないため、p_ON (注文ID) と p_ENO (連番) を組み合わせる
	if orderID, ok := data["p_ON"]; ok {
		execution.OrderID = orderID
		if executionNo, ok := data["p_ENO"]; ok {
			execution.ExecutionID = fmt.Sprintf("%s-%s", orderID, executionNo) // 注文ID-約定番号
		} else {
			execution.ExecutionID = orderID // 約定番号がない場合は注文IDのみ
		}
	} else {
		return nil, errors.New("EC event missing p_ON (OrderID)")
	}

	// Symbol
	if val, ok := data["p_IC"]; ok { // p_IC は銘柄コード
		execution.Symbol = val
	} else {
		return nil, errors.New("EC event missing p_IC (Symbol)")
	}

	// TradeType
	if val, ok := data["p_ST"]; ok { // p_ST は売買区分 (1:買, 2:売)
		switch val {
		case "1":
			execution.TradeType = model.TradeTypeBuy
		case "2":
			execution.TradeType = model.TradeTypeSell
		default:
			return nil, errors.Errorf("invalid p_ST (TradeType) in EC event: %s", val)
		}
	} else {
		return nil, errors.New("EC event missing p_ST (TradeType)")
	}

	// Quantity
	if val, ok := data["p_EXSR"]; ok { // p_EXSR は約定数量
		qty, err := parseInt(val)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid p_EXSR (Quantity) in EC event: %s", val)
		}
		execution.Quantity = qty
	} else {
		return nil, errors.New("EC event missing p_EXSR (Quantity)")
	}

	// Price
	if val, ok := data["p_EXPR"]; ok { // p_EXPR は約定単価
		price, err := parseFloat(val)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid p_EXPR (Price) in EC event: %s", val)
		}
		execution.Price = price
	} else {
		return nil, errors.New("EC event missing p_EXPR (Price)")
	}

	// ExecutedAt
	if val, ok := data["p_EXDT"]; ok { // p_EXDT は約定日時 YYYYMMDDhhmmss
		executedAt, err := parseTime(val)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid p_EXDT (ExecutedAt) in EC event: %s", val)
		}
		execution.ExecutedAt = executedAt
	} else {
		h.logger.Warn("EC event missing p_EXDT (ExecutedAt), using current time")
		execution.ExecutedAt = time.Now()
	}

	// Commission (optional): ECイベントログには見当たらないため、0とする
	execution.Commission = 0

	return execution, nil
}

// parseInt は文字列をintにパースするヘルパー関数
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// parseFloat は文字列をfloat64にパースするヘルパー関数
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// parseTime は文字列をtime.Timeにパースするヘルパー関数
func parseTime(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
		"20060102150405",            // YYYYMMDDhhmmss
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse time string: %s", s)
}
