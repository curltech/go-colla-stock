package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// 主要报表类型说明
// 代码	类型	说明
// 1	合并报表	上市公司最新报表(默认)
// 2	单季合并	单一季度的合并报表
// 3	调整单季合并表	调整后的单季合并报表(如果有)
// 4	调整合并报表	本年度公布上年同期的财务报表数据,报告期为上年度
// 5	调整前合并报表	数据发生变更,将原数据进行保留,即调整前的原数据
// 6	母公司报表	该公司母公司的财务报表数据
// 7	母公司单季表	母公司的单季度表
// 8	母公司调整单季表	母公司调整后的单季表
// 9	母公司调整表	该公司母公司的本年度公布上年同期的财务报表数据
// 10	母公司调整前报表	母公司调整之前的原始财务报表数据
// 11	调整前合并报表	调整之前合并报表原数据
// 12	母公司调整前报表	母公司报表发生变更前保留的原数据
type SheetRequest struct {
	TsCode     string `json:"ts_code,omitempty"`     // str	Y	股票代码
	AnnDate    string `json:"ann_date,omitempty"`    // str	N	公告日期
	StartDate  string `json:"start_date,omitempty"`  // str	N	公告开始日期
	EndDate    string `json:"end_date,omitempty"`    // str	N	公告结束日期
	Period     string `json:"period,omitempty"`      // str	N	报告期(每个季度最后一天的日期,比如20171231表示年报)
	ReportType string `json:"report_type,omitempty"` // str	N	报告类型： 参考下表说明
	CompType   string `json:"comp_type,omitempty"`   // str	N	公司类型：1一般工商业 2银行 3保险 4证券
}

// 获取上市公司财务利润表数据,用户需要至少800积分才可以调取,具体请参阅积分获取办法 https://tushare.pro/document/1?doc_id=13
func (ts *IncomeService) Income(params SheetRequest) (tsRsp *TushareResponse, err error) {
	resp, err := FastPost(&TushareRequest{
		ApiName: "income",
		Token:   token,
		Params:  params,
		Fields:  reflectFields(entity.Income{}),
	})
	return resp, err
}

func assembleIncomeData(tsRsp *TushareResponse) []*entity.Income {
	tsData := []*entity.Income{}
	for _, data := range tsRsp.Data.Items {
		body, err := ReflectResponseData(tsRsp.Data.Fields, data)
		if err == nil {
			n := new(entity.Income)
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
