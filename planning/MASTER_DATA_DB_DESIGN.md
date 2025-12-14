# マスターデータにおけるデータベース設計

## 1. 目的
本ドキュメントは、APIからダウンロードしたマスターデータを、アプリケーションで効率的に利用するためのデータベース設計と、データ更新フローを定義することを目的とする。

## 2. 背景・課題
- `DownloadMasterData` APIは、利用可能な全銘柄（約4000）のマスターデータを一括で返す。ダウンロード時に銘柄を絞り込むことはできない。
- 一方、アプリケーションが実際に監視・取引対象とする銘柄は、そのうちの一部（例: 200銘柄）である。
- 毎回全銘柄のデータをDBに保存・更新するのは非効率であり、ストレージと処理時間の無駄につながる。

## 3. 設計要件
- 監視対象となる銘柄のマスターデータのみをデータベースに永続化する。
- 監視対象の銘柄リストは、柔軟に変更可能であること。
- 日々のマスターデータ更新を、差分更新によって効率的に行うこと。

## 4. 設計上の論点（検討事項）

### 4.1. 監視対象銘柄リストの管理方法
アプリケーションが「どの銘柄を監視すべきか」を知るための方法を定義する必要がある。

- **案1: 設定ファイルで管理**
  - `config.yaml` や `.env` のようなファイルに、監視対象の銘柄コードのリストを記載する。
  - **メリット**: 実装が容易。起動時にリストを読み込むだけ。
  - **デメリット**: 監視対象の変更にアプリケーションの再起動が必要になる場合がある。
- **案2: データベースの専用テーブルで管理**
  - `watched_stocks` のようなテーブルを作成し、そこに銘柄コードを保存する。
  - **メリット**: アプリケーションを停止せずに、API経由や直接DBを編集することで動的に監視対象を変更できる。将来的な拡張性が高い。
  - **デメリット**: 初期実装のコストがやや高い。

### 4.2. データ取得・更新フロー
監視対象リストと、APIからダウンロードした全件データをどのように突き合わせ、DBを更新するかのフロー。

1. **(毎日定時実行)** APIから全件マスターデータをダウンロードする。
2. DBまたは設定ファイルから「監視対象銘柄リスト」を取得する。
3. ダウンロードした全件データの中から、「監視対象銘柄リスト」に含まれる銘柄のデータのみを抽出する。
4. 抽出したデータセットを使い、対象のマスターデータテーブル（`stock_masters` など）に対して **UPSERT** 処理を実行する。
   - 既にDBに存在する銘柄は、情報が更新される。
   - DBに存在しない銘柄（監視対象に新たに追加された銘柄）は、新規に挿入される。
5. (オプション) 監視対象から外された銘柄のデータをDBから削除 (`DELETE`)、または無効化 (`is_active = false` のようにフラグを立てる) する処理を検討する。

### 4.3. データベース スキーマ
マスターデータを格納するためのテーブル設計。`DownloadMasterData` で取得できるデータ種別ごとにテーブルを分けるのが基本。

- `stock_masters` (株式銘柄マスタ)
- `stock_market_masters` (株式市場マスタ)
- `tick_rules` (呼値マスタ)
- ...など。

**例: `stock_masters` テーブル**
```sql
CREATE TABLE stock_masters (
    issue_code VARCHAR(10) PRIMARY KEY, -- 銘柄コード (主キー)
    issue_name VARCHAR(255) NOT NULL,
    issue_name_kana VARCHAR(255),
    market VARCHAR(50),
    industry_code VARCHAR(10),
    industry_name VARCHAR(255),
    -- その他、必要なフィールドを追加
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```
- 主キー (`issue_code`) を設定することが、UPSERT処理の前提となる。

**例: `tick_rules` 及び `tick_levels` テーブル (正規化案)**
APIから返される呼値マスタ(`CLMYobine`)は、1つのJSONオブジェクト内に複数の価格帯と呼値が横持ちで含まれる非正規化な形式となっている。
```json
{
  "sCLMID": "CLMYobine",
  "sYobineTaniNumber": "101",
  "sTekiyouDay": "20140101",
  "sKizunPrice_1": "3000.000000",
  "sYobineTanka_1": "1.000000",
  "sKizunPrice_2": "5000.000000",
  "sYobineTanka_2": "5.000000",
  ...
}
```
これをこのままDBに保存すると冗長なため、ドメインモデルでは`TickRule`(親)と`TickLevel`(子)の正規化された親子関係に変換して永続化する。

