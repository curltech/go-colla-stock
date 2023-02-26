package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Share struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string `xorm:"unique notnull" json:"ts_code,omitempty"` // str	Y	TS代码
	Symbol            string `json:"symbol,omitempty"`                        // str 股票代码
	Name              string `xorm:"index notnull" json:"name,omitempty"`     // str 股票名称
	Area              string `json:"area,omitempty"`                          // str 所在地域
	Industry          string `xorm:"index" json:"industry,omitempty"`         // str 所属行业
	Sector            string `xorm:"index" json:"sector,omitempty"`           // str 所属细分行业行业
	Fullname          string `json:"fullname,omitempty"`                      // str 股票全称
	Enname            string `json:"enname,omitempty"`                        // str 英文全称
	Market            string `json:"market,omitempty"`                        // str 市场类型 (主板/中小板/创业板/科创板/CDR)
	Exchange          string `json:"exchange,omitempty"`                      // str 交易所代码
	CurrType          string `json:"curr_type,omitempty"`                     // str 交易货币
	ListStatus        string `json:"list_status,omitempty"`                   // str 上市状态: L上市 D退市 P暂停上市
	ListDate          string `json:"list_date,omitempty"`                     // str 上市日期
	DelistDate        string `json:"delist_date,omitempty"`                   // str 退市日期
	IsHs              string `json:"is_hs,omitempty"`                         // str 是否沪深港通标的,N否 H沪股通 S深股通
	PinYin            string `xorm:"index" json:"pin_yin,omitempty"`
}

func (Share) TableName() string {
	return "stk_share"
}

func (Share) KeyName() string {
	return "ShareId"
}

func (Share) IdName() string {
	return entity.FieldName_Id
}
