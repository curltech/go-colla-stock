package stock

import (
	reflect1 "github.com/curltech/go-colla-core/util/reflect"
	"math"
	"sort"
)

type AggregationType int

const (
	Aggregation_SUM AggregationType = iota + 1
	Aggregation_MAX
	Aggregation_MIN
	Aggregation_MEAN
	Aggregation_MEDIAN
	Aggregation_STDDEV
	Aggregation_RSD
	Aggregation_COV
	Aggregation_COR
	Aggregation_COUNT
	Aggregation_STD
	Aggregation_MINMAX
	Aggregation_PCA
)

var AggregationTypes = []AggregationType{
	Aggregation_SUM,
	Aggregation_MAX,
	Aggregation_MIN,
	Aggregation_MEDIAN,
	Aggregation_MEAN,
	Aggregation_STDDEV,
	Aggregation_RSD,
	Aggregation_COV,
	Aggregation_COR,
	Aggregation_COUNT,
	Aggregation_STD,
	Aggregation_MINMAX,
	Aggregation_PCA,
}

type Stat struct {
	Sum               interface{}
	Min               interface{}
	Max               interface{}
	Count             interface{}
	Median            []interface{}
	Mean              interface{}
	Stddev            interface{}
	Mad               interface{}
	Covs              map[string]interface{}
	Cors              map[string]interface{}
	Rsd               interface{}
	Pca               []interface{}
	Std               []interface{}
	Minmax            []interface{}
	Data              []interface{}
	Colnames          []string
	jsonMap           map[string]string
	typMap            map[string]string
	std_colnames      []string
	reserved_colnames []string
	sample            interface{}
}

func CreateStat(data []interface{}, std_colnames []string) *Stat {
	if len(data) <= 0 {
		return nil
	}

	jsonMap, typMap, _ := GetJsonMap(data[0])
	stat := &Stat{
		Data:         data,
		std_colnames: std_colnames,
		sample:       data[0],
		jsonMap:      jsonMap,
		typMap:       typMap,
	}

	return stat
}

func (this *Stat) CalSum() interface{} {
	if this.Sum != nil {
		return this.Sum
	}
	this.Sum = reflect1.New(this.sample)
	this.Max = reflect1.New(this.sample)
	this.Min = reflect1.New(this.sample)
	this.Count = reflect1.New(this.sample)
	i := 0
	for _, qp := range this.Data {
		for _, colname := range this.std_colnames {
			fieldname := this.jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := this.typMap[colname]
			if fieldtyp == "" {
				continue
			}
			v, _ := reflect1.GetValue(qp, fieldname)
			val, _ := v.(float64)

			if i == 0 {
				reflect1.SetValue(this.Sum, fieldname, val)
				reflect1.SetValue(this.Max, fieldname, val)
				reflect1.SetValue(this.Min, fieldname, val)
				reflect1.SetValue(this.Count, fieldname, 1.0)
			} else {
				v, _ = reflect1.GetValue(this.Max, fieldname)
				aggreVal, _ := v.(float64)
				if val > aggreVal {
					reflect1.SetValue(this.Max, fieldname, val)
				}
				v, _ = reflect1.GetValue(this.Min, fieldname)
				aggreVal, _ = v.(float64)
				if val < aggreVal {
					reflect1.SetValue(this.Min, fieldname, val)
				}
				v, _ = reflect1.GetValue(this.Sum, fieldname)
				aggreVal, _ = v.(float64)
				reflect1.SetValue(this.Sum, fieldname, aggreVal+val)
				v, _ = reflect1.GetValue(this.Count, fieldname)
				aggreVal, _ = v.(float64)
				reflect1.SetValue(this.Count, fieldname, aggreVal+1)
			}
		}
		i++
	}
	return this.Sum
}

