# 株式取引システム 全体フロー詳細ドキュメント

## 概要

本ドキュメントは、株式取引システムの全体的な処理フローを時系列順に詳細に説明します。
システム起動から取引実行、エラーハンドリング、システム終了まで、すべての主要フローを網羅しています。

---

## 1. システム起動フロー

### 1.1 アプリケーション初期化

```
[起動] cmd/myapp/main.go
    ↓
[設定] .env ファイル読み込み
    ↓
[DI] internal/infrastructure/container 初期化
    ↓
[DB] PostgreSQL 接続確立
    ↓
[Migration] データベーススキーマ更新
    ↓
[API Client] 立花証券APIクライアント初期化
    ↓
[Session] セッション管理初期化
    ↓
[Goa] サービス層・ハンドラー初期化
    ↓
[HTTP] サーバー起動 (ポート8080)
```

### 1.2 詳細処理

| ステップ | 処理内容 | 関連ファイル | 備考 |
|---------|----------|-------------|------|
| 設定読み込み | 環境変数・認証情報取得 | `.env` | 立花API認証情報含む |
| DIコンテナ | 依存関係注入設定 | `internal/infrastructure/container/` | 全サービス初期化 |
| DB接続 | PostgreSQL接続確立 | `compose.yaml` | ヘルスチェック実行 |
| マイグレーション | テーブル作成・更新 | `migrations/` | 自動実行 |
| APIクライアント | 立花証券API接続準備 | `internal/infrastructure/client/` | 認証未実行 |
| セッション管理 | 日付ベース管理初期化 | `date_based_session_manager.go` | ファイルベース |
| Goaサービス | REST APIエンドポイント準備 | `internal/handler/web/` | 全エンドポイント有効化 |

---

## 2. ログイン・認証フロー

### 2.1 認証シーケンス

```
[Client] GET /trade/session
    ↓
[Goa Handler] リクエスト受信・バリデーション
    ↓
[Trade Service] GetSession() メソッド呼び出し
    ↓
[Session Manager] 認証状態確認
    ↓
[判定] セッション有効性チェック
    ├─ [有効] → 既存セッション返却
    └─ [無効] → 新規ログイン実行
        ↓
    [Auth Client] 立花証券APIログインリクエスト
        ↓
    [立花API] 認証情報検証・セッション発行
        ↓
    [File Save] セッション情報ファイル保存
        ↓
    [Response] JSON形式でクライアントに返却
```

### 2.2 セッション管理詳細

#### セッションファイル構造
```json
{
  "session": {
    "session_id": "ABC123...",
    "user_id": "ycu05273",
    "login_time": "2026-01-02T09:00:00Z",
    "p_no": 12345
  },
  "date": "20260102",
  "created_at": "2026-01-02T09:00:00Z",
  "last_used_at": "2026-01-02T15:30:00Z"
}
```

#### ファイル保存場所
- **パス**: `./data/sessions/tachibana_session_YYYYMMDD.json`
- **権限**: `0600` (所有者のみ読み書き可能)
- **ローテーション**: 営業日ごとに自動作成

---

## 3. 注文発行フロー

### 3.1 注文処理シーケンス

```
[Client] POST /trade/orders
    ↓
[Request Body] JSON形式注文データ
    {
      "symbol": "1301",
      "trade_type": "BUY",
      "order_type": "STOP_LIMIT",
      "quantity": 100,
      "price": 1500,
      "trigger_price": 1480
    }
    ↓
[Goa Handler] パラメータ検証・型変換
    ↓
[Trade Service] PlaceOrder() メソッド呼び出し
    ↓
[Session Manager] 認証状態確認・セッション取得
    ↓
[Order Client] 注文パラメータ変換
    ↓
[Parameter Mapping] システム形式 → 立花API形式
    ↓
[立花API] 注文リクエスト送信
    ↓
[Response Parse] 注文受付結果解析
    ↓
[Database] orders テーブルに注文情報保存
    ↓
[Client Response] 注文結果をJSON形式で返却
```

### 3.2 注文タイプ別パラメータマッピング

