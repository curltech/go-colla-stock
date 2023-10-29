package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type QPerformance struct {
	entity.BaseEntity  `xorm:"extends"`
	TsCode             string  `xorm:"index notnull" json:"ts_code"`
	SecurityName       string  `json:"security_name"`
	Industry           string  `xorm:"index" json:"industry,omitempty"` // str 所属行业
	Sector             string  `xorm:"index" json:"sector,omitempty"`   // str 所属细分行业行业
	QDate              string  `xorm:"index notnull" json:"qdate"`      //业绩对应的季度
	NDate              string  `json:"ndate"`                           //业绩发布时间对应的季度
	TradeDate          int64   `xorm:"index notnull" json:"trade_date"` //价格对应的日期
	Source             string  `xorm:"index notnull" json:"source"`     //业绩来源：预测，快报，业绩报告
	LineType           int64   `xorm:"index" json:"line_type"`          //价格来源：wmqy,day
	Pe                 float64 `json:"pe"`
	Peg                float64 `json:"peg"`
	ShareNumber        float64 `json:"share_number"` //通过换手率和成交量算出
	High               float64 `json:"high"`
	Close              float64 `json:"close"`
	MarketValue        float64 `json:"market_value"`
	YearNetProfit      float64 `json:"year_net_profit"`
	YearOperateIncome  float64 `json:"year_operate_income"`  //营收
	TotalOperateIncome float64 `json:"total_operate_income"` //营收
	PctChgHigh         float64 `json:"pct_chg_high"`
	PctChgClose        float64 `json:"pct_chg_close"`
	PctChgMarketValue  float64 `json:"pct_chg_market_value"`
	WeightAvgRoe       float64 `json:"weight_avg_roe"`
	GrossProfitMargin  float64 `json:"gross_profit_margin"`
	ParentNetProfit    float64 `json:"parent_net_profit"`
	BasicEps           float64 `json:"basic_eps"`
	OrLastMonth        float64 `json:"or_last_month"`
	NpLastMonth        float64 `json:"np_last_month"`
	YoySales           float64 `json:"yoy_sales"`
	YoyDeduNp          float64 `json:"yoy_dedu_np"`
	Cfps               float64 `json:"cfps"`
	DividendYieldRatio float64 `json:"dividend_yield_ratio"`
}

func (QPerformance) TableName() string {
	return "stk_qperformance"
}

func (QPerformance) KeyName() string {
	return entity.FieldName_Id
}

func (QPerformance) IdName() string {
	return entity.FieldName_Id
}
