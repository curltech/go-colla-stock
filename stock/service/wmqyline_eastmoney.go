package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
*
获取某只股票的日线数据
*/
func (svc *WmqyLineService) GetWmqyLine(secId string, beg int64, end int, limit int, klt int, previous *entity.WmqyLine) ([]*entity.WmqyLine, error) {
	klines, err := GetDayLineService().GetKLine(secId, int(beg), end, limit, klt)
	if err != nil {
		return nil, err
	}
	wmqyLines := make([]*entity.WmqyLine, 0)
	for _, kline := range klines {
		wmqyLine, _ := strToWmqyLine(secId, kline)
		if wmqyLine != nil {
			wmqyLine.LineType = klt
			if klt == 104 {
				wmqyLine.QDate = stock.GetQTradeDate(wmqyLine.TradeDate)
			} else if klt == 105 {
				year := wmqyLine.TradeDate / 10000
				month := wmqyLine.TradeDate - year*10000
				month = month / 100
				if month < 7 {
					wmqyLine.QDate = fmt.Sprint(year) + "06"
				} else {
					wmqyLine.QDate = fmt.Sprint(year) + "12"
				}
			} else if klt == 106 {
				wmqyLine.QDate = fmt.Sprint(wmqyLine.TradeDate / 10000)
			} else if klt == 103 {
				wmqyLine.QDate = fmt.Sprint(wmqyLine.TradeDate / 100)
			} else if klt == 102 {
				wmqyLine.QDate = stock.GetWTradeDate(wmqyLine.TradeDate)
			}
			if previous != nil && previous.Open != 0.0 {
				wmqyLine.PctChgOpen = wmqyLine.Open/previous.Open - 1
			}
			if previous != nil && previous.High != 0.0 {
				wmqyLine.PctChgHigh = wmqyLine.High/previous.High - 1
			}
			if previous != nil && previous.Low != 0.0 {
				wmqyLine.PctChgLow = wmqyLine.Low/previous.Low - 1
			}
			if previous != nil && previous.Close != 0.0 {
				wmqyLine.PctChgClose = wmqyLine.Close/previous.Close - 1
			}
			if previous != nil && previous.Amount != 0.0 {
				wmqyLine.PctChgAmount = wmqyLine.Amount/previous.Amount - 1
			}
			if previous != nil && previous.Vol != 0.0 {
				wmqyLine.PctChgVol = wmqyLine.Vol/previous.Vol - 1
			}
			if previous != nil {
				wmqyLine.PreClose = previous.Close
			}
			previous = wmqyLine
			wmqyLines = append(wmqyLines, wmqyLine)
		}
	}

	return wmqyLines, err
}

func (svc *WmqyLineService) findByTradeDate(ts_code string, line_type int, startDate int64, endDate int64) ([]*entity.WmqyLine, error) {
	cond := &entity.WmqyLine{}
	cond.TsCode = ts_code
	cond.LineType = line_type
	wmqyLines := make([]*entity.WmqyLine, 0)
	conds := "? <= tradedate"
	paras := make([]interface{}, 0)
	paras = append(paras, startDate)
	if endDate > 0 {
		conds += " and tradedate <= ?"
		paras = append(paras, endDate)
	}
	err := svc.Find(&wmqyLines, cond, "tradedate", 0, 0, conds, paras...)

	return wmqyLines, err
}

func strToWmqyLine(secId string, kline string) (*entity.WmqyLine, error) {
	kls := strings.Split(kline, ",")
	wmqyLine := &entity.WmqyLine{}
	wmqyLine.TsCode = secId
	//"trade_date,open,close,high,low,vol,amount,nil,pct_chg%,change,turnover%"
	tradeDate, err := strconv.ParseInt(strings.ReplaceAll(kls[0], "-", ""), 10, 64)
	if err != nil {
		logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
		return nil, err
	}
	wmqyLine.TradeDate = tradeDate
	wmqyLine.Open, err = strToFloat(kls[1])
	if err != nil {
		return nil, err
	}
	wmqyLine.Close, err = strToFloat(kls[2])
	if err != nil {
		return nil, err
	}
	wmqyLine.High, err = strToFloat(kls[3])
	if err != nil {
		return nil, err
	}
	wmqyLine.Low, err = strToFloat(kls[4])
	if err != nil {
		return nil, err
	}
	wmqyLine.Vol, err = strToFloat(kls[5])
	if err != nil {
		return nil, err
	}
	wmqyLine.Amount, err = strToFloat(kls[6])
	if err != nil {
		return nil, err
	}
	pctChg, err := strToFloat(kls[8])
	if err != nil {
		return nil, err
	}
	wmqyLine.PctChgClose = pctChg
	wmqyLine.ChgClose, err = strToFloat(kls[9])
	if err != nil {
		return nil, err
	}
	wmqyLine.Turnover, err = strToFloat(kls[10])
	if err != nil {
		return nil, err
	}

	return wmqyLine, nil
}

