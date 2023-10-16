package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"strings"
)

/**
任何买卖点都是有以下几个部分组成：记录在inoutpoint表中，每个单独的条件记录在filtercond表中
1.走势(Trend)：表示股票的最近走势，分成三种，上涨，下跌，盘整，组合以后可以是下跌进入盘整，盘整进入上涨等共6种，体现在公式上代表accpctchgclose值
2.时机(Occasion)：表示当前的各条均线状态，最关键的是均线汇聚，均线的差，均线的交叉，体现在maclose
3.当天的价量(Situation)：涨跌，红绿，涨跌幅，上影线，下影线，宽幅震荡，窄幅震荡，放量，缩量，上穿均线，下穿均线，体现在close，preclose，maclose
4.后续跟进(Following)：第二天的走势，代表当天的判断是否稳固，某些特殊情形也看第三天，比如连续两天的小阳和连续小阴
5.预测准确性(Forecast)：分位短期和长期的准确性，体现在futurepctchgclose
上涨趋势中，越近的均线越高，长期均线向上走，站上13日线开始，直到稳定跌破13线结束
下降趋势中，越近的均线越低，长期均线向下走，13日线跌破开始，直到稳定站上13日线结束
如何区分上涨中的回调和下降，以及下降过程中的反弹与上涨，就是看长期均线的走势，回调时13日还是上走，反弹时13日还是向下
上涨初期：上涨趋势，34下穿55线之前，或者34高出55线不多为初期
下跌初期：下跌趋势，34上穿55线之前，或者55高出34线不多为初期
*/

type InOutPoint struct {
	Data      []*entity.DayLine        `json:"data,omitempty"`
	CondValue []map[string]interface{} `json:"cond_value,omitempty"`
	Count     int64                    `json:"count,omitempty"`
}

func (svc *DayLineService) FindAllInOutEvent(ts_code string, eventCode string, startDate int64, endDate int64, compareValue float64, condNum int) map[string][]*entity.DayLine {
	ecs := strings.Split(eventCode, ",")
	InOutEventMap := make(map[string][]*entity.DayLine)
	for _, ec := range ecs {
		inOutPoint, _ := svc.FindInOutEvent(ts_code, 0, ec, nil, startDate, endDate, compareValue, condNum, 0, 0, 0)
		if inOutPoint != nil {
			InOutEventMap[ec] = inOutPoint.Data
		}
	}

	return InOutEventMap
}

func (svc *DayLineService) FindInOutEvent(ts_code string, tradedate int64, eventCode string, filterParas []interface{}, startDate int64, endDate int64, compareValue float64, condNum int, from int, limit int, count int64) (*InOutPoint, error) {
	efs, ok := GetEventFilterService().GetCacheEventFilter()[eventCode]
	if !ok {
		return nil, errors.New("eventCode err")
	}
	con, fields := svc.buildConds(efs, compareValue, condNum)
	if con == "" {
		return nil, errors.New("content err")
	}
	inOutPoint, err := svc.FindFlexPoint(ts_code, tradedate, fields, con, filterParas, startDate, endDate, from, limit, count)

	return inOutPoint, err
}

/*
*
compareValue是负数，表示放宽比较值，缺省是-0.01；num是负数，表示放宽条件个数，缺省是0
*/
func (svc *DayLineService) buildConds(efs []*entity.EventFilter, compareValue float64, condNum int) (string, []string) {
	cond := ""
	fields := make([]string, 0)
	k := 0
	//默认>0，计算汇总值的时候用不上
	//condParas := ">0"
	for _, ef := range efs {
		//condCode := ef.CondCode                   //rise
		codeAlias := ef.CodeAlias                 //minRise
		condContent := "(" + ef.CondContent + ")" //pctchgclose-0.02
		condName := ef.CondName                   //涨幅
		condAlias := ef.CondAlias                 //最低涨幅
		//完整表达式的字段，带参数
		fields = append(fields, "'"+condContent+"' as "+codeAlias+"_cond")
		value := condContent + " as " + codeAlias
		//表达式的计算值
		fields = append(fields, value)
		if condName != "" {
			name := "'" + condName + "' as " + codeAlias + "_name"
			fields = append(fields, name)
		}
		if condAlias != "" {
			alias := "'" + condAlias + "' as " + codeAlias + "_alias"
			fields = append(fields, alias)
		}
		condParas := ef.CondParas
		if condParas == "" {
			condParas = ">" + fmt.Sprint(compareValue-0.01)
		}
		result := "case when " + condContent + condParas + " then 1 else 0 end as " + codeAlias + "_result"
		fields = append(fields, result)
		if k != 0 {
			cond = cond + "+"
		}
		cond = cond + "(case when " + condContent + condParas + " then 1 else 0 end)"
		k++
	}
	cond = "(" + cond + ")>=" + fmt.Sprint(k+condNum)

	return cond, fields
}

func (svc *DayLineService) FindFlexPoint(ts_code string, tradedate int64, fields []string, eventContent string, filterParas []interface{}, startDate int64, endDate int64, from int, limit int, count int64) (*InOutPoint, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	dayLines := make([]*entity.DayLine, 0)
	conds += " and ma3close is not null and ma3close!=0 and (high-low)!=0"
	inOutPoint := &InOutPoint{}
	var err error
	if tradedate != 0 {
		conds = conds + " and tradedate=?"
		paras = append(paras, tradedate)
	}
	if eventContent != "" {
		conds = conds + " and " + eventContent
		if filterParas != nil && len(filterParas) > 0 {
			paras = append(paras, filterParas...)
		}
	}
	if startDate != 0 {
		conds = conds + " and tradedate>?"
		paras = append(paras, startDate)
	}
	if endDate != 0 {
		conds += " and tradedate<=?"
		paras = append(paras, endDate)
	}
	condiBean := &entity.DayLine{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, err
		}
		inOutPoint.Count = count
	}

	err = svc.Find(&dayLines, nil, "tscode,tradedate desc", from, limit, conds, paras...)
	if err != nil {
		return inOutPoint, err
	}
	inOutPoint.Data = dayLines
	if fields == nil || len(fields) == 0 {
		return inOutPoint, nil
	}
	sql := "select tscode as ts_code,tradedate as trade_date," + strings.Join(fields, ",") + " from stk_dayline where " + conds + " order by tscode,tradedate desc"
	if limit > 0 {
		sql = sql + " limit " + fmt.Sprint(limit)
	}
	sql = sql + " offset " + fmt.Sprint(from)
	results, err := svc.Query(sql, paras...)
	if err != nil {
		return inOutPoint, err
	}
	reals := make([]map[string]interface{}, 0)
	for _, result := range results {
		real := make(map[string]interface{})
		for f, r := range result {
			strVal := string(r)
			if f == "trade_date" {
				v, _ := convert.ToObject(strVal, "int64")
				tradedate := v.(int64)
				real[f] = tradedate
			} else if f == "ts_code" || strings.HasSuffix(f, "_cond") ||
				strings.HasSuffix(f, "_name") ||
				strings.HasSuffix(f, "_alias") ||
				strings.HasSuffix(f, "_paras") {
				real[f] = strVal
			} else {
				v, _ := convert.ToObject(strVal, "float64")
				r := v.(float64)
				real[f] = r
			}
		}
		reals = append(reals, real)
	}
	inOutPoint.CondValue = reals

	return inOutPoint, nil
}
