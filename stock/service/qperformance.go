package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	entity "github.com/curltech/go-colla-stock/stock/entity"
	"math"
	"strings"
)

type QPerformanceService struct {
	service.OrmBaseService
}

var qperformanceService = &QPerformanceService{}

func GetQPerformanceService() *QPerformanceService {
	return qperformanceService
}

func (this *QPerformanceService) GetSeqName() string {
	return seqname
}

func (this *QPerformanceService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.QPerformance{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *QPerformanceService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.QPerformance, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *QPerformanceService) Search(keyword string, tscode string, terms []int, sourceOptions []string, startDate string, endDate string, orderby string, from int, limit int, count int64) ([]*entity.QPerformance, int64, error) {
	if keyword == "" && tscode == "" {
		if limit == 0 {
			limit = 20
		}
	}
	//计算StartDate，前提是terms有值
	if startDate == "" && terms != nil && len(terms) > 0 {
		term := 0
		if terms != nil {
			for _, t := range terms {
				if t == 0 {
					term = 0
					break
				} else {
					if t > term {
						term = t
					}
				}
			}
		}
		tscodes := tscode
		if tscodes == "" && keyword != "" {
			shares, err := GetShareService().Search(keyword, 0, 10)
			if err != nil {
				return nil, 0, err
			}
			i := 0
			for _, share := range shares {
				if i != 0 {
					tscodes += ","
				}
				tscodes += share.TsCode
			}
		}
		qterm, err := this.GetQTermBySql(tscodes, term)
		if err == nil {
			startDate = qterm.StartDate
		}
	}
	paras := make([]interface{}, 0)
	conds := "1=1"
	if startDate != "" {
		conds += " and qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		conds += " and qdate<=?"
		paras = append(paras, endDate)
	}
	if sourceOptions != nil && len(sourceOptions) > 0 {
		sourceConds, sourceParas := stock.InBuildStr("source", strings.Join(sourceOptions, ","), ",")
		conds += " and " + sourceConds
		paras = append(paras, sourceParas...)
	}
	if tscode != "" {
		tscodeCond, tscodeParams := stock.InBuildStr("tscode", tscode, ",")
		conds += " and " + tscodeCond
		paras = append(paras, tscodeParams...)
	} else {
		if keyword != "" {
			conds += " and tscode in (select tscode from stk_share where name like ? or tscode like ? or pinyin like ?)"
			paras = append(paras, keyword+"%")
			paras = append(paras, keyword+"%")
			paras = append(paras, strings.ToLower(keyword)+"%")
		}
	}
	qperformances := make([]*entity.QPerformance, 0)
	var err error
	condiBean := &entity.QPerformance{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "tscode,qdate desc,tradedate desc"
	}
	err = this.Find(&qperformances, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return qperformances, count, nil
}

type PerformanceType int64

const (
	PerformanceType_Performance PerformanceType = iota + 1
	PerformanceType_Express
	PerformanceType_Forecast
)

type StdType int64

const (
	StdType_None StdType = iota + 1
	StdType_Std
	StdType_MinMax
)

type LineType int64

const (
	LineType_Day LineType = iota + 1
	LineType_Wmqy
)

/**
组合当前价格和股票季度正式业绩数据
*/
func (this *QPerformanceService) FindQPerformance(lineType LineType, ts_code string, startDate string, endDate string) (map[string][]interface{}, error) {
	lineName := "stk_dayline"
	if lineType == LineType_Wmqy {
		lineName = "stk_wmqyline"
	}
	sql := "select w.tscode as ts_code,sp.securitynameabbr as security_name" +
		",sp.qdate as qdate,sp.ndate as ndate,w.tradedate as trade_date,'performance' as source" +
		",case when sp.BasicEps!=0 then" +
		" case right(sp.qdate,1)" +
		" when '4' then w.close/sp.BasicEps" +
		" when '3' then w.close/(sp.BasicEps/3*4)" +
		" when '2' then w.close/(sp.BasicEps*2)" +
		" when '1' then w.close/(sp.BasicEps*4)" +
		" end" +
		" else 0 end as pe" +
		",case when w.turnover!=0 then w.vol*10000/w.turnover else 0 end as share_number" +
		",case when w.turnover!=0 then w.vol*10000*w.close/w.turnover else 0 end as market_value" +
		",w.high as high,w.pctchghigh as pct_chg_high,w.close as close,w.pctchgclose as pct_chg_close" +
		",sp.WeightAvgRoe as weight_avg_roe,sp.GrossProfitMargin as gross_profit_margin" +
		",sp.totaloperateincome as total_operate_income,sp.ParentNetProfit as parent_net_profit,sp.BasicEps as basic_eps,sp.OrLastMonth as or_last_month" +
		",sp.NpLastMonth as np_last_month,sp.YoySales as yoy_sales,sp.YoyDeduNp as yoy_dedu_np" +
		",sp.Cfps as cfps,sp.DividendYieldRatio as dividend_yield_ratio" +
		" from " + lineName + " w join stk_performance sp on w.tscode=sp.securitycode"
	conds, tscodeParas := stock.InBuildStr("w.tscode", ts_code, ",")
	paras := make([]interface{}, 0)
	if lineType == LineType_Wmqy {
		sql = sql + " where w.qdate = sp.qdate"
	} else {
		sql = sql + " where w.tradedate = ?"
		tscode := "000001"
		if ts_code != "" {
			tscode = tscodeParas[0].(string)
		}
		dayLines, err := GetDayLineService().findMaxTradeDate(tscode)
		if err != nil {
			return nil, err
		}
		if dayLines == nil || len(dayLines) == 0 || dayLines[0] == nil {
			return nil, errors.New("NoTradeDate")
		}
		paras = append(paras, dayLines[0].TradeDate)
	}
	if ts_code != "" {
		sql = sql + " and " + conds
		paras = append(paras, tscodeParas...)
	}
	if startDate != "" {
		sql = sql + " and sp.qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		sql = sql + " and sp.qdate<=?"
		paras = append(paras, endDate)
	}
	sql = sql + " order by w.tscode,sp.qdate"
	results, err := GetDayLineService().Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	jsonMap, _, _ := stock.GetJsonMap(entity.QPerformance{})
	qpMap := make(map[string][]interface{}, 0)
	for _, result := range results {
		qp := &entity.QPerformance{}
		qp.LineType = int64(lineType)
		for colname, v := range result {
			err = reflect.Set(qp, jsonMap[colname], string(v))
			if err != nil {
				logger.Sugar.Errorf("Set colname %v value %v error", colname, string(v))
			}
		}
		ps, ok := qpMap[qp.TsCode]
		if !ok {
			ps = make([]interface{}, 0)
		}
		ps = append(ps, qp)
		qpMap[qp.TsCode] = ps
	}

	return qpMap, nil
}

/**
对股票季度业绩数据进行异常值处理和标准化操作
*/
func (this *QPerformanceService) Std(ps []interface{}, stdType StdType, isWinsorize bool) []*entity.QPerformance {
	qps := make([]*entity.QPerformance, 0)
	if stdType == 0 || stdType == StdType_None {
		for _, qp := range ps {
			qps = append(qps, qp.(*entity.QPerformance))
		}
		return qps
	}
	_, _, heads := stock.GetJsonMap(entity.QPerformance{})
	std_colnames := heads[12:]
	reserved_colnames := heads[0:12]
	stat := stock.CreateStat(ps, std_colnames)
	stds, minmaxs := stat.CalStd(reserved_colnames, isWinsorize)
	if stdType == StdType_MinMax {
		for _, qp := range minmaxs {
			qps = append(qps, qp.(*entity.QPerformance))
		}
	} else if stdType == StdType_Std {
		for _, qp := range stds {
			qps = append(qps, qp.(*entity.QPerformance))
		}
	}

	return qps
}

/**
组合当前价格和股票季度业绩、预测和快报数据，形成完整的股票季度业绩数据，进行标准化操作
*/
func (this *QPerformanceService) findWmqyQPerformance(ts_code string, startDate string, endDate string, stdType StdType, isWinsorize bool) (map[string][]*entity.QPerformance, error) {
	previousMap := make(map[string]*entity.QPerformance, 0)
	if startDate == "" {
		qp, _ := this.findMaxQDate(ts_code, int64(LineType_Wmqy))
		if qp != nil {
			startDate = qp.QDate
		}
	}
	this.deleteQPerformance(ts_code, LineType_Wmqy, startDate)
	qp, _ := this.findMaxQDate(ts_code, int64(LineType_Wmqy))
	if qp != nil {
		previousMap[ts_code] = qp
	}
	qpMap, err := this.FindQPerformance(LineType_Wmqy, ts_code, startDate, endDate)
	if err != nil {
		return nil, err
	}
	qeMap, err := this.FindQExpress(LineType_Wmqy, ts_code, startDate, endDate)
	if err != nil {
		return nil, err
	}
	this.merge(qeMap, qpMap)

	qfMap, err := this.FindQForecast(LineType_Wmqy, ts_code, startDate, endDate)
	if err != nil {
		return nil, err
	}
	this.merge(qfMap, qpMap)

	this.Compute(qpMap, previousMap)

	stds := this.StdMap(qpMap, stdType, isWinsorize)

	return stds, nil
}

func (this *QPerformanceService) findDayQPerformance(ts_code string, stdType StdType, isWinsorize bool) (map[string][]*entity.QPerformance, error) {
	var dqMap map[string][]interface{}
	var err error
	performanceQDate, err := GetPerformanceService().findMaxQDate(ts_code)
	if err != nil {
		logger.Sugar.Errorf("Performance findMaxQDate error:%v", err.Error())
	}
	expressQDate, err := GetExpressService().findMaxQDate(ts_code)
	if err != nil {
		logger.Sugar.Errorf("Express findMaxQDate error:%v", err.Error())
	}
	var startDate string
	if performanceQDate >= expressQDate {
		startDate = performanceQDate
	} else {
		startDate = expressQDate
	}
	forecastQDate, err := GetForecastService().findMaxQDate(ts_code)
	if err != nil {
		logger.Sugar.Errorf("Forecast findMaxQDate error:%v", err.Error())
	}
	if forecastQDate >= startDate {
		startDate = forecastQDate
	}
	endDate := startDate
	if performanceQDate == startDate {
		dqMap, err = this.FindQPerformance(LineType_Day, ts_code, startDate, endDate)
		if err != nil {
			return nil, err
		}
	} else if expressQDate == startDate {
		dqMap, err = this.FindQExpress(LineType_Day, ts_code, startDate, endDate)
		if err != nil {
			return nil, err
		}
	} else if forecastQDate == startDate {
		dqMap, err = this.FindQForecast(LineType_Day, ts_code, startDate, endDate)
		if err != nil {
			return nil, err
		}

	}
	previousMap := make(map[string]*entity.QPerformance, 0)
	qp, err := this.findMaxQDate(ts_code, 0)
	if qp != nil {
		previousMap[ts_code] = qp
	}
	this.Compute(dqMap, previousMap)

	stds := this.StdMap(dqMap, stdType, isWinsorize)

	return stds, nil
}

/**
查询股票季度业绩数据，并进行标准化处理
*/
func (this *QPerformanceService) FindStdQPerformance(ts_code string, terms []int, startDate string, endDate string, stdType StdType, isWinsorize bool) (map[string][]*entity.QPerformance, error) {
	qps, _, err := this.Search("", ts_code, terms, nil, startDate, endDate, "", 0, 0, 1)
	if err != nil {
		return nil, err
	}
	qpMap := make(map[string][]interface{}, 0)
	for _, qp := range qps {
		ps, ok := qpMap[qp.TsCode]
		if !ok {
			ps = make([]interface{}, 0)
		}
		ps = append(ps, qp)
		qpMap[qp.TsCode] = ps
	}

	stds := this.StdMap(qpMap, stdType, isWinsorize)

	return stds, nil
}

/**
批量标准化股票季度业绩数据
*/
func (this *QPerformanceService) StdMap(psMap map[string][]interface{}, stdType StdType, isWinsorize bool) map[string][]*entity.QPerformance {
	stds := make(map[string][]*entity.QPerformance)
	for tscode, ps := range psMap {
		std := this.Std(ps, stdType, isWinsorize)
		stds[tscode] = std
	}

	return stds
}

/**
合并当前价格和股票季度业绩、预测、快报数据，按照时间排序
*/
func (this *QPerformanceService) merge(qeMap map[string][]interface{}, qpMap map[string][]interface{}) {
	for ts_code, qes := range qeMap {
		qps, ok := qpMap[ts_code]
		qdate := ""
		var tradeDate int64
		if ok && len(qps) > 0 {
			qp := qps[len(qps)-1].(*entity.QPerformance)
			qdate = qp.QDate
			tradeDate = qp.TradeDate
		}
		if len(qes) > 0 {
			for i := len(qes) - 1; i >= 0; i-- {
				qe := qes[i].(*entity.QPerformance)
				if qe.QDate > qdate {
					qps = append(qps, qe)
					qpMap[ts_code] = qps
				} else if qe.QDate == qdate {
					if qe.TradeDate > tradeDate {
						qps = append(qps, qe)
						qpMap[ts_code] = qps
					} else {
						break
					}
				} else {
					break
				}
			}
		}
	}
}

/**
计算并推导缺失的当前价格和股票季度业绩数据项
*/
func (this *QPerformanceService) Compute(qpMap map[string][]interface{}, previousMap map[string]*entity.QPerformance) {
	_, tscodeMap := GetShareService().GetCacheShare()
	for tscode, qps := range qpMap {
		var previous *entity.QPerformance
		if previousMap != nil {
			previous, _ = previousMap[tscode]
		}
		for _, q := range qps {
			qp := q.(*entity.QPerformance)
			if stock.Equal(qp.BasicEps, 0.0) && !stock.Equal(qp.ShareNumber, 0.0) && !stock.Equal(qp.ParentNetProfit, 0.0) {
				qp.BasicEps = qp.ParentNetProfit / qp.ShareNumber
				if stock.Equal(qp.Pe, 0.0) {
					if strings.HasSuffix(qp.QDate, "Q4") {
						qp.Pe = qp.Close / qp.BasicEps
					} else if strings.HasSuffix(qp.QDate, "Q3") {
						qp.Pe = qp.Close / (qp.BasicEps * 4 / 3)
					} else if strings.HasSuffix(qp.QDate, "Q2") {
						qp.Pe = qp.Close / (qp.BasicEps * 2)
					} else if strings.HasSuffix(qp.QDate, "Q1") {
						qp.Pe = qp.Close / (qp.BasicEps * 4)
					}
				}
			}
			if !stock.Equal(qp.ParentNetProfit, 0.0) {
				if strings.HasSuffix(qp.QDate, "Q4") {
					qp.YearNetProfit = qp.ParentNetProfit
				} else if strings.HasSuffix(qp.QDate, "Q3") {
					qp.YearNetProfit = qp.ParentNetProfit * 4 / 3
				} else if strings.HasSuffix(qp.QDate, "Q2") {
					qp.YearNetProfit = qp.ParentNetProfit * 2
				} else if strings.HasSuffix(qp.QDate, "Q1") {
					qp.YearNetProfit = qp.ParentNetProfit * 4
				}
			}
			if !stock.Equal(qp.TotalOperateIncome, 0.0) {
				if strings.HasSuffix(qp.QDate, "Q4") {
					qp.YearOperateIncome = qp.TotalOperateIncome
				} else if strings.HasSuffix(qp.QDate, "Q3") {
					qp.YearOperateIncome = qp.TotalOperateIncome * 4 / 3
				} else if strings.HasSuffix(qp.QDate, "Q2") {
					qp.YearOperateIncome = qp.TotalOperateIncome * 2
				} else if strings.HasSuffix(qp.QDate, "Q1") {
					qp.YearOperateIncome = qp.TotalOperateIncome * 4
				}
			}
			if !stock.Equal(qp.WeightAvgRoe, 0.0) {
				if strings.HasSuffix(qp.QDate, "Q4") {
					qp.WeightAvgRoe = qp.WeightAvgRoe
				} else if strings.HasSuffix(qp.QDate, "Q3") {
					qp.WeightAvgRoe = qp.WeightAvgRoe * 4 / 3
				} else if strings.HasSuffix(qp.QDate, "Q2") {
					qp.WeightAvgRoe = qp.WeightAvgRoe * 2
				} else if strings.HasSuffix(qp.QDate, "Q1") {
					qp.WeightAvgRoe = qp.WeightAvgRoe * 4
				}
			}
			if previous != nil && !stock.Equal(previous.MarketValue, 0.0) {
				qp.PctChgMarketValue = (qp.MarketValue - previous.MarketValue) / previous.MarketValue
			}
			if previous != nil && stock.Equal(qp.WeightAvgRoe, 0.0) && !stock.Equal(previous.YearNetProfit, 0) {
				qp.WeightAvgRoe = (qp.YearNetProfit * previous.WeightAvgRoe) / previous.YearNetProfit
			}
			if previous != nil && stock.Equal(qp.YoySales, 0.0) {
				qp.YoySales = previous.YoySales
			}
			qp.Peg = qp.Pe / (10 * math.Pow(1+qp.YoyDeduNp/100, 5))

			if previous != nil && !stock.Equal(previous.TotalOperateIncome, 0) && stock.Equal(qp.OrLastMonth, 0.0) {
				if previous.QDate < qp.QDate {
					qp.OrLastMonth = 100 * (qp.TotalOperateIncome - previous.TotalOperateIncome) / previous.TotalOperateIncome
				} else if previous.QDate == qp.QDate {
					qp.OrLastMonth = previous.OrLastMonth
				}
			}
			if previous != nil && !stock.Equal(previous.TotalOperateIncome, 0) && stock.Equal(qp.NpLastMonth, 0.0) {
				if previous.QDate < qp.QDate {
					qp.NpLastMonth = 100 * (qp.ParentNetProfit - previous.ParentNetProfit) / previous.ParentNetProfit
				} else if previous.QDate == qp.QDate {
					qp.NpLastMonth = previous.NpLastMonth
				}
			}
			if previous != nil && stock.Equal(qp.GrossprofitMargin, 0.0) {
				qp.GrossprofitMargin = previous.GrossprofitMargin
			}
			if qp.Industry == "" {
				share, ok := tscodeMap[qp.TsCode]
				if ok {
					qp.Industry = share.Industry
				}
			}
			if qp.Sector == "" {
				share, ok := tscodeMap[qp.TsCode]
				if ok {
					qp.Sector = share.Sector
				}
			}

			previous = qp
		}
	}
}

/**
组合当前价格和股票季度业绩快报数据
*/
func (this *QPerformanceService) FindQExpress(lineType LineType, ts_code string, startDate string, endDate string) (map[string][]interface{}, error) {
	lineName := "stk_dayline"
	if lineType == LineType_Wmqy {
		lineName = "stk_wmqyline"
	}
	sql := "select w.tscode as ts_code,sp.securitynameabbr as security_name" +
		",sp.qdate as qdate,sp.ndate as ndate,w.tradedate as trade_date,'express' as source" +
		",case when sp.BasicEps!=0 then" +
		" case right(sp.qdate,1)" +
		" when '4' then w.close/sp.BasicEps" +
		" when '3' then w.close/(sp.BasicEps/3*4)" +
		" when '2' then w.close/(sp.BasicEps*2)" +
		" when '1' then w.close/(sp.BasicEps*4)" +
		" end" +
		" else 0 end as pe" +
		",case when w.turnover!=0 then w.vol*10000/w.turnover else 0 end as share_number" +
		",case when w.turnover!=0 then w.vol*10000*w.close/w.turnover else 0 end as market_value" +
		",w.high as high,w.pctchghigh as pct_chg_high,w.close as close,w.pctchgclose as pct_chg_close" +
		",sp.WeightAvgRoe as weight_avg_roe,sp.totaloperateincome as total_operate_income" +
		",case when sp.totaloperateincome != 0 then sp.parentnetprofit / sp.totaloperateincome end as gross_profit_margin" +
		",sp.ParentNetProfit as parent_net_profit,sp.BasicEps as basic_eps,sp.OrLastMonth as or_last_month" +
		",sp.NpLastMonth as np_last_month,sp.YoySales as yoy_sales,sp.YoyNetProfit as yoy_dedu_np" +
		" from " + lineName + " w join stk_express sp on w.tscode=sp.securitycode"
	conds, tscodeParas := stock.InBuildStr("w.tscode", ts_code, ",")
	paras := make([]interface{}, 0)
	if lineType == LineType_Wmqy {
		sql = sql + " where w.qdate = sp.qdate"
	} else {
		sql = sql + " where w.tradedate = ?"
		tscode := "000001"
		if ts_code != "" {
			tscode = tscodeParas[0].(string)
		}
		dayLines, err := GetDayLineService().findMaxTradeDate(tscode)
		if err != nil {
			return nil, err
		}
		if dayLines == nil || len(dayLines) == 0 || dayLines[0] == nil {
			return nil, errors.New("NoTradeDate")
		}
		paras = append(paras, dayLines[0].TradeDate)
	}
	if ts_code != "" {
		sql = sql + " and " + conds
		paras = append(paras, tscodeParas...)
	}
	if startDate != "" {
		sql = sql + " and sp.qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		sql = sql + " and sp.qdate<=?"
		paras = append(paras, endDate)
	}
	sql = sql + " order by w.tscode,sp.qdate"
	results, err := GetDayLineService().Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	jsonMap, _, _ := stock.GetJsonMap(entity.QPerformance{})
	qpMap := make(map[string][]interface{}, 0)
	for _, result := range results {
		qp := &entity.QPerformance{}
		qp.LineType = int64(lineType)
		for colname, v := range result {
			err = reflect.Set(qp, jsonMap[colname], string(v))
			if err != nil {
				//logger.Sugar.Errorf("Set colname %v value %v error", colname, string(v))
			}
		}
		ps, ok := qpMap[qp.TsCode]
		if !ok {
			ps = make([]interface{}, 0)
		}
		ps = append(ps, qp)
		qpMap[qp.TsCode] = ps
	}

	return qpMap, nil
}

/**
组合当前价格和股票季度业绩预测数据
*/
func (this *QPerformanceService) FindQForecast(lineType LineType, ts_code string, startDate string, endDate string) (map[string][]interface{}, error) {
	lineName := "stk_dayline"
	if lineType == LineType_Wmqy {
		lineName = "stk_wmqyline"
	}
	forecastView := "(select securitycode,securitynameabbr,qdate,max(ndate) as ndate" +
		",sum(case predictfinancecode when '003' then forecastavg end) as BasicEps" +
		",sum(case predictfinancecode when '004' then forecastavg end) as ParentNetProfit" +
		",sum(case predictfinancecode when '004' then increaseavg end) as YoyDeduNp" +
		",sum(case predictfinancecode when '006' then forecastavg end) as TotalOperateIncome" +
		",sum(case predictfinancecode when '006' then increaseavg end) as YoySales" +
		" from public.stk_forecast" +
		" group by securitycode,securitynameabbr,qdate)"
	sql := "select w.tscode as ts_code,sp.securitynameabbr as security_name" +
		",sp.qdate as qdate,sp.ndate as ndate,w.tradedate as trade_date,'forecast' as source" +
		",case when sp.BasicEps!=0 then" +
		" case right(sp.qdate,1)" +
		" when '4' then w.close/sp.BasicEps" +
		" when '3' then w.close/(sp.BasicEps/3*4)" +
		" when '2' then w.close/(sp.BasicEps*2)" +
		" when '1' then w.close/(sp.BasicEps*4)" +
		" end" +
		" else 0 end as pe" +
		",case when w.turnover!=0 then w.vol*10000/w.turnover else 0 end as share_number" +
		",case when w.turnover!=0 then w.vol*10000*w.close/w.turnover else 0 end as market_value" +
		",w.high as high,w.pctchghigh as pct_chg_high,w.close as close,w.pctchgclose as pct_chg_close" +
		",sp.TotalOperateIncome as total_operate_income,sp.ParentNetProfit as parent_net_profit,sp.BasicEps as basic_eps" +
		",sp.YoyDeduNp as yoy_dedu_np,sp.YoySales as yoy_sales" +
		" from " + lineName + " w join "
	sql = sql + forecastView + " sp on w.tscode=sp.securitycode"
	conds, tscodeParas := stock.InBuildStr("w.tscode", ts_code, ",")
	paras := make([]interface{}, 0)
	if lineType == LineType_Wmqy {
		sql = sql + " where w.qdate = sp.qdate"
	} else {
		sql = sql + " where w.tradedate = ?"
		tscode := "000001"
		if ts_code != "" {
			tscode = (tscodeParas[0]).(string)
		}
		dayLines, err := GetDayLineService().findMaxTradeDate(tscode)
		if err != nil {
			return nil, err
		}
		if dayLines == nil || len(dayLines) == 0 || dayLines[0] == nil {
			return nil, errors.New("NoTradeDate")
		}
		paras = append(paras, dayLines[0].TradeDate)
	}
	if ts_code != "" {
		sql = sql + " and " + conds
		paras = append(paras, tscodeParas...)
	}
	if startDate != "" {
		sql = sql + " and sp.qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		sql = sql + " and sp.qdate<=?"
		paras = append(paras, endDate)
	}
	sql = sql + " order by w.tscode,sp.qdate"
	results, err := GetDayLineService().Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	jsonMap, _, _ := stock.GetJsonMap(entity.QPerformance{})
	qpMap := make(map[string][]interface{}, 0)
	for _, result := range results {
		qp := &entity.QPerformance{}
		qp.LineType = int64(lineType)
		for colname, v := range result {
			if len(v) == 0 {
				logger.Sugar.Debugf("Set colname %v no value", colname)
				continue
			}
			err = reflect.Set(qp, jsonMap[colname], string(v))
			if err != nil {
				logger.Sugar.Errorf("Set colname %v value %v error", colname, string(v))
			}
		}

		ps, ok := qpMap[qp.TsCode]
		if !ok {
			ps = make([]interface{}, 0)
		}
		ps = append(ps, qp)
		qpMap[qp.TsCode] = ps
	}

	return qpMap, nil
}

/**
刷新所有股票的季度业绩数据
*/
func (this *QPerformanceService) RefreshWmqyQPerformance(startDate string) error {
	processLog := GetProcessLogService().StartLog("qperformance", "RefreshWmqyQPerformance", "")
	routinePool := thread.CreateRoutinePool(10, this.AsyncUpdateWmqyQPerformance, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	for _, ts_code := range ts_codes {
		para := make([]interface{}, 0)
		para = append(para, ts_code)
		para = append(para, startDate)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(processLog, "", "")
	return nil
}

func (this *QPerformanceService) RefreshDayQPerformance() error {
	processLog := GetProcessLogService().StartLog("qperformance", "RefreshDayQPerformance", "")
	routinePool := thread.CreateRoutinePool(10, this.AsyncUpdateDayQPerformance, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	for _, ts_code := range ts_codes {
		para := make([]interface{}, 0)
		para = append(para, ts_code)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(processLog, "", "")
	return nil
}

/**
查询最新的股票季度业绩数据的时间
*/
func (this *QPerformanceService) findMaxQDate(ts_code string, lineType int64) (*entity.QPerformance, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	qperformances := make([]*entity.QPerformance, 0)
	if lineType != 0 {
		conds = conds + " and linetype=?"
		paras = append(paras, lineType)
	}
	err := this.Find(&qperformances, nil, "qdate desc,tradedate desc", 0, 1, conds, paras...)
	if err != nil {
		return nil, err
	}
	if len(qperformances) > 0 {
		return qperformances[0], nil
	}

	return nil, nil
}

/**
查询最新的股票季度业绩数据的时间
*/
func (this *QPerformanceService) findMinQDate(ts_code string, lineType int64) (*entity.QPerformance, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	qperformances := make([]*entity.QPerformance, 0)
	if lineType != 0 {
		conds = conds + " and linetype=?"
		paras = append(paras, lineType)
	}
	err := this.Find(&qperformances, nil, "qdate,tradedate desc", 0, 1, conds, paras...)
	if err != nil {
		return nil, err
	}
	if len(qperformances) > 0 {
		return qperformances[0], nil
	}

	return nil, nil
}

/**
删除股票的季度业绩数据
*/
func (this *QPerformanceService) deleteQPerformance(ts_code string, lineType LineType, startDate string) error {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	qperformance := &entity.QPerformance{}
	if lineType != 0 {
		if lineType == LineType_Wmqy {
			conds = conds + " and (linetype=? or linetype is null)"
		} else {
			conds = conds + " and linetype=?"
		}
		paras = append(paras, int64(lineType))
	}
	if startDate != "" {
		conds = conds + " and qdate>=?"
		paras = append(paras, startDate)
	}
	_, err := this.Delete(qperformance, conds, paras...)
	if err != nil {
		return err
	}

	return nil
}

func (this *QPerformanceService) AsyncUpdateWmqyQPerformance(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	startDate := (para.([]interface{}))[1].(string)
	this.GetUpdateWmqyQPerformance(tscode, startDate)
}

func (this *QPerformanceService) AsyncUpdateDayQPerformance(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	this.GetUpdateDayQPerformance(tscode)
}

/**
更新股票季度业绩数据，并返回结果
*/
func (this *QPerformanceService) GetUpdateWmqyQPerformance(tscode string, startDate string) ([]interface{}, error) {
	processLog := GetProcessLogService().StartLog("qperformance", "GetUpdateQPerformance", tscode)
	ps, err := this.updateWmqyQPerformance(tscode, startDate)
	if err != nil {
		GetProcessLogService().EndLog(processLog, "", err.Error())
		return nil, err
	}
	return ps, err
}

func (this *QPerformanceService) GetUpdateDayQPerformance(tscode string) ([]interface{}, error) {
	//processLog := GetProcessLogService().StartLog("qperformance", "GetUpdateQPerformance", tscode)
	this.deleteQPerformance(tscode, LineType_Day, "")
	ps, err := this.updateDayQPerformance(tscode)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return nil, err
	}
	GetQStatService().GetUpdateQStat(tscode)
	return ps, err
}

/**
更新股票季度业绩数据，并返回结果
*/
func (this *QPerformanceService) updateWmqyQPerformance(tscode string, startDate string) ([]interface{}, error) {
	qperformanceMap, err := this.findWmqyQPerformance(tscode, startDate, "", StdType_None, false)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, err.Error())
		return nil, err
	}
	if len(qperformanceMap) <= 0 {
		logger.Sugar.Errorf("Error:%v", "qperformances len is 0")
		return nil, errors.New("")
	}
	ps := make([]interface{}, 0)
	for _, qperformances := range qperformanceMap {
		for _, qperformance := range qperformances {
			ps = append(ps, qperformance)
		}
	}
	_, err = this.Upsert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}

	return ps, err
}

/**
更新股票季度业绩数据，并返回结果
*/
func (this *QPerformanceService) updateDayQPerformance(tscode string) ([]interface{}, error) {
	qperformanceMap, err := this.findDayQPerformance(tscode, StdType_None, false)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, err.Error())
		return nil, err
	}
	if len(qperformanceMap) <= 0 {
		logger.Sugar.Errorf("Error:%v", "qperformances len is 0")
		return nil, errors.New("")
	}
	ps := make([]interface{}, 0)
	for _, qperformances := range qperformanceMap {
		for _, qperformance := range qperformances {
			ps = append(ps, qperformance)
		}
	}
	_, err = this.Upsert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}

	return ps, err
}

