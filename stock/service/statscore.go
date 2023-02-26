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
	entity "github.com/curltech/go-colla-stock/stock/entity"
	"math"
	"strings"
	"time"
)

type StatScoreService struct {
	service.OrmBaseService
}

var statScoreService = &StatScoreService{}

func GetStatScoreService() *StatScoreService {
	return statScoreService
}

func (this *StatScoreService) GetSeqName() string {
	return seqname
}

func (this *StatScoreService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.StatScore{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *StatScoreService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.StatScore, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *StatScoreService) Search(keyword string, tscode string, terms []int, scoreOptions []string, orderby string, from int, limit int, count int64) ([]*entity.StatScore, int64, error) {
	termConds, termParas := stock.InBuildInt("term", terms)
	paras := make([]interface{}, 0)
	conds := termConds
	paras = append(paras, termParas...)
	if keyword != "" {
		conds += " and tscode in (select tscode from stk_share where name like ? or tscode like ? or pinyin like ?)"
		paras = append(paras, keyword+"%")
		paras = append(paras, keyword+"%")
		paras = append(paras, strings.ToLower(keyword)+"%")
	}
	if tscode != "" {
		tscodeConds, tscodeParas := stock.InBuildStr("tscode", tscode, ",")
		conds = conds + " and " + tscodeConds
		paras = append(paras, tscodeParas...)
	}
	statScores := make([]*entity.StatScore, 0)
	var err error
	condiBean := &entity.StatScore{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "totalscore"
	}
	i := 0
	if scoreOptions == nil || len(scoreOptions) == 0 {
		orderby = "term"
	} else {
		for _, scoreOption := range scoreOptions {
			if scoreOption == "TotalScore" {
				orderby = "totalscore"
				break
			} else {
				if i == 0 {
					orderby = scoreOption
				} else {
					orderby = orderby + "+" + scoreOption
				}
			}
			i++
		}
		orderby = orderby + " desc"
	}
	err = this.Find(&statScores, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	i = 1
	for _, statScore := range statScores {
		statScore.Id = uint64(from + i)
		i++
	}

	return statScores, count, nil
}

/**
计算股票季度业绩统计数据中位数，最大值和最小值，并返回结果，方便进行去极值和标准化
*/
func (this *StatScoreService) FindStatScoreMedian() (map[string]*entity.StatScore, error) {
	jsonHeads := []string{"RiskScore", "AccScore", "ProsScore", "TrendScore", "IncreaseScore", "CorrScore", "PriceScore", "CorrScore", "TotalScore"}
	jsonMap, _, _ := stock.GetJsonMap(entity.StatScore{})
	sql := "select term as term"
	for _, jsonHead := range jsonHeads {
		fieldname := jsonMap[jsonHead]
		sql = sql + ",percentile_cont(0.5) within group(order by " + fieldname + ") as median_" + jsonHead
		sql = sql + ",max(" + fieldname + ") as max_" + jsonHead
		sql = sql + ",min(" + fieldname + ") as min_" + jsonHead
	}
	sql = sql + " from stk_statscore"
	sql = sql + " group by term"
	sql = sql + " order by term"
	results, err := this.Query(sql, nil)
	if err != nil {
		return nil, err
	}
	medianMap := make(map[string]*entity.StatScore, 0)
	for _, result := range results {
		median := &entity.StatScore{}
		max := &entity.StatScore{}
		min := &entity.StatScore{}
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
		key := fmt.Sprint(median.Term)
		medianMap["median"+":"+key] = median
		medianMap["max"+":"+key] = max
		medianMap["min"+":"+key] = min
	}

	return medianMap, nil
}

func (this *StatScoreService) deleteStatScore(ts_code string) error {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	statScore := &entity.StatScore{}
	_, err := this.Delete(statScore, conds, paras...)
	if err != nil {
		return err
	}

	return nil
}

/**
刷新所有股票的季度业绩统计评分数据
*/
func (this *StatScoreService) RefreshStatScore() error {
	processLog := GetProcessLogService().StartLog("qstatscore", "RefreshStatScore", "")
	routinePool := thread.CreateRoutinePool(NetRoutinePoolSize, this.AsyncUpdateStatScore, nil)
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

func (this *StatScoreService) AsyncUpdateStatScore(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	this.GetUpdateStatScore(tscode)
}

func (this *StatScoreService) GetUpdateStatScore(ts_code string) (map[string]map[int]*entity.StatScore, error) {
	//processLog := GetProcessLogService().StartLog("statscore", "GetUpdateStatScore", ts_code)
	statScoresMap, err := this.updateStatScore(ts_code)
	if err != nil {
		//GetProcessLogService().EndLog(processLog, "", err.Error())
		return statScoresMap, err
	}
	return statScoresMap, err
}

/**
刷新股票季度业绩评分数据，在计算评分之前需要对所有股票的业绩统计数据进行minmax标准化
*/
func (this *StatScoreService) updateStatScore(ts_code string) (map[string]map[int]*entity.StatScore, error) {
	qstatService := GetQStatService()
	//err := qstatService.GetQStatMedian()
	//if err != nil {
	//	logger.Sugar.Errorf("ts_code:%v Error:%v", ts_code, err.Error())
	//	return nil, err
	//}
	terms := qstatService.Terms
	qstatMap, err := qstatService.FindQStat(ts_code, terms, "", "")
	if err != nil {
		logger.Sugar.Errorf("ts_code:%v Error:%v", ts_code, err.Error())
	}

	_, shares := GetShareService().GetCacheShare()
	statScoresMap := make(map[string]map[int]*entity.StatScore, 0)
	for tscode, qstatsMap := range qstatMap {
		//最新一年期的累计涨幅，代表股票的景气度
		var lastAccVal *AccValue
		statScoreMap, ok := statScoresMap[tscode]
		if !ok {
			statScoreMap = make(map[int]*entity.StatScore, 0)
			statScoresMap[tscode] = statScoreMap
		}
		terms := []int{1, 0, 3, 5, 8, 10, 15}
		for _, term := range terms {
			qstats, ok := qstatsMap[term]
			if !ok {
				//logger.Sugar.Errorf("ts_code:%v Error:%v", ts_code, errors.New("No term stat data"))
				continue
			}
			if len(qstats) == 0 {
				logger.Sugar.Errorf("ts_code:%v Error:%v", ts_code, errors.New("%v term stat data is len 0"))
				continue
			}
			statScore, ok := statScoreMap[term]
			if !ok {
				statScore = &entity.StatScore{}
				statScoreMap[term] = statScore
			}
			statScore.Term = term
			stat := qstats[0].(*entity.QStat)
			statScore.ReportNumber = stat.ReportNumber
			startDate := stat.StartDate
			statScore.StartDate = startDate
			statScore.EndDate = stat.EndDate
			statScore.TradeDate = stat.TradeDate
			statScore.TsCode = tscode
			statScore.SecurityName = stat.SecurityName
			share, ok := shares[tscode]
			if ok {
				statScore.Industry = share.Industry
				statScore.Sector = share.Sector
				statScore.Area = share.Area
				statScore.Market = share.Market
				listDate := share.ListDate
				if listDate == "" {
					listDate = stat.ActualStartDate
					v, err := convert.ToObject(listDate[0:4], "int64")
					if err == nil && v != nil {
						statScore.ListDate = v.(int64) * 10000
					}
				} else {
					v, err := convert.ToObject(listDate, "int64")
					if err == nil && v != nil {
						statScore.ListDate = v.(int64)
					}
				}
				//计算上市时间，如果小于统计时间，退出
				diff := statScore.TradeDate/10000 - statScore.ListDate/10000
				if diff < int64(term) {
					delete(statScoreMap, term)
					break
				}
				statScore.ListStatus = share.ListStatus
			}
			for _, qstat := range qstats {
				q := qstat.(*entity.QStat)
				source := q.Source
				sourceName := q.SourceName
				if source == "last" {
					if lastAccVal == nil {
						lastAccVal = &AccValue{}
					}
					lastAccVal.PctChgMarketValue = q.PctChgMarketValue
					lastAccVal.YoySales = q.YoySales / 100
					lastAccVal.YoyDeduNp = q.YoyDeduNp / 100
					lastAccVal.OrLastMonth = q.OrLastMonth / 100
					lastAccVal.NpLastMonth = q.NpLastMonth / 100
					lastAccVal.Pe = q.Pe
					lastAccVal.Peg = q.Peg
				} else if source == "acc" {
					statScore.AccPctChgMarketValue = q.MarketValue / 100
					statScore.AccYoySales = q.YearOperateIncome / 100
					statScore.AccYoyDeduNp = q.YearNetProfit / 100
				} else if source == "rsd" {
					statScore.RsdPctChgMarketValue = q.PctChgMarketValue
					statScore.RsdYoySales = q.YoySales
					statScore.RsdYoyDeduNp = q.YoyDeduNp
					statScore.RsdOrLastMonth = q.OrLastMonth
					statScore.RsdNpLastMonth = q.NpLastMonth
					statScore.RsdPe = q.Pe
					statScore.RsdWeightAvgRoe = q.WeightAvgRoe
					statScore.RsdGrossprofitMargin = q.GrossprofitMargin
				} else if source == "mean" {
					statScore.MeanPctChgMarketValue = q.PctChgMarketValue
					statScore.MeanYoySales = q.YoySales / 100
					statScore.MeanYoyDeduNp = q.YoyDeduNp / 100
					statScore.MeanOrLastMonth = q.OrLastMonth / 100
					statScore.MeanNpLastMonth = q.NpLastMonth / 100
					statScore.MeanWeightAvgRoe = q.WeightAvgRoe / 100
					statScore.MeanGrossprofitMargin = q.GrossprofitMargin / 100
					statScore.MeanPe = q.Pe
					statScore.MeanPeg = q.Peg
				} else if source == "median" {
					statScore.MedianPctChgMarketValue = q.PctChgMarketValue
					statScore.MedianYoySales = q.YoySales / 100
					statScore.MedianYoyDeduNp = q.YoyDeduNp / 100
					statScore.MedianOrLastMonth = q.OrLastMonth / 100
					statScore.MedianNpLastMonth = q.NpLastMonth / 100
					statScore.MedianWeightAvgRoe = q.WeightAvgRoe / 100
					statScore.MedianGrossprofitMargin = q.GrossprofitMargin / 100
					statScore.MedianPe = q.Pe
					statScore.MedianPeg = q.Peg
				} else if source == "corr" && sourceName == "MarketValue" {
					statScore.CorrYearOperateIncome = q.YearOperateIncome
					statScore.CorrYearNetProfit = q.YearNetProfit
					statScore.CorrYoySales = q.YoySales
					statScore.CorrYoyDeduNp = q.YoyDeduNp
					statScore.CorrWeightAvgRoe = q.WeightAvgRoe
					statScore.CorrGrossprofitMargin = q.GrossprofitMargin
				}
			}
			if lastAccVal != nil {
				statScore.LastPctChgMarketValue = lastAccVal.PctChgMarketValue
				statScore.LastYoyDeduNp = lastAccVal.YoyDeduNp
				statScore.LastYoySales = lastAccVal.YoySales
				statScore.LastOrLastMonth = lastAccVal.OrLastMonth
				statScore.LastNpLastMonth = lastAccVal.NpLastMonth
				statScore.LastMeanPe = lastAccVal.Pe
				statScore.LastMeanPeg = lastAccVal.Peg
			}

			this.calAccScore(statScore)
			this.calProsScore(statScore)
			this.calPriceScore(statScore)
			this.calStableScore(statScore)
			this.calRiskScore(statScore)
			this.calIncreaseScore(statScore)
			this.calCorrScore(statScore)
			this.calOperationScore(statScore)
			this.calTrendScore(statScore)
		}
	}

	this.deleteStatScore(ts_code)
	ps := make([]interface{}, 0)
	for _, statScoreMap := range statScoresMap {
		for _, statScore := range statScoreMap {
			ps = append(ps, statScore)
		}
	}
	this.Insert(ps...)
	return statScoresMap, nil
}

type AccValue struct {
	PctChgMarketValue float64
	YoyDeduNp         float64
	YoySales          float64
	OrLastMonth       float64
	NpLastMonth       float64
	Pe                float64
	Peg               float64
}

func (this *StatScoreService) minmaxScore(val float64) float64 {
	return 10 * val
}

func (this *StatScoreService) maxminScore(val float64) float64 {
	return 10 * (1 - val)
}

/**
计算各项评分，进行汇总，每个评分项的总分为0-100，每个子项的总分为0-10,5为一般情况
*/
func (this *StatScoreService) calStableScore(statScore *entity.StatScore) {
	score := this.rsdScore(statScore.RsdPctChgMarketValue)
	statScore.StableScore += score

	score = this.rsdScore(statScore.RsdYoySales)
	statScore.StableScore += score

	score = this.rsdScore(statScore.RsdYoyDeduNp)
	statScore.StableScore += score

	score = this.rsdScore(statScore.RsdOrLastMonth)
	statScore.StableScore += score

	score = this.rsdScore(statScore.RsdNpLastMonth)
	statScore.StableScore += score

	score = this.rsdScore(statScore.RsdPe)
	statScore.StableScore += score

	score = this.rsdScore(statScore.RsdWeightAvgRoe)
	statScore.StableScore += score * 0.5

	score = this.rsdScore(statScore.RsdGrossprofitMargin)
	statScore.StableScore += score * 0.5

	statScore.StableScore = statScore.StableScore * 10 / 7

	if statScore.StableScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*业绩极其稳定\n"
	} else if statScore.StableScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*业绩稳定\n"
	} else if statScore.StableScore < 20 {
		statScore.BadTip = statScore.BadTip + "*业绩不稳定\n"
	}

	statScore.TotalScore += statScore.StableScore
}

func (this *StatScoreService) rsdScore(val float64) float64 {
	if stock.Equal(val, 0.0) {
		return 5
	}
	abs := math.Abs(val)
	scoreRules := []float64{0.2, 0.5, 1, 2, 3}
	for i, limit := range scoreRules {
		if abs <= limit {
			score := float64(5 - i)
			return score * 2
		}
	}

	return 0.0
}

func (this *StatScoreService) calIncreaseScore(statScore *entity.StatScore) {
	score := this.increaseScore(statScore.MeanPctChgMarketValue)
	statScore.IncreaseScore += score

	score = this.increaseScore(statScore.MeanYoySales)
	statScore.IncreaseScore += score

	score = this.increaseScore(statScore.MeanYoyDeduNp)
	statScore.IncreaseScore += score

	score = this.increaseScore(statScore.MedianPctChgMarketValue)
	statScore.IncreaseScore += score

	score = this.increaseScore(statScore.MedianYoySales)
	statScore.IncreaseScore += score

	score = this.increaseScore(statScore.MedianYoyDeduNp)
	statScore.IncreaseScore += score

	statScore.IncreaseScore = statScore.IncreaseScore * 10 / 6

	if statScore.IncreaseScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*平均业绩很好\n"
	} else if statScore.IncreaseScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*平均业绩好\n"
	} else if statScore.IncreaseScore < 20 {
		statScore.BadTip = statScore.BadTip + "*平均业绩差\n"
	}

	statScore.TotalScore += statScore.IncreaseScore
}

