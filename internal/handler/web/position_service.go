package web

import (
	"context"
	"log/slog"
	"stock-bot/gen/position"
	"stock-bot/internal/app"
	"stock-bot/internal/infrastructure/client"
)

// PositionService implements the position.Service interface.
type PositionService struct {
	positionUseCase app.PositionUseCase
	logger          *slog.Logger
    session         *client.Session
}

// NewPositionService creates a new PositionService.
func NewPositionService(positionUseCase app.PositionUseCase, logger *slog.Logger, session *client.Session) *PositionService {
	return &PositionService{
		positionUseCase: positionUseCase,
		logger:          logger,
        session:         session,
	}
}

// List implements the list action.
func (s *PositionService) List(ctx context.Context, p *position.ListPayload) (res *position.StockbotPositionCollection, err error) {
	s.logger.Info("ListPositions called", slog.String("filterType", p.Type))

	filterType := "all" // デフォルト
	if p.Type != "" {
		filterType = p.Type
	}

	// Call the use case
	appPositions, err := s.positionUseCase.ListPositions(ctx, s.session, filterType)

	if err != nil {
		s.logger.Error("Failed to list positions from use case", slog.Any("error", err))
		return nil, err
	}



	// Map app.Position to position.PositionResult (Goa generated type)

	// And collect into a PositionCollection

	resultPositions := make([]*position.PositionResult, len(appPositions))

	for i, ap := range appPositions {

		// Create pointers for optional fields

		currentPrice := &ap.CurrentPrice

		unrealizedPL := &ap.UnrealizedPL

		unrealizedPLRate := &ap.UnrealizedPLRate

		openedDate := &ap.OpenedDate



		resultPositions[i] = &position.PositionResult{

			Symbol:            ap.Symbol,

			PositionType:      string(ap.PositionType), // Convert custom type to string

			Quantity:          ap.Quantity,

			AverageCost:       ap.AverageCost,

			CurrentPrice:      currentPrice,

			UnrealizedPl:      unrealizedPL,

			UnrealizedPlRate:  unrealizedPLRate,

			OpenedDate:        openedDate,

		}

	}



	res = &position.StockbotPositionCollection{

		Positions: resultPositions,

	}



	s.logger.Info("ListPositions successful", slog.Int("count", len(resultPositions)))

	return res, nil

}
