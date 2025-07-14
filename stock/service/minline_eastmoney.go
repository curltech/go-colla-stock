package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/eastmoney"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/robfig/cron"
	"strconv"
	"strings"
)

// 获取某只股票最新的日期和分钟
func (svc *MinLineService) findMaxTradeDate(tsCode string) (int64, int64, error) {
	cond := &entity.MinLine{}
	cond.TsCode = tsCode
	minLines := make([]*entity.MinLine, 0)
	err := svc.Find(&minLines, cond, "tradedate desc,trademinute desc", 0, 1, "")
	if err != nil {
		return 0, 0, err
	}
	if len(minLines) > 0 {
		return minLines[0].TradeDate, minLines[0].TradeMinute, nil
	}

	return 0, 0, nil
}

// GetMinLine 获取过去的分钟数据
func (svc *MinLineService) GetMinLine(secId string, beg int, limit int, klt int) ([]*entity.MinLine, error) {
	klines, err := GetDayLineService().GetKLine(secId, beg, 0, limit, klt)
	if err != nil {
		return nil, err
	}
	minLines := make([]*entity.MinLine, 0)
	for _, kline := range klines {
		minLine, _ := strToMinLine(secId, kline)
		if minLine != nil && minLine.TradeDate >= int64(beg) {
			minLines = append(minLines, minLine)
		}
	}

	return minLines, err
}

type TodayMinLineResponseData struct {
	Code        string   `json:"code,omitempty"`
	Market      int      `json:"market,omitempty"`
	Name        string   `json:"name,omitempty"`
	Decimal     int      `json:"decimal,omitempty"`     //小数位
	TrendsTotal int      `json:"trendsTotal,omitempty"` //总记录数
	PrePrice    float64  `json:"prePrice,omitempty"`
	PreClose    float64  `json:"preClose,omitempty"`
	Trends      []string `json:"trends,omitempty"` //数据
}

type TodayMinLineResponseResult struct {
	Rc   int                       `json:"rc,omitempty"`
	Rt   int                       `json:"rt,omitempty"`
	Svr  int                       `json:"svr,omitempty"`
	Lt   int                       `json:"lt,omitempty"`
	Full int                       `json:"full,omitempty"`
	Data *TodayMinLineResponseData `json:"data,omitempty"`
}

// GetTodayMinLine 获取当日的分钟数据
func (svc *MinLineService) GetTodayMinLine(secId string) ([]*entity.MinLine, error) {
	params := eastmoney.CreateDayLineRequestParam()
	params.SecId = getSecId(secId)
	params.Fields1 = "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13"
	params.Fields2 = "f51,f52,f53,f54,f55,f56,f57,f58"
	resp, err := eastmoney.TodayFastGet(*params)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	r := &TodayMinLineResponseResult{}
	err = json.Unmarshal(resp, r)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return nil, err
	}
	if r.Data == nil || r.Data.Trends == nil {
		logger.Sugar.Errorf("secId:%v Error: %s", secId, errors.New("NoTrends"))
		return nil, errors.New("NoTrends")
	}
	minLines := make([]*entity.MinLine, 0)
	for _, trend := range r.Data.Trends {
		minLine, _ := strToTodayMinLine(secId, trend)
		if minLine != nil {
			minLines = append(minLines, minLine)
		}
	}

	return minLines, err
}

