// internal/infrastructure/client/dto/master/request/get_master_data.go
package request

import (
	"stock-bot/internal/infrastructure/client/dto"
)

// GetMasterDataRequest は、マスタ情報問合取得のリクエストを表すDTO
type GetMasterDataRequest struct {
	dto.RequestBase        // 共通フィールド
	CLMID           string `json:"sCLMID"`        // 機能ID (固定値: "CLMMfdsGetMasterData")
	TargetCLMID     string `json:"sTargetCLMID"`  // 対象機能ID (カンマ区切りで複数指定可能)
	TargetColumn    string `json:"sTargetColumn"` // 対象項目 (カンマ区切りで複数指定可能)
}