#### STOP_LIMIT注文の場合
```go
// システム内部形式
{
  "price": 1500,         // 指値価格
  "trigger_price": 1480  // 逆指値価格
}

// 立花API形式
{
  "OrderPrice": 1500,      // 指値価格
  "GyakusasiPrice": 1480   // 逆指値価格
}
```

#### 成行注文の場合
```go
// システム内部形式
{
  "price": null
}

// 立花API形式
{
  "OrderPrice": 0
}
```

### 3.3 データベース保存

```sql
INSERT INTO orders (
    order_id, symbol, trade_type, order_type,
    quantity, price, trigger_price, 
    position_account_type, order_status,
    created_at
) VALUES (
    'ORD123456', '1301', 'BUY', 'STOP_LIMIT',
    100, 1500.0, 1480.0,
    'CASH', 'PENDING',
    NOW()
);
```

---

## 4. ポジション・残高照会フロー

### 4.1 ポジション取得シーケンス

```
[Client] GET /trade/positions
    ↓
[Goa Handler] リクエスト受信
    ↓
[Trade Service] GetPositions() メソッド呼び出し
    ↓
[Session Manager] セッション有効性確認
    ↓
[Position Client] 立花証券API呼び出し
    ├─ GetGenbutuKabuList() (現物株式)
    └─ GetShinyouTategyokuList() (信用建玉)
    ↓
[Data Transform] 立花API形式 → システム内部形式
    ↓
[Response] JSON形式でポジション一覧返却
```

### 4.2 残高取得シーケンス

```
[Client] GET /trade/balance
    ↓
[Balance Service] GetBalance() メソッド呼び出し
    ↓
[Balance Client] GetZanKaiSummary() API呼び出し
    ↓
[Data Parse] 残高情報解析
    ├─ 現金残高
    ├─ 買付余力
    ├─ 評価額
    └─ 損益
    ↓
[Response] 残高情報をJSON形式で返却
```

---

## 5. 価格情報取得フロー

### 5.1 価格履歴取得シーケンス

```
[Client] GET /trade/price-history/{symbol}
    ↓
[Goa Handler] パスパラメータ抽出
    ↓
[Price Service] GetPriceHistory() メソッド呼び出し
    ↓
[Price Client] 立花証券API価格情報リクエスト
    ↓
[立花API] 価格履歴データ取得
    ↓
[Data Format] OHLCV形式に整形
    {
      "date": "2026-01-02",
      "open": 1450.0,
      "high": 1520.0,
      "low": 1440.0,
      "close": 1500.0,
      "volume": 150000
    }
    ↓
[Response] 価格履歴をJSON配列で返却
```

### 5.2 リアルタイム価格取得

```
[Client] GET /trade/price/{symbol}
    ↓
[Price Service] GetCurrentPrice() メソッド呼び出し
    ↓
[Price Client] GetPriceInfo() API呼び出し
    ↓
[Cache Check] 価格キャッシュ確認
    ├─ [Hit] → キャッシュから返却
    └─ [Miss] → API呼び出し実行
    ↓
[Response] 現在価格情報返却
```

---

## 6. マスターデータ同期フロー

### 6.1 マスターデータ更新シーケンス

```
[Trigger] システム起動時 OR POST /master/update
    ↓
[Master Service] UpdateMasterData() メソッド呼び出し
    ↓
[Master Client] 立花証券APIマスターデータ取得
    ├─ 銘柄マスター取得
    ├─ 呼値ルール取得
    └─ 市場情報取得
    ↓
[Database Transaction] マスターデータ更新
    ├─ stock_masters テーブル更新
    ├─ tick_rules テーブル更新
    ├─ tick_levels テーブル更新
    └─ stock_market_masters テーブル更新
    ↓
[Commit] トランザクションコミット
    ↓
[Response] 更新完了レスポンス返却
```

### 6.2 更新対象テーブル

| テーブル | 内容 | 更新頻度 |
|---------|------|----------|
| `stock_masters` | 銘柄基本情報 | 日次 |
| `tick_rules` | 呼値ルール | 不定期 |
| `tick_levels` | 価格帯別呼値 | 不定期 |
| `stock_market_masters` | 市場別銘柄情報 | 日次 |

