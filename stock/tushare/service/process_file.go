package service

import (
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"sync"
)

type PreprocessPara struct {
	entity.ProcessLog
	locker *sync.Mutex
}

var DataType = []string{"lday", "minline", "fzline"}
var PreprocessType = []string{"parse", "minmax", "standard"}

var minline_header = []string{"id", "trade_date", "trade_minute", "open", "high", "low", "close", "amount", "vol"}

const date_key = "trade_date"
const minute_key = "trade_minute"

// 输出序列之间的间隔
const sequence_stride = 240

// 过去5天，共5*240=1200个数据点
var past = 5
var batch_size = 128

// 序列里面每个元素的间隔
var step = 1
var epochs = 10

/**
解析原始stock的数据
*/
//func ParsePath(src string, target string) error {
//	src = stock.StdPath(src)
//	target = stock.StdPath(target)
//
//	daySvc := service.GetDayLineService()
//	err := daySvc.ParsePath(src+"\\sz\\lday", target+"\\lday")
//	if err != nil {
//		return err
//	}
//	err = daySvc.ParsePath(src+"\\sh\\lday", target+"\\lday")
//	if err != nil {
//		return err
//	}
//	minSvc := service.GetMinLineService()
//	err = minSvc.ParsePath(src+"\\sz\\minline", target+"\\minline")
//	if err != nil {
//		return err
//	}
//	err = minSvc.ParsePath(src+"\\sh\\minline", target+"\\minline")
//	if err != nil {
//		return err
//	}
//	fzSvc := service.GetFminLineService()
//	err = fzSvc.ParsePath(src+"\\sz\\fzline", target+"\\fzline")
//	if err != nil {
//		return err
//	}
//	err = minSvc.ParsePath(src+"\\sh\\fzline", target+"\\fzline")
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func Split(x_data dataframe.DataFrame, y_data dataframe.DataFrame) (dataframe.DataFrame, dataframe.DataFrame, dataframe.DataFrame, dataframe.DataFrame) {
	// 选择列
	x_features := x_data.Select([]int{1, 2, 3, 4, 5, 6})
	split_fraction := 0.715
	train_split := int(split_fraction*(float64(x_data.Nrow()))/sequence_stride) * sequence_stride
	// 划分训练集和验证集
	x_train_data := x_features.Filter(
		dataframe.F{0, "id", series.Less, train_split},
	)
	x_val_data := x_features.FilterAggregation(
		dataframe.And,
		dataframe.F{0, "id", series.Less, x_features.Nrow() - sequence_stride},
		dataframe.F{0, "id", series.GreaterEq, train_split - (past*sequence_stride - sequence_stride)},
	)

	// 选择列
	y_features := y_data.Select([]int{1, 2, 3, 4, 5, 6})
	start := past
	end := int(train_split/sequence_stride) + 1
	// 划分训练集和验证集
	y_train_data := y_features.FilterAggregation(
		dataframe.And,
		dataframe.F{0, "id", series.GreaterEq, start},
		dataframe.F{0, "id", series.Less, end},
	)
	y_val_data := y_features.Filter(
		dataframe.F{0, "id", series.GreaterEq, end},
	)

	return x_train_data, x_val_data, y_train_data, y_val_data
}
