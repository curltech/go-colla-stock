package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

/**
抖音的评论， https://www.douyin.com/aweme/v1/web/comment/list/
https://www.douyin.com/aweme/v1/web/comment/list/reply/
comment_id
'cursor': page * 20,
'count': 20,
*/

type Comment struct {
	entity.BaseEntity `xorm:"extends"`
	CommentId         string `xorm:"varchar(255) index notnull" json:"cid"`
	Text              string `json:"text"`
	AwemeId           string `json:"aweme_id"` //视频编号
	DiggCount         int64  `json:"digg_count"`
	UserId            string `json:"uid"`
	Nickname          string `json:"nickname"`
	Signature         string `json:"signature"`
	ReplyId           string `json:"reply_id"`
	UserDigged        int64  `json:"user_digged"`
	ReplyComment      string `json:"reply_comment"`
	ReplyCommentTotal int64  `json:"reply_comment_total"`
	IsHot             bool   `json:"is_hot"`
	ItemCommentTotal  int64  `json:"item_comment_total"`
}

func (Comment) TableName() string {
	return "social_comment"
}

func (Comment) KeyName() string {
	return entity.FieldName_Id
}

func (Comment) IdName() string {
	return entity.FieldName_Id
}
