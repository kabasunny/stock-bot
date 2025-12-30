# WebSocket 接続テスト手順書 (2025-12-30)

## 1. 目的

本テストの目的は、証券会社のAPIからWebSocket (EVENT I/F) を通じてリアルタイムに送信される、**約定通知**、**時価情報**、**ステータス通知**などのメッセージの具体的なフォーマット（キーと値のペア）を特定することです。

特に、これまで未確認であった以下の`p_cmd`に対応するメッセージの生データを収集します。
- 約定通知 (`EX`)
- 時価情報 (`FD`)
- ステータス通知 (`ST`)

収集したデータは、`agent.go`内のプレースホルダー関数 `handleExecution`, `handlePriceData`, `handleStatus` の具体的な実装に不可欠な情報となります。

### 1.1. 本テストの戦略的整理

本テストは、個別の関数を検証する**ユニットテスト**ではなく、システム全体を動かして外部APIとの連携を確認する**結合テスト**です。

現在の最優先事項は、実際のAPIがどのようなデータを返すかという「正解」を知ることです。この「正解データ」を収集しない限り、堅牢なユニットテストを作成することができません。

したがって、以下の戦略的な順序で開発を進めます。
1.  **【今回】結合テストの実施**: まず、本番に近い環境で実際にAPIと通信し、生の「お手本データ」を確実に収集する。
2.  **【今後】モックベースのユニットテスト構築**: 収集したデータを基に、APIの応答を模倣する「モック」を作成し、それを利用して高速かつ安定したユニットテストを再構築・拡充する。

これにより、システムの信頼性と開発効率の両方を最大化します。


## 2. 準備

1.  **コードの確認**: `internal/infrastructure/client/event_client_impl.go` が更新され、WebSocketから受信した全てのメッセージを生データ（raw aessage）のままログ出力する機能が追加されていることを確認してください。具体的には、ログに `Received raw WebSocket message` というメッセージが出力されるようになっています。

2.  **環境の起動**:
    - Docker環境を起動します。
      ```shell
      docker compose up -d
      ```
    - データベースのマイグレーションを実行します。
      ```shell
      go run ./cmd/migrator/main.go
      ```

## 3. 実行手順

### ステップ1: アプリケーションの起動

ターミナルを開き、以下のコマンドでアプリケーションサーバーを起動します。`--skip-sync`フラグを付与することで、起動時の初期同期をスキップし、迅速にWebSocketの待ち受けを開始します。

```shell
go run ./cmd/myapp/main.go --skip-sync
```

サーバーが正常に起動し、ログに `Connecting to WebSocket` および `Successfully connected to WebSocket` と表示されることを確認してください。

### ステップ2: リアルタイムイベントの監視

アプリケーションを起動したターミナルのログ出力を継続的に監視します。WebSocketサーバーから何らかのメッセージ（時価情報など）が送られてくると、以下のようなログが出力されるはずです。

```json
INFO Received raw WebSocket message message="p_cmd^BST^Ap_no^B1^A..."
```

### ステップ3: 約定通知イベントの生成

約定通知メッセージを意図的に発生させるため、別のターミナルを開き、`Invoke-WebRequest`（PowerShell）または`curl`（bash/zsh）を使用して、実際に市場で取引されている銘柄に対して**ごく少量の**成行注文を発行します。

**【注意】** 実際に費用が発生する操作です。テストする銘柄コード（`symbol`）、数量（`quantity`）には十分に注意してください。

