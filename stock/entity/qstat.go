package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type QStat struct {
	entity.BaseEntity  `xorm:"extends"`
	TsCode             string  `xorm:"index notnull" json:"ts_code"`
	SecurityName       string  `json:"security_name"`
	Industry           string  `xorm:"index" json:"industry,omitempty"`  // str 所属行业
	Sector             string  `xorm:"index" json:"sector,omitempty"`    // str 所属细分行业行业
	StartDate          string  `xorm:"index notnull" json:"start_date"`  //统计指标对应的起始季度
	EndDate            string  `json:"end_date"`                         //统计指标对应的结束季度
	TradeDate          int64   `xorm:"index notnull" json:"trade_date"`  //价格对应的日期
	ActualStartDate    string  `json:"actual_start_date"`                //实际的起始季度
	Term               int     `xorm:"index notnull" json:"term"`        //年限
	Source             string  `xorm:"index notnull" json:"source"`      //统计指标
	SourceName         string  `xorm:"index notnull" json:"source_name"` //统计指标对应的字段名
	ReportNumber       int     `json:"report_number"`                    //原始业绩数据的份数
	Pe                 float64 `json:"pe"`
	Peg                float64 `json:"peg"`
	ShareNumber        float64 `json:"share_number"` //通过换手率和成交量算出
	High               float64 `json:"high"`
	Close              float64 `json:"close"`
	MarketValue        float64 `json:"market_value"`
	YearOperateIncome  float64 `json:"year_operate_income"` //营收
	YearNetProfit      float64 `json:"year_net_profit"`
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

func (QStat) TableName() string {
	return "stk_qstat"
}

func (QStat) KeyName() string {
	return entity.FieldName_Id
}

func (QStat) IdName() string {
	return entity.FieldName_Id
}
