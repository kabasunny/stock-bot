package eventprocessing

import (
	"context"
	"fmt"
	"log/slog"
)

// StatusEventHandlerImpl はステータス通知イベントハンドラーの実装
type StatusEventHandlerImpl struct {
	logger *slog.Logger
}

// NewStatusEventHandler は新しいステータスイベントハンドラーを作成する
func NewStatusEventHandler(logger *slog.Logger) *StatusEventHandlerImpl {
	return &StatusEventHandlerImpl{
		logger: logger,
	}
}

// HandleEvent はイベントハンドラーインターフェースの実装
func (h *StatusEventHandlerImpl) HandleEvent(ctx context.Context, eventType string, data map[string]string) error {
	if eventType != "ST" {
		return fmt.Errorf("unsupported event type for status handler: %s", eventType)
	}

	return h.HandleStatusUpdate(ctx, data)
}

// HandleStatusUpdate はステータス更新を処理する
func (h *StatusEventHandlerImpl) HandleStatusUpdate(ctx context.Context, status map[string]string) error {
	h.logger.Warn("received unhandled Status Notification (ST) event", "data", status)

	// TODO: セッション状態などを確認し、必要に応じて再接続などの処理を行う
	// 例: "session inactive" を検知してエージェントを安全に停止させるなど

	return nil
}