```sql
CREATE TABLE tick_rules (
    tick_unit_number VARCHAR(10) PRIMARY KEY, -- 呼値の単位番号 (主キー)
    applicable_date VARCHAR(8),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE tick_levels (
    id BIGSERIAL PRIMARY KEY,
    tick_rule_unit_number VARCHAR(10) NOT NULL REFERENCES tick_rules(tick_unit_number), -- 外部キー
    lower_price NUMERIC NOT NULL, -- 基準値段の下限
    upper_price NUMERIC NOT NULL, -- 基準値段の上限
    tick_value NUMERIC NOT NULL,  -- 呼値
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

## 5. 次のステップ

### 完了済み

- [x] **監視対象銘柄リストの管理方法**を決定する。（**CSVファイルで管理**する方針に決定）
  - `watched_stocks.csv` ファイルをプロジェクトルートに配置。
  - `internal/config/config.go` で `watched_stocks.csv` を読み込むロジックを実装済み。
- [x] `MasterRepository` に、UPSERT処理を行うメソッド（例: `UpsertStockMasters(models []*model.StockMaster) error`）を定義し、実装。（`StockMaster` モデルのUPSERTロジックを実装済み）
  - **課題解決**: `gorm` が `model.StockMaster` 内のリレーションフィールド (`TickRules`) を参照し、外部キーがないために発生していた500エラーを解消するため、`UpsertStockMasters` メソッドにおいて、`model.StockMaster` オブジェクトを直接GORMに渡すのではなく、**リレーションフィールドを含まない `map[string]interface{}` の形式に変換してからUPSERT処理を実行**するように修正した。これにより、GORMは`stock_masters`テーブルの主要なデータのみを安全に処理できるようになった。
- [x] `MasterUseCase` に、全体の更新フロー（`DownloadAndStoreMasterData`）を実装。
  - `masterClient.DownloadMasterData` を呼び出し、全件データを取得。
  - `config.WatchedStocks` に基づきデータをフィルタリング。
  - `response.ResStockMaster` と `response.ResStockMarketMaster` から `model.StockMaster` への変換。
  - `masterRepo.UpsertStockMasters` を呼び出し、DBに保存。
- [x] **マスターデータ更新のトリガー**を決定。（**APIエンドポイント `POST /master/update`** で手動トリガーする方針に決定）
  - `design/design.go` に `master` サービス `update` メソッドの定義を追加済み。
  - Goaコード生成とハンドラ実装、統合テストで以下の結果となるため、原因特定と修正が必要
    >curl -i -X POST http://localhost:8080/master/update
     HTTP/1.1 500 Internal Server Error
     Content-Type: application/json
     Date: Sun, 14 Dec 2025 04:09:54 GMT
     Content-Length: 324

    {"name":"fault","id":"bgC_7fLg","message":"failed to upsert stock masters: failed to upsert stock masters: invalid field found for struct stock-bot/domain/model.StockMaster's field TickRules: define a valid foreign key for relations or implement the Valuer/Scanner interface","temporary":false,"timeout":false,"fault":true}


### 残りのタスク

- [ ] **`TickRule` モデルの再設計**:
  - `domain/model/master_tick_rule.go` を修正し、`TickRule`と`TickLevel`の新しいスキーマに合わせた構造に変更する。主キーは`TickUnitNumber`とする。
- [ ] **テーブルスキーマの定義**:
  - `migrations` に `tick_rules` と `tick_levels` テーブルを作成するSQLを追加する。
- [ ] **`MasterRepository` の拡張**:
  - `UpsertTickRulesAndLevels` のような、親子関係をトランザクション内で処理するメソッドを実装する。
- [ ] **`MasterUseCase` の拡張**:
  - `DownloadAndStoreMasterData` 内で、`CLMYobine` データを新しいモデルに変換し、リポジトリのメソッドを呼び出すロジックを追加する。
- [ ] 他のマスターデータ（`StockIssueRegulation`など）についても、同様に永続化の設計と実装を行う。
- [ ] **平日実施タスク**:
  - `Master` サービス (`GetMasterDataQuery`) の再検証。
  - `WebSocket` 接続テストの再開。
