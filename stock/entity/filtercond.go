package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

// FilterCond 股票发生的事件
type FilterCond struct {
	entity.BaseEntity `xorm:"extends"`
	CondCode          string  `json:"cond_code,omitempty"`
	CondType          string  `json:"cond_type,omitempty"` //Trend，Occasion，Situation
	Name              string  `json:"name,omitempty"`
	Content           string  `xorm:"varchar(32000)" json:"content,omitempty"`
	CondParas         string  `json:"cond_paras,omitempty"`
	Score             float64 `json:"score"` //满足条件的评分，正数表示上涨，负数表示下跌
	Descr             string  `xorm:"varchar(32000)" json:"descr,omitempty"`
}

func (FilterCond) TableName() string {
	return "stk_filtercond"
}

func (FilterCond) KeyName() string {
	return entity.FieldName_Id
}

func (FilterCond) IdName() string {
	return entity.FieldName_Id
}
