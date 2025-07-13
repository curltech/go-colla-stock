package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/eastmoney"
	"github.com/curltech/go-colla-stock/stock/entity"
	"math"
	"os"
	"strconv"
	"strings"
)

// GetKLine 获取某只股票的分钟线到年线的数据
func (svc *DayLineService) GetKLine(secId string, beg int, end int, limit int, klt int) ([]string, error) {
	params := eastmoney.CreateDayLineRequestParam()
	params.SecId = getSecId(secId)
	params.Fields1 = "f1,f2,f3,f4,f5,f6"
	params.Fields2 = "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61"
	params.Klt = klt
	params.Beg = beg
	if end <= 0 {
		params.End = 20500101
	} else {
		params.End = end
	}
	if limit <= 0 {
		params.Lmt = 10000
	} else {
		params.Lmt = limit
	}
	params.Underscore = "1638513559443"
	resp, err := eastmoney.DayLineFastGet(*params)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	r := &eastmoney.DayLineResponseResult{}
	err = json.Unmarshal(resp, r)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	if r.Data == nil || r.Data.Klines == nil {
		logger.Sugar.Errorf("secId:%v Error: %s", secId, errors.New("NoKlines"))
		return nil, errors.New("NoKlines")
	}

	return r.Data.Klines, nil
}

// GetDayLine 获取某只股票的日线数据
func (svc *DayLineService) GetDayLine(secId string, beg int, end int, limit int, previous *entity.DayLine) ([]*entity.DayLine, error) {
	klines, err := svc.GetKLine(secId, beg, end, limit, 101)
	if err != nil {
		return nil, err
	}
	dayLines := make([]*entity.DayLine, 0)
	for _, kline := range klines {
		dayLine, _ := strToDayLine(secId, kline)
		if dayLine != nil {
			if previous != nil && previous.Open != 0.0 {
				dayLine.PctChgOpen = dayLine.Open/previous.Open - 1
			}
			if previous != nil && previous.High != 0.0 {
				dayLine.PctChgHigh = dayLine.High/previous.High - 1
			}
			if previous != nil && previous.Low != 0.0 {
				dayLine.PctChgLow = dayLine.Low/previous.Low - 1
			}
			if previous != nil && previous.Close != 0.0 {
				dayLine.PctChgClose = dayLine.Close/previous.Close - 1
			}
			if previous != nil && previous.Amount != 0.0 {
				dayLine.PctChgAmount = dayLine.Amount/previous.Amount - 1
			}
			if previous != nil && previous.Vol != 0.0 {
				dayLine.PctChgVol = dayLine.Vol/previous.Vol - 1
			}
			if previous != nil {
				dayLine.PreClose = previous.Close
			}
			previous = dayLine
			dayLines = append(dayLines, dayLine)
		}
	}

	return dayLines, err
}

// 获取某只股票最新的日期
func (svc *DayLineService) findByTradeDate(tsCode string, startDate int64, endDate int64) ([]*entity.DayLine, error) {
	cond := &entity.DayLine{}
	cond.TsCode = tsCode
	dayLines := make([]*entity.DayLine, 0)
	conds := "? <= tradedate"
	paras := make([]interface{}, 0)
	paras = append(paras, startDate)
	if endDate > 0 {
		conds += " and tradedate <= ?"
		paras = append(paras, endDate)
	}
	err := svc.Find(&dayLines, cond, "tradedate", 0, 0, conds, paras...)

	return dayLines, err
}

func strToDayLine(secId string, kline string) (*entity.DayLine, error) {
	kls := strings.Split(kline, ",")
	dayLine := &entity.DayLine{}
	dayLine.TsCode = secId
	//"trade_date,open,close,high,low,vol,amount,nil,pct_chg%,change,turnover%"
	tradeDate, err := strconv.ParseInt(strings.ReplaceAll(kls[0], "-", ""), 10, 64)
	if err != nil {
		logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
		return nil, err
	}
	dayLine.TradeDate = tradeDate
	dayLine.Open, err = strToFloat(kls[1])
	if err != nil {
		return nil, err
	}
	dayLine.Close, err = strToFloat(kls[2])
	if err != nil {
		return nil, err
	}
	dayLine.High, err = strToFloat(kls[3])
	if err != nil {
		return nil, err
	}
	dayLine.Low, err = strToFloat(kls[4])
	if err != nil {
		return nil, err
	}
	dayLine.Vol, err = strToFloat(kls[5])
	if err != nil {
		return nil, err
	}
	dayLine.Amount, err = strToFloat(kls[6])
	if err != nil {
		return nil, err
	}
	pctChg, err := strToFloat(kls[8])
	if err != nil {
		return nil, err
	}
	dayLine.PctChgClose = pctChg
	dayLine.ChgClose, err = strToFloat(kls[9])
	if err != nil {
		return nil, err
	}
	dayLine.Turnover, err = strToFloat(kls[10])
	if err != nil {
		return nil, err
	}

	return dayLine, nil
}

