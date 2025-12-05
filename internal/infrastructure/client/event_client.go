// internal/infrastructure/client/event_client.go
package client

import (
	"context"
	"net/http"
)

// EventClient is an interface for handling real-time events via WebSocket.
type EventClient interface {
	Connect(ctx context.Context, urlString string, jar http.CookieJar) error
	Close()
	ReadMessages(ctx context.Context) (<-chan map[string]string, <-chan error)
}
