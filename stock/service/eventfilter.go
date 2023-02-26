package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

/**
同步表结构，服务继承基本服务的方法
*/
type EventFilterService struct {
	service.OrmBaseService
}

var eventFilterService = &EventFilterService{}

func GetEventFilterService() *EventFilterService {
	return eventFilterService
}

func (this *EventFilterService) GetSeqName() string {
	return seqname
}

func (this *EventFilterService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.EventFilter{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *EventFilterService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.EventFilter, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

var eventFilterCache map[string][]*entity.EventFilter = nil

func (this *EventFilterService) GetCacheEventFilter() map[string][]*entity.EventFilter {
	if eventFilterCache == nil {
		eventFilterCache = make(map[string][]*entity.EventFilter, 0)
		eventFilters := make([]*entity.EventFilter, 0)
		svc := GetEventFilterService()
		err := svc.Find(&eventFilters, nil, "", 0, 0, "")
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
		i := 0
		for _, eventFilter := range eventFilters {
			efs, ok := eventFilterCache[eventFilter.EventCode]
			if !ok {
				efs = make([]*entity.EventFilter, 0)
			}
			efs = append(efs, eventFilter)
			eventFilterCache[eventFilter.EventCode] = efs
			i++
		}
	}

	return eventFilterCache
}

func (this *EventFilterService) RefreshCacheEventFilter() {
	eventFilterCache = nil
}

func init() {
	service.GetSession().Sync(new(entity.EventFilter))
	eventFilterService.OrmBaseService.GetSeqName = eventFilterService.GetSeqName
	eventFilterService.OrmBaseService.FactNewEntity = eventFilterService.NewEntity
	eventFilterService.OrmBaseService.FactNewEntities = eventFilterService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("inoutfilter", eventFilterService)
}
