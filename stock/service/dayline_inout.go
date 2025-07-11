package service

import (
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
)

type InOutPoint struct {
	Data  []*entity.DayLine `json:"data,omitempty"`
	Count int64             `json:"count,omitempty"`
}

// FindFlexPoint 最基本的查询买卖点的方法，最为灵活
// 条件包括tsCode，tradeDate，filterContent（filterParas），startDate，endDate
// 如果filterContent中有?，则filterParas中必须有对应的参数值
func (svc *DayLineService) FindFlexPoint(tsCode string, tradeDate int64, filterContent string, filterParas []interface{}, startDate int64, endDate int64, from int, limit int, count int64) (*InOutPoint, error) {
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	dayLines := make([]*entity.DayLine, 0)
	conds += " and ma3close is not null and ma3close!=0 and (high-low)!=0"
	inOutPoint := &InOutPoint{}
	var err error
	if tradeDate != 0 {
		conds = conds + " and tradedate=?"
		paras = append(paras, tradeDate)
	}
	if filterContent != "" {
		conds = conds + " and " + filterContent
		if filterParas != nil && len(filterParas) > 0 {
			paras = append(paras, filterParas...)
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
			return nil, err
		}
		inOutPoint.Count = count
	}

	err = svc.Find(&dayLines, nil, "tscode,tradedate desc", from, limit, conds, paras...)
	if err != nil {
		return inOutPoint, err
	}
	inOutPoint.Data = dayLines

	return inOutPoint, nil
}
