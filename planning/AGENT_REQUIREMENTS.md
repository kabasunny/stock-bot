# エージェント要件定義

## 0. 前提となるシステム構成とツール群

エージェントは、独立した複数の「ツール」を適切なタイミングで利用し、全体のワークフローを指揮する「司令塔」として機能する。

### 0.1. ツール群の定義

| ツール名 | 実装言語 | 主な役割 |
| :--- | :--- | :--- |
| **モデルメーカー** | Python | **[重い処理]** 長期データから特徴量を生成し、機械学習モデルを構築・更新する。 |
| **シグナルメーカー** | Python | **[軽い処理]** 学習済みモデルを使い、その日の**売買シグナル**を生成し、以下のバイナリフォーマットでファイルに出力する。<br>- **レコード**: 銘柄コード(`uint16`) + 売買区分(`uint8`) の3バイト構成。<br>- **売買区分**: `0x01`=BUY, `0x02`=SELL。 |
| **パラメータオプティマイザー** | Go | バックテストを行い、**最適な利確・損切りパラメータ**を計算・提供する。 |
| **トレードサービス** | Go (Goa) | 証券会社APIと通信し、**注文実行**や**DBへの状態永続化**を行う。 |

### 0.2. エージェントのワークフロー

エージェントは、これらのツールを以下のように利用して取引を実行する。

1.  **準備フェーズ (取引開始前など)**:
    -   **パラメータ設定**: **パラメータオプティマイザー**が算出した最適パラメータを `config.yaml` から読み込む。
    -   **シグナル生成**: **シグナルメーカー**を実行し、その日の売買シグナルが記述された「シグナルファイル」を取得する。

2.  **実行フェーズ (取引時間中など)**:
    -   **意思決定**: 取得した「シグナルファイル」と「最適パラメータ」に基づき、具体的な注文内容（銘柄、数量、注文方法など）を決定する。
    -   **注文実行**: **トレードサービス**を呼び出し、意思決定した内容で発注を指示する。


本ドキュメントは、株式自動取引システムの「エージェント(Agent)」コンポーネントの要件を定義する。

## 1. ユースケースストーリー: エージェントの一日

「状態管理」が実際の取引でどのように機能するのかを、具体的なストーリー形式で以下に示す。

---
### **とあるエージェントの一日**

#### **【9:00前】 起動と準備**
1.  エージェントは起動すると、まず自身の状態を**「準備中」**にします。
2.  次に、**状態管理**の一環として、「トレードサービス」を通じて証券会社に現在の正確な**口座情報**を問い合わせます。
    *   「現在の現金残高は1,000万円か…」
    *   「前日から持ち越しているA社の株が100株あるな…」
3.  エージェントはこれらの情報を自身の「記憶（状態）」にしっかりと記録します。
4.  準備が完了すると、自身の状態を**「監視中」**に変え、取引開始を待ちます。

#### **【9:30】 シグナル発生**
1.  エージェントは、定期巡回（tick）で「シグナルファイル」をチェックし、「**B社の株を買い**」というシグナルを見つけます。
2.  ここで**意思決定ロジック**が働きます。
    *   **（状態確認①）**：「B社の株は、今すでに保有しているだろうか？」
        *   → 自身の「記憶（状態）」を確認。「いや、持っていないな。よし、買ってOKだ。」
    *   **（状態確認②）**：「B社の株を買うお金は足りるだろうか？」
        *   → 起動時に記憶した現金残高と現在の株価を比較。「うん、十分足りるな。」
3.  全ての確認が取れたため、エージェントは「トレードサービス」に**B株の買い注文**を指示します。
4.  発注後、すぐに自身の状態を更新します。「現在、B株の買い注文が**発注中**」と記憶します。

#### **【10:00】 約定の確認**
1.  エージェントは「トレードサービス」を通じて、「B株の買い注文が約定した」ことを知ります。
2.  すぐさま自身の状態を更新します。
    *   「発注中リスト」からB株の注文を削除。
    *   「**保有ポジションリスト**」に「B株 100株（取得単価 XXX円）」を追加。
    *   「現金残高」から買付代金を差し引く。
    *   → これでエージェントは、「自分がいま何を持っているか」を正確に把握できました。

#### **【11:00】 ２つ目のシグナル**
1.  エージェントは、巡回中に再び「**B社の株を買い**」というシグナルを見つけました。
2.  **意思決定ロジック**が働きます。
    *   **（状態確認）**：「B社の株は、今すでに保有しているだろうか？」
        *   → 自身の「記憶（状態）」を確認。「おっと、さっき買ったばかりだ。重複して買うのはリスクが高いから、このシグナルは**見送ろう**。」
    *   → **状態管理機能がなければ、ここで不要な重複買いをしてしまっていました。**

