package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type ProcessLog struct {
	entity.BaseEntity `xorm:"extends"`
	Name              string     `xorm:"varchar(256)" json:"name,omitempty"`
	BizCode           string     `xorm:"varchar(256)" json:"biz_code,omitempty"`
	SchedualDate      *time.Time `json:"schedual_date,omitempty"`
	MethodName        string     `xorm:"varchar(512)" json:"method_name,omitempty"`
	ErrorCode         string     `xorm:"varchar(32)" json:"error_code,omitempty"`
	ErrorMsg          string     `xorm:"varchar(5120)" json:"error_msg,omitempty"`
	StartDate         *time.Time `json:"start_date,omitempty"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	Elapse            int64      `json:"elapse,omitempty"`
}

func (ProcessLog) TableName() string {
	return "stk_processlog"
}

func (ProcessLog) KeyName() string {
	return entity.FieldName_Id
}

func (ProcessLog) IdName() string {
	return entity.FieldName_Id
}
