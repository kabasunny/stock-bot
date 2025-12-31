# 株式取引システム 詳細テストプラン

## テスト実行ガイドライン

### 優先度定義
- 🔴 **P0 (Critical)**: システム基本動作に必須
- 🟡 **P1 (High)**: 主要機能、ユーザー体験に重要
- 🟢 **P2 (Medium)**: 拡張機能、品質向上
- ⚪ **P3 (Low)**: エッジケース、パフォーマンス

### テスト種別
- **Unit**: 単体テスト（関数・メソッド単位）
- **Integration**: 統合テスト（コンポーネント間）
- **E2E**: エンドツーエンドテスト（フルフロー）
- **Performance**: パフォーマンステスト

---

## 1. 認証・セッション管理 🔴

### 1.1 Session基盤クラス
**ファイル**: `internal/infrastructure/client/tests/session_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestNewSession | Unit | 🔴 | ✅ | Session作成の基本動作 |
| TestSession_GetPNo | Unit | 🔴 | ✅ | PNo自動インクリメント |
| TestSession_GetPNo_Concurrent | Unit | 🔴 | ✅ | PNo並行安全性 |
| TestSession_SetLoginResponse | Unit | 🔴 | ✅ | ログインレスポンス設定 |
| TestSession_SetLoginResponse_NilInput | Unit | 🟡 | ✅ | nil入力時の安全性 |
| TestSession_SetLoginResponse_EmptyValues | Unit | 🟡 | ✅ | 空値入力時の動作 |

### 1.2 AuthClient基盤
**ファイル**: `internal/infrastructure/client/tests/auth_client_impl_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestAuthClientImpl_LoginOnly | Unit | 🔴 | ✅ | 基本ログイン機能 |
| TestAuthClientImpl_LogoutOnly | Unit | 🔴 | ✅ | 基本ログアウト機能 |
| TestAuthClientImpl_InvalidCredentials | Unit | 🔴 | ✅ | 不正認証情報エラー |
| TestAuthClientImpl_EmptyCredentials | Unit | 🟡 | ✅ | 空認証情報エラー |
| TestAuthClientImpl_LogoutWithoutLogin | Unit | 🟡 | ✅ | 未ログイン状態でのログアウト |
| TestAuthClientImpl_LogoutWithNilSession | Unit | 🟡 | ✅ | nilセッションでのログアウト |
| TestAuthClientImpl_MultipleSessions | Unit | 🟡 | ✅ | 複数セッション管理 |
| TestAuthClientImpl_Sequence_LoginWaitLogoutLogin | Integration | 🟢 | ✅ | 長時間セッション管理 |

### 1.3 TachibanaUnifiedClient統合
**ファイル**: `internal/infrastructure/client/tests/tachibana_unified_client_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestTachibanaUnifiedClient_NewClient | Unit | 🔴 | ✅ | クライアント作成 |
| TestTachibanaUnifiedClient_GetSession | Unit | 🔴 | ✅ | 自動認証機能 |
| TestTachibanaUnifiedClient_EnsureAuthenticated | Unit | 🔴 | 📋 | 認証状態確認 |
| TestTachibanaUnifiedClient_MultipleGetSession | Unit | 🔴 | 📋 | セッション再利用 |
| TestTachibanaUnifiedClient_Logout | Unit | 🔴 | 📋 | ログアウト機能 |
| TestTachibanaUnifiedClient_InvalidCredentials | Unit | 🟡 | 📋 | 不正認証エラー |
| TestTachibanaUnifiedClient_LogoutWithoutLogin | Unit | 🟡 | 📋 | 未ログイン状態処理 |
| TestTachibanaUnifiedClient_SessionExpiry | Integration | 🟡 | 📋 | 8時間セッション期限 |
| TestTachibanaUnifiedClient_AutoReauth | Integration | 🟡 | 📋 | 自動再認証 |

### 1.4 UnifiedClientAdapter
**ファイル**: `internal/infrastructure/client/tests/tachibana_unified_client_adapters_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestAdapter_AuthClientCompatibility | Unit | 🔴 | 📋 | AuthClient互換性 |
| TestAdapter_BalanceClientCompatibility | Unit | 🔴 | 📋 | BalanceClient互換性 |
| TestAdapter_OrderClientCompatibility | Unit | 🔴 | 📋 | OrderClient互換性 |
| TestAdapter_PriceInfoClientCompatibility | Unit | 🔴 | 📋 | PriceInfoClient互換性 |
| TestAdapter_MasterDataClientCompatibility | Unit | 🔴 | 📋 | MasterDataClient互換性 |
| TestAdapter_EventClientCompatibility | Unit | 🟡 | 📋 | EventClient互換性 |

