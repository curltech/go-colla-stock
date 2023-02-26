package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

/**
文件名即股票代码
每32个字节为一天数据
每4个字节为一个字段，每个字段内低字节在前
00 ~ 03 字节：年月日, 整型
04 ~ 07 字节：开盘价*1000， 整型
08 ~ 11 字节：最高价*1000,  整型
12 ~ 15 字节：最低价*1000,  整型
16 ~ 19 字节：收盘价*1000,  整型
20 ~ 23 字节：成交额（元），float型
24 ~ 27 字节：成交量（手），整型
28 ~ 31 字节：上日收盘*1000, 整型
*/

type StockLine struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string  `xorm:"varchar(255) index notnull" json:"ts_code"` // str	股票代码
	TradeDate         int64   `xorm:"index notnull" json:"trade_date"`
	ShareNumber       float64 `json:"share_number"`
	Open              float64 `json:"open"`   // float	(1/5/15/30/60分钟) D日 W周 M月 开盘价
	High              float64 `json:"high"`   // float	(1/5/15/30/60分钟) D日 W周 M月 最高价
	Low               float64 `json:"low"`    // float	(1/5/15/30/60分钟) D日 W周 M月 最低价
	Close             float64 `json:"close"`  // float	(1/5/15/30/60分钟) D日 W周 M月 收盘价
	Vol               float64 `json:"vol"`    // float	(1/5/15/30/60分钟) D日 W周 M月 成交量 (手)
	Amount            float64 `json:"amount"` // float	(1/5/15/30/60分钟) D日 W周 M月 成交额 (千元)
	Turnover          float64 `json:"turnover"`
	PreClose          float64 `json:"pre_close"` // float	上一(1/5/15/30/60分钟) D日 W周 M月 收盘价
	MainNetInflow     float64 `json:"main_net_inflow"`
	SmallNetInflow    float64 `json:"small_net_inflow"`
	MiddleNetInflow   float64 `json:"middle_net_inflow"`
	LargeNetInflow    float64 `json:"large_net_inflow"`
	SuperNetInflow    float64 `json:"super_net_inflow"`
}

