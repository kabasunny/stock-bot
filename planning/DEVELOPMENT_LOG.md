## 開発進捗（2025-12-22）

### 発注後のエージェント内部状態更新
- **状態更新ロジックの実装**: `agent.go`の`tick`メソッドにおいて、`tradeService.PlaceOrder`が成功した直後に、返された`order`オブジェクトをエージェントの内部状態（`a.state`）に即座に追加する処理を実装しました。
- **スレッドセーフな状態変更**: `state.go`に、Mutex（ロック）を利用して安全に単一の注文情報を追加するための`AddOrder`メソッドを新規に実装しました。
- **目的達成**: これにより、エージェントは自身で発行した注文（未約定注文）を次の`tick`を待たずに即座に認識し、その後の意思決定（例: 同一銘柄への連続発注防止など）に正しく反映できるようになりました。

---
# 開発ログ

**注意**: このファイルは、プロジェクトの詳細な開発経緯、デバッグの記録、および過去の決定事項を時系列で記録するものです。現在のプロジェクトの全体像や次のアクションプランについては、`@planning/SYSTEM_DESIGN_MEMO.md` を参照してください。

---
## 開発進捗（2025-12-22）

### 注文リクエスト生成ロジックの実装
- **意思決定ロジックの実装**: `agent.go` の `tick` メソッド内に、シグナルファイルから読み込んだ売買指示（BUY/SELL）に基づき、注文リクエストを生成するロジックを実装しました。
- **注文内容の決定**:
    - **BUYシグナル**: `agent_config.yaml` で設定された `lot_size` に基づいて注文数量を決定します。重複買いを避けるため、すでにポジションを保有している銘柄の買いシグナルは無視します。
    - **SELLシグナル**: 保有しているポジションの全数量を売却するリクエストを生成します。ポジションがない銘柄の売りシグナルは無視します。
- **注文実行**: 生成されたリクエストを `tradeService.PlaceOrder` メソッドに渡し、注文を発行します。成功・失敗の結果はログに出力されます。
- **ビルド確認**: 上記の変更後、`go build ./...` を実行し、コンパイルエラーや依存関係の問題がないことを確認しました。

---
## 開発進捗（2025-12-21）

### `TEST-005` の保留とエージェント開発への移行
本番環境におけるセッションの無通信タイムアウト仕様を解明するテスト(`TEST-005`)は、原因不明の失敗が続いたため一旦保留としました。
これに伴い、次の開発フェーズである **エージェント(Agent)の要件定義と実装** に着手しました。

### エージェントの要件定義と骨格実装
- **要件定義**: エージェントが利用するツール群（モデルメーカー、シグナルメーカー等）とその役割を`planning/AGENT_REQUIREMENTS.md`に定義し、アーキテクチャの共通認識を確立しました。
- **設定/シグナル読込実装**: `agent_config.yaml`とバイナリ形式のシグナルファイルを読み込む機能を実装し、テストを完了しました。
- **実行ループ実装と動作確認**: エージェントのメインループを実装し、`main.go`から起動・安全に停止する仕組みを構築。ダミーデータを用いて、定期的な実行とファイル読み込みが成功することを確認しました。
- **状態管理機能とサービス連携**: エージェントが保有ポジションや残高等の内部状態をスレッドセーフに管理する機能を実装しました。また、外部APIクライアントとの連携を抽象化する`TradeService`インターフェースを導入し、エージェント起動時に実際の口座情報を取得して内部状態を同期する機能の実装と動作確認を完了しました。

---
## 開発進捗（2025-12-17）

### Goaサービスの追加開発 - master.update の進捗
本日は、`master` サービスの `update` メソッドの追加開発を進めました。
1.  **ユースケーステストの確認と修正**: `GetStock` メソッドのテスト（`TestGetStock_Success`, `TestGetStock_NotFound`, `TestGetStock_RepoError`）が、実装の変更（ローカルDBからの取得）に合わせて更新されていなかった問題を修正しました。`DownloadAndStoreMasterData` のユースケースとテストは既に存在し、正常に機能することを確認しました。
2.  **DI設定の確認**: `cmd/myapp/main.go` における `master` サービスの依存性注入（DI）設定は正しく、変更は不要であることを確認しました。
3.  **ハンドラの実装と修正**: `internal/handler/web/master_service.go` に `Update` メソッドを追加しましたが、Goaが `Payload(Empty)` で生成するインターフェースとの不一致によりコンパイルエラーが発生しました。これを修正し、正しいメソッドシグネチャに調整しました。
4.  **Goaコードの再生成**: 上記の修正を反映させるため `goa gen` を実行し、Goa生成コードを最新の状態に更新しました。

### APIログイン問題の発生
上記作業完了後、統合テストのためにアプリケーションサーバーを起動しようとした際に、APIログインエラーが発生しました。
*   **エラー内容**: `result code 10033: 電話番号認証が認証されない、ユーザID、暗証番号のご入力間違いが弊社規程回数を超えたため、現在ログイン停止中です。(ログイン停止の解除は、コールセンターまでお電話下さい。)`
*   **影響**: Tachibana APIアカウントがロックされているため、現在、本番環境およびデモ環境へのログインができません。これにより、`POST /master/update` エンドポイントの統合テストを含め、APIとの連携が必要な機能のテストをこれ以上進めることができません。
*   **今後の対応**: ユーザー様より、明日以降に証券会社のコールセンターに連絡してログイン停止の解除を依頼する予定であるとのご指示がありました。アカウントロックが解除され次第、統合テストを再開します。

---
### 開発進捗 (2025-12-16)

#### API統合テストの実施とそれに伴うデバッグ
`SYSTEM_DESIGN_MEMO.md`のアクションプランに基づき、実装済みAPIエンドポイントの統合テストを実施した。この過程で複数のバグが顕在化し、それらを段階的に修正した。

1.  **APIのデモ環境と本番環境の仕様調査**:
    *   **本番環境の2FA**: `price_info_client`のテストを本番環境に対して実行した結果、`result code 10088`エラーによりログインに失敗。これがAPIの仕様である「電話番号認証」によるものであることを特定した。認証用の電話番号等の詳細情報を`SYSTEM_DESIGN_MEMO.md`に記録した。
    *   **手動テスト用のヘルパー実装**: 本番環境でのテストを可能にするため、手動で取得したセッション情報（Cookie等）をクライアントに設定するヘルパー関数 `SetLoginStateForTest` と、それを使用するテストファイル `price_info_client_impl_prod_test.go` を作成した。

