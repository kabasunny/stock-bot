# テスト環境戦略・制約事項

## 🚨 証券会社APIテストの制約事項

### 時間帯制約
- **取引時間**: 平日 9:00-15:00 (昼休み 11:30-12:30除く)
- **夜間メンテナンス**: 平日 17:00-翌8:00 (API停止)
- **週末メンテナンス**: 土日 (API完全停止)
- **祝日**: 市場休場日 (API制限あり)

### 環境制約
- **デモ環境**: 本番同様の時間制約あり
- **電話認証**: 3分間有効 (再認証必要)
- **セッション**: 8時間有効期限
- **レート制限**: 1秒間に数回のAPI呼び出し制限

### データ制約
- **マスターデータ**: 日次更新 (前日夜間)
- **価格データ**: リアルタイム (取引時間のみ)
- **注文データ**: 取引時間外は制限あり

---

## 🎯 テスト戦略の分類

### Layer 1: 環境非依存テスト (24時間実行可能) 🟢
**優先度**: 最高 - 常時実行可能

#### 1.1 純粋単体テスト
- **Session管理**: NewSession, GetPNo, SetLoginResponse
- **型変換**: Request/Response変換ロジック
- **バリデーション**: 入力値検証ロジック
- **計算ロジック**: 損益計算、平均単価計算
- **状態管理**: State更新、スレッドセーフ

#### 1.2 モックテスト
- **HTTPクライアント**: モックレスポンス使用
- **データベース**: インメモリDB使用
- **WebSocket**: モック接続使用

```go
// 例: モックを使用した環境非依存テスト
func TestOrderClient_NewOrder_Mock(t *testing.T) {
    // モックサーバーでレスポンス制御
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"ResultCode": "0", "OrderNumber": "123456"}`))
    }))
    defer mockServer.Close()
    
    // テスト実行...
}
```

### Layer 2: 環境依存・時間制約あり (平日日中のみ) 🟡
**優先度**: 高 - 取引時間内実行

#### 2.1 認証系テスト
- **ログイン/ログアウト**: 実際のAPI認証
- **セッション管理**: 実セッション取得
- **自動再認証**: 8時間期限テスト

#### 2.2 情報取得系テスト
- **残高照会**: GetZanKaiSummary
- **ポジション取得**: GetGenbutuKabuList
- **価格情報**: GetPriceInfo (リアルタイム)
- **マスターデータ**: GetStockInfo

#### 2.3 参照系API統合テスト
- **GET /trade/session**
- **GET /trade/balance**
- **GET /trade/positions**
- **GET /master/stocks/{symbol}**

### Layer 3: 高リスク・実取引影響 (慎重実行) 🔴
**優先度**: 中 - デモ環境推奨

#### 3.1 注文系テスト
- **注文発行**: NewOrder (少額・キャンセル前提)
- **注文訂正**: CorrectOrder
- **注文キャンセル**: CancelOrder

#### 3.2 取引フローテスト
- **POST /trade/orders**
- **PUT /trade/orders/{id}**
- **DELETE /trade/orders/{id}**

---

## 📅 時間帯別テスト実行計画

### 平日 9:00-11:30 (前場)
**実行可能**: Layer 1 + Layer 2 + Layer 3
- 全テスト実行可能
- 注文系テストは少額で実行
- リアルタイム価格データ取得テスト

### 平日 11:30-12:30 (昼休み)
**実行可能**: Layer 1 + Layer 2 (制限あり)
- 認証・情報取得系のみ
- 注文系テストは避ける

### 平日 12:30-15:00 (後場)
**実行可能**: Layer 1 + Layer 2 + Layer 3
- 全テスト実行可能
- 14:30以降は注文系テスト避ける (大引け前)

### 平日 15:00-17:00 (アフター)
**実行可能**: Layer 1 + Layer 2 (制限あり)
- 認証・残高照会は可能
- 価格データは前日終値

### 平日 17:00-翌8:00 (夜間)
**実行可能**: Layer 1のみ
- 環境非依存テストのみ
- モック・単体テスト中心

### 土日・祝日
**実行可能**: Layer 1のみ
- 完全にオフライン
- 開発・リファクタリング時間

---

## 🛠️ テスト実装戦略

### Phase 1: 基盤構築 (環境非依存優先)
**期間**: 1週間 (土日含む24時間実行可能)

1. **モックフレームワーク構築**
   ```go
   // internal/infrastructure/client/tests/mock_server.go
   type MockTachibanaServer struct {
       responses map[string]string
   }
   ```

2. **テストヘルパー拡張**
   ```go
   // 環境判定ヘルパー
   func IsMarketOpen() bool
   func IsDemoEnvironment() bool
   func SkipIfMarketClosed(t *testing.T)
   ```

3. **Layer 1テスト完全実装** (80項目)
   - Session, 型変換, バリデーション
   - 計算ロジック, 状態管理
   - モックベーステスト

### Phase 2: 環境依存テスト (平日日中実行)
**期間**: 1週間 (平日のみ実行)

1. **Layer 2テスト実装** (40項目)
   - 認証系統合テスト
   - 情報取得系統合テスト
   - API統合テスト

2. **時間帯別テスト実行**
   ```bash
   # 平日日中のみ実行
   go test -tags=integration ./...
   
   # 夜間・休日は単体テストのみ
   go test -tags=unit ./...
   ```

### Phase 3: 高リスクテスト (慎重実行)
**期間**: 1週間 (デモ環境推奨)

1. **Layer 3テスト実装** (20項目)
   - 注文系テスト (少額)
   - 取引フローテスト
   - E2Eテスト

---

## 🏷️ テストタグ戦略

### ビルドタグによる分類
```go
//go:build unit
// +build unit

