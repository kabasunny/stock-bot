// internal/infrastructure/client/dto/master/response/date_info.go
package response

// DateInfo は、日付情報を表すDTO
type DateInfo struct {
	CLMID                   string `json:"sCLMID"`                // 機能ID (CLMDateZyouhou)
	DayKey                  string `json:"sDayKey"`               // 日付KEY (001：当日基準, 002：翌日基準（夕場）)
	PreviousBusinessDay1    string `json:"sMaeEigyouDay_1"`       // 1営業日前 (YYYYMMDD)
	PreviousBusinessDay2    string `json:"sMaeEigyouDay_2"`       // 2営業日前 (YYYYMMDD)
	PreviousBusinessDay3    string `json:"sMaeEigyouDay_3"`       // 3営業日前 (YYYYMMDD)
	CurrentDay              string `json:"sTheDay"`               // 当日日付 (YYYYMMDD)
	NextBusinessDay1        string `json:"sYokuEigyouDay_1"`      // 翌1営業日 (YYYYMMDD)
	NextBusinessDay2        string `json:"sYokuEigyouDay_2"`      // 翌2営業日 (YYYYMMDD)
	NextBusinessDay3        string `json:"sYokuEigyouDay_3"`      // 翌3営業日 (YYYYMMDD)
	NextBusinessDay4        string `json:"sYokuEigyouDay_4"`      // 翌4営業日 (YYYYMMDD)
	NextBusinessDay5        string `json:"sYokuEigyouDay_5"`      // 翌5営業日 (YYYYMMDD)
	NextBusinessDay6        string `json:"sYokuEigyouDay_6"`      // 翌6営業日 (YYYYMMDD)
	NextBusinessDay7        string `json:"sYokuEigyouDay_7"`      // 翌7営業日 (YYYYMMDD)
	NextBusinessDay8        string `json:"sYokuEigyouDay_8"`      // 翌8営業日 (YYYYMMDD)
	NextBusinessDay9        string `json:"sYokuEigyouDay_9"`      // 翌9営業日 (YYYYMMDD)
	NextBusinessDay10       string `json:"sYokuEigyouDay_10"`     // 翌10営業日 (YYYYMMDD)
	StockSettlementDate     string `json:"sKabuUkewatasiDay"`     // 株式受渡日 (YYYYMMDD)
	StockTempSettlementDate string `json:"sKabuKariUkewatasiDay"` // 株式仮決受渡日 (YYYYMMDD)
	BondSettlementDate      string `json:"sBondUkewatasiDay"`     // 債券受渡日 (YYYYMMDD)
}