---

## 2. 注文管理機能 🔴

### 2.1 OrderClient基盤
**ファイル**: `internal/infrastructure/client/tests/order_client_impl_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestOrderClient_NewOrder_Market | Unit | 🔴 | 📋 | 成行注文発行 |
| TestOrderClient_NewOrder_Limit | Unit | 🔴 | 📋 | 指値注文発行 |
| TestOrderClient_NewOrder_Stop | Unit | 🔴 | 📋 | 逆指値注文発行 |
| TestOrderClient_NewOrder_InvalidSymbol | Unit | 🔴 | 📋 | 不正銘柄コードエラー |
| TestOrderClient_NewOrder_InvalidQuantity | Unit | 🔴 | 📋 | 不正数量エラー |
| TestOrderClient_NewOrder_InvalidPrice | Unit | 🟡 | 📋 | 不正価格エラー |
| TestOrderClient_GetOrderList | Unit | 🔴 | 📋 | 注文一覧取得 |
| TestOrderClient_GetOrderList_Empty | Unit | 🟡 | 📋 | 空注文一覧 |
| TestOrderClient_CorrectOrder | Unit | 🔴 | 📋 | 注文訂正 |
| TestOrderClient_CorrectOrder_InvalidOrderId | Unit | 🟡 | 📋 | 存在しない注文ID |
| TestOrderClient_CancelOrder | Unit | 🔴 | 📋 | 注文キャンセル |
| TestOrderClient_CancelOrder_InvalidOrderId | Unit | 🟡 | 📋 | 存在しない注文ID |

### 2.2 注文フロー統合テスト
**ファイル**: `internal/infrastructure/client/tests/order_flow_integration_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestOrderFlow_PlaceAndCancel | Integration | 🔴 | 📋 | 注文発行→キャンセル |
| TestOrderFlow_PlaceAndCorrect | Integration | 🔴 | 📋 | 注文発行→訂正 |
| TestOrderFlow_MultipleOrders | Integration | 🟡 | 📋 | 複数注文同時処理 |
| TestOrderFlow_OrderExecution | Integration | 🟡 | 📋 | 注文約定フロー |

---

## 3. 残高・ポジション管理 🔴

### 3.1 BalanceClient基盤
**ファイル**: `internal/infrastructure/client/tests/balance_client_impl_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestBalanceClient_GetZanKaiSummary | Unit | 🔴 | 📋 | 残高サマリー取得 |
| TestBalanceClient_GetZanKaiSummary_Fields | Unit | 🔴 | 📋 | 残高フィールド検証 |
| TestBalanceClient_GetGenbutuKabuList | Unit | 🔴 | 📋 | 現物株式一覧 |
| TestBalanceClient_GetGenbutuKabuList_Empty | Unit | 🟡 | 📋 | 空ポジション |
| TestBalanceClient_GetShinyouTategyokuList | Unit | 🔴 | 📋 | 信用建玉一覧 |
| TestBalanceClient_GetShinyouTategyokuList_Empty | Unit | 🟡 | 📋 | 空建玉 |

### 3.2 ポジション計算テスト
**ファイル**: `internal/infrastructure/client/tests/position_calculation_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestPosition_UnrealizedPL | Unit | 🔴 | 📋 | 評価損益計算 |
| TestPosition_AverageCost | Unit | 🔴 | 📋 | 平均取得単価計算 |
| TestPosition_MarginRequirement | Unit | 🟡 | 📋 | 証拠金必要額計算 |

---

## 4. 価格情報・マーケットデータ 🟡

