package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Rhythmic struct {
	entity.BaseEntity `xorm:"extends"`
	Name              string `xorm:"varchar(255)" json:"collection,omitempty"` // 词牌名
	Notes             string `xorm:"varchar(32000)" json:"notes,omitempty"`    // 注释
	RhythmicType      string `xorm:"varchar(255)" json:"poemType,omitempty"`   // 诗词曲
}

func (Rhythmic) TableName() string {
	return "pm_rhythmic"
}

func (Rhythmic) KeyName() string {
	return entity.FieldName_Id
}

func (Rhythmic) IdName() string {
	return entity.FieldName_Id
}
