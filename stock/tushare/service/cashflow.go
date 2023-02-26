package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// 获取新股上市列表数据,单次最大2000条,总量不限制,用户需要至少120积分才可以调取,具体请参阅积分获取办法 https://tushare.pro/document/1?doc_id=13
func (ts *CashFlowService) CashFlow(params SheetRequest) (tsRsp *TushareResponse, err error) {
	resp, err := FastPost(&TushareRequest{
		ApiName: "cashflow",
		Token:   token,
		Params:  params,
		Fields:  reflectFields(entity.CashFlow{}),
	})
	return resp, err
}

func assembleCashFlow(tsRsp *TushareResponse) []*entity.CashFlow {
	tsData := []*entity.CashFlow{}
	for _, data := range tsRsp.Data.Items {
		body, err := ReflectResponseData(tsRsp.Data.Fields, data)
		if err == nil {
			n := new(entity.CashFlow)
			err = json.Unmarshal(body, &n)
			if err == nil {
				tsData = append(tsData, n)
			}
		}
	}
	return tsData
}

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