// GetFinanceFlow 获取某只股票的日线资金流动数据，最多6个月的数据
func (svc *DayLineService) GetFinanceFlow(secId string, beg int, limit int) ([]*entity.DayLine, error) {
	params := eastmoney.CreateFinanceFlowRequestParam()
	params.SecId = getSecId(secId)
	params.Fields1 = "f1,f2,f3,f7"
	params.Fields2 = "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65"
	params.Klt = 101
	params.Beg = beg
	params.End = 20500101
	params.Lmt = limit
	params.Underscore = "1638372426494"
	resp, err := eastmoney.FinanceFlowFastGet(*params)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	r := &eastmoney.DayLineResponseResult{}
	err = json.Unmarshal(resp, r)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	if r.Data == nil || r.Data.Klines == nil {
		logger.Sugar.Errorf("secId:%v Error: %s", secId, errors.New("NoKlines"))
		return nil, errors.New("NoKlines")
	}
	dayLines := make([]*entity.DayLine, 0)
	for _, kline := range r.Data.Klines {
		dayLine, _ := strToFinanceFlow(secId, kline)
		if dayLine != nil {
			dayLines = append(dayLines, dayLine)
		}
	}

	return dayLines, err
}

func strToFinanceFlow(secId string, kline string) (*entity.DayLine, error) {
	dayLine := &entity.DayLine{}
	dayLine.TsCode = secId
	kls := strings.Split(kline, ",")
	//"trade_date,主力净流入/净额,小单净流入/净额,中单净流入/净额,大单净流入/净额,超大单净流入/净额,主力净流入/净占比%,
	//小单净流入/净占比%,中单净流入/净占比%,大单净流入/净占比%,超大单净流入/净占比%,close,pct_chg,,"
	tradeDate, err := strconv.ParseInt(strings.ReplaceAll(kls[0], "-", ""), 10, 64)
	if err != nil {
		logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
		return nil, err
	}
	dayLine.TradeDate = tradeDate

	dayLine.MainNetInflow, err = strToFloat(kls[1])
	if err != nil {
		return nil, err
	}
	dayLine.SmallNetInflow, err = strToFloat(kls[2])
	if err != nil {
		return nil, err
	}
	dayLine.MiddleNetInflow, err = strToFloat(kls[3])
	if err != nil {
		return nil, err
	}
	dayLine.LargeNetInflow, err = strToFloat(kls[4])
	if err != nil {
		return nil, err
	}
	dayLine.SuperNetInflow, err = strToFloat(kls[5])
	if err != nil {
		return nil, err
	}
	dayLine.PctMainNetInflow, err = strToFloat(kls[6])
	if err != nil {
		return nil, err
	}
	dayLine.PctSmallNetInflow, err = strToFloat(kls[7])
	if err != nil {
		return nil, err
	}
	dayLine.PctMiddleNetInflow, err = strToFloat(kls[8])
	if err != nil {
		return nil, err
	}
	dayLine.PctLargeNetInflow, err = strToFloat(kls[9])
	if err != nil {
		return nil, err
	}
	dayLine.PctSuperNetInflow, err = strToFloat(kls[10])
	if err != nil {
		return nil, err
	}

	return dayLine, nil
}

func strToFloat(value string) (float64, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logger.Sugar.Errorf("string format error:%v", value)
		return 0, err
	}
	return f, nil
}

