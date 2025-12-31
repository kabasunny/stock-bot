package main

import (
    "encoding/binary"
    "fmt"
    "os"
    "path/filepath"
    "stock-bot/internal/agent" // agentパッケージのTradeSignal, BuySignal, SellSignalを使用
)

func main() {
    // signals ディレクトリが存在しない場合は作成
    signalsDir := "signals"
    if _, err := os.Stat(signalsDir); os.IsNotExist(err) {
        err = os.Mkdir(signalsDir, 0755)
        if err != nil {
            fmt.Printf("Failed to create signals directory: %v\n", err)
            return
        }
    }

    filePath := filepath.Join(signalsDir, "test_signal.bin")
    file, err := os.Create(filePath)
    if err != nil {
        fmt.Printf("Failed to create signal file: %v\n", err)
        return
        }
    defer file.Close()

    // ダミーシグナルデータ
    records := []agent.SignalRecord{
        {Symbol: 7203, Signal: agent.BuySignal},  // トヨタ BUY
        {Symbol: 9984, Signal: agent.SellSignal}, // ソフトバンクG SELL
        {Symbol: 6758, Signal: agent.BuySignal},  // ソニー BUY
    }

    for _, rec := range records {
        err := binary.Write(file, binary.LittleEndian, rec.Symbol)
        if err != nil {
            fmt.Printf("Failed to write symbol: %v\n", err)
            return
        }
        err = binary.Write(file, binary.LittleEndian, rec.Signal)
        if err != nil {
            fmt.Printf("Failed to write signal: %v\n", err)
            return
        }
    }

    fmt.Printf("Successfully created dummy signal file: %s\n", filePath)
}
