# マスターデータDB設計と管理方針

## 1. はじめに

本ドキュメントは、`SYSTEM_DESIGN_MEMO.md`で提起された「マスターデータ管理の非効率性」という課題に対応するため、DB設計とデータ管理フローを定義することを目的とする。

**課題**: `DownloadMasterData` APIは全市場の全銘柄（数千件）のデータを一括で返す。これを効率的にデータベースに保存・更新する仕組みが必要となる。

**方針**: `DownloadMasterData` APIから返される**全銘柄のマスターデータを、日次バッチ等で効率的にデータベースへ保存・更新する**仕組みを構築する。

## 2. DBスキーマ設計案

### `stock_masters` テーブル
`domain/model/master_stock.go` のモデルに対応する、株式のマスターデータを格納するテーブル。

- **役割**: 全銘柄の最新マスターデータを保持する。
- **カラム**: `domain/model/master_stock.go` の `StockMaster` 構造体フィールドに対応するカラムを定義する。
    - `issue_code` (VARCHAR(255) PRIMARY KEY): 銘柄コードを主キーとする。
    - `stock_name` (VARCHAR(255)): 銘柄名称
    - `issue_name_short` (VARCHAR(255)): 銘柄名略称
    - `issue_name_kana` (VARCHAR(255)): 銘柄名（カナ）
    - `issue_name_english` (VARCHAR(255)): 銘柄名（英語表記）
    - `market_code` (VARCHAR(255)): 優先市場コード
    - `industry_code` (VARCHAR(255)): 業種コード
    - `industry_name` (VARCHAR(255)): 業種コード名
    - `trading_unit` (INTEGER): 売買単位
    - `listed_shares_outstanding` (BIGINT): 上場発行株数
    - `upper_limit` (DOUBLE PRECISION): 値幅上限
    - `lower_limit` (DOUBLE PRECISION): 値幅下限
    - `created_at` (TIMESTAMPTZ): レコード作成日時 (GORM管理)
    - `updated_at` (TIMESTAMPTZ): レコード更新日時 (GORM管理)

*(注: `TickRules`は関連テーブルとして別に管理される。他のマスタデータ（市場マスタ、呼値マスタ等）も同様に、それぞれ専用のテーブルを作成して管理することを想定する)*

## 3. データ更新・管理フロー

### 3.1. `DownloadMasterData` 関数の役割
`DownloadMasterData` 関数は、APIからの全マスターデータを安定して取得する責務を持つ。

1.  **ストリーミング処理**: `DownloadMasterData` は、APIからのレスポンスをストリーミングで処理し、タイムアウトやメモリ枯渇を防ぐ。Goプログラムでは、`bytes.Buffer` を使用してチャンクデータを蓄積し、JSONオブジェクトの終端を示す`}`を検出することで、手動で各JSONオブジェクトをデコードする。
2.  **全件返却**: APIから返された全てのマスターデータ（株式、市場、呼値など）を `*response.ResDownloadMaster` 構造体に格納して呼び出し元に返却する。

### 3.2. データ更新バッチ（`UseCase`）の実装
マスターデータを定期的に更新するためのバッチ処理（または`UseCase`）を実装する。

**処理フロー**:
1.  **全マスターデータ取得**: `MasterDataClient.DownloadMasterData()` を呼び出し、最新の全マスターデータを取得する。
2.  **ドメインモデルへの変換**: `*response.ResDownloadMaster` の中から`ResStockMaster`のスライスを取り出し、`model.StockMaster`のスライスに変換する。この際、APIレスポンスの文字列をGoの適切な型（`int`, `float64`, `int64`など）に`strconv`パッケージを用いて変換する。`ResStockMarketMaster`から取得する`UpperLimit`, `LowerLimit`もここでマッピングする。変換エラーはログ出力に留め、処理は続行する。
3.  **データベース更新**: 変換された `[]*model.StockMaster` を `MasterRepository.UpsertStockMasters()` を通じて`stock_masters`テーブルに一括で更新（Insert or Update）する。

このバッチ処理は、アプリケーション起動時や、1日1回などのスケジュールで実行することを想定する。

## 4. `MasterRepository` に追加が必要なメソッド
上記フローを実現するため、`MasterRepository` インターフェースに以下のメソッドを実装する。

- `UpsertStockMasters(ctx context.Context, masters []*model.StockMaster) error`: 複数の株式マスターデータを一括でUpsertする。（`internal/infrastructure/repository/master_repository_impl.go` にDTOパターンを適用して実装済み）
- `UpsertTickRules(ctx context.Context, tickRules []*model.TickRule) error`: 複数の呼値データを一括でUpsertする。（`TickRules`は`stock_masters`とは別のテーブルに保存）

