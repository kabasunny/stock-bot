// client/tests/master_data_client_impl_test.go
package tests

import (
	"context"
	"stock-bot/internal/infrastructure/client"
	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
	"stock-bot/internal/infrastructure/client/dto/master/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require" // 追加
)

func TestGetMasterDataQuery(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	// リクエストパラメータの設定
	req := request.ReqGetMasterData{
		TargetCLMID:  "CLMIssueMstKabu,CLMOrderErrReason",
		TargetColumn: "sIssueCode,sIssueName,sErrReasonCode,sErrReasonText",
	}

	// API呼び出し
	res, err := c.GetMasterDataQuery(context.Background(), session, req)
	if err != nil {
		t.Fatalf("API呼び出しエラー: %v", err)
	}

	// レスポンスの検証
	assert.NotNil(t, res)
	assert.Equal(t, "CLMMfdsGetMasterData", res.CLMID)

	// StockMasterの検証
	if len(res.StockMaster) > 0 {
		assert.NotEmpty(t, res.StockMaster[0].IssueCode)
		assert.NotEmpty(t, res.StockMaster[0].IssueName)
	}

	// ErrorReasonの検証
	if len(res.ErrorReason) > 0 {
		assert.NotEmpty(t, res.ErrorReason[0].ErrorCode)
		assert.NotEmpty(t, res.ErrorReason[0].ErrorText)
	}
}

func TestGetNewsHeader(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	// リクエストパラメータの設定 (必須パラメータのみ)
	req := request.ReqGetNewsHead{
		Offset: "0",  // レコード取得位置
		Limit:  "10", // レコード取得件数最大
	}

	// API呼び出し
	res, err := c.GetNewsHeader(context.Background(), session, req)
	if err != nil {
		t.Fatalf("API呼び出しエラー: %v", err)
	}

	// レスポンスの検証
	assert.NotNil(t, res)
	assert.Equal(t, "CLMMfdsGetNewsHead", res.CLMID)

	// ニュースヘッダーが存在するかチェック
	if len(res.CLMMfdsNewsHead) > 0 {
		// 最初のニュースヘッダーの検証
		assert.NotEmpty(t, res.CLMMfdsNewsHead[0].PID, "ニュースIDが空でないこと")
		assert.NotEmpty(t, res.CLMMfdsNewsHead[0].PDT, "ニュース日付が空でないこと")
		assert.NotEmpty(t, res.CLMMfdsNewsHead[0].PTM, "ニュース時刻が空でないこと")
		assert.NotEmpty(t, res.CLMMfdsNewsHead[0].PHDL, "ニュースヘッドラインが空でないこと")
	} else {
		t.Log("ニュースヘッダーが存在しません")
	}
}
func TestGetNewsBody(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	// リクエストパラメータの設定 (必須パラメータのみ)
	// TODO: 実際に存在するニュースIDを設定する必要がある
	newsID := "20230315121900_NYU8165"
	req := request.ReqGetNewsBody{
		NewsID: newsID, // 適切なニュースIDを設定
	}

	// API呼び出し
	res, err := c.GetNewsBody(context.Background(), session, req)
	if err != nil {
		t.Fatalf("API呼び出しエラー: %v", err)
	}

	// レスポンスの検証
	assert.NotNil(t, res)
	assert.Equal(t, "CLMMfdsGetNewsBody", res.CLMID)

	// ニュース本文が存在するかチェック
	if len(res.CLMMfdsNewsBody) > 0 {
		// 最初のニュース本文の検証
		assert.Equal(t, newsID, res.CLMMfdsNewsBody[0].PID, "ニュースIDが一致すること")
		assert.NotEmpty(t, res.CLMMfdsNewsBody[0].PDT, "ニュース日付が空でないこと")
		assert.NotEmpty(t, res.CLMMfdsNewsBody[0].PTM, "ニュース時刻が空でないこと")
		assert.NotEmpty(t, res.CLMMfdsNewsBody[0].PHDL, "ニュースヘッドラインが空でないこと")
		assert.NotEmpty(t, res.CLMMfdsNewsBody[0].PTX, "ニュース本文が空でないこと")
	} else {
		t.Log("ニュース本文が存在しません")
	}
}