**PowerShellの例:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/order -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"symbol":"7203","trade_type":"BUY", "order_type":"MARKET", "quantity":1}'
```

**curlの例:**
```shell
curl -X POST -H "Content-Type: application/json" -d '{"symbol":"7203", "trade_type":"BUY", "order_type":"MARKET", "quantity":1}' http://localhost:8080/order
```

注文が成功すると、数秒〜数十秒以内に取引所で約定し、その結果がWebSocket経由で通知されるはずです。

## 4. 収集すべき情報

ステップ1のアプリケーションログから、`Received raw WebSocket message` というログをすべて収集・保存してください。

特に、以下の`p_cmd`の値を持つメッセージが重要です。

-   `p_cmd^BEX`: **（最重要）** 約定通知。注文ID、約定価格、数量、手数料などの情報が含まれると予想されます。
-   `p_cmd^BFD`: 時価配信。現在値、気配値などの情報が含まれると予想されます。
-   `p_cmd^BST`: ステータス通知。接続状態などを示すメッセージと予想されます。
-   その他、予期しない`p_cmd`を持つすべてのメッセージ。

収集した生のメッセージ文字列は、後の解析のためにテキストファイル等にまとめてください。

## 5. テスト後の作業

収集したログデータを基に、`p_cmd`ごとのメッセージ構造を解析し、`agent.go`の`handleExecution`、`handlePriceData`、`handleStatus`関数を実装します。これにより、エージェントがリアルタイムイベントに正しく反応できるようになります。

---

## 6. 本日の開発進捗と成果 (2025-12-30)

本日のWebSocket接続テストおよび関連機能の実装において、以下の課題を特定し、解決または進展させました。

### 6.1. 解決済みの課題と実装された機能

1.  **APIログインエラーの詳細化**:
    *   **課題**: ログイン失敗時に`result code`が空で、原因特定が困難でした。
    *   **解決**: `internal/infrastructure/client/auth_client_impl.go`を修正し、APIからの生レスポンス全体をログに出力するように変更。これにより、具体的なエラーコード`10089`とメッセージ「電話認証後、3分以内にログインしてください。」を特定できました。

2.  **WebSocket接続エラーの解消**:
    *   **課題1**: `malformed ws or wss URL` エラーが発生していました。
    *   **解決1**: `internal/infrastructure/client/event_client_impl.go`を修正し、`https://`スキーマを`wss://`に変換するロジックを追加しました。
    *   **課題2**: `websocket: bad handshake` エラーが発生していました。
    *   **解決2**:
        *   ログインレスポンスから`sUrlEventWebSocket`を正しく取得するため、`internal/infrastructure/client/dto/auth/response/login.go`の`ResLogin`構造体に`SUrlEventWebSocket`フィールドを追加。
        *   `internal/infrastructure/client/session.go`の`SetLoginResponse`メソッドで、`session.EventURL`に`res.SUrlEventWebSocket`の値を設定するように修正。
        *   `fmt`パッケージのインポート漏れを修正(`internal/infrastructure/client/event_client_impl.go`)。

3.  **WebSocket接続後の「parameter error」解消とデータ受信開始**:
    *   **課題**: 接続後すぐに`close 1000 (normal)`エラーが発生し、WebSocketからのメッセージが受信されませんでした。
    *   **解決**: Pythonサンプルコード(`e_api_websocket_receive_tel.py`)の分析に基づき、`internal/infrastructure/client/event_client_impl.go`の`Connect`メソッドで、WebSocket URLに以下の必須クエリパラメータを追加しました。
        *   `p_rid=22`, `p_board_no=1000`, `p_eno=0`
        *   `p_evt_cmd=ST,KP,FD,EC`
        *   `p_gyou_no=` (銘柄数に応じた連番), `p_mkt_code=` (銘柄数に応じた`00`)
        *   これにより、APIからの`KP` (キープアライブ) メッセージ受信が開始され、接続が維持されるようになりました。

4.  **`KP` (キープアライブ) イベントの警告解消**:
    *   **課題**: `KP`イベントに対して`unhandled websocket event command`の警告が出ていました。
    *   **解決**: `internal/agent/agent.go`の`watchEvents`メソッドの`switch cmd`ステートメントに`KP`ケースを追加し、Debugログを出力するように修正。

5.  **`FD` (時価配信) イベントの受信と処理の基盤構築**:
    *   **課題1**: `FD`イベントが受信されず、エージェントがリアルタイム価格を取得できませんでした。
    *   **解決1**: `agent_config.yaml`の`target_symbols`をPythonサンプルコードで使用されている`'8697'`に変更したところ、`FD`イベントの受信に成功しました。
    *   **課題2**: `FD`イベントから価格を抽出し、エージェントの内部状態に反映するメカニズムがありませんでした。
    *   **解決2**:
        *   `internal/agent/state.go`に`prices`マップと`UpdatePrice`, `GetPrice`メソッドを追加し、銘柄ごとの現在価格を管理できるようにしました。
        *   `internal/agent/agent.go`の`Agent`構造体に`gyouNoToSymbol`マップを追加し、WebSocketイベントの行番号から銘柄コードを特定できるようにしました。
        *   `internal/agent/agent.go`の`handlePriceData`関数を実装し、`FD`イベントデータ(`p_N_DPP`)から価格を抽出して`State`を更新するようにしました。

6.  **エージェントの価格情報参照ロジックの改善**:
    *   **課題**: エージェントが`agent_config.yaml`の`target_symbols`に関わらず、シグナルファイルに記載された銘柄の価格をREST API経由で取得しようとして失敗していました。
    *   **解決**:
        *   `internal/agent/agent.go`の`checkSignalsForEntry`関数を修正し、`agent_config.yaml`の`target_symbols`に含まれる銘柄のシグナルのみを処理するようにフィルタリング。
        *   `checkSignalsForEntry`および`checkPositionsForExit`関数内で、価格取得を`a.tradeService.GetPrice`から`a.state.GetPrice`に変更し、WebSocketからのリアルタイム価格を利用するようにしました。