func strToMinLine(secId string, kline string) (*entity.MinLine, error) {
	kls := strings.Split(kline, ",")
	minLine := &entity.MinLine{}
	minLine.TsCode = secId
	//"trade_date,open,close,high,low,vol,amount,nil"
	tradeDates := strings.Split(kls[0], " ")
	tradeDate, err := strconv.ParseInt(strings.ReplaceAll(tradeDates[0], "-", ""), 10, 64)
	if err != nil {
		logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
		return nil, err
	}
	minLine.TradeDate = tradeDate
	var tradeMinutes []string
	if len(tradeDates) > 1 {
		tradeMinutes = strings.Split(tradeDates[1], ":")
		if len(tradeMinutes) > 1 {
			hour, err := strconv.ParseInt(tradeMinutes[0], 10, 64)
			if err != nil {
				logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
				return nil, err
			}
			minute, err := strconv.ParseInt(tradeMinutes[1], 10, 64)
			if err != nil {
				logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
				return nil, err
			}
			minLine.TradeMinute = hour*60 + minute
		}
	}
	minLine.Open, err = strToFloat(kls[1])
	if err != nil {
		return nil, err
	}
	minLine.Close, err = strToFloat(kls[2])
	if err != nil {
		return nil, err
	}
	minLine.High, err = strToFloat(kls[3])
	if err != nil {
		return nil, err
	}
	minLine.Low, err = strToFloat(kls[4])
	if err != nil {
		return nil, err
	}
	minLine.Vol, err = strToFloat(kls[5])
	if err != nil {
		return nil, err
	}
	minLine.Amount, err = strToFloat(kls[6])
	if err != nil {
		return nil, err
	}
	minLine.Turnover, err = strToFloat(kls[10])
	if err != nil {
		return nil, err
	}
	return minLine, nil
}

func strToTodayMinLine(secId string, kline string) (*entity.MinLine, error) {
	kls := strings.Split(kline, ",")
	minLine := &entity.MinLine{}
	minLine.TsCode = secId
	//"trade_date,open,close,high,low,vol,amount,nil"
	tradeDates := strings.Split(kls[0], " ")
	tradeDate, err := strconv.ParseInt(strings.ReplaceAll(tradeDates[0], "-", ""), 10, 64)
	if err != nil {
		logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
		return nil, err
	}
	minLine.TradeDate = tradeDate
	var tradeMinutes []string
	if len(tradeDates) > 1 {
		tradeMinutes = strings.Split(tradeDates[1], ":")
		if len(tradeMinutes) > 1 {
			hour, err := strconv.ParseInt(tradeMinutes[0], 10, 64)
			if err != nil {
				logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
				return nil, err
			}
			minute, err := strconv.ParseInt(tradeMinutes[1], 10, 64)
			if err != nil {
				logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
				return nil, err
			}
			minLine.TradeMinute = hour*60 + minute
		}
	}
	minLine.Open, err = strToFloat(kls[1])
	if err != nil {
		return nil, err
	}
	minLine.Close, err = strToFloat(kls[2])
	if err != nil {
		return nil, err
	}
	minLine.High, err = strToFloat(kls[3])
	if err != nil {
		return nil, err
	}
	minLine.Low, err = strToFloat(kls[4])
	if err != nil {
		return nil, err
	}
	minLine.Vol, err = strToFloat(kls[5])
	if err != nil {
		return nil, err
	}
	minLine.Amount, err = strToFloat(kls[6])
	if err != nil {
		return nil, err
	}
	minLine.PreClose, err = strToFloat(kls[7])
	if err != nil {
		return nil, err
	}
	return minLine, nil
}

// GetTodayFinanceFlow 获取当日的资金流向分钟数据
func (svc *MinLineService) GetTodayFinanceFlow(secId string) ([]*entity.MinLine, error) {
	params := eastmoney.CreateFinanceFlowRequestParam()
	params.SecId = getSecId(secId)
	params.Fields1 = "f1,f2,f3,f7"
	params.Fields2 = "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65"
	params.Klt = 1
	params.Lmt = 0
	params.Underscore = "1638371480346"
	resp, err := eastmoney.TodayFfFastGet(*params)
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
	minLines := make([]*entity.MinLine, 0)
	for _, kline := range r.Data.Klines {
		minLine, _ := strToMinLineFinanceFlow(secId, kline)
		if minLine != nil {
			minLines = append(minLines, minLine)
		}
	}

	return minLines, err
}

