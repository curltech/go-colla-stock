package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Forecast struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string  `xorm:"varchar(255)" json:"ts_code,omitempty"` // str	Y	TS代码
	AnnDate           string  `json:"ann_date,omitempty"`                    // str	公告日期
	EndDate           string  `json:"end_date,omitempty"`                    // str	报告期
	Type              string  `json:"type,omitempty"`                        // str	业绩预告类型(预增/预减/扭亏/首亏/续亏/续盈/略增/略减)
	PChangeMin        float64 `json:"p_change_min,omitempty"`                // float	预告净利润变动幅度下限(%)
	PChangeMax        float64 `json:"p_change_max,omitempty"`                // float	预告净利润变动幅度上限(%)
	NetProfitMin      float64 `json:"net_profit_min,omitempty"`              // float	预告净利润下限(万元)
	NetProfitMax      float64 `json:"net_profit_max,omitempty"`              // float	预告净利润上限(万元)
	LastParentNet     float64 `json:"last_parent_net,omitempty"`             // float	上年同期归属母公司净利润
	FirstAnnDate      string  `json:"first_ann_date,omitempty"`              // str	首次公告日
	Summary           string  `json:"summary,omitempty"`                     // str	业绩预告摘要
	ChangeReason      string  `json:"change_reason,omitempty"`               // str	业绩变动原因
}

func (Forecast) TableName() string {
	return "stk_forecast"
}

func (Forecast) KeyName() string {
	return entity.FieldName_Id
}

func (Forecast) IdName() string {
	return entity.FieldName_Id
}
