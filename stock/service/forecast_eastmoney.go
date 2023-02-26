package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/collection"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/eastmoney"
	"github.com/curltech/go-colla-stock/stock/entity"
	"net/url"
	"strings"
)

type ForecastRequestParam struct {
	Callback    string `json:"callback"`
	SortColumns string `json:"sortColumns"` //排序的字段，逗号分隔
	SortTypes   string `json:"sortTypes"`   //-1降序，1升序，逗号分隔
	PageSize    string `json:"pageSize"`    //每页的记录数
	PageNumber  int    `json:"pageNumber"`  //页数
	ReportName  string `json:"reportName"`  //获取的数据类型
	Columns     string `json:"columns"`     //获取的字段，ALL
	Token       string `json:"token"`
	Filter      string `json:"filter,omitempty"` //条件
}

type ForecastResponse struct {
	SecuCode            string  `json:"SECUCODE"`
	SecurityCode        string  `json:"SECURITY_CODE"`
	SecurityNameAbbr    string  `json:"SECURITY_NAME_ABBR"`
	TradeMarketCode     string  `json:"TRADE_MARKET_CODE"`
	TradeMarket         string  `json:"TRADE_MARKET"`
	SecurityTypeCode    string  `json:"SECURITY_TYPE_CODE"`
	SecurityType        string  `json:"SECURITY_TYPE"`
	NoticeDate          string  `json:"NOTICE_DATE"`
	OrgCode             string  `json:"ORG_CODE"`
	ReportDate          string  `json:"REPORT_DATE"`
	QDate               string  `json:"QDATE"`
	NDate               string  `json:"NDATE"`
	PredictFinanceCode  string  `json:"PREDICT_FINANCE_CODE"`
	PredictFinance      string  `xorm:"varchar(1024)" json:"PREDICT_FINANCE"`
	PredictAmtLower     float64 `json:"PREDICT_AMT_LOWER"`
	PredictAmtUpper     float64 `json:"PREDICT_AMT_UPPER"`
	AddAmpLower         float64 `json:"ADD_AMP_LOWER"`
	AddAmpUpper         float64 `json:"ADD_AMP_UPPER"`
	PredictContent      string  `xorm:"varchar(32000)" json:"PREDICT_CONTENT"`
	ChangeReasonExplain string  `xorm:"varchar(32000)" json:"CHANGE_REASON_EXPLAIN"`
	PredictType         string  `json:"PREDICT_TYPE"`
	PreYearSamePeriod   float64 `json:"PREYEAR_SAME_PERIOD"`
	IncreaseAvg         float64 `json:"INCREASE_JZ"`
	ForecastAvg         float64 `json:"FORECAST_JZ"`
	ForecastState       string  `json:"FORECAST_STATE"`
	IsLatest            string  `json:"IS_LATEST"`
}

type ForecastResponseData struct {
	eastmoney.ReportResponseData
	Data []*ForecastResponse `json:"data,omitempty"`
}

type ForecastResponseResult struct {
	eastmoney.ReportResponseResult
	Result *ForecastResponseData `json:"result,omitempty"`
}

func ForecastFastGet(requestParam ForecastRequestParam) ([]byte, error) {
	resp, err := eastmoney.FastGet("http://datacenter-web.eastmoney.com/securities/api/data/v1/get", requestParam)
	if err != nil {
		fmt.Println("Get fail:", err.Error())
		return nil, err
	}
	respStr := string(resp)
	respStr = strings.TrimPrefix(respStr, requestParam.Callback)
	respStr = strings.TrimPrefix(respStr, "(")
	respStr = strings.TrimSuffix(respStr, ");")
	resp = []byte(respStr)

	return resp, nil
}

func (this *ForecastService) updateForecast(securityCode string, reportDate string, page int) (*ForecastResponseResult, []interface{}, error) {
	params := &ForecastRequestParam{
		Callback:   "jQuery112308629035897216608_1638369560260",
		Token:      "894050c76af8597a853f5b408b759f5d",
		ReportName: "RPT_PUBLIC_OP_NEWPREDICT",
		Columns:    "ALL",
	}
	params.SortColumns = "SECURITY_CODE,REPORT_DATE"
	params.SortTypes = "1,1"
	params.PageNumber = page
	params.PageSize = "500"
	if securityCode != "" {
		params.Filter = "(SECURITY_CODE%3D%22" + securityCode + "%22)"
	}
	if reportDate != "" {
		params.Filter = params.Filter + "(REPORT_DATE%3E%27" + url.QueryEscape(reportDate) + "%27)"
		//params.Filter = params.Filter + url.QueryEscape("(REPORT_DATE>'"+reportDate+"')")
	}
	resp, err := ForecastFastGet(*params)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, nil, err
	}
	r := &ForecastResponseResult{}
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
	for _, fr := range r.Result.Data {
		fr.QDate = stock.GetQReportDate(fr.ReportDate)
		if fr.NoticeDate != "" {
			fr.NDate = stock.GetQReportDate(fr.NoticeDate)
		} else {
			fr.NDate = fr.QDate
		}
		m := collection.StructToMap(fr, nil)
		p := &entity.Forecast{}
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

func (this *ForecastService) RefreshForecast() error {
	processLog := GetProcessLogService().StartLog("forecast", "RefreshForecast", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, this.AsyncUpdateForecast, nil)
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

func (this *ForecastService) AsyncUpdateForecast(para interface{}) {
	securityCode := para.(string)
	this.GetUpdateForecast(securityCode)
}

func (this *ForecastService) GetUpdateForecast(securityCode string) ([]interface{}, error) {
	//processLog := GetProcessLogService().StartLog("forecast", "GetUpdateForecast", securityCode)
	ps, err := this.UpdateForecast(securityCode)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	return ps, err
}

func (this *ForecastService) UpdateForecast(securityCode string) ([]interface{}, error) {
	reportDate, _ := this.findMaxReportDate(securityCode)
	qdate, _ := this.findMaxQDate(securityCode)
	result, ps, err := this.updateForecast(securityCode, reportDate, 1)
	if err != nil {
		//logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	if result == nil || result.Result == nil || result.Result.Pages <= 1 {
		return nil, err
	}
	for i := 2; i <= result.Result.Pages; i++ {
		this.updateForecast(securityCode, reportDate, i)
	}
	GetQPerformanceService().GetUpdateWmqyQPerformance(securityCode, qdate)

	return ps, nil
}