type WmqyLineBuf struct {
	Buf  []interface{}
	Lock sync.Mutex
}

/*
*
刷新所有股票的日线数据，beg为负数的时候从已有的最新数据开始更新
*/
func (svc *WmqyLineService) RefreshWmqyLine(beg int64) error {
	processLog := GetProcessLogService().StartLog("wmqyLine", "RefreshWmqyLine", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, svc.AsyncUpdateWmqyLine, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	wmqy := &WmqyLineBuf{Buf: make([]interface{}, 0), Lock: sync.Mutex{}}
	for _, ts_code := range ts_codes {
		para := make([]interface{}, 0)
		para = append(para, ts_code)
		para = append(para, beg)
		para = append(para, 10000)
		para = append(para, wmqy)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	wmqy.Lock.Lock()
	defer wmqy.Lock.Unlock()
	if len(wmqy.Buf) > 0 {
		_, err := svc.Upsert(wmqy.Buf...)
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
		wmqy.Buf = make([]interface{}, 0)
	}
	GetProcessLogService().EndLog(processLog, "", "")
	return nil
}

func (svc *WmqyLineService) AsyncUpdateWmqyLine(para interface{}) {
	secId := (para.([]interface{}))[0].(string)
	beg := (para.([]interface{}))[1].(int64)
	limit := (para.([]interface{}))[2].(int)
	wmqy := (para.([]interface{}))[3].(*WmqyLineBuf)
	svc.GetUpdateWmqyLine(secId, beg, limit, wmqy)
}

func (svc *WmqyLineService) GetUpdateWmqyLine(secId string, beg int64, limit int, wmqy *WmqyLineBuf) ([]interface{}, error) {
	//processLog := GetProcessLogService().StartLog("wmqyLine", "GetUpdateWmqyLine", secId)
	ps, err := svc.UpdateWmqyLine(secId, beg, limit, wmqy)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	return ps, err
}

func (svc *WmqyLineService) UpdateWmqyLine(secId string, beg int64, limit int, wmqy *WmqyLineBuf) ([]interface{}, error) {
	ps := make([]interface{}, 0)
	for klt := 102; klt <= 106; klt++ {
		wmqyLines, err := svc.getWmqyLine(secId, klt, beg, limit)
		if err != nil {
			logger.Sugar.Errorf("Error:%v", err.Error())
			return ps, err
		}
		if len(wmqyLines) <= 0 {
			logger.Sugar.Errorf("Error:%v", "wmqyLines len is 0")
			return ps, err
		}
		for _, wmqyLine := range wmqyLines {
			if !stock.Equal(wmqyLine.Turnover, 0) {
				wmqyLine.ShareNumber = wmqyLine.Amount / wmqyLine.Turnover
			}
			ps = append(ps, wmqyLine)
			svc.deleteWmqyLine(wmqyLine.TsCode, wmqyLine.QDate)
		}
		svc.batchUpdate(wmqy, wmqyLines)
	}
	if wmqy == nil {
		_, err := svc.Upsert(ps...)
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
	}

	return ps, nil
}

func (svc *WmqyLineService) batchUpdate(wmqy *WmqyLineBuf, wmqyLines []*entity.WmqyLine) {
	if wmqy != nil {
		wmqy.Lock.Lock()
		defer wmqy.Lock.Unlock()
		for _, wmqyLine := range wmqyLines {
			wmqy.Buf = append(wmqy.Buf, wmqyLine)
		}
		if len(wmqy.Buf) > 1000 {
			_, err := svc.Upsert(wmqy.Buf...)
			if err != nil {
				logger.Sugar.Errorf("Error: %s", err.Error())
				return
			}

			wmqy.Buf = make([]interface{}, 0)
		}
	}
}

func (svc *WmqyLineService) getWmqyLine(secId string, klt int, beg int64, limit int) ([]*entity.WmqyLine, error) {
	currentDate := stock.CurrentDate()
	today := currentDate
	var todayReportDate string
	if klt == 104 {
		todayReportDate = stock.GetQTradeDate(today)
	} else if klt == 105 {
		year := today / 10000
		month := today - year*10000
		month = month / 100
		if month < 7 {
			todayReportDate = fmt.Sprint(year) + "06"
		} else {
			todayReportDate = fmt.Sprint(year) + "12"
		}
	} else if klt == 106 {
		todayReportDate = fmt.Sprint(today / 10000)
	} else if klt == 103 {
		todayReportDate = fmt.Sprint(today / 100)
	} else if klt == 102 {
		todayReportDate = stock.GetWTradeDate(today)
	}
	var end = 0
	previous, previous1, err := svc.findMaxTradeDate(secId, klt)
	if previous != nil {
		if todayReportDate == previous.QDate {
			svc.Delete(previous, "")
			previous = previous1
		}
	}
	if beg < 0 {
		if previous != nil {
			beg = previous.TradeDate
		} else {
			beg = 0
		}
		if beg > 0 {
			beg++
		}
	}
	if beg > 0 && beg > currentDate {
		logger.Sugar.Errorf("Error:%v", errors.New("data is updated"))
		return nil, errors.New("")
	}
	wmqyLines, err := svc.GetWmqyLine(secId, beg, end, limit, klt, previous)

	return wmqyLines, err
}

func (svc *WmqyLineService) findAggregation(startDate int64, endDate int64) (map[stock.AggregationType]map[string]*entity.WmqyLine, error) {
	sql := "select tscode TsCode,linetype LineType"
	jsonMap, _, _ := stock.GetJsonMap(entity.WmqyLine{})
	for _, colname := range DaylineHeader[2:] {
		fieldname, ok := jsonMap[colname]
		if ok {
			lower_fieldname := strings.ToLower(fieldname)
			sql = sql + ",max(" + lower_fieldname + ")" + " max_" + fieldname
			sql = sql + ",min(" + lower_fieldname + ")" + " min_" + fieldname
			sql = sql + ",avg(" + lower_fieldname + ")" + " avg_" + fieldname
			sql = sql + ",stddev(" + lower_fieldname + ")" + " stddev_" + fieldname
		}
	}
	sql = sql + " from stk_dayline"
	conds := " where ? <= tradedate"
	paras := make([]interface{}, 0)
	paras = append(paras, startDate)
	if endDate > 0 {
		conds += " and tradedate <= ?"
		paras = append(paras, endDate)
	}
	sql = sql + conds
	sql = sql + " group by tscode,linetype"
	results, err := svc.Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	aggreWmqylines := make(map[stock.AggregationType]map[string]*entity.WmqyLine)
	aggreWmqylines[stock.Aggregation_MAX] = make(map[string]*entity.WmqyLine, 0)
	aggreWmqylines[stock.Aggregation_MIN] = make(map[string]*entity.WmqyLine, 0)
	aggreWmqylines[stock.Aggregation_MEAN] = make(map[string]*entity.WmqyLine, 0)
	aggreWmqylines[stock.Aggregation_STDDEV] = make(map[string]*entity.WmqyLine, 0)
	for _, result := range results {
		maxWmqyline := &entity.WmqyLine{}
		minWmqyline := &entity.WmqyLine{}
		meanWmqyline := &entity.WmqyLine{}
		stddevWmqyline := &entity.WmqyLine{}
		for colname, v := range result {
			if colname == "TsCode" {
				reflect.Set(maxWmqyline, "TsCode", string(v))
			} else {
				fieldnames := strings.Split(colname, "_")
				if fieldnames[0] == "max" {
					reflect.Set(maxWmqyline, fieldnames[1], string(v))
				}
				if fieldnames[0] == "min" {
					reflect.Set(minWmqyline, fieldnames[1], string(v))
				}
				if fieldnames[0] == "mean" {
					reflect.Set(meanWmqyline, fieldnames[1], string(v))
				}
				if fieldnames[0] == "stddev" {
					reflect.Set(stddevWmqyline, fieldnames[1], string(v))
				}
			}
		}
		aggreWmqylines[stock.Aggregation_MAX][maxWmqyline.TsCode] = maxWmqyline
		aggreWmqylines[stock.Aggregation_MIN][minWmqyline.TsCode] = minWmqyline
		aggreWmqylines[stock.Aggregation_MEAN][meanWmqyline.TsCode] = meanWmqyline
		aggreWmqylines[stock.Aggregation_STDDEV][stddevWmqyline.TsCode] = stddevWmqyline
	}

	return aggreWmqylines, nil
}

/*
*
minmax: "C:\stock\data\minmax\lday"
standard: "C:\stock\data\standard\lday"
*/
func (svc *WmqyLineService) StdPath(minmax string, standard string, startDate int64, endDate int64) error {
	stock.Mkdir(minmax + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate))
	stock.Mkdir(standard + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate))
	aggreStrs := make(map[stock.AggregationType][]string)
	for _, typ := range stock.AggregationTypes {
		aggreStrs[typ] = make([]string, 0)
	}

	routinePool := thread.CreateRoutinePool(10, svc.AsyncStdFile, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	for _, ts_code := range ts_codes {
		para := make([]interface{}, 0)
		para = append(para, minmax)
		para = append(para, standard)
		para = append(para, startDate)
		para = append(para, endDate)
		para = append(para, ts_code)
		para = append(para, aggreStrs)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	raw := "ts_code," + strings.Join(WmqylineHeader[3:], ",") + "\n"
	max := raw + strings.Join(aggreStrs[stock.Aggregation_MAX], "\n")
	maxFileName := minmax + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\max.csv"
	err := ioutil.WriteFile(maxFileName, []byte(max), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v max failure!", maxFileName)
	}

	min := raw + strings.Join(aggreStrs[stock.Aggregation_MIN], "\n")
	minFileName := minmax + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\min.csv"
	err = ioutil.WriteFile(minFileName, []byte(min), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v min failure!", minFileName)
	}

	mean := raw + strings.Join(aggreStrs[stock.Aggregation_MEAN], "\n")
	meanFileName := standard + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\mean.csv"
	err = ioutil.WriteFile(meanFileName, []byte(mean), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v mean failure!", meanFileName)
	}

	std := raw + strings.Join(aggreStrs[stock.Aggregation_STDDEV], "\n")
	stdFileName := standard + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\std.csv"
	err = ioutil.WriteFile(stdFileName, []byte(std), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v std failure!", stdFileName)
	}
	return nil
}

func (svc *WmqyLineService) AsyncStdFile(para interface{}) {
	minmax := (para.([]interface{}))[0].(string)
	standard := (para.([]interface{}))[1].(string)
	startDate := (para.([]interface{}))[2].(int64)
	endDate := (para.([]interface{}))[3].(int64)
	filename := (para.([]interface{}))[4].(string)
	aggreStrs := (para.([]interface{}))[5].(map[stock.AggregationType][]string)
	aggres, err := svc.StdFile(minmax, standard, startDate, endDate, filename)
	if err != nil {
		return
	}

	for key, aggre := range aggres {
		raw := strings.TrimSuffix(filename, ".csv") + ","
		i := 0
		jsonMap, _, _ := stock.GetJsonMap(aggre)
		for _, colname := range WmqylineHeader[3:] {
			fieldname := jsonMap[colname]
			v, _ := reflect.GetValue(aggre, fieldname)
			raw = raw + fmt.Sprint(v)
			if i < len(WmqylineHeader[3:])-1 {
				raw = raw + ","
			}
			i++
		}
		aggreStrs[key] = append(aggreStrs[key], raw)
	}
}

func (svc *WmqyLineService) StdFile(minmax string, standard string, _startDate int64, _endDate int64, filename string) (map[stock.AggregationType]*entity.QPerformance, error) {
	var err error
	ts_code := strings.TrimSuffix(filename, ".csv")
	startDate := stock.GetQTradeDate(_startDate)
	endDate := stock.GetQTradeDate(_endDate)
	wmqyLineMap, err := GetQPerformanceService().FindQPerformance(LinetypeWmqy, ts_code, startDate, endDate)
	if err != nil {
		return nil, err
	}
	wmqyLines, _ := wmqyLineMap[ts_code]
	qps := GetQPerformanceService().Std(wmqyLines, StdtypeMinmax, false)
	if len(qps) < 30 {
		logger.Sugar.Warnf("ts_code:%v dayLines len: %v is less than 30!", ts_code, len(qps))
		return nil, nil
	}
	ps := make([]interface{}, len(qps))
	for qp := range qps {
		ps = append(ps, qp)
	}
	path := string(os.PathSeparator) + fmt.Sprint(_startDate) + "-" + fmt.Sprint(_endDate) + string(os.PathSeparator) + ts_code + ".csv"
	stat := stock.CreateStat(ps, WmqylineHeader[3:])
	stds, minmaxs := stat.CalStd(WmqylineHeader[0:3], false)
	minmaxFileName := minmax + path
	err = ioutil.WriteFile(minmaxFileName, []byte(stock.ToCsv(WmqylineHeader[2:], minmaxs)), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v minmax failure!", minmaxFileName)
	}
	stdFileName := standard + path
	err = ioutil.WriteFile(stdFileName, []byte(stock.ToCsv(WmqylineHeader[2:], stds)), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v std failure!", stdFileName)
	}
	as := make(map[stock.AggregationType]*entity.QPerformance)
	as[stock.Aggregation_SUM] = stat.Sum.(*entity.QPerformance)
	as[stock.Aggregation_MAX] = stat.Max.(*entity.QPerformance)
	as[stock.Aggregation_MIN] = stat.Min.(*entity.QPerformance)
	as[stock.Aggregation_MEAN] = stat.Mean.(*entity.QPerformance)
	as[stock.Aggregation_MEDIAN] = stat.Median[1].(*entity.QPerformance)
	as[stock.Aggregation_STDDEV] = stat.Stddev.(*entity.QPerformance)
	as[stock.Aggregation_RSD] = stat.Rsd.(*entity.QPerformance)

	return as, nil
}
