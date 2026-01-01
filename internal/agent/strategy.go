package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// StrategySignal は戦略の取引シグナルを表す
type StrategySignal struct {
	ShouldTrade  bool
	OrderRequest *PlaceOrderRequest
}

// Strategy は取引戦略のインターフェース
type Strategy interface {
	// Name は戦略名を返す
	Name() string

	// Evaluate はシグナルを評価する
	Evaluate(ctx context.Context, data *MarketData) (*StrategySignal, error)

	// GetExecutionInterval は実行間隔を返す
	GetExecutionInterval() time.Duration

	// GetRiskLimits はリスク制限を返す
	GetRiskLimits() *RiskLimits
}

// MarketData は市場データを表す
type MarketData struct {
	Balance   *Balance
	Positions []*Position
	Orders    []*Order
	Timestamp time.Time
}

// RiskLimits はリスク制限を表す
type RiskLimits struct {
	MaxLossAmount         float64 // 最大損失額
	MaxLossPercent        float64 // 最大損失率
	MaxPositions          int     // 最大ポジション数
	MaxOrderAmount        float64 // 最大注文金額
	DailyLossLimit        float64 // 日次損失制限
	MaxPositionsPerSymbol int     // 銘柄あたりの最大ポジション数
}

// SimpleStrategy は簡単な戦略実装
type SimpleStrategy struct {
	name              string
	targetSymbol      string
	orderQuantity     uint
	executionInterval time.Duration
	riskLimits        *RiskLimits
	logger            *slog.Logger
}

// NewSimpleStrategy は新しいSimpleStrategyを作成する
func NewSimpleStrategy(name, targetSymbol string, orderQuantity uint, logger *slog.Logger) *SimpleStrategy {
	return &SimpleStrategy{
		name:              name,
		targetSymbol:      targetSymbol,
		orderQuantity:     orderQuantity,
		executionInterval: 30 * time.Second,
		riskLimits: &RiskLimits{
			MaxLossAmount:         100000, // 10万円
			MaxLossPercent:        5.0,    // 5%
			MaxPositions:          3,      // 最大3ポジション
			MaxOrderAmount:        500000, // 50万円
			DailyLossLimit:        50000,  // 日次5万円
			MaxPositionsPerSymbol: 1,      // 銘柄あたり1ポジション
		},
		logger: logger,
	}
}

// Name は戦略名を返す
func (s *SimpleStrategy) Name() string {
	return s.name
}

// GetExecutionInterval は実行間隔を返す
func (s *SimpleStrategy) GetExecutionInterval() time.Duration {
	return s.executionInterval
}

// GetRiskLimits はリスク制限を返す
func (s *SimpleStrategy) GetRiskLimits() *RiskLimits {
	return s.riskLimits
}

// Evaluate はシグナルを評価する
func (s *SimpleStrategy) Evaluate(ctx context.Context, data *MarketData) (*StrategySignal, error) {
	s.logger.Debug("Evaluating strategy",
		"strategy", s.name,
		"target_symbol", s.targetSymbol,
		"cash", data.Balance.Cash,
		"positions_count", len(data.Positions))

	// リスクチェック
	if err := s.checkRiskLimits(data); err != nil {
		s.logger.Warn("Risk limit exceeded", "error", err)
		return &StrategySignal{ShouldTrade: false}, nil
	}

	// 既存ポジションチェック
	hasPosition := s.hasPositionForSymbol(data.Positions, s.targetSymbol)
	if hasPosition {
		s.logger.Debug("Already have position for symbol", "symbol", s.targetSymbol)
		return &StrategySignal{ShouldTrade: false}, nil
	}

	// 現金残高チェック
	requiredCash := float64(s.orderQuantity) * 3000 // 仮定価格3000円
	if data.Balance.Cash < requiredCash {
		s.logger.Debug("Insufficient cash",
			"required", requiredCash,
			"available", data.Balance.Cash)
		return &StrategySignal{ShouldTrade: false}, nil
	}

	// 簡単な買いシグナル生成
	return &StrategySignal{
		ShouldTrade: true,
		OrderRequest: &PlaceOrderRequest{
			Symbol:              s.targetSymbol,
			TradeType:           "BUY",
			OrderType:           "MARKET",
			Quantity:            s.orderQuantity,
			Price:               0, // 成行注文
			PositionAccountType: "CASH",
		},
	}, nil
}

