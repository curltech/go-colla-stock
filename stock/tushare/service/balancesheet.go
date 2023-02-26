package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// 获取上市公司资产负债表,用户需要至少800积分才可以调取,具体请参阅积分获取办法 https://tushare.pro/document/1?doc_id=13
func (ts *BalanceSheetService) BalanceSheet(params SheetRequest) (tsRsp *TushareResponse, err error) {
	resp, err := FastPost(&TushareRequest{
		ApiName: "balancesheet",
		Token:   token,
		Params:  params,
		Fields:  reflectFields(entity.BalanceSheet{}),
	})
	return resp, err
}

func assembleBalanceSheet(tsRsp *TushareResponse) []*entity.BalanceSheet {
	tsData := []*entity.BalanceSheet{}
	for _, data := range tsRsp.Data.Items {
		body, err := ReflectResponseData(tsRsp.Data.Fields, data)
		if err == nil {
			n := new(entity.BalanceSheet)
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
type BalanceSheetService struct {
	service.OrmBaseService
}

var balanceSheetService = &BalanceSheetService{}

func GetBalanceSheetService() *BalanceSheetService {
	return balanceSheetService
}

var seqname = "seq_stock"

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
