package service

import (
	"fmt"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"strings"
)

// EventCondService 同步表结构，服务继承基本服务的方法
type EventCondService struct {
	service.OrmBaseService
}

var eventCondService = &EventCondService{}

func GetEventCondService() *EventCondService {
	return eventCondService
}

func (svc *EventCondService) GetSeqName() string {
	return seqname
}

func (svc *EventCondService) NewEntity(data []byte) (interface{}, error) {
	eventCond := &entity.EventCond{}
	if data == nil {
		return eventCond, nil
	}
	err := message.Unmarshal(data, eventCond)
	if err != nil {
		return nil, err
	}

	return eventCond, err
}

func (svc *EventCondService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.EventCond, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (svc *EventCondService) findMaxTradeDate(tscode string) (int64, error) {
	conds, paras := stock.InBuildStr("tscode", tscode, ",")
	ps := make([]*entity.EventCond, 0)
	err := svc.Find(&ps, nil, "tradedate desc", 0, 1, conds, paras...)
	if err != nil {
		return 0, err
	}
	if len(ps) > 0 {
		return ps[0].TradeDate, nil
	}

	return 0, nil
}

func (svc *EventCondService) FindGroupby(tscode string, startDate int64, endDate int64, eventCode string, eventType string, orderby string, from int, limit int, count int64) ([]*entity.EventCond, int64, error) {
	conds, paras := stock.InBuildStr("tscode", tscode, ",")
	if startDate != 0 {
		conds = conds + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	if endDate != 0 {
		conds = conds + " and tradedate<=?"
		paras = append(paras, endDate)
	}
	if eventCode != "" {
		conds = conds + " and eventcode=?"
		paras = append(paras, eventCode)
	}
	if eventType != "" {
		conds = conds + " and eventType=?"
		paras = append(paras, eventType)
	}
	conds = "select tscode as ts_code,tradedate as trade_date,eventcode as event_code" +
		",eventname as event_name,eventtype as event_type,sum(score) as score from stk_eventcond where " + conds
	conds = conds + " group by tscode,tradedate,eventcode,eventname,eventtype"
	if orderby == "" {
		orderby = "tscode,tradedate desc"
	}
	conds = conds + " order by " + orderby
	if count == 0 {
		sql := "select count(*) as count from (" + conds + ") t"
		results, err := svc.Query(sql, paras...)
		if err != nil {
			return nil, count, err
		}
		if results != nil && len(results) > 0 {
			result := results[0]
			v, _ := result["count"]
			strVal := string(v)
			val, _ := convert.ToObject(strVal, "int64")
			count = val.(int64)
		}
	}
	sql := "select * from (" + conds + ") t offset " + fmt.Sprint(from)
	if limit > 0 {
		sql = sql + " limit " + fmt.Sprint(limit)
	}
	results, err := svc.Query(sql, paras...)
	if err != nil {
		return nil, count, err
	}
	jsonMap, _, _ := stock.GetJsonMap(entity.EventCond{})
	ps := make([]*entity.EventCond, 0)
	_, shareMap := shareService.GetCacheShare()
	for _, result := range results {
		p := &entity.EventCond{}
		for k, v := range result {
			strVal := string(v)
			if k == "score" {
				val, _ := convert.ToObject(strVal, "float64")
				p.Score = val.(float64)
			} else {
				fieldname, ok := jsonMap[k]
				if ok {
					err = reflect.Set(p, fieldname, strVal)
					if err != nil {
						return nil, 0, err
					}
				}
			}
		}
		if p.Name == "" {
			share, ok := shareMap[p.TsCode]
			if ok {
				p.Name = share.Name
			}
		}
		ps = append(ps, p)
	}

	return ps, count, nil
}

func (svc *EventCondService) RefreshEventCond() error {
	processLog := GetProcessLogService().StartLog("EventCond", "RefreshEventCond", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, svc.AsyncUpdateEventCond, nil)
	defer routinePool.Release()
	tsCodes, _ := GetShareService().GetCacheShare()
	for _, tsCode := range tsCodes {
		para := make([]interface{}, 0)
		para = append(para, tsCode)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(processLog, "", "")

	return nil
}

func (svc *EventCondService) AsyncUpdateEventCond(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)

	svc.GetUpdateEventCond(tscode)
}

func (svc *EventCondService) GetUpdateEventCond(tscode string) []interface{} {
	eventMap := GetEventService().GetCacheEvent()
	startDate, _ := svc.findMaxTradeDate(tscode)
	ps := make([]interface{}, 0)
	for _, event := range eventMap {
		if event.EventType == "in" || event.EventType == "out" {
			inOutPoint, err := GetDayLineService().FindInOutEvent(tscode, 0, event.EventCode, nil, startDate, 0, 0, 0, 0)
			if err != nil {
				continue
			}
			cvs := inOutPoint.CondValue
			for _, cv := range cvs {
				tradeDate := cv["trade_date"].(int64)
				for k, v := range cv {
					if k == "ts_code" || k == "trade_date" {

					} else if !strings.HasSuffix(k, "_cond") &&
						!strings.HasSuffix(k, "_name") &&
						!strings.HasSuffix(k, "_alias") &&
						!strings.HasSuffix(k, "_result") &&
						!strings.HasSuffix(k, "_paras") {
						ec := &entity.EventCond{TsCode: tscode, TradeDate: tradeDate}
						ec.EventCode = event.EventCode
						ec.EventName = event.EventName
						ec.CondCode = k
						ec.CondValue, _ = v.(float64)
						ec.CondContent, _ = cv[k+"_cond"].(string)
						ec.CondParas, _ = cv[k+"_paras"].(string)
						ec.CondName, _ = cv[k+"_name"].(string)
						ec.CondAlias, _ = cv[k+"_alias"].(string)
						ec.CondResult, _ = cv[k+"_result"].(float64)
						ec.EventType = event.EventType
						ps = append(ps, ec)
					}
				}
			}
		}
	}
	_, err := svc.Insert(ps...)
	if err != nil {
		return nil
	}

	return ps
}

func (svc *EventCondService) Search(tscode string, startDate int64, endDate int64, eventCode string, eventType string, orderby string, from int, limit int, count int64) ([]*entity.EventCond, int64, error) {
	conds, paras := stock.InBuildStr("tscode", tscode, ",")
	if startDate != 0 {
		conds = conds + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	if endDate != 0 {
		conds = conds + " and tradedate<=?"
		paras = append(paras, endDate)
	}
	if eventCode != "" {
		conds = conds + " and eventcode=?"
		paras = append(paras, eventCode)
	}
	if eventType != "" {
		conds = conds + " and eventType=?"
		paras = append(paras, eventType)
	}
	var err error
	condiBean := &entity.EventCond{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "tscode,tradedate desc"
	}
	ps := make([]*entity.EventCond, 0)
	err = svc.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	return ps, count, nil
}

func init() {
	err := service.GetSession().Sync(new(entity.EventCond))
	if err != nil {
		return
	}
	eventCondService.OrmBaseService.GetSeqName = eventCondService.GetSeqName
	eventCondService.OrmBaseService.FactNewEntity = eventCondService.NewEntity
	eventCondService.OrmBaseService.FactNewEntities = eventCondService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("eventcond", eventCondService)
}
