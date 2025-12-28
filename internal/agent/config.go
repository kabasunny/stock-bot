package agent

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// AgentConfig はエージェント固有の設定を保持する構造体
type AgentConfig struct {
	Agent            AgentSettings    `yaml:"agent"`
	StrategySettings StrategySettings `yaml:"strategy_settings"`
	API              struct {
		GoWrapperURL    string `yaml:"go_wrapper_url"`
		PythonSignalURL string `yaml:"python_signal_url"` // シグナルメーカーが将来的にHTTPの場合
	} `yaml:"api"`
}

// AgentSettings はエージェントの基本設定
type AgentSettings struct {
	Strategy          string        `yaml:"strategy"`
	ExecutionInterval time.Duration `yaml:"execution_interval"`
	LogLevel          string        `yaml:"log_level"`
	Timezone          string        `yaml:"timezone"`
}

// StrategySettings は各戦略の設定をまとめる
type StrategySettings struct {
	Swingtrade SwingtradeSettings `yaml:"swingtrade"`
	Daytrade   struct {
		// デイトレード戦略用の設定
	} `yaml:"daytrade"`
}

// SwingtradeSettings はスイングトレード戦略固有の設定
type SwingtradeSettings struct {
	TargetSymbols             []string  `yaml:"target_symbols"`
	TradeRiskPercentage       float64   `yaml:"trade_risk_percentage"`
	MaxPositionSizePercentage float64   `yaml:"max_position_size_percentage"`
	UnitSize                  int       `yaml:"unit_size"`
	ProfitTakeRate            float64   `yaml:"profit_take_rate"`
	StopLossRate              float64   `yaml:"stop_loss_rate"`
	TrailingStopTriggerRate   float64   `yaml:"trailing_stop_trigger_rate"`
	TrailingStopRate          float64   `yaml:"trailing_stop_rate"`
	SignalFilePattern         string    `yaml:"signal_file_pattern"` // シグナルファイルのパターンを追加
	ATRPeriod                 int       `yaml:"atr_period"`          // New: ATR期間
	RiskPerATR                float64   `yaml:"risk_per_atr"`        // New: ATR単位でのリスク量
	StopLossATRMultiplier     float64   `yaml:"stop_loss_atr_multiplier"` // New: ATRを基準とした損切り幅の乗数
}

// LoadAgentConfig は指定されたYAMLファイルからエージェントの設定を読み込む
func LoadAgentConfig(configPath string) (*AgentConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read agent config file %s: %w", configPath, err)
	}

	var cfg AgentConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal agent config from %s: %w", configPath, err)
	}

	// デフォルト値やバリデーション (必要に応じて追加)
	if cfg.Agent.ExecutionInterval == 0 {
		cfg.Agent.ExecutionInterval = 1 * time.Minute // デフォルトは1分
	}
	if cfg.Agent.LogLevel == "" {
		cfg.Agent.LogLevel = "info"
	}
	if cfg.Agent.Timezone == "" {
		cfg.Agent.Timezone = "Asia/Tokyo"
	}
	if cfg.StrategySettings.Swingtrade.SignalFilePattern == "" {
		// デフォルトのシグナルファイルパターン
		cfg.StrategySettings.Swingtrade.SignalFilePattern = "./signals/*.bin"
	}
	if cfg.StrategySettings.Swingtrade.TradeRiskPercentage == 0 {
		cfg.StrategySettings.Swingtrade.TradeRiskPercentage = 0.02 // デフォルトは2%
	}
	if cfg.StrategySettings.Swingtrade.MaxPositionSizePercentage == 0 {
		cfg.StrategySettings.Swingtrade.MaxPositionSizePercentage = 0.25 // デフォルトは25%
	}
	if cfg.StrategySettings.Swingtrade.UnitSize == 0 {
		cfg.StrategySettings.Swingtrade.UnitSize = 100 // デフォルトは100株
	}
	if cfg.StrategySettings.Swingtrade.ATRPeriod == 0 {
		cfg.StrategySettings.Swingtrade.ATRPeriod = 14 // デフォルトは14期間
	}
	if cfg.StrategySettings.Swingtrade.RiskPerATR == 0 {
		cfg.StrategySettings.Swingtrade.RiskPerATR = 0.5 // デフォルトは0.5ATR
	}
	if cfg.StrategySettings.Swingtrade.StopLossATRMultiplier == 0 {
		cfg.StrategySettings.Swingtrade.StopLossATRMultiplier = 2.0 // デフォルトは2.0ATR
	}


	return &cfg, nil
}
