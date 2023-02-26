package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/poem/entity"
)

/**
同步表结构，服务继承基本服务的方法
*/
type RhythmicService struct {
	service.OrmBaseService
}

var rhythmicService = &RhythmicService{}

func GetRhythmicService() *RhythmicService {
	return rhythmicService
}

func (this *RhythmicService) GetSeqName() string {
	return seqname
}

func (this *RhythmicService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Rhythmic{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *RhythmicService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Rhythmic, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func init() {
	service.GetSession().Sync(new(entity.Rhythmic))
	rhythmicService.OrmBaseService.GetSeqName = rhythmicService.GetSeqName
	rhythmicService.OrmBaseService.FactNewEntity = rhythmicService.NewEntity
	rhythmicService.OrmBaseService.FactNewEntities = rhythmicService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("rhythmic", rhythmicService)
}