func (this *StatScoreService) increaseScore(val float64) float64 {
	scoreRules := []float64{0.05, 0.2, 0.3, 0.5, 0.7, 1}
	for i, limit := range scoreRules {
		if val <= limit {
			score := float64(i)
			return score * 10 / 6
		}
	}

	return 10
}

func (this *StatScoreService) calOperationScore(statScore *entity.StatScore) {
	score := this.operationScore(statScore.MeanWeightAvgRoe)
	statScore.OperationScore += score

	score = this.operationScore(statScore.MeanGrossprofitMargin)
	statScore.OperationScore += score

	score = this.operationScore(statScore.MedianWeightAvgRoe)
	statScore.OperationScore += score

	score = this.operationScore(statScore.MedianGrossprofitMargin)
	statScore.OperationScore += score

	statScore.OperationScore = statScore.OperationScore * 2.5

	if statScore.OperationScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*运营质量很好\n"
	} else if statScore.OperationScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*运营质量好\n"
	} else if statScore.OperationScore < 20 {
		statScore.BadTip = statScore.BadTip + "*运营质量差\n"
	}

	statScore.TotalScore += statScore.OperationScore
}

func (this *StatScoreService) operationScore(val float64) float64 {
	scoreRules := []float64{0.05, 0.15, 0.35, 0.7, 1}
	for i, limit := range scoreRules {
		if val <= limit {
			score := float64(i)
			return score * 2
		}
	}

	return 10
}

