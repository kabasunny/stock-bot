package agent

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// TradeSignal は売買区分を表す型
type TradeSignal uint8

const (
	BuySignal  TradeSignal = 0x01
	SellSignal TradeSignal = 0x02
)

// SignalRecord はシグナルファイル内の1レコードを表す
type SignalRecord struct {
	Symbol uint16
	Signal TradeSignal
}

// ReadSignalFile は指定されたバイナリファイルからシグナルを読み込む
// ファイルフォーマット: [銘柄コード(uint16)][売買区分(uint8)] が連続
func ReadSignalFile(filePath string) ([]*SignalRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open signal file %s: %w", filePath, err)
	}
	defer file.Close()

	signals := make([]*SignalRecord, 0)
	// リトルエンディアンで読み込む
	// 3バイトずつ読み込む
	for {
		record := &SignalRecord{}
		// 銘柄コードの読み込み
		err := binary.Read(file, binary.LittleEndian, &record.Symbol)
		if err != nil {
			if err == io.EOF {
				break // ファイルの終端に達したら終了
			}
			return nil, fmt.Errorf("failed to read symbol from signal file: %w", err)
		}

		// 売買区分の読み込み
		err = binary.Read(file, binary.LittleEndian, &record.Signal)
		if err != nil {
			// Symbolを読んだ直後にEOFになるケース
			if err == io.EOF {
				return nil, fmt.Errorf("incomplete record found at end of signal file: symbol read but signal is missing")
			}
			return nil, fmt.Errorf("failed to read signal from signal file: %w", err)
		}
		
		// 不正なシグナル値でないかチェック
		if record.Signal != BuySignal && record.Signal != SellSignal {
			return nil, fmt.Errorf("invalid trade signal value found: 0x%02x", record.Signal)
		}

		signals = append(signals, record)
	}

	return signals, nil
}
