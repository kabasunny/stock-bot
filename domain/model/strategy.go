package model

import (
	"time"
)

// StrategyType は戦略の種別
type StrategyType string

const (
	StrategyTypeSwing  StrategyType = "swing"
	StrategyTypeDay    StrategyType = "day"
	StrategyTypeScalp  StrategyType = "scalp"
	StrategyTypeCustom StrategyType = "custom"
)

// StrategyStatus は戦略の状態
type StrategyStatus string

const (
	StrategyStatusActive   StrategyStatus = "active"
	StrategyStatusInactive StrategyStatus = "inactive"
	StrategyStatusPaused   StrategyStatus = "paused"
	StrategyStatusStopped  StrategyStatus = "stopped"
)

// Strategy は取引戦略を表すドメインモデル
type Strategy struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        StrategyType   `json:"type"`
	Status      StrategyStatus `json:"status"`
	Description string         `json:"description"`

	// 戦略設定
	Config StrategyConfig `json:"config"`

	// リスク管理
	RiskLimits RiskLimits `json:"risk_limits"`

	// 実行統計
	Statistics StrategyStats `json:"statistics"`

	// メタデータ
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
}

// StrategyConfig は戦略の設定
type StrategyConfig struct {
	// 対象銘柄
	TargetSymbols []string `json:"target_symbols"`

	// 実行間隔
	ExecutionInterval time.Duration `json:"execution_interval"`

	// シグナル設定
	SignalSettings SignalSettings `json:"signal_settings"`

	// 注文設定
	OrderSettings OrderSettings `json:"order_settings"`

	// カスタムパラメータ
	CustomParams map[string]interface{} `json:"custom_params"`
}

// SignalSettings はシグナルの設定
type SignalSettings struct {
	SignalFilePattern string            `json:"signal_file_pattern"`
	SignalThreshold   float64           `json:"signal_threshold"`
	CustomFilters     map[string]string `json:"custom_filters"`
}

// OrderSettings は注文の設定
type OrderSettings struct {
	DefaultOrderType       OrderType           `json:"default_order_type"`
	DefaultPositionAccount PositionAccountType `json:"default_position_account"`
	MaxOrderAmount         float64             `json:"max_order_amount"`
	MaxOrderAmountPercent  float64             `json:"max_order_amount_percent"`
	MinOrderAmount         float64             `json:"min_order_amount"`
	OrderSizeCalculation   string              `json:"order_size_calculation"` // "fixed", "percent", "kelly"
}

// RiskLimits はリスク制限
type RiskLimits struct {
	// 最大損失額
	MaxLossAmount float64 `json:"max_loss_amount"`

	// 最大損失率
	MaxLossPercent float64 `json:"max_loss_percent"`

	// 最大ポジション数
	MaxPositions int `json:"max_positions"`

	// 銘柄あたりの最大ポジション数
	MaxPositionsPerSymbol int `json:"max_positions_per_symbol"`

	// 最大レバレッジ
	MaxLeverage float64 `json:"max_leverage"`

	// ドローダウン制限
	MaxDrawdown float64 `json:"max_drawdown"`

	// 日次損失制限
	DailyLossLimit float64 `json:"daily_loss_limit"`
}

// StrategyStats は戦略の実行統計
type StrategyStats struct {
	// 実行回数
	ExecutionCount int `json:"execution_count"`

	// 注文統計
	TotalOrders      int `json:"total_orders"`
	SuccessfulOrders int `json:"successful_orders"`
	FailedOrders     int `json:"failed_orders"`

	// 損益統計
	TotalPL      float64 `json:"total_pl"`
	RealizedPL   float64 `json:"realized_pl"`
	UnrealizedPL float64 `json:"unrealized_pl"`

	// 勝率統計
	WinCount  int     `json:"win_count"`
	LossCount int     `json:"loss_count"`
	WinRate   float64 `json:"win_rate"`

	// 最大ドローダウン
	MaxDrawdown     float64 `json:"max_drawdown"`
	CurrentDrawdown float64 `json:"current_drawdown"`

	// 最終実行時刻
	LastExecutedAt time.Time `json:"last_executed_at"`
}

// IsActive は戦略がアクティブかどうかを判定
func (s *Strategy) IsActive() bool {
	return s.Status == StrategyStatusActive
}

// CanExecute は戦略が実行可能かどうかを判定
func (s *Strategy) CanExecute() bool {
	return s.Status == StrategyStatusActive
}

// UpdateStatistics は統計を更新
func (s *Strategy) UpdateStatistics(pl float64, isWin bool) {
	s.Statistics.ExecutionCount++
	s.Statistics.TotalPL += pl

	if isWin {
		s.Statistics.WinCount++
	} else {
		s.Statistics.LossCount++
	}

	totalTrades := s.Statistics.WinCount + s.Statistics.LossCount
	if totalTrades > 0 {
		s.Statistics.WinRate = float64(s.Statistics.WinCount) / float64(totalTrades)
	}

	s.Statistics.LastExecutedAt = time.Now()
	s.UpdatedAt = time.Now()
}

// CheckRiskLimits はリスク制限をチェック
func (s *Strategy) CheckRiskLimits() []string {
	violations := make([]string, 0)

	// 最大損失額チェック
	if s.RiskLimits.MaxLossAmount > 0 && s.Statistics.TotalPL < -s.RiskLimits.MaxLossAmount {
		violations = append(violations, "max_loss_amount_exceeded")
	}

	// 最大ドローダウンチェック
	if s.RiskLimits.MaxDrawdown > 0 && s.Statistics.CurrentDrawdown > s.RiskLimits.MaxDrawdown {
		violations = append(violations, "max_drawdown_exceeded")
	}

	return violations
}

// Activate は戦略をアクティブ化
func (s *Strategy) Activate() {
	s.Status = StrategyStatusActive
	s.UpdatedAt = time.Now()
}

// Deactivate は戦略を非アクティブ化
func (s *Strategy) Deactivate() {
	s.Status = StrategyStatusInactive
	s.UpdatedAt = time.Now()
}

// Pause は戦略を一時停止
func (s *Strategy) Pause() {
	s.Status = StrategyStatusPaused
	s.UpdatedAt = time.Now()
}

// Stop は戦略を停止
func (s *Strategy) Stop() {
	s.Status = StrategyStatusStopped
	s.UpdatedAt = time.Now()
}
