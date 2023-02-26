package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

/**

 */
type StatScore struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string `xorm:"index notnull" json:"ts_code"`
	SecurityName      string `json:"security_name"`
	StartDate         string `xorm:"index notnull" json:"start_date"` //统计指标对应的起始季度
	EndDate           string `json:"end_date"`                        //统计指标对应的结束季度
	Term              int    `xorm:"index notnull" json:"term"`       //年限
	ReportNumber      int    `json:"report_number"`                   //原始业绩报告数目
	TradeDate         int64  `xorm:"index notnull" json:"trade_date"` //价格对应的日期
	//评分指标类型：贵贱，稳定性，增速，累积，相关性，风险，景气程度7大类指标
	Industry                 string  `json:"industry"`
	Sector                   string  `xorm:"index" json:"sector,omitempty"` // str 所属细分行业行业
	Area                     string  `json:"area"`
	Market                   string  `json:"market"`
	ListDate                 int64   `json:"list_date"`
	ListStatus               string  `json:"list_status"`
	RiskScore                float64 `xorm:"index notnull" json:"risk_score"`
	RsdOrLastMonth           float64 `json:"rsd_or_last_month"`
	RsdNpLastMonth           float64 `json:"rsd_np_last_month"`
	RsdPctChgMarketValue     float64 `json:"rsd_pct_chg_market_value"` //*稳定性（rsd）和增速指标（mean，accumulate）
	RsdYoySales              float64 `json:"rsd_yoy_sales"`            //*稳定性和增速指标（mean，accumulate），越高分越高
	RsdYoyDeduNp             float64 `json:"rsd_yoy_dedu_np"`          //*稳定性和增速指标（mean，accumulate）
	RsdPe                    float64 `json:"rsd_pe"`                   //*贵贱指标，越高分越低
	RsdWeightAvgRoe          float64 `json:"rsd_weight_avg_roe"`       //*稳定性和增速指标
	RsdGrossprofitMargin     float64 `json:"rsd_gross_profit_margin"`  //*稳定性和增速指标
	StableScore              float64 `xorm:"index notnull" json:"stable_score"`
	MeanPctChgMarketValue    float64 `json:"mean_pct_chg_market_value"` //*稳定性（rsd）和增速指标（mean，accumulate）
	MeanYoySales             float64 `json:"mean_yoy_sales"`            //*稳定性和增速指标（mean，accumulate），越高分越高
	MeanYoyDeduNp            float64 `json:"mean_yoy_dedu_np"`          //*稳定性和增速指标（mean，accumulate）
	MeanOrLastMonth          float64 `json:"mean_or_last_month"`
	MeanNpLastMonth          float64 `json:"mean_np_last_month"`
	MeanPe                   float64 `json:"mean_pe"`                  //*贵贱指标，越高分越低
	MeanWeightAvgRoe         float64 `json:"mean_weight_avg_roe"`      //*稳定性和增速指标
	MeanGrossprofitMargin    float64 `json:"mean_gross_profit_margin"` //*稳定性和增速指标
	MedianPctChgMarketValue  float64 `json:"median_pct_chg_market_value"`
	MedianYoySales           float64 `json:"median_yoy_sales"`
	MedianYoyDeduNp          float64 `json:"median_yoy_dedu_np"`
	MedianOrLastMonth        float64 `json:"median_or_last_month"`
	MedianNpLastMonth        float64 `json:"median_np_last_month"`
	MedianWeightAvgRoe       float64 `json:"median_weight_avg_roe"`
	MedianGrossprofitMargin  float64 `json:"median_gross_profit_margin"`
	IncreaseScore            float64 `xorm:"index notnull" json:"increase_score"`
	AccPctChgMarketValue     float64 `json:"acc_pct_chg_market_value"` //*稳定性（rsd）和增速指标（mean，accumulate）
	AccYoySales              float64 `json:"acc_yoy_sales"`            //*稳定性和增速指标（mean，accumulate），越高分越高
	AccYoyDeduNp             float64 `json:"acc_yoy_dedu_np"`          //*稳定性和增速指标（mean，accumulate）
	AccScore                 float64 `xorm:"index notnull" json:"acc_score"`
	MedianPe                 float64 `json:"median_pe"` //*贵贱指标，越高分越低
	MeanPeg                  float64 `json:"mean_peg"`
	MedianPeg                float64 `json:"median_peg"`
	PriceScore               float64 `xorm:"index notnull" json:"price_score"`
	CorrYoySales             float64 `json:"corr_yoy_sales"`           //*相关性指标，越高分越低
	CorrYoyDeduNp            float64 `json:"corr_yoy_dedu_np"`         //*相关性指标，越高分越低
	CorrYearNetProfit        float64 `json:"corr_year_net_profit"`     //*相关性指标
	CorrYearOperateIncome    float64 `json:"corr_year_operate_income"` //*相关性指标
	CorrWeightAvgRoe         float64 `json:"corr_weight_avg_roe"`
	CorrGrossprofitMargin    float64 `json:"corr_gross_profit_margin"`
	CorrScore                float64 `xorm:"index notnull" json:"corr_score"`
	LastPctChgMarketValue    float64 `json:"last_pct_chg_market_value"` //*稳定性（rsd）和增速指标（mean，accumulate）
	LastYoySales             float64 `json:"last_yoy_sales"`            //*稳定性和增速指标（mean，accumulate），越高分越高
	LastYoyDeduNp            float64 `json:"last_yoy_dedu_np"`
	LastOrLastMonth          float64 `json:"last_or_last_month"`
	LastNpLastMonth          float64 `json:"last_np_last_month"`
	LastMeanPe               float64 `json:"last_mean_pe"`
	LastMeanPeg              float64 `json:"last_mean_peg"`
	ProsScore                float64 `xorm:"index notnull" json:"pros_score"`
	TrendScore               float64 `xorm:"index notnull" json:"trend_score"`
	OperationScore           float64 `xorm:"index notnull" json:"operation_score"`
	TotalScore               float64 `xorm:"index notnull" json:"total_score"`
	BadTip                   string  `xorm:"varchar(32000)" json:"bad_tip"`
	GoodTip                  string  `xorm:"varchar(32000)" json:"good_tip"`
	PercentileRiskScore      float64 `xorm:"index" json:"percentile_risk_score"`
	PercentileStableScore    float64 `xorm:"index" json:"percentile_stable_score"`
	PercentileIncreaseScore  float64 `xorm:"index" json:"percentile_increase_score"`
	PercentileAccScore       float64 `xorm:"index" json:"percentile_acc_score"`
	PercentilePriceScore     float64 `xorm:"index" json:"percentile_price_score"`
	PercentileCorrScore      float64 `xorm:"index" json:"percentile_corr_score"`
	PercentileProsScore      float64 `xorm:"index" json:"percentile_pros_score"`
	PercentileTrendScore     float64 `xorm:"index" json:"percentile_trend_score"`
	PercentileOperationScore float64 `xorm:"index" json:"percentile_operation_score"`
	PercentileTotalScore     float64 `xorm:"index" json:"percentile_total_score"`
}

func (StatScore) TableName() string {
	return "stk_statscore"
}

func (StatScore) KeyName() string {
	return entity.FieldName_Id
}

func (StatScore) IdName() string {
	return entity.FieldName_Id
}
