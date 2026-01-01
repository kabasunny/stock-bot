# 立花証券セッション管理戦略（日付ベース）

## セッション管理の基本方針

### **1. 日付ベースセッション管理**
- 当日の営業日に紐づくセッションの有無で判定
- 営業日が変わったら新しいセッションが必要
- 同一営業日内はセッション維持

### **2. エラー駆動再認証**
- 証券会社からのセッションエラー時のみ再ログイン
- 自動的なセッション回復
- エラー種別による適切な対応

## セッションファイル構造

### **ディレクトリ構造**
```
./data/sessions/
├── tachibana_session_2024-01-01.json  # 月曜日のセッション
├── tachibana_session_2024-01-02.json  # 火曜日のセッション
├── tachibana_session_2024-01-03.json  # 水曜日のセッション
└── ...
```

### **セッションファイル内容**
```json
{
  "session": {
    "ResultCode": "0",
    "ResultText": "OK",
    "SecondPassword": "****",
    "RequestURL": "https://...",
    "MasterURL": "https://...",
    "PriceURL": "https://...",
    "EventURL": "wss://..."
  },
  "date": "2024-01-01",
  "created_at": "2024-01-01T09:00:00Z",
  "last_used_at": "2024-01-01T15:30:00Z"
}
```

## 運用フロー

### **1. 初回起動（月曜日）**
```bash
# 当日（2024-01-01）のセッションファイルなし
./goa-service.exe
# → 初回ログイン（電話認証）
# → ./data/sessions/tachibana_session_2024-01-01.json 作成
```

### **2. 同日再起動**
```bash
# 当日（2024-01-01）のセッションファイルあり
./goa-service.exe  
# → セッション復元
# → 継続利用
```

### **3. 翌営業日起動**
```bash
# 新しい営業日（2024-01-02）
./goa-service.exe
# → 前日のセッション無効
# → 新しい日付のログイン
# → ./data/sessions/tachibana_session_2024-01-02.json 作成
```

### **4. 週末・祝日の処理**
```bash
# 土曜日起動
./goa-service.exe
# → 前の金曜日（2024-01-05）のセッションを使用
# → 新しいログインは不要
```

## 実装例

### **日付ベースセッション管理**
```go
// getCurrentBusinessDate は現在の営業日を取得
func (dsm *DateBasedSessionManager) getCurrentBusinessDate() string {
    now := time.Now()
    
    // 土日の場合は前の金曜日を返す
    for now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
        now = now.AddDate(0, 0, -1)
    }
    
    return now.Format("2006-01-02")
}

// checkDateChange は日付変更をチェック
func (dsm *DateBasedSessionManager) checkDateChange(ctx context.Context) error {
    currentDate := dsm.getCurrentBusinessDate()
    
    if currentDate == dsm.currentDate {
        return nil // 日付変更なし
    }
    
    // 営業日が変わった場合、新しいセッションが必要
    dsm.logger.Info("Business date changed", 
        "old_date", dsm.currentDate, 
        "new_date", currentDate)
    
    dsm.invalidateSession()
    dsm.currentDate = currentDate
    
    // 新しい日付のセッション復元を試行
    if err := dsm.loadTodaysSession(); err == nil {
        return nil // 既存セッション復元成功
    }
    
    // 新しい日付のログインが必要
    return dsm.performReloginIfNeeded(ctx)
}
```

### **セッションファイル管理**
```go
// getSessionFilePath は指定日付のセッションファイルパスを取得
func (dsm *DateBasedSessionManager) getSessionFilePath(date string) string {
    filename := fmt.Sprintf("tachibana_session_%s.json", date)
    return filepath.Join(dsm.sessionDir, filename)
}

// loadTodaysSession は当日のセッションファイルから復元
func (dsm *DateBasedSessionManager) loadTodaysSession() error {
    filePath := dsm.getSessionFilePath(dsm.currentDate)
    
    data, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("no session file for date %s", dsm.currentDate)
    }
    
    var sessionData SessionData
    if err := json.Unmarshal(data, &sessionData); err != nil {
        return fmt.Errorf("failed to unmarshal session: %w", err)
    }
    
    // 日付の整合性チェック
    if sessionData.Date != dsm.currentDate {
        return fmt.Errorf("session date mismatch: expected %s, got %s", 
            dsm.currentDate, sessionData.Date)
    }
    
    dsm.session = sessionData.Session
    dsm.isAuthenticated = true
    dsm.sessionDate = sessionData.Date
    
    return nil
}
```

## 利点

### **1. 営業日対応**
- 営業日ごとの適切なセッション管理
- 土日・祝日での前営業日セッション利用
- 日付変更時の自動対応

### **2. セッション履歴**
- 過去のセッション情報保持
- デバッグ・監査用途
- 自動クリーンアップ機能

### **3. 運用効率**
- 営業日あたり1回のログイン
- 同日内の完全なセッション維持
- エラー時の自動回復

### **4. 実際の運用例**

#### **月曜日（新しい週）**
```
09:00 - サービス起動
09:01 - 当日セッションなし → ログイン（電話認証）
09:02 - セッション保存（tachibana_session_2024-01-01.json）
09:03 - エージェント開始
```

#### **火曜日（同じ週）**
```
09:00 - サービス起動  
09:01 - 営業日変更検知（2024-01-01 → 2024-01-02）
09:02 - 新しい日付のセッションなし → ログイン
09:03 - セッション保存（tachibana_session_2024-01-02.json）
```

#### **土曜日（週末）**
```
10:00 - サービス起動
10:01 - 営業日判定（土曜日 → 金曜日 2024-01-05）
10:02 - 金曜日のセッション復元
10:03 - ログイン不要で継続利用
```

この実装により、立花証券の営業日ベースのセッション管理要件に完全対応し、最小限のログイン回数で安定した運用が可能になります。