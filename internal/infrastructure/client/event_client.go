package client

import (
	"context"
)

// EventClient is an interface for handling real-time events via WebSocket.
type EventClient interface {
	// Connect establishes a WebSocket connection using the provided session information.
	// It returns a read-only channel for incoming messages, a read-only channel for errors,
	// and an error if the initial connection fails.
	Connect(ctx context.Context, session *Session) (<-chan []byte, <-chan error, error)

	// Close terminates the WebSocket connection.
	Close()
}
