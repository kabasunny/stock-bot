// internal/infrastructure/client/master_data_client.go
package client

import (
	"context"
	"stock-bot/internal/infrastructure/client/dto/master/request"
	"stock-bot/internal/infrastructure/client/dto/master/response"
)

// MasterDataClient は、マスタデータ関連の API を扱うインターフェース
type MasterDataClient interface {
	// DownloadMasterData は、各種マスタ情報をリアルタイム配信でダウンロード
	DownloadMasterData(ctx context.Context, session *Session, req request.ReqDownloadMaster) (*response.ResDownloadMaster, error)
	// GetMasterDataQuery は、指定したマスタ情報を取得（複数指定、項目指定可能）
	GetMasterDataQuery(ctx context.Context, session *Session, req request.ReqGetMasterData) (*response.ResGetMasterData, error)
	// GetNewsHeader は、指定した条件のニュースヘッダーを取得
	GetNewsHeader(ctx context.Context, session *Session, req request.ReqGetNewsHead) (*response.ResGetNewsHeader, error)
	// GetNewsBody は、指定したニュースIDのニュース本文を取得
	GetNewsBody(ctx context.Context, session *Session, req request.ReqGetNewsBody) (*response.ResGetNewsBody, error)
	// GetIssueDetail は、BPS, EPS, 配当等の情報を取得
	GetIssueDetail(ctx context.Context, session *Session, req request.ReqGetIssueDetail) (*response.ResGetIssueDetail, error)
	// GetMarginInfo は、証金残情報を取得
	GetMarginInfo(ctx context.Context, session *Session, req request.ReqGetMarginInfo) (*response.ResGetMarginInfo, error)
	// GetCreditInfo は、信用残情報を取得
GetCreditInfo(ctx context.Context, session *Session, req request.ReqGetCreditInfo) (*response.ResGetCreditInfo, error)
	// GetMarginPremiumInfo は、逆日歩情報を取得
	GetMarginPremiumInfo(ctx context.Context, session *Session, req request.ReqGetMarginPremiumInfo) (*response.ResGetMarginPremiumInfo, error)
}
