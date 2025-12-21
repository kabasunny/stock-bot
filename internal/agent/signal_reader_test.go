package agent_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"stock-bot/internal/agent"
	"github.com/stretchr/testify/assert"
)

// createDummySignalFile はテスト用のバイナリシグナルファイルを作成するヘルパー関数
func createDummySignalFile(t *testing.T, dir string, filename string, records []agent.SignalRecord) string {
	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	assert.NoError(t, err)
	defer file.Close()

	buf := new(bytes.Buffer)
	for _, rec := range records {
		err := binary.Write(buf, binary.LittleEndian, rec.Symbol)
		assert.NoError(t, err)
		err = binary.Write(buf, binary.LittleEndian, rec.Signal)
		assert.NoError(t, err)
	}
	_, err = file.Write(buf.Bytes())
	assert.NoError(t, err)

	return path
}

func TestReadSignalFile_Success(t *testing.T) {
	// テストデータ
	records := []agent.SignalRecord{
		{Symbol: 7203, Signal: agent.BuySignal},  // トヨタ
		{Symbol: 9984, Signal: agent.SellSignal}, // ソフトバンクG
		{Symbol: 6758, Signal: agent.BuySignal},  // ソニー
	}

	// ダミーファイル作成
	tmpDir := t.TempDir()
	filePath := createDummySignalFile(t, tmpDir, "signals.bin", records)

	// 読み込み実行
	readRecords, err := agent.ReadSignalFile(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, readRecords)
	assert.Equal(t, len(records), len(readRecords))

	// 内容の検証
	for i, expected := range records {
		assert.Equal(t, expected.Symbol, readRecords[i].Symbol)
		assert.Equal(t, expected.Signal, readRecords[i].Signal)
	}
}

func TestReadSignalFile_EmptyFile(t *testing.T) {
	// 空のファイルを作成
	tmpDir := t.TempDir()
	filePath := createDummySignalFile(t, tmpDir, "empty.bin", []agent.SignalRecord{})

	// 読み込み実行
	readRecords, err := agent.ReadSignalFile(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, readRecords)
	assert.Empty(t, readRecords)
}

func TestReadSignalFile_FileNotFound(t *testing.T) {
	readRecords, err := agent.ReadSignalFile("non_existent_file.bin")
	assert.Error(t, err)
	assert.Nil(t, readRecords)
	assert.Contains(t, err.Error(), "failed to open signal file")
}

func TestReadSignalFile_IncompleteRecord(t *testing.T) {
	// 不完全なレコードを持つファイルを作成 (Symbolのみ)
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "incomplete.bin")
	file, err := os.Create(path)
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, uint16(1234))
	assert.NoError(t, err)
	_, err = file.Write(buf.Bytes())
	assert.NoError(t, err)
	file.Close()

	// 読み込み実行
	readRecords, err := agent.ReadSignalFile(path)
	assert.Error(t, err)
	assert.Nil(t, readRecords)
	assert.Contains(t, err.Error(), "incomplete record found")
}

func TestReadSignalFile_InvalidSignalValue(t *testing.T) {
	// 不正なシグナル値を持つレコード
	records := []agent.SignalRecord{
		{Symbol: 7203, Signal: 0x03}, // 0x01, 0x02 以外
	}

	// ダミーファイル作成
	tmpDir := t.TempDir()
	filePath := createDummySignalFile(t, tmpDir, "invalid_signal.bin", records)

	// 読み込み実行
	readRecords, err := agent.ReadSignalFile(filePath)
	assert.Error(t, err)
	assert.Nil(t, readRecords)
	assert.Contains(t, err.Error(), "invalid trade signal value found")
}
