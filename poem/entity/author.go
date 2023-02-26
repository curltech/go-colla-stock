package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Author struct {
	entity.BaseEntity `xorm:"extends"`
	Name              string `xorm:"varchar(255)" json:"name,omitempty"`    // 作者名
	Notes             string `xorm:"varchar(32000)" json:"notes,omitempty"` // 注释
	Dynasty           string `xorm:"varchar(255)" json:"dynasty,omitempty"` // 朝代
}

func (Author) TableName() string {
	return "pm_author"
}

func (Author) KeyName() string {
	return entity.FieldName_Id
}

func (Author) IdName() string {
	return entity.FieldName_Id
}
