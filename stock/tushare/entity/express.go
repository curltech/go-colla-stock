package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Express struct {
	entity.BaseEntity     `xorm:"extends"`
	TsCode                string  `xorm:"varchar(255)" json:"ts_code,omitempty"` // str	Y	TS代码
	AnnDate               string  `json:"ann_date,omitempty"`                    //	str	公告日期
	EndDate               string  `json:"end_date,omitempty"`                    //	str	报告期
	Revenue               float64 `json:"revenue,omitempty"`                     //	float	营业收入(元)
	OperateProfit         float64 `json:"operate_profit,omitempty"`              //	float	营业利润(元)
	TotalProfit           float64 `json:"total_profit,omitempty"`                //	float	利润总额(元)
	NIncome               float64 `json:"n_income,omitempty"`                    //	float	净利润(元)
	TotalAssets           float64 `json:"total_assets,omitempty"`                //	float	总资产(元)
	TotalHldrEqyExcMinInt float64 `json:"total_hldr_eqy_exc_min_int,omitempty"`  //	float	股东权益合计(不含少数股东权益)(元)
	DilutedEps            float64 `json:"diluted_eps,omitempty"`                 //	float	每股收益(摊薄)(元)
	DilutedRoe            float64 `json:"diluted_roe,omitempty"`                 //	float	净资产收益率(摊薄)(%)
	YoyNetProfit          float64 `json:"yoy_net_profit,omitempty"`              //	float	去年同期修正后净利润
	Bps                   float64 `json:"bps,omitempty"`                         //	float	每股净资产
	YoySales              float64 `json:"yoy_sales,omitempty"`                   //	float	同比增长率:营业收入
	YoyOp                 float64 `json:"yoy_op,omitempty"`                      //	float	同比增长率:营业利润
	YoyTp                 float64 `json:"yoy_tp,omitempty"`                      //	float	同比增长率:利润总额
	YoyDeduNp             float64 `json:"yoy_dedu_np,omitempty"`                 //	float	同比增长率:归属母公司股东的净利润
	YoyEps                float64 `json:"yoy_eps,omitempty"`                     //	float	同比增长率:基本每股收益
	YoyRoe                float64 `json:"yoy_roe,omitempty"`                     //	float	同比增减:加权平均净资产收益率
	GrowthAssets          float64 `json:"growth_assets,omitempty"`               //	float	比年初增长率:总资产
	YoyEquity             float64 `json:"yoy_equity,omitempty"`                  //	float	比年初增长率:归属母公司的股东权益
	GrowthBps             float64 `json:"growth_bps,omitempty"`                  //	float	比年初增长率:归属于母公司股东的每股净资产
	OrLastYear            float64 `json:"or_last_year,omitempty"`                //	float	去年同期营业收入
	OpLastYear            float64 `json:"op_last_year,omitempty"`                //	float	去年同期营业利润
	TpLastYear            float64 `json:"tp_last_year,omitempty"`                //	float	去年同期利润总额
	NpLastYear            float64 `json:"np_last_year,omitempty"`                //	float	去年同期净利润
	EpsLastYear           float64 `json:"eps_last_year,omitempty"`               //	float	去年同期每股收益
	OpenNetAssets         float64 `json:"open_net_assets,omitempty"`             //	float	期初净资产
	OpenBps               float64 `json:"open_bps,omitempty"`                    //	float	期初每股净资产
	PerfSummary           string  `json:"perf_summary,omitempty"`                //	str	业绩简要说明
	IsAudit               int64   `json:"is_audit,omitempty"`                    //	int	是否审计： 1是 0否
	Remark                string  `json:"remark,omitempty"`                      //	str	备注
}

func (Express) TableName() string {
	return "stk_forecast"
}

func (Express) KeyName() string {
	return entity.FieldName_Id
}

func (Express) IdName() string {
	return entity.FieldName_Id
}
