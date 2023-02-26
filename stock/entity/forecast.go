package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Forecast struct {
	entity.BaseEntity   `xorm:"extends"`
	SecuCode            string  `json:"secu_code"`
	SecurityCode        string  `xorm:"index notnull" json:"security_code"`
	SecurityNameAbbr    string  `json:"security_name_abbr"`
	TradeMarketCode     string  `json:"trade_market_code"`
	TradeMarket         string  `json:"trade_market"`
	SecurityTypeCode    string  `json:"security_type_code"`
	SecurityType        string  `json:"security_type"`
	NoticeDate          string  `json:"notice_date"`
	OrgCode             string  `json:"org_code"`
	ReportDate          string  `xorm:"index notnull" json:"report_date"`
	QDate               string  `xorm:"index notnull" json:"qdate"`
	NDate               string  `xorm:"index notnull" json:"ndate"`
	PredictFinanceCode  string  `json:"predict_finance_code"`
	PredictFinance      string  `xorm:"varchar(1024)" json:"predict_finance"`
	PredictAmtLower     float64 `json:"predict_amt_lower"`
	PredictAmtUpper     float64 `json:"predict_amt_upper"`
	AddAmpLower         float64 `json:"add_amp_lower"`
	AddAmpUpper         float64 `json:"add_amp_upper"`
	PredictContent      string  `xorm:"varchar(32000)" json:"predict_content"`
	ChangeReasonExplain string  `xorm:"varchar(32000)" json:"change_reason_explain"`
	PredictType         string  `json:"predict_type"`
	PreYearSamePeriod   float64 `json:"pre_year_same_period"`
	IncreaseAvg         float64 `json:"increase_avg"`
	ForecastAvg         float64 `json:"forecast_avg"`
	ForecastState       string  `json:"forecast_state"`
	IsLatest            string  `json:"is_latest"`
}

func (Forecast) TableName() string {
	return "stk_forecast"
}

func (Forecast) KeyName() string {
	return entity.FieldName_Id
}

func (Forecast) IdName() string {
	return entity.FieldName_Id
}
