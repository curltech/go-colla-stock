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
	"github.com/curltech/go-colla-stock/stock/entity"
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

func (svc *QPerformanceService) GetSeqName() string {
	return seqname
}

func (svc *QPerformanceService) NewEntity(data []byte) (interface{}, error) {
	qperformance := &entity.QPerformance{}
	if data == nil {
		return qperformance, nil
	}
	err := message.Unmarshal(data, qperformance)
	if err != nil {
		return nil, err
	}

	return qperformance, err
}

func (svc *QPerformanceService) NewEntities(data []byte) (interface{}, error) {
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

// / 标准查询，tsCode和tradeDate，qDate必有一个不为空
func (svc *QPerformanceService) FindByCondContent(tsCode string, qDate string, tradeDate int64, condContent string, condParas []interface{}, orderby string, from int, limit int, count int64) ([]*entity.QPerformance, int64, error) {
	if tsCode == "" && qDate == "" && tradeDate == 0 {
		return nil, 0, errors.New("tsCode, qDate and tradeDate can't all be empty")
	}
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	if condContent != "" {
		conds = conds + " and " + condContent
		if condParas != nil && len(condParas) > 0 {
			paras = append(paras, condParas...)
		}
	}
	if qDate != "" {
		conds = conds + " and qdate=?"
		paras = append(paras, qDate)
	}
	if tradeDate != 0 {
		conds = conds + " and tradeDate=?"
		paras = append(paras, tradeDate)
	}
	var err error
	condiBean := &entity.QPerformance{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "tscode,qdate desc"
	}
	ps := make([]*entity.QPerformance, 0)
	if limit == 0 {
		limit = 10
	}
	err = svc.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return ps, count, nil
}

func (svc *QPerformanceService) Search(keyword string, terms []int, sourceOptions []string, startDate string, endDate string, condContent string, condParas []interface{}, orderby string, from int, limit int, count int64) ([]*entity.QPerformance, int64, error) {
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
		tscodes := ""
		if keyword != "" {
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
				i++
			}
		}
		qterm, err := svc.GetQTermBySql(tscodes, term)
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

	if keyword != "" {
		conds += " and tscode in (select tscode from stk_share where name like ? or tscode like ? or pinyin like ?)"
		paras = append(paras, keyword+"%")
		paras = append(paras, keyword+"%")
		paras = append(paras, strings.ToLower(keyword)+"%")
	}
	if condContent != "" {
		conds = conds + " and " + condContent
		if condParas != nil && len(condParas) > 0 {
			paras = append(paras, condParas...)
		}
	}
	qperformances := make([]*entity.QPerformance, 0)
	var err error
	condiBean := &entity.QPerformance{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "tscode,qdate desc,tradedate desc"
	}
	err = svc.Find(&qperformances, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return qperformances, count, nil
}

type PerformanceType int64

type StdType int64

const (
	StdtypeNone StdType = iota + 1
	StdtypeStd
	StdtypeMinmax
)

type LineType int64

const (
	LinetypeDay LineType = iota + 1
	LinetypeWmqy
)

// FindQPerformance 组合当前价格和股票季度正式业绩数据
func (svc *QPerformanceService) FindQPerformance(lineType LineType, tsCode string, startDate string, endDate string) (map[string][]interface{}, error) {
	lineName := "stk_dayline"
	if lineType == LinetypeWmqy {
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
	conds, tscodeParas := stock.InBuildStr("w.tscode", tsCode, ",")
	paras := make([]interface{}, 0)
	if lineType == LinetypeWmqy {
		sql = sql + " where w.qdate = sp.qdate"
	} else {
		sql = sql + " where w.tradedate = ?"
		tscode := "000001"
		if tsCode != "" {
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
	if tsCode != "" {
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
	qpMap := make(map[string][]interface{})
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

// Std 对股票季度业绩数据进行异常值处理和标准化操作
func (svc *QPerformanceService) Std(ps []interface{}, stdType StdType, isWinsorize bool) []*entity.QPerformance {
	qps := make([]*entity.QPerformance, 0)
	if stdType == 0 || stdType == StdtypeNone {
		for _, qp := range ps {
			qps = append(qps, qp.(*entity.QPerformance))
		}
		return qps
	}
	_, _, heads := stock.GetJsonMap(entity.QPerformance{})
	stdColnames := heads[12:]
	reservedColnames := heads[0:12]
	stat := stock.CreateStat(ps, stdColnames)
	stds, minmaxs := stat.CalStd(reservedColnames, isWinsorize)
	if stdType == StdtypeMinmax {
		for _, qp := range minmaxs {
			qps = append(qps, qp.(*entity.QPerformance))
		}
	} else if stdType == StdtypeStd {
		for _, qp := range stds {
			qps = append(qps, qp.(*entity.QPerformance))
		}
	}

	return qps
}

// 组合当前价格和股票季度业绩、预测和快报数据，形成完整的股票季度业绩数据，进行标准化操作
func (svc *QPerformanceService) findWmqyQPerformance(tsCode string, startDate string, endDate string, stdType StdType, isWinsorize bool) (map[string][]*entity.QPerformance, error) {
	previousMap := make(map[string]*entity.QPerformance)
	if startDate == "" {
		qp, _ := svc.findMaxQDate(tsCode, int64(LinetypeWmqy))
		if qp != nil {
			startDate = qp.QDate
		}
	}
	err := svc.deleteQPerformance(tsCode, LinetypeWmqy, startDate)
	if err != nil {
		return nil, err
	}
	qp, _ := svc.findMaxQDate(tsCode, int64(LinetypeWmqy))
	if qp != nil {
		previousMap[tsCode] = qp
	}
	qpMap, err := svc.FindQPerformance(LinetypeWmqy, tsCode, startDate, endDate)
	if err != nil {
		return nil, err
	}
	qeMap, err := svc.FindQExpress(LinetypeWmqy, tsCode, startDate, endDate)
	if err != nil {
		return nil, err
	}
	svc.merge(qeMap, qpMap)

	qfMap, err := svc.FindQForecast(LinetypeWmqy, tsCode, startDate, endDate)
	if err != nil {
		return nil, err
	}
	svc.merge(qfMap, qpMap)

	svc.Compute(qpMap, previousMap)

	stds := svc.StdMap(qpMap, stdType, isWinsorize)

	return stds, nil
}

func (svc *QPerformanceService) findDayQPerformance(tsCode string, stdType StdType, isWinsorize bool) (map[string][]*entity.QPerformance, error) {
	var dqMap map[string][]interface{}
	var err error
	performanceQDate, err := GetPerformanceService().findMaxQDate(tsCode)
	if err != nil {
		logger.Sugar.Errorf("Performance findMaxQDate error:%v", err.Error())
	}
	expressQDate, err := GetExpressService().findMaxQDate(tsCode)
	if err != nil {
		logger.Sugar.Errorf("Express findMaxQDate error:%v", err.Error())
	}
	var startDate string
	if performanceQDate >= expressQDate {
		startDate = performanceQDate
	} else {
		startDate = expressQDate
	}
	forecastQDate, err := GetForecastService().findMaxQDate(tsCode)
	if err != nil {
		logger.Sugar.Errorf("Forecast findMaxQDate error:%v", err.Error())
	}
	if forecastQDate >= startDate {
		startDate = forecastQDate
	}
	endDate := startDate
	if performanceQDate == startDate {
		dqMap, err = svc.FindQPerformance(LinetypeDay, tsCode, startDate, endDate)
		if err != nil {
			return nil, err
		}
	} else if expressQDate == startDate {
		dqMap, err = svc.FindQExpress(LinetypeDay, tsCode, startDate, endDate)
		if err != nil {
			return nil, err
		}
	} else if forecastQDate == startDate {
		dqMap, err = svc.FindQForecast(LinetypeDay, tsCode, startDate, endDate)
		if err != nil {
			return nil, err
		}

	}
	previousMap := make(map[string]*entity.QPerformance)
	qp, err := svc.findMaxQDate(tsCode, 0)
	if qp != nil {
		previousMap[tsCode] = qp
	}
	svc.Compute(dqMap, previousMap)

	stds := svc.StdMap(dqMap, stdType, isWinsorize)

	return stds, nil
}

// FindStdQPerformance 查询股票季度业绩数据，并进行标准化处理
func (svc *QPerformanceService) FindStdQPerformance(tsCode string, terms []int, startDate string, endDate string, stdType StdType, isWinsorize bool) (map[string][]*entity.QPerformance, error) {
	qps, _, err := svc.Search(tsCode, terms, nil, startDate, endDate, "", nil, "", 0, 0, 1)
	if err != nil {
		return nil, err
	}
	qpMap := make(map[string][]interface{})
	for _, qp := range qps {
		ps, ok := qpMap[qp.TsCode]
		if !ok {
			ps = make([]interface{}, 0)
		}
		ps = append(ps, qp)
		qpMap[qp.TsCode] = ps
	}

	stds := svc.StdMap(qpMap, stdType, isWinsorize)

	return stds, nil
}

// StdMap 批量标准化股票季度业绩数据
func (svc *QPerformanceService) StdMap(psMap map[string][]interface{}, stdType StdType, isWinsorize bool) map[string][]*entity.QPerformance {
	stds := make(map[string][]*entity.QPerformance)
	for tscode, ps := range psMap {
		std := svc.Std(ps, stdType, isWinsorize)
		stds[tscode] = std
	}

	return stds
}

// 合并当前价格和股票季度业绩、预测、快报数据，按照时间排序
func (svc *QPerformanceService) merge(qeMap map[string][]interface{}, qpMap map[string][]interface{}) {
	for tsCode, qes := range qeMap {
		qps, ok := qpMap[tsCode]
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
					qpMap[tsCode] = qps
				} else if qe.QDate == qdate {
					if qe.TradeDate > tradeDate {
						qps = append(qps, qe)
						qpMap[tsCode] = qps
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

// Compute 计算并推导缺失的当前价格和股票季度业绩数据项
func (svc *QPerformanceService) Compute(qpMap map[string][]interface{}, previousMap map[string]*entity.QPerformance) {
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
			if previous != nil && stock.Equal(qp.GrossProfitMargin, 0.0) {
				qp.GrossProfitMargin = previous.GrossProfitMargin
			}
			if qp.Industry == "" {
				share := GetShareService().GetCacheShare(qp.TsCode)
				if share != nil {
					qp.Industry = share.Industry
				}
			}
			if qp.Sector == "" {
				share := GetShareService().GetCacheShare(qp.TsCode)
				if share != nil {
					qp.Sector = share.Sector
				}
			}

			previous = qp
		}
	}
}

// FindQExpress 组合当前价格和股票季度业绩快报数据
func (svc *QPerformanceService) FindQExpress(lineType LineType, tsCode string, startDate string, endDate string) (map[string][]interface{}, error) {
	lineName := "stk_dayline"
	if lineType == LinetypeWmqy {
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
	conds, tscodeParas := stock.InBuildStr("w.tscode", tsCode, ",")
	paras := make([]interface{}, 0)
	if lineType == LinetypeWmqy {
		sql = sql + " where w.qdate = sp.qdate"
	} else {
		sql = sql + " where w.tradedate = ?"
		tscode := "000001"
		if tsCode != "" {
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
	if tsCode != "" {
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
	qpMap := make(map[string][]interface{})
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

// FindQForecast 组合当前价格和股票季度业绩预测数据
func (svc *QPerformanceService) FindQForecast(lineType LineType, tsCode string, startDate string, endDate string) (map[string][]interface{}, error) {
	lineName := "stk_dayline"
	if lineType == LinetypeWmqy {
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
	conds, tscodeParas := stock.InBuildStr("w.tscode", tsCode, ",")
	paras := make([]interface{}, 0)
	if lineType == LinetypeWmqy {
		sql = sql + " where w.qdate = sp.qdate"
	} else {
		sql = sql + " where w.tradedate = ?"
		tscode := "000001"
		if tsCode != "" {
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
	if tsCode != "" {
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
	qpMap := make(map[string][]interface{})
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

// RefreshWmqyQPerformance 刷新所有股票的季度业绩数据
func (svc *QPerformanceService) RefreshWmqyQPerformance(startDate string) error {
	processLog := GetProcessLogService().StartLog("qperformance", "RefreshWmqyQPerformance", "")
	routinePool := thread.CreateRoutinePool(10, svc.AsyncUpdateWmqyQPerformance, nil)
	defer routinePool.Release()
	tsCodes, _ := GetShareService().GetShareCache()
	for _, tsCode := range tsCodes {
		para := make([]interface{}, 0)
		para = append(para, tsCode)
		para = append(para, startDate)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(processLog, "", "")
	return nil
}

func (svc *QPerformanceService) RefreshDayQPerformance() error {
	processLog := GetProcessLogService().StartLog("qperformance", "RefreshDayQPerformance", "")
	routinePool := thread.CreateRoutinePool(10, svc.AsyncUpdateDayQPerformance, nil)
	defer routinePool.Release()
	tsCodes, _ := GetShareService().GetShareCache()
	for _, tsCode := range tsCodes {
		para := make([]interface{}, 0)
		para = append(para, tsCode)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(processLog, "", "")
	return nil
}

// 查询最新的股票季度业绩数据的时间
func (svc *QPerformanceService) findMaxQDate(tsCode string, lineType int64) (*entity.QPerformance, error) {
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	qperformances := make([]*entity.QPerformance, 0)
	if lineType != 0 {
		conds = conds + " and linetype=?"
		paras = append(paras, lineType)
	}
	err := svc.Find(&qperformances, nil, "qdate desc,tradedate desc", 0, 1, conds, paras...)
	if err != nil {
		return nil, err
	}
	if len(qperformances) > 0 {
		return qperformances[0], nil
	}

	return nil, nil
}

// 查询最新的股票季度业绩数据的时间
func (svc *QPerformanceService) findMinQDate(tsCode string, lineType int64) (*entity.QPerformance, error) {
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	qperformances := make([]*entity.QPerformance, 0)
	if lineType != 0 {
		conds = conds + " and linetype=?"
		paras = append(paras, lineType)
	}
	err := svc.Find(&qperformances, nil, "qdate,tradedate desc", 0, 1, conds, paras...)
	if err != nil {
		return nil, err
	}
	if len(qperformances) > 0 {
		return qperformances[0], nil
	}

	return nil, nil
}

// 删除股票的季度业绩数据
func (svc *QPerformanceService) deleteQPerformance(tsCode string, lineType LineType, startDate string) error {
	conds, paras := stock.InBuildStr("tscode", tsCode, ",")
	qperformance := &entity.QPerformance{}
	if lineType != 0 {
		if lineType == LinetypeWmqy {
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
	_, err := svc.Delete(qperformance, conds, paras...)
	if err != nil {
		return err
	}

	return nil
}

func (svc *QPerformanceService) AsyncUpdateWmqyQPerformance(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	startDate := (para.([]interface{}))[1].(string)
	_, err := svc.GetUpdateWmqyQPerformance(tscode, startDate)
	if err != nil {
		return
	}
}

func (svc *QPerformanceService) AsyncUpdateDayQPerformance(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	_, err := svc.GetUpdateDayQPerformance(tscode)
	if err != nil {
		return
	}
}

// GetUpdateWmqyQPerformance 更新股票季度业绩数据，并返回结果
func (svc *QPerformanceService) GetUpdateWmqyQPerformance(tscode string, startDate string) ([]interface{}, error) {
	processLog := GetProcessLogService().StartLog("qperformance", "GetUpdateQPerformance", tscode)
	ps, err := svc.updateWmqyQPerformance(tscode, startDate)
	if err != nil {
		GetProcessLogService().EndLog(processLog, "", err.Error())
		return nil, err
	}
	return ps, err
}

func (svc *QPerformanceService) GetUpdateDayQPerformance(tscode string) ([]interface{}, error) {
	//processLog := GetProcessLogService().StartLog("qperformance", "GetUpdateQPerformance", tscode)
	err := svc.deleteQPerformance(tscode, LinetypeDay, "")
	if err != nil {
		return nil, err
	}
	ps, err := svc.updateDayQPerformance(tscode)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return nil, err
	}
	_, err = GetQStatService().GetUpdateQStat(tscode)
	if err != nil {
		return nil, err
	}
	return ps, err
}

// 更新股票季度业绩数据，并返回结果
func (svc *QPerformanceService) updateWmqyQPerformance(tscode string, startDate string) ([]interface{}, error) {
	qperformanceMap, err := svc.findWmqyQPerformance(tscode, startDate, "", StdtypeNone, false)
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
	_, err = svc.Upsert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}

	return ps, err
}

// 更新股票季度业绩数据，并返回结果
func (svc *QPerformanceService) updateDayQPerformance(tscode string) ([]interface{}, error) {
	qperformanceMap, err := svc.findDayQPerformance(tscode, StdtypeNone, false)
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
	_, err = svc.Upsert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}

	return ps, err
}

// FindAccBySql 通过数据库sql计算累计涨幅
func (svc *QPerformanceService) FindAccBySql(tsCode string, startDate string) (map[string][]interface{}, error) {
	startDate, _ = stock.AddQuarter(startDate, -1)
	qpMap, err := svc.FindStdQPerformance(tsCode, nil, startDate, "", StdtypeNone, false)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tsCode, err.Error())
		return nil, err
	}
	qsMap := make(map[string][]interface{})
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

func (svc *QPerformanceService) GetQTermBySql(tscode string, term int) (*QTerm, error) {
	qp, err := svc.findMaxQDate(tscode, 0)
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

// GetQTerm 在所有的数据中分选出各种不同的term对应的数据集，原始数据降序排列
func (svc *QPerformanceService) GetQTerm(qpMap map[string][]*entity.QPerformance, terms []int) (map[string]map[int]*QTerm, error) {
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

// FindAllQStatBySql 通过数据库sql计算股票季度业绩全部统计数据，并返回结果
func (svc *QPerformanceService) FindAllQStatBySql(tsCode string, startDate string, endDate string) map[string][]interface{} {
	qpMap, err := svc.FindQStatBySql("sum", tsCode, startDate, endDate, "")
	if err != nil {
		qpMap = make(map[string][]interface{})
	}
	stats := make([]map[string][]interface{}, 0)
	maxMap, err := svc.FindQStatBySql("max", tsCode, startDate, endDate, "")
	if err == nil {
		stats = append(stats, maxMap)
	}
	minMap, err := svc.FindQStatBySql("min", tsCode, startDate, endDate, "")
	if err == nil {
		stats = append(stats, minMap)
	}
	meanMap, err := svc.FindQStatBySql("mean", tsCode, startDate, endDate, "")
	if err == nil {
		stats = append(stats, meanMap)
	}
	medianMap, err := svc.FindQStatBySql("median", tsCode, startDate, endDate, "")
	if err == nil {
		stats = append(stats, medianMap)
	}
	stddevMap, err := svc.FindQStatBySql("stddev", tsCode, startDate, endDate, "")
	if err == nil {
		stats = append(stats, stddevMap)
	}
	rsdMap, err := svc.FindQStatBySql("rsd", tsCode, startDate, endDate, "")
	if err == nil {
		stats = append(stats, rsdMap)
	}
	jsonMap, _, jsonHeads := stock.GetJsonMap(entity.QPerformance{})
	for _, jsonHead := range jsonHeads[14:] {
		fieldname := jsonMap[jsonHead]
		corrMap, err := svc.FindQStatBySql("corr", tsCode, startDate, endDate, fieldname)
		if err == nil {
			stats = append(stats, corrMap)
		}
	}
	accMap, err := svc.FindAccBySql(tsCode, startDate)
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

// FindQStatBySql 通过数据库sql计算股票季度业绩某种统计数据，并返回结果
func (svc *QPerformanceService) FindQStatBySql(aggregationType string, tsCode string, startDate string, endDate string, sourceName string) (map[string][]interface{}, error) {
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
	conds, tscodeParas := stock.InBuildStr("tscode", tsCode, ",")
	paras := make([]interface{}, 0)
	if tsCode != "" {
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
	results, err := svc.Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	qpMap := make(map[string][]interface{})
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

// FindPercentRank 计算股票的在本股票历史上以及在同行业历史上的位置,不支持tscode多只股票
func (svc *QPerformanceService) FindPercentRank(rankType string, tsCode string, tradeDate int64, startDate string, endDate string, from int, limit int, count int64) ([]*entity.QPerformance, error) {
	in, inParas := stock.InBuildStr("ts_code", tsCode, ",")
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
	share := GetShareService().GetCacheShare(tsCode)
	if share != nil {
		if rankType == "tscode" {
			sql = sql + " where tscode='" + tsCode + "'"
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
	results, err := svc.Query(sql, paras...)
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
	err := service.GetSession().Sync(new(entity.QPerformance))
	if err != nil {
		return
	}
	qperformanceService.OrmBaseService.GetSeqName = qperformanceService.GetSeqName
	qperformanceService.OrmBaseService.FactNewEntity = qperformanceService.NewEntity
	qperformanceService.OrmBaseService.FactNewEntities = qperformanceService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("qperformance", qperformanceService)
}