/**
通过数据库sql计算累计涨幅
*/
func (this *QPerformanceService) FindAccBySql(ts_code string, startDate string) (map[string][]interface{}, error) {
	startDate, _ = stock.AddQuarter(startDate, -1)
	qpMap, err := this.FindStdQPerformance(ts_code, nil, startDate, "", StdType_None, false)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", ts_code, err.Error())
		return nil, err
	}
	qsMap := make(map[string][]interface{}, 0)
	for tscode, qps := range qpMap {
		marketValue := 0.0
		yearNetProfit := 0.0
		yearOperateIncome := 0.0
		for _, qp := range qps {
			if qp.MarketValue > 0.0 && stock.Equal(marketValue, 0.0) {
				marketValue = qp.MarketValue
			}
			if qp.YearNetProfit > 0.0 && stock.Equal(yearNetProfit, 0.0) {
				yearNetProfit = qp.YearNetProfit
			}
			if qp.YearOperateIncome > 0.0 && stock.Equal(yearOperateIncome, 0.0) {
				yearOperateIncome = qp.YearOperateIncome
			}
			if !stock.Equal(marketValue, 0.0) && !stock.Equal(yearNetProfit, 0.0) && !stock.Equal(yearNetProfit, 0.0) {
				break
			}
		}
		end := qps[len(qps)-1]
		qs := &entity.QStat{}
		if !stock.Equal(marketValue, 0.0) {
			qs.PctChgMarketValue = (end.MarketValue-marketValue)/math.Abs(marketValue) + 1
		}
		if !stock.Equal(yearNetProfit, 0.0) {
			qs.YoyDeduNp = (end.YearNetProfit-yearNetProfit)/math.Abs(yearNetProfit) + 1
		}
		if !stock.Equal(yearOperateIncome, 0.0) {
			qs.YoySales = (end.YearOperateIncome-yearOperateIncome)/math.Abs(yearOperateIncome) + 1
		}
		ps, ok := qsMap[tscode]
		if !ok {
			ps = make([]interface{}, 0)
		}
		qs.TsCode = tscode
		qs.SecurityName = end.SecurityName
		qs.Industry = end.Industry
		qs.Sector = end.Sector
		qs.Source = "acc"
		ps = append(ps, qs)
		qsMap[tscode] = ps
	}

	return qsMap, nil
}

