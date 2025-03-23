// internal/infrastructure/client/dto/master/response/system_status.go
package response

// ResSystemStatus は、システムステータスの情報を表すDTO
type ResSystemStatus struct {
	CLMID           string `json:"sCLMID"`           // 機能ID (CLMSystemStatus)
	SystemStatusKey string `json:"sSystemStatusKey"` // システム状態KEY (固定値: "001")
	LoginAllowed    string `json:"sLoginKyokaKubun"` // ログイン許可区分 (0：不許可, 1：許可, 2：不許可(サービス時間外), 9：管理者のみ可(テスト中))
	SystemStatus    string `json:"sSystemStatus"`    // システム状態 (0：閉局, 1：開局, 2：一時停止)
	CreateTime      string `json:"sCreateTime"`      // 作成時刻
	UpdateTime      string `json:"sUpdateTime"`      // 更新時刻
	UpdateNumber    string `json:"sUpdateNumber"`    // 更新通番
	DeleteFlag      string `json:"sDeleteFlag"`      // 削除フラグ
	DeleteTime      string `json:"sDeleteTime"`      // 削除時刻
}