// RefreshDayLine 刷新所有股票的日线数据，beg为负数的时候从已有的最新数据开始更新
func (svc *DayLineService) RefreshDayLine(beg int64) error {
	processLog := GetProcessLogService().StartLog("dayline", "RefreshDayLine", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, svc.AsyncUpdateDayLine, nil)
	defer routinePool.Release()
	tsCodes, _ := GetShareService().GetCacheShare()
	for _, tsCode := range tsCodes {
		para := make([]interface{}, 0)
		para = append(para, tsCode)
		para = append(para, beg)
		para = append(para, 10000)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(processLog, "", "")

	return nil
}

func (svc *DayLineService) AsyncUpdateDayLine(para interface{}) {
	secId := (para.([]interface{}))[0].(string)
	beg := (para.([]interface{}))[1].(int64)
	limit := (para.([]interface{}))[2].(int)

	_, err := svc.GetUpdateDayline(secId, beg, limit)
	if err != nil {
		return
	}
}

// GetUpdateDayline 当天15点之前缺当天资金流数据
func (svc *DayLineService) GetUpdateDayline(secId string, beg int64, limit int) ([]*entity.DayLine, error) {
	processLog := GetProcessLogService().StartLog("dayline", "GetUpdateDayline", secId)
	ps, err := svc.UpdateDayline(secId, beg, limit)
	if err != nil {
		GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	if len(ps) > 0 {
		_, err = svc.UpdateStat(secId, 0)
		if err != nil {
			return nil, err
		}
	}
	return ps, err
}

func (svc *DayLineService) deleteDayline(secId string, beg int64) error {
	dayline := &entity.DayLine{}
	dayline.TsCode = secId
	dayline.TradeDate = beg
	_, err := svc.Delete(dayline, "")

	return err
}

func (svc *DayLineService) UpdateDayline(secId string, beg int64, limit int) ([]*entity.DayLine, error) {
	previous, err := svc.findMaxTradeDate(secId)
	var pre *entity.DayLine = nil
	if previous == nil || len(previous) == 0 {
		beg = 0
	} else if previous[0] != nil {
		pre = previous[0]
		if beg < 0 {
			beg = stock.AddDay(previous[0].TradeDate, 1)
		}
	}
	today := stock.CurrentDate()
	if beg > 0 && beg > today {
		return nil, errors.New("data is updated")
	}

	dayLines, err := svc.GetDayLine(secId, int(beg), 0, limit, pre)
	if err != nil {
		return nil, err
	}
	if len(dayLines) <= 0 {
		return nil, errors.New("dayLines len is 0")
	}
	ps, err := svc.UpdateFinanceFlow(dayLines, secId, int(beg), limit)
	if err != nil {
		return ps, err
	}

	return ps, err
}

func (svc *DayLineService) UpdateFinanceFlow(dayLines []*entity.DayLine, secId string, beg int, limit int) ([]*entity.DayLine, error) {
	ffs, err := svc.GetFinanceFlow(secId, beg, limit)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
	}
	dls := make(map[string]*entity.DayLine)
	for _, ff := range ffs {
		dls[ff.TsCode+":"+fmt.Sprint(ff.TradeDate)] = ff
	}
	ps := make([]interface{}, 0)
	for _, dayLine := range dayLines {
		key := dayLine.TsCode + ":" + fmt.Sprint(dayLine.TradeDate)
		ff, exist := dls[key]
		if exist {
			dayLine.MainNetInflow = ff.MainNetInflow
			dayLine.PctMainNetInflow = ff.PctMainNetInflow
			dayLine.SmallNetInflow = ff.SmallNetInflow
			dayLine.PctSmallNetInflow = ff.PctSmallNetInflow
			dayLine.MiddleNetInflow = ff.MiddleNetInflow
			dayLine.PctMiddleNetInflow = ff.PctMiddleNetInflow
			dayLine.LargeNetInflow = ff.LargeNetInflow
			dayLine.PctLargeNetInflow = ff.PctLargeNetInflow
			dayLine.SuperNetInflow = ff.SuperNetInflow
			dayLine.PctSuperNetInflow = ff.PctSuperNetInflow
		} else {
			//logger.Sugar.Warnf("key:%v not exist", key)
		}
		if !stock.Equal(dayLine.Turnover, 0) {
			dayLine.ShareNumber = dayLine.Amount / dayLine.Turnover
		}
		if stock.Equal(dayLine.PreClose, 0) {
			if !stock.Equal(dayLine.ChgClose, 0) {
				dayLine.PreClose = dayLine.Close - dayLine.ChgClose
			} else {
				dayLine.PreClose = dayLine.Close
			}
		}
		ps = append(ps, dayLine)
	}
	_, err = svc.Upsert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return dayLines, err
	}

	return dayLines, nil
}

