package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

/**
事件定义
*/
type Event struct {
	entity.BaseEntity `xorm:"extends"`
	EventCode         string     `json:"event_code,omitempty"`
	EventType         string     `json:"event_type,omitempty"` //in,out,report
	EventName         string     `json:"event_name,omitempty"`
	Content           string     `xorm:"varchar(32000)" json:"content,omitempty"`
	ContentParas      string     `xorm:"varchar(255)" json:"content_paras,omitempty"`
	Descr             string     `json:"descr,omitempty"`
	Status            string     `json:"status,omitempty"`
	StatusDate        *time.Time `json:"status_date,omitempty"`
}

func (Event) TableName() string {
	return "stk_event"
}

func (Event) KeyName() string {
	return entity.FieldName_Id
}

func (Event) IdName() string {
	return entity.FieldName_Id
}