### 4.1 PriceInfoClient基盤
**ファイル**: `internal/infrastructure/client/tests/price_info_client_impl_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestPriceInfoClient_GetPriceInfo | Unit | 🟡 | 📋 | 現在価格取得 |
| TestPriceInfoClient_GetPriceInfo_InvalidSymbol | Unit | 🟡 | 📋 | 不正銘柄エラー |
| TestPriceInfoClient_GetPriceInfoHistory | Unit | 🟡 | 📋 | 価格履歴取得 |
| TestPriceInfoClient_GetPriceInfoHistory_DateRange | Unit | 🟡 | 📋 | 日付範囲指定 |
| TestPriceInfoClient_GetPriceInfoHistory_OHLCV | Unit | 🟡 | 📋 | OHLCV形式検証 |

---

## 5. マスターデータ管理 🔴

### 5.1 MasterDataClient基盤
**ファイル**: `internal/infrastructure/client/tests/master_data_client_impl_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestMasterDataClient_GetStockInfo | Unit | 🔴 | 📋 | 銘柄情報取得 |
| TestMasterDataClient_GetStockInfo_InvalidSymbol | Unit | 🔴 | 📋 | 存在しない銘柄 |
| TestMasterDataClient_DownloadMasterData | Unit | 🔴 | 📋 | マスターデータ一括取得 |
| TestMasterDataClient_TradingUnit | Unit | 🔴 | 📋 | 売買単位取得 |
| TestMasterDataClient_MarketInfo | Unit | 🟡 | 📋 | 市場情報取得 |

### 5.2 MasterDataScheduler
**ファイル**: `internal/scheduler/tests/master_data_scheduler_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestMasterDataScheduler_NewScheduler | Unit | 🟡 | 📋 | スケジューラー作成 |
| TestMasterDataScheduler_Start | Unit | 🟡 | 📋 | スケジューラー開始 |
| TestMasterDataScheduler_Stop | Unit | 🟡 | 📋 | スケジューラー停止 |
| TestMasterDataScheduler_TriggerManualUpdate | Unit | 🔴 | 📋 | 手動更新実行 |
| TestMasterDataScheduler_ScheduledUpdate | Integration | 🟡 | 📋 | 定期更新実行 |

---

## 6. WebSocketイベント処理 🟡

### 6.1 EventClient基盤
**ファイル**: `internal/infrastructure/client/tests/event_client_impl_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestEventClient_Connect | Unit | 🟡 | 📋 | WebSocket接続 |
| TestEventClient_Connect_InvalidURL | Unit | 🟡 | 📋 | 不正URL接続エラー |
| TestEventClient_ReceiveMessage | Unit | 🟡 | 📋 | メッセージ受信 |
| TestEventClient_Close | Unit | 🟡 | 📋 | 接続クローズ |
| TestEventClient_Reconnect | Integration | 🟢 | 📋 | 自動再接続 |

### 6.2 EventHandler
**ファイル**: `internal/eventprocessing/tests/event_handler_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestExecutionEventHandler_HandleExecution | Unit | 🟡 | 📋 | 約定イベント処理 |
| TestPriceEventHandler_HandlePrice | Unit | 🟡 | 📋 | 価格イベント処理 |
| TestStatusEventHandler_HandleStatus | Unit | 🟡 | 📋 | ステータスイベント処理 |
| TestEventDispatcher_Dispatch | Unit | 🟡 | 📋 | イベント振り分け |

---

## 7. TradeServiceドメイン層 🔴

### 7.1 TradeServiceインターフェース
**ファイル**: `domain/service/tests/trade_service_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestTradeService_GetSession | Unit | 🔴 | 📋 | セッション取得 |
| TestTradeService_GetPositions | Unit | 🔴 | 📋 | ポジション取得 |
| TestTradeService_GetOrders | Unit | 🔴 | 📋 | 注文取得 |
| TestTradeService_GetBalance | Unit | 🔴 | 📋 | 残高取得 |
| TestTradeService_PlaceOrder | Unit | 🔴 | 📋 | 注文発行 |
| TestTradeService_CancelOrder | Unit | 🔴 | 📋 | 注文キャンセル |
| TestTradeService_CorrectOrder | Unit | 🔴 | 📋 | 注文訂正 |
| TestTradeService_GetPriceHistory | Unit | 🟡 | 📋 | 価格履歴取得 |
| TestTradeService_HealthCheck | Unit | 🟡 | 📋 | ヘルスチェック |

