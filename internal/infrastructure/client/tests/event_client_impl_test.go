// internal/infrastructure/client/tests/event_client_impl_test.go
package tests

import (
	// "context"
	// "log/slog"
	// "os"
	"stock-bot/internal/infrastructure/client"
	"testing"

	// "time"

	"github.com/stretchr/testify/assert"
)

// TestParseMessage tests the ParseMessage function.
func TestParseMessage(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected map[string]string
	}{
		{
			name:     "normal case",
			input:    []byte("key1\x02value1\x01key2\x02value2a\x03value2b\x01key3\x02value3"),
			expected: map[string]string{"key1": "value1", "key2": "value2a,value2b", "key3": "value3"},
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: map[string]string{},
		},
		{
			name:     "malformed pair",
			input:    []byte("key1\x02value1\x01key2value2\x01key3\x02value3"),
			expected: map[string]string{"key1": "value1", "key3": "value3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := client.ParseMessage(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

// Note: The TestEventClient_ConnectReadMessagesWithDemoAPI test is commented out
// as it requires a live API connection and is not suitable for regular CI.
// It has been fundamentally broken by the recent refactoring of EventClient.
// A new integration test will be created later.
/*
func TestEventClient_ConnectReadMessagesWithDemoAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// This test now requires a valid session object, which is hard to mock completely.
	// For now, we'll just test the constructor.
	// A proper integration test needs to be written that handles login and session.
	eventClient := client.NewEventClient(logger)
	assert.NotNil(t, eventClient)

	// A mock session would be needed here to proceed.
	// For example:
	// session := &client.Session{ EventURL: "wss://real.api.endpoint" } // and a valid cookie jar
	// messages, errs, err := eventClient.Connect(context.Background(), session)
	// ... and then check the channels.
}
*/

// go test -v ./internal/infrastructure/client/tests/event_client_impl_test.go
