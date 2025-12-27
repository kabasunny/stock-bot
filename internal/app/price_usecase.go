package app

import (
	"context"
	"stock-bot/gen/price" // Generated Goa service types
)

// HistoricalPriceItem represents a single historical price data point for the app layer.
type HistoricalPriceItem struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume uint64
}

// HistoricalPriceResult represents the historical price information for a stock for the app layer.
type HistoricalPriceResult struct {
	Symbol  string
	History []*HistoricalPriceItem
}

// PriceUseCase is the interface that wraps the basic Price methods.
type PriceUseCase interface {
	// Get retrieves the current price for a specified stock symbol.
	Get(ctx context.Context, symbol string) (*price.StockbotPrice, error)
	// GetHistory retrieves historical price data for a specified stock symbol.
	GetHistory(ctx context.Context, symbol string, days uint) (*HistoricalPriceResult, error)
}
