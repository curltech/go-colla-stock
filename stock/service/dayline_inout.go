package service

import (
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// FindFlexPoint 最基本的查询买卖点的方法，最为灵活
// 条件包括tsCode，tradeDate，filterContent（filterParas），startDate，endDate
// 如果filterContent中有?，则filterParas中必须有对应的参数值
func (svc *DayLineService) FindFlexPoint(tsCode string, tradeDate int64, condContent string, condParas []interface{}, startDate int64, endDate int64, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	dayLines := make([]*entity.DayLine, 0)
	conds += " and ma3close is not null and ma3close!=0 and (high-low)!=0"
	var err error
	if tradeDate != 0 {
		conds = conds + " and tradedate=?"
		paras = append(paras, tradeDate)
	}
	if condContent != "" {
		conds = conds + " and " + condContent
		if condParas != nil && len(condParas) > 0 {
			paras = append(paras, condParas...)
		}
	}
	if startDate != 0 {
		conds = conds + " and tradedate>?"
		paras = append(paras, startDate)
	}
	if endDate != 0 {
		conds += " and tradedate<=?"
		paras = append(paras, endDate)
	}
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, 0, err
		}
	}

	err = svc.Find(&dayLines, nil, "tscode,tradedate desc", from, limit, conds, paras...)

	return dayLines, count, nil
}
