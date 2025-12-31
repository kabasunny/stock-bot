package client

import (
	"bytes"
	"context"
	"fmt" // 追加
	"log/slog"
	"net/http"
	"net/url" // 追加
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/websocket"
)

// eventClient implements the EventClient interface.
type eventClient struct {
	conn   *websocket.Conn
	logger *slog.Logger
	mu     sync.Mutex
}

// NewEventClient creates a new EventClient.
func NewEventClient(logger *slog.Logger) EventClient {
	return &eventClient{
		logger: logger.WithGroup("event_client"),
	}
}

// Connect establishes a WebSocket connection and starts a message reading goroutine.
func (c *eventClient) Connect(ctx context.Context, session *Session, symbols []string) (<-chan []byte, <-chan error, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil, nil, errors.New("event client is already connected")
	}

	if session.EventURL == "" {
		return nil, nil, errors.New("WebSocket EventURL is not available in the session")
	}

	// Convert http/https to ws/wss
	wsURL := strings.Replace(session.EventURL, "https://", "wss://", 1)
	wsURL = strings.Replace(wsURL, "http://", "ws://", 1)

	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		c.logger.Error("Failed to parse WebSocket URL", "error", err, "url", wsURL)
		return nil, nil, errors.Wrap(err, "failed to parse WebSocket URL")
	}

	// クエリパラメータを追加
	query := parsedURL.Query()
	
	// Pythonサンプルコードで必須とされているパラメータを追加
	query.Set("p_rid", "22")
	query.Set("p_board_no", "1000")
	query.Set("p_eno", "0") // 配信開始したいイベント通知番号、0なら全て

	// p_evt_cmd はST, KP, FD, ECを追加。必要に応じてNS, RR, SS, USも追加可能
	// ST: エラーステータス情報配信, KP: キープアライブ情報配信, FD: 時価情報配信, EC: 注文約定通知イベント配信
	query.Set("p_evt_cmd", "ST,KP,FD,EC") 

	if len(symbols) > 0 {
		// 銘柄コードを追加
		query.Set("p_issue_code", strings.Join(symbols, ",")) 

		// 行番号 (p_gyou_no) を生成
		var gyouNos []string
		for i := 1; i <= len(symbols); i++ {
			gyouNos = append(gyouNos, fmt.Sprintf("%d", i))
		}
		query.Set("p_gyou_no", strings.Join(gyouNos, ","))

		// 市場コード (p_mkt_code) を生成
		var mktCodes []string
		for i := 0; i < len(symbols); i++ {
			mktCodes = append(mktCodes, "00") // 00:東証
		}
		query.Set("p_mkt_code", strings.Join(mktCodes, ","))
	}
	parsedURL.RawQuery = query.Encode()

	finalWSURL := parsedURL.String()

	// Prepare request headers, especially the Cookie
	header := http.Header{}
	if session.CookieJar != nil {
		eventURL, err := url.Parse(session.EventURL)
		if err != nil {
			c.logger.Error("Failed to parse EventURL for cookie retrieval", "error", err, "eventURL", session.EventURL)
			// エラーを返すか、処理を続行するかは設計によるが、ここではエラーを返す
			return nil, nil, errors.Wrap(err, "failed to parse EventURL for cookie retrieval")
		}
		cookies := session.CookieJar.Cookies(eventURL)
		for _, cookie := range cookies {
			header.Add("Cookie", cookie.String())
		}
	}
	c.logger.Info("Connecting to WebSocket", "url", finalWSURL)

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, finalWSURL, header)
	if err != nil {
		c.logger.Error("Failed to connect to WebSocket", "error", err)
		return nil, nil, errors.Wrap(err, "failed to dial websocket")
	}
	c.conn = conn
	c.logger.Info("Successfully connected to WebSocket")

	messages := make(chan []byte)
	errs := make(chan error, 1)

	// Start a goroutine to read messages
	go func() {
		defer func() {
			close(messages)
			close(errs)
			c.Close() // Ensure connection is closed when loop exits
		}()

		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				messageType, message, err := c.conn.ReadMessage()
				if err != nil {
					// Check if it's a clean close
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						c.logger.Error("WebSocket read error", "error", err)
						errs <- err
					} else {
						c.logger.Info("WebSocket closed cleanly", "error", err)
					}
					return
				}

				if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
					// Log the raw message for debugging and analysis
					c.logger.Info("Received raw WebSocket message", "message", string(message))
					messages <- message
				}
			}
		}
	}()

	return messages, errs, nil
}

// Close terminates the WebSocket connection.
func (c *eventClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.logger.Info("Closing WebSocket connection")
		// Send a clean close message to the server
		_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = c.conn.Close()
		c.conn = nil
	}
}

// ParseMessage はWebSocketから受信した独自形式のメッセージをパースする
func ParseMessage(raw []byte) map[string]string {
	result := make(map[string]string)
	if len(raw) == 0 {
		return result
	}

	records := bytes.Split(raw, []byte{'\x01'})
	for _, record := range records {
		if len(record) == 0 {
			continue
		}
		parts := bytes.SplitN(record, []byte{'\x02'}, 2)
		if len(parts) != 2 {
			continue
		}
		key := string(parts[0])
		// 値 부분은 `\x03`으로 더 분할될 수 있다
		valueParts := bytes.Split(parts[1], []byte{'\x03'})
		var valueStrings []string
		for _, v := range valueParts {
			valueStrings = append(valueStrings, string(v))
		}
		result[key] = strings.Join(valueStrings, ",")
	}
	return result
}