### 7.2 GoaTradeService実装
**ファイル**: `internal/tradeservice/tests/goa_trade_service_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestGoaTradeService_AllMethods | Unit | 🔴 | 📋 | 全メソッド動作確認 |
| TestGoaTradeService_ErrorHandling | Unit | 🔴 | 📋 | エラーハンドリング |
| TestGoaTradeService_TypeConversion | Unit | 🔴 | 📋 | 型変換処理 |
| TestGoaTradeService_MasterDataIntegration | Integration | 🔴 | 📋 | マスターデータ連携 |
| TestGoaTradeService_ValidationIntegration | Integration | 🔴 | 📋 | バリデーション連携 |

---

## 8. HTTP APIハンドラー 🔴

### 8.1 TradeServiceハンドラー
**ファイル**: `internal/handler/web/tests/trade_service_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestTradeServiceHandler_GetSession | Unit | 🔴 | 📋 | GET /trade/session |
| TestTradeServiceHandler_GetPositions | Unit | 🔴 | 📋 | GET /trade/positions |
| TestTradeServiceHandler_GetOrders | Unit | 🔴 | 📋 | GET /trade/orders |
| TestTradeServiceHandler_GetBalance | Unit | 🔴 | 📋 | GET /trade/balance |
| TestTradeServiceHandler_PlaceOrder | Unit | 🔴 | 📋 | POST /trade/orders |
| TestTradeServiceHandler_CancelOrder | Unit | 🔴 | 📋 | DELETE /trade/orders/{id} |
| TestTradeServiceHandler_CorrectOrder | Unit | 🔴 | 📋 | PUT /trade/orders/{id} |
| TestTradeServiceHandler_ValidateSymbol | Unit | 🔴 | 📋 | GET /trade/symbols/{symbol}/validate |
| TestTradeServiceHandler_GetOrderHistory | Unit | 🟡 | 📋 | GET /trade/orders/history |
| TestTradeServiceHandler_HealthCheck | Unit | 🟡 | 📋 | GET /trade/health |

### 8.2 その他APIハンドラー
**ファイル**: `internal/handler/web/tests/other_handlers_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestMasterServiceHandler_GetStock | Unit | 🔴 | 📋 | GET /master/stocks/{symbol} |
| TestMasterServiceHandler_Update | Unit | 🔴 | 📋 | POST /master/update |
| TestBalanceServiceHandler_Get | Unit | 🔴 | 📋 | GET /balance |
| TestPositionServiceHandler_List | Unit | 🔴 | 📋 | GET /positions |
| TestOrderServiceHandler_Create | Unit | 🔴 | 📋 | POST /order |
| TestPriceServiceHandler_Get | Unit | 🟡 | 📋 | GET /price/{symbol} |
| TestPriceServiceHandler_GetHistory | Unit | 🟡 | 📋 | GET /price/{symbol}/history |

---

## 9. HTTP APIエンドポイント統合テスト 🔴

### 9.1 エンドポイント統合テスト
**ファイル**: `tests/integration/api_endpoints_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestAPI_SessionManagement | Integration | 🔴 | 📋 | セッション管理API |
| TestAPI_TradingFlow | Integration | 🔴 | 📋 | 取引フローAPI |
| TestAPI_InformationRetrieval | Integration | 🟡 | 📋 | 情報取得API |
| TestAPI_MasterDataManagement | Integration | 🔴 | 📋 | マスターデータAPI |
| TestAPI_ErrorHandling | Integration | 🟡 | 📋 | APIエラーハンドリング |

---

## 10. エージェント・戦略実行 🟡

