package eventprocessing

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
)

// WebSocketEventService はWebSocketイベントの監視と処理を担当するサービス
type WebSocketEventService struct {
	eventClient     client.EventClient
	eventDispatcher service.EventDispatcher
	logger          *slog.Logger
}

// NewWebSocketEventService は新しいWebSocketイベントサービスを作成する
func NewWebSocketEventService(
	eventClient client.EventClient,
	eventDispatcher service.EventDispatcher,
	logger *slog.Logger,
) *WebSocketEventService {
	return &WebSocketEventService{
		eventClient:     eventClient,
		eventDispatcher: eventDispatcher,
		logger:          logger,
	}
}

// StartEventWatcher はWebSocketイベントの監視を開始する
func (s *WebSocketEventService) StartEventWatcher(ctx context.Context, session *client.Session, targetSymbols []string) error {
	s.logger.Info("starting WebSocket event watcher...")

	if len(targetSymbols) == 0 {
		s.logger.Warn("no target symbols defined for WebSocket subscription, connecting without symbols")
	}

	messages, errs, err := s.eventClient.Connect(ctx, session, targetSymbols)
	if err != nil {
		return fmt.Errorf("failed to connect to event stream: %w", err)
	}

	s.logger.Info("WebSocket event watcher connected")

	go s.processEvents(ctx, messages, errs)

	return nil
}

// processEvents はWebSocketイベントを処理する
func (s *WebSocketEventService) processEvents(ctx context.Context, messages <-chan []byte, errs <-chan error) {
	defer s.eventClient.Close()

	for {
		select {
		case msgBytes, ok := <-messages:
			if !ok {
				s.logger.Info("message channel closed, stopping event processing")
				return
			}

			if err := s.handleMessage(ctx, msgBytes); err != nil {
				s.logger.Error("failed to handle WebSocket message", "error", err)
			}

		case err, ok := <-errs:
			if !ok {
				s.logger.Info("error channel closed, stopping event processing")
				return
			}

			s.logger.Error("received error from event stream", "error", err)
			// TODO: エラー内容に応じた再接続処理などを検討
			return

		case <-ctx.Done():
			s.logger.Info("context done, stopping event processing")
			return
		}
	}
}

// handleMessage はWebSocketメッセージを解析し、適切なハンドラーに振り分ける
func (s *WebSocketEventService) handleMessage(ctx context.Context, msgBytes []byte) error {
	// 1. パース処理
	parsedMsg, err := s.parseEventMessage(msgBytes)
	if err != nil {
		return fmt.Errorf("failed to parse websocket event: %w", err)
	}

	// 2. イベント種別を取得
	cmd, ok := parsedMsg["p_cmd"]
	if !ok {
		s.logger.Warn("p_cmd not found in websocket event", "message", parsedMsg)
		return nil
	}

	s.logger.Info("received websocket event", "command", cmd)

	// 3. イベントディスパッチャーに振り分け
	switch cmd {
	case "FD": // 時価配信データ (Feed Data)
		return s.eventDispatcher.DispatchEvent(ctx, "FD", parsedMsg)
	case "ST": // ステータス通知 (Status)
		return s.eventDispatcher.DispatchEvent(ctx, "ST", parsedMsg)
	case "EC": // 約定通知 (Execution)
		return s.eventDispatcher.DispatchEvent(ctx, "EC", parsedMsg)
	case "KP": // キープアライブ (Keep Alive)
		s.logger.Debug("received keep alive event", "details", parsedMsg)
		return nil
	default:
		s.logger.Warn("unhandled websocket event command", "command", cmd, "details", parsedMsg)
		return nil
	}
}

// parseEventMessage はWebSocketのカスタムフォーマットメッセージをパースする
func (s *WebSocketEventService) parseEventMessage(msg []byte) (map[string]string, error) {
	result := make(map[string]string)

	// メッセージが空、または改行コードのみの場合を除外
	if len(bytes.TrimSpace(msg)) == 0 {
		return result, nil
	}

	pairs := bytes.Split(msg, []byte{0x01}) // ^A で分割
	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}

		kv := bytes.SplitN(pair, []byte{0x02}, 2) // ^B でキーと値に分割
		if len(kv) != 2 {
			// キーだけのペア（例: `p_no^B1`の後ろに`^A`がない場合など）も許容する
			if len(kv) == 1 && len(kv[0]) > 0 {
				result[string(kv[0])] = ""
				continue
			}
			// 不正な形式のペアは無視する
			s.logger.Warn("invalid key-value pair format in websocket message", "pair", string(pair))
			continue
		}

		key := string(kv[0])
		value := string(kv[1])
		result[key] = value
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("message parsing resulted in no key-value pairs: %s", string(msg))
	}

	return result, nil
}