#### **【14:30】 予期せぬ急落と自動損切り**
1.  この日のシグナルはもうありません。
2.  しかし、エージェントは常に自身の「**保有ポジション**」を監視しています。B株の株価が下落し、`config.yaml`で決められた「損切り率（-2%）」に達してしまいました。
3.  **意思決定ロジック**（リスク管理部分）が作動し、「保有しているB株を全て**損切りのために売る**」と即座に判断。
4.  エージェントは「トレードサービス」に売り注文を出し、損失の拡大を防ぎました。

---

このように、エージェントが**過去の行動や現在の状況を記憶（状態管理）**することで、単純にシグナルに従うだけでなく、より安全で賢明な取引判断が可能になります。


## 1. 概要

エージェントはシステム全体の「頭脳」として、以下の役割を担う。

-   **役割**: 主体的な意思決定を行う中央司令塔。
-   **機能**: 取引戦略のメインループとして、Go APIラッパー（データ取得/注文）とPythonサービス（シグナル問合せ）を呼び出し、全体のワークフローを指揮する。

## 2. アーキテクチャ方針

-   **共通プラットフォーム + 戦略モジュール**:
    -   異なる取引戦略（デイトレード、スイングトレード等）で共通して利用可能な「プラットフォーム」としてのコア機能を実装する。
    -   個別の取引戦略は、このプラットフォーム上で動作する「戦略モジュール」として実装し、交換可能（pluggable）な設計を目指す。
    -   初期実装では、共通プラットフォームの構築と、「スイングトレード」戦略モジュールの実装を目標とする。シグナル生成については、当面Pythonサービスとのリアルタイム連携ではなく、外部から提供されるシグナルファイルを読み込む方式を採用する。

## 3. 機能要件 (共通プラットフォーム)

### 3.1. 設定管理 (Configuration)
-   **要件**:
    -   取引戦略のパラメータ（取引対象銘柄リスト、ロットサイズ、利確/損切り幅など）を外部ファイル（例: `config.yaml`）から読み込めること。
    -   APIのエンドポイントや認証情報、ログレベルなどのシステム設定も同ファイルで管理できること。
-   **目的**: コードを変更せずにエージェントの振る舞いを調整できるようにするため。

### 3.2. 状態管理 (State Management)
-   **要件**:
    -   現在のポジション（保有銘柄、数量、平均取得単価）をメモリ上で管理できること。
    -   発注中の注文情報（未約定注文）を管理できること。
    -   口座の資金余力（現金、買付可能額）を管理できること。
    -   エージェント自身の内部状態（例: `初期化中`, `待機中`, `シグナル待ち`, `発注中`, `エラー停止`）を管理できること。
-   **目的**: 意思決定の基礎となる、エージェントが置かれている状況を正確に把握するため。

### 3.3. 実行ループ (Execution Loop)
-   **要件**:
    -   市場の取引時間に合わせて、定期的にワークフローを実行するスケジューラを持つこと。
    -   実行間隔は設定ファイルで変更可能であること。
    -   ワークフロー（シグナルファイル読み込み→意思決定→発注）を制御するメインループを実装すること。
-   **目的**: 取引戦略を自律的かつ継続的に実行するため。

### 3.4. 意思決定ロジック (Decision Logic)
-   **要件**:
    -   戦略モジュール（シグナルファイル）から渡されたシグナル（銘柄コード、売買区分）と、現在の状態（ポジション有無、資金余力、市場状況など）を組み合わせて、最終的な行動を決定するコアロジックを持つこと。
    -   **シグナル解釈と注文方法の決定**:
        -   シグナルファイルから受け取った「銘柄コード」と「売買区分（BUY/SELL）」に基づき、現在の市場状況やリスク許容度に応じて、最適な「注文種別（MARKET/LIMIT）」および「数量」を決定する。
        -   指値注文の場合、適切な指値価格を算出する。
    -   **リスク管理（自動利確・損切り）**:
        -   保有ポジションに対し、設定ファイル（`config.yaml`）で定義された利確率（`profit_take_rate`）と損切り率（`stop_loss_rate`）に基づき、自動的に利確・損切り注文を出す判断を行う。
        -   利確・損切り判断は、戦略モジュールからの明確な手仕舞いシグナルがない場合でも実行される。
    -   戦略モジュールからの指示と、リスク管理ロジックの間で競合が発生した場合の優先順位付けロジックを持つこと。