2.  **`/master/stocks/{symbol}` エンドポイントのデバッグ**:
    *   **バグ1（データソースの誤り）**: `GET /master/stocks/7203` を実行したところ、`Stock master not found` エラーが発生。調査の結果、`GetStock`ユースケースがローカルDBではなく、毎回外部APIを呼び出しているという根本的なバグを発見し、ローカルDB (`masterRepo`) を参照するように修正した。
    *   **バグ2（GORMリレーションエラー）**: 上記修正後、`invalid field found for struct ... TickRules` というGORMのエラーに遭遇。`StockMaster`モデルと`TickRule`モデルの関連付け定義が不適切であることが原因であったため、`StockMaster`の`TickRules`フィールドに`gorm:"-"`タグを追加し、DB読み込み時にこのリレーションを一旦無視することで問題を解決した。
    *   **問題3（文字化け）**: APIからの応答で日本語が文字化けする問題が発生。当初はAPIのエンコーディングが`EUC-JP`であると仮説を立て修正したが解決せず。最終的に、`curl -o`でファイルに出力し、それをテキストエディタで確認することで、サーバーからの応答（UTF-8）は正しく、ターミナルの表示に問題があったことを切り分けた。

3.  **開発効率の改善**:
    *   開発中に毎回マスターデータを同期する非効率性を解消するため、`cmd/myapp/main.go`に`-skip-sync`コマンドラインフラグを追加。これにより、開発時のサーバー起動時間を大幅に短縮した。

4.  **統合テストの完了**:
    *   上記デバッグを経て、`/master/stocks/{symbol}`、`/balance`、`/positions`、`/order` の4つのエンドポイント全てで、`200 OK` または `201 Created` が返却され、期待通りのJSONデータ（または`order_id`）が得られることを確認。「API統合テストの拡充」タスクを完了した。

---

## 6. 調査と決定事項 (2025-12-02)

### Issue 1: リアルタイムイベントの受信方式の特定
-   **調査**: 立花証券APIのドキュメントおよび公式GitHubリポジトリ(`e-shiten-jp/e_api_websocket_receive_tel.py`)のサンプルコードを調査・解析した。
-   **結論**: APIはリアルタイム配信用に **WebSocket (`EVENT I/F`)** を提供している。リアルタイム性と効率性を考慮し、本システムではこのWebSocket方式を採用する。

#### WebSocket (`EVENT I/F`) の仕様
-   **接続URL**:
    1.  通常のログインAPI(`auth_client`)を呼び出し、認証を行う。
    2.  レスポンスに含まれるWebSocket専用の **仮想URL (`sUrlEventWebSocket`)** を取得する。
    3.  この仮想URLに対し、購読したい銘柄コードや情報種別(`p_evt_cmd=FD`等)をクエリパラメータとして付加し、接続する。
-   **データ形式**:
    -   一般的なJSONではなく、**特殊な制御文字で区切られた独自のテキスト形式**である。
    -   `\x01` (`^A`): 項目全体の区切り
    -   `\x02` (`^B`): 項目名と値の区切り
    -   `\x03` (`^C`): 項目内で複数の値を区切る
    -   この仕様に基づき、Go側で専用のパーサーを実装する必要がある。

### Issue 2: Go-Python間の連携インターフェース設計
-   **方針**: 上記WebSocketの採用に伴い、シグナル系統の連携方式を具体化する。
-   **Go -> Python**: GoのWebSocketクライアントがリアルタイムデータを受信する都度、パース処理を行い、案2の通りPython側のWeb APIエンドポイント (例: `POST /api/signal`) へHTTP POSTでプッシュ通知する。
-   **Python -> Go**: 従来の方針通り、Go側で注文受付用のHTTP API (例: `POST /api/order`) を用意する。

### 次のアクション: GoによるWebSocketクライアントの実装
上記方針に基づき、Go側で `EVENT I/F` をハンドリングするクライアントの実装に着手する。

1.  **ファイル作成**:
    *   `internal/infrastructure/client/event_client.go` (インターフェース)
    *   `internal/infrastructure/client/event_client_impl.go` (実装)
2.  **接続処理の実装**: ログイン機能と連携し、取得した仮想URLを使ってWebSocketサーバーに接続する処理を実装する。
3.  **パーサーの実装**: 受信した独自形式のメッセージを制御文字で分割・解析し、Goのデータ構造（`map`や`struct`）に変換するパーサーを実装する。
4.  **イベントループの実装**: サーバーから継続的にメッセージを受信し、パーサーを通して処理するイベントループを実装する。
5.  **アプリケーションへの統合**: 実装したクライアントをアプリケーション全体に組み込み、受信データを後続処理（Pythonへの通知など）へ連携させる。

### 開発進捗 (2025-12-02)

#### Issue 1: リアルタイムイベントの受信方式の特定 (進捗)
-   GoによるWebSocketクライアント (`EventClient`) の実装に着手し、`event_client.go` および `event_client_impl.go` を作成した。
-   WebSocketメッセージの独自形式を解析するパーサー (`ParseMessage`) の単体テストは**PASS**した。
-   デモAPIへのWebSocket接続テスト (`TestEventClient_ConnectReadMessagesWithDemoAPI`) を実装したが、依然として `websocket: bad handshake` エラーで**FAIL**している。
-   これまでに `Origin` ヘッダーと `User-Agent` ヘッダーの追加を試みたが、エラーは解消されていない。

#### 次のステップ (2025-12-03 以降)
-   引き続き `websocket: bad handshake` エラーの原因を詳細に調査する。APIドキュメントの再確認、Pythonサンプルコードのより深い分析、または`gorilla/websocket`とAPIサーバー間の通信プロトコルの詳細な比較が必要となる可能性がある。

### 開発進捗 (2025-12-03)

#### `websocket: bad handshake` エラーの深掘り調査

-   **問題**: `Subprotocol`ヘッダーを追加後も、依然として `websocket: bad handshake` エラーが解消しない。
-   **仮説1: 認証Cookieの欠落**:
    -   **調査**: 公式PythonサンプルおよびGoの参考実装(`tsuchinaga/go-tachibana-e-api`)を再度調査。ログイン時に取得した認証情報(`Cookie`)が、後続のWebSocketハンドシェイクリクエストに含まれていないことが原因である可能性が高いと判断。
    -   **修正**: `TachibanaClientImpl`が`CookieJar`を持つ共有の`http.Client`インスタンスを一元管理するよう、大規模なリファクタリングを実施。
        1.  `tachibana_client.go`: `TachibanaClientImpl`に`httpClient *http.Client`フィールドを追加し、`NewTachibanaClient`で`CookieJar`と共に初期化するよう修正。
        2.  `util.go`: `SendRequest`, `SendPostRequest`が、引数で渡された共有`http.Client`インスタンスを使用するよう修正。
        3.  `auth_client_impl.go`, `balance_client_impl.go`, `master_data_client_impl.go`, `order_client_impl.go`, `price_info_client_impl.go`: `SendRequest`等の呼び出し時に、共有`httpClient`を渡すよう全ファイルを修正。
        4.  `event_client_impl.go`: WebSocket接続時に`CookieJar`を`websocket.Dialer`に設定するよう修正。
