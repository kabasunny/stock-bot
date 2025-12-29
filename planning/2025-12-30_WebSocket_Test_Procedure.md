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
Invoke-WebRequest -Uri http://localhost:8080/order -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"symbol":"7203","trade_type":"BUY","order_type":"MARKET","quantity":1}'
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