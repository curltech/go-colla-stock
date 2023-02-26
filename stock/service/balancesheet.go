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
type BalanceSheetService struct {
	service.OrmBaseService
}

var balanceSheetService = &BalanceSheetService{}

func GetBalanceSheetService() *BalanceSheetService {
	return balanceSheetService
}

func (this *BalanceSheetService) GetSeqName() string {
	return seqname
}

func (this *BalanceSheetService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.BalanceSheet{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *BalanceSheetService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.BalanceSheet, 0)
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
	service.GetSession().Sync(new(entity.BalanceSheet))
	balanceSheetService.OrmBaseService.GetSeqName = balanceSheetService.GetSeqName
	balanceSheetService.OrmBaseService.FactNewEntity = balanceSheetService.NewEntity
	balanceSheetService.OrmBaseService.FactNewEntities = balanceSheetService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("balancesheet", balanceSheetService)
}