func (this *Stat) CalMedian() []interface{} {
	if this.Median != nil {
		return this.Median
	}
	this.Median = make([]interface{}, 3)
	this.Median[0] = reflect1.New(this.sample)
	this.Median[1] = reflect1.New(this.sample)
	this.Median[2] = reflect1.New(this.sample)
	this.CalSum()
	sorted := make(map[string][]float64, 0)
	for i := 0; i < len(this.Data); i++ {
		qp := this.Data[i]
		for j := 0; j < len(this.std_colnames); j++ {
			colname := this.std_colnames[j]
			fieldname := this.jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := this.typMap[colname]
			if fieldtyp == "" {
				continue
			}
			v, _ := reflect1.GetValue(qp, fieldname)
			val, _ := v.(float64)
			if i == 0 {
				sorted[fieldname] = make([]float64, 0)
			}
			sorted[fieldname] = append(sorted[fieldname], val)
		}
	}
	//求中位数
	for j := 0; j < len(this.std_colnames); j++ {
		colname := this.std_colnames[j]
		fieldname := this.jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := this.typMap[colname]
		if fieldtyp == "" {
			continue
		}
		arr := sorted[fieldname]
		sort.Float64s(arr)
		length := len(arr)
		if length%4 == 0 {
			v := (arr[length/4-1] + arr[length/4]) / 2.0
			reflect1.SetValue(this.Median[0], fieldname, v)
		} else {
			v := arr[int(math.Floor(float64(length/2)))]
			reflect1.SetValue(this.Median[0], fieldname, v)
		}
		if len(arr)%2 == 0 {
			v := (arr[length/2-1] + arr[length/2]) / 2.0
			reflect1.SetValue(this.Median[1], fieldname, v)
		} else {
			v := arr[int(math.Floor(float64(length/2)))]
			reflect1.SetValue(this.Median[1], fieldname, v)
		}
		if len(arr)*3%4 == 0 {
			v := (arr[length*3/4-1] + arr[length*3/4]) / 2.0
			reflect1.SetValue(this.Median[2], fieldname, v)
		} else {
			v := arr[int(math.Floor(float64(length*3/4)))]
			reflect1.SetValue(this.Median[2], fieldname, v)
		}
	}

	return this.Median
}

func (this *Stat) CalMean() interface{} {
	if this.Mean != nil {
		return this.Mean
	}
	this.Mean = reflect1.New(this.sample)
	this.CalSum()
	for _, colname := range this.std_colnames {
		fieldname := this.jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := this.typMap[colname]
		if fieldtyp == "" {
			continue
		}
		v, _ := reflect1.GetValue(this.Sum, fieldname)
		sumVal, _ := v.(float64)
		v, _ = reflect1.GetValue(this.Count, fieldname)
		countVal, _ := v.(float64)
		if countVal != 0.0 {
			reflect1.SetValue(this.Mean, fieldname, sumVal/countVal)
		}
	}
	return this.Mean
}

func (this *Stat) CalStddev() interface{} {
	if this.Stddev != nil {
		return this.Stddev
	}
	this.Stddev = reflect1.New(this.sample)
	this.CalMean()
	for _, qp := range this.Data {
		for _, colname := range this.std_colnames {
			fieldname := this.jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := this.typMap[colname]
			if fieldtyp == "" {
				continue
			}
			val := 0.0
			v, _ := reflect1.GetValue(qp, fieldname)
			val, _ = v.(float64)
			v, _ = reflect1.GetValue(this.Mean, fieldname)
			meanVal, _ := v.(float64)
			diff := val - meanVal
			v, _ = reflect1.GetValue(this.Stddev, fieldname)
			stddevVal, _ := v.(float64)
			reflect1.SetValue(this.Stddev, fieldname, stddevVal+diff*diff)
		}
	}
	for _, colname := range this.std_colnames {
		fieldname := this.jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := this.typMap[colname]
		if fieldtyp == "" {
			continue
		}
		v, _ := reflect1.GetValue(this.Stddev, fieldname)
		stddevVal, _ := v.(float64)
		if len(this.Data) != 1 {
			stddevVal = stddevVal / float64(len(this.Data)-1)
		}
		stddevVal = math.Sqrt(stddevVal)
		reflect1.SetValue(this.Stddev, fieldname, stddevVal)
	}
	return this.Stddev
}

