package tradeservice

import (
	"log/slog"
	"stock-bot/domain/service"
	"stock-bot/internal/infrastructure/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoaTradeService_GetSession(t *testing.T) {
	// Setup
	session := &client.Session{}
	logger := slog.Default()

	service := NewGoaTradeService(
		nil, // balanceClient not needed for this test
		nil, // orderClient not needed for this test
		nil, // priceClient not needed for this test
		nil, // orderRepo not needed for this test
		session,
		logger,
	)

	// Execute
	result := service.GetSession()

	// Assert
	assert.Equal(t, session, result)
}

func TestGoaTradeService_ImplementsTradeService(t *testing.T) {
	// This test ensures that GoaTradeService implements service.TradeService interface
	var _ service.TradeService = (*GoaTradeService)(nil)
}
