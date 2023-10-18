package service

import (
	"errors"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// ShareEventService 同步表结构，服务继承基本服务的方法
type ShareEventService struct {
	service.OrmBaseService
}

var shareEventService = &ShareEventService{}

func GetShareEventService() *ShareEventService {
	return shareEventService
}

func (svc *ShareEventService) GetSeqName() string {
	return seqname
}

func (svc *ShareEventService) NewEntity(data []byte) (interface{}, error) {
	shareEvent := &entity.ShareEvent{}
	if data == nil {
		return shareEvent, nil
	}
	err := message.Unmarshal(data, shareEvent)
	if err != nil {
		return nil, err
	}

	return shareEvent, err
}

func (svc *ShareEventService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.ShareEvent, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (svc *ShareEventService) GetMine(userName string, tsCode string) ([]interface{}, error) {
	tscode := "000001"
	if tsCode != "" {
		tscode = tsCode
	}
	dayLines, err := GetDayLineService().findMaxTradeDate(tscode)
	if err != nil {
		return nil, err
	}
	if dayLines == nil || len(dayLines) == 0 || dayLines[0] == nil {
		return nil, errors.New("NoTradeDate")
	}
	sql := "select s.tscode as ts_code,s.name,s.industry,sd.close,sd.pctchgclose as pct_chg_close" +
		",sd.pctchgvol as pct_chg_vol" +
		" from stk_share s join stk_shareEvent ss on s.tscode = ss.tscode" +
		" join stk_dayline sd on sd.tscode=s.tscode" +
		" where ss.username = ?  and sd.tradedate = ?"
	paras := make([]interface{}, 0)
	paras = append(paras, userName)
	paras = append(paras, dayLines[0].TradeDate)
	if tsCode != "" {
		sql = sql + " and ss.tscode= ?"
		paras = append(paras, tsCode)
	}
	sql = sql + " order by s.industry"
	results, err := svc.Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	ps := make([]interface{}, 0)
	jsonMap, _, _ := stock.GetJsonMap(UserShare{})
	var i int64
	for _, result := range results {
		qp := &UserShare{}
		for colname, v := range result {
			err = reflect.Set(qp, jsonMap[colname], string(v))
			if err != nil {
				logger.Sugar.Errorf("Set colname %v value %v error", colname, string(v))
			}
		}
		i++
		qp.Id = i
		ps = append(ps, qp)
	}

	return ps, nil
}

func init() {
	err := service.GetSession().Sync(new(entity.ShareEvent))
	if err != nil {
		return
	}
	shareEventService.OrmBaseService.GetSeqName = shareEventService.GetSeqName
	shareEventService.OrmBaseService.FactNewEntity = shareEventService.NewEntity
	shareEventService.OrmBaseService.FactNewEntities = shareEventService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("shareevent", shareEventService)
}