/**
中位数偏差的中位数
*/
func (this *Stat) CalMad() interface{} {
	if this.Mad != nil {
		return this.Mad
	}
	this.Mad = reflect1.New(this.sample)
	this.CalMedian()
	sorted := make(map[string][]float64, 0)
	i := 0
	for _, qp := range this.Data {
		for _, colname := range this.std_colnames {
			fieldname := this.jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := this.typMap[colname]
			if fieldtyp == "" {
				continue
			}
			val := 0.0
			v, _ := reflect1.GetValue(qp, fieldname)
			val, _ = v.(float64)
			v, _ = reflect1.GetValue(this.Median[1], fieldname)
			medianVal, _ := v.(float64)
			diff := math.Abs(val - medianVal)
			if i == 0 {
				sorted[fieldname] = make([]float64, 0)
			}
			sorted[fieldname] = append(sorted[fieldname], diff)
		}
		i++
	}
	//求中位数
	for j := 0; j < len(this.std_colnames); j++ {
		colname := this.std_colnames[j]
		fieldname := this.jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := this.typMap[colname]
		if fieldtyp == "" {
			continue
		}
		arr := sorted[fieldname]
		sort.Float64s(arr)
		if len(arr)%2 == 0 {
			v := (arr[len(arr)/2-1] + arr[len(arr)/2]) / 2.0
			reflect1.SetValue(this.Mad, fieldname, v)
		} else {
			v := arr[int(math.Floor(float64(len(arr)/2)))]
			reflect1.SetValue(this.Mad, fieldname, v)
		}
	}
	return this.Mad
}

func (this *Stat) CalRsd() interface{} {
	if this.Rsd != nil {
		return this.Rsd
	}
	this.Rsd = reflect1.New(this.sample)
	this.CalStddev()
	for _, colname := range this.std_colnames {
		fieldname := this.jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := this.typMap[colname]
		if fieldtyp == "" {
			continue
		}
		v, _ := reflect1.GetValue(this.Stddev, fieldname)
		stddevVal, _ := v.(float64)
		v, _ = reflect1.GetValue(this.Mean, fieldname)
		meanVal, _ := v.(float64)
		if !Equal(meanVal, 0) {
			reflect1.SetValue(this.Rsd, fieldname, stddevVal/meanVal)
		}
	}
	return this.Rsd
}

func (this *Stat) CalCov(x_colname string) interface{} {
	var cov interface{}
	if this.Covs != nil {
		cov, ok := this.Covs[x_colname]
		if ok {
			return cov
		}
	} else {
		this.Covs = make(map[string]interface{}, 0)
	}
	cov = reflect1.New(this.sample)
	this.Covs[x_colname] = cov
	this.CalMean()
	x_fieldname := this.jsonMap[x_colname]
	for i := 0; i < len(this.Data); i++ {
		qp := this.Data[i]
		v, _ := reflect1.GetValue(qp, x_fieldname)
		x_val, _ := v.(float64)
		//if !ok {
		//	logger.Sugar.Errorf("error: fieldname:%v transfer to float64", x_colname)
		//}
		for j := 0; j < len(this.std_colnames); j++ {
			colname := this.std_colnames[j]
			fieldname := this.jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := this.typMap[colname]
			if fieldtyp == "" {
				continue
			}
			v, _ := reflect1.GetValue(qp, fieldname)
			y_val, _ := v.(float64)
			v, _ = reflect1.GetValue(this.Mean, x_fieldname)
			x_meanVal, _ := v.(float64)
			v, _ = reflect1.GetValue(this.Mean, fieldname)
			meanVal, _ := v.(float64)
			if !Equal(y_val, meanVal) {
				diff := (x_val - x_meanVal) * (y_val - meanVal)
				if i == 0 {
					reflect1.SetValue(cov, fieldname, diff)
				} else {
					v, _ = reflect1.GetValue(cov, fieldname)
					covVal, _ := v.(float64)
					reflect1.SetValue(cov, fieldname, covVal+diff)
				}
			}
		}
	}
	for j := 0; j < len(this.std_colnames); j++ {
		colname := this.std_colnames[j]
		fieldname := this.jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := this.typMap[colname]
		if fieldtyp == "" {
			continue
		}
		v, _ := reflect1.GetValue(cov, fieldname)
		covVal, _ := v.(float64)
		if len(this.Data) != 1 {
			covVal = covVal / float64(len(this.Data)-1)
		}
		reflect1.SetValue(cov, fieldname, covVal)
	}
	return cov
}

