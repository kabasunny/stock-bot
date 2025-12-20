// internal/infrastructure/client/event_client_impl_test.go
package tests

import (
	"context"
	"net/url"
	"testing"
	"time"

	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"

	"stock-bot/internal/infrastructure/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseMessage tests the parseMessage function.
func TestParseMessage(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected map[string]string
	}{
		{
			name:  "normal case",
			input: []byte("key1\x02value1\x01key2\x02value2a\x03value2b\x01key3\x02value3"),
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2a,value2b",
				"key3": "value3",
			},
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
		{
			name:     "no value separator",
			input:    []byte("key1value1"),
			expected: map[string]string{},
		},
		{
			name:     "trailing separators",
			input:    []byte("key1\x02value1\x01\x01"),
			expected: map[string]string{"key1": "value1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := client.ParseMessage(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

// Helper to construct WebSocket URL for the test
func makeWebSocketURLForTest(t *testing.T, eventURL string, pGyouNo, sIssueCode, sSizyouC string) string {
	t.Helper()
	// Mimic func_make_websocket_url from Python sample
	strURL := eventURL
	strURL += "?"
	strURL += "p_rid=22"              // Fixed value for market data
	strURL += "&" + "p_board_no=1000" // Fixed value
	strURL += "&" + "p_gyou_no=" + pGyouNo
	strURL += "&" + "p_mkt_code=" + sSizyouC
	strURL += "&" + "p_eno=0"            // Event notification number, 0 for all
	strURL += "&" + "p_evt_cmd=ST,KP,FD" // Status, Keep-Alive, Market Data
	strURL += "&" + "p_issue_code=" + sIssueCode
	return strURL
}

func TestEventClient_ConnectReadMessagesWithDemoAPI(t *testing.T) {
	// Skip test if not running with demo API credentials (or if it's explicitly skipped)
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. Create a TachibanaClient test instance
	tc := client.CreateTestClient(t)
	require.NotNil(t, tc, "TachibanaClient should not be nil")

	// 2. Login to the demo API
	loginReq := request_auth.ReqLogin{
		UserId:   tc.GetUserIDForTest(),
		Password: tc.GetPasswordForTest(),
	}
	t.Logf("Logging in with UserID: %s", loginReq.UserId)
	session, err := tc.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err, "Login to demo API should be successful")
	require.NotNil(t, session, "Login session should not be nil")
	require.Equal(t, "0", session.ResultCode, "Login ResultCode should be 0 for success")
	defer func() {
		t.Log("Logging out from demo API.")
		logoutReq := request_auth.ReqLogout{}
		_, logoutErr := tc.LogoutWithPost(context.Background(), session, logoutReq)
		assert.NoError(t, logoutErr, "Logout should be successful")
	}()

	// 3. Get Event URL from login info and convert it to WebSocket scheme
	eventURL := session.EventURL
	require.NotEmpty(t, eventURL, "EventURL should not be empty after login")
	t.Logf("Received HTTP Event URL: %s", eventURL)

	parsedURL, err := url.Parse(eventURL)
	require.NoError(t, err)

	if parsedURL.Scheme == "https" {
		parsedURL.Scheme = "wss"
	} else {
		parsedURL.Scheme = "ws"
	}
	eventWebSocketURL := parsedURL.String()
	require.NotEmpty(t, eventWebSocketURL, "Constructed WebSocket URL should not be empty")
	t.Logf("Converted to WebSocket Event URL: %s", eventWebSocketURL)

	// 4. Construct the full WebSocket connection URL
	testPGyouNo := "1"
	testIssueCode := "8411" // Example: Mizuho FG
	testSizyouC := "00"     // TSE
	fullWSURL := makeWebSocketURLForTest(t, eventWebSocketURL, testPGyouNo, testIssueCode, testSizyouC)
	t.Logf("Full WebSocket Connection URL: %s", fullWSURL)

	// 5. Create EventClient and connect
	eventClient := client.NewEventClient()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // Increased timeout for external API
	defer cancel()

	// Retrieve the cookie jar after login
	jar := session.CookieJar
	require.NotNil(t, jar, "CookieJar should not be nil after login")

	err = eventClient.Connect(ctx, fullWSURL, jar)
	require.NoError(t, err, "Connecting to WebSocket should be successful")
	t.Log("WebSocket connected successfully.")

	// 6. Read messages from WebSocket
	msgCh, errCh := eventClient.ReadMessages(ctx)

	var receivedCount int
	// We expect at least a keep-alive message (KP) or initial price data (FD).
	t.Log("Waiting for messages from WebSocket...")
	for receivedCount < 2 { // Wait for at least 2 messages
		select {
		case msg, ok := <-msgCh:
			if !ok {
				t.Fatalf("Message channel closed unexpectedly. Received %d messages.", receivedCount)
			}
			t.Logf("Received message: %v", msg)
			assert.NotEmpty(t, msg, "Received message map should not be empty")
			receivedCount++
		case err := <-errCh:
			t.Fatalf("Received error from WebSocket: %v", err)
		case <-ctx.Done():
			if receivedCount > 0 {
				t.Logf("Context timed out, but received %d messages successfully. Test passes.", receivedCount)
				return
			}
			t.Fatal("Test timed out waiting for any message from WebSocket.")
		}
	}
	require.Greater(t, receivedCount, 0, "Should have received at least one message")
	t.Logf("Successfully received %d messages.", receivedCount)

	// 7. Close the EventClient
	eventClient.Close()
	t.Log("WebSocket client closed.")
}

// go test -v ./internal/infrastructure/client/tests/event_client_impl_test.go