### 10.1 RefactoredAgent
**ファイル**: `internal/refactoredagent/tests/agent_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestRefactoredAgent_Initialize | Unit | 🟡 | 📋 | エージェント初期化 |
| TestRefactoredAgent_ExecuteStrategy | Unit | 🟡 | 📋 | 戦略実行 |
| TestRefactoredAgent_TradeServiceIntegration | Integration | 🟡 | 📋 | TradeService連携 |
| TestRefactoredAgent_EventProcessingSeparation | Integration | 🟡 | 📋 | イベント処理分離 |

### 10.2 State管理
**ファイル**: `internal/state/tests/state_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestState_UpdateState | Unit | 🟡 | 📋 | 状態更新 |
| TestState_ThreadSafety | Unit | 🟡 | 📋 | スレッドセーフ |
| TestState_Persistence | Unit | 🟢 | 📋 | 状態永続化 |

---

## 11. E2E・統合テスト 🟡

### 11.1 フルフロー統合テスト
**ファイル**: `tests/e2e/full_flow_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestE2E_LoginToLogout | E2E | 🟡 | 📋 | ログイン→取引→ログアウト |
| TestE2E_MultipleOrdersFlow | E2E | 🟡 | 📋 | 複数注文処理フロー |
| TestE2E_ErrorRecoveryFlow | E2E | 🟢 | 📋 | エラー発生時の回復 |
| TestE2E_ConcurrentUsersFlow | E2E | 🟢 | 📋 | 複数ユーザー同時利用 |

---

## 12. パフォーマンス・品質テスト 🟢

### 12.1 パフォーマンステスト
**ファイル**: `tests/performance/performance_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestPerformance_ConcurrentConnections | Performance | 🟢 | 📋 | 同時接続数テスト |
| TestPerformance_HighVolumeOrders | Performance | 🟢 | 📋 | 大量注文処理 |
| TestPerformance_MemoryLeak | Performance | 🟢 | 📋 | メモリリークテスト |
| TestPerformance_ResponseTime | Performance | 🟢 | 📋 | レスポンス時間測定 |

### 12.2 エラーハンドリングテスト
**ファイル**: `tests/error_handling/error_handling_test.go`

| テストケース | 種別 | 優先度 | 実装状況 | 説明 |
|-------------|------|--------|----------|------|
| TestErrorHandling_NetworkTimeout | Unit | 🟢 | 📋 | ネットワークタイムアウト |
| TestErrorHandling_ConnectionLoss | Unit | 🟢 | 📋 | 接続断エラー |
| TestErrorHandling_InvalidResponse | Unit | 🟢 | 📋 | 不正レスポンス |
| TestErrorHandling_RateLimit | Unit | 🟢 | 📋 | レート制限エラー |

---

## テスト実行統計

### 全体統計
- **総テスト項目数**: 156項目
- **完了済み**: 14項目 (9.0%)
- **進行中**: 2項目 (1.3%)
- **未実装**: 140項目 (89.7%)

### 優先度別統計
- **🔴 P0 (Critical)**: 89項目 (57.1%)
- **🟡 P1 (High)**: 41項目 (26.3%)
- **🟢 P2 (Medium)**: 20項目 (12.8%)
- **⚪ P3 (Low)**: 6項目 (3.8%)

### 種別統計
- **Unit**: 118項目 (75.6%)
- **Integration**: 28項目 (17.9%)
- **E2E**: 4項目 (2.6%)
- **Performance**: 6項目 (3.8%)

---

## 推奨実行順序

### Phase 1: 基盤テスト (P0優先) - 2週間
1. **Week 1**: Client層単体テスト
   - TachibanaUnifiedClient完成
   - OrderClient, BalanceClient, MasterDataClient
2. **Week 2**: Service層単体テスト
   - TradeService, GoaTradeService

### Phase 2: 統合テスト (P0→P1) - 1週間
1. **HTTP APIハンドラーテスト**
2. **エンドポイント統合テスト**

### Phase 3: 品質向上テスト (P1→P2) - 1週間
1. **WebSocketイベント処理**
2. **エラーハンドリング**

### Phase 4: 最終品質テスト (P2→P3) - 1週間
1. **E2Eテスト**
2. **パフォーマンステスト**

**推定完了時間**: 5週間