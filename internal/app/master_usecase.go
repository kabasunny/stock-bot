package app

import (
	"context"
	"errors"
	"stock-bot/internal/infrastructure/client"
)

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = errors.New("not found")

// StockMasterResult holds basic master data for a single stock.
type StockMasterResult struct {
	Symbol       string
	Name         string
	NameKana     string
	Market       string
	IndustryCode string
	IndustryName string
}

// MasterUseCase defines the interface for master data related use cases.
type MasterUseCase interface {
	GetStock(ctx context.Context, symbol string) (*StockMasterResult, error)
	DownloadAndStoreMasterData(ctx context.Context, session *client.Session) error
}