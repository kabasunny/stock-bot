# ログイン仕様に関するテスト計画書

## 1. 目的

ログイン、ログアウト、およびセッション管理に関連するAPIの挙動を明確にし、不明点を解消する。
特に、以下の点についての仕様を特定することを目的とする。

- 公開情報より手動認証後、ログインまでの猶予は3分である
- 上記の以外のシチュエーションの時間や、挙動を洗い出す
- 

## 2. テスト項目一覧

---

### テストID: TEST-001
- **テスト内容**: 本番URL連続2回ログイン
- **前提条件**: 初回ログイン前に手動認証後2分以内
- **実行手順**:
  1. 手動で電話認証を行う。
  2. 1回目のログインを実行し、`p_no`が`1`であることを確認する。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LoginOnly
     ```
  3. 直後に2回目のログインを実行する。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LoginOnly
     ```
- **期待される結果**: 2回目のログインも成功し、新しいセッションが開始される。ログに出力される`p_no`が再び`1`になる。
- **実行結果**: 1回目のログイン、および直後の2回目のログインともに成功しました。それぞれのログイン後、ログには「`ログイン成功後の p_no: 1`」と表示されました。
- **ステータス**: 成功

---

### テストID: TEST-002
- **テスト内容**: タイムラグを伴った本番URL連続2回ログイン
- **前提条件**: 初回ログイン前に手動認証済み。
- **実行手順**:
  1. 手動で電話認証を行う。
  2. 1回目のログインを実行する。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LoginOnly
     ```
  3. **5分間待機する。**
  4. 2回目のログインを実行する。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LoginOnly
     ```
- **期待される結果**: 1回目のログインから5分後には電話認証の有効期間（3分）が切れているため、2回目のログインは失敗する。
- **実行結果**: 1回目のログインは成功しましたが、5分以上経過してからの2回目のログインは「`login failed with result code 10089: 当社に登録の電話番号から認証電話番号へかけた後、3分以内にログインしてください。`」というエラーで失敗しました。
- **ステータス**: 成功

---

### テストID: TEST-003
- **テスト内容**: ログアウト後の再ログイン
- **前提条件**: 初回ログイン前の手動認証から2分以内に行う
- **実行手順**:
  1. 手動で電話認証を行う。
  2. ログインとログアウトを実行する（`TestAuthClientImpl_LogoutOnly`は内部でログインも行います）。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LogoutOnly
     ```
  3. 直後に再ログインを実行する。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LoginOnly
     ```
- **期待される結果**: すべての操作が電話認証の有効期間（3分）内に行われれば、再ログインは成功する。`p_no`は`1`になる。
- **実行結果**: `TestAuthClientImpl_LogoutOnly`（ログイン→ログアウト）と、その直後の`TestAuthClientImpl_LoginOnly`（再ログイン）が共に成功した。
- **ステータス**: 成功

---

### テストID: TEST-004
- **テスト内容**: タイムラグを伴ったログアウト後の再ログイン
- **前提条件**: 初回ログイン前に手動認証済み。
- **実行手順**:
  1. 手動で電話認証を行う。
  2. 新しく追加したシーケンステストを実行する。このテストは内部でログイン、5分待機、ログアウトを行います。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_Sequence_LoginWaitLogoutLogin
     ```
  3. テストのログに出力される`p_no`の挙動と、ログアウトAPIの結果を確認する。
- **期待される結果**: ログインから5分後には電話認証の有効期間（3分）が切れているため、ログアウト後の再ログインは失敗する。
- **実行結果**: `TestAuthClientImpl_Sequence_LoginWaitLogoutLogin` テストを実行。ログイン→5分待機→ログアウトのシーケンスは成功。その後の再ログインは、期待通り `result code 10089` エラーで失敗した。
- **ステータス**: 成功（APIの挙動を解明）

---

### テストID: TEST-005
- **テスト内容**: 無通信タイムアウトの確認
- **前提条件**: 本番環境に接続可能。
- **実行手順**:
  1. `price_info_client_impl_prod_test.go`に、ログイン・30分待機・株価照会を連続して行う新しいテスト`TestPriceInfo_Sequence_LoginWaitGetPrice`を追加する。（次のステップで実施します）
  2. 電話認証を行う。
  3. 追加したテストを実行する。このテストは完了まで30分以上かかる。
     ```bash
     go test -v ./internal/infrastructure/client/tests/price_info_client_impl_prod_test.go -run TestPriceInfo_Sequence_LoginWaitGetPrice -timeout 40m
     ```
- **期待される結果**: もし無通信タイムアウトの仕様があれば、30分後の株価照会はセッションエラーとなり失敗する。もし無ければ成功する。どちらの結果になるかを確認することが目的。
- **実行結果**:
- **ステータス**: 未実施

---

### テストID: TEST-006
- **テスト内容**: セッションの競合（二重ログイン）
- **前提条件**: 本番環境に接続可能。ターミナルを2つ開けること。
- **実行手順**:
  1. **【ターミナルA】** で電話認証を行い、ログインする。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LoginOnly
     ```
  2. **【ターミナルB】** で電話認証を行い、ログインする。
     ```bash
     go test -v ./internal/infrastructure/client/tests/auth_client_impl_test.go -run TestAuthClientImpl_LoginOnly
     ```
  3. **【ターミナルA】** に戻り、セッションが有効か確認するため、株価照会のテストを実行する。
     ```bash
     go test -v ./internal/infrastructure/client/tests/price_info_client_impl_prod_test.go -run TestPriceInfoClientImpl_Production
     ```
- **期待される結果**: ターミナルBのログインによってターミナルAのセッションは無効化されているはず。そのため、ターミナルAでの最後の株価照会テストは、ログインの段階で失敗する（再度、電話認証を求められる）か、何らかのセッションエラーが発生する。
- **実行結果**: `TestAuthClientImpl_MultipleSessions` テストにより、新しいログインセッションが確立されると、それ以前のセッションはサーバー側で無効化されることが確認された。
- **ステータス**: 成功

---







### (以降、必要に応じてテスト項目を追加)
