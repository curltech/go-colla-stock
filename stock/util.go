package stock

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/util/collection"
	"github.com/curltech/go-colla-core/util/convert"
	"math"
	"os"
	"reflect"
	"strings"
	"time"
)

func BytesToInt64(bys []byte) int64 {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	binary.Read(bytebuff, binary.LittleEndian, &data)

	return int64(data)
}

func BytesToFloat64(bys []byte) float64 {
	bytebuff := bytes.NewBuffer(bys)
	var data float32
	binary.Read(bytebuff, binary.LittleEndian, &data)

	return float64(data)
}

func BytesToInt(bys []byte) uint16 {
	bytebuff := bytes.NewBuffer(bys)
	var data uint16
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return data
}

//
//func BytesToFloat64(bys []byte) float64 {
//	bytebuff := bytes.NewBuffer(bys)
//	var data float64
//	binary.Read(bytebuff, binary.LittleEndian, &data)
//	return data
//}

func StdPath(path string) string {
	p := strings.TrimSuffix(path, string(os.PathSeparator))
	p = strings.TrimSuffix(p, "/")

	return p
}

func Mkdir(path string) error {
	err := os.MkdirAll(path, 0766)
	return err
}

func Rename(path string, newpath string) error {
	err := os.Rename(path, newpath)
	return err
}

func CurrentDate() int64 {
	timeStr := time.Now().Format("20060102")
	v, _ := convert.ToObject(timeStr, "int64")
	return v.(int64)
}

func AddDay(tradeDate int64, days int) int64 {
	year := tradeDate / 10000
	month := (tradeDate - year*10000) / 100
	day := tradeDate - year*10000 - month*100
	timeTradeDate := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC).AddDate(0, 0, days)
	timeStr := timeTradeDate.Format("20060102")
	v, _ := convert.ToObject(timeStr, "int64")

	return v.(int64)
}

func CurrentMinute() int64 {
	t := time.Now()
	tradeMinute := t.Hour()*60 + t.Minute()

	return int64(tradeMinute)
}

func AddYear(qdate string, diff int) (string, error) {
	ss := strings.Split(qdate, "Q")
	if len(ss) == 2 {
		v, err := convert.ToObject(ss[0], "int")
		if err == nil {
			year := v.(int) + diff
			qdate := fmt.Sprint(year) + "Q" + ss[1]
			return qdate, nil
		} else {
			return "", err
		}
	}

	return "", errors.New("error qdate format")
}

func DiffYear(sqdate string, eqdate string) (int, error) {
	ss := strings.Split(sqdate, "Q")
	if len(ss) == 2 {
		es := strings.Split(eqdate, "Q")
		if len(es) == 2 {
			v, err := convert.ToObject(ss[0], "int")
			if err == nil {
				syear := v.(int)
				v, err = convert.ToObject(es[0], "int")
				if err == nil {
					eyear := v.(int)
					return eyear - syear, nil
				}
			}
		}
	}

	return 0, errors.New("error qdate format")
}

func AddQuarter(qdate string, diff int) (string, error) {
	ss := strings.Split(qdate, "Q")
	if len(ss) == 2 {
		v, err := convert.ToObject(ss[1], "int")
		if err == nil {
			quarter := v.(int) + diff
			if quarter > 4 {
				year := quarter / 4
				quarter = quarter % 4
				v, err := convert.ToObject(ss[0], "int")
				if err == nil {
					if quarter == 0 {
						year--
						quarter = 4
					}
					year = v.(int) + year
					qdate := fmt.Sprint(year) + "Q" + fmt.Sprint(quarter)
					return qdate, nil
				}
			} else if quarter < 1 {
				year := quarter / 4
				quarter = quarter % 4
				v, err := convert.ToObject(ss[0], "int")
				if err == nil {
					if quarter <= 0 {
						year--
						quarter = 4 + quarter
					}
					year = v.(int) + year
					qdate := fmt.Sprint(year) + "Q" + fmt.Sprint(quarter)
					return qdate, nil
				}
			} else {
				qdate := ss[0] + "Q" + fmt.Sprint(quarter)
				return qdate, nil
			}
		} else {
			return "", err
		}
	}

	return "", errors.New("error qdate format")
}

func GetQTradeDate(tradeDate int64) string {
	var qDate string
	if tradeDate <= 0 {
		t := time.Now()
		year := t.Year()
		month := int(t.Month())
		qDate = fmt.Sprint(year, "Q", (month+2)/3)
	} else {
		year := tradeDate / 10000
		month := tradeDate/100 - year*100
		qDate = fmt.Sprint(year, "Q", (month+2)/3)
	}

	return qDate
}

func GetQReportDate(reportDate string) string {
	t, _ := time.Parse("2006-01-02 00:00:00", reportDate)
	year := t.Year()
	month := int(t.Month())
	qDate := fmt.Sprint(year, "Q", (month+2)/3)

	return qDate
}

