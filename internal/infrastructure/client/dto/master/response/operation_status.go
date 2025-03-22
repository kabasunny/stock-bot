// internal/infrastructure/client/dto/master/response/operation_status.go
package response

// OperationStatus は、運用ステータス別状態の情報を表すDTOです。
type OperationStatus struct {
	CLMID             string `json:"sCLMID"`            // 機能ID (CLMUnyouStatus)
	SystemAccountType string `json:"sSystemKouzaKubun"` // システム口座区分
	OperationCategory string `json:"sUnyouCategory"`    // 運用カテゴリ
	OperationUnit     string `json:"sUnyouUnit"`        // 運用単位
	BusinessDayFlag   string `json:"sEigyouDayC"`       // 営業日区分
	OperationStatus   string `json:"sUnyouStatus"`      // 運用ステータス
	TargetBusiness    string `json:"sTaisyouGyoumu"`    // 対象業務
	BusinessStatus    string `json:"sGyoumuZyoutai"`    // 業務別状態
	CreateTime        string `json:"sCreateTime"`       // 作成時刻
	UpdateTime        string `json:"sUpdateTime"`       // 更新時刻
	UpdateNumber      string `json:"sUpdateNumber"`     // 更新通番
	DeleteFlag        string `json:"sDeleteFlag"`       // 削除フラグ
	DeleteTime        string `json:"sDeleteTime"`       // 削除時刻
	EventName         string `json:"sEventName"`        // イベント名
	EstimatedTime     string `json:"sMeyasuTime"`       // 目安時刻
}