-   **仮説2: `Origin`ヘッダーの形式不備**:
    -   **調査**: 上記修正後もエラーが解消せず。公式Pythonサンプルの`Origin`ヘッダーがパス情報を含まない (`https://<hostname>`) のに対し、こちらの実装ではパス情報まで含めてしまっている (`https://<hostname>/<path>`) ことを発見。これが原因である可能性を特定。
    -   **修正**: `event_client_impl.go`を修正し、`Origin`ヘッダーが`scheme`と`host`のみで構成されるよう修正。
-   **結果**: 上記2つの仮説に基づき大規模な修正を行ったが、テスト結果は変わらず `websocket: bad handshake` エラーが継続。

#### 新たな可能性と今後のアクション

-   **新たな可能性（API稼働時間）**: ユーザーからの指摘により、エラーの根本原因が技術的な問題ではなく、**APIの稼働時間（取引時間外）**である可能性が浮上した。リアルタイムAPIは、市場が閉まっている時間帯には接続を拒否する仕様であることが多い。
-   **次のアクションプラン**:
    1.  **最優先事項**: 平日の取引時間中（例: 9:00〜15:00 JST）に、現在のコードのまま再度テスト(`TestEventClient_ConnectReadMessagesWithDemoAPI`)を実行し、接続が成功するかどうかを確認する。
    2.  **次善手（取引時間中でも失敗した場合）**: もし取引時間中でも`bad handshake`エラーが解消されない場合は、原因の切り分けのため、「Cookieが本当に必要か」を再検証する。具体的には、`eventClient.Connect`に`nil`の`CookieJar`を渡してテストを実行し、挙動の変化を確認する。

### 開発進捗 (2025-12-06)

#### アーキテクチャの再定義とGoa導入
- **アーキテクチャの再定義**: 議論を経て、システム全体の設計を「エージェント中心モデル」に更新。Go APIラッパー、Pythonシグナル生成サービス、そして全体の司令塔となるエージェントの3層構造を定義した。短期計画としてエージェントをGoで実装し、長期目標としてRustへの移行を目指す方針を固めた。
- **ドキュメント更新**: 上記の新アーキテクチャに合わせて、`SYSTEM_DESIGN_MEMO.md`および`README.md`を全面的に更新した。
- **ディレクトリ構造の変更**: 新しいアーキテクチャの責務を明確にするため、`internal/interface`ディレクトリを廃止し、`internal/handler`（Webリクエスト処理層）と`internal/agent`（エージェントロジック層）に再編成した。
- **Goaフレームワーク導入**:
    1. Goaツールをインストール。
    2. APIの設計図として`design/design.go`を作成。
    3. `goa gen`コマンドでコードを自動生成。
    4. サービス実装の雛形として`internal/handler/web/order_service.go`を作成。
    5. アプリケーションのエントリーポイントとして`cmd/myapp/main.go`を作成。
- **サーバー起動とAPIテスト**:
    - 複数回にわたるコンパイルエラーのデバッグ（`import`パス、`WaitGroup`の使用法、Goaの`Logger`インターフェース、`Muxer`の`Handle`メソッドなど）を経て、`go run ./cmd/myapp/main.go`による**サーバー起動に成功**した。
    - `Invoke-WebRequest`コマンドを使用し、`POST /order`エンドポイントのテストを実施。HTTPステータス`201 Created`とダミーの注文ID `{"order_id":"order-12345"}`が返却されることを確認し、**APIが正常に動作していることを確認した**。

#### 次回のアクションプラン (2025-12-07 以降)
1.  **つなぎこみ実装**: `order_service.go`のダミー処理を、実際の`OrderUsecase`を呼び出すロジックに置き換える。
2.  **WebSocket接続テスト**: `websocket: bad handshake`エラーのデバッグを、平日の取引時間中に実施する。

### 開発進捗 (2025-12-07)

#### `Order` サービスのバックエンド実装
TDD（テスト駆動開発）に基づく標準手順に沿って、`POST /order` APIのバックエンド実装を推進した。

1.  **開発標準手順の策定:**
    *   TDDに基づいたGoaサービス実装の標準手順を新たに策定し、本ドキュメントに追記した。今後、他のGoaサービスを実装する際もこの手順に統一する。

2.  **ユースケースの実装と単体テスト:**
    *   `OrderUseCase` の振る舞いを定義する単体テスト (`order_usecase_impl_test.go`) を先行して作成した。
    *   コンパイルエラーとテスト失敗を段階的に修正し、テストをすべてパスする `OrderUseCase` の実装 (`order_usecase_impl.go`) を完了させた (`go test ./internal/app/...` は `PASS`)。

3.  **依存性注入 (DI) とハンドラのつなぎこみ:**
    *   `cmd/myapp/main.go` を修正し、`OrderClient` → `OrderUseCase` → `OrderService` (ハンドラ) の依存関係を正しく注入した。
    *   `internal/handler/web/order_service.go` を修正し、APIリクエストを `OrderUseCase` に連携するようにした。

4.  **統合テストと課題の特定:**
    *   サーバーを起動し、`POST /order` API の統合テストを実施。
    *   結果、`TachibanaClient` が未ログイン状態だったため、「`not logged in`」エラーが発生することを確認。アプリケーションのライフサイクルにおけるログイン状態管理の必要性が明らかになった。

#### 次回のアクションプラン (2025-12-08 以降)

1.  **最優先: 起動時ログイン処理の実装**
    *   **対象ファイル:** `cmd/myapp/main.go`
    *   **内容:** `TachibanaClient` の初期化後、サーバーがリクエストの受付を開始する前に `tachibanaClient.Login()` を呼び出す処理を追加する。ログインに失敗した場合は、エラーをログに出力してアプリケーションを終了させる。
    *   **目的:** 「`not logged in`」エラーを解消し、統合テストを成功させる。

2.  **統合テストの再実行**
    *   上記修正後、再度 `go run ./cmd/myapp/main.go` でサーバーを起動し、`POST /order` API を呼び出して、HTTPステータス `201` が返ってくることを確認する。

