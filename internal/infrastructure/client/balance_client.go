// internal/infrastructure/client/balance_client.go
package client

import (
	"context"
	"stock-bot/internal/infrastructure/client/dto/balance/request"
	"stock-bot/internal/infrastructure/client/dto/balance/response"
)

// BalanceClient は、残高・余力関連の API を扱うインターフェース
type BalanceClient interface {
	// GetGenbutuKabuList は、現物保有銘柄の一覧を取得
	GetGenbutuKabuList(ctx context.Context) (*response.ResGenbutuKabuList, error)
	// GetShinyouTategyokuList は、信用建玉の一覧を取得
	GetShinyouTategyokuList(ctx context.Context) (*response.ResShinyouTategyokuList, error)
	// GetZanKaiKanougaku は、現物買付可能額、信用新規建可能額などを取得
	GetZanKaiKanougaku(ctx context.Context, req request.ReqZanKaiKanougaku) (*response.ResZanKaiKanougaku, error)
	// GetZanKaiKanougakuSuii は、現物や信用の可能額、委託保証金率等の推移を過去6営業日に遡って取得
	GetZanKaiKanougakuSuii(ctx context.Context, req request.ReqZanKaiKanougakuSuii) (*response.ResZanKaiKanougakuSuii, error)
	// GetZanKaiSummary 可能額サマリーを取得
	GetZanKaiSummary(ctx context.Context) (*response.ResZanKaiSummary, error)
	// GetZanKaiGenbutuKaitukeSyousai は、指定営業日の現物株式買付可能額詳細を取得
	GetZanKaiGenbutuKaitukeSyousai(ctx context.Context, tradingDay int) (*response.ResZanKaiGenbutuKaitukeSyousai, error)
	// GetZanKaiSinyouSinkidateSyousai は、指定営業日の信用新規建て可能額詳細を取得
	GetZanKaiSinyouSinkidateSyousai(ctx context.Context, tradingDay int) (*response.ResZanKaiSinyouSinkidateSyousai, error)
	// GetZanRealHosyoukinRitu は、リアルタイムの委託保証金率等を取得
	GetZanRealHosyoukinRitu(ctx context.Context, req request.ReqZanRealHosyoukinRitu) (*response.ResZanRealHosyoukinRitu, error)
	// GetZanShinkiKanoIjiritu は、信用新規建て可能維持率を取得
	GetZanShinkiKanoIjiritu(ctx context.Context, req request.ReqZanShinkiKanoIjiritu) (*response.ResZanShinkiKanoIjiritu, error)
	// GetZanUriKanousuu は、指定銘柄の売却可能数量を取得
	GetZanUriKanousuu(ctx context.Context, req request.ReqZanUriKanousuu) (*response.ResZanUriKanousuu, error)
}