// 環境非依存テスト

//go:build integration
// +build integration

// 環境依存テスト (平日日中のみ)

//go:build e2e
// +build e2e

// 高リスクテスト (デモ環境推奨)
```

### 実行コマンド例
```bash
# 常時実行可能 (CI/CD)
go test -tags=unit ./...

# 平日日中のみ実行
go test -tags=integration ./...

# デモ環境での慎重実行
go test -tags=e2e ./...

# 全テスト実行 (平日日中・デモ環境)
go test -tags="unit integration e2e" ./...
```

---

## ⚠️ テスト実行時の注意事項

### 1. 環境確認
```go
func TestMain(m *testing.M) {
    if !IsMarketOpen() && hasIntegrationTag() {
        log.Println("Market is closed. Skipping integration tests.")
        os.Exit(0)
    }
    os.Exit(m.Run())
}
```

### 2. 電話認証管理
- 3分以内にテスト完了
- 認証失敗時の適切なスキップ
- 複数テスト間での認証共有

### 3. レート制限対応
```go
// テスト間の適切な間隔
time.Sleep(1 * time.Second)
```

### 4. データクリーンアップ
- テスト注文の確実なキャンセル
- テストデータの残存防止

---

## 📊 修正されたテスト進捗計画

### Week 1: 基盤テスト (24時間実行可能)
- **Layer 1テスト**: 80項目
- **モックフレームワーク**: 構築
- **実行時間**: 制約なし

### Week 2: 統合テスト (平日日中のみ)
- **Layer 2テスト**: 40項目
- **実行時間**: 平日 9:00-15:00
- **1日あたり**: 8項目 (5日間)

### Week 3: E2Eテスト (慎重実行)
- **Layer 3テスト**: 20項目
- **実行環境**: デモ環境推奨
- **実行時間**: 平日 10:00-14:00

### Week 4: 品質保証・最適化
- **パフォーマンステスト**: 6項目
- **エラーハンドリング**: 10項目
- **ドキュメント整備**: 完成

---

**重要**: この戦略により、環境制約を考慮した現実的なテスト実装が可能になります。