type QTerm struct {
	ActualStartDate string
	StartDate       string
	EndDate         string
	Term            int
	TradeDate       int64
}

func (this *QPerformanceService) GetQTermBySql(tscode string, term int) (*QTerm, error) {
	qp, err := this.findMaxQDate(tscode, 0)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, err.Error())
		return nil, err
	}
	if qp == nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, errors.New("qp is nil"))
		return nil, errors.New("qperformance is nil")
	}

	endDate := qp.QDate
	tradeDate := qp.TradeDate
	qp, err = GetQPerformanceService().findMinQDate(tscode, 0)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, err.Error())
		return nil, err
	}
	if qp == nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, errors.New("qp is nil"))
		return nil, errors.New("qperformance is nil")
	}
	actualStartDate := qp.QDate
	startDate := ""
	if term > 0 {
		startDate, _ = stock.AddYear(endDate, -term)
		startDate, _ = stock.AddQuarter(startDate, 1)
	} else {
		startDate = actualStartDate
	}
	qterm := &QTerm{ActualStartDate: actualStartDate, StartDate: startDate, EndDate: endDate, Term: term, TradeDate: tradeDate}

	return qterm, nil
}

/**
在所有的数据中分选出各种不同的term对应的数据集，原始数据降序排列
*/
func (this *QPerformanceService) GetQTerm(qpMap map[string][]*entity.QPerformance, terms []int) (map[string]map[int]*QTerm, error) {
	qtermMap := make(map[string]map[int]*QTerm)
	for tscode, qps := range qpMap {
		qterms, ok := qtermMap[tscode]
		if !ok {
			qterms = make(map[int]*QTerm)
			qtermMap[tscode] = qterms
		}
		maxQp := qps[0]
		endDate := maxQp.QDate
		tradeDate := maxQp.TradeDate
		minQp := qps[len(qps)-1]
		actualStartDate := minQp.QDate

		startDate := ""
		for _, term := range terms {
			if term > 0 {
				startDate, _ = stock.AddYear(endDate, -term)
				startDate, _ = stock.AddQuarter(startDate, 1)
			} else {
				startDate = actualStartDate
			}
			qterm, ok := qterms[term]
			if !ok {
				if actualStartDate > startDate {
					//logger.Sugar.Errorf("tscode:%v Error:%v", tscode, "ActualStartDate>StartDate")
					continue
				}
				qterm = &QTerm{ActualStartDate: actualStartDate, StartDate: startDate, EndDate: endDate, Term: term, TradeDate: tradeDate}
				qterms[term] = qterm
			}

		}
	}

	return qtermMap, nil
}

