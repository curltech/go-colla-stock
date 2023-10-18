package service

import (
	"fmt"
	"github.com/curltech/go-colla-core/config"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"os"
	"strings"
)

// FindModelData 获取某只股票最新的日期
func (svc *DayLineService) FindModelData(ts_code string, industry string, startDate int64) ([]*entity.DayLine, error) {
	var conds string
	var paras []interface{}
	if ts_code != "" {
		conds, paras = stock.InBuildStr("tscode", ts_code, ",")
	} else if industry != "" {
		conds = "tscode in (select tscode from stk_share where industry=?)"
		paras = append(paras, industry)
	}
	cond := &entity.DayLine{}
	dayLines := make([]*entity.DayLine, 0)
	if startDate > 0 {
		conds = conds + " and tradedate >= ?"
		paras = append(paras, startDate)
	}
	err := svc.Find(&dayLines, cond, "tscode,tradedate", 0, 0, conds, paras...)

	return dayLines, err
}

func (svc *DayLineService) WriteAllFile(startDate int64) error {
	v, _ := config.Get("stock.src")
	src := v.(string)
	src = src + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(stock.CurrentDate())
	err := stock.Mkdir(src)
	if err != nil {
		return err
	}
	routinePool := thread.CreateRoutinePool(10, svc.AsyncWriteFile, nil)
	defer routinePool.Release()
	tsCodes, _ := GetShareService().GetCacheShare()
	for _, tsCode := range tsCodes {
		para := make([]interface{}, 0)
		para = append(para, src)
		para = append(para, tsCode)
		para = append(para, startDate)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)
	return nil
}

func (svc *DayLineService) AsyncWriteFile(para interface{}) {
	src := (para.([]interface{}))[0].(string)
	tsCode := (para.([]interface{}))[1].(string)
	startDate := (para.([]interface{}))[2].(int64)
	err := svc.WriteFile(src, tsCode, startDate)
	if err != nil {
		return
	}
}