3.  **新規タスクの起票: `TachibanaClient` のセッション自動管理機能の実装**
    *   アプリケーションの長期的な安定稼働のため、より堅牢なセッション管理メカニズムを実装する必要がある。
    *   **具体的な検討事項:**
        *   セッションの有効期限が切れる前の定期的な再ログイン処理。
        *   API呼び出し時に認証エラーが返された場合の、動的な再ログインとリクエストのリトライ処理。
    *   このタスクは、本件の完了後、新たなIssueとして計画・管理する。

### 開発標準手順

### リファクタリング標準手順 (2025-12-09 追記)
レイヤー間の責務移動など、アーキテクチャの健全性を維持するためのリファクタリングは、以下の手順書に従って実施する。
- **`planning/REFACTORING_PROCEDURE.md`**

### 共通操作におけるユーザーとの連携方針

プロジェクトのビルド、アプリケーションサーバーの起動、および `curl` コマンドなどによるAPIエンドポイントのテストといった、システムの状態を変更したり、外部との連携を伴う共通操作については、以下の原則に基づきユーザーに実行を依頼する。

-   **ビルド操作**: `go build` 等のビルドコマンド。
-   **アプリケーションサーバーの起動**: `go run` やコンパイル済みバイナリの実行など。
-   **外部連携コマンド**: `curl`, `Invoke-WebRequest` など、APIエンドポイントへのリクエスト送信。

これは、ユーザー環境への影響を最小限に抑え、各ステップにおいてユーザーの明示的な承認を得るためのものである。

### Goaサービス実装の標準手順 (2025-12-07 追記)

GoaでAPIサービスを実装する際の標準的な手順を以下に定める。これは、テスト駆動開発(TDD)のアプローチを取り入れ、堅牢なシステム構築を目指すものである。すべてのGoaサービス実装において、この手順に統一して開発を進めること。

**ゴール**: 特定のAPIエンドポイントが、クライアントから受け取った情報に基づき、インフラ層のクライアントを呼び出し、必要な処理を実行して、その結果を返す。

**前提**: Goaの設計ファイル (`design/design.go`) にAPI定義が完了しており、`goa gen` によってコードが自動生成されていること。また、インフラ層の外部APIクライアントは単体テストが完了していること。

#### ステップ1: ユースケースの「振る舞い」をテストで定義する (TDD)
目的: `UseCase` が持つべき振る舞いをテストで定義する。

1.  **テストファイル作成**: `internal/app/<service_name>_usecase_impl_test.go` を新規作成。
2.  **テスト内容**:
    *   モック（例: `OrderRepository`, `TachibanaOrderClient`）を準備し、`UseCase` が依存するコンポーネントが期待通りに呼び出されることを検証する。
    *   成功ケース、失敗ケース、バリデーションエラーなど、主要なシナリオに対するテストケースを記述する。
3.  **実行**: `go test ./internal/app/...` を実行。テストはコンパイルエラーまたは失敗するはず。これが、次の実装の明確なゴールとなる。

#### ステップ2: テストをパスさせるユースケースを実装する
目的: ステップ1で書いたテストをパスさせる。

1.  **実装ファイル作成**: `internal/app/<service_name>_usecase_impl.go` を新規作成。
2.  **実装内容**:
    *   `<Service>UseCaseImpl` 構造体を定義し、依存する `Repository` や `Client` をフィールドに持つ。
    *   `Execute<Service>` メソッド（または対応するメソッド）を実装する。この中で、インフラ層のクライアントを呼び出し、ビジネスロジックを実行する。
3.  **実行**: `go test ./internal/app/...` を実行し、**ステップ1のテストがすべてパスする**まで実装を修正する。

#### ステップ3: アプリケーション起動時の依存性注入 (DI)
目的: アプリケーション起動時に、各コンポーネントを正しく組み立てる。

1.  **ファイル修正**: `cmd/myapp/main.go` を修正。
2.  **実装内容**:
    *   インフラ層のクライアント、リポジトリのインスタンスを作成。
    *   上記を `New<Service>UseCaseImpl` に渡して `UseCase` のインスタンスを作成。
    *   作成した `UseCase` を `web.New<Service>Service` に渡して `Service` (ハンドラ) のインスタンスを作成。
    *   Goaサーバーに `Service` を登録する。
3.  **実行**: `go run ./cmd/myapp/main.go` を実行し、コンパイルエラーや起動時エラーが出ないことを確認する。

#### ステップ4: ハンドラとユースケースのつなぎこみ
目的: APIハンドラから、DIされたユースケースを呼び出す。

1.  **ファイル修正**: `internal/handler/web/<service_name>_service.go` を修正。
2.  **実装内容**:
    *   Goaの `Payload` を `app.Params` に変換する。
    *   `s.usecase.Execute<Service>(...)` を呼び出す。
    *   結果をGoaの `Result` に変換して返す。

#### ステップ5: 統合テスト
目的: APIエンドポイントを実際に呼び出し、システム全体が正しく連携して動作することを確認する。

1.  **実行**:
    1.  `go run ./cmd/myapp/main.go` でサーバーを起動。
    2.  `curl` などのツールで対象のAPIエンドポイントを呼び出す。
2.  **確認**:
    *   期待通りのHTTPステータスコードとレスポンスボディが返ってくること。
    *   (必要に応じて) データベースや外部システムのログなどで、処理が正しく行われたことを確認する。



Invoke-WebRequest -Uri http://localhost:8080/order -Method POST  -Headers @{"Content-Type"="application/json"} -Body '{"symbol": "7203", "trade_type": "BUY", "order_type": "MARKET", "quantity": 100}'

### 開発進捗 (2025-12-08)

#### `POST /order` APIの統合テスト成功とデバッグの軌跡
`not logged in`エラーの解消から始まり、`order failed with result code : `という500エラーの解決まで、段階的なデバッグを経て`POST /order` APIの統合テストを成功させた。

1.  **起動時ログイン処理の実装:**
    *   `main.go`に`tachibanaClient.Login()`を呼び出す処理を追加し、「`not logged in`」エラーを解消。
    *   `config.Config`のフィールド名（`UserID` -> `TachibanaUserID`）の不整合を修正し、ログイン処理を正常に完了させた。

2.  **注文API (500エラー) のデバッグ:**
    *   `order_client_impl_neworder_test.go`が成功することから、APIサーバー経由のリクエストとテストのリクエスト内容の差異を調査。
    *   **原因1 (SecondPasswordの欠落):** `order_usecase_impl.go`で第二パスワードが設定されていなかったため、テストコードに倣いログインパスワードを渡すように修正。しかし、エラーは解消しなかった。
    *   **原因2 (必須フィールドの不足):** さらに比較した結果、逆指値関連の複数のフィールド（`GyakusasiOrderType`など）がリクエストに不足していることが根本原因であると特定。`order_usecase_impl.go`でこれらのフィールドにテストコードと同じデフォルト値を設定したところ、**API呼び出しが成功し、HTTPステータス `201` と注文IDが返却されることを確認した。**