-   **目的**: 戦略モジュールの判断を、実際の取引アクションに変換し、エージェントが自律的なリスク管理を伴う取引判断を行えるようにするため。

### 3.5. API連携 (API Integration)
-   **要件**:
    -   Goaで定義された各APIクライアント（`BalanceClient`, `OrderClient`, `PriceInfoClient`など）を呼び出せること。
    -   取引シグナルが記述された外部ファイル（`.bin`形式などを想定）をパース（解析）して、取引指示を読み込めること。
-   **目的**: 外部コンポーネントと連携し、ワークフローを実行するため。

### 3.6. エラーハンドリングと耐障害性 (Error Handling & Resilience)
-   **要件**:
    -   API呼び出しの失敗（ネットワークエラー、APIエラーレスポンス）に対して、設定可能なリトライ処理を行えること。
    -   セッション切れを検知した場合、再ログインを試みること（本番環境の電話認証制約を考慮し、失敗時は管理者に通知する）。
    -   予期せぬエラーが発生した場合、安全に停止（例: 全ての未約定注文をキャンセルし、エージェントを停止）するフェールセーフ機構を持つこと。
-   **目的**: システムの安定稼働と、想定外の事態による損失を防ぐため。

### 3.7. ロギング (Logging)
-   **要件**:
    -   エージェントの全ての主要な決定と行動（シグナル受信、発注、エラー、状態遷移など）を、構造化ログ（JSON形式など）として記録できること。
    -   ログレベル（DEBUG, INFO, WARN, ERROR）を設定ファイルで変更できること。
-   **目的**: 稼働状況の監視と、問題発生時の原因究明を容易にするため。

## 4. 未解決の課題
-   （ここに議論で出てきた課題などを追記していく）







# ポジションサイジングと意思決定ロジックの参考例と解説

このドキュメントでは、提供されたGoコードベースから抽出した、トレーディングにおけるポジションサイジング（建玉量の決定）と意思決定（エントリーおよびエグジット）に関するロジ-ックについて解説します。

## 概要

このトレーディングシステムは、事前に外部から与えられた売買シグナル（買いの日付）に基づいて動作します。システムのコアロジックは、シグナルを生成することではなく、「どのように取引を実行し、リスクを管理するか」という点にあります。

主な意思決定のフローは以下の通りです。

1.  **シグナルの受信**: 取引すべき銘柄と日付のリストを受け取ります。
2.  **ポジションサイズの決定**: ポートフォリオ全体のリスクを考慮し、ATR（Average True Range）を用いて1取引あたりの適切な建玉量を計算します。
3.  **エントリー**: シグナルで指定された日の始値で株式を購入します。
4.  **エグジット**: あらかじめ定義された2種類の方法（静的なストップロス、または利益が出た後に有効になるトレーリングストップ）に基づいて、保有ポジションを売却します。

以下に、このロジックを実装している主要な関数とコードを解説します。

---

## 1. ポジションサイジングのロジック (`determine_position_size.go`)

ポジションサイズは、有名なリスク管理手法である「Van Tharpモデル」に基づいて決定されます。これは、1回のトレードで許容できる最大損失額を事前に決め、それに基づいて建玉量を調整するアプローチです。

このロジックは `DeterminePositionSize` 関数に実装されています。

### コード