---

## 7. WebSocketイベント処理フロー

### 7.1 イベント処理シーケンス

```
[WebSocket] 立花証券イベントAPI接続
    ↓
[Event Receive] 約定・価格更新イベント受信
    ↓
[Event Dispatcher] イベント種別判定
    ├─ [約定] → ExecutionEventHandler
    ├─ [価格] → PriceEventHandler
    └─ [ステータス] → StatusEventHandler
    ↓
[Database Update] 対応テーブル更新
    ├─ executions テーブル (約定情報)
    ├─ orders テーブル (注文ステータス)
    └─ price_cache (価格キャッシュ)
    ↓
[Event Complete] イベント処理完了
```

### 7.2 イベント種別と処理内容

| イベント種別 | 処理内容 | 更新対象 |
|-------------|----------|----------|
| 約定通知 | 約定情報をDBに保存 | `executions` テーブル |
| 価格更新 | 価格キャッシュ更新 | メモリキャッシュ |
| 注文ステータス | 注文状態更新 | `orders` テーブル |
| システム通知 | ログ出力 | ログファイル |

---

## 8. エラーハンドリング・復旧フロー

### 8.1 エラー処理シーケンス

```
[Error Occur] APIエラー・ネットワーク障害等発生
    ↓
[Error Handler] エラー種別判定
    ├─ [認証エラー] → セッション再取得
    ├─ [ネットワークエラー] → リトライ処理
    ├─ [バリデーションエラー] → エラーレスポンス
    └─ [システムエラー] → ログ出力・アラート
    ↓
[Recovery Action] 復旧処理実行
    ↓
[Log Output] 構造化ログでエラー詳細記録
    ↓
[Client Response] エラーレスポンス返却
```

### 8.2 エラー種別と対応

| エラー種別 | 原因 | 対応処理 | リトライ |
|-----------|------|----------|---------|
| 認証エラー | セッション期限切れ | 自動再ログイン | あり (3回) |
| ネットワークエラー | 通信障害 | 指数バックオフリトライ | あり (5回) |
| バリデーションエラー | 不正パラメータ | エラーレスポンス返却 | なし |
| レート制限エラー | API呼び出し過多 | 待機後リトライ | あり (無制限) |
| システムエラー | 内部処理エラー | ログ出力・アラート | なし |

### 8.3 ログ出力形式

```json
{
  "timestamp": "2026-01-02T15:30:00Z",
  "level": "ERROR",
  "message": "API call failed",
  "error": {
    "type": "NetworkError",
    "code": "TIMEOUT",
    "details": "Connection timeout after 30s"
  },
  "context": {
    "user_id": "ycu05273",
    "session_id": "ABC123...",
    "api_endpoint": "/api/orders",
    "retry_count": 2
  }
}
```

---

## 9. セッション管理・日次切り替えフロー

### 9.1 日次切り替えシーケンス

```
[Date Check] 営業日変更検知
    ↓
[Session Invalidate] 古いセッション無効化
    ↓
[File Create] 新しいセッションファイル作成
    ↓
[Re-login] 立花証券APIに再ログイン
    ↓
[Session Save] 新セッション保存・アクティブ化
    ↓
[Cleanup] 古いセッションファイル削除 (7日以上前)
```

### 9.2 営業日判定ロジック

```go
func getCurrentBusinessDate() string {
    now := time.Now().In(jst)
    
    // 土日は前営業日
    if now.Weekday() == time.Saturday {
        now = now.AddDate(0, 0, -1)
    } else if now.Weekday() == time.Sunday {
        now = now.AddDate(0, 0, -2)
    }
    
    // 祝日判定（簡易版）
    // 実際は祝日カレンダーAPIと連携
    
    return now.Format("20060102")
}
```

---

## 10. システム終了フロー

### 10.1 Graceful Shutdown シーケンス

