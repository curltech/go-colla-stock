package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

/**
股票组合统计
*/
type PortfolioStat struct {
	entity.BaseEntity  `xorm:"extends"`
	TsCode             string  `xorm:"index notnull" json:"ts_code"`
	SecurityName       string  `json:"security_name"`
	TargetTsCode       string  `xorm:"index notnull" json:"target_ts_code"`
	TargetSecurityName string  `json:"target_security_name"`
	StartDate          int64   `xorm:"index notnull" json:"start_date"` //统计指标对应的起始季度
	EndDate            int64   `json:"end_date"`                        //统计指标对应的结束季度
	Term               int64   `json:"term"`
	Source             string  `xorm:"index notnull" json:"source"`      //统计指标
	SourceName         string  `xorm:"index notnull" json:"source_name"` //统计指标对应的字段名
	StatValue          float64 `json:"stat_value"`
}

func (PortfolioStat) TableName() string {
	return "stk_portfoliostat"
}

func (PortfolioStat) KeyName() string {
	return entity.FieldName_Id
}

func (PortfolioStat) IdName() string {
	return entity.FieldName_Id
}