func getSecId(secId string) string {
	_, shares := GetShareService().GetCacheShare()
	share, exist := shares[secId]
	if exist {
		if strings.HasSuffix(share.Symbol, ".SH") {
			return "1." + secId
		}
	}
	return "0." + secId
}

func (svc *DayLineService) findAggregation(startDate int64, endDate int64) (map[stock.AggregationType]map[string]*entity.DayLine, error) {
	sql := "select tscode TsCode"
	jsonMap, _, _ := stock.GetJsonMap(entity.DayLine{})
	for _, colname := range DaylineHeader[2:] {
		fieldname, ok := jsonMap[colname]
		if ok {
			lowerFieldname := strings.ToLower(fieldname)
			sql = sql + ",max(" + lowerFieldname + ")" + " max_" + fieldname
			sql = sql + ",min(" + lowerFieldname + ")" + " min_" + fieldname
			sql = sql + ",avg(" + lowerFieldname + ")" + " avg_" + fieldname
			sql = sql + ",stddev(" + lowerFieldname + ")" + " stddev_" + fieldname
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
	sql = sql + " group by tscode"
	results, err := svc.Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	aggreDaylines := make(map[stock.AggregationType]map[string]*entity.DayLine)
	aggreDaylines[stock.Aggregation_MAX] = make(map[string]*entity.DayLine)
	aggreDaylines[stock.Aggregation_MIN] = make(map[string]*entity.DayLine)
	aggreDaylines[stock.Aggregation_MEAN] = make(map[string]*entity.DayLine)
	aggreDaylines[stock.Aggregation_STDDEV] = make(map[string]*entity.DayLine)
	for _, result := range results {
		maxDayline := &entity.DayLine{}
		minDayline := &entity.DayLine{}
		meanDayline := &entity.DayLine{}
		stddevDayline := &entity.DayLine{}
		for colname, v := range result {
			if colname == "TsCode" {
				err = reflect.Set(maxDayline, "TsCode", string(v))
				if err != nil {
					return nil, err
				}
			} else {
				fieldnames := strings.Split(colname, "_")
				if fieldnames[0] == "max" {
					err = reflect.Set(maxDayline, fieldnames[1], string(v))
					if err != nil {
						return nil, err
					}
				}
				if fieldnames[0] == "min" {
					err = reflect.Set(minDayline, fieldnames[1], string(v))
					if err != nil {
						return nil, err
					}
				}
				if fieldnames[0] == "mean" {
					err = reflect.Set(meanDayline, fieldnames[1], string(v))
					if err != nil {
						return nil, err
					}
				}
				if fieldnames[0] == "stddev" {
					err = reflect.Set(stddevDayline, fieldnames[1], string(v))
					if err != nil {
						return nil, err
					}
				}
			}
		}
		aggreDaylines[stock.Aggregation_MAX][maxDayline.TsCode] = maxDayline
		aggreDaylines[stock.Aggregation_MIN][minDayline.TsCode] = minDayline
		aggreDaylines[stock.Aggregation_MEAN][meanDayline.TsCode] = meanDayline
		aggreDaylines[stock.Aggregation_STDDEV][stddevDayline.TsCode] = stddevDayline
	}

	return aggreDaylines, nil
}

// StdPath src: "C:\stock\data\origin\lday"
// minmax: "C:\stock\data\minmax\lday"
// standard: "C:\stock\data\standard\lday"
func (svc *DayLineService) StdPath(src string, minmax string, standard string, startDate int64, endDate int64) error {
	err := stock.Mkdir(minmax + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate))
	if err != nil {
		return err
	}
	err = stock.Mkdir(standard + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate))
	if err != nil {
		return err
	}
	aggreStrs := make(map[stock.AggregationType][]string)
	for _, typ := range stock.AggregationTypes {
		aggreStrs[typ] = make([]string, 0)
	}

	routinePool := thread.CreateRoutinePool(10, svc.AsyncStdFile, nil)
	defer routinePool.Release()
	if src != "" {
		files, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, file := range files {
			filename := file.Name()
			hasSuffix := strings.HasSuffix(filename, ".csv")
			if hasSuffix {
				para := make([]interface{}, 0)
				para = append(para, src)
				para = append(para, minmax)
				para = append(para, standard)
				para = append(para, startDate)
				para = append(para, endDate)
				para = append(para, filename)
				para = append(para, aggreStrs)
				routinePool.Invoke(para)
			}
		}
	} else {
		tsCodes, _ := GetShareService().GetCacheShare()
		for _, tsCode := range tsCodes {
			para := make([]interface{}, 0)
			para = append(para, src)
			para = append(para, minmax)
			para = append(para, standard)
			para = append(para, startDate)
			para = append(para, endDate)
			para = append(para, tsCode)
			para = append(para, aggreStrs)
			routinePool.Invoke(para)
		}
	}
	routinePool.Wait(nil)
	raw := "ts_code," + strings.Join(DaylineHeader[2:], ",") + "\n"
	max := raw + strings.Join(aggreStrs[stock.Aggregation_MAX], "\n")
	maxFileName := minmax + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\max.csv"
	err = os.WriteFile(maxFileName, []byte(max), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v max failure!", maxFileName)
	} else {
		logger.Sugar.Infof("%v max completely!", maxFileName)
	}

	min := raw + strings.Join(aggreStrs[stock.Aggregation_MIN], "\n")
	minFileName := minmax + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\min.csv"
	err = os.WriteFile(minFileName, []byte(min), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v min failure!", minFileName)
	} else {
		logger.Sugar.Infof("%v min completely!", minFileName)
	}

	mean := raw + strings.Join(aggreStrs[stock.Aggregation_MEAN], "\n")
	meanFileName := standard + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\mean.csv"
	err = os.WriteFile(meanFileName, []byte(mean), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v mean failure!", meanFileName)
	} else {
		logger.Sugar.Infof("%v mean completely!", meanFileName)
	}

	std := raw + strings.Join(aggreStrs[stock.Aggregation_STDDEV], "\n")
	stdFileName := standard + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + "\\std.csv"
	err = os.WriteFile(stdFileName, []byte(std), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v std failure!", stdFileName)
	} else {
		logger.Sugar.Infof("%v std completely!", stdFileName)
	}
	return nil
}

