package agent_test // agentパッケージをテストするのでagent_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"stock-bot/internal/agent" // agentパッケージをインポート
	"github.com/stretchr/testify/assert"

)

func TestLoadAgentConfig(t *testing.T) {
	// 一時的なconfig.yamlファイルを作成
	testConfigContent := `
agent:
  strategy: swingtrade
  execution_interval: 1m30s
  log_level: debug
  timezone: "America/New_York"
strategy_settings:
  swingtrade:
    target_symbols:
      - "AAPL"
      - "GOOG"
    trade_risk_percentage: 0.15
    unit_size: 50
    profit_take_rate: 3.5
    stop_loss_rate: 1.0
    signal_file_pattern: "./test_signals/*.bin"
api:
  go_wrapper_url: "http://localhost:8081"
  python_signal_url: "http://localhost:5001"
`
	// 一時ファイルに書き出し
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	err := os.WriteFile(configPath, []byte(testConfigContent), 0644)
	assert.NoError(t, err)

	// 設定をロード
	cfg, err := agent.LoadAgentConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// 値の検証
	assert.Equal(t, "swingtrade", cfg.Agent.Strategy)
	assert.Equal(t, 1*time.Minute+30*time.Second, cfg.Agent.ExecutionInterval)
	assert.Equal(t, "debug", cfg.Agent.LogLevel)
	assert.Equal(t, "America/New_York", cfg.Agent.Timezone)

	assert.ElementsMatch(t, []string{"AAPL", "GOOG"}, cfg.StrategySettings.Swingtrade.TargetSymbols)
	assert.Equal(t, 0.15, cfg.StrategySettings.Swingtrade.TradeRiskPercentage)
	assert.Equal(t, 50, cfg.StrategySettings.Swingtrade.UnitSize)
	assert.Equal(t, 3.5, cfg.StrategySettings.Swingtrade.ProfitTakeRate)
	assert.Equal(t, 1.0, cfg.StrategySettings.Swingtrade.StopLossRate)
	assert.Equal(t, "./test_signals/*.bin", cfg.StrategySettings.Swingtrade.SignalFilePattern)

	assert.Equal(t, "http://localhost:8081", cfg.API.GoWrapperURL)
	assert.Equal(t, "http://localhost:5001", cfg.API.PythonSignalURL)
}

func TestLoadAgentConfig_DefaultValues(t *testing.T) {
	// デフォルト値が適用される一時的なconfig.yamlファイルを作成 (一部の項目を省略)
	testConfigContent := `
agent:
  strategy: daytrade
  # execution_interval は省略
  # log_level は省略
  # timezone は省略
strategy_settings:
  swingtrade:
    target_symbols:
      - "MSFT"
    # profit_take_rate は省略
    # stop_loss_rate は省略
    # signal_file_pattern は省略
    # trade_risk_percentage は省略
    # unit_size は省略
api:
  go_wrapper_url: "http://localhost:8080"
  # python_signal_url は省略
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config_default.yaml")
	err := os.WriteFile(configPath, []byte(testConfigContent), 0644)
	assert.NoError(t, err)

	cfg, err := agent.LoadAgentConfig(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// デフォルト値の検証
	assert.Equal(t, "daytrade", cfg.Agent.Strategy)
	assert.Equal(t, 1*time.Minute, cfg.Agent.ExecutionInterval) // デフォルト値
	assert.Equal(t, "info", cfg.Agent.LogLevel)               // デフォルト値
	assert.Equal(t, "Asia/Tokyo", cfg.Agent.Timezone)         // デフォルト値

	assert.ElementsMatch(t, []string{"MSFT"}, cfg.StrategySettings.Swingtrade.TargetSymbols)
	assert.Equal(t, 0.02, cfg.StrategySettings.Swingtrade.TradeRiskPercentage) // デフォルト値
	assert.Equal(t, 100, cfg.StrategySettings.Swingtrade.UnitSize)           // デフォルト値
	// float64のデフォルト値は0.0なので、ここではテストしないか、初期化された構造体の値を期待する
	assert.Equal(t, 0.0, cfg.StrategySettings.Swingtrade.ProfitTakeRate)
	assert.Equal(t, 0.0, cfg.StrategySettings.Swingtrade.StopLossRate)
	assert.Equal(t, "./signals/*.bin", cfg.StrategySettings.Swingtrade.SignalFilePattern) // デフォルト値

	assert.Equal(t, "http://localhost:8080", cfg.API.GoWrapperURL)
	assert.Equal(t, "", cfg.API.PythonSignalURL) // デフォルト値は空文字列
}

func TestLoadAgentConfig_FileNotFound(t *testing.T) {
	cfg, err := agent.LoadAgentConfig("non_existent_config.yaml")
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to read agent config file")
}

func TestLoadAgentConfig_InvalidYaml(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid_config.yaml")
	err := os.WriteFile(configPath, []byte("agent: \n  strategy: : invalid_value"), 0644) // 不正なYAML
	assert.NoError(t, err)

	cfg, err := agent.LoadAgentConfig(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to unmarshal agent config")
}
