package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Poem struct {
	entity.BaseEntity `xorm:"extends"`
	Collection        string `xorm:"varchar(255)" json:"collection,omitempty"`   // 书名
	Title             string `xorm:"varchar(255)" json:"title,omitempty"`        // 标题
	Chapter           string `xorm:"varchar(255)" json:"chapter,omitempty"`      // 章节
	Section           string `xorm:"varchar(255)" json:"section,omitempty"`      // 部分
	Notes             string `xorm:"varchar(32000)" json:"notes,omitempty"`      // 注释
	Author            string `xorm:"varchar(255)" json:"author,omitempty"`       // 作者
	Rhythmic          string `xorm:"varchar(255)" json:"rhythmic,omitempty"`     // 词牌
	Paragraphs        string `xorm:"varchar(32000)" json:"paragraphs,omitempty"` // 段落
	PoemType          string `xorm:"varchar(255)" json:"poemType,omitempty"`     // 诗词曲
	Dynasty           string `xorm:"varchar(255)" json:"dynasty,omitempty"`      // 朝代
}

func (Poem) TableName() string {
	return "pm_poem"
}

func (Poem) KeyName() string {
	return entity.FieldName_Id
}

func (Poem) IdName() string {
	return entity.FieldName_Id
}