func (svc *DayLineService) AsyncStdFile(para interface{}) {
	//src := (para.([]interface{}))[0].(string)
	minmax := (para.([]interface{}))[1].(string)
	standard := (para.([]interface{}))[2].(string)
	startDate := (para.([]interface{}))[3].(int64)
	endDate := (para.([]interface{}))[4].(int64)
	filename := (para.([]interface{}))[5].(string)
	aggreStrs := (para.([]interface{}))[6].(map[stock.AggregationType][]string)
	aggres, err := svc.StdFile(startDate, endDate, filename, minmax, standard)
	if err != nil {
		return
	}

	for key, aggre := range aggres {
		raw := strings.TrimSuffix(filename, ".csv") + ","
		i := 0
		jsonMap, _, _ := stock.GetJsonMap(aggre)
		for _, colname := range DaylineHeader[2:] {
			fieldname := jsonMap[colname]
			v, _ := reflect.GetValue(aggre, fieldname)
			raw = raw + fmt.Sprint(v)
			if i < len(DaylineHeader[2:])-1 {
				raw = raw + ","
			}
			i++
		}
		aggreStrs[key] = append(aggreStrs[key], raw)
	}
}

func (svc *DayLineService) LoadFile(src string, startDate int64, endDate int64, filename string) ([]*entity.DayLine, error) {
	c, err := os.ReadFile(src + string(os.PathSeparator) + filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(c), "\n")
	i := 0
	var head []string
	dayLines := make([]*entity.DayLine, 0)
	jsonMap, typMap, _ := stock.GetJsonMap(entity.DayLine{})
	for _, line := range lines {
		if line == "" {
			continue
		}
		if i == 0 {
			head = strings.Split(line, ",")
		} else {
			values := strings.Split(line, ",")
			j := 0
			var dayLine *entity.DayLine
			for _, colname := range head {
				fieldname := jsonMap[colname]
				if fieldname == "" {
					j++
					continue
				}
				fieldtyp := typMap[colname]
				if fieldtyp == "" {
					j++
					continue
				}
				if colname == "trade_date" {
					v, _ := convert.ToObject(values[j], fieldtyp)
					tradeDate := v.(int64)
					if tradeDate >= startDate && tradeDate <= endDate {
						dayLine = &entity.DayLine{}
						err = reflect.SetValue(dayLine, fieldname, tradeDate)
						if err != nil {
							return nil, err
						}
						dayLines = append(dayLines, dayLine)
					}
				} else {
					if dayLine != nil {
						val, _ := convert.ToObject(values[j], fieldtyp)
						err = reflect.SetValue(dayLine, fieldname, val)
						if err != nil {
							return nil, err
						}
					}
				}
				j++
			}
		}
		i++
	}
	//df := dataframe.ReadCSV(file, dataframe.DetectTypes(true),
	//	dataframe.DefaultType(series.Float),
	//	dataframe.WithTypes(map[string]series.Type{
	//		"trade_date":   series.Int,
	//		"trade_minute": series.Int,
	//	}))
	//df = df.FilterAggregation(
	//	dataframe.And,
	//	dataframe.F{1, "trade_date", series.GreaterEq, startDate},
	//	dataframe.F{1, "trade_date", series.LessEq, endDate},
	//)
	//dayLines := make([]*entity.DayLine, 0)
	//for i := 0; i < df.Nrow(); i++ {
	//	dayLine := &entity.DayLine{}
	//	for _, colname := range DaylineHeader[2:] {
	//		val := df.Col(colname).Elem(i).Float()
	//		reflect.SetValue(dayLine, colname, val)
	//		dayLines = append(dayLines, dayLine)
	//	}
	//}

	return dayLines, nil
}