#### 次回のアクションプラン (2025-12-09 以降)

`POST /order` APIの基本的なE2Eテストが成功したため、次の実装フェーズに進む。

1.  **リポジトリ層の実装:**
    *   **対象:** `internal/infrastructure/repository`
    *   **内容:** 現在ダミー実装となっている`OrderRepository`を、PostgreSQLと`gorm`を使用して永続化する実装に置き換える。データベース接続のセットアップも含む。
    *   **目的:** 注文情報をデータベースに保存し、システムの永続化基盤を構築する。

2.  **WebSocket接続テストの再開:**
    *   **対象:** `TestEventClient_ConnectReadMessagesWithDemoAPI`
    *   **内容:** `websocket: bad handshake`エラーのデバッグを、平日の取引時間中に再度実施する。
    *   **目的:** リアルタイムの株価・約定情報を受信する機能を確立する。

3.  **Goaサービスの追加開発:**
    *   **対象:** `Balance`, `Position`などの未実装サービス
    *   **内容:** `SYSTEM_DESIGN_MEMO.md`に記載した「Goaサービス実装の標準手順」に従い、他のサービスの開発に着手する。
    *   **目的:** アプリケーションのコア機能を拡充する。


まずユースケースがセカンドパスワードを扱う必要があるのか？行き当たりばったりで浅はかな解決に走り、アーキテクチャの基本から外れていないか

非常に鋭いご指摘です。そして、その疑問はアーキテクチャを考える上で極めて重要です。

  結論から言うと、理想的な設計では、ユースケースがパスワードのようなインフラ層の詳細を直接扱うべきでは
  ありません。

  ご指摘の通り、UseCaseはアプリケーションのビジネスルールを記述する層であり、「パスワードを使ってAPIリ
  クエストを組み立てる」という処理は、本来インフラ層が担当すべき責務です。

  では、なぜ今回はパスワードを渡す形になったのか

  これは、APIを早期に動作させることを優先し、変更範囲を最小限に抑えた結果、アーキテクチャ的な理想から
  は少し妥協した形になっているためです。

  現状の設計では、UseCaseが「証券会社APIのリクエストDTO（ReqNewOrder）を組み立てる」という責務を持って
  います。そして、そのリクエストの一部としてSecondPasswordが必須であるため、UseCaseはどうしてもパスワ
  ードを知る必要がありました。

  より良い設計（今後のリファクタリング案）

  本来あるべき姿は、OrderClientインターフェースの抽象度を上げ、パスワードを完全にインフラ層に隠蔽する
  ことです。

  例えば、OrderClientのNewOrderメソッドが、UseCaseからOrderParamsのようなビジネス的な情報だけを受け取
  ります。
   1 // in OrderClient interface
   2 NewOrder(ctx context.Context, params app.OrderParams) (*response.ResNewOrder, error)
  そして、OrderClientの実装（order_client_impl.go）の中で、paramsと、設定ファイルから取得したSecondPas
  swordを組み合わせて、実際のAPIリクエスト（ReqNewOrder）を組み立てます。

  こうすれば、UseCaseはパスワードの存在を一切知ることなく、自身のビジネスロジックに集中できます。

  結論

  あなたの「ユースケースがパスワードを扱うべきか」という問いは正しく、現在の実装はその点で改善の余地が
  あります。
  今回の修正は、まずエラーを解消するための現実的な一歩と捉えていただければ幸いです。将来的には、この部
  分をリファクタリングして、よりクリーンな関心の分離を目指すべきだと考えています。

---

### 開発進捗 (2025-12-09)

#### アーキテクチャ改善計画：責務の分離リファクタリング

- **課題の特定**: `POST /order` APIのデバッグ過程で、`OrderUseCase` がインフラ層の詳細である `SecondPassword` を扱っている問題が明らかになった。これは「関心の分離」の原則に反しており、技術的負債となる。
- **標準手順の策定**: このようなレイヤー間の責務移動を伴うリファクタリングを安全かつ一貫して行うため、新たに `planning/REFACTORING_PROCEDURE.md` を作成した。

### 開発進捗 (2025-12-11)

#### `OrderClient` 関連テストの修正
- **課題**: `SecondPassword` の責務を `UseCase` 層から `Infrastructure` 層へ移譲するリファクタリング (`2025-12-09` 実施) の影響で、`OrderClient` を利用している複数のテスト (`cancelorder`, `cancelorderall`, `correctorder`) でコンパイルエラーが発生していた。
- **修正**:
    1. `NewOrder` メソッドの呼び出し部分を、新しい `client.NewOrderParams` 構造体を使うように修正し、すべてのコンパイルエラーを解消。
    2. `order_client_impl_cancelorder_test.go` で発生していた実行時エラー（APIエラーコード `13001`, `11121`）を調査。原因が逆指値注文のパラメータにあると特定。
    3. ユーザーの指示に基づき、テストの意図（特定のリクエストを生成すること）を維持するため、`NewOrderParams` の値は元のテストコードの値を保持するように最終調整。これにより、テストはコンパイル可能だが、APIの仕様により実行時には失敗する可能性がある状態となった。
- **結論**: `OrderClient` に関連するテストは、リファクタリング後のインターフェースに準拠した形に修正され、コンパイル可能な状態に復旧した。

#### `OrderClient` メソッドの SecondPassword 責務移譲の完了
- **課題**: `NewOrder` メソッドに適用した `SecondPassword` の責務移譲が、`OrderClient` の他のメソッド (`CorrectOrder`, `CancelOrder`, `CancelOrderAll`) に対して未完了であった。
- **修正**:
    1. `internal/infrastructure/client/order_client.go` 内の `OrderClient` インターフェースを更新し、`CorrectOrderParams`, `CancelOrderParams`, `CancelOrderAllParams` の各構造体を定義し、対応するメソッドのシグネチャを変更。
    2. `internal/infrastructure/client/order_client_impl.go` 内で、変更されたインターフェースに合わせて `CorrectOrder`, `CancelOrder`, `CancelOrderAll` の実装を修正し、`SecondPassword` の扱いを内部にカプセル化。
    3. 関連するテストファイル (`order_client_impl_cancelorder_test.go`, `order_client_impl_correctorder_test.go`, `order_client_impl_cancelorderall_test.go`) を更新されたインターフェースに合わせるように修正。
    4. `internal/app/order_usecase_impl.go` はこれらのメソッドを使用していないため、変更は不要であることを確認。
