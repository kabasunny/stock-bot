# マスターデータ改善計画（後回し）

## 概要
マスターデータは低頻度更新（日次1回）で十分。Goaサービスとして適切な設計を実装する。

## 実装予定機能

### 1. 自動更新スケジューラー
- 日次午前2時に自動実行
- `internal/scheduler/master_data_scheduler.go` - 実装済み
- main.goでの統合が必要

### 2. 注文前バリデーション
- 銘柄コードの妥当性チェック
- 売買単位チェック
- TradeServiceでの活用

### 3. API拡張
```
GET /trade/symbols/{symbol}/validate  # 銘柄バリデーション
```

## 実装ファイル（作成済み）
- `internal/scheduler/master_data_scheduler.go`
- 各種修正ファイル（一時的に作成済み）

## 優先度
- 低（注文機能の完成後に実装）
- 現在の手動更新で十分機能する

## 注意事項
- 現在の実装では手動更新（POST /master/update）が利用可能
- 基本的な銘柄マスター取得（GET /master/stocks/{symbol}）は動作中