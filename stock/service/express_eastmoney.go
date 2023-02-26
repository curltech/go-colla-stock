package service

import (
	"errors"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/collection"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/eastmoney"
	"github.com/curltech/go-colla-stock/stock/entity"
)

type ExpressResponse struct {
	SecurityCode         string  `json:"SECURITY_CODE"`
	SecurityNameAbbr     string  `json:"SECURITY_NAME_ABBR"`
	TradeMarketCode      string  `json:"TRADE_MARKET_CODE"`
	TradeMarket          string  `json:"TRADE_MARKET"`
	SecurityTypeCode     string  `json:"SECURITY_TYPE_CODE"`
	SecurityType         string  `json:"SECURITY_TYPE"`
	NewestDate           string  `json:"UPDATE_DATE"`
	ReportDate           string  `json:"REPORT_DATE"`
	BasicEps             float64 `json:"BASIC_EPS"`               //每股收益
	DeductBasicEps       float64 `json:"DEDUCT_BASIC_EPS"`        //每股扣非收益
	TotalOperateIncome   float64 `json:"TOTAL_OPERATE_INCOME"`    //营收
	TotalOperateIncomeSq float64 `json:"TOTAL_OPERATE_INCOME_SQ"` //去年同期(元)
	ParentNetProfit      float64 `json:"PARENT_NETPROFIT"`        //归母净利润
	ParentNetProfitSq    float64 `json:"PARENT_NETPROFIT_SQ"`     //去年同期(元)
	ParentBvps           float64 `json:"PARENT_BVPS"`             //每股净资产
	WeightAvgRoe         float64 `json:"WEIGHTAVG_ROE"`           //净资产收益率
	YoySales             float64 `json:"YSTZ"`                    //收入同比增长
	YoyNetProfit         float64 `json:"JLRTBZCL"`                //净利润同比增长
	OrLastMonth          float64 `json:"DJDYSHZ"`                 //收入季度环比增长
	NpLastMonth          float64 `json:"DJDJLHZ"`                 //利润季度环比增长
	PublishName          string  `json:"PUBLISHNAME"`
	NoticeDate           string  `json:"NOTICE_DATE"`
	OrgCode              string  `json:"ORG_CODE"`
	Market               string  `json:"MARKET"`
	IsNew                string  `json:"ISNEW"`
	QDate                string  `json:"QDATE"`
	NDate                string  `json:"NDATE"`
	DataType             string  `json:"DATATYPE"`
	DataYear             string  `json:"DATAYEAR"`
	DateMmDd             string  `json:"DATEMMDD"`
	EITime               string  `json:"EITIME"`
	SecuCode             string  `json:"SECUCODE"`
}

type ExpressResponseData struct {
	eastmoney.ReportResponseData
	Data []*ExpressResponse `json:"data,omitempty"`
}

type ExpressResponseResult struct {
	eastmoney.ReportResponseResult
	Result *ExpressResponseData `json:"result,omitempty"`
}

func (this *ExpressService) updateExpress(securityCode string, qDate string, page int) (*ExpressResponseResult, []interface{}, error) {
	params := eastmoney.CreateRequestParam()
	params.St = "SECURITY_CODE,QDATE"
	params.Sr = "1,1"
	params.P = page
	params.Ps = "500"
	params.Type = "RPT_FCI_PERFORMANCEE"
	if securityCode != "" {
		params.Filter = "(SECURITY_CODE%3D%22" + securityCode + "%22)"
	}
	if qDate != "" {
		params.Filter = params.Filter + "(QDATE%3E%22" + qDate + "%22)"
	}
	resp, err := eastmoney.ReportFastGet(*params)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, nil, err
	}
	r := &ExpressResponseResult{}
	err = json.Unmarshal(resp, r)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, nil, err
	}
	if !r.Success {
		logger.Sugar.Errorf("Error: %s", r.Message)
		return nil, nil, errors.New(r.Message)
	}
	if r.Result == nil || r.Result.Data == nil {
		logger.Sugar.Errorf("Error: %s", "r.Result.Data is nil")
		return nil, nil, errors.New("r.Result.Data is nil")
	}
	ps := make([]interface{}, 0)
	for _, er := range r.Result.Data {
		if er.NoticeDate != "" {
			er.NDate = stock.GetQReportDate(er.NoticeDate)
		} else {
			er.NDate = er.QDate
		}
		m := collection.StructToMap(er, nil)
		p := &entity.Express{}
		err = collection.MapToStruct(m, p)
		if err == nil {
			ps = append(ps, p)
		}
	}
	_, err = this.Insert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return r, ps, err
	}

	return r, ps, err
}

func (this *ExpressService) RefreshExpress() error {
	processLog := GetProcessLogService().StartLog("express", "RefreshExpress", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, this.AsyncUpdateExpress, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	for _, securityCode := range ts_codes {
		routinePool.Invoke(securityCode)
	}
	routinePool.Wait(nil)
	GetShareService().RefreshCacheShare()
	GetProcessLogService().EndLog(processLog, "", "")
	return nil
}

func (this *ExpressService) AsyncUpdateExpress(para interface{}) {
	securityCode := para.(string)
	this.GetUpdateExpress(securityCode)
}

func (this *ExpressService) GetUpdateExpress(securityCode string) ([]interface{}, error) {
	//processLog := GetProcessLogService().StartLog("express", "GetUpdateExpress", securityCode)
	ps, err := this.UpdateExpress(securityCode)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	return ps, err
}

func (this *ExpressService) UpdateExpress(securityCode string) ([]interface{}, error) {
	qdate, _ := this.findMaxQDate(securityCode)
	result, ps, err := this.updateExpress(securityCode, qdate, 1)
	if err != nil {
		//logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	if result == nil || result.Result == nil || result.Result.Pages <= 1 {
		return nil, err
	}
	for i := 2; i <= result.Result.Pages; i++ {
		this.updateExpress(securityCode, qdate, i)
	}
	GetQPerformanceService().GetUpdateWmqyQPerformance(securityCode, qdate)

	return ps, nil
}

func (this *ExpressService) updateShare(express *entity.Express) error {
	cond := &entity.Share{TsCode: express.SecurityCode}
	exist, err := this.Get(cond, false, "", "")
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
	}
	if !exist {
		share := &entity.Share{
			TsCode:   express.SecurityCode,
			Name:     express.SecurityNameAbbr,
			Market:   express.TradeMarket,
			Symbol:   express.SecuCode,
			Industry: express.PublishName,
			Area:     express.SecurityType,
		}
		_, err = GetShareService().Insert(share)
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
	} else if express.PublishName == "" {
		cond.TsCode = express.SecurityCode
		cond.Name = express.SecurityNameAbbr
		cond.Market = express.TradeMarket
		cond.Symbol = express.SecuCode
		cond.Industry = express.PublishName
		cond.Area = express.SecurityType
		_, err = GetShareService().Update(cond, nil, "")
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
	}
	return nil
}
