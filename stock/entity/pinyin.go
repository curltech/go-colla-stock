package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type PinYin struct {
	entity.BaseEntity `xorm:"extends"`
	ChineseChar       string `xorm:"unique notnull" json:"chinese_char,omitempty"`
	PinYin            string `json:"pin_yin,omitempty"`
	FirstChar         string `json:"first_char,omitempty"`
}

func (PinYin) TableName() string {
	return "stk_pinyin"
}

func (PinYin) KeyName() string {
	return entity.FieldName_Id
}

func (PinYin) IdName() string {
	return entity.FieldName_Id
}