7.  **`EC` (約定通知) イベントの受信と処理の基盤構築**:
    *   **課題1**: `EC`イベントが受信されない問題が一時的に発生していましたが、`p_evt_cmd=EC`のみを購読することでデモ環境でも受信できることを確認しました。
    *   **課題2**: `EC`イベントを受信した際に、データベースの`executions`テーブルに`symbol`や`trade_type`カラムがない、`total executed quantity exceeds order quantity`といったデータベースエラーが発生しました。
    *   **解決2**:
        *   `executions`テーブルに`symbol`、`trade_type`カラムを追加し、既存の`execution_time`、`execution_price`、`execution_quantity`カラムを`executed_at`、`price`、`quantity`にリネームするマイグレーションを適用しました。
        *   `internal/agent/agent.go`の`handleExecution`関数を修正し、`EC`イベントのキー名(`p_ON`, `p_IC`, `p_ST`, `p_NT`, `p_EXPR`, `p_EXDT`)に合わせて`model.Execution`に正しくマッピングするようにしました。
        *   `internal/agent/agent.go`の`parseTime`ヘルパー関数を修正し、`YYYYMMDDhhmmss`形式のタイムスタンプもパースできるようにしました。
        *   `internal/infrastructure/repository/order_repository_impl.go`の`UpdateOrderStatusByExecution`関数を修正し、約定数量の重複計上を防ぎ、`EC`イベントの約定を正しくデータベースに永続化するようにしました。

8.  **`EC`イベント処理時のログレベル調整**:
    *   **課題**: データベースをクリーンアップした後にWebSocketから過去の`EC`イベントが再送されると、データベースに存在しない注文IDに対して`order with ID ... not found`エラーが大量に発生していました。
    *   **解決**: `internal/agent/agent.go`の`handleExecution`関数のシグネチャを`error`を返すように変更し、エラーパスで`errors.New`または`errors.Wrapf`を使うように修正しました。また、`watchEvents`メソッド内の`case "EC":`ブロックで、`handleExecution`が返すエラーを処理し、`order with ID ... not found`エラーの場合は`WARN`レベルでログを出力するようにしました。
    *   **成果**: これにより、過去の注文に対する`EC`イベントによるログは`WARN`レベルに抑制され、本当に重要なエラーのみが`ERROR`レベルで表示されるようになり、ログの可読性が大幅に向上しました。

### 6.2. 今後の課題

*   **`ST` (ステータス通知) イベントの処理**: `handleStatus`プレースホルダーを実装し、APIからの重要なステータス変更（例: エラー、セッション無効化）を処理する必要があります。
*   **REST APIによる価格取得の廃止**: `GoaTradeService.GetPrice`のようなREST API経由の現在価格取得は不要になったため、コードのクリーンアップを検討。

---

## 7. 追加の開発進捗と成果 (2025-12-30 - 継続)

これまでの作業で、いくつかの重要な課題が解決され、システムの安定性が向上しました。

### 7.1. 解決済みの課題と実装された機能

1.  **約定数量超過エラー (`total executed quantity exceeds order quantity`) の解決**:
    *   **課題**: `model.Execution`の`Quantity`が`EC`イベントの`p_NT`から誤ってパースされ、合計約定数量が注文数量を超過するエラーが発生していました。
    *   **解決**: `internal/agent/agent.go`の`handleExecution`関数を修正し、`model.Execution.Quantity`に`p_NT`ではなく、実際の約定数量を示す`p_EXSR`を使用するように変更しました。
    *   **成果**: データベースをクリーンアップ (`go run ./cmd/migrator/main.go`) し、`agent_config.yaml`に新しい`target_symbols`を追加した後の**新規注文では約定数量超過エラーが発生しないことを確認しました**。ログに出ていた過去の約定に関するエラーは、データベースクリア前の古いイベントに対するものです。

2.  **価格情報取得失敗エラー (`failed to get price for exit check from state`) への対応 (部分的に解決)**:
    *   **課題**: エージェントが`a.tradeService.GetPrice`を呼び出し、`a.state`に価格情報が存在しない場合にエラーが発生していました。
    *   **解決**:
        *   `internal/agent/agent.go`の`checkPositionsForExit`関数を修正し、`a.tradeService.GetPrice`の代わりに`a.state.GetPrice`を使用するように変更しました。
        *   `internal/agent/trade_service.go`インターフェースおよび`internal/agent/goa_trade_service.go`実装から、冗長な`GetPrice`メソッドを削除しました。
        *   `agent_config.yaml`の`StrategySettings.Swingtrade.TargetSymbols`に、エラーが出ていた全ての銘柄（`6504`, `6505`, `9001`, `6501`, `6502`など）を追加し、WebSocket経由で`FD`イベントを購読するように設定しました。
    *   **成果**: `target_symbols`がWebSocket購読に正しく反映されていることは確認できましたが、**市場が閉まっている（昼休み中）ため、APIから`FD`イベントが送信されず、`a.state.prices`が更新されない状態が続いています。** これにより、「`failed to get price for exit check from state`」エラーが引き続き発生しています。これは市場が開場するまで完全な検証はできません。

