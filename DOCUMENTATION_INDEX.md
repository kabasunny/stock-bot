# ドキュメント一覧

株式取引システムの包括的なドキュメントです。目的に応じて適切なドキュメントを参照してください。

## 📋 概要・入門

| ドキュメント | 説明 | 対象者 |
|-------------|------|--------|
| [README.md](README.md) | プロジェクト概要とクイックスタート | 全員 |
| [SYSTEM_ARCHITECTURE_OVERVIEW.md](SYSTEM_ARCHITECTURE_OVERVIEW.md) | システム全体のアーキテクチャ図解 | 開発者・アーキテクト |
| [SYSTEM_FLOW_DOCUMENTATION.md](SYSTEM_FLOW_DOCUMENTATION.md) | システム全体の処理フロー詳細 | 開発者・運用者 |

## 🏗️ アーキテクチャ・設計

| ドキュメント | 説明 | 対象者 |
|-------------|------|--------|
| [CURRENT_ARCHITECTURE.md](CURRENT_ARCHITECTURE.md) | 現在の詳細な技術仕様 | 開発者 |
| [SESSION_MANAGEMENT_ARCHITECTURE.md](SESSION_MANAGEMENT_ARCHITECTURE.md) | セッション管理の設計思想 | 開発者 |
| [MULTI_BROKER_ARCHITECTURE.md](MULTI_BROKER_ARCHITECTURE.md) | マルチブローカー対応設計 | アーキテクト |

## 📊 戦略・セッション管理

| ドキュメント | 説明 | 対象者 |
|-------------|------|--------|
| [TACHIBANA_SESSION_STRATEGY.md](TACHIBANA_SESSION_STRATEGY.md) | 立花証券セッション戦略 | 開発者 |
| [TACHIBANA_DATE_BASED_SESSION_STRATEGY.md](TACHIBANA_DATE_BASED_SESSION_STRATEGY.md) | 日付ベースセッション戦略 | 開発者 |

## 🧪 テスト・品質

| ドキュメント | 説明 | 対象者 |
|-------------|------|--------|
| [TEST_PLAN.md](TEST_PLAN.md) | 包括的なテスト戦略と進捗 | 開発者・QA |
| [DETAILED_TEST_PLAN.md](DETAILED_TEST_PLAN.md) | 詳細なテスト仕様 | QA・テスター |
| [TEST_PROGRESS_TRACKER.md](TEST_PROGRESS_TRACKER.md) | テスト進捗管理 | プロジェクトマネージャー |

## 🔄 リファクタリング・改善

| ドキュメント | 説明 | 対象者 |
|-------------|------|--------|
| [REFACTORING_PLAN.md](REFACTORING_PLAN.md) | リファクタリング計画 | 開発者 |
| [ARCHITECTURE_REFACTORING_PROGRESS.md](ARCHITECTURE_REFACTORING_PROGRESS.md) | アーキテクチャ改善進捗 | アーキテクト |

## 📈 改善・TODO

| ドキュメント | 説明 | 対象者 |
|-------------|------|--------|
| [TODO_MASTER_DATA_IMPROVEMENTS.md](TODO_MASTER_DATA_IMPROVEMENTS.md) | マスターデータ改善項目 | 開発者 |
| [LIGHTWEIGHT_VERSION_PROPOSAL.md](LIGHTWEIGHT_VERSION_PROPOSAL.md) | 個人用軽量版システム提案 | 個人投資家・コスト重視 |
| [TEST_ENVIRONMENT_STRATEGY.md](TEST_ENVIRONMENT_STRATEGY.md) | テスト環境戦略 | DevOps・開発者 |

## 📝 API・仕様

| ファイル | 説明 | 対象者 |
|---------|------|--------|
| [api_commands.md](api_commands.md) | API コマンド一覧 | 開発者・運用者 |

## 🎯 ドキュメント選択ガイド

### 初めてプロジェクトに参加する場合
1. [README.md](README.md) - プロジェクト概要を把握
2. [SYSTEM_ARCHITECTURE_OVERVIEW.md](SYSTEM_ARCHITECTURE_OVERVIEW.md) - 全体像を理解
3. [SYSTEM_FLOW_DOCUMENTATION.md](SYSTEM_FLOW_DOCUMENTATION.md) - 処理フローを学習
4. [CURRENT_ARCHITECTURE.md](CURRENT_ARCHITECTURE.md) - 詳細な実装を学習

### 開発を始める場合
1. [SYSTEM_FLOW_DOCUMENTATION.md](SYSTEM_FLOW_DOCUMENTATION.md) - 処理フローを理解
2. [CURRENT_ARCHITECTURE.md](CURRENT_ARCHITECTURE.md) - 技術仕様を確認
3. [TEST_PLAN.md](TEST_PLAN.md) - テスト戦略を理解
4. [SESSION_MANAGEMENT_ARCHITECTURE.md](SESSION_MANAGEMENT_ARCHITECTURE.md) - セッション管理を学習

### テスト・品質保証の場合
1. [TEST_PLAN.md](TEST_PLAN.md) - テスト戦略全体を把握
2. [DETAILED_TEST_PLAN.md](DETAILED_TEST_PLAN.md) - 詳細なテスト仕様を確認
3. [TEST_PROGRESS_TRACKER.md](TEST_PROGRESS_TRACKER.md) - 進捗を管理

### アーキテクチャ設計・改善の場合
1. [SYSTEM_ARCHITECTURE_OVERVIEW.md](SYSTEM_ARCHITECTURE_OVERVIEW.md) - 現在の設計を把握
2. [MULTI_BROKER_ARCHITECTURE.md](MULTI_BROKER_ARCHITECTURE.md) - 拡張設計を理解
3. [REFACTORING_PLAN.md](REFACTORING_PLAN.md) - 改善計画を確認

### 運用・保守の場合
1. [README.md](README.md) - セットアップ手順を確認
2. [api_commands.md](api_commands.md) - API操作方法を学習
3. [TEST_ENVIRONMENT_STRATEGY.md](TEST_ENVIRONMENT_STRATEGY.md) - 環境戦略を理解

## 📊 ドキュメント状態

| カテゴリ | 完成度 | 最終更新 |
|---------|--------|----------|
| 概要・入門 | ✅ 完成 | 2026-01-01 |
| アーキテクチャ | ✅ 完成 | 2026-01-01 |
| テスト・品質 | ✅ 完成 | 2026-01-01 |
| API・仕様 | 🔄 更新中 | 2025-12-30 |
| 改善・TODO | 🔄 更新中 | 2025-12-30 |

## 🔄 ドキュメント更新ルール

### 更新頻度
- **README.md**: 機能追加時
- **アーキテクチャ文書**: 設計変更時
- **テスト文書**: テスト追加・変更時
- **TODO文書**: 課題発見・解決時

### 更新責任者
- **概要文書**: プロジェクトリーダー
- **技術文書**: 担当開発者
- **テスト文書**: QA担当者
- **改善文書**: アーキテクト

### レビュープロセス
1. ドキュメント更新
2. 関係者レビュー
3. 承認・マージ
4. 関係者への通知

## 📞 サポート

ドキュメントに関する質問や改善提案は以下まで：

- **GitHub Issues**: 技術的な質問・バグ報告
- **Pull Request**: ドキュメント改善提案
- **Wiki**: FAQ・よくある質問

---

**注意**: このドキュメント一覧は定期的に更新されます。新しいドキュメントが追加された場合は、このファイルも併せて更新してください。