package service

import (
	"errors"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// FindByCondContent 最基本的查询买卖点的方法，最为灵活
// 条件包括tsCode，tradeDate，condContent（condParas)
// tsCode，tradeDate必有一个不为空
func (svc *DayLineService) FindByCondContent(tsCode string, tradeDate int64, condContent string, condParas []interface{}, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	if tsCode == "" && tradeDate == 0 {
		return nil, 0, errors.New("tsCode and tradeDate can't both be empty")
	}

	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	dayLines := make([]*entity.DayLine, 0)
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