func (svc *DayLineService) WriteFile(src string, tsCode string, startDate int64) error {
	if src == "" {
		v, _ := config.Get("stock.src")
		src = v.(string)
		src = src + string(os.PathSeparator) + fmt.Sprint(startDate) + "-" + fmt.Sprint(stock.CurrentDate())
		err := stock.Mkdir(src)
		if err != nil {
			return err
		}
	}
	dayLines, err := svc.FindModelData(tsCode, "", startDate)
	if err != nil {
		logger.Sugar.Errorf("%v FindModelData failure!", tsCode)
	}
	raw := "id,ts_code,trade_date,open,high,low,turnover" +
		",main_net_inflow,small_net_inflow,middle_net_inflow,large_net_inflow,super_net_inflow" +
		",pct_main_net_inflow,pct_small_net_inflow" +
		",pct_middle_net_inflow,pct_large_net_inflow,pct_super_net_inflow,pct_chg_open,pct_chg_high" +
		",pct_chg_low,pct_chg_close,pct_chg_amount,pct_chg_vol,ma3_close,ma5_close,ma10_close" +
		",ma13_close,ma20_close,ma21_close,ma30_close,ma34_close,ma55_close,ma60_close,ma90_close" +
		",ma120_close,ma144_close,ma233_close,ma240_close,max3_close,max5_close,max10_close,max13_close" +
		",max20_close,max21_close,max30_close,max34_close,max55_close,max60_close,max90_close,max120_close" +
		",max144_close,max233_close,max240_close,min3_close,min5_close,min10_close,min13_close,min20_close" +
		",min21_close,min30_close,min34_close,min55_close,min60_close,min90_close,min120_close,min144_close" +
		",min233_close,min240_close,before1_ma3_close,before1_ma5_close,before1_ma10_close" +
		",before1_ma13_close,before1_ma20_close,before1_ma21_close,before1_ma30_close,before1_ma34_close" +
		",before1_ma55_close,before1_ma60_close,before3_ma3_close,before3_ma5_close,before3_ma10_close" +
		",before3_ma13_close,before3_ma20_close,before3_ma21_close,before3_ma30_close,before3_ma34_close" +
		",before3_ma55_close,before3_ma60_close,before5_ma3_close,before5_ma5_close,before5_ma10_close" +
		",before5_ma13_close,before5_ma20_close,before5_ma21_close,before5_ma30_close,before5_ma34_close" +
		",before5_ma55_close,before5_ma60_close,acc3_pct_chg_close,acc5_pct_chg_close,acc10_pct_chg_close" +
		",acc13_pct_chg_close,acc20_pct_chg_close,acc21_pct_chg_close,acc30_pct_chg_close" +
		",acc34_pct_chg_close,acc55_pct_chg_close,acc60_pct_chg_close,acc90_pct_chg_close" +
		",acc120_pct_chg_close,acc144_pct_chg_close,acc233_pct_chg_close,acc240_pct_chg_close" +
		",future1_pct_chg_close,future3_pct_chg_close,future5_pct_chg_close,future10_pct_chg_close" +
		",future13_pct_chg_close,future20_pct_chg_close,future21_pct_chg_close,future30_pct_chg_close" +
		",future34_pct_chg_close,future55_pct_chg_close,future60_pct_chg_close,future90_pct_chg_close" +
		",future120_pct_chg_close,future144_pct_chg_close,future233_pct_chg_close,future240_pct_chg_close"
	heads := strings.Split(raw, ",")
	jsonMap, _, _ := stock.GetJsonMap(&entity.DayLine{})
	raw = raw + "\n"
	filename := src + string(os.PathSeparator) + tsCode + ".csv"
	err = os.WriteFile(filename, []byte(raw), 0644)
	if err != nil {
		logger.Sugar.Errorf("%v write file failure!", tsCode)
	}
	var file *os.File
	if len(dayLines) > 0 {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			logger.Sugar.Errorf("%v open file failure!", tsCode)
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {

			}
		}(file)
	}
	lineNum := 0
	for _, dayLine := range dayLines {
		if lineNum <= 5 {
			lineNum++
			continue
		}
		i := 0
		raw = ""
		shareNumber := dayLine.ShareNumber
		closePrice := dayLine.Close
		marketValue := shareNumber * closePrice
		amount := dayLine.Amount
		dayLine.Open = (dayLine.Open - closePrice) / closePrice
		dayLine.High = (dayLine.High - closePrice) / closePrice
		dayLine.Low = (dayLine.Low - closePrice) / closePrice
		dayLine.Vol = dayLine.Vol / shareNumber
		dayLine.Amount = amount / marketValue
		for _, colname := range heads {
			fieldname := jsonMap[colname]
			if colname == "turnover" {
				v, _ := reflect.GetValue(dayLine, fieldname)
				val := v.(float64)
				nval := 0.0
				if !stock.Equal(val, 0.0) {
					nval = val / 100
				}
				err = reflect.SetValue(dayLine, fieldname, nval)
				if err != nil {
					return err
				}
			} else if strings.HasSuffix(colname, "_close") && !strings.HasSuffix(colname, "pct_chg_close") && colname != "chg_close" {
				v, _ := reflect.GetValue(dayLine, fieldname)
				val := v.(float64)
				nval := 0.0
				if !stock.Equal(val, 0.0) {
					nval = (val - closePrice) / closePrice
				}
				err = reflect.SetValue(dayLine, fieldname, nval)
				if err != nil {
					return err
				}
			} else if strings.HasSuffix(colname, "_net_inflow") && !strings.HasPrefix(colname, "pct_") {
				v, _ := reflect.GetValue(dayLine, fieldname)
				val := v.(float64)
				nval := 0.0
				if !stock.Equal(val, 0.0) {
					nval = (val - amount) / amount
				}
				err = reflect.SetValue(dayLine, fieldname, nval)
				if err != nil {
					return err
				}
			}
			v, _ := reflect.GetValue(dayLine, fieldname)
			raw = raw + fmt.Sprint(v)
			if i < len(heads)-1 {
				raw = raw + ","
			}
			i++
		}
		lineNum++
		if lineNum < len(dayLines) {
			raw = raw + "\n"
		}
		if _, err = file.WriteString(raw); err != nil {
			logger.Sugar.Errorf("%v append file failure!", tsCode)
		} else {
			logger.Sugar.Infof("tscode:%v lineNum:%v", tsCode, lineNum)
		}
	}

	return nil
}