3.  **`ST` (ステータス通知) イベントのログレベル調整と内容の記録**:
    *   **課題**: `handleStatus`がプレースホルダーの実装であったため、`ST`イベントの内容が把握しにくい状況でした。
    *   **解決**: `internal/agent/agent.go`の`handleStatus`関数を修正し、受信した`ST`イベントの全データを`WARN`レベルでログに出力するように変更しました。
    *   **成果**: 実際に`p_err:database i/o error.`という内容の`ST`イベントが確認され、その内容がログに記録されることを確認しました。このエラー自体の意味は不明ですが、今後の分析のための基盤は整いました。

### 7.2. 残された課題と次の検証ステップ

1.  **口座区分不一致エラー (`errno: 11481`) の解決**:
    *   **課題**: 売却注文時、「`選択した口座区分がお預かり銘柄と不一致のため、このご注文はお受けできません。`」というエラーが発生していました。これは、`PlaceOrder`が常に`GenkinShinyouKubun: "0"`（現物）をAPIに送信していたため、信用取引で保有しているポジション（`6504`, `6505`, `9001`など）を売却しようとすると口座区分不一致となったためです。
    *   **解決**:
        *   `domain/model/position.go`に`PositionAccountType`フィールドを追加し、ポジションが現物 (`CASH`) か信用 (`MARGIN`) かを区別できるようにしました。
        *   データベースに`position_account_type`カラムを追加する新しいマイグレーションを適用しました。
        *   `internal/agent/trade_service.go`の`PlaceOrderRequest`に`PositionAccountType`フィールドを追加しました。
        *   `internal/agent/goa_trade_service.go`の`GoaTradeService.GetPositions`を修正し、APIからポジションを取得する際に`PositionAccountType`を現物 (`CASH`) または信用 (`MARGIN`) に設定するようにしました。
        *   `internal/agent/goa_trade_service.go`の`GoaTradeService.PlaceOrder`を修正し、`PlaceOrderRequest`の`PositionAccountType`に基づいて`GenkinShinyouKubun`を`"0"`（現物）または`"1"`（信用）に設定するようにしました。
        *   `internal/agent/agent.go`の`placeExitOrder`および`checkSignalsForEntry`の`PlaceOrder`呼び出し箇所を修正し、`PlaceOrderRequest`に`model.Position`から取得した`PositionAccountType`を渡すようにしました（買い注文の場合は`model.PositionAccountTypeCash`をデフォルトとしました）。
    *   **成果**: これにより、口座区分不一致エラーは解消されると期待されます。

2.  **銘柄市場マスタデータエラー (`errno: 11108`)**:
    *   **課題**: 銘柄`6502`の売却注文時に「`銘柄市場マスタにデータがありません`」というエラーが発生しました。
    *   **現状**: このエラーは、APIが銘柄`6502`を認識していないことを示唆しています。これはエージェント側のコードで直接解決できる問題ではなく、銘柄コードが有効であるか、APIの提供範囲内であるかを確認する必要があります。

3.  **市場開場後の価格情報 (`FD`) および歴史的価格 (`GetPriceHistory`) の検証**:
    *   現在発生している「`failed to get price for exit check from state`」エラーは、市場が閉まっているためにAPIから時価情報（`FD`イベント）が提供されていないことが原因と強く推測されます。同様に、`GoaTradeService.GetPriceHistory`も現在データを提供していません。
    *   **市場が開場した後、再度アプリケーションを起動し、これらのエラーが解消されているかを確認する必要があります。**

4.  **`ST`イベントの内容解析と適切な処理の実装**:
    *   「`p_err:database i/o error.`」という内容の`ST`イベントが確認されました。これが何を意味し、エージェントがどのように反応すべきか（例: セッションの再確立、特定のエラー通知）は、今後の分析と、APIドキュメント（もしあれば）の参照によって決定する必要があります。

**結論**:
最もクリティカルであった約定数量超過エラーと、口座区分不一致エラーに対するコード側の修正は完了しました。価格情報に関するエラーは市場の状況に依存するため、開場を待つ必要があります。銘柄`6502`のエラーはAPI側の問題の可能性があります。これらの変更は、エージェントの安定性と正確性を大きく向上させました。