// checkRiskLimits はリスク制限をチェックする
func (s *SimpleStrategy) checkRiskLimits(data *MarketData) error {
	// ポジション数チェック
	if len(data.Positions) >= s.riskLimits.MaxPositions {
		return fmt.Errorf("max positions exceeded: %d >= %d",
			len(data.Positions), s.riskLimits.MaxPositions)
	}

	// 銘柄あたりのポジション数チェック
	symbolPositions := s.countPositionsForSymbol(data.Positions, s.targetSymbol)
	if symbolPositions >= s.riskLimits.MaxPositionsPerSymbol {
		return fmt.Errorf("max positions per symbol exceeded for %s: %d >= %d",
			s.targetSymbol, symbolPositions, s.riskLimits.MaxPositionsPerSymbol)
	}

	// 注文金額チェック
	orderAmount := float64(s.orderQuantity) * 3000 // 仮定価格
	if orderAmount > s.riskLimits.MaxOrderAmount {
		return fmt.Errorf("order amount exceeds limit: %.2f > %.2f",
			orderAmount, s.riskLimits.MaxOrderAmount)
	}

	return nil
}

// hasPositionForSymbol は指定銘柄のポジションを持っているかチェックする
func (s *SimpleStrategy) hasPositionForSymbol(positions []*Position, symbol string) bool {
	for _, pos := range positions {
		if pos.Symbol == symbol {
			return true
		}
	}
	return false
}

// countPositionsForSymbol は指定銘柄のポジション数を数える
func (s *SimpleStrategy) countPositionsForSymbol(positions []*Position, symbol string) int {
	count := 0
	for _, pos := range positions {
		if pos.Symbol == symbol {
			count++
		}
	}
	return count
}

// SwingStrategy はスイング戦略の実装
type SwingStrategy struct {
	*SimpleStrategy
	holdingPeriod time.Duration
}

// NewSwingStrategy は新しいSwingStrategyを作成する
func NewSwingStrategy(name, targetSymbol string, orderQuantity uint, logger *slog.Logger) *SwingStrategy {
	simple := NewSimpleStrategy(name, targetSymbol, orderQuantity, logger)
	simple.executionInterval = 1 * time.Hour // 1時間間隔

	return &SwingStrategy{
		SimpleStrategy: simple,
		holdingPeriod:  24 * time.Hour, // 24時間保持
	}
}

// Evaluate はスイング戦略のシグナルを評価する
func (s *SwingStrategy) Evaluate(ctx context.Context, data *MarketData) (*StrategySignal, error) {
	s.logger.Debug("Evaluating swing strategy", "strategy", s.name)

	// 基本チェックは親クラスに委譲
	signal, err := s.SimpleStrategy.Evaluate(ctx, data)
	if err != nil || !signal.ShouldTrade {
		return signal, err
	}

	// スイング戦略固有のロジック
	// 例: 時間帯チェック（9:00-11:30のみエントリー）
	now := time.Now()
	if now.Hour() < 9 || (now.Hour() == 11 && now.Minute() > 30) || now.Hour() > 11 {
		s.logger.Debug("Outside trading hours for swing strategy")
		return &StrategySignal{ShouldTrade: false}, nil
	}

	return signal, nil
}

// DayTradingStrategy はデイトレード戦略の実装
type DayTradingStrategy struct {
	*SimpleStrategy
}

// NewDayTradingStrategy は新しいDayTradingStrategyを作成する
func NewDayTradingStrategy(name, targetSymbol string, orderQuantity uint, logger *slog.Logger) *DayTradingStrategy {
	simple := NewSimpleStrategy(name, targetSymbol, orderQuantity, logger)
	simple.executionInterval = 5 * time.Minute // 5分間隔

	return &DayTradingStrategy{
		SimpleStrategy: simple,
	}
}

// Evaluate はデイトレード戦略のシグナルを評価する
func (s *DayTradingStrategy) Evaluate(ctx context.Context, data *MarketData) (*StrategySignal, error) {
	s.logger.Debug("Evaluating day trading strategy", "strategy", s.name)

	// 基本チェックは親クラスに委譲
	signal, err := s.SimpleStrategy.Evaluate(ctx, data)
	if err != nil || !signal.ShouldTrade {
		return signal, err
	}

	// デイトレード戦略固有のロジック
	// 例: 取引時間チェック（9:00-14:30のみ）
	now := time.Now()
	if now.Hour() < 9 || now.Hour() > 14 || (now.Hour() == 14 && now.Minute() > 30) {
		s.logger.Debug("Outside day trading hours")
		return &StrategySignal{ShouldTrade: false}, nil
	}

	return signal, nil
}