## 5. 次のアクションプラン (完了済み)
1.  本ファイル (`planning/MASTER_DATA_DB_DESIGN.md`) の内容を最新化。（完了）
2.  `DownloadMasterData` のリファクタリングとテスト成功。（完了）
3.  マスターデータ更新用の `UseCase` (`DownloadAndStoreMasterData`) および、それを呼び出すバッチ処理を実装。（完了）
4.  `MasterRepository` に `UpsertStockMasters` メソッドを実装。（完了）
5.  `cmd/migrator/main.go` を作成し、Goプログラムからマイグレーションを実行する仕組みを導入。（完了）
6.  `main.go` に `MasterUseCase.DownloadAndStoreMasterData` の呼び出しを追加。（完了）

---

## 実装から得られた知見（データベース永続化編）

本セクションでは、開発過程で遭遇したGORMの挙動やマイグレーションに関する特殊なノウハウを記録する。

### 1. GORMでの一括Upsert (`Create` + `OnConflict`) とリレーションフィールドの問題

- **課題**: GORMのドメインモデル (`model.StockMaster`) がリレーションフィールド (`TickRules []model.TickRule`) を持つ場合、`.Clauses(clause.OnConflict{...}).Create(&models)` の形式で一括Upsertしようとすると、GORMが`TickRules`を物理テーブルのカラムとして扱おうとし、`invalid field found for struct ... TickRules`エラーが発生した。`.Omit("TickRules")`や`.Select(...)`もこの一括Upsertの文脈では期待通りに機能しなかった。
- **原因**: GORMは、モデル内にリレーションを示すフィールドが存在すると、デフォルトでそれらをDB操作の対象として含めようとする。しかし、一括Upsert (`Create` + `OnConflict`) の文脈では、子テーブルへの挿入/更新を自動で行う機能が限定的である、またはモデル定義のタグと内部挙動の間にミスマッチが生じたため、エラーが発生した。
- **解決策 (DTOパターン)**:
    - リポジトリ層の内部に、リレーションフィールドを持たない「データベース保存専用の構造体（DTO: `dbStockMaster`）」を定義した。
    - `UpsertStockMasters`メソッド内で、引数で受け取ったドメインモデル (`[]*model.StockMaster`) を、このDTOのリスト (`[]dbStockMaster`) に変換。この変換時にリレーションフィールド (`TickRules`) は物理的にコピーの対象から除外し、単純に無視した。
    - GORMの`Create`メソッドには、このリレーションを持たない`[]dbStockMaster`のリストを渡すことで、GORMが単純な構造体としてデータを扱い、`invalid field`エラーを確実に回避できた。
- **ノウハウ**: GORMの一括Upsertとリレーションフィールドの組み合わせは複雑な挙動を示す場合がある。その際は、DTOを導入して永続化層とドメインモデルを切り離し、永続化層の関心事を明確にすることで、問題を構造的に解決でき、かつGORMの一括Upsertのパフォーマンスを維持できる。

### 2. GORMのモデル定義と主キーの規約

- **課題**: GORMによるUpsertで`ERROR: column "id" does not exist`エラーが発生した。GORMはデフォルトで`RETURNING "id"`というSQLを生成しようとした。
- **原因**: `domain/model/master_base.go`に`ID uint gorm:"primarykey"`というフィールドが定義されており、これが`stock_masters`テーブルにも`id`という名前のオートインクリメント主キーカラムが存在することをGORMに期待させてしまっていた。しかし、`stock_masters`テーブルの実際の主キーは`issue_code`であった。
- **解決策**:
    - `domain/model/master_base.go`から`ID`フィールドを削除し、`MasterBase`は`CreatedAt`, `UpdatedAt`, `DeletedAt`といったGORMが自動管理するタイムスタンプフィールドのみを持つようにした。
    - 各ドメインモデル（例: `model.StockMaster`）は自身の主キーを`gorm:"primaryKey"`タグで明示的に定義する形となる。`stock_masters`テーブルの主キーは`issue_code`として正しく機能するようになった。
- **ノウハウ**: GORMは`ID`という名前の`primarykey`フィールドを特別扱いする傾向がある。カスタム主キーを持つテーブルのモデルを定義する場合、共通の基底構造体で`ID`を`primarykey`として定義せず、各モデルで明示的に主キーを指定する必要がある。

### 3. データベーススキーマの不整合とマイグレーションの重要性

