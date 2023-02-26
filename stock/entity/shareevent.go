package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

/**
代表某只股票的发生了某种事件，这种事件关于季报，买入卖出机会
*/
type ShareEvent struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string  `json:"ts_code,omitempty"`
	TradeDate         int64   `json:"trade_date,omitempty"`
	EventCode         string  `json:"event_code,omitempty"`
	EventName         string  `json:"event_name,omitempty"`
	Score             float64 `json:"score"`
	Descr             string  `json:"Descr,omitempty"`
	Pe                float64 `json:"pe"`
	Peg               float64 `json:"peg"`
	PercentilePe      float64 `json:"percentile_pe"`
	PercentilePeg     float64 `json:"percentile_peg"`
}

func (ShareEvent) TableName() string {
	return "stk_shareevent"
}

func (ShareEvent) KeyName() string {
	return entity.FieldName_Id
}

func (ShareEvent) IdName() string {
	return entity.FieldName_Id
}
