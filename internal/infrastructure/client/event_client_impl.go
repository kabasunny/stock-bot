// internal/infrastructure/client/event_client_impl.go
package client

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/websocket"
)

var (
	sepA = []byte("\x01") // Item separator
	sepB = []byte("\x02") // Key-value separator
	sepC = []byte("\x03") // Value-value separator
)

type eventClientImpl struct {
	conn     *websocket.Conn
	mu       sync.Mutex
	isClosed bool
}

// NewEventClient creates a new EventClient.
func NewEventClient() EventClient {
	return &eventClientImpl{}
}

func (c *eventClientImpl) Connect(ctx context.Context, urlString string, jar http.CookieJar) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return errors.New("already connected")
	}

	header := http.Header{}
	u, err := url.Parse(urlString)
	if err != nil {
		return errors.Wrap(err, "failed to parse websocket url")
	}
	originURL := *u
	if originURL.Scheme == "wss" {
		originURL.Scheme = "https"
	} else {
		originURL.Scheme = "http"
	}
	origin := originURL.Scheme + "://" + originURL.Host
	header.Set("Origin", origin)

	// Use a custom dialer to specify the subprotocol and cookie jar.
	dialer := websocket.Dialer{
		Subprotocols: []string{"e-api-stream"},
		Jar:          jar,
	}

	conn, _, err := dialer.DialContext(ctx, urlString, header)
	if err != nil {
		return errors.Wrap(err, "failed to connect to websocket")
	}
	log.Println("WebSocket connected")

	c.conn = conn
	c.isClosed = false
	return nil
}

func (c *eventClientImpl) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed || c.conn == nil {
		return
	}

	log.Println("Closing WebSocket connection")
	c.isClosed = true
	// Send close message and ignore error, as the connection might be already gone.
	_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.conn.Close()
}

// ParseMessage parses the custom message format from the WebSocket server.
// The format uses ^A, ^B, and ^C as delimiters.
// 項目A1^B値B1^A項目A2^B値B21^CB22^CB23^A...
func ParseMessage(msg []byte) map[string]string {
	result := make(map[string]string)
	pairs := bytes.Split(msg, sepA)
	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}
		kv := bytes.SplitN(pair, sepB, 2)
		if len(kv) == 2 {
			key := string(kv[0])
			// Values separated by ^C are joined with a comma.
			value := bytes.ReplaceAll(kv[1], sepC, []byte(","))
			result[key] = string(value)
		}
	}
	return result
}

func (c *eventClientImpl) ReadMessages(ctx context.Context) (<-chan map[string]string, <-chan error) {
	msgCh := make(chan map[string]string)
	errCh := make(chan error, 1)

	go func() {
		defer close(msgCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				log.Println("Context done, stopping ReadMessages.")
				return
			default:
				c.mu.Lock()
				if c.isClosed {
					c.mu.Unlock()
					return
				}
				c.mu.Unlock()

				_, message, err := c.conn.ReadMessage()
				if err != nil {
					if !c.isClosed {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
							log.Printf("Unexpeced close error: %v", err)
							errCh <- errors.Wrap(err, "websocket read error")
						}
					}
					return
				}

				if len(message) > 0 {
					parsedMsg := ParseMessage(message)
					if len(parsedMsg) > 0 {
						msgCh <- parsedMsg
					}
				}
			}
		}
	}()

	return msgCh, errCh
}
