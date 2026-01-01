package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStrategy_IsActive(t *testing.T) {
	strategy := &Strategy{
		Status: StrategyStatusActive,
	}
	assert.True(t, strategy.IsActive())

	strategy.Status = StrategyStatusInactive
	assert.False(t, strategy.IsActive())

	strategy.Status = StrategyStatusPaused
	assert.False(t, strategy.IsActive())
}

func TestStrategy_CanExecute(t *testing.T) {
	strategy := &Strategy{
		Status: StrategyStatusActive,
	}
	assert.True(t, strategy.CanExecute())

	strategy.Status = StrategyStatusInactive
	assert.False(t, strategy.CanExecute())

	strategy.Status = StrategyStatusPaused
	assert.False(t, strategy.CanExecute())

	strategy.Status = StrategyStatusStopped
	assert.False(t, strategy.CanExecute())
}

func TestStrategy_UpdateStatistics(t *testing.T) {
	strategy := &Strategy{
		Statistics: StrategyStats{
			ExecutionCount: 0,
			TotalPL:        0.0,
			WinCount:       0,
			LossCount:      0,
			WinRate:        0.0,
		},
	}

	// Test win
	strategy.UpdateStatistics(5000.0, true)
	assert.Equal(t, 1, strategy.Statistics.ExecutionCount)
	assert.Equal(t, 5000.0, strategy.Statistics.TotalPL)
	assert.Equal(t, 1, strategy.Statistics.WinCount)
	assert.Equal(t, 0, strategy.Statistics.LossCount)
	assert.Equal(t, 1.0, strategy.Statistics.WinRate)

	// Test loss
	strategy.UpdateStatistics(-2000.0, false)
	assert.Equal(t, 2, strategy.Statistics.ExecutionCount)
	assert.Equal(t, 3000.0, strategy.Statistics.TotalPL)
	assert.Equal(t, 1, strategy.Statistics.WinCount)
	assert.Equal(t, 1, strategy.Statistics.LossCount)
	assert.Equal(t, 0.5, strategy.Statistics.WinRate)

	// Verify timestamp was updated
	assert.True(t, strategy.Statistics.LastExecutedAt.After(time.Time{}))
	assert.True(t, strategy.UpdatedAt.After(time.Time{}))
}

func TestStrategy_CheckRiskLimits(t *testing.T) {
	strategy := &Strategy{
		RiskLimits: RiskLimits{
			MaxLossAmount: 10000.0,
			MaxDrawdown:   5000.0,
		},
		Statistics: StrategyStats{
			TotalPL:         -15000.0, // Exceeds max loss
			CurrentDrawdown: 6000.0,   // Exceeds max drawdown
		},
	}

	violations := strategy.CheckRiskLimits()
	assert.Len(t, violations, 2)
	assert.Contains(t, violations, "max_loss_amount_exceeded")
	assert.Contains(t, violations, "max_drawdown_exceeded")
}

func TestStrategy_CheckRiskLimits_NoViolations(t *testing.T) {
	strategy := &Strategy{
		RiskLimits: RiskLimits{
			MaxLossAmount: 10000.0,
			MaxDrawdown:   5000.0,
		},
		Statistics: StrategyStats{
			TotalPL:         -5000.0, // Within limit
			CurrentDrawdown: 3000.0,  // Within limit
		},
	}

	violations := strategy.CheckRiskLimits()
	assert.Len(t, violations, 0)
}

func TestStrategy_Activate(t *testing.T) {
	strategy := &Strategy{
		Status: StrategyStatusInactive,
	}

	strategy.Activate()
	assert.Equal(t, StrategyStatusActive, strategy.Status)
	assert.True(t, strategy.UpdatedAt.After(time.Time{}))
}

func TestStrategy_Deactivate(t *testing.T) {
	strategy := &Strategy{
		Status: StrategyStatusActive,
	}

	strategy.Deactivate()
	assert.Equal(t, StrategyStatusInactive, strategy.Status)
	assert.True(t, strategy.UpdatedAt.After(time.Time{}))
}

func TestStrategy_Pause(t *testing.T) {
	strategy := &Strategy{
		Status: StrategyStatusActive,
	}

	strategy.Pause()
	assert.Equal(t, StrategyStatusPaused, strategy.Status)
	assert.True(t, strategy.UpdatedAt.After(time.Time{}))
}

func TestStrategy_Stop(t *testing.T) {
	strategy := &Strategy{
		Status: StrategyStatusActive,
	}

	strategy.Stop()
	assert.Equal(t, StrategyStatusStopped, strategy.Status)
	assert.True(t, strategy.UpdatedAt.After(time.Time{}))
}

func TestStrategyType_Constants(t *testing.T) {
	assert.Equal(t, StrategyType("swing"), StrategyTypeSwing)
	assert.Equal(t, StrategyType("day"), StrategyTypeDay)
	assert.Equal(t, StrategyType("scalp"), StrategyTypeScalp)
	assert.Equal(t, StrategyType("custom"), StrategyTypeCustom)
}

func TestStrategyStatus_Constants(t *testing.T) {
	assert.Equal(t, StrategyStatus("active"), StrategyStatusActive)
	assert.Equal(t, StrategyStatus("inactive"), StrategyStatusInactive)
	assert.Equal(t, StrategyStatus("paused"), StrategyStatusPaused)
	assert.Equal(t, StrategyStatus("stopped"), StrategyStatusStopped)
}