func (svc *DayLineService) StdFile(startDate int64, endDate int64, filename string, minmax string, standard string) (map[stock.AggregationType]*entity.DayLine, error) {
	var dayLines []*entity.DayLine
	var err error
	tsCode := strings.TrimSuffix(filename, ".csv")
	if strings.HasSuffix(filename, ".csv") {
		dayLines, err = svc.LoadFile("", startDate, endDate, filename)
		if err != nil {
			return nil, err
		}
	} else {
		dayLines, err = GetDayLineService().findByTradeDate(tsCode, startDate, endDate)
		if err != nil {
			return nil, err
		}
	}
	if len(dayLines) < 100 {
		logger.Sugar.Warnf("ts_code:%v dayLines len: %v is less than 100!", tsCode, len(dayLines))
		return nil, nil
	}
	path := string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(endDate) + string(os.PathSeparator) + tsCode + ".csv"
	aggres, stds, minmaxs := svc.Aggregation(dayLines)
	minmaxFileName := minmax + path
	err = os.WriteFile(minmaxFileName, []byte(stock.ToCsv(DaylineHeader[1:], minmaxs)), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v minmax failure!", minmaxFileName)
	} else {
		logger.Sugar.Infof("%v minmax completely!", minmaxFileName)
	}
	stdFileName := standard + path
	err = os.WriteFile(stdFileName, []byte(stock.ToCsv(DaylineHeader[1:], stds)), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v std failure!", stdFileName)
	} else {
		logger.Sugar.Infof("%v std completely!", stdFileName)
	}
	return aggres, nil
}