```go
package trading

import (
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"time"
)

// determinePositionSize は、ATRに基づきポジションサイズとエントリー価格、エントリーコストを決定
func DeterminePositionSize(param *ml_stockdata.Parameter, portfolioValue int, availableFundsInt int, entryPrice float64, commissionRate *float64, dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) (float64, float64, error) {

	const unitSize = 100 // 単元数
	availableFunds := float64(availableFundsInt)

	// ATRを計算
	atr := calculateATR(dailyData, signalDate)
	// fmt.Println("ATR:", atr)

	// 許容損失額を計算 (ポートフォリオ価値のストップロス割合)
	allowedLoss := float64(portfolioValue) * (param.RiskPercentage / 100)

	// ストップロス幅をATRの2倍に設定（過去の価格変動の2倍の幅でストップロスを設定）
	stopLossAmount := atr * param.ATRMultiplier

	// 初期ポジションサイズを計算
	initialPositionSize := allowedLoss / stopLossAmount
	// fmt.Println("positionSize before unit size:", positionSize)

	// ポジションサイズを調整して最小単元の倍数にする
	initialPositionSize = math.Floor(initialPositionSize/float64(unitSize)) * float64(unitSize)
	// fmt.Println("positionSize after unit size:", positionSize)

	// 手数料を加味してエントリーコストを計算
	initialEntryCost := entryPrice * initialPositionSize
	commission := initialEntryCost * (*commissionRate / 100)
	initialTotalEntryCost := initialEntryCost + commission
	// fmt.Println("totalEntryCost:", totalEntryCost)

	// 使用可能な資金に対してエントリーコストが足りるか、
	// かつ、ポートフォリオのリスク許容範囲を超えないようにポジションサイズを調整
	maxPositionSize := initialPositionSize

	if initialTotalEntryCost > availableFunds {
		// 利用可能資金を超える場合、ポジションサイズを縮小
		maxPositionSize = math.Floor((availableFunds/(entryPrice*(1+(*commissionRate/100))))/float64(unitSize)) * float64(unitSize)
		// エントリーコストを再計算
		initialEntryCost = entryPrice * maxPositionSize
		commission = initialEntryCost * (*commissionRate / 100)
		initialTotalEntryCost = initialEntryCost + commission
	}

	// ポートフォリオのリスク許容範囲を超えないようにポジションサイズを調整
	riskLimitEntryCost := float64(portfolioValue) * param.RiskPercentage
	if initialTotalEntryCost > riskLimitEntryCost {
		maxPositionSize = math.Floor((riskLimitEntryCost/(entryPrice*(1+(*commissionRate/100))))/float64(unitSize)) * float64(unitSize)
		// エントリーコストを再計算
		initialEntryCost = entryPrice * maxPositionSize
		commission = initialEntryCost * (*commissionRate / 100)
		initialTotalEntryCost = initialEntryCost + commission
	}

	// 最終的なポジションサイズとエントリーコスト
	positionSize := maxPositionSize
	totalEntryCost := initialTotalEntryCost

	// ポジションサイズがゼロ以下の場合、エントリーしない
	if positionSize <= 0 {
		return 0, 0, nil
	}

	return positionSize, totalEntryCost, nil

}

// calculateATR は、過去一定期間のATR（Average True Range）を計算する
func calculateATR(dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) float64 {
	// ATRの計算ロジック（過去n日間のTrue Range平均）
	n := 14 // 計算に使用する日数
	trueRanges := make([]float64, 0, n)

	// signalDate以前のn日間のデータを収集
	for i := len(*dailyData) - 1; i >= 1; i-- { // i >= 1 に変更 (yesterdayDataのために最低2つのデータが必要)
		data := (*dailyData)[i]
		date, _ := time.Parse("2006-01-02", data.Date)

		// signalDateより後のデータはスキップ
		if date.After(signalDate) {
			continue
		}

		// signalDate当日のデータもスキップ
		if date.Equal(signalDate) {
			continue
		}

		yesterdayData := (*dailyData)[i-1]
		trueRange := calculateTrueRange(data, yesterdayData)
		trueRanges = append([]float64{trueRange}, trueRanges...) // 先頭に追加
		//trueRanges = append(trueRanges, trueRange)
		if len(trueRanges) >= n {
			break
		}
	}

	if len(trueRanges) == 0 {
		fmt.Println("ATR計算に必要なデータが不足しています。エントリーを見送ります。")
		return 0 // ATRが計算できない場合は、0を返す（ポジションサイズが0になる）
	}

	// ATRを計算
	sum := 0.0
	for _, tr := range trueRanges {
		sum += tr
	}
	atr := sum / float64(len(trueRanges))
	return atr
}

// calculateTrueRange は、前日と当日のデータに基づいてTrue Rangeを計算
func calculateTrueRange(today, yesterday ml_stockdata.InMLDailyData) float64 {
	highLow := today.High - today.Low
	highClose := math.Abs(today.High - yesterday.Close)
	lowClose := math.Abs(today.Low - yesterday.Close)
	trueRange := math.Max(highLow, math.Max(highClose, lowClose))
	return trueRange
}
```

### 解説

1.  **ATRの計算 (`calculateATR`)**:
    *   `calculateATR` は、過去14日間のATR（Average True Range）を計算します。ATRは株価のボラティリティ（変動幅）を測る指標です。
    *   `calculateTrueRange` で1日あたりの真の変動幅を計算し、その平均をATRとしています。

