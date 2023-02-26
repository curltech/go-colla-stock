package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

/**
同步表结构，服务继承基本服务的方法
*/
type IncomeService struct {
	service.OrmBaseService
}

var incomeService = &IncomeService{}

func GetIncomeService() *IncomeService {
	return incomeService
}

func (this *IncomeService) GetSeqName() string {
	return seqname
}

func (this *IncomeService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Income{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *IncomeService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Income, 0)
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
	service.GetSession().Sync(new(entity.Income))
	incomeService.OrmBaseService.GetSeqName = incomeService.GetSeqName
	incomeService.OrmBaseService.FactNewEntity = incomeService.NewEntity
	incomeService.OrmBaseService.FactNewEntities = incomeService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("income", incomeService)
}
