package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/tushare/entity"
)

// 获取上市公司财务指标数据,为避免服务器压力,现阶段每次请求最多返回60条记录,可通过设置日期多次请求获取更多数据,用户需要至少800积分才可以调取,具体请参阅积分获取办法 https://tushare.pro/document/1?doc_id=13
func (ts *FinancialIndicatorService) FinancialIndicator(params FinanceRequest) (tsRsp *TushareResponse, err error) {
	resp, err := FastPost(&TushareRequest{
		ApiName: "financialIndicator",
		Token:   token,
		Params:  params,
		Fields:  reflectFields(entity.FinancialIndicator{}),
	})
	return resp, err
}

func assembleFinancialIndicator(tsRsp *TushareResponse) []*entity.FinancialIndicator {
	tsData := []*entity.FinancialIndicator{}
	for _, data := range tsRsp.Data.Items {
		body, err := ReflectResponseData(tsRsp.Data.Fields, data)
		if err == nil {
			n := new(entity.FinancialIndicator)
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
type FinancialIndicatorService struct {
	service.OrmBaseService
}

var financialIndicatorService = &FinancialIndicatorService{}

func GetFinancialIndicatorService() *FinancialIndicatorService {
	return financialIndicatorService
}

func (this *FinancialIndicatorService) GetSeqName() string {
	return seqname
}

func (this *FinancialIndicatorService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.FinancialIndicator{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *FinancialIndicatorService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.FinancialIndicator, 0)
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
	service.GetSession().Sync(new(entity.FinancialIndicator))
	financialIndicatorService.OrmBaseService.GetSeqName = financialIndicatorService.GetSeqName
	financialIndicatorService.OrmBaseService.FactNewEntity = financialIndicatorService.NewEntity
	financialIndicatorService.OrmBaseService.FactNewEntities = financialIndicatorService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("financialindicator", financialIndicatorService)
}
