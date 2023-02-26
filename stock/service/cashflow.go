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
type CashFlowService struct {
	service.OrmBaseService
}

var cashFlowService = &CashFlowService{}

func GetCashFlowService() *CashFlowService {
	return cashFlowService
}

func (this *CashFlowService) GetSeqName() string {
	return seqname
}

func (this *CashFlowService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.CashFlow{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *CashFlowService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.CashFlow, 0)
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
	service.GetSession().Sync(new(entity.CashFlow))
	cashFlowService.OrmBaseService.GetSeqName = cashFlowService.GetSeqName
	cashFlowService.OrmBaseService.FactNewEntity = cashFlowService.NewEntity
	cashFlowService.OrmBaseService.FactNewEntities = cashFlowService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("cashflow", cashFlowService)
}
