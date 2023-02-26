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
type FilterCondService struct {
	service.OrmBaseService
}

var filterCondService = &FilterCondService{}

func GetFilterCondService() *FilterCondService {
	return filterCondService
}

func (this *FilterCondService) GetSeqName() string {
	return seqname
}

func (this *FilterCondService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.FilterCond{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *FilterCondService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.FilterCond, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

var filterCondCache map[string]*entity.FilterCond = nil

func GetCacheFilterCond() map[string]*entity.FilterCond {
	if filterCondCache == nil {
		filterCondCache = make(map[string]*entity.FilterCond, 0)
		filterConds := make([]*entity.FilterCond, 0)
		svc := GetFilterCondService()
		err := svc.Find(&filterConds, nil, "", 0, 0, "")
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
		i := 0
		for _, filterCond := range filterConds {
			filterCondCache[filterCond.CondCode] = filterCond
			i++
		}
	}

	return filterCondCache
}

func RefreshCacheFilterCond() {
	filterCondCache = nil
}

func init() {
	service.GetSession().Sync(new(entity.FilterCond))
	filterCondService.OrmBaseService.GetSeqName = filterCondService.GetSeqName
	filterCondService.OrmBaseService.FactNewEntity = filterCondService.NewEntity
	filterCondService.OrmBaseService.FactNewEntities = filterCondService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("filtercond", filterCondService)
}
