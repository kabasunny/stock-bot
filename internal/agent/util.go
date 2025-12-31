package agent

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"stock-bot/domain/service"
	"time"
)

// FindSignalFile は指定されたパターンに一致するシグナルファイルを探し、最も新しい更新日時を持つファイルを返す
func FindSignalFile(pattern string) (string, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to glob pattern %q: %w", pattern, err)
	}
	if len(files) == 0 {
		return "", nil // ファイルが見つからなくてもエラーではない
	}

	var latestFile string
	var latestModTime time.Time

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			// ファイルが見つからない、またはアクセスできない場合はスキップ
			// ただし、エラーとしてログに出力する方が良い場合もあるが、ここでは堅牢性を優先しスキップ
			continue
		}

		if latestFile == "" || info.ModTime().After(latestModTime) {
			latestModTime = info.ModTime()
			latestFile = file
		}
	}

	if latestFile == "" {
		// globでファイルが見つかったが、os.Statで全て失敗した場合
		return "", fmt.Errorf("no accessible signal files found matching pattern %q", pattern)
	}

	return latestFile, nil
}

// calculateATR は指定された期間のATR (Average True Range) を計算する
func calculateATR(history []*service.HistoricalPrice, atrPeriod int) (float64, error) {
	if len(history) < atrPeriod+1 { // ATR計算には少なくともATRPeriod+1個のデータが必要
		return 0, fmt.Errorf("not enough historical data to calculate ATR for period %d. Requires at least %d data points, got %d", atrPeriod, atrPeriod+1, len(history))
	}

	trueRanges := make([]float64, len(history)-1)
	for i := 1; i < len(history); i++ {
		// TR = max(H - L, abs(H - C_prev), abs(L - C_prev))
		highLow := history[i].High - history[i].Low
		highPrevClose := math.Abs(history[i].High - history[i-1].Close)
		lowPrevClose := math.Abs(history[i].Low - history[i-1].Close)

		tr := math.Max(highLow, math.Max(highPrevClose, lowPrevClose))
		trueRanges[i-1] = tr
	}

	// 最初のATR値は単純移動平均 (SMA) で計算
	initialATRSum := 0.0
	for i := 0; i < atrPeriod; i++ {
		initialATRSum += trueRanges[i]
	}
	atr := initialATRSum / float64(atrPeriod)

	// それ以降のATR値は Wilder の平滑化方法で計算
	for i := atrPeriod; i < len(trueRanges); i++ {
		atr = (atr*float64(atrPeriod-1) + trueRanges[i]) / float64(atrPeriod)
	}

	return atr, nil
}