func GetWTradeDate(tradeDate int64) string {
	var t time.Time
	if tradeDate <= 0 {
		t = time.Now()
	} else {
		t, _ = time.Parse("20060102", fmt.Sprint(tradeDate))
	}
	year, week := t.ISOWeek()
	var wDate string
	if week < 10 {
		wDate = fmt.Sprint(year, "W0", week)
	} else {
		wDate = fmt.Sprint(year, "W", week)
	}

	return wDate
}

func ReflectFields(d interface{}) (fieldParam string) {
	fields := []string{}
	tp := reflect.TypeOf(d)
	v := reflect.ValueOf(d)
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if f.Name == "BaseEntity" {
			continue
		}
		jsonTag := f.Tag.Get("json")
		value := v.Field(i).Interface()
		if strings.Contains(jsonTag, "omitempty") {
			if value == "" {
				continue
			}
			if value == 0 {
				continue
			}
		}
		t1 := fmt.Sprintf("%s=%s", removeOmitEmpty(jsonTag), fmt.Sprint(value))
		fields = append(fields, t1)
	}
	fieldParam = strings.Join(fields, "&")
	return
}

func removeOmitEmpty(tag string) string {
	// remove omitEmpty
	if strings.HasSuffix(tag, "omitempty") {
		idx := strings.Index(tag, ",")
		if idx > 0 {
			tag = tag[:idx]
		} else {
			tag = ""
		}
	}
	return tag
}

func ToCsv(head []string, data []interface{}) string {
	raw := strings.Join(head, ",") + "\n"
	i := 0
	for _, d := range data {
		jsonMap, _, _ := GetJsonMap(d)
		line := ""
		lineValue := collection.StructToMap(d, nil)
		j := 0
		for _, h := range head {
			f := jsonMap[h]
			if f != "" {
				value := lineValue[f]
				if value != nil {
					line = line + fmt.Sprint(value)
				}
			}
			if j < len(head)-1 {
				line = line + ","
			}
			j++
		}
		if i < len(data)-1 {
			line = line + "\n"
		}
		raw = raw + line
	}

	return raw
}

func GetJsonMap(data interface{}) (map[string]string, map[string]string, []string) {
	jsonMap := make(map[string]string)
	typMap := make(map[string]string)
	jsonHeads := make([]string, 0)
	tp := reflect.TypeOf(data)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	v := reflect.ValueOf(data)
	v = reflect.Indirect(v)
	for j := 0; j < tp.NumField(); j++ {
		f := tp.Field(j)
		if f.Type.Kind() == reflect.Struct {
			value := v.Field(j).Interface()
			ns, ts, hs := GetJsonMap(value)
			if ns != nil && len(ns) > 0 {
				jsonHeads = append(jsonHeads, hs...)
				for n, fname := range ns {
					jsonMap[n] = fname
				}
			}
			if ts != nil && len(ts) > 0 {
				for t, tname := range ts {
					typMap[t] = tname
				}
			}
		} else {
			jsonTag := f.Tag.Get("json")
			jsonTag = removeOmitEmpty(jsonTag)
			jsonMap[jsonTag] = f.Name
			jsonHeads = append(jsonHeads, jsonTag)
			typMap[jsonTag] = f.Type.Name()
		}
	}

	return jsonMap, typMap, jsonHeads
}

func InBuildStr(colname string, paraStr string, sep string) (string, []interface{}) {
	paraStr = strings.Trim(paraStr, sep)
	paras := make([]interface{}, 0)
	if paraStr == "" {
		return "1=1", paras
	}
	conds := colname + " in ("
	eles := strings.Split(paraStr, sep)
	i := 0
	for _, ele := range eles {
		if i == 0 {
			conds = conds + "?"
		} else {

			conds = conds + ",?"
		}
		paras = append(paras, ele)
		i++
	}
	conds = conds + ")"

	return conds, paras
}

func InBuildInt(colname string, paraInt []int) (string, []interface{}) {
	paras := make([]interface{}, 0)
	if paraInt == nil && len(paraInt) == 0 {
		return "1=1", paras
	}
	conds := colname + " in ("
	eles := paraInt
	i := 0
	for _, ele := range eles {
		if i == 0 {
			conds = conds + "?"
		} else {

			conds = conds + ",?"
		}
		paras = append(paras, ele)
		i++
	}
	conds = conds + ")"

	return conds, paras
}

func CalTpr(apr float64, period float64) float64 {
	return math.Pow(1+apr, period)
}

func CalApr(tpr float64, period float64) float64 {
	apr := float64(1)
	tmp := float64(0)
	for math.Abs(tmp-apr) > 0.0000000001 {
		tmp = apr
		apr = apr - (math.Pow(apr, period)-tpr)/(period*math.Pow(apr, period-1))
	}
	if math.IsNaN(apr) {
		return -99
	}
	return apr - 1
}

func Equal(src float64, target float64) bool {
	if src > target {
		return src-target < 0.0000000001
	} else {
		return target-src < 0.0000000001
	}
}

func CreateI18n(data interface{}) string {
	jsonMap, _, jsonHeads := GetJsonMap(data)
	i18n := ""
	for _, head := range jsonHeads {
		i18n += "\"" + head + "\": " + "\"" + jsonMap[head] + "\",\n"
	}

	return i18n
}

func init() {
}
