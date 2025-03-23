// internal/infrastructure/client/dto/master/response/tick_rule.go
package response

// ResTickRule は、呼値の情報を表すDTOです。
type ResTickRule struct {
	CLMID           string `json:"sCLMID"`            // 機能ID (CLMYobine)
	TickUnitNumber  string `json:"sYobineTaniNumber"` // 呼値の単位番号
	ApplicableDate  string `json:"sTekiyouDay"`       // 適用日 (YYYYMMDD)
	BasePrice1      string `json:"sKizunPrice_1"`     // 基準値段1
	BasePrice2      string `json:"sKizunPrice_2"`     // 基準値段2
	BasePrice3      string `json:"sKizunPrice_3"`     // 基準値段3
	BasePrice4      string `json:"sKizunPrice_4"`     // 基準値段4
	BasePrice5      string `json:"sKizunPrice_5"`     // 基準値段5
	BasePrice6      string `json:"sKizunPrice_6"`     // 基準値段6
	BasePrice7      string `json:"sKizunPrice_7"`     // 基準値段7
	BasePrice8      string `json:"sKizunPrice_8"`     // 基準値段8
	BasePrice9      string `json:"sKizunPrice_9"`     // 基準値段9
	BasePrice10     string `json:"sKizunPrice_10"`    // 基準値段10
	BasePrice11     string `json:"sKizunPrice_11"`    // 基準値段11
	BasePrice12     string `json:"sKizunPrice_12"`    // 基準値段12
	BasePrice13     string `json:"sKizunPrice_13"`    // 基準値段13
	BasePrice14     string `json:"sKizunPrice_14"`    // 基準値段14
	BasePrice15     string `json:"sKizunPrice_15"`    // 基準値段15
	BasePrice16     string `json:"sKizunPrice_16"`    // 基準値段16
	BasePrice17     string `json:"sKizunPrice_17"`    // 基準値段17
	BasePrice18     string `json:"sKizunPrice_18"`    // 基準値段18
	BasePrice19     string `json:"sKizunPrice_19"`    // 基準値段19
	BasePrice20     string `json:"sKizunPrice_20"`    // 基準値段20
	TickValue1      string `json:"sYobineTanka_1"`    // 呼値単価1
	TickValue2      string `json:"sYobineTanka_2"`    // 呼値単価2
	TickValue3      string `json:"sYobineTanka_3"`    // 呼値単価3
	TickValue4      string `json:"sYobineTanka_4"`    // 呼値単価4
	TickValue5      string `json:"sYobineTanka_5"`    // 呼値単価5
	TickValue6      string `json:"sYobineTanka_6"`    // 呼値単価6
	TickValue7      string `json:"sYobineTanka_7"`    // 呼値単価7
	TickValue8      string `json:"sYobineTanka_8"`    // 呼値単価8
	TickValue9      string `json:"sYobineTanka_9"`    // 呼値単価9
	TickValue10     string `json:"sYobineTanka_10"`   // 呼値単価10
	TickValue11     string `json:"sYobineTanka_11"`   // 呼値単価11
	TickValue12     string `json:"sYobineTanka_12"`   // 呼値単価12
	TickValue13     string `json:"sYobineTanka_13"`   // 呼値単価13
	TickValue14     string `json:"sYobineTanka_14"`   // 呼値単価14
	TickValue15     string `json:"sYobineTanka_15"`   // 呼値単価15
	TickValue16     string `json:"sYobineTanka_16"`   // 呼値単価16
	TickValue17     string `json:"sYobineTanka_17"`   // 呼値単価17
	TickValue18     string `json:"sYobineTanka_18"`   // 呼値単価18
	TickValue19     string `json:"sYobineTanka_19"`   // 呼値単価19
	TickValue20     string `json:"sYobineTanka_20"`   // 呼値単価20
	DecimalPlaces1  string `json:"sDecimal_1"`        // 小数点桁数1
	DecimalPlaces2  string `json:"sDecimal_2"`        // 小数点桁数2
	DecimalPlaces3  string `json:"sDecimal_3"`        // 小数点桁数3
	DecimalPlaces4  string `json:"sDecimal_4"`        // 小数点桁数4
	DecimalPlaces5  string `json:"sDecimal_5"`        // 小数点桁数5
	DecimalPlaces6  string `json:"sDecimal_6"`        // 小数点桁数6
	DecimalPlaces7  string `json:"sDecimal_7"`        // 小数点桁数7
	DecimalPlaces8  string `json:"sDecimal_8"`        // 小数点桁数8
	DecimalPlaces9  string `json:"sDecimal_9"`        // 小数点桁数9
	DecimalPlaces10 string `json:"sDecimal_10"`       // 小数点桁数10
	DecimalPlaces11 string `json:"sDecimal_11"`       // 小数点桁数11
	DecimalPlaces12 string `json:"sDecimal_12"`       // 小数点桁数12
	DecimalPlaces13 string `json:"sDecimal_13"`       // 小数点桁数13
	DecimalPlaces14 string `json:"sDecimal_14"`       // 小数点桁数14
	DecimalPlaces15 string `json:"sDecimal_15"`       // 小数点桁数15
	DecimalPlaces16 string `json:"sDecimal_16"`       // 小数点桁数16
	DecimalPlaces17 string `json:"sDecimal_17"`       // 小数点桁数17
	DecimalPlaces18 string `json:"sDecimal_18"`       // 小数点桁数18
	DecimalPlaces19 string `json:"sDecimal_19"`       // 小数点桁数19
	DecimalPlaces20 string `json:"sDecimal_20"`       // 小数点桁数20
	CreateDate      string `json:"sCreateDate"`       // 作成日時
	UpdateDate      string `json:"sUpdateDate"`       // 更新日時
}