2.  **1トレードあたりの許容損失額の計算**:
    *   `allowedLoss := float64(portfolioValue) * (param.RiskPercentage / 100)`
    *   ポートフォリオ全体の価値（`portfolioValue`）に対して、パラメータで指定された一定の割合（`RiskPercentage`）を掛け合わせ、この1回の取引で失ってもよい最大金額を算出します。例えば、ポートフォリオが100万円で`RiskPercentage`が2%なら、許容損失額は2万円となります。

3.  **1株あたりのリスク額（ストップロス幅）**:
    *   `stopLossAmount := atr * param.ATRMultiplier`
    *   ATRにパラメータで指定された倍率（`ATRMultiplier`）を掛けて、1株あたりの想定リスク（エントリー価格から損切りラインまでの幅）を決定します。ATRが大きい（変動が激しい）銘柄ほど、リスク額は大きくなります。

4.  **初期ポジションサイズの計算**:
    *   `initialPositionSize := allowedLoss / stopLossAmount`
    *   「1トレードあたりの許容損失額」を「1株あたりのリスク額」で割ることで、購入すべき株数（ポジションサイズ）を計算します。

5.  **ポジションサイズの調整**:
    *   `initialPositionSize = math.Floor(initialPositionSize/float64(unitSize)) * float64(unitSize)`
    *   計算されたポジションサイズを、取引所の最小取引単位（`unitSize`、ここでは100株）の倍数になるように切り捨てます。
    *   さらに、算出したポジションの合計コスト（`initialTotalEntryCost`）が、利用可能な現金（`availableFunds`）やポートフォリオのリスク許容範囲（`riskLimitEntryCost`）を超えないように、ポジションサイズを下方修正します。これにより、資金不足や過大なリスクを避けます。

最終的に、すべての制約を満たした `positionSize`（株数）と `totalEntryCost`（合計コスト）が返されます。

---

## 2. エグジット（売却）の意思決定ロジック (`stop_order_utils.go`)

ポジションをいつ手仕舞うかは、`findExitDate` 関数によって決定されます。この関数は、静的なストップロスと、利益が出た後に発動するトレーリングストップの2つのルールを実装しています。

### コード

```go
package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"time"
)

// roundUp 関数: 四捨五入（切り上げ）
func roundUp(value float64) float64 {
	return float64(int(value*10+1)) / 10
}

// roundDown 関数: 四捨五入（切り捨て）
func roundDown(value float64) float64 {
	return float64(int(value*10)) / 10
}

// findExitDate: 売却日と売却価格を決定
func findExitDate(data []ml_stockdata.InMLDailyData, purchaseDate time.Time, purchasePrice float64, param *ml_stockdata.Parameter) (time.Time, float64, error) {
	var endDate time.Time
	var endPrice float64

	// パラメータを保存
	stopLossPercentage := param.StopLossPercentage
	trailingStopTrigger := param.TrailingStopTrigger
	trailingStopUpdate := param.TrailingStopUpdate

	// ストップロスとトレーリングストップの閾値を計算
	stopLossThreshold := roundDown(purchasePrice * (1 - stopLossPercentage/100))
	trailingStopTriggerPrice := roundUp(purchasePrice * (1 + trailingStopTrigger/100))

	// purchaseDate 以降のデータを取得（スライスを最適化）
	var startIndex int
	for i, day := range data {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}
		if !parsedDate.Before(purchaseDate) {
			startIndex = i
			break
		}
	}
	filteredData := data[startIndex:]

	// トレーリングストップの監視
	for _, day := range filteredData {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}

		openPrice := day.Open
		lowPrice := day.Low
		closePrice := day.Close

		// ストップロスに到達した場合
		if lowPrice <= stopLossThreshold || openPrice <= stopLossThreshold {
			endPrice = stopLossThreshold
			endDate = parsedDate
			break
		}

		// トレーリングストップのトリガーをチェック
		if closePrice >= trailingStopTriggerPrice {
			trailingStopTriggerPrice = roundUp(closePrice * (1 + trailingStopTrigger/100))
			stopLossThreshold = roundDown(closePrice * (1 - trailingStopUpdate/100))
		}
	}

	// 途中で売却しなかった場合、最終データを採用
	if endDate.IsZero() {
		lastIndex := len(filteredData) - 1
		endPrice = filteredData[lastIndex].Close
		endDate, _ = parseDate(filteredData[lastIndex].Date)
	}

	return endDate, endPrice, nil
}
```

### 解説