/**
通过数据库sql计算股票季度业绩全部统计数据，并返回结果
*/
func (this *QPerformanceService) FindAllQStatBySql(ts_code string, startDate string, endDate string) map[string][]interface{} {
	qpMap, err := this.FindQStatBySql("sum", ts_code, startDate, endDate, "")
	if err != nil {
		qpMap = make(map[string][]interface{}, 0)
	}
	stats := make([]map[string][]interface{}, 0)
	maxMap, err := this.FindQStatBySql("max", ts_code, startDate, endDate, "")
	if err == nil {
		stats = append(stats, maxMap)
	}
	minMap, err := this.FindQStatBySql("min", ts_code, startDate, endDate, "")
	if err == nil {
		stats = append(stats, minMap)
	}
	meanMap, err := this.FindQStatBySql("mean", ts_code, startDate, endDate, "")
	if err == nil {
		stats = append(stats, meanMap)
	}
	medianMap, err := this.FindQStatBySql("median", ts_code, startDate, endDate, "")
	if err == nil {
		stats = append(stats, medianMap)
	}
	stddevMap, err := this.FindQStatBySql("stddev", ts_code, startDate, endDate, "")
	if err == nil {
		stats = append(stats, stddevMap)
	}
	rsdMap, err := this.FindQStatBySql("rsd", ts_code, startDate, endDate, "")
	if err == nil {
		stats = append(stats, rsdMap)
	}
	jsonMap, _, jsonHeads := stock.GetJsonMap(entity.QPerformance{})
	for _, jsonHead := range jsonHeads[14:] {
		fieldname := jsonMap[jsonHead]
		corrMap, err := this.FindQStatBySql("corr", ts_code, startDate, endDate, fieldname)
		if err == nil {
			stats = append(stats, corrMap)
		}
	}
	accMap, err := this.FindAccBySql(ts_code, startDate)
	if err == nil {
		stats = append(stats, accMap)
	}

	for tscode, qps := range qpMap {
		for _, stat := range stats {
			if stat != nil && len(stat) > 0 {
				ss, ok := stat[tscode]
				if ok && len(ss) > 0 {
					qps = append(qps, ss...)
					qpMap[tscode] = qps
				}
			}
		}
	}

	return qpMap
}

