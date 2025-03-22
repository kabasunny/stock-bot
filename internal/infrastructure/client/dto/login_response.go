// infrastructure/client/dto/login_response.go
package dto

type LoginResponse struct {
	ResultCode                  string `json:"sResultCode"`
	ResultText                  string `json:"sResultText"`
	ZyoutoekiKazeiC             string `json:"sZyoutoekiKazeiC"`
	SecondPasswordOmit          string `json:"sSecondPasswordOmit"`
	LastLoginDate               string `json:"sLastLoginDate"`
	SogoKouzaKubun              string `json:"sSogoKouzaKubun"`
	HogoAdukariKouzaKubun       string `json:"sHogoAdukariKouzaKubun"`
	FurikaeKouzaKubun           string `json:"sFurikaeKouzaKubun"`
	GaikokuKouzaKubun           string `json:"sGaikokuKouzaKubun"`
	MRFKouzaKubun               string `json:"sMRFKouzaKubun"`
	TokuteiKouzaKubunGenbutu    string `json:"sTokuteiKouzaKubunGenbutu"`
	TokuteiKouzaKubunSinyou     string `json:"sTokuteiKouzaKubunSinyou"`
	TokuteiKouzaKubunTousin     string `json:"sTokuteiKouzaKubunTousin"`
	TokuteiHaitouKouzaKubun     string `json:"sTokuteiHaitouKouzaKubun"`
	TokuteiKanriKouzaKubun      string `json:"sTokuteiKanriKouzaKubun"`
	SinyouKouzaKubun            string `json:"sSinyouKouzaKubun"`
	SakopKouzaKubun             string `json:"sSakopKouzaKubun"`
	MMFKouzaKubun               string `json:"sMMFKouzaKubun"`
	TyukokufKouzaKubun          string `json:"sTyukokufKouzaKubun"`
	KawaseKouzaKubun            string `json:"sKawaseKouzaKubun"`
	HikazeiKouzaKubun           string `json:"sHikazeiKouzaKubun"`
	KinsyouhouMidokuFlg         string `json:"sKinsyouhouMidokuFlg"`
	RequestURL                  string `json:"sUrlRequest"`
	MasterURL                   string `json:"sUrlMaster"`
	PriceURL                    string `json:"sUrlPrice"`
	EventURL                    string `json:"sUrlEvent"`
	UpdateInformWebDocument     string `json:"sUpdateInformWebDocument"`
	UpdateInformAPISpecFunction string `json:"sUpdateInformAPISpecFunction"`
}
