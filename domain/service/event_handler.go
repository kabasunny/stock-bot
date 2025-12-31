package service

import (
	"context"
	"stock-bot/domain/model"
)

// EventHandler はWebSocketイベントを処理するためのインターフェース
type EventHandler interface {
	// HandleEvent は受信したイベントデータを処理する
	HandleEvent(ctx context.Context, eventType string, data map[string]string) error
}

// ExecutionEventHandler は約定通知イベントを処理するインターフェース
type ExecutionEventHandler interface {
	EventHandler
	// HandleExecution は約定通知を処理し、注文状態を更新する
	HandleExecution(ctx context.Context, execution *model.Execution) error
}

// PriceEventHandler は価格データイベントを処理するインターフェース
type PriceEventHandler interface {
	EventHandler
	// HandlePriceUpdate は価格更新を処理し、状態を更新する
	HandlePriceUpdate(ctx context.Context, symbol string, price float64) error
}

// StatusEventHandler はステータス通知イベントを処理するインターフェース
type StatusEventHandler interface {
	EventHandler
	// HandleStatusUpdate はステータス更新を処理する
	HandleStatusUpdate(ctx context.Context, status map[string]string) error
}

// EventDispatcher はイベントを適切なハンドラーに振り分けるインターフェース
type EventDispatcher interface {
	// RegisterHandler はイベントタイプに対するハンドラーを登録する
	RegisterHandler(eventType string, handler EventHandler)
	// DispatchEvent はイベントを適切なハンドラーに振り分ける
	DispatchEvent(ctx context.Context, eventType string, data map[string]string) error
}
