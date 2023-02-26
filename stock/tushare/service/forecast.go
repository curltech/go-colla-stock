package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

type ForecastRequest struct {
	TsCode    string `json:"ts_code,omitempty"`    // str	Y	股票代码
	AnnDate   string `json:"ann_date,omitempty"`   // str	N	公告日期
	StartDate string `json:"start_date,omitempty"` // str	N	公告开始日期
	EndDate   string `json:"end_date,omitempty"`   // str	N	公告结束日期
	Period    string `json:"period,omitempty"`     // str	N	报告期(每个季度最后一天的日期,比如20171231表示年报)
	Type      string `json:"type,omitempty"`       // str	N	预告类型(预增/预减/扭亏/首亏/续亏/续盈/略增/略减)
}

// 获取业绩预告数据,用户需要至少800积分才可以调取,具体请参阅积分获取办法 https://tushare.pro/document/1?doc_id=13
func (ts *ForecastService) Forecast(params ForecastRequest) (tsRsp *TushareResponse, err error) {
	resp, err := FastPost(&TushareRequest{
		ApiName: "forecast",
		Token:   token,
		Params:  params,
		Fields:  reflectFields(entity.Forecast{}),
	})
	return resp, err
}

func assembleForecast(tsRsp *TushareResponse) []*entity.Forecast {
	tsData := []*entity.Forecast{}
	for _, data := range tsRsp.Data.Items {
		body, err := ReflectResponseData(tsRsp.Data.Fields, data)
		if err == nil {
			n := new(entity.Forecast)
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
type ForecastService struct {
	service.OrmBaseService
}

var forecastService = &ForecastService{}

func GetForecastService() *ForecastService {
	return forecastService
}

func (this *ForecastService) GetSeqName() string {
	return seqname
}

func (this *ForecastService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Forecast{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ForecastService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Forecast, 0)
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
	service.GetSession().Sync(new(entity.Forecast))
	forecastService.OrmBaseService.GetSeqName = forecastService.GetSeqName
	forecastService.OrmBaseService.FactNewEntity = forecastService.NewEntity
	forecastService.OrmBaseService.FactNewEntities = forecastService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("cashflow", forecastService)
}
