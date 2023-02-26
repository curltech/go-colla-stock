package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Performance struct {
	entity.BaseEntity  `xorm:"extends"`
	SecurityCode       string  `xorm:"index notnull" json:"security_code"`
	SecurityNameAbbr   string  `json:"security_name_abbr"`
	TradeMarketCode    string  `json:"trade_market_code"`
	TradeMarket        string  `json:"trade_market"`
	SecurityTypeCode   string  `json:"security_type_code"`
	SecurityType       string  `json:"security_type"`
	NewestDate         string  `json:"newest_date"`
	ReportDate         string  `xorm:"index notnull" json:"report_date"`
	BasicEps           float64 `json:"basic_eps"`            //每股收益
	DeductBasicEps     float64 `json:"deduct_basic_eps"`     //每股扣非收益
	TotalOperateIncome float64 `json:"total_operate_income"` //营收
	ParentNetProfit    float64 `json:"parent_net_profit"`    //归母净利润
	WeightAvgRoe       float64 `json:"weight_avg_roe"`       //净资产收益率
	YoySales           float64 `json:"yoy_sales"`            //应收同比增长
	YoyDeduNp          float64 `json:"yoy_dedu_np"`          //扣非净利润同比增长
	Bps                float64 `json:"bps"`                  //每股净资产
	Cfps               float64 `json:"cfps"`                 //每股经营现金流量(元)
	GrossProfitMargin  float64 `json:"gross_profit_margin"`  //销售毛利率(%)
	OrLastMonth        float64 `json:"or_last_month"`        //营业收入季度环比增长(%)
	NpLastMonth        float64 `json:"np_last_month"`        //净利润季度环比增长(%)
	AssignDscrpt       string  `json:"assign_dscrpt"`
	PayYear            string  `json:"pay_year"`
	PublishName        string  `json:"publish_name"`
	DividendYieldRatio float64 `json:"dividend_yield_ratio"` //股息率
	NoticeDate         string  `json:"notice_date"`
	OrgCode            string  `json:"org_code"`
	TradeMarketZJG     string  `json:"trade_market_zjg"`
	IsNew              string  `json:"is_new"`
	QDate              string  `xorm:"index notnull" json:"qdate"`
	NDate              string  `json:"ndate"`
	DataType           string  `json:"data_type"`
	DataYear           string  `json:"data_year"`
	DateMmDd           string  `json:"date_mm_dd"`
	EITime             string  `json:"ei_time"`
	SecuCode           string  `json:"secu_code"`
}

func (Performance) TableName() string {
	return "stk_performance"
}

func (Performance) KeyName() string {
	return entity.FieldName_Id
}

func (Performance) IdName() string {
	return entity.FieldName_Id
}