func (this *Stat) CalCor(x_colname string) interface{} {
	var cor interface{}
	if this.Cors != nil {
		cor, ok := this.Cors[x_colname]
		if ok {
			return cor
		}
	} else {
		this.Cors = make(map[string]interface{}, 0)
	}
	cor = reflect1.New(this.sample)
	this.Cors[x_colname] = cor
	cov := this.CalCov(x_colname)
	x_fieldname := this.jsonMap[x_colname]
	for j := 0; j < len(this.std_colnames); j++ {
		colname := this.std_colnames[j]
		fieldname := this.jsonMap[colname]
		if fieldname == "" {
			continue
		}
		fieldtyp := this.typMap[colname]
		if fieldtyp == "" {
			continue
		}
		v, _ := reflect1.GetValue(this.Stddev, x_fieldname)
		x_stddevVal, _ := v.(float64)
		v, _ = reflect1.GetValue(this.Stddev, fieldname)
		stddevVal, _ := v.(float64)
		v, _ = reflect1.GetValue(cov, fieldname)
		covVal, _ := v.(float64)
		if !Equal(x_stddevVal, 0) && !Equal(stddevVal, 0) {
			reflect1.SetValue(cor, fieldname, covVal/(x_stddevVal*stddevVal))
		}
	}
	return cor
}

func (this *Stat) CalPca() []interface{} {
	if this.Pca != nil {
		return this.Pca
	}
	this.Pca = make([]interface{}, 0)
	for i := 0; i < len(this.std_colnames); i++ {
		x_colname := this.std_colnames[i]
		cor := this.CalCor(x_colname)
		this.Pca = append(this.Pca, cor)
	}

	return this.Pca
}

func (this *Stat) CalCors() map[string]interface{} {
	if this.Cors == nil {
		this.CalPca()
	}
	return this.Cors
}

func (this *Stat) CalStd(reserved_colnames []string, isWinsorize bool) ([]interface{}, []interface{}) {
	if this.Std != nil {
		return this.Std, this.Minmax
	}
	this.CalStddev()
	this.Std = make([]interface{}, 0)
	this.Minmax = make([]interface{}, 0)
	for _, qp := range this.Data {
		std := reflect1.New(this.sample)
		minmax := reflect1.New(this.sample)
		for _, colname := range reserved_colnames {
			fieldname := this.jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := this.typMap[colname]
			if fieldtyp == "" {
				continue
			}
			v, _ := reflect1.GetValue(qp, fieldname)
			reflect1.SetValue(std, fieldname, v)
			reflect1.SetValue(minmax, fieldname, v)
		}
		for _, colname := range this.std_colnames {
			fieldname := this.jsonMap[colname]
			if fieldname == "" {
				continue
			}
			fieldtyp := this.typMap[colname]
			if fieldtyp == "" {
				continue
			}
			val := 0.0
			v, _ := reflect1.GetValue(qp, fieldname)
			val, _ = v.(float64)
			v, _ = reflect1.GetValue(this.Mean, fieldname)
			meanVal, _ := v.(float64)
			v, _ = reflect1.GetValue(this.Stddev, fieldname)
			stddevVal, _ := v.(float64)
			stdVal := (val - meanVal) / stddevVal
			if stddevVal != 0 {
				if isWinsorize {
					winsorizeVal := 3 * stddevVal
					diff := val - meanVal
					if diff < -winsorizeVal {
						stdVal = (-winsorizeVal - meanVal) / stddevVal
					} else if diff > winsorizeVal {
						stdVal = (winsorizeVal - meanVal) / stddevVal
					}
				}
				reflect1.SetValue(std, fieldname, stdVal)
			}
			v, _ = reflect1.GetValue(this.Min, fieldname)
			minVal, _ := v.(float64)
			v, _ = reflect1.GetValue(this.Max, fieldname)
			maxVal, _ := v.(float64)
			if maxVal != minVal {
				if isWinsorize {
					maxDiff := (maxVal - minVal) * 0.025
					maxVal = maxVal - maxDiff
					minVal = minVal + maxDiff
					if val < minVal {
						val = minVal
					} else if val > maxVal {
						val = maxVal
					}
				}
				minmaxVal := (val - minVal) / (maxVal - minVal)
				reflect1.SetValue(minmax, fieldname, minmaxVal)
			}
		}
		this.Std = append(this.Std, std)
		this.Minmax = append(this.Minmax, minmax)
	}

	return this.Std, this.Minmax
}