func strToMinLineFinanceFlow(secId string, kline string) (*entity.MinLine, error) {
	minLine := &entity.MinLine{}
	minLine.TsCode = secId
	kls := strings.Split(kline, ",")
	tradeDates := strings.Split(kls[0], " ")
	//"trade_date,主力净流入/净额,小单净流入/净额,中单净流入/净额,大单净流入/净额,超大单净流入/净额"
	tradeDate, err := strconv.ParseInt(strings.ReplaceAll(tradeDates[0], "-", ""), 10, 64)
	if err != nil {
		logger.Sugar.Errorf("tradeDate format error:%v", tradeDates[0])
		return nil, err
	}
	minLine.TradeDate = tradeDate
	var tradeMinutes []string
	if len(tradeDates) > 1 {
		tradeMinutes = strings.Split(tradeDates[1], ":")
		if len(tradeMinutes) > 1 {
			hour, err := strconv.ParseInt(tradeMinutes[0], 10, 64)
			if err != nil {
				logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
				return nil, err
			}
			minute, err := strconv.ParseInt(tradeMinutes[1], 10, 64)
			if err != nil {
				logger.Sugar.Errorf("tradeDate format error:%v", kls[0])
				return nil, err
			}
			minLine.TradeMinute = hour*60 + minute
		}
	}

	minLine.MainNetInflow, err = strToFloat(kls[1])
	if err != nil {
		return nil, err
	}
	minLine.SmallNetInflow, err = strToFloat(kls[2])
	if err != nil {
		return nil, err
	}
	minLine.MiddleNetInflow, err = strToFloat(kls[3])
	if err != nil {
		return nil, err
	}
	minLine.LargeNetInflow, err = strToFloat(kls[4])
	if err != nil {
		return nil, err
	}
	minLine.SuperNetInflow, err = strToFloat(kls[5])
	if err != nil {
		return nil, err
	}

	return minLine, nil
}