- **結論**: `OrderClient` のすべての関連メソッドにおいて `SecondPassword` の管理責務が `Infrastructure` 層に完全に移譲され、リファクタリングが完了した。

#### リポジトリ層の静的コードレビュー完了
- **課題**: リポジトリ層の実装が `gorm` を使用して適切に行われているか、静的に確認する必要があった。
- **レビュー結果**:
    1. `OrderRepository` (`domain/repository/order_repository.go`, `domain/model/order.go`, `internal/infrastructure/repository/order_repository_impl.go`) をレビューし、インターフェース、`gorm` タグ付きモデル、`gorm` ベースの実装が適切であることを確認。
    2. `PositionRepository` (`domain/repository/position_repository.go`, `domain/model/position.go`, `internal/infrastructure/repository/position_repository_impl.go`) をレビューし、同様に適切であることを確認。
    3. `SignalRepository` (`domain/repository/signal_repository.go`, `domain/model/signal.go`, `internal/infrastructure/repository/signal_repository_impl.go`) をレビューし、同様に適切であることを確認。
    4. `MasterRepository` (`domain/repository/master_repository.go`, `domain/model/master_*.go`, `internal/infrastructure/repository/master_repository_impl.go`) をレビューし、同様に適切であることを確認。`FindByIssueCode` メソッドは `entityType` に基づいて適切なモデルを検索する汎用的な実装であり、既存コードに修正すべき論理的欠陥はなかった。
- **結論**: すべてのリポジトリコンポーネント（インターフェース、`gorm` タグ付きモデル、`gorm` ベースの実装）は、静的コードの観点から完成していると判断される。

### 開発進捗 (2025-12-12)

#### リポジトリ層の統合と永続化の実現
ダミー実装だったリポジトリ層を、実際のデータベース（PostgreSQL）に接続する実装に置き換え、アプリケーションの永続化基盤を構築した。

1.  **開発用データベース環境の構築:**
    *   `docker-compose.yml` を新規に作成し、PostgreSQLコンテナを定義。開発環境のデータベースをDockerで簡単に起動できるようにした。
    *   `.env` ファイルにデータベース接続情報（`DB_HOST`, `DB_USER`等）を設定する方法を明確化し、接続問題を解決した。

2.  **`main.go` へのGORM統合:**
    *   アプリケーション起動時に、`gorm` を用いてPostgreSQLに接続する処理を `cmd/myapp/main.go` に実装。
    *   `db.AutoMigrate` を使用し、`Order`, `Position`, `Signal`, `StockMaster` などのドメインモデルに基づいて、データベーススキーマが自動的に生成・更新されるようにした。
    *   `OrderUseCase` に注入するリポジトリを、ダミーの `dummyOrderRepo` から `gorm` ベースの `repository_impl.NewOrderRepository` に置き換えた。

3.  **コンパイルエラーと実行時エラーの修正:**
    *   `main.go` で発生していた、モデル名 (`StockMaster` 等) やリポジトリのコンストラクタ名 (`NewOrderRepository`) の不一致によるコンパイルエラーを修正した。
    *   `.env` の設定不備に起因するデータベース接続エラー (`lookup db: no such host`) を特定し、ユーザーが設定を修正することで解決に導いた。

4.  **統合の最終確認:**
    *   上記修正後、`go run ./cmd/myapp/main.go` を実行し、アプリケーションが正常に起動、データベース接続、スキーマのマイグレーション、APIへのログインを完了し、HTTPサーバーがリッスン状態になることを確認した。

#### 次回のアクションプラン (2025-12-13 以降)

1.  **WebSocket接続テストの再開 (最優先)**:
    *   **対象:** `TestEventClient_ConnectReadMessagesWithDemoAPI`
    *   **内容:** 平日の取引時間中に `websocket: bad handshake`エラーのデバッグを再開する。
    *   **目的:** リアルタイムの株価・約定情報を受信する機能を確立する。

2.  **Goaサービスの追加開発:**
    *   **対象:** `Balance`, `Position`などの未実装サービス
    *   **内容:** `SYSTEM_DESIGN_MEMO.md`に記載した「Goaサービス実装の標準手順」に従い、他のサービスの開発に着手する。
    *   **目的:** アプリケーションのコア機能を拡充する。

### 開発進捗 (2025-12-13)

#### データベースマイグレーションの導入
-   **課題**: 既存の `gorm.AutoMigrate` は開発初期には便利だが、本番環境での運用には不向きであった。
-   **解決策**: `golang-migrate/migrate` ツールを導入し、バージョン管理されたSQLファイルによるマイグレーションシステムを構築した。
-   **具体的な変更**:
    1.  `golang-migrate/migrate` CLIツールをインストール。
    2.  プロジェクトルートに `migrations` ディレクトリを作成し、初期スキーマ (`000001_create_initial_tables.up.sql`, `.down.sql`) を生成。
    3.  既存の `domain/model` 定義から、PostgreSQL用の `CREATE TABLE` および `DROP TABLE` SQLを生成し、マイグレーションファイルに記述。
    4.  `cmd/myapp/main.go` から `db.AutoMigrate(...)` の呼び出しを削除。
    5.  `README.md` を更新し、マイグレーションの実行方法に関する説明を追加。
    6.  `migrations/README.md` を作成し、ディレクトリの目的と利用方法を解説。

#### Goaサービス「Order」のテスト修正
-   **課題**: `OrderUseCase` のモック（`OrderClientMock`）が、`SecondPassword` の責務移譲に伴う `client.OrderClient` インターフェースの変更に追従できておらず、コンパイルエラーが発生していた。
-   **解決策**: `internal/app/tests/order_usecase_impl_test.go` 内の `OrderClientMock` のメソッドシグネチャを、新しい `client....Params` 型に合わせて修正。

#### Goaサービス「Balance」の追加
-   **目的**: 口座の残高サマリーを取得する `GET /balance` エンドポイントを実装。
-   **実装詳細**:
    1.  `design/design.go` に `balance` サービスを定義。主要な残高情報（買付可能額、保証金率など）を `BalanceResult` として抽出。
    2.  `goa gen` でコードを生成。
    3.  `internal/app/tests/balance_usecase_impl_test.go` で `BalanceUseCase` のテスト（成功、クライアントエラー、パースエラー）を定義。
    4.  `internal/app/balance_usecase.go` と `internal/app/balance_usecase_impl.go` でユースケースを実装。`client.BalanceClient.GetZanKaiSummary` を呼び出し、API応答文字列を適切な型にパース。
    5.  `internal/handler/web/balance_service.go` でGoaハンドラを実装。
    6.  `cmd/myapp/main.go` に `BalanceUseCase` と `BalanceService` をDIし、エンドポイントをマウント。
    7.  **デバッグと修正**: `balance.BalanceResult` が `balance.StockbotBalance` という名前で生成されていたため、ハンドラコードを修正。
    8.  `curl` コマンドによる統合テストで動作を確認。

