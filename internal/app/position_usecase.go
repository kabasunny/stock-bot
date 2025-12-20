package app

import (
	"context"
	"stock-bot/internal/infrastructure/client"
)

// PositionType defines the type of a holding.
type PositionType string

const (
	PositionTypeCash        PositionType = "CASH"
	PositionTypeMarginLong  PositionType = "MARGIN_LONG"
	PositionTypeMarginShort PositionType = "MARGIN_SHORT"
)

// Position represents a single unified trading position.
type Position struct {
	Symbol           string
	PositionType     PositionType
	Quantity         float64
	AverageCost      float64
	CurrentPrice     float64
	UnrealizedPL     float64
	UnrealizedPLRate float64
	OpenedDate       string
}

// PositionUseCase defines the interface for position-related use cases.
type PositionUseCase interface {
	ListPositions(ctx context.Context, session *client.Session, filterType string) ([]*Position, error)
}