package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// HistoricalPrice は履歴価格データの単一のエントリを表します。
type HistoricalPrice struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

// PriceHistoryReader はCSVファイルから履歴価格データを読み込むための構造体です。
type PriceHistoryReader struct {
	dataDir string
}

// NewPriceHistoryReader は新しいPriceHistoryReaderのインスタンスを作成します。
func NewPriceHistoryReader(dataDir string) *PriceHistoryReader {
	return &PriceHistoryReader{
		dataDir: dataDir,
	}
}

// ReadHistory は指定されたシンボルの履歴価格データを読み込み、パースします。
func (r *PriceHistoryReader) ReadHistory(symbol string) ([]*HistoricalPrice, error) {
	filePath := filepath.Join(r.dataDir, fmt.Sprintf("%s.csv", symbol))
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open history file for %s: %w", symbol, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// ヘッダー行をスキップ
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header for %s: %w", symbol, err)
	}

	var history []*HistoricalPrice
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record for %s: %w", symbol, err)
		}

		if len(record) != 6 {
			return nil, fmt.Errorf("invalid record length for %s: expected 6, got %d", symbol, len(record))
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse date '%s' for %s: %w", record[0], symbol, err)
		}
		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse open price '%s' for %s: %w", record[1], symbol, err)
		}
		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse high price '%s' for %s: %w", record[2], symbol, err)
		}
		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse low price '%s' for %s: %w", record[3], symbol, err)
		}
		close, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse close price '%s' for %s: %w", record[4], symbol, err)
		}
		volume, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse volume '%s' for %s: %w", record[5], symbol, err)
		}

		history = append(history, &HistoricalPrice{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	return history, nil
}
