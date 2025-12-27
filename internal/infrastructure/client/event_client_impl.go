package client

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
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
func (c *eventClient) Connect(ctx context.Context, session *Session) (<-chan []byte, <-chan error, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil, nil, errors.New("event client is already connected")
	}

	if session.EventURL == "" {
		return nil, nil, errors.New("WebSocket EventURL is not available in the session")
	}

	// Prepare request headers, especially the Cookie
	header := http.Header{}
	if session.CookieJar != nil {
		cookies := session.CookieJar.Cookies(nil) // URL is nil, gets all cookies
		for _, cookie := range cookies {
			header.Add("Cookie", cookie.String())
		}
	}
	c.logger.Info("Connecting to WebSocket", "url", session.EventURL)

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, session.EventURL, header)
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