/**
通过数据库sql计算股票季度业绩某种统计数据，并返回结果
*/
func (this *QPerformanceService) FindQStatBySql(aggregationType string, ts_code string, startDate string, endDate string, sourceName string) (map[string][]interface{}, error) {
	jsonMap, _, jsonHeads := stock.GetJsonMap(entity.QPerformance{})
	sql := "select tscode as ts_code,SecurityName as security_name,industry as industry,sector as sector"
	for _, jsonHead := range jsonHeads[14:] {
		fieldname := jsonMap[jsonHead]
		if aggregationType == "corr" {
			if sourceName == "" {
				sql = sql + "," + aggregationType + "(marketvalue," + fieldname + ") as " + jsonHead
			} else {
				sql = sql + "," + aggregationType + "(" + sourceName + "," + fieldname + ") as " + jsonHead
			}
		} else if aggregationType == "median" {
			sql = sql + ",percentile_cont(0.5) within group(order by " + fieldname + ") as " + jsonHead
		} else if aggregationType == "rsd" {
			sql = sql + ",case when avg(" + fieldname + ")!=0 then stddev(" + fieldname + ")/avg(" + fieldname + ") else 0 end as " + jsonHead
		} else if aggregationType == "mean" {
			sql = sql + ",avg(" + fieldname + ") as " + jsonHead
		} else {
			sql = sql + "," + aggregationType + "(" + fieldname + ") as " + jsonHead
		}
	}
	sql = sql + " from stk_qperformance where 1=1"
	conds, tscodeParas := stock.InBuildStr("tscode", ts_code, ",")
	paras := make([]interface{}, 0)
	if ts_code != "" {
		sql = sql + " and " + conds
		paras = append(paras, tscodeParas...)
	}
	if startDate != "" {
		sql = sql + " and qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		sql = sql + " and qdate<=?"
		paras = append(paras, endDate)
	}

	sql = sql + " group by tscode,securityname,industry,sector"
	sql = sql + " having count(tscode)>0"
	sql = sql + " order by tscode"
	results, err := this.Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	qpMap := make(map[string][]interface{}, 0)
	for _, result := range results {
		qp := &entity.QStat{}
		for colname, v := range result {
			s := string(v)
			if s != "" {
				err = reflect.Set(qp, jsonMap[colname], s)
				if err != nil {
					logger.Sugar.Errorf("Set colname %v value %v error:%v", colname, s, err.Error())
				}
			}
		}
		ps, ok := qpMap[qp.TsCode]
		if !ok {
			ps = make([]interface{}, 0)
		}
		qp.Source = aggregationType
		qp.SourceName = sourceName
		ps = append(ps, qp)
		qpMap[qp.TsCode] = ps
	}

	return qpMap, nil
}

/**
计算股票的在本股票历史上以及在同行业历史上的位置,不支持tscode多只股票
*/
func (this *QPerformanceService) FindPercentRank(rankType string, ts_code string, tradeDate int64, startDate string, endDate string, from int, limit int, count int64) ([]*entity.QPerformance, error) {
	_, shares := GetShareService().GetCacheShare()
	in, inParas := stock.InBuildStr("ts_code", ts_code, ",")
	jsonMap, _, jsonHeads := stock.GetJsonMap(entity.QPerformance{})
	sql := "select id as id,tscode as ts_code,securityname as security_name,industry as industry,sector as sector"
	sql = sql + ",qdate as qdate,ndate as ndate,tradedate as trade_date,source as source,linetype as line_type"
	for _, jsonHead := range jsonHeads[14:] {
		fieldname := jsonMap[jsonHead]
		if rankType == "tscode" {
			sql = sql + ",percent_rank() over (partition by tscode order by " + fieldname + " asc) as " + jsonHead
		} else if rankType == "industry" {
			sql = sql + ",percent_rank() over (partition by industry order by " + fieldname + " asc) as " + jsonHead
		} else if rankType == "sector" {
			sql = sql + ",percent_rank() over (partition by sector order by " + fieldname + " asc) as " + jsonHead
		}
	}
	sql = sql + " from stk_qperformance"
	share, ok := shares[ts_code]
	if ok {
		if rankType == "tscode" {
			sql = sql + " where tscode='" + ts_code + "'"
		} else if rankType == "industry" {
			sql = sql + " where industry='" + share.Industry + "'"
		} else if rankType == "sector" {
			sql = sql + " where sector='" + share.Sector + "'"
		}
	}
	paras := make([]interface{}, 0)
	if tradeDate > 0 {
		sql = sql + " and tradedate = ?"
		paras = append(paras, tradeDate)
	}
	if startDate != "" {
		sql = sql + " and qdate >= ?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		sql = sql + " and qdate <= ?"
		paras = append(paras, endDate)
	}
	sql = "select * from (" + sql + ") t"
	sql = sql + " where " + in
	paras = append(paras, inParas...)
	sql = sql + " order by ts_code,qdate desc,trade_date desc"
	if from > 0 {
		sql = sql + " offset " + fmt.Sprint(from)
	}
	if limit > 0 {
		sql = sql + " limit " + fmt.Sprint(limit)
	}
	results, err := this.Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	qps := make([]*entity.QPerformance, 0)
	for _, result := range results {
		qp := &entity.QPerformance{}
		for colname, v := range result {
			s := string(v)
			if s != "" {
				err = reflect.Set(qp, jsonMap[colname], s)
				if err != nil {
					logger.Sugar.Errorf("Set colname %v value %v error:%v", colname, s, err.Error())
				}
			}
		}
		qps = append(qps, qp)
	}

	return qps, nil
}

func init() {
	service.GetSession().Sync(new(entity.QPerformance))
	qperformanceService.OrmBaseService.GetSeqName = qperformanceService.GetSeqName
	qperformanceService.OrmBaseService.FactNewEntity = qperformanceService.NewEntity
	qperformanceService.OrmBaseService.FactNewEntities = qperformanceService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("qperformance", qperformanceService)
}