func (this *StatScoreService) calCorrScore(statScore *entity.StatScore) {
	score := this.corrScore(statScore.CorrYearOperateIncome)
	statScore.CorrScore += score

	score = this.corrScore(statScore.CorrYearNetProfit)
	statScore.CorrScore += score

	score = this.corrScore(statScore.CorrYoySales)
	statScore.CorrScore += score

	score = this.corrScore(statScore.CorrYoyDeduNp)
	statScore.CorrScore += score

	score = this.corrScore(statScore.CorrWeightAvgRoe)
	statScore.CorrScore += score

	score = this.corrScore(statScore.CorrGrossprofitMargin)
	statScore.CorrScore += score

	statScore.CorrScore = statScore.CorrScore * 10 / 6

	if statScore.CorrScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*市值与业绩相关性很强\n"
	} else if statScore.CorrScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*市值与业绩相关性强\n"
	} else if statScore.CorrScore < 20 {
		statScore.BadTip = statScore.BadTip + "*市值与业绩相关性差\n"
	}

	statScore.TotalScore += statScore.CorrScore
}

func (this *StatScoreService) corrScore(val float64) float64 {
	scoreRules := []float64{0.4, 0.6, 0.7, 0.8, 0.9}
	for i, limit := range scoreRules {
		if val <= limit {
			score := float64(i)
			return score * 2
		}
	}

	return 10
}