func TestGetIssueDetail(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	// リクエストパラメータの設定
	targetIssueCode := "6501,7203" // 銘柄コード指定
	req := request.ReqGetIssueDetail{
		TargetIssueCodes: targetIssueCode,
	}

	// API呼び出し
	res, err := c.GetIssueDetail(context.Background(), session, req)
	if err != nil {
		t.Fatalf("API呼び出しエラー: %v", err)
	}

	// レスポンスの検証
	assert.NotNil(t, res)
	assert.Equal(t, "CLMMfdsGetIssueDetail", res.CLMID)

	// 銘柄詳細情報が存在するかチェック
	if len(res.CLMMfdsIssueDetail) > 0 {
		// 最初の銘柄詳細情報の検証
		assert.NotEmpty(t, res.CLMMfdsIssueDetail[0].IssueCode, "銘柄コードが空でないこと")
		// 他のフィールドも必要に応じて検証
	} else {
		t.Log("銘柄詳細情報が存在しません")
	}
}

func TestGetMarginInfo(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	// リクエストパラメータの設定
	targetIssueCode := "6501,7203" // 銘柄コード指定
	req := request.ReqGetMarginInfo{
		TargetIssueCodes: targetIssueCode,
	}

	// API呼び出し
	res, err := c.GetMarginInfo(context.Background(), session, req)
	if err != nil {
		t.Fatalf("API呼び出しエラー: %v", err)
	}

	// レスポンスの検証
	assert.NotNil(t, res)
	assert.Equal(t, "CLMMfdsGetSyoukinZan", res.CLMID)

	// 銘柄詳細情報が存在するかチェック
	if len(res.CLMMfdsSyoukinZan) > 0 {
		// 最初の銘柄詳細情報の検証
		assert.NotEmpty(t, res.CLMMfdsSyoukinZan[0].IssueCode, "銘柄コードが空でないこと")
		// 他のフィールドも必要に応じて検証
	} else {
		t.Log("証金残情報が存在しません")
	}
}

func TestGetCreditInfo(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	// リクエストパラメータの設定
	targetIssueCode := "6501,7203" // 銘柄コード指定
	req := request.ReqGetCreditInfo{
		TargetIssueCodes: targetIssueCode,
	}

	// API呼び出し
	res, err := c.GetCreditInfo(context.Background(), session, req)
	if err != nil {
		t.Fatalf("API呼び出しエラー: %v", err)
	}

	// レスポンスの検証
	assert.NotNil(t, res)
	assert.Equal(t, "CLMMfdsGetShinyouZan", res.CLMID)

	// 銘柄詳細情報が存在するかチェック
	if len(res.CLMMfdsShinyouZan) > 0 {
		// 最初の銘柄詳細情報の検証
		assert.NotEmpty(t, res.CLMMfdsShinyouZan[0].IssueCode, "銘柄コードが空でないこと")
		// 他のフィールドも必要に応じて検証
	} else {
		t.Log("信用残情報が存在しません")
	}
}

func TestGetMarginPremiumInfo(t *testing.T) {
	// テスト用の TachibanaClient を作成
	c := client.CreateTestClient(t)

	// ログイン
	loginReq := request_auth.ReqLogin{
		UserId:   c.GetUserIDForTest(),
		Password: c.GetPasswordForTest(),
	}
	session, err := c.LoginWithPost(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotNil(t, session)

	// リクエストパラメータの設定
	targetIssueCode := "6501,7203" // 銘柄コード指定
	req := request.ReqGetMarginPremiumInfo{
		TargetIssueCodes: targetIssueCode,
	}

	// API呼び出し
	res, err := c.GetMarginPremiumInfo(context.Background(), session, req)
	if err != nil {
		t.Fatalf("API呼び出しエラー: %v", err)
	}

	// レスポンスの検証
	assert.NotNil(t, res)
	assert.Equal(t, "CLMMfdsGetHibuInfo", res.CLMID)

	// 銘柄詳細情報が存在するかチェック
	if len(res.CLMMfdsHibuInfo) > 0 {
		// 最初の銘柄詳細情報の検証
		assert.NotEmpty(t, res.CLMMfdsHibuInfo[0].IssueCode, "銘柄コードが空でないこと")
		// 他のフィールドも必要に応じて検証
		assert.NotEmpty(t, res.CLMMfdsHibuInfo[0].PBWRQ, "逆日歩が空でないこと")
	} else {
		t.Log("逆日歩情報が存在しません")
	}
}

// go test -v ./internal/infrastructure/client/tests/master_data_client_impl_test.go