func (svc *DayLineService) Aggregation(dayLines []*entity.DayLine) (map[stock.AggregationType]*entity.DayLine, []interface{}, []interface{}) {
	colnames := DaylineHeader[2:]
	aggres := make(map[stock.AggregationType]*entity.DayLine)
	for _, typ := range stock.AggregationTypes {
		aggres[typ] = &entity.DayLine{}
	}
	jsonMap, typMap, _ := stock.GetJsonMap(entity.DayLine{})
	stds := make([]interface{}, 0)
	minmaxs := make([]interface{}, 0)
	i := 0
	for _, dayLine := range dayLines {
		for _, colname := range colnames {
			fieldname := jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := typMap[colname]
			if fieldtyp == "" {
				continue
			}
			v, _ := reflect.GetValue(dayLine, fieldname)
			val := v.(float64)
			for _, typ := range stock.AggregationTypes {
				if typ == stock.Aggregation_STDDEV {
					continue
				}
				aggre := aggres[typ]
				v, _ = reflect.GetValue(aggre, fieldname)
				aggreVal := v.(float64)
				if i == 0 {
					err := reflect.SetValue(aggre, fieldname, val)
					if err != nil {
						return nil, nil, nil
					}
				}
				switch typ {
				case stock.Aggregation_MAX:
					if i != 0 {
						if val > aggreVal {
							err := reflect.SetValue(aggre, fieldname, val)
							if err != nil {
								return nil, nil, nil
							}
						}
					}
				case stock.Aggregation_MIN:
					if i != 0 {
						if val < aggreVal {
							err := reflect.SetValue(aggre, fieldname, val)
							if err != nil {
								return nil, nil, nil
							}
						}
					}
				case stock.Aggregation_SUM:
					if i != 0 {
						err := reflect.SetValue(aggre, fieldname, aggreVal+val)
						if err != nil {
							return nil, nil, nil
						}
					}
				case stock.Aggregation_COUNT:
					if i == 0 {
						err := reflect.SetValue(aggre, fieldname, 1.0)
						if err != nil {
							return nil, nil, nil
						}
					} else {
						err := reflect.SetValue(aggre, fieldname, aggreVal+1)
						if err != nil {
							return nil, nil, nil
						}
					}
				default:
				}
			}
		}
		i++
	}
	max, ok := aggres[stock.Aggregation_MAX]
	min, ok := aggres[stock.Aggregation_MIN]
	sum, ok := aggres[stock.Aggregation_SUM]
	count, ok := aggres[stock.Aggregation_COUNT]
	mean, ok := aggres[stock.Aggregation_MEAN]
	for _, colname := range colnames {
		fieldname := jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := typMap[colname]
		if fieldtyp == "" {
			continue
		}
		v, _ := reflect.GetValue(sum, fieldname)
		sumVal := v.(float64)
		v, _ = reflect.GetValue(count, fieldname)
		countVal := v.(float64)
		if countVal != 0.0 {
			err := reflect.SetValue(mean, fieldname, sumVal/countVal)
			if err != nil {
				return nil, nil, nil
			}
		}
	}
	stddev, ok := aggres[stock.Aggregation_STDDEV]
	if ok {
		for _, dayLine := range dayLines {
			for _, colname := range colnames {
				fieldname := jsonMap[colname]
				if fieldname == "" {
					continue
				}
				fieldtyp := typMap[colname]
				if fieldtyp == "" {
					continue
				}
				val := 0.0
				v, _ := reflect.GetValue(dayLine, fieldname)
				val = v.(float64)
				v, _ = reflect.GetValue(mean, fieldname)
				meanVal := v.(float64)
				diff := val - meanVal
				v, _ = reflect.GetValue(stddev, fieldname)
				stddevVal := v.(float64)
				err := reflect.SetValue(stddev, fieldname, stddevVal+diff*diff)
				if err != nil {
					return nil, nil, nil
				}
			}
		}
		for _, colname := range colnames {
			fieldname := jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := typMap[colname]
			if fieldtyp == "" {
				continue
			}
			v, _ := reflect.GetValue(stddev, fieldname)
			stddevVal := v.(float64)
			v, _ = reflect.GetValue(count, fieldname)
			countVal := v.(float64)
			if countVal != 1.0 {
				stddevVal = stddevVal / (countVal - 1.0)
				err := reflect.SetValue(stddev, fieldname, math.Sqrt(stddevVal))
				if err != nil {
					return nil, nil, nil
				}
			}
		}
	}
	for _, dayLine := range dayLines {
		id := dayLine.Id
		tradeDate := dayLine.TradeDate
		std := &entity.DayLine{}
		minmax := &entity.DayLine{}
		std.Id = id
		std.TradeDate = tradeDate
		minmax.Id = id
		minmax.TradeDate = tradeDate
		for _, colname := range colnames {
			fieldname := jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := typMap[colname]
			if fieldtyp == "" {
				continue
			}
			val := 0.0
			v, _ := reflect.GetValue(dayLine, fieldname)
			val = v.(float64)
			v, _ = reflect.GetValue(mean, fieldname)
			meanVal := v.(float64)
			v, _ = reflect.GetValue(stddev, fieldname)
			stddevVal := v.(float64)
			if stddevVal != 0 {
				err := reflect.SetValue(std, fieldname, (val-meanVal)/stddevVal)
				if err != nil {
					return nil, nil, nil
				}
			}
			v, _ = reflect.GetValue(min, fieldname)
			minVal := v.(float64)
			v, _ = reflect.GetValue(max, fieldname)
			maxVal := v.(float64)
			if maxVal != minVal {
				err := reflect.SetValue(minmax, fieldname, (val-minVal)/(maxVal-minVal))
				if err != nil {
					return nil, nil, nil
				}
			}
		}
		stds = append(stds, std)
		minmaxs = append(minmaxs, minmax)
	}

	return aggres, stds, minmaxs
}
