package entity

import "github.com/curltech/go-colla-core/entity"

/**
周月季度年数据
*/
type WmqyLine struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string  `xorm:"varchar(255) index notnull" json:"ts_code"` // str	股票代码
	TradeDate         int64   `xorm:"index notnull" json:"trade_date"`
	QDate             string  `xorm:"varchar(255) index notnull" json:"qdate"`
	ShareNumber       float64 `json:"share_number"`
	Open              float64 `json:"open"`   // float	(1/5/15/30/60分钟) D日 W周 M月 开盘价
	High              float64 `json:"high"`   // float	(1/5/15/30/60分钟) D日 W周 M月 最高价
	Low               float64 `json:"low"`    // float	(1/5/15/30/60分钟) D日 W周 M月 最低价
	Close             float64 `json:"close"`  // float	(1/5/15/30/60分钟) D日 W周 M月 收盘价
	Vol               float64 `json:"vol"`    // float	(1/5/15/30/60分钟) D日 W周 M月 成交量 (手)
	Amount            float64 `json:"amount"` // float	(1/5/15/30/60分钟) D日 W周 M月 成交额 (千元)
	Turnover          float64 `json:"turnover"`
	PreClose          float64 `json:"pre_close"` // float	上一(1/5/15/30/60分钟) D日 W周 M月 收盘价
	ChgClose          float64 `json:"chg_close"`
	PctChgOpen        float64 `json:"pct_chg_open"`
	PctChgHigh        float64 `json:"pct_chg_high"`
	PctChgLow         float64 `json:"pct_chg_low"`
	PctChgClose       float64 `json:"pct_chg_close"` // float	(1/5/15/30/60分钟) D日 W周 M月 涨跌幅
	PctChgAmount      float64 `json:"pct_chg_amount"`
	PctChgVol         float64 `json:"pct_chg_vol"`
	LineType          int     `xorm:"index notnull" json:"line_type"`
}

func (WmqyLine) TableName() string {
	return "stk_wmqyline"
}

func (WmqyLine) KeyName() string {
	return entity.FieldName_Id
}

func (WmqyLine) IdName() string {
	return entity.FieldName_Id
}