- **課題**: Goの`model.StockMaster`に`issue_name_short`などの新しいカラムを追加した後、アプリケーション実行時に`ERROR: column "issue_name_short" of relation "stock_masters" does not exist`エラーが発生した。
- **原因**: Goのモデル（コード）は更新されたにもかかわらず、データベースのテーブルスキーマ（実際のDB構造）を更新していなかったため。GoのコードとDBのテーブル構造が一致していなかった。
- **解決策**:
    - `golang-migrate/migrate`ツールを使用し、`ALTER TABLE`文を含む新しいマイグレーションファイルを生成した。
    - `ALTER TABLE stock_masters ADD COLUMN ...`形式で、不足しているカラムをデータベースに追加するSQLを記述し、適用した。
- **ノウハウ**: Goのコードでドメインモデル（特にDBに永続化されるモデル）に構造的な変更（カラムの追加、型の変更など）を加える場合、必ずそれに対応するデータベースマイグレーションファイルを生成・適用し、**コードとDBのスキーマを常に一致させる**必要がある。このプロセスを怠ると、`column does not exist`のような予期せぬ実行時エラーの原因となる。

#### 追記：マイグレーションの一本化による管理の単純化
- **経緯**: 上記のように、開発の進行に伴い`ALTER TABLE`文を含む複数のマイグレーションファイルが作成された結果、マイグレーションの依存関係が複雑化し、環境差異によるエラーのリスクが高まった。
- **対策**: 開発初期段階の安定性を確保するため、リリースまではマイグレーションファイルを単一の初期スキーマファイルに統合する方針を採用。すべてのスキーマ定義を`CREATE TABLE`文に集約し、どの環境でも一度のマイグレーションで最新のスキーマを確実に構築できるようにした。

### 4. マイグレーション実行の自動化

- **課題**: `migrate`CLIツールのコマンドは長く、データベース接続文字列を直接記述する必要があるため、手動での実行は煩雑でエラーを起こしやすい。また、`migrate`CLIツールがインストールされていない環境では実行できない。
- **原因**: マイグレーションの実行がシェルコマンドと手動設定に依存していたため、ポータビリティと使いやすさが低かった。
- **解決策**: `.env`ファイルを読み込み、`golang-migrate/migrate`ライブラリをGoプログラム (`cmd/migrator/main.go`) で直接実行する仕組みを導入した。
- **ノウハウ**: `go run ./cmd/migrator/main.go`コマンド一つで、`.env`ファイルの設定に基づき、Goプログラム自身が環境に依存せず簡単にマイグレーションを実行できるようになった。これにより、デプロイや開発環境のセットアップが大幅に簡素化され、「他の環境ですぐに立ち上げたい」というニーズに応えられるようになった。

### 5. `ON CONFLICT`句の再発問題と恒久対策

- **過去の状況**: `TickRules`テーブルへのデータ挿入時に`ERROR: there is no unique or exclusion constraint matching the ON CONFLICT specification (SQLSTATE 42P10)`エラーが、特定の環境で再発した。これは`tick_rules`モデルの`TickUnitNumber`が`PRIMARY KEY`として定義されているにもかかわらず発生していた。
- **分析**:
    - 別の開発環境ではこの問題は発生していなかったため、環境ごとのデータベース状態の差異が原因と推察された。
    - GORMの`ON CONFLICT`句は、PostgreSQLが`UNIQUE`制約または`EXCLUSION`制約を明示的に認識している必要がある。`PRIMARY KEY`のみでは、GORMやDBドライバの挙動により不十分な場合があった。
    - 複数のマイグレーションファイルが存在したことで、環境によって`UNIQUE`制約が存在しないDBが構築されてしまうことが、根本原因の一つと考えられた。
- **恒久対策**:
    1.  **明示的なインデックス作成**: マイグレーションファイルに`CREATE UNIQUE INDEX IF NOT EXISTS idx_tick_rules_tick_unit_number ON tick_rules(tick_unit_number);`を追加し、`ON CONFLICT`句が依存するユニーク制約を明示的に保証した。
    2.  **マイグレーションの一本化**: さらに、上記の「3. マイグレーションによるスキーマ管理」の追記で述べた通り、マイグレーションファイルを一つに統合。これにより、環境差異によるスキーマの不整合リスクを排除し、問題の再発を根本的に防ぐ構成とした。
- **運用ルール**: データベースのスキーマで問題が発生した場合の最も確実な復旧方法は、Docker Volumeを含めてデータベースを完全に再作成し、統合されたマイグレーションを再適用することである。

