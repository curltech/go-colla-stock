package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type T struct {
}

type Express struct {
	entity.BaseEntity    `xorm:"extends"`
	SecurityCode         string  `xorm:"index notnull" json:"security_code"`
	SecurityNameAbbr     string  `json:"security_name_abbr"`
	TradeMarketCode      string  `json:"trade_market_code"`
	TradeMarket          string  `json:"trade_market"`
	SecurityTypeCode     string  `json:"security_type_code"`
	SecurityType         string  `json:"security_type"`
	NewestDate           string  `json:"newest_date"`
	ReportDate           string  `xorm:"index notnull" json:"report_date"`
	BasicEps             float64 `json:"basic_eps"`               //每股收益
	TotalOperateIncome   float64 `json:"total_operate_income"`    //营收
	TotalOperateIncomeSq float64 `json:"total_operate_income_sq"` //去年同期(元)
	ParentNetProfit      float64 `json:"parent_net_profit"`       //归母净利润
	ParentNetProfitSq    float64 `json:"parent_net_profit_sq"`    //去年同期(元)
	ParentBvps           float64 `json:"parent_bvps"`             //每股净资产
	WeightAvgRoe         float64 `json:"weight_avg_roe"`          //净资产收益率
	YoySales             float64 `json:"yoy_sales"`               //收入同比增长
	YoyNetProfit         float64 `json:"yoy_net_profit"`          //净利润同比增长
	OrLastMonth          float64 `json:"or_last_month"`           //收入季度环比增长
	NpLastMonth          float64 `json:"np_last_month"`           //利润季度环比增长
	PublishName          string  `json:"publish_name"`
	NoticeDate           string  `json:"notice_date"`
	OrgCode              string  `json:"org_code"`
	Market               string  `json:"market"`
	IsNew                string  `json:"is_new"`
	QDate                string  `xorm:"index notnull" json:"qdate"`
	NDate                string  `xorm:"index notnull" json:"ndate"`
	DataType             string  `json:"data_type"`
	EITime               string  `json:"ei_time"`
	SecuCode             string  `json:"secu_code"`
}

func (Express) TableName() string {
	return "stk_express"
}

func (Express) KeyName() string {
	return entity.FieldName_Id
}

func (Express) IdName() string {
	return entity.FieldName_Id
}