#### Goaサービス「Position」の追加
-   **目的**: 現在保有しているポジション（建玉）の一覧を取得する `GET /positions` エンドポイントを実装。
-   **実装詳細**:
    1.  `design/design.go` に `position` サービスを定義。現物と信用のポジションを統合した `PositionResult` および `PositionCollection` を定義。`type` パラメータによるフィルタリングをサポート。
    2.  `goa gen` でコードを生成。
    3.  `internal/app/tests/position_usecase_impl_test.go` で `PositionUseCase` のテスト（`all`, `cash`, `margin` フィルタリング、クライアントエラー）を定義。
    4.  `internal/app/position_usecase.go` と `internal/app/position_usecase_impl.go` でユースケースを実装。`client.BalanceClient.GetGenbutuKabuList` と `GetShinyouTategyokuList` を呼び出し、統一された `Position` 構造体に変換。
    5.  `internal/handler/web/position_service.go` でGoaハンドラを実装。
    6.  `cmd/myapp/main.go` に `PositionUseCase` と `PositionService` をDIし、エンドポイントをマウント。
    7.  **デバッグと修正**: `PositionUseCaseBalanceClientMock` が `client.BalanceClient` インターフェースの全メソッドを実装していなかった点を修正。
    8.  **デバッグと修正**: `position.PositionCollection` が `position.StockbotPositionCollection` という名前で生成されていたため、ハンドラコードを修正。
    9.  `curl` コマンドによる統合テストで動作を確認。

#### Goaサービス「Master」の追加 (途中)
-   **目的**: 個別銘柄のマスタデータ（PER, PBR等の詳細情報）を取得する `GET /master/stocks/{symbol}` エンドポイントを実装。
-   **実装詳細**:
    1.  `design/design.go` に `master` サービスを定義。`get_stock_detail` メソッドと `StockDetailResult`（PER, PBR等の財務指標を含む）を定義。
    2.  `goa gen` でコードを生成。
    3.  `MasterUseCase` とハンドラの開発を進めたが、統合テストで `Stock detail not found` エラーが発生。
    4.  **原因調査**: 立花証券APIのPythonサンプルコードを分析した結果、`GetIssueDetail` はデモ環境で期待通りの詳細データを返却しない可能性が高いと判明。代わりに `GetMasterDataQuery` を使用し、基本的な銘柄情報のみを取得する方針に転換。
    5.  **再設計と実装（現在デバッグ中）**:
        -   `design/design.design.go` を修正し、`get_stock` メソッドと `StockMasterResult`（銘柄コード、名称、市場、業種コード名など基本的な情報のみを含む）を定義。
        -   `goa gen` を再実行。
        -   `internal/app/tests/master_usecase_impl_test.go` を、`GetMasterDataQuery` をモックし、新しい `StockMasterResult` のフィールドをアサートするように全面修正。
        -   `internal/app/master_usecase.go` および `internal/app/master_usecase_impl.go` を修正し、`GetMasterDataQuery` を呼び出してレスポンスから銘柄情報を抽出し、`StockMasterResult` を返すように変更。
        -   `internal/handler/web/master_service.go` を修正し、新しい `get_stock` メソッドと `master.StockbotStockMaster` 型を使用するように変更。
        -   **現在デバッグ中**: `ResStockMaster` 内のフィールド名 (`YusenSizyou` -> `PreferredMarket`, `GyousyuCode` -> `IndustryCode`, `GyousyuName` -> `IndustryName`) の不一致や、Goa生成型名（`StockbotStockMaster`）のミスマッチ、テストのモック引数不一致、構文エラーなど、複数のコンパイル／実行時エラーを修正中。

---

#### 次回のアクションプラン (2025-12-16 以降)

1.  **WebSocket接続テストの再開 (最優先)**:
    *   **対象:** `TestEventClient_ConnectReadMessagesWithDemoAPI`
    *   **内容:** 平日の取引時間中に `websocket: bad handshake`エラーのデバッグを再開する。
    *   **目的:** リアルタイムの株価・約定情報を受信する機能を確立する。

2.  **API統合テストの拡充**:
    *   **対象**: 起動したアプリケーション全体
    *   **内容**: `curl`や`Invoke-WebRequest`などのツールを使用し、実装済みの各APIエンドポイント（`/order`, `/balance`, `/positions`, `/master/stocks/{symbol}`）が、永続化されたDBと連携して正しく動作するかを体系的にテストする。
    *   **目的**: 各サービスのE2Eでの動作を保証する。

3.  **Goaサービスの追加開発:**
    *   **対象:** `design/design.go` に定義されている未実装のサービス
    *   **内容:** `SYSTEM_DESIGN_MEMO.md`に記載した「Goaサービス実装の標準手順」に従い、他のサービスの開発に着手する。
    *   **目的:** アプリケーションのコア機能を拡充する。

### 開発進捗 (2025-12-14)

#### マスターデータ同期機能の実装完了
長期間にわたるデバッグの末、アプリケーション起動時にマスターデータをダウンロードし、データベースに保存する一連の機能が正常に動作することを確認した。

1.  **API接続の課題解決**:
    *   **ログイン404エラー**: 原因は、`.env`ファイルに設定されたAPIのベースURL (`TACHIBANA_BASE_URL`) のバージョンが古かった (`v4r7`) ことであった。最新のバージョン (`v4r8`) に修正したことで解決した。
    *   **マスターデータ取得エラー**: `DownloadMasterData` APIが返す巨大なストリーミングデータ（改行なしの連続したJSON）が、Go標準ライブラリの`bufio.Scanner`のバッファ上限を超えてしまう問題があった。これは、Pythonサンプルを参考に、チャンクで読み込み`}`を区切り文字として手動でJSONをパースするロジックを実装することで解決した。

2.  **データベース永続化の課題解決**:
    *   **GORMとリレーションのUpsert問題**: リレーション (`TickRules`) を持つGORMモデルをそのまま一括Upsertしようとすると `invalid field` エラーが発生した。`.Omit()`や`.Select()`も期待通りに機能しなかったため、最終的にリポジトリ層でリレーションフィールドを持たないDB保存用のDTO (`dbStockMaster`) にデータを詰め替える「DTOパターン」を採用することで、GORMの一括Upsert機能を活かしつつ問題を構造的に解決した。
    *   **マイグレーションの課題**: モデルとDBスキーマの不整合（カラム不足）や、GORMの主キー規約（`id`カラムの自動探索）に起因するエラーが発生。これらは、`golang-migrate/migrate`を使ってスキーマを修正し、モデル定義から不要な`ID`フィールドを削除することで解決した。

