package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

// EventFilter 事件的过滤器，事件与条件是多对多的关系，本表就是二者的关联表
// 相同事件代码的过滤器由多行filtercond记录组成，并设置自己的参数
type EventFilter struct {
	entity.BaseEntity `xorm:"extends"`
	EventCode         string  `json:"event_code,omitempty"`
	EventName         string  `json:"event_name,omitempty"`
	CondCode          string  `xorm:"varchar(255)" json:"cond_code,omitempty"`
	CodeAlias         string  `xorm:"varchar(255)" json:"code_alias,omitempty"`
	CondName          string  `xorm:"varchar(255)" json:"cond_name,omitempty"`
	CondAlias         string  `xorm:"varchar(255)" json:"cond_alias,omitempty"`
	CondContent       string  `xorm:"varchar(32000)" json:"cond_content,omitempty"`
	CondParas         string  `xorm:"varchar(255)" json:"cond_paras,omitempty"`
	Score             float64 `json:"score"`
	Descr             string  `xorm:"varchar(32000)" json:"descr,omitempty"`
}

func (EventFilter) TableName() string {
	return "stk_eventfilter"
}

func (EventFilter) KeyName() string {
	return entity.FieldName_Id
}

func (EventFilter) IdName() string {
	return entity.FieldName_Id
}