func (svc *MinLineService) RefreshMinLine(beg int64) error {
	processLog := GetProcessLogService().StartLog("minline", "RefreshMinLine", "")
	routinePool := thread.CreateRoutinePool(10, svc.AsyncUpdateMinLine, nil)
	defer routinePool.Release()
	tsCodes, _ := GetShareService().GetShareCache()
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

func (svc *MinLineService) AsyncUpdateMinLine(para interface{}) {
	secId := (para.([]interface{}))[0].(string)
	beg := (para.([]interface{}))[1].(int64)
	limit := (para.([]interface{}))[2].(int)
	_, err := svc.GetUpdateMinLine(secId, beg, limit)
	if err != nil {
		return
	}
}

func (svc *MinLineService) GetUpdateMinLine(secId string, beg int64, limit int) ([]*entity.MinLine, error) {
	processLog := GetProcessLogService().StartLog("minline", "GetUpdateMinLine", secId)
	ps, err := svc.UpdateMinLine(secId, beg, limit)
	if err != nil {
		GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	GetProcessLogService().EndLog(processLog, "", "")
	return ps, err
}

func (svc *MinLineService) UpdateMinLine(secId string, beg int64, limit int) ([]*entity.MinLine, error) {
	var minute int64 = 0
	if beg < 0 {
		beg, minute, _ = svc.findMaxTradeDate(secId)
		if beg > 0 {
			if minute >= 900 {
				beg++
				minute = 0
			} else {
				/**
				不删除，增量
				*/
				//minLine := &entity.MinLine{}
				//minLine.TsCode = secId
				//minLine.TradeDate = beg
				//svc.Delete(minLine, "")
			}
		}
	}
	today := stock.CurrentDate()
	if beg > 0 && beg > today {
		return nil, errors.New("data is updated")
	}
	minLines, err := svc.GetMinLine(secId, int(beg), limit, 1)
	if err != nil {
		return nil, errors.New("")
	}
	if len(minLines) <= 0 {
		return minLines, nil
	}
	//找新增加的数据,minute代表已经存在的分钟数据
	i := 0
	for _, minline := range minLines {
		if minute == minline.TradeMinute {
			i++
			break
		}
		i++
	}
	if minute != 0 {
		minLines = minLines[i:]
	}
	if len(minLines) > 0 {
		//对新的数据更新到数据库
		return svc.UpdateTodayFinanceFlow(minLines, secId)
	}

	return minLines, nil
}

func (svc *MinLineService) RefreshTodayMinLine() error {
	processLog := GetProcessLogService().StartLog("minline", "RefreshTodayMinLine", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, svc.AsyncUpdateTodayMinLine, nil)
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

func (svc *MinLineService) AsyncUpdateTodayMinLine(para interface{}) {
	secId := (para.([]interface{}))[0].(string)
	_, err := svc.GetUpdateTodayMinLine(secId)
	if err != nil {
		return
	}
}

func (svc *MinLineService) GetUpdateTodayMinLine(secId string) ([]*entity.MinLine, error) {
	//processLog := GetProcessLogService().StartLog("minline", "GetUpdateTodayMinLine", secId)
	ps, err := svc.UpdateTodayMinLine(secId)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	return ps, err
}

// UpdateTodayMinLine 获取今天的分钟数据，并更新数据库
func (svc *MinLineService) UpdateTodayMinLine(secId string) ([]*entity.MinLine, error) {
	minLines, err := svc.GetTodayMinLine(secId)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return nil, err
	}
	if len(minLines) > 0 {
		minLine := minLines[0]
		//当前已有的最大日期和分钟
		beg, minute, _ := svc.findMaxTradeDate(secId)
		if beg > minLine.TradeDate { //当前有的比下载的大，说明数据已经存在
			logger.Sugar.Errorf("today:%v minline data exist", beg)
			return nil, errors.New("today minline data exist")
		} else if beg == minLine.TradeDate { //当前有的与下载的日期相同，说明当天数据已经存在部分或者全部
			//在下载的数据中寻找最大的分钟的记录位置
			i := 0
			for _, minline := range minLines {
				if minute == minline.TradeMinute {
					i++
					break
				}
				i++
			}
			if minute != 0 {
				minLines = minLines[i:]
			}
			if len(minLines) > 0 {
				//对新的数据更新到数据库
				return svc.UpdateTodayFinanceFlow(minLines, secId)
			}
		} else {
			if len(minLines) > 0 {
				//对新的数据更新到数据库
				return svc.UpdateTodayFinanceFlow(minLines, secId)
			}
		}

		return minLines, nil
	}
	return nil, errors.New("minLines len 0")
}

func (svc *MinLineService) UpdateTodayFinanceFlow(minLines []*entity.MinLine, secId string) ([]*entity.MinLine, error) {
	ffs, err := svc.GetTodayFinanceFlow(secId)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
	}
	mls := make(map[string]*entity.MinLine)
	for _, ff := range ffs {
		mls[ff.TsCode+":"+fmt.Sprint(ff.TradeDate)+":"+fmt.Sprint(ff.TradeMinute)] = ff
	}
	ps := make([]interface{}, 0)
	for _, minLine := range minLines {
		key := minLine.TsCode + ":" + fmt.Sprint(minLine.TradeDate) + ":" + fmt.Sprint(minLine.TradeMinute)
		ff, exist := mls[key]
		if exist {
			minLine.MainNetInflow = ff.MainNetInflow
			minLine.SmallNetInflow = ff.SmallNetInflow
			minLine.MiddleNetInflow = ff.MiddleNetInflow
			minLine.LargeNetInflow = ff.LargeNetInflow
			minLine.SuperNetInflow = ff.SuperNetInflow
		} else {
			//logger.Sugar.Warnf("key:%v not exist", key)
		}
		if !stock.Equal(minLine.Turnover, 0) {
			minLine.ShareNumber = minLine.Amount / minLine.Turnover
		}
		ps = append(ps, minLine)
	}
	_, err = svc.Upsert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error: %s", err.Error())
		return minLines, err
	}

	return minLines, nil
}

func (svc *MinLineService) Cron() *cron.Cron {
	c := cron.New()
	err := c.AddFunc("0 0/10 9-15 * * 1-5", svc.RefreshToday)
	if err != nil {
		return nil
	}
	c.Start()

	return c
}

func (svc *MinLineService) RefreshToday() {
	err := svc.RefreshTodayMinLine()
	if err != nil {
		return
	}
	today := stock.CurrentDate()
	err = GetDayLineService().RefreshDayLine(today)
	if err != nil {
		return
	}
}
