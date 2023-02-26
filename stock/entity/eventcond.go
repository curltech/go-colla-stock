package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

/**
代表某只股票的发生了某种事件的条件满足的细节
*/
type EventCond struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string  `xorm:"index notnull" json:"ts_code,omitempty"`
	Name              string  `xorm:"-" json:"name,omitempty"`
	TradeDate         int64   `xorm:"index notnull" json:"trade_date,omitempty"`
	EventCode         string  `xorm:"index notnull" json:"event_code,omitempty"`
	EventType         string  `xorm:"index notnull" json:"event_type,omitempty"` //in,out,report
	EventName         string  `json:"event_name,omitempty"`
	CondCode          string  `json:"cond_code,omitempty"`
	CondName          string  `json:"cond_name,omitempty"`
	CondAlias         string  `json:"cond_alias,omitempty"`
	CondContent       string  `json:"cond_content,omitempty"`
	CondParas         string  `json:"cond_paras,omitempty"`
	CondValue         float64 `json:"cond_value,omitempty"`
	CondResult        float64 `json:"cond_result,omitempty"`
	Score             float64 `json:"score"` //事件条件的评分，正数表示上涨，负数表示下跌
	Descr             string  `json:"descr,omitempty"`
}

func (EventCond) TableName() string {
	return "stk_eventcond"
}

func (EventCond) KeyName() string {
	return entity.FieldName_Id
}

func (EventCond) IdName() string {
	return entity.FieldName_Id
}
