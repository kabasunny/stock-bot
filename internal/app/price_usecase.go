package app

import (
	"context"
	"stock-bot/gen/price" // Generated Goa service types
)

// PriceUseCase is the interface that wraps the basic Price methods.
type PriceUseCase interface {
	// Get retrieves the current price for a specified stock symbol.
	Get(ctx context.Context, symbol string) (*price.StockbotPrice, error)
}
