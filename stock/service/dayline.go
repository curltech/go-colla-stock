package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"io/ioutil"
	"os"
	"strings"
)

/**
同步表结构，服务继承基本服务的方法
*/
type DayLineService struct {
	service.OrmBaseService
}

var dayLineService = &DayLineService{}

func GetDayLineService() *DayLineService {
	return dayLineService
}

func (this *DayLineService) GetSeqName() string {
	return seqname
}

func (this *DayLineService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Share{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *DayLineService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.DayLine, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

/**
读目录下的数据
*/
func (this *DayLineService) ParsePath(src string, target string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	routinePool := thread.CreateRoutinePool(10, this.AsyncParseFile, nil)
	defer routinePool.Release()
	for _, file := range files {
		filename := file.Name()
		hasSuffix := strings.HasSuffix(filename, ".day")
		if hasSuffix {
			para := make([]string, 0)
			para = append(para, src)
			para = append(para, target)
			para = append(para, filename)
			routinePool.Invoke(para)
		}
	}
	routinePool.Wait(nil)
	stock.Rename(src, src+"-"+fmt.Sprint(stock.CurrentDate()))
	stock.Mkdir(src)
	return nil
}

func (this *DayLineService) AsyncParseFile(para interface{}) {
	src := (para.([]string))[0]
	target := (para.([]string))[1]
	filename := (para.([]string))[2]
	this.ParseFile(src, target, filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

func (this *DayLineService) ParseFile(src string, target string, filename string, flag int) error {
	shareId := strings.TrimSuffix(filename, ".day")
	logger.Sugar.Infof("shareId:%v", shareId)
	content, err := ioutil.ReadFile(src + string(os.PathSeparator) + filename)
	if err != nil {
		return err
	}
	targetFileName := target + string(os.PathSeparator) + shareId + ".csv"
	dayLines := this.ParseByte(shareId, content)
	raw := stock.ToCsv(DaylineHeader[2:], dayLines)
	file, err := os.OpenFile(targetFileName, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write([]byte(raw))
	logger.Sugar.Infof("Parse day file %v record %v completely!", targetFileName, len(dayLines))
	this.batchSave(dayLines)

	return nil
}

var DaylineHeader = []string{"id", "trade_date", "open", "high", "low", "close", "amount", "vol", "turnover",
	"main_net_inflow", "small_net_inflow", "middle_net_inflow", "large_net_inflow", "super_net_inflow",
	"pct_chg_open", "pct_chg_high", "pct_chg_low", "pct_chg_close", "pct_chg_amount", "pct_chg_vol",
	"pct_main_net_inflow", "pct_small_net_inflow", "pct_middle_net_inflow", "pct_large_net_inflow", "pct_super_net_inflow"}

func (this *DayLineService) batchSave(dayLines []interface{}) error {
	batch := 1000
	dls := make([]interface{}, 0)
	for i := 0; i < len(dayLines); i = i + batch {
		for j := 0; j < batch; j++ {
			if i+j < len(dayLines) {
				dayLine := dayLines[i+j]
				dls = append(dls, dayLine)
			}
		}
		_, err := this.Insert(dls...)
		if err != nil {
			logger.Sugar.Errorf("Insert database error:%v", err.Error())
			return err
		} else {
			logger.Sugar.Infof("Insert database record:%v", len(dls))
		}
		dls = make([]interface{}, 0)
	}

	return nil
}

func (this *DayLineService) ParseByte(shareId string, content []byte) []interface{} {
	dayLines := make([]interface{}, 0)
	var previous *entity.DayLine = nil
	for i := 0; i < len(content); i = i + 32 {
		dayLine := entity.DayLine{}
		dayLine.TsCode = shareId
		dayLine.TradeDate = stock.BytesToInt64(content[i : i+4])
		dayLine.Open = float64(stock.BytesToInt64(content[i+4:i+8])) / 100
		if previous != nil && previous.Open != 0.0 {
			dayLine.PctChgOpen = dayLine.Open/previous.Open - 1
		}
		dayLine.High = float64(stock.BytesToInt64(content[i+8:i+12])) / 100
		if previous != nil && previous.High != 0.0 {
			dayLine.PctChgHigh = dayLine.High/previous.High - 1
		}
		dayLine.Low = float64(stock.BytesToInt64(content[i+12:i+16])) / 100
		if previous != nil && previous.Low != 0.0 {
			dayLine.PctChgLow = dayLine.Low/previous.Low - 1
		}
		dayLine.Close = float64(stock.BytesToInt64(content[i+16:i+20])) / 100
		if previous != nil && previous.Close != 0.0 {
			dayLine.PctChgClose = dayLine.Close/previous.Close - 1
		}
		dayLine.Amount = stock.BytesToFloat64(content[i+20 : i+24])
		if previous != nil && previous.Amount != 0.0 {
			dayLine.PctChgAmount = dayLine.Amount/previous.Amount - 1
		}
		dayLine.Vol = float64(stock.BytesToInt64(content[i+24:i+28])) / 100
		if previous != nil && previous.Vol != 0.0 {
			dayLine.PctChgVol = dayLine.Vol/previous.Vol - 1
		}
		dayLines = append(dayLines, &dayLine)
		previous = &dayLine
	}
	return dayLines
}

func (this *DayLineService) findMaxTradeDate(ts_code string) ([]*entity.DayLine, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	err := this.Find(&dayLines, nil, "tradedate desc", 0, 2, conds, paras...)
	if err != nil {
		return nil, err
	}
	if len(dayLines) > 0 {
		return dayLines, nil
	}

	return nil, nil
}

func (this *DayLineService) findMaxMaTradeDate(ts_code string, fieldname string) (*entity.DayLine, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	conds = conds + " and " + fieldname + "!=0 and " + fieldname + " is not null"
	dayLines := make([]*entity.DayLine, 0)
	err := this.Find(&dayLines, nil, "tradedate desc", 0, 4, conds, paras...)
	if err != nil {
		return nil, err
	}
	if len(dayLines) > 0 {
		return dayLines[0], nil
	}

	return nil, nil
}

func (this *DayLineService) Search(ts_code string, industry string, sector string, startDate int64, endDate int64, orderby string, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	if industry != "" {
		conds += " and tscode in (select tscode from stk_share where industry = ?)"
		paras = append(paras, industry)
	}
	if sector != "" {
		conds += " and tscode in (select tscode from stk_share where sector = ?)"
		paras = append(paras, sector)
	}
	if startDate != 0 {
		conds = conds + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	if endDate != 0 {
		conds = conds + " and tradedate<=?"
		paras = append(paras, endDate)
	}
	var err error
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "tscode,tradedate"
	}
	err = this.Find(&dayLines, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return dayLines, count, nil
}

var dayCounts = []string{"1", "3", "5", "10", "13", "20", "21", "30", "34", "55", "60", "90", "120", "144", "233", "240"}

/**
获取某时间点前limit条数据，如果没有日期范围的指定，就是返回最新的回溯limit条数据
*/
func (this *DayLineService) FindPreceding(ts_code string, endDate int64, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	conds += " and ma3close is not null and ma3close!=0"
	if endDate != 0 {
		conds += " and tradedate<=?"
		paras = append(paras, endDate)
	}
	var err error
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	err = this.Find(&dayLines, nil, "tscode,tradedate desc", from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	length := len(dayLines)
	ps := make([]*entity.DayLine, length)
	for i := length; i > 0; i-- {
		ps[length-i] = dayLines[i-1]
	}
	if len(ps) > 0 {
		logger.Sugar.Infof("from %v to %v datline data", ps[0].TradeDate, ps[len(ps)-1].TradeDate)
	} else {
		logger.Sugar.Errorf("dayline len 0")
	}
	return ps, count, nil
}

/*
获取某时间点后limit条数据，如果没有日期范围的指定，就是返回最早limit条数据
*/
func (this *DayLineService) FindFollowing(ts_code string, startDate int64, endDate int64, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	conds += " and ma3close is not null and ma3close!=0"
	if startDate != 0 {
		conds = conds + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	if endDate != 0 {
		conds = conds + " and tradedate<=?"
		paras = append(paras, endDate)
	}
	var err error
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	err = this.Find(&dayLines, nil, "tscode,tradedate", from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	if len(dayLines) > 0 {
		logger.Sugar.Infof("from %v to %v datline data", dayLines[0].TradeDate, dayLines[len(dayLines)-1].TradeDate)
	}
	return dayLines, count, nil
}

/*
获取某时间点前后limit条数据
*/
func (this *DayLineService) FindRange(ts_code string, startDate int64, endDate int64, limit int) ([]*entity.DayLine, error) {
	preceding, _, err := this.FindPreceding(ts_code, startDate, 0, limit, 0)
	if err != nil {
		return nil, err
	}
	var daylines []*entity.DayLine
	if endDate != 0 && endDate > startDate {
		daylines, _, err = this.FindFollowing(ts_code, startDate, endDate, 0, 0, 0)
		if err != nil {
			return nil, err
		}
	} else {
		endDate = startDate
	}
	following, _, err := this.FindFollowing(ts_code, endDate, 0, 0, limit, 0)
	if err != nil {
		return nil, err
	}
	if daylines != nil && len(daylines) > 0 {
		if daylines[0].TradeDate == startDate {
			daylines = daylines[1:]
		}
		preceding = append(preceding, daylines...)
	}
	if len(following) > 0 && following[0].TradeDate == endDate {
		following = following[1:]
	}
	preceding = append(preceding, following...)

	return preceding, nil
}

/*
获取某时间点后，前后dayCount内的是最高点的，limit条数据
*/
func (this *DayLineService) FindHighest(ts_code string, dayCount string, startDate int64, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	if dayCount == "" {
		return dayLines, count, errors.New("")
	}
	for _, _dayCount := range dayCounts {
		if dayCount == "1" {
			conds += " and chgclose>0"
		} else {
			conds += " and acc" + _dayCount + "pctchgclose>0"
		}
		conds += " and future" + _dayCount + "pctchgclose<0"
		if _dayCount == dayCount {
			break
		}
	}
	conds += " and ma3close is not null and ma3close!=0"
	if startDate != 0 {
		conds = conds + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	var err error
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	err = this.Find(&dayLines, nil, "tscode,tradedate", from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return dayLines, count, nil
}

/*
获取某时间点后，前后dayCount内的是最高点的，limit条数据
*/
func (this *DayLineService) FindLowest(ts_code string, dayCount string, startDate int64, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	if dayCount == "" {
		return dayLines, count, errors.New("")
	}
	for _, _dayCount := range dayCounts {
		if dayCount == "1" {
			conds += " and chgclose<0"
		} else {
			conds += " and acc" + _dayCount + "pctchgclose<0"
		}
		conds += " and future" + _dayCount + "pctchgclose>0"
		if _dayCount == dayCount {
			break
		}
	}
	conds += " and ma3close is not null and ma3close!=0"
	if startDate != 0 {
		conds = conds + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	var err error
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	err = this.Find(&dayLines, nil, "tscode,tradedate", from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return dayLines, count, nil
}

/*
获取某时间点后，两dayCount均线交叉的，limit条数据
*/
func (this *DayLineService) FindMaCross(ts_code string, srcDayCount string, targetDayCount string, startDate int64, cross string, from int, limit int, count int64) ([]*entity.DayLine, int64, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	if srcDayCount == "1" || srcDayCount == "3" || targetDayCount == "1" || targetDayCount == "3" {
		return dayLines, count, errors.New("")
	}
	for _, _srcDayCount := range dayCounts[2:] {
		if srcDayCount != "" {
			_srcDayCount = srcDayCount
		}
		for _, _targetDayCount := range dayCounts[2:] {
			if targetDayCount != "" {
				_targetDayCount = targetDayCount
			}
			if _srcDayCount == _targetDayCount {
				continue
			}
			if cross == "up" {
				conds += " and ma" + _srcDayCount + "close>ma" + _targetDayCount + "close"
				conds += " and (ma" + _srcDayCount + "close-close/(acc" + _srcDayCount + "pctchgclose+1)+close*(1+future1pctchgclose))/" + _srcDayCount + "<(ma" + _targetDayCount + "close-close/(acc" + _targetDayCount + "pctchgclose+1)+close*(1+future1pctchgclose))/" + _targetDayCount
			} else {
				conds += " and ma" + _srcDayCount + "close<ma" + _targetDayCount + "close"
				conds += " and (ma" + _srcDayCount + "close-close/(acc" + _srcDayCount + "pctchgclose+1)+close*(1+future1pctchgclose))/" + _srcDayCount + ">(ma" + _targetDayCount + "close-close/(acc" + _targetDayCount + "pctchgclose+1)+close*(1+future1pctchgclose))/" + _targetDayCount
			}
			if targetDayCount != "" {
				break
			}
		}
		if srcDayCount != "" {
			break
		}
	}
	conds += " and ma3close is not null and ma3close!=0"
	if startDate != 0 {
		conds = conds + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	var err error
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	err = this.Find(&dayLines, nil, "tscode,tradedate", from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return dayLines, count, nil
}

/**
查找收盘价相关性的股票
*/
func (this *DayLineService) FindCorr(ts_code string, startDate int64, from int, limit int, orderby string, count int64) ([]*entity.PortfolioStat, int64, error) {
	paras := make([]interface{}, 0)
	sql := "select src.tscode as ts_code,target.tscode as target_ts_code,corr(src.pctchgclose,target.pctchgclose) as stat_value" +
		" from (select tscode,tradedate,pctchgclose as pctchgclose from stk_dayline"
	if startDate != 0 {
		sql = sql + " where tradedate>=?"
		paras = append(paras, startDate)
	}
	sql = sql + ") src" +
		" join (select tscode as tscode,tradedate as tradedate,pctchgclose as pctchgclose" +
		" from stk_dayline where tscode not like '8%' and tscode not like '4%'" +
		" and tscode in (select tscode from stk_dayline group by tscode having count(*)>500)"
	if startDate != 0 {
		sql = sql + " and tradedate>=?"
		paras = append(paras, startDate)
	}
	sql = sql + ") target" +
		" on src.tradedate=target.tradedate" +
		" where src.tscode=?"
	paras = append(paras, ts_code)
	sql = sql + " group by src.tscode,target.tscode"
	var err error
	if count == 0 {
		shares, _ := GetShareService().GetCacheShare()
		count = int64(len(shares))
	}
	sql = "select * from (" + sql + ") t where stat_value is not null"
	if orderby == "" || orderby == "desc" {
		sql = sql + " order by stat_value desc"
	} else if orderby == "asc" {
		sql = sql + " order by stat_value"
	}
	if from > 0 {
		sql = sql + " offset " + fmt.Sprint(from)
	}
	if limit > 0 {
		sql = sql + " limit " + fmt.Sprint(limit)
	}
	results, err := this.Query(sql, paras...)
	if err != nil {
		return nil, count, err
	}
	corrs := make([]*entity.PortfolioStat, 0)
	jsonMap, _, _ := stock.GetJsonMap(entity.PortfolioStat{})
	_, shares := GetShareService().GetCacheShare()
	for _, result := range results {
		corr := &entity.PortfolioStat{}
		for colname, v := range result {
			err = reflect.Set(corr, jsonMap[colname], string(v))
			if err != nil {
				logger.Sugar.Errorf("Set colname %v value %v error", colname, string(v))
			}
		}
		share, ok := shares[corr.TsCode]
		if ok {
			corr.SecurityName = share.Name
		}
		share, ok = shares[corr.TargetTsCode]
		if ok {
			corr.TargetSecurityName = share.Name
		}
		corr.StartDate = startDate
		corr.EndDate = stock.CurrentDate()
		corr.Source = "corr"
		corr.SourceName = "pctchgclose"
		corrs = append(corrs, corr)
	}

	return corrs, count, nil
}

/**
刷新所有股票的日线统计数据，包括移动平均，累计增长，均值，相对标准差
*/
func (this *DayLineService) RefreshStat(startDate int64) error {
	processLog := GetProcessLogService().StartLog("", "RefreshStat", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, this.AsyncUpdateStat, nil)
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

func (this *DayLineService) AsyncUpdateStat(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	startDate := (para.([]interface{}))[1].(int64)
	this.UpdateStat(tscode, startDate)
}

func (this *DayLineService) UpdateStat(tscode string, startDate int64) (int64, error) {
	//processLog := GetProcessLogService().StartLog("dayline", "updateStat", tscode)
	affected, err := this.updateStat(tscode, startDate)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return 0, err
	}
	//GetProcessLogService().EndLog(processLog, "", "")
	//processLog = GetProcessLogService().StartLog("dayline", "UpdateBeforeMa", tscode)
	affected, err = this.UpdateBeforeMa(tscode, startDate)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return 0, err
	}
	//GetProcessLogService().EndLog(processLog, "", "")

	return affected, err
}

/**
计算移动平均
avg(close) over (partition by tscode
 order by tradedate asc
 rows between 5 preceding and current row)
其中rows between and子句是最重要的部分，用于定义取值窗口：
5 preceding：过去5行
5 following：未来5行
current row：当前行
unbounded：代替数字代表所有行
*/
func (this *DayLineService) updateStat(tscode string, startDate int64) (int64, error) {
	jsonMap, _, jsonHeads := stock.GetJsonMap(&entity.DayLine{})
	sql := "update stk_dayline d set "
	updateFields := ""
	selectFields := ""
	i := 0
	for _, jsonHead := range jsonHeads[33:] {
		prefix := ""
		suffix := ""
		if strings.HasPrefix(jsonHead, "max") {
			prefix = "max"
			suffix = strings.TrimPrefix(jsonHead, "max")
		} else if strings.HasPrefix(jsonHead, "ma") {
			prefix = "ma"
			suffix = strings.TrimPrefix(jsonHead, "ma")
		} else if strings.HasPrefix(jsonHead, "min") {
			prefix = "min"
			suffix = strings.TrimPrefix(jsonHead, "min")
		} else if strings.HasPrefix(jsonHead, "acc") {
			prefix = "acc"
			suffix = strings.TrimPrefix(jsonHead, "acc")
		} else if strings.HasPrefix(jsonHead, "future") {
			prefix = "future"
			suffix = strings.TrimPrefix(jsonHead, "future")
		}

		pos := strings.Index(suffix, "_")
		if pos < 0 {
			continue
		}
		v, err := convert.ToObject(suffix[0:pos], "int64")
		if err != nil || v == nil {
			continue
		}
		past, _ := v.(int64)
		if past <= 0 {
			continue
		}
		fieldname, ok := jsonMap[jsonHead]
		if !ok {
			continue
		}
		//过去past天收盘价的均值
		ma := "avg(Close) over (partition by tscode order by tradedate rows between " + fmt.Sprint(past) + " preceding and 1 preceding) as " + fieldname
		//过去past天最高收盘价
		max := "max(Close) over (partition by tscode order by tradedate rows between " + fmt.Sprint(past) + " preceding and 1 preceding) as " + fieldname
		//过去past天最低收盘价
		min := "min(Close) over (partition by tscode order by tradedate rows between " + fmt.Sprint(past) + " preceding and 1 preceding) as " + fieldname
		//过去第past天收盘价
		acc := "avg(Close) over (partition by tscode order by tradedate rows between " + fmt.Sprint(past) + " preceding and " + fmt.Sprint(past) + " preceding) as " + fieldname
		//未来第past天收盘价
		future := "avg(Close) over (partition by tscode order by tradedate rows between " + fmt.Sprint(past) + " following and " + fmt.Sprint(past) + " following) as " + fieldname

		if i > 0 {
			updateFields = updateFields + ","
			selectFields = selectFields + ","
		}
		if prefix == "ma" { //过去past天收盘价的均值
			updateFields = updateFields + fieldname + "=stat." + fieldname
			selectFields = selectFields + ma
		} else if prefix == "max" { //过去past天最高收盘价
			updateFields = updateFields + fieldname + "=stat." + fieldname
			selectFields = selectFields + max
		} else if prefix == "min" { //过去past天最低收盘价
			updateFields = updateFields + fieldname + "=stat." + fieldname
			selectFields = selectFields + min
		} else if prefix == "acc" { //今天的收盘价/过去第past天收盘价-1，即累计涨幅
			updateFields = updateFields + fieldname + "=case when stat." + fieldname + "!=0 then (stat.close/stat." + fieldname + "-1) else d." + fieldname + " end"
			selectFields = selectFields + acc
		} else if prefix == "future" { //过去第past天收盘价/今天的收盘价-1，即累计涨幅
			updateFields = updateFields + fieldname + "=case when stat.close!=0 and stat." + fieldname + "!=0 then (stat." + fieldname + "/stat.close-1) else d." + fieldname + " end"
			selectFields = selectFields + future
		}
		i++
	}
	sql = sql + updateFields + " from (select tscode,tradedate,close," + selectFields
	sql += " from stk_dayline where "
	cond, tscodeParas := stock.InBuildStr("tscode", tscode, ",")
	paras := make([]interface{}, 0)
	if cond != "" {
		sql += cond
		paras = append(paras, tscodeParas...)
	}
	sql += ") stat where stat.tscode = d.tscode and stat.tradedate=d.tradedate"
	if startDate == 0 {
		max, err := this.findMaxMaTradeDate(tscode, "ma3close")
		if err == nil && max != nil {
			//由于要计算future的统计值，所以提前1年计算，更好的办法是统计future的命令单独分开
			startDate = max.TradeDate - 10000
		}
	}
	if startDate > 0 {
		sql += " and stat.tradedate > ?"
		paras = append(paras, startDate)
	}
	//sql += " order by stat.tscode,stat.tradedate desc"
	processLog := GetProcessLogService().StartLog("dayline", "updateStat Exec", tscode)
	result, err := this.Exec(sql, paras...)
	if err != nil {
		return 0, err
	}
	GetProcessLogService().EndLog(processLog, "", "")
	if result == nil {
		return 0, errors.New("result is nil")
	}

	return result.RowsAffected()
}

/**
刷新所有股票的日线统计数据，包括移动平均，累计增长，均值，相对标准差
*/
func (this *DayLineService) RefreshBeforeMa(startDate int64) error {
	preocessLog := GetProcessLogService().StartLog("", "RefreshBeforeMa", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, this.AsyncUpdateBeforeMa, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	for _, ts_code := range ts_codes {
		para := make([]interface{}, 0)
		para = append(para, ts_code)
		para = append(para, startDate)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetProcessLogService().EndLog(preocessLog, "", "")
	return nil
}

func (this *DayLineService) AsyncUpdateBeforeMa(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	startDate := (para.([]interface{}))[1].(int64)
	this.UpdateBeforeMa(tscode, startDate)
}

/**
填充过去1，3，5天的各移动平均值
*/
func (this *DayLineService) UpdateBeforeMa(tscode string, startDate int64) (int64, error) {
	jsonMap, _, jsonHeads := stock.GetJsonMap(&entity.DayLine{})
	sql := "update stk_dayline d set "
	updateFields := ""
	selectFields := ""
	i := 0
	for _, jsonHead := range jsonHeads[78:108] {
		prefix := ""
		suffix := ""
		if strings.HasPrefix(jsonHead, "before") {
			prefix = "before"
			suffix = strings.TrimPrefix(jsonHead, "before")
		}

		pos := strings.Index(suffix, "_")
		if pos < 0 {
			continue
		}
		v, err := convert.ToObject(suffix[0:pos], "int64")
		if err != nil || v == nil {
			continue
		}
		past, _ := v.(int64)
		if past <= 0 {
			continue
		}
		macolname := suffix[pos+1:]
		mafieldname, ok := jsonMap[macolname]
		if !ok {
			continue
		}
		fieldname, ok := jsonMap[jsonHead]
		if !ok {
			continue
		}
		//过去第past天ma
		before := "avg(" + mafieldname + ") over (partition by tscode order by tradedate rows between " + fmt.Sprint(past) + " preceding and " + fmt.Sprint(past) + " preceding) as " + fieldname

		if i > 0 {
			updateFields = updateFields + ","
			selectFields = selectFields + ","
		}
		if prefix == "before" { //
			updateFields = updateFields + fieldname + "=stat." + fieldname
			selectFields = selectFields + before
		}
		i++
	}
	sql = sql + updateFields + " from (select tscode,tradedate,close," + selectFields
	sql += " from stk_dayline where "
	cond, tscodeParas := stock.InBuildStr("tscode", tscode, ",")
	paras := make([]interface{}, 0)
	if cond != "" {
		sql += cond
		paras = append(paras, tscodeParas...)
	}
	sql += ") stat where stat.tscode = d.tscode and stat.tradedate=d.tradedate"
	if startDate == 0 {
		max, err := this.findMaxMaTradeDate(tscode, "before3ma21close")
		if err == nil && max != nil {
			startDate = max.TradeDate
		}
	}
	if startDate > 0 {
		sql += " and stat.tradedate > ?"
		paras = append(paras, startDate)
	}
	//sql += " order by stat.tscode,stat.tradedate desc"
	result, err := this.Exec(sql, paras...)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, errors.New("result is nil")
	}

	return result.RowsAffected()
}

func init() {
	service.GetSession().Sync(new(entity.DayLine))
	dayLineService.OrmBaseService.GetSeqName = dayLineService.GetSeqName
	dayLineService.OrmBaseService.FactNewEntity = dayLineService.NewEntity
	dayLineService.OrmBaseService.FactNewEntities = dayLineService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("dayline", dayLineService)
}