3.  **開発効率の改善**:
    *   マイグレーションを簡単かつ確実に行うため、`.env`ファイルを読み込んで`migrate`ライブラリを直接実行するGoプログラム (`cmd/migrator/main.go`) を作成した。これにより、`go run`コマンド一つで誰でもマイグレーションを実行できるようになった。

#### マスターデータ同期機能におけるDB問題の再発

-   **課題の再発**: マスターデータ同期機能の実行時に、`TickRules`テーブルへのデータ挿入で`ERROR: there is no unique or exclusion constraint matching the ON CONFLICT specification (SQLSTATE 42P10)`エラーが再発した。これは、`tick_rules`モデルの`TickUnitNumber`が`PRIMARY KEY`として定義されているにもかかわらず発生している。
-   **環境差異とデータベース状態**:
    -   別の開発環境ではこの`ON CONFLICT`エラーは発生しておらず、現在の環境でのみ再発している。この事実は、両環境間でデータベーススキーマの状態に差異があることを強く示唆している。
    -   これまでのデバッグ過程で、マイグレーションの失敗によるデータベースの「ダーティ」状態の発生や、`ALTER TABLE ADD COLUMN`の重複エラーなど、スキーマの不整合に起因する問題が複数回発生している。
    -   `ON CONFLICT`句が期待通りに機能するためには、PostgreSQLが`UNIQUE`制約または`EXCLUSION`制約を明示的に認識している必要があり、`PRIMARY KEY`のみでは不十分な場合がある。
-   **現在のデータベース状態の評価**: 現在の環境で発生している一連のデータベース関連エラーは、過去のマイグレーション失敗や不完全な適用により、データベースのスキーマがアプリケーションコードやマイグレーションファイルが期待する状態と一致していないことに起因すると考えられる。開発段階にあるとはいえ、このような不整合なデータベースの状態を維持しようとすることは、デバッグを困難にし、さらなる問題を引き起こす可能性が高いため、現状のデータベース設定は価値がないと判断する。
-   **今後の対応方針**:
    -   `migrations/20251212203023_create_initial_tables.up.sql`に`tick_rules`テーブルの`tick_unit_number`カラムに対して明示的な`CREATE UNIQUE INDEX IF NOT EXISTS idx_tick_rules_tick_unit_number ON tick_rules(tick_unit_number);`を追加した。
    -   今後の開発を確実に行うため、問題が再発した場合は、現在のデータベースを完全に破棄し、クリーンな状態からマイグレーションを再適用することを基本的な運用方針とする。

### 開発進捗 (2025-12-15)

#### マイグレーション管理の安定化とDBの正常化
- **課題**: `2025-12-14`に記録されたDB問題（環境差異による`ON CONFLICT`エラーの再発）の根本原因が、開発の進行に伴うマイグレーション管理の複雑化にあると判断。環境ごとの適用状態の差異が、スキーマの不整合を引き起こしていた。
- **解決策**: 開発初期の安定性と再現性を高めるため、マイグレーションファイルを単一の初期スキーマファイルに統合するリファクタリングを実施。
    1.  `..._add_fields_to_stock_masters.up.sql`の内容を、`..._create_initial_tables.up.sql`の`CREATE TABLE`文にマージした。
    2.  不要になった古いマイグレーションファイルを削除した。
    3.  動作検証として、`docker-compose down -v`でデータベースを完全にクリーンアップした後、`go run ./cmd/migrator/main.go`を実行。統合されたマイグレーションが正常に適用されることを確認した。
- **結論**: これにより、どの開発環境でも一度のマイグレーションで最新のスキーマを確実に構築できるようになり、環境差異に起因するデータベース問題が構造的に解決された。アプリケーションも、クリーンなDB上で正常に起動し、マスターデータを同期できることを確認済み。

---

## 実装から得られた知見（APIクライアント編）

本セクションでは、開発過程で遭遇した立花証券APIの特殊な仕様や、それに対する実装上のノウハウを記録する。

### 1. APIの環境差異とURLのバージョン管理

- **課題**: テスト環境とローカル環境で同じコードにも関わらず、ローカルでのみログインAPIが404エラーを返した。
- **原因**: APIのベースURLにバージョン情報（例: `v4r8`）が含まれており、ローカルの`.env`ファイルに設定されたURLのバージョンが古かった (`v4r7`)。
- **ノウハウ**:
    - APIへの接続テストが失敗する場合、コードのロジックだけでなく、`.env`ファイルに設定されたエンドポイントURL (`TACHIBANA_BASE_URL`など) が、テスト対象の環境で有効なものであるかを最初に確認する必要がある。
    - APIのバージョンアップに伴い、URLも変更される可能性があることを常に念頭に置く。

### 2. マスターデータ取得APIの特殊なストリーミング仕様

- **課題**: 全件マスターデータを取得する`DownloadMasterData` APIを呼び出すと、`bufio.Scanner: token too long`エラーが発生し、ストリームを最後まで読み取れなかった。
- **原因**: このAPIは、数万行に及ぶデータを、**改行なしの単一の巨大なライン、あるいは連続したJSONオブジェクト**としてストリーミング配信する特殊な仕様となっている。Go標準ライブラリの`bufio.Scanner`は改行をデリミタとしており、この形式に対応できない。
- **解決策**: 公式のPythonサンプルコードのロジックを参考に、以下の手動パーシング処理を実装した。
    1. レスポンスボディを固定長のチャンク（例: 4096バイト）で読み込む。
    2. 読み込んだバイト列を一時的なバッファ (`bytes.Buffer`) に蓄積する。
    3. バッファ内にJSONオブジェクトの終端文字 (`}`) が存在するかを検索する。
    4. 終端文字が見つかった場合、そこまでを一つのJSONオブジェクト候補として切り出し、`json.Unmarshal`でデコードを試みる。
    5. デコードが成功した場合、バッファからその部分を削除し、次のオブジェクトの処理に移る。
    6. これを、APIから`CLMEventDownloadComplete`という完了通知オブジェクトが送られてくるまで繰り返す。
- **ノウハウ**:
    - ストリーミングAPIを扱う際は、データがどのような形式・区切り文字で送られてくるかを正確に把握することが極めて重要である。
    - 標準ライブラリで対応できない特殊な形式の場合、公式のサンプルコード（もしあれば）の挙動を模倣した、より低レベルなバイト/チャンク処理の実装が必要となる。