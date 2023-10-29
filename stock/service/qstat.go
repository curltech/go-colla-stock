package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	entity "github.com/curltech/go-colla-stock/stock/entity"
	"math"
	"strings"
	"sync"
)

type QStatService struct {
	service.OrmBaseService
	Terms     []int
	MedianMap map[string]*entity.QStat
	Locker    sync.Mutex
}

var qstatService = &QStatService{Terms: []int{0, 1, 3, 5, 8, 10, 15}, Locker: sync.Mutex{}}

func GetQStatService() *QStatService {
	return qstatService
}

func (this *QStatService) GetSeqName() string {
	return seqname
}

func (this *QStatService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.QStat{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *QStatService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.QStat, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *QStatService) Search(keyword string, terms []int, sourceOptions []string, from int, limit int, count int64) ([]*entity.QStat, int64, error) {
	termConds, termParas := stock.InBuildInt("term", terms)
	paras := make([]interface{}, 0)
	conds := termConds
	paras = append(paras, termParas...)
	conds += " and source!='sum' and source!='stddev'"
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
	qstats := make([]*entity.QStat, 0)
	var err error
	condiBean := &entity.QStat{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	orderby := "tscode,term,startdate desc"
	err = this.Find(&qstats, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}

	return qstats, count, nil
}

/*
*
查询股票季度业绩统计数据
*/
func (this *QStatService) FindQStat(ts_code string, terms []int, source string, sourceName string) (map[string]map[int][]interface{}, error) {
	tscodeConds, tscodeParas := stock.InBuildStr("tscode", ts_code, ",")
	termConds, termParas := stock.InBuildInt("term", terms)
	paras := make([]interface{}, 0)
	sql := tscodeConds + " and " + termConds
	paras = append(paras, tscodeParas...)
	paras = append(paras, termParas...)
	if source != "" {
		sql = sql + " and source=?"
		paras = append(paras, source)
	}
	if sourceName != "" {
		sql = sql + " and sourcename=?"
		paras = append(paras, sourceName)
	}
	qstats := make([]*entity.QStat, 0)
	err := this.Find(&qstats, nil, "tscode,term,source,sourcename", 0, 0, sql, paras...)
	if err != nil {
		return nil, err
	}
	qpMap := make(map[string]map[int][]interface{}, 0)
	for _, qstat := range qstats {
		psMap, ok := qpMap[qstat.TsCode]
		if !ok {
			psMap = make(map[int][]interface{}, 0)
			qpMap[qstat.TsCode] = psMap
		}
		ps, ok := psMap[qstat.Term]
		if !ok {
			ps = make([]interface{}, 0)
		}
		ps = append(ps, qstat)
		psMap[qstat.Term] = ps
	}

	return qpMap, nil
}

/*
*
计算股票季度业绩统计数据中位数，最大值和最小值，并返回结果，方便进行去极值和标准化
*/
func (this *QStatService) FindQStatMedian() (map[string]*entity.QStat, error) {
	jsonMap, _, jsonHeads := stock.GetJsonMap(entity.QStat{})
	sql := "select source as source,sourcename as source_name,term as term"
	for _, jsonHead := range jsonHeads[17:] {
		fieldname := jsonMap[jsonHead]
		sql = sql + ",percentile_cont(0.5) within group(order by " + fieldname + ") as median_" + jsonHead
		sql = sql + ",percentile_cont(0.975) within group(order by " + fieldname + ") as max_" + jsonHead
		sql = sql + ",percentile_cont(0.025) within group(order by " + fieldname + ") as min_" + jsonHead
	}
	sql = sql + " from stk_qstat"
	sql = sql + " group by source,sourcename,term"
	sql = sql + " order by source,sourcename,term"
	paras := make([]interface{}, 0)
	results, err := this.Query(sql, paras...)
	if err != nil {
		return nil, err
	}
	medianMap := make(map[string]*entity.QStat, 0)
	for _, result := range results {
		median := &entity.QStat{}
		max := &entity.QStat{}
		min := &entity.QStat{}
		for colname, v := range result {
			s := string(v)
			if s != "" {
				if strings.HasPrefix(colname, "median_") {
					colname = strings.TrimPrefix(colname, "median_")
					err = reflect.Set(median, jsonMap[colname], s)
				} else if strings.HasPrefix(colname, "max_") {
					colname = strings.TrimPrefix(colname, "max_")
					err = reflect.Set(max, jsonMap[colname], s)
				} else if strings.HasPrefix(colname, "min_") {
					colname = strings.TrimPrefix(colname, "min_")
					err = reflect.Set(min, jsonMap[colname], s)
				} else {
					err = reflect.Set(median, jsonMap[colname], s)
					err = reflect.Set(max, jsonMap[colname], s)
					err = reflect.Set(min, jsonMap[colname], s)
				}
				if err != nil {
					logger.Sugar.Errorf("Set colname %v value %v error:%v", colname, s, err.Error())
				}
			}
		}
		key := median.Source + ":" + median.SourceName + ":" + fmt.Sprint(median.Term)
		medianMap["median"+":"+key] = median
		medianMap["max"+":"+key] = max
		medianMap["min"+":"+key] = min
	}

	return medianMap, nil
}

/*
*
删除股票季度业绩统计数据
*/
func (this *QStatService) deleteQStat(ts_code string) error {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	qstat := &entity.QStat{}
	_, err := this.Delete(qstat, conds, paras...)
	if err != nil {
		return err
	}

	return nil
}

/*
*
对某个统计指标计算minmax标准化
*/
func (this *QStatService) MinmaxStd(q *entity.QStat, fieldname string) float64 {
	key := q.Source + ":" + q.SourceName + ":" + fmt.Sprint(q.Term)

	max, ok := this.MedianMap["max"+":"+key]
	if !ok {
		logger.Sugar.Errorf("MedianMap no max key:%v", key)
		return 0
	}
	v, err := reflect.GetValue(max, fieldname)
	if err != nil {
		logger.Sugar.Errorf("Max GetValue no value fieldname:%v", fieldname)
		return 0
	}
	maxVal, ok := v.(float64)
	if !ok {
		return 0
	}
	min, ok := this.MedianMap["min"+":"+key]
	if !ok {
		logger.Sugar.Errorf("MedianMap no min key:%v", key)
		return 0
	}
	v, err = reflect.GetValue(min, fieldname)
	if err != nil {
		logger.Sugar.Errorf("Min GetValue no value fieldname:%v", fieldname)
		return 0
	}
	minVal, ok := v.(float64)
	if !ok {
		return 0
	}
	v, err = reflect.GetValue(q, fieldname)
	if err != nil {
		logger.Sugar.Errorf("QStat GetValue no value fieldname:%v", fieldname)
		return 0
	}
	val, ok := v.(float64)
	if !ok {
		return 0
	}
	if val > maxVal {
		val = maxVal
	}
	if val < minVal {
		val = minVal
	}
	if stock.Equal(maxVal, minVal) {
		return 1
	}
	return (val - minVal) / (maxVal - minVal)
}

/*
*
计算股票所有统计数据的Min，Max，并保存
*/
func (this *QStatService) GetQStatMedian() error {
	this.Locker.Lock()
	defer this.Locker.Unlock()
	var err error
	if this.MedianMap == nil {
		this.MedianMap, err = this.FindQStatMedian()
		if err != nil {
			logger.Sugar.Errorf("Error:%v", err.Error())
			return err
		}
		logger.Sugar.Infof("FindQStatMedian successfully!")
	}
	return nil
}

/*
*
刷新所有股票的季度业绩统计数据
*/
func (this *QStatService) RefreshQStat() error {
	processLog := GetProcessLogService().StartLog("qstat", "RefreshQStat", "")
	routinePool := thread.CreateRoutinePool(10, this.AsyncUpdateQStat, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	for _, ts_code := range ts_codes {
		para := make([]interface{}, 0)
		para = append(para, ts_code)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	GetQStatService().Locker.Lock()
	defer GetQStatService().Locker.Unlock()
	GetQStatService().MedianMap = nil

	GetProcessLogService().EndLog(processLog, "", "")

	return nil
}

func (this *QStatService) AsyncUpdateQStat(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	this.GetUpdateQStat(tscode)
}

/*
*
更新股票季度业绩统计数据，并返回结果
*/
func (this *QStatService) GetUpdateQStat(tscode string) ([]interface{}, error) {
	var ps []interface{}
	var err error
	//processLog := GetProcessLogService().StartLog("qstat", "GetUpdateQStat", tscode)
	this.deleteQStat(tscode)
	//for _, term := range this.Terms {
	//	ps, err = this.UpdateQStatBySql(tscode, term)
	//}
	ps, err = this.UpdateQStat(tscode, this.Terms)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	GetStatScoreService().GetUpdateStatScore(tscode)

	return ps, err
}

/*
*
通过内存更新股票季度业绩统计数据，并返回结果
*/
func (this *QStatService) UpdateQStat(tscode string, terms []int) ([]interface{}, error) {
	//获取降序的记录
	qpMap, err := GetQPerformanceService().FindStdQPerformance(tscode, nil, "", "", StdtypeNone, false)
	if err != nil {
		return nil, err
	}
	//获取所有的term的参数，数据降序排列
	qtermMap, err := GetQPerformanceService().GetQTerm(qpMap, terms)
	if err != nil {
		return nil, err
	}
	if qtermMap == nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, "qterm is nil")
		return nil, errors.New("qterm is nil")
	}
	//获取每个term的统计数据，数据降序排列
	statMap := this.FindAllQStat(qpMap, qtermMap)
	if len(statMap) <= 0 {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, "qstats len is 0")
		return nil, errors.New("qstats len is 0")
	}
	ps := make([]interface{}, 0)
	for _, qstats := range statMap {
		for _, qstat := range qstats {
			ps = append(ps, qstat)
		}
	}
	_, err = this.Insert(ps...)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error: %s", tscode, err.Error())
		return nil, err
	}
	//更新排名统计数据
	//go this.UpdatePercentRank(tscode, qtermMap)

	return ps, err
}

func (this *QStatService) UpdatePercentRank(tscode string, qtermMap map[string]map[int]*QTerm) ([]interface{}, error) {
	ps := make([]interface{}, 0)
	//更新排名统计数据
	for ts_code, qterms := range qtermMap {
		for _, qterm := range qterms {
			qstats, err := this.FindPercentRank(ts_code, qterm)
			if err == nil {
				for _, qstat := range qstats {
					ps = append(ps, qstat)
				}
			}
		}
	}

	_, err := this.Insert(ps...)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error: %s", tscode, err.Error())
		return nil, err
	}

	return ps, err
}

func (this *QStatService) FindPercentRank(tscode string, qterm *QTerm) ([]*entity.QStat, error) {
	qss := make([]*entity.QStat, 0)
	qstats, err := this.findPercentRank("tscode", tscode, qterm)
	if err == nil {
		qss = append(qss, qstats...)
	}
	qstats, err = this.findPercentRank("industry", tscode, qterm)
	if err == nil {
		qss = append(qss, qstats...)
	}
	qstats, err = this.findPercentRank("sector", tscode, qterm)
	if err == nil {
		qss = append(qss, qstats...)
	}

	return qss, nil
}

func (this *QStatService) findPercentRank(rankType string, tscode string, qterm *QTerm) ([]*entity.QStat, error) {
	qstats := make([]*entity.QStat, 0)
	startDate := qterm.StartDate
	endDate := qterm.EndDate
	qps, err := GetQPerformanceService().FindPercentRank(rankType, tscode, 0, startDate, endDate, 0, 1, 0)
	if err != nil {
		return nil, err
	}
	reportNumber := len(qps)
	for _, qp := range qps {
		qstat := this.toQStat(qp, qterm)
		qstat.Source = "rank"
		qstat.SourceName = rankType
		qstat.Id = 0
		qstat.TsCode = qp.TsCode
		qstat.SecurityName = qp.SecurityName
		qstat.Industry = qp.Industry
		qstat.Sector = qp.Sector
		qstat.ReportNumber = reportNumber
		qstats = append(qstats, qstat)
	}

	return qstats, nil
}

/*
*
通过内存计算股票季度业绩全部统计数据，并返回结果，原始数据降序排列
*/
func (this *QStatService) FindAllQStat(qpMap map[string][]*entity.QPerformance, qtermMap map[string]map[int]*QTerm) map[string][]interface{} {
	jsonMap, _, jsonHeads := stock.GetJsonMap(entity.QPerformance{})
	qstatMap := make(map[string][]interface{}, 0)
	for tscode, qps := range qpMap {
		if len(qps) == 0 {
			continue
		}
		qstats, ok := qstatMap[tscode]
		if !ok {
			qstats = make([]interface{}, 0)
		}
		qterms, ok := qtermMap[tscode]
		if !ok {
			continue
		}
		//最新的数据，所有的累计值都是与last值的比较
		last := qps[0]
		psMap := make(map[int][]interface{}, 0)
		//每个term的前一个记录，用于计算累计增长
		preMap := make(map[int]interface{}, 0)
		for _, qp := range qps {
			qdate := qp.QDate
			for term, qterm := range qterms {
				//如果当前计算的term的开始日期小于记录日期，为当前term的有效记录
				if qterm.StartDate <= qdate {
					ps, ok := psMap[term]
					if !ok {
						ps = make([]interface{}, 0)
					}
					ps = append(ps, qp)
					psMap[term] = ps
				} else {
					_, ok := preMap[term]
					if !ok {
						preMap[term] = qp
					}
				}
			}
		}
		for term, qterm := range qterms {
			ps, ok := psMap[term]
			if ok {
				reportNumber := len(ps)
				stat := stock.CreateStat(ps, jsonHeads[14:])
				sum := stat.CalSum()
				qstat := this.toQStat(sum, qterm)
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.Source = "sum"
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				max := stat.Max
				qstat = this.toQStat(max, qterm)
				qstat.Source = "max"
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				min := stat.Min
				qstat = this.toQStat(min, qterm)
				qstat.Source = "min"
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				mean := stat.CalMean()
				qstat = this.toQStat(mean, qterm)
				qstat.Source = "mean"
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				medians := stat.CalMedian()
				qstat = this.toQStat(medians[1], qterm)
				qstat.Source = "median"
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				stddev := stat.CalStddev()
				qstat = this.toQStat(stddev, qterm)
				qstat.Source = "stddev"
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				rsd := stat.CalRsd()
				qstat = this.toQStat(rsd, qterm)
				qstat.Source = "rsd"
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				cor := stat.CalCor("market_value")
				qstat = this.toQStat(cor, qterm)
				qstat.Id = 0
				qstat.TsCode = last.TsCode
				qstat.SecurityName = last.SecurityName
				qstat.Industry = last.Industry
				qstat.Sector = last.Sector
				qstat.Source = "corr"
				fieldname, _ := jsonMap["market_value"]
				qstat.SourceName = fieldname
				qstat.ReportNumber = reportNumber
				qstats = append(qstats, qstat)
				var pre *entity.QPerformance
				v, ok := preMap[term]
				if ok {
					pre = v.(*entity.QPerformance)
				}
				acc := this.FindAcc(ps, qterm, pre)
				if acc != nil {
					qstats = append(qstats, acc)
				}
				lastQStat := &entity.QStat{}
				bs, err := json.Marshal(last)
				if err == nil {
					err = json.Unmarshal(bs, lastQStat)
					if err != nil {
						logger.Sugar.Errorf("lastQStat set value fail")
					}
					lastQStat.Id = 0
					lastQStat.Source = "last"
					lastQStat.SourceName = last.Source
					lastQStat.StartDate = last.QDate
					lastQStat.EndDate = last.NDate
					lastQStat.ActualStartDate = qterm.ActualStartDate
					lastQStat.Term = qterm.Term
					lastQStat.ReportNumber = 1
				}
				qstats = append(qstats, lastQStat)
				qstatMap[tscode] = qstats
			}
		}
	}

	return qstatMap
}

func (this *QStatService) toQStat(val interface{}, qterm *QTerm) *entity.QStat {
	qstat := &entity.QStat{}
	bs, err := json.Marshal(val)
	if err == nil {
		err = json.Unmarshal(bs, qstat)
	}

	qstat.ActualStartDate = qterm.ActualStartDate
	qstat.StartDate = qterm.StartDate
	qstat.EndDate = qterm.EndDate
	qstat.Term = qterm.Term
	qstat.TradeDate = qterm.TradeDate

	return qstat
}

// 原始数据降序排列
func (this *QStatService) FindAcc(ps []interface{}, qterm *QTerm, pre *entity.QPerformance) *entity.QStat {
	jsonMap, _, jsonHeads := stock.GetJsonMap(entity.QPerformance{})
	if pre == nil {
		pre = ps[len(ps)-1].(*entity.QPerformance)
	}
	end := ps[0].(*entity.QPerformance)
	qs := &entity.QStat{}
	qs.TsCode = end.TsCode
	qs.SecurityName = end.SecurityName
	qs.Industry = end.Industry
	qs.Sector = end.Sector
	qs.Source = "acc"
	qs.TradeDate = qterm.TradeDate
	qs.StartDate = qterm.StartDate
	qs.EndDate = qterm.EndDate
	qs.ActualStartDate = qterm.ActualStartDate
	qs.Term = qterm.Term
	qs.ReportNumber = len(ps)
	t := qterm.Term
	var err error
	if t == 0 {
		startDate := qs.StartDate
		endDate := qs.EndDate
		t, err = stock.DiffYear(startDate, endDate)
		if err != nil || t == 0 {
			t = 1
		}
	}
	for _, jsonHead := range jsonHeads[14:] {
		fieldname := jsonMap[jsonHead]
		v, _ := reflect.GetValue(pre, fieldname)
		startVal, _ := v.(float64)
		if !stock.Equal(startVal, 0.0) {
			v, _ = reflect.GetValue(end, fieldname)
			endVal, _ := v.(float64)
			diff := endVal - startVal
			apr := (diff / math.Abs(startVal)) / float64(t) // stock.CalApr((endVal-startVal)/math.Abs(startVal)+1, float64(t))
			reflect.SetValue(qs, fieldname, apr*100)
		}
	}
	return qs
}

/*
*
通过数据库sql更新股票季度业绩统计数据，并返回结果
*/
func (this *QStatService) UpdateQStatBySql(tscode string, term int) ([]interface{}, error) {
	qterm, err := GetQPerformanceService().GetQTermBySql(tscode, term)
	if err != nil {
		return nil, err
	}
	if qterm == nil {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, "qterm is nil")
		return nil, errors.New("qterm is nil")
	}
	if qterm.ActualStartDate > qterm.StartDate {
		//logger.Sugar.Errorf("tscode:%v Error:%v", tscode, "ActualStartDate>StartDate")
		return nil, errors.New("ActualStartDate>StartDate")
	}
	statMap := GetQPerformanceService().FindAllQStatBySql(tscode, qterm.StartDate, "")
	if len(statMap) <= 0 {
		logger.Sugar.Errorf("tscode:%v Error:%v", tscode, "qstats len is 0")
		return nil, errors.New("qstats len is 0")
	}
	ps := make([]interface{}, 0)
	for _, qstats := range statMap {
		for _, stat := range qstats {
			qstat := stat.(*entity.QStat)
			qstat.ActualStartDate = qterm.ActualStartDate
			qstat.StartDate = qterm.StartDate
			qstat.EndDate = qterm.EndDate
			qstat.Term = term
			qstat.TradeDate = qterm.TradeDate
			//计算累计的年化增长
			if qstat.Source == "acc" {
				t := qstat.Term
				if t == 0 {
					startDate := qstat.StartDate
					end := qstat.EndDate
					t, err = stock.DiffYear(startDate, end)
					if err != nil {
						t = 1
					}
				}
				qstat.PctChgMarketValue = stock.CalApr(qstat.PctChgMarketValue, float64(t))
				qstat.YoySales = stock.CalApr(qstat.YoySales, float64(t))
				yoyDeduNp := stock.CalApr(qstat.YoyDeduNp, float64(t))
				qstat.YoyDeduNp = yoyDeduNp
			}
			ps = append(ps, qstat)
		}
	}
	last, err := GetQPerformanceService().findMaxQDate(tscode, 0)
	if err == nil {
		lastQStat := &entity.QStat{}
		bs, err := json.Marshal(last)
		if err == nil {
			err = json.Unmarshal(bs, lastQStat)
			if err != nil {
				logger.Sugar.Errorf("latestQStat set value fail")
			}
			lastQStat.Source = "last"
			lastQStat.SourceName = last.Source
			lastQStat.StartDate = last.QDate
			lastQStat.EndDate = last.NDate
			ps = append(ps, lastQStat)
		}
	}
	_, err = this.Upsert(ps...)
	if err != nil {
		logger.Sugar.Errorf("tscode:%v Error: %s", tscode, err.Error())
		return nil, err
	}

	return ps, err
}

func init() {
	service.GetSession().Sync(new(entity.QStat))
	qstatService.OrmBaseService.GetSeqName = qstatService.GetSeqName
	qstatService.OrmBaseService.FactNewEntity = qstatService.NewEntity
	qstatService.OrmBaseService.FactNewEntities = qstatService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("qstat", qstatService)
}
