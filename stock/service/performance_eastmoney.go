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

type PerformanceResponse struct {
	SecurityCode       string  `json:"SECURITY_CODE"`
	SecurityNameAbbr   string  `json:"SECURITY_NAME_ABBR"`
	TradeMarketCode    string  `json:"TRADE_MARKET_CODE"`
	TradeMarket        string  `json:"TRADE_MARKET"`
	SecurityTypeCode   string  `json:"SECURITY_TYPE_CODE"`
	SecurityType       string  `json:"SECURITY_TYPE"`
	NewestDate         string  `json:"UPDATE_DATE"`
	ReportDate         string  `json:"REPORTDATE"`
	BasicEps           float64 `json:"BASIC_EPS"`            //每股收益
	DeductBasicEps     float64 `json:"DEDUCT_BASIC_EPS"`     //每股扣非收益
	TotalOperateIncome float64 `json:"TOTAL_OPERATE_INCOME"` //营收
	ParentNetProfit    float64 `json:"PARENT_NETPROFIT"`     //归母净利润
	WeightAvgRoe       float64 `json:"WEIGHTAVG_ROE"`        //净资产收益率
	YoySales           float64 `json:"YSTZ"`                 //应收同比增长
	YoyDeduNp          float64 `json:"SJLTZ"`                //扣非净利润同比增长
	Bps                float64 `json:"BPS"`                  //每股净资产
	Cfps               float64 `json:"MGJYXJJE"`             //每股经营现金流量(元)
	GrossprofitMargin  float64 `json:"XSMLL"`                //销售毛利率(%)
	OrLastMonth        float64 `json:"YSHZ"`                 //营业收入季度环比增长(%)
	NpLastMonth        float64 `json:"SJLHZ"`                //净利润季度环比增长(%)
	AssignDscrpt       string  `json:"ASSIGNDSCRPT"`
	PayYear            string  `json:"PAYYEAR"`
	PublishName        string  `json:"PUBLISHNAME"`
	DividendYieldRatio float64 `json:"ZXGXL"` //股息率
	NoticeDate         string  `json:"NOTICE_DATE"`
	OrgCode            string  `json:"ORG_CODE"`
	TradeMarketZJG     string  `json:"TRADE_MARKET_ZJG"`
	IsNew              string  `json:"ISNEW"`
	QDate              string  `json:"QDATE"`
	NDate              string  `json:"NDATE"`
	DataType           string  `json:"DATATYPE"`
	DataYear           string  `json:"DATAYEAR"`
	DateMmDd           string  `json:"DATEMMDD"`
	EITime             string  `json:"EITIME"`
	SecuCode           string  `json:"SECUCODE"`
}

type PerformanceResponseData struct {
	eastmoney.ReportResponseData
	Data []*PerformanceResponse `json:"data,omitempty"`
}

type PerformanceResponseResult struct {
	eastmoney.ReportResponseResult
	Result *PerformanceResponseData `json:"result,omitempty"`
}

func (svc *PerformanceService) updatePerformance(securityCode string, qDate string, page int) (*PerformanceResponseResult, []interface{}, error) {
	params := eastmoney.CreateRequestParam()
	params.St = "SECURITY_CODE,REPORTDATE"
	params.Sr = "1,1"
	params.P = page
	params.Ps = "500"
	params.Type = "RPT_LICO_FN_CPD"
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
	r := &PerformanceResponseResult{}
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
	for _, pr := range r.Result.Data {
		if pr.NoticeDate != "" {
			pr.NDate = stock.GetQReportDate(pr.NoticeDate)
		} else {
			pr.NDate = pr.QDate
		}
		m := collection.StructToMap(pr, nil)
		p := &entity.Performance{}
		err = collection.MapToStruct(m, p)
		if err == nil {
			go svc.updateShare(p)
			ps = append(ps, p)
		}
	}
	_, err = svc.Insert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return r, ps, err
	}
	return r, ps, err
}

func (svc *PerformanceService) updateShare(performance *entity.Performance) error {
	cond := &entity.Share{TsCode: performance.SecurityCode}
	exist, err := GetShareService().Get(cond, false, "", "")
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
	}
	if !exist {
		share := &entity.Share{
			TsCode:   performance.SecurityCode,
			Name:     performance.SecurityNameAbbr,
			Market:   performance.TradeMarket,
			Symbol:   performance.SecuCode,
			Industry: performance.PublishName,
			Area:     performance.SecurityType,
		}
		_, err = GetShareService().Insert(share)
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
	} else if performance.PublishName == "" {
		cond.TsCode = performance.SecurityCode
		cond.Name = performance.SecurityNameAbbr
		cond.Market = performance.TradeMarket
		cond.Symbol = performance.SecuCode
		cond.Industry = performance.PublishName
		cond.Area = performance.SecurityType
		_, err = GetShareService().Update(cond, nil, "")
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
	}
	return nil
}

func (svc *PerformanceService) RefreshPerformance() error {
	processLog := GetProcessLogService().StartLog("performance", "RefreshPerformance", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, svc.AsyncUpdatePerformance, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetShareCache()
	for _, securityCode := range ts_codes {
		routinePool.Invoke(securityCode)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(processLog, "", "")
	return nil
}

func (svc *PerformanceService) AsyncUpdatePerformance(para interface{}) {
	securityCode := para.(string)
	svc.GetUpdatePerformance(securityCode)
}

func (svc *PerformanceService) GetUpdatePerformance(securityCode string) ([]interface{}, error) {
	//processLog := GetProcessLogService().StartLog("performance", "GetUpdatePerformance", securityCode)
	ps, err := svc.UpdatePerformance(securityCode)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	return ps, err
}

func (svc *PerformanceService) UpdatePerformance(securityCode string) ([]interface{}, error) {
	qdate, _ := svc.findMaxQDate(securityCode)
	result, ps, err := svc.updatePerformance(securityCode, qdate, 1)
	if err != nil {
		//logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	if result == nil || result.Result == nil || result.Result.Pages <= 1 {
		logger.Sugar.Errorf("Error: %s", "result.Result is nil")
		return nil, err
	}
	for i := 2; i <= result.Result.Pages; i++ {
		svc.updatePerformance(securityCode, qdate, i)
	}
	GetQPerformanceService().GetUpdateWmqyQPerformance(securityCode, qdate)

	return ps, nil
}