type DayLine struct {
	StockLine          `xorm:"extends"`
	PctMainNetInflow   float64 `json:"pct_main_net_inflow"`
	PctSmallNetInflow  float64 `json:"pct_small_net_inflow"`
	PctMiddleNetInflow float64 `json:"pct_middle_net_inflow"`
	PctLargeNetInflow  float64 `json:"pct_large_net_inflow"`
	PctSuperNetInflow  float64 `json:"pct_super_net_inflow"`
	ChgClose           float64 `json:"chg_close"`
	PctChgOpen         float64 `json:"pct_chg_open"`
	PctChgHigh         float64 `json:"pct_chg_high"`
	PctChgLow          float64 `json:"pct_chg_low"`
	PctChgClose        float64 `json:"pct_chg_close"`
	PctChgAmount       float64 `json:"pct_chg_amount"`
	PctChgVol          float64 `json:"pct_chg_vol"`
	Ma3Close           float64 `json:"ma3_close"`
	Ma5Close           float64 `json:"ma5_close"`
	Ma10Close          float64 `json:"ma10_close"`
	Ma13Close          float64 `json:"ma13_close"`
	Ma20Close          float64 `json:"ma20_close"`
	Ma21Close          float64 `json:"ma21_close"`
	Ma30Close          float64 `json:"ma30_close"`
	Ma34Close          float64 `json:"ma34_close"`
	Ma55Close          float64 `json:"ma55_close"`
	Ma60Close          float64 `json:"ma60_close"`
	Ma90Close          float64 `json:"ma90_close"`
	Ma120Close         float64 `json:"ma120_close"`
	Ma144Close         float64 `json:"ma144_close"`
	Ma233Close         float64 `json:"ma233_close"`
	Ma240Close         float64 `json:"ma240_close"`
	Max3Close          float64 `json:"max3_close"`
	Max5Close          float64 `json:"max5_close"`
	Max10Close         float64 `json:"max10_close"`
	Max13Close         float64 `json:"max13_close"`
	Max20Close         float64 `json:"max20_close"`
	Max21Close         float64 `json:"max21_close"`
	Max30Close         float64 `json:"max30_close"`
	Max34Close         float64 `json:"max34_close"`
	Max55Close         float64 `json:"max55_close"`
	Max60Close         float64 `json:"max60_close"`
	Max90Close         float64 `json:"max90_close"`
	Max120Close        float64 `json:"max120_close"`
	Max144Close        float64 `json:"max144_close"`
	Max233Close        float64 `json:"max233_close"`
	Max240Close        float64 `json:"max240_close"`
	Min3Close          float64 `json:"min3_close"`
	Min5Close          float64 `json:"min5_close"`
	Min10Close         float64 `json:"min10_close"`
	Min13Close         float64 `json:"min13_close"`
	Min20Close         float64 `json:"min20_close"`
	Min21Close         float64 `json:"min21_close"`
	Min30Close         float64 `json:"min30_close"`
	Min34Close         float64 `json:"min34_close"`
	Min55Close         float64 `json:"min55_close"`
	Min60Close         float64 `json:"min60_close"`
	Min90Close         float64 `json:"min90_close"`
	Min120Close        float64 `json:"min120_close"`
	Min144Close        float64 `json:"min144_close"`
	Min233Close        float64 `json:"min233_close"`
	Min240Close        float64 `json:"min240_close"`
	//1天前的各均线
	Before1Ma3Close  float64 `json:"before1_ma3_close"`
	Before1Ma5Close  float64 `json:"before1_ma5_close"`
	Before1Ma10Close float64 `json:"before1_ma10_close"`
	Before1Ma13Close float64 `json:"before1_ma13_close"`
	Before1Ma20Close float64 `json:"before1_ma20_close"`
	Before1Ma21Close float64 `json:"before1_ma21_close"`
	Before1Ma30Close float64 `json:"before1_ma30_close"`
	Before1Ma34Close float64 `json:"before1_ma34_close"`
	Before1Ma55Close float64 `json:"before1_ma55_close"`
	Before1Ma60Close float64 `json:"before1_ma60_close"`
	//3天前的各均线
	Before3Ma3Close  float64 `json:"before3_ma3_close"`
	Before3Ma5Close  float64 `json:"before3_ma5_close"`
	Before3Ma10Close float64 `json:"before3_ma10_close"`
	Before3Ma13Close float64 `json:"before3_ma13_close"`
	Before3Ma20Close float64 `json:"before3_ma20_close"`
	Before3Ma21Close float64 `json:"before3_ma21_close"`
	Before3Ma30Close float64 `json:"before3_ma30_close"`
	Before3Ma34Close float64 `json:"before3_ma34_close"`
	Before3Ma55Close float64 `json:"before3_ma55_close"`
	Before3Ma60Close float64 `json:"before3_ma60_close"`
	//5天前的各均线
	Before5Ma3Close      float64 `json:"before5_ma3_close"`
	Before5Ma5Close      float64 `json:"before5_ma5_close"`
	Before5Ma10Close     float64 `json:"before5_ma10_close"`
	Before5Ma13Close     float64 `json:"before5_ma13_close"`
	Before5Ma20Close     float64 `json:"before5_ma20_close"`
	Before5Ma21Close     float64 `json:"before5_ma21_close"`
	Before5Ma30Close     float64 `json:"before5_ma30_close"`
	Before5Ma34Close     float64 `json:"before5_ma34_close"`
	Before5Ma55Close     float64 `json:"before5_ma55_close"`
	Before5Ma60Close     float64 `json:"before5_ma60_close"`
	Acc3PctChgClose      float64 `json:"acc3_pct_chg_close"`
	Acc5PctChgClose      float64 `json:"acc5_pct_chg_close"`
	Acc10PctChgClose     float64 `json:"acc10_pct_chg_close"`
	Acc13PctChgClose     float64 `json:"acc13_pct_chg_close"`
	Acc20PctChgClose     float64 `json:"acc20_pct_chg_close"`
	Acc21PctChgClose     float64 `json:"acc21_pct_chg_close"`
	Acc30PctChgClose     float64 `json:"acc30_pct_chg_close"`
	Acc34PctChgClose     float64 `json:"acc34_pct_chg_close"`
	Acc55PctChgClose     float64 `json:"acc55_pct_chg_close"`
	Acc60PctChgClose     float64 `json:"acc60_pct_chg_close"`
	Acc90PctChgClose     float64 `json:"acc90_pct_chg_close"`
	Acc120PctChgClose    float64 `json:"acc120_pct_chg_close"`
	Acc144PctChgClose    float64 `json:"acc144_pct_chg_close"`
	Acc233PctChgClose    float64 `json:"acc233_pct_chg_close"`
	Acc240PctChgClose    float64 `json:"acc240_pct_chg_close"`
	Future1PctChgClose   float64 `json:"future1_pct_chg_close"`
	Future3PctChgClose   float64 `json:"future3_pct_chg_close"`
	Future5PctChgClose   float64 `json:"future5_pct_chg_close"`
	Future10PctChgClose  float64 `json:"future10_pct_chg_close"`
	Future13PctChgClose  float64 `json:"future13_pct_chg_close"`
	Future20PctChgClose  float64 `json:"future20_pct_chg_close"`
	Future21PctChgClose  float64 `json:"future21_pct_chg_close"`
	Future30PctChgClose  float64 `json:"future30_pct_chg_close"`
	Future34PctChgClose  float64 `json:"future34_pct_chg_close"`
	Future55PctChgClose  float64 `json:"future55_pct_chg_close"`
	Future60PctChgClose  float64 `json:"future60_pct_chg_close"`
	Future90PctChgClose  float64 `json:"future90_pct_chg_close"`
	Future120PctChgClose float64 `json:"future120_pct_chg_close"`
	Future144PctChgClose float64 `json:"future144_pct_chg_close"`
	Future233PctChgClose float64 `json:"future233_pct_chg_close"`
	Future240PctChgClose float64 `json:"future240_pct_chg_close"`
}

func (DayLine) TableName() string {
	return "stk_dayline"
}

func (DayLine) KeyName() string {
	return entity.FieldName_Id
}

func (DayLine) IdName() string {
	return entity.FieldName_Id
}
