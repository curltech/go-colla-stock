package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// WmqyLineService 同步表结构，服务继承基本服务的方法
type WmqyLineService struct {
	service.OrmBaseService
}

var wmqyLineService = &WmqyLineService{}

func GetWmqyLineService() *WmqyLineService {
	return wmqyLineService
}

func (svc *WmqyLineService) GetSeqName() string {
	return seqname
}

func (svc *WmqyLineService) NewEntity(data []byte) (interface{}, error) {
	wmqyLine := &entity.WmqyLine{}
	if data == nil {
		return wmqyLine, nil
	}
	err := message.Unmarshal(data, wmqyLine)
	if err != nil {
		return nil, err
	}

	return wmqyLine, err
}

func (svc *WmqyLineService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.WmqyLine, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

var WmqylineHeader = []string{"ts_code", "security_name", "qdate", "share_number", "ep", "high", "close", "pct_chg_high",
	"pct_chg_close", "weight_avg_roe", "gross_profit_margin", "parent_net_profit", "basic_eps",
	"or_last_month", "np_last_month", "yoy_sales", "yoy_dedu_np", "cfps", "dividend_yield_ratio"}

func (svc *WmqyLineService) findMaxTradeDate(tsCode string, lineType int) (*entity.WmqyLine, *entity.WmqyLine, error) {
	cond := &entity.WmqyLine{}
	cond.TsCode = tsCode
	cond.LineType = lineType
	wmqyLines := make([]*entity.WmqyLine, 0)
	err := svc.Find(&wmqyLines, cond, "tradedate desc", 0, 2, "")
	if err != nil {
		return nil, nil, err
	}
	if len(wmqyLines) > 1 {
		return wmqyLines[0], wmqyLines[1], nil
	} else if len(wmqyLines) > 0 {
		return wmqyLines[0], nil, nil
	}

	return nil, nil, nil
}

// FindPreceding 获取某时间点前limit条数据，如果没有日期范围的指定，就是返回最新的回溯limit条数据
func (svc *WmqyLineService) FindPreceding(tsCode string, lineType int, endDate string, from int, limit int, count int64) ([]*entity.WmqyLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	wmqyLines := make([]*entity.WmqyLine, 0)
	conds += " and close is not null and close!=0"
	if lineType != 0 {
		conds = conds + " and linetype=?"
		paras = append(paras, lineType)
	}
	if endDate != "" {
		conds += " and qdate<=?"
		paras = append(paras, endDate)
	}
	var err error
	condiBean := &entity.WmqyLine{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	err = svc.Find(&wmqyLines, nil, "tscode,linetype,qdate,tradedate desc", from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	length := len(wmqyLines)
	ps := make([]*entity.WmqyLine, length)
	qdate := ""
	for i := length; i > 0; i-- {
		wmqyLine := wmqyLines[i-1]
		if qdate == wmqyLine.QDate {
			svc.Delete(nil, "id=?", wmqyLine.Id)
			count--
		} else {
			qdate = wmqyLine.QDate
			ps[length-i] = wmqyLine
		}
	}
	if len(ps) > 0 {
		logger.Sugar.Infof("from %v to %v WmqyLine data", ps[0].QDate, ps[len(ps)-1].QDate)
	} else {
		logger.Sugar.Errorf("WmqyLine len 0")
	}
	return ps, count, nil
}

// FindFollowing 获取某时间点后limit条数据，如果没有日期范围的指定，就是返回最早limit条数据
func (svc *WmqyLineService) FindFollowing(tsCode string, lineType int, startDate string, endDate string, from int, limit int, count int64) ([]*entity.WmqyLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	wmqyLines := make([]*entity.WmqyLine, 0)
	conds += " and close is not null and close!=0"
	if lineType != 0 {
		conds = conds + " and linetype=?"
		paras = append(paras, lineType)
	}
	if startDate != "" {
		conds = conds + " and qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		conds = conds + " and qdate<=?"
		paras = append(paras, endDate)
	}
	var err error
	condiBean := &entity.WmqyLine{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	err = svc.Find(&wmqyLines, nil, "tscode,linetype,qdate", from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	if len(wmqyLines) > 0 {
		logger.Sugar.Infof("from %v to %v wmqyLines data", wmqyLines[0].QDate, wmqyLines[len(wmqyLines)-1].QDate)
	}
	return wmqyLines, count, nil
}

func (svc *WmqyLineService) deleteWmqyLine(tscode string, qdate string) error {
	wmqyline := &entity.WmqyLine{}
	conds := "tscode=? and qdate=?"
	paras := make([]interface{}, 0)
	paras = append(paras, tscode)
	paras = append(paras, qdate)
	_, err := svc.Delete(wmqyline, conds, paras...)

	return err
}

func init() {
	err := service.GetSession().Sync(new(entity.WmqyLine))
	if err != nil {
		return
	}
	wmqyLineService.OrmBaseService.GetSeqName = wmqyLineService.GetSeqName
	wmqyLineService.OrmBaseService.FactNewEntity = wmqyLineService.NewEntity
	wmqyLineService.OrmBaseService.FactNewEntities = wmqyLineService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("wmqyline", wmqyLineService)
}