1.  **初期損切りラインの設定**:
    *   `stopLossThreshold := roundDown(purchasePrice * (1 - stopLossPercentage/100))`
    *   購入価格（`purchasePrice`）から、パラメータで指定された静的な損切り率（`stopLossPercentage`）だけ低い価格を、最初の損切りラインとして設定します。

2.  **トレーリングストップの発動条件**:
    *   `trailingStopTriggerPrice := roundUp(purchasePrice * (1 + trailingStopTrigger/100))`
    *   購入価格が一定の利益率（`trailingStopTrigger`）に達したら、トレーリングストップが発動します。この価格が `trailingStopTriggerPrice` です。

3.  **日々の価格監視**:
    *   購入日以降、日足データのループ処理に入ります。
    *   `if lowPrice <= stopLossThreshold || openPrice <= stopLossThreshold`:
        *   その日の安値または始値が現在の損切りライン (`stopLossThreshold`) を下回った場合、その時点で損切りが確定し、ループを抜けます。
    *   `if closePrice >= trailingStopTriggerPrice`:
        *   終値がトレーリングストップの発動価格を上回った場合、トレーリングストップが有効になります。
        *   損切りライン (`stopLossThreshold`) が、現在の終値から一定割合（`trailingStopUpdate`）だけ低い価格に引き上げられます。
        *   同時に、次のトレーリングストップ発動価格もさらに高い価格に更新されます。
    *   これにより、利益が伸びる限り損切りラインを切り上げ続け、損失を限定しつつ利益を確保しようとします。

4.  **最終的な売却**:
    *   もし期間中に一度も損切りラインに掛からなかった場合は、データ期間の最終日の終値で売却されます。

---

## 3. 全体戦略の実行 (`trading_strategy.go`)

`TradingStrategy` 関数は、これまでのロジックを統合し、バックテスト全体を管理するオーケストレーターです。

### コード（抜粋）

```go
// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて最適なパラメータの組み合わせを見つける
func TradingStrategy(response *ml_stockdata.InMLStockResponse, totalFunds *int, param *ml_stockdata.Parameter, commissionRate *float64, options ...bool) (ml_stockdata.OptimizedResult, error) {
    // ... (変数の初期化) ...

	// シグナルを日付順、優先順にソート
	sort.Slice(signals, func(i, j int) bool {
		// ...
	})

    // ... (ポートフォリオ変数の初期化) ...

	// ---- シグナルの処理 ----
	for _, signal := range signals {
        // ... (既存ポジションの決済処理) ...

        // ... (利用可能資金の計算) ...

		for _, symbolData := range response.SymbolData {
			if symbolData.Symbol != signal.Symbol {
				continue
			}
            // singleTradingStrategyは内部でfindPurchaseDateとfindExitDateを呼び出す
			purchaseDate, exitDate, profitLoss, entryPrice, exitPrice, err := singleTradingStrategy(
				&symbolData.DailyData, signal.SignalDate, param,
			)
			if err != nil {
				continue
			}

            // ポジションサイズを決定
			positionSize, entryCost, err := DeterminePositionSize(param, portfolioValue, availableFunds, entryPrice, commissionRate, &symbolData.DailyData, signal.SignalDate)
			if err != nil || entryCost == 0 {
				continue
			}

            // ... (取引記録の作成とアクティブトレードへの追加) ...
		}
	}

    // ... (最終的なパフォーマンス指標の計算) ...

	return result, nil
}
```

### 解説

1.  **シグナルのソート**:
    *   受け取ったすべてのシグナルを、日付と優先度に基づいてソートします。これにより、時系列に沿ってバックテストが実行されます。

2.  **シグナルのループ処理**:
    *   ソートされたシグナルを一つずつ処理します。
    *   ループの冒頭で、現在のシグナル日付より前に決済されるべき既存のポジションを処理し、ポートフォリオの価値（`portfolioValue`）と利用可能資金（`availableFunds`）を更新します。

3.  **ポジションサイズ決定とエントリー**:
    *   `singleTradingStrategy` を呼び出して、エントリー価格や仮の exitDate を取得します。
    *   次に、`DeterminePositionSize` を呼び出して、計算されたエントリー価格と現在のポートフォリオ状況に基づき、実際に購入する株数を決定します。
    *   資金が不足している場合や、計算されたポジションサイズが0の場合は、取引を見送ります。

4.  **結果の記録**:
    *   取引が実行された場合、その内容（エントリー日、コスト、株数など）を `activeTrades` マップに記録します。
    *   ループが完了した後、すべての取引結果を集計し、勝率、平均利益/損失、最大ドローダウンなどのパフォーマンス指標を計算して返します。