```
[Signal] SIGTERM/SIGINT シグナル受信
    ↓
[Shutdown Start] Graceful Shutdown開始
    ↓
[HTTP Stop] HTTPサーバー停止 (既存リクエスト完了待ち)
    ↓
[WebSocket Close] WebSocket接続切断
    ↓
[API Logout] 立花証券APIからログアウト
    ↓
[DB Close] データベース接続切断
    ↓
[Resource Cleanup] リソース解放
    ↓
[System Exit] システム終了 (exit code 0)
```

### 10.2 終了処理詳細

| ステップ | 処理内容 | タイムアウト |
|---------|----------|-------------|
| HTTPサーバー停止 | 既存リクエスト完了待ち | 5秒 |
| WebSocket切断 | イベント接続クリーンアップ | 3秒 |
| APIログアウト | セッション無効化 | 10秒 |
| DB接続切断 | コネクションプール解放 | 5秒 |
| リソース解放 | メモリ・ファイルハンドル解放 | 2秒 |

---

## 11. パフォーマンス・監視

### 11.1 パフォーマンス指標

| 指標 | 目標値 | 測定方法 |
|------|--------|----------|
| API応答時間 | < 500ms | HTTPミドルウェア |
| 同時接続数 | > 100 | 負荷テスト |
| スループット | > 1000 req/sec | ベンチマーク |
| メモリ使用量 | < 2MB | プロファイリング |
| DB接続数 | < 10 | 接続プール監視 |

### 11.2 監視項目

- **ヘルスチェック**: `/trade/health` エンドポイント
- **メトリクス**: Prometheus形式でメトリクス出力
- **ログ**: 構造化ログ (JSON形式)
- **アラート**: エラー率・応答時間閾値監視

---

## 12. セキュリティ

### 12.1 セキュリティ対策

| 項目 | 対策内容 |
|------|----------|
| 認証情報 | 環境変数で管理・平文保存禁止 |
| セッション | ファイル権限600・暗号化推奨 |
| 通信 | HTTPS必須・証明書検証 |
| ログ | 機密情報マスキング |
| API | レート制限・入力値検証 |

### 12.2 セキュリティチェックリスト

- [ ] 認証情報の環境変数化
- [ ] セッションファイルの暗号化
- [ ] HTTPS通信の強制
- [ ] ログの機密情報マスキング
- [ ] APIレート制限の実装
- [ ] 入力値検証の強化

---

## 付録

### A. 関連ファイル一覧

| カテゴリ | ファイルパス | 説明 |
|---------|-------------|------|
| エントリーポイント | `cmd/myapp/main.go` | メインアプリケーション |
| 設定 | `.env` | 環境変数設定 |
| DIコンテナ | `internal/infrastructure/container/` | 依存関係注入 |
| ハンドラー | `internal/handler/web/` | HTTP APIハンドラー |
| サービス層 | `internal/tradeservice/` | ビジネスロジック |
| クライアント | `internal/infrastructure/client/` | 外部API接続 |
| マイグレーション | `migrations/` | データベーススキーマ |

### B. API エンドポイント一覧

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/trade/session` | セッション情報取得 |
| GET | `/trade/positions` | ポジション一覧取得 |
| GET | `/trade/orders` | 注文一覧取得 |
| POST | `/trade/orders` | 注文発行 |
| DELETE | `/trade/orders/{id}` | 注文キャンセル |
| GET | `/trade/balance` | 残高情報取得 |
| GET | `/trade/price-history/{symbol}` | 価格履歴取得 |
| POST | `/master/update` | マスターデータ更新 |
| GET | `/trade/health` | ヘルスチェック |

### C. データベーステーブル一覧

| テーブル名 | 説明 | 主要カラム |
|-----------|------|-----------|
| `orders` | 注文情報 | order_id, symbol, trade_type, order_type |
| `executions` | 約定情報 | execution_id, order_id, execution_price |
| `positions` | ポジション情報 | symbol, position_type, quantity |
| `stock_masters` | 銘柄マスター | issue_code, issue_name, trading_unit |
| `signals` | 取引シグナル | symbol, signal_type, generated_at |

---

**ドキュメント作成日**: 2026-01-02  
**バージョン**: 1.0  
**最終更新**: システム全体フロー詳細化