func (this *StatScoreService) calPriceScore(statScore *entity.StatScore) {
	score := this.peScore(statScore.MeanPe)
	statScore.PriceScore += score

	score = this.peScore(statScore.MedianPe)
	statScore.PriceScore += score

	score = this.pegScore(statScore.MeanPeg)
	statScore.PriceScore += score

	score = this.pegScore(statScore.MedianPeg)
	statScore.PriceScore += score

	statScore.PriceScore = statScore.PriceScore * 2.5
	if statScore.PriceScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*价格很便宜\n"
	} else if statScore.PriceScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*价格便宜\n"
	} else if statScore.PriceScore < 20 {
		statScore.BadTip = statScore.BadTip + "*价格贵\n"
	}
	statScore.TotalScore += statScore.PriceScore
}

func (this *StatScoreService) peScore(val float64) float64 {
	if val < 0 {
		return 0
	}
	scoreRules := []float64{10, 20, 50, 70, 100}
	if val < 0 {
		return 0.0
	}
	for i, limit := range scoreRules {
		if val <= limit {
			score := float64(5 - i)
			return score * 2
		}
	}

	return 10.0
}

func (this *StatScoreService) pegScore(val float64) float64 {
	if val < 0 {
		return 0.0
	}
	scoreRules := []float64{0.5, 1, 2, 4, 5}
	for i, limit := range scoreRules {
		if val <= limit {
			score := float64(5 - i)
			return score * 2
		}
	}

	return 10.0
}

/**
趋势表示最近的市场情况，市值，pe的下降，价格更为便宜
*/
func (this *StatScoreService) calTrendScore(statScore *entity.StatScore) {
	score := this.increaseScore(statScore.LastPctChgMarketValue)
	statScore.TrendScore += score
	score = this.peScore(statScore.LastMeanPe)
	statScore.TrendScore += score
	score = this.pegScore(statScore.LastMeanPeg)
	statScore.TrendScore += score
	statScore.TrendScore = statScore.TrendScore * 10 / 3
	if statScore.TrendScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*最近市值下降很快\n"
	} else if statScore.TrendScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*最近市值下降快\n"
	} else if statScore.TrendScore < 20 {
		statScore.BadTip = statScore.BadTip + "*最近市值稳定\n"
	}
	//statScore.TotalScore += statScore.TrendScore
}

/**
景气度表示最新的业绩情况，销售和利润增长
*/
func (this *StatScoreService) calProsScore(statScore *entity.StatScore) {
	score := this.increaseScore(statScore.LastYoyDeduNp)
	statScore.ProsScore += score
	score = this.increaseScore(statScore.LastYoySales)
	statScore.ProsScore += score
	score = this.increaseScore(statScore.LastOrLastMonth)
	statScore.ProsScore += score
	score = this.increaseScore(statScore.LastNpLastMonth)
	statScore.ProsScore += score
	statScore.ProsScore = statScore.ProsScore * 2.5
	if statScore.ProsScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*最新业绩很好\n"
	} else if statScore.ProsScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*最新业绩好\n"
	} else if statScore.ProsScore < 20 {
		statScore.BadTip = statScore.BadTip + "*最新业绩差\n"
	}
	statScore.TotalScore += statScore.ProsScore
}

/**
累计增长
*/
func (this *StatScoreService) calAccScore(statScore *entity.StatScore) {
	score := this.increaseScore(statScore.AccPctChgMarketValue)
	statScore.AccScore += score
	score = this.increaseScore(statScore.AccYoyDeduNp)
	statScore.AccScore += score
	score = this.increaseScore(statScore.AccYoySales)
	statScore.AccScore += score
	statScore.AccScore = statScore.AccScore * 10 / 3
	if statScore.AccScore > 90 {
		statScore.GoodTip = statScore.GoodTip + "*累计业绩很好\n"
	} else if statScore.AccScore > 70 {
		statScore.GoodTip = statScore.GoodTip + "*累计业绩好\n"
	} else if statScore.AccScore < 20 {
		statScore.BadTip = statScore.BadTip + "*累计业绩差\n"
	}
	statScore.TotalScore += statScore.AccScore
}

func (this *StatScoreService) calRiskScore(statScore *entity.StatScore) {
	score := this.reportNumberScore(statScore)
	if score <= 2 {
		statScore.RiskScore = 0.0
	} else {
		score = this.industryScore(statScore)
		statScore.RiskScore += score
		score = this.areaScore(statScore)
		statScore.RiskScore += score
		score = this.marketScore(statScore)
		statScore.RiskScore += score
		score = this.listScore(statScore)
		statScore.RiskScore += score
		statScore.RiskScore = statScore.RiskScore * 2.5
		statScore.TotalScore += statScore.RiskScore
	}
}

func (this *StatScoreService) reportNumberScore(statScore *entity.StatScore) float64 {
	reportNumber := statScore.ReportNumber
	term := statScore.Term
	if term == 0 {
		term, _ = stock.DiffYear(statScore.StartDate, statScore.EndDate)
	}
	if term == 0 {
		logger.Sugar.Errorf("term is 0")
		return 5.0
	}
	pct := float64(reportNumber) / float64(term*4)
	scoreRules := []float64{0.5, 0.6, 0.7, 0.9, 1}
	score := 10.0
	for i, limit := range scoreRules {
		if pct <= limit {
			score = float64(i) * 2
			break
		}
	}
	if score <= 2 {
		statScore.BadTip = statScore.BadTip + "*业绩报告极度缺失\n"
	} else if score <= 4 {
		statScore.BadTip = statScore.BadTip + "*业绩报告缺失\n"
	}

	return score
}

func (this *StatScoreService) industryScore(statScore *entity.StatScore) float64 {
	blackRules := []string{"中成药", "渔业", "种植业"}
	grayRules := []string{"航空", "工程机械"}
	greenRules := []string{"医疗保健", "互联网", "IT设备", "生物制药", "半导体", "元器件", "环境保护"}
	industry := statScore.Industry
	for _, black := range blackRules {
		if industry == black {
			statScore.BadTip = statScore.BadTip + "*高风险行业\n"
			return 0.0
		}
	}
	for _, gray := range grayRules {
		if industry == gray {
			statScore.BadTip = statScore.BadTip + "*风险行业\n"
			return 2.0
		}
	}
	for _, green := range greenRules {
		if industry == green {
			return 7.0
		}
	}

	return 5.0
}

func (this *StatScoreService) areaScore(statScore *entity.StatScore) float64 {
	blackRules := []string{"辽宁", "新疆", "黑龙江", "三板股", "吉林", "内蒙"}
	grayRules := []string{"甘肃", "河北", "青海", "宁夏"}
	greenRules := []string{"深圳", "广东", "江苏", "浙江", "上海"}
	area := statScore.Area
	for _, black := range blackRules {
		if area == black {
			statScore.BadTip = statScore.BadTip + "*高风险地区\n"
			return 0.0
		}
	}
	for _, gray := range grayRules {
		if area == gray {
			statScore.BadTip = statScore.BadTip + "*风险地区\n"
			return 2.0
		}
	}
	for _, green := range greenRules {
		if area == green {
			return 7.0
		}
	}

	return 5.0
}

func (this *StatScoreService) marketScore(statScore *entity.StatScore) float64 {
	blackRules := []string{"深交所风险警示板", "老三板", "科创板", "新三板", "上交所科创板", "北交所"}
	grayRules := []string{"中小板"}
	greenRules := []string{}
	market := statScore.Market
	for _, black := range blackRules {
		if market == black {
			statScore.BadTip = statScore.BadTip + "*高风险板块\n"
			return 0.0
		}
	}
	for _, gray := range grayRules {
		if market == gray {
			statScore.BadTip = statScore.BadTip + "*风险板块\n"
			return 2.0
		}
	}
	for _, green := range greenRules {
		if market == green {
			return 7.0
		}
	}

	return 5.0
}

func (this *StatScoreService) listScore(statScore *entity.StatScore) float64 {
	listStatus := statScore.ListStatus
	listDate := statScore.ListDate
	if listStatus == "D" || listStatus == "P" {
		statScore.BadTip = statScore.BadTip + "*退市\n"
		return 0
	}
	if listDate == 0 {
		statScore.BadTip = statScore.BadTip + "*上市时间太短\n"
		return 0
	}
	var t = time.Now()
	diff := int64(t.Year()) - listDate/10000
	if diff > 10 {
		diff = 10
	}
	if diff < 3 {
		statScore.BadTip = statScore.BadTip + "*上市时间太短\n"
	} else if diff < 5 {
		statScore.BadTip = statScore.BadTip + "*上市时间较短\n"
	}

	return float64(diff)
}

/**
计算股票季度业绩统计评分数据中位数，最大值和最小值，并返回结果，方便进行去极值和标准化
*/
func (this *StatScoreService) CreateScorePercentile() (int64, error) {
	jsonHeads := []string{"percentile_risk_score", "percentile_stable_score", "percentile_acc_score", "percentile_pros_score", "percentile_trend_score", "percentile_increase_score",
		"percentile_corr_score", "percentile_operation_score", "percentile_price_score", "percentile_total_score"}
	jsonMap, _, _ := stock.GetJsonMap(entity.StatScore{})
	sql := "update stk_statscore ss set "
	updateFields := ""
	selectFields := ""
	i := 0
	for _, jsonHead := range jsonHeads {
		fieldname := jsonMap[jsonHead]
		originFieldName := strings.TrimPrefix(fieldname, "Percentile")
		originJsonHead := strings.TrimPrefix(jsonHead, "percentile_")
		if i > 0 {
			updateFields = updateFields + ","
			selectFields = selectFields + ","
		}
		updateFields = updateFields + fieldname + "=case when stat.max_" + originJsonHead + "!=stat.min_" + originJsonHead
		updateFields = updateFields + " then (ss." + originFieldName + "-stat.min_" + originJsonHead + ")/"
		updateFields = updateFields + "(stat.max_" + originJsonHead + "-stat.min_" + originJsonHead + ") else 0 end"
		selectFields = selectFields + "max(" + originFieldName + ") as max_" + originJsonHead + ",min(" + originFieldName + ") as min_" + originJsonHead
		i++
	}
	sql = sql + updateFields + " from (select term," + selectFields
	sql += " from stk_statscore where 1=1"
	sql = sql + " group by term"
	sql += ") stat where stat.term=ss.term"
	paras := make([]interface{}, 0)
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
	service.GetSession().Sync(new(entity.StatScore))
	statScoreService.OrmBaseService.GetSeqName = statScoreService.GetSeqName
	statScoreService.OrmBaseService.FactNewEntity = statScoreService.NewEntity
	statScoreService.OrmBaseService.FactNewEntities = statScoreService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("statScore", statScoreService)
}
