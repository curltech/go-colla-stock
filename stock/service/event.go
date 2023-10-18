package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// EventService 同步表结构，服务继承基本服务的方法
type EventService struct {
	service.OrmBaseService
}

var eventService = &EventService{}

func GetEventService() *EventService {
	return eventService
}

func (svc *EventService) GetSeqName() string {
	return seqname
}

func (svc *EventService) NewEntity(data []byte) (interface{}, error) {
	event := &entity.Event{}
	if data == nil {
		return event, nil
	}
	err := message.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}

	return event, err
}

func (svc *EventService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Event, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

var eventCache map[string]*entity.Event = nil

func (svc *EventService) GetCacheEvent() map[string]*entity.Event {
	if eventCache == nil {
		eventCache = make(map[string]*entity.Event, 0)
		events := make([]*entity.Event, 0)
		svc := GetEventService()
		condiBean := &entity.Event{}
		condiBean.Status = "Enabled"
		err := svc.Find(&events, condiBean, "", 0, 0, "")
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
		for _, event := range events {
			eventCache[event.EventCode] = event
		}
	}

	return eventCache
}

func (svc *EventService) RefreshCacheEvent() {
	eventCache = nil
}

func init() {
	err := service.GetSession().Sync(new(entity.Event))
	if err != nil {
		return
	}
	eventService.OrmBaseService.GetSeqName = eventService.GetSeqName
	eventService.OrmBaseService.FactNewEntity = eventService.NewEntity
	eventService.OrmBaseService.FactNewEntities = eventService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("event", eventService)
}
