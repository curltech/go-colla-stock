package eastmoney

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"strings"
	"time"
)

type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

func (t Time) Now() Time {
	return Time(time.Now())
}

func (t Time) ParseTime(tt time.Time) Time {
	return Time(tt)
}

func (t Time) format() string {
	return time.Time(t).Format(timeFormat)
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.format()), nil
}

func (t *Time) FromDB(b []byte) error {
	if nil == b || len(b) == 0 {
		t = nil
		return nil
	}
	var now time.Time
	var err error
	now, err = time.ParseInLocation(timeFormat, string(b), time.Local)
	if nil == err {
		*t = Time(now)
		return nil
	}
	now, err = time.ParseInLocation("2006-01-02T15:04:05Z", string(b), time.Local)
	if nil == err {
		*t = Time(now)
		return nil
	}
	panic("自己定义个layout日期格式处理一下数据库里面的日期型数据解析!")
	return err
}

func (t *Time) ToDB() ([]byte, error) {
	if nil == t {
		return nil, nil
	}
	return []byte(time.Time(*t).Format(timeFormat)), nil
}

func (t *Time) Value() (driver.Value, error) {
	if nil == t {
		return nil, nil
	}
	return time.Time(*t).Format(timeFormat), nil
}

type ReportRequestParam struct {
	Callback string `json:"callback"`
	St       string `json:"st"`   //排序的字段，逗号分隔
	Sr       string `json:"sr"`   //-1降序，1升序，逗号分隔
	Ps       string `json:"ps"`   //每页的记录数
	P        int    `json:"p"`    //页数
	Type     string `json:"type"` //获取的数据类型
	Sty      string `json:"sty"`  //获取的字段，ALL
	Token    string `json:"token"`
	Filter   string `json:"filter,omitempty"` //条件
}

type ReportResponseData struct {
	Pages int `json:"pages,omitempty"`
	Count int `json:"count,omitempty"`
}

type ReportResponseResult struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Version string `json:"version,omitempty"`
}

var report_url = "http://datacenter-web.eastmoney.com/api/data/get"
var report_callback = "jQuery112300570549541785732_1638366842955"
var report_token = "894050c76af8597a853f5b408b759f5d"
var performance_type = "RPT_LICO_FN_CPD"
var report_sty = "ALL"

func CreateRequestParam() *ReportRequestParam {
	return &ReportRequestParam{Callback: report_callback, Token: report_token, Sty: report_sty}
}

/*
*
空格    -    %20 （URL中的空格可以用+号或者编码值表示）
"          -    %22
#         -    %23
%        -    %25
&         -    %26
(          -    %28
)          -    %29
+         -    %2B
,          -    %2C
/          -    %2F
:          -    %3A
;          -    %3B
<         -    %3C
=         -    %3D
>         -    %3E
?         -    %3F
@       -    %40
\          -    %5C
|          -    %7C
{          -    %7B
}          -    %7D
*/
func FastGet(urlStr string, requestParam interface{}) ([]byte, error) {
	args := stock.ReflectFields(requestParam)
	uri := urlStr + "?" + args
	//logger.Sugar.Infof("fasthttp get: %v", uri)
	status, resp, err := fasthttp.Get(nil, uri)
	if err != nil {
		logger.Sugar.Errorf("fasthttp Get fail:", err.Error())
		status, resp, err = fasthttp.Get(nil, uri)
		if err != nil {
			logger.Sugar.Errorf("second fasthttp Get fail:", err.Error())
			return nil, err
		}
	}
	if status != fasthttp.StatusOK {
		logger.Sugar.Errorf("fasthttp Get status:%s", status)
		return nil, errors.New(fmt.Sprint(status))
	}

	return resp, nil
}

func Get(url string, requestParam interface{}) ([]byte, error) {
	args := stock.ReflectFields(requestParam)
	uri := url + "?" + args
	logger.Sugar.Infof("http get: %v", uri)
	resp, err := http.Get(uri)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logger.Sugar.Errorf("http Get fail:%v", err.Error())
		resp, err = http.Get(uri)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			logger.Sugar.Errorf("second http Get fail:%v", err.Error())
			return nil, err
		}
	}
	if resp.StatusCode != http.StatusOK {
		logger.Sugar.Errorf("http Get status:%s", resp.StatusCode)
		return nil, errors.New(fmt.Sprint(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Sugar.Errorf("http Get fail:%v", err.Error())
		return nil, err
	}
	return body, nil
}

func ReportFastGet(requestParam ReportRequestParam) ([]byte, error) {
	resp, err := FastGet(report_url, requestParam)
	if err != nil {
		fmt.Println("fasthttp Get fail:", err.Error())
		return nil, err
	}
	respStr := string(resp)
	respStr = strings.TrimPrefix(respStr, requestParam.Callback)
	respStr = strings.TrimPrefix(respStr, "(")
	respStr = strings.TrimSuffix(respStr, ");")
	resp = []byte(respStr)

	return resp, nil
}

type DayLineRequestParam struct {
	Cb         string `json:"cb"`
	SecId      string `json:"secid"`            //股票代码
	Ut         string `json:"ut"`               //token
	Fields1    string `json:"fields1"`          //
	Fields2    string `json:"fields2"`          //
	Klt        int    `json:"klt"`              //每隔时长获取一次记录，1代表一分，5代表5分钟，101代表每天，102代表每周，103代表每月，104代表每季度，105代表每半年，106代表每年
	Fqt        int    `json:"fqt"`              //
	Smplmt     string `json:"smplmt,omitempty"` //
	Lmt        int    `json:"lmt"`              //获取记录数
	Beg        int    `json:"beg,omitempty"`    //开始日期
	End        int    `json:"end"`              //终止日期 20500101
	Underscore string `json:"_"`                //
	Fields     string `json:"fields,omitempty"` //
	Fltt       int    `json:"fltt,omitempty"`   //
	Invt       int    `json:"invt,omitempty"`   //
}

func CreateDayLineRequestParam() *DayLineRequestParam {
	return &DayLineRequestParam{Cb: dayline_callback, Ut: dayline_token}
}

func CreateFinanceFlowRequestParam() *DayLineRequestParam {
	return &DayLineRequestParam{Cb: day_fflow_callback, Ut: day_fflow_token}
}

/*
*
rt:17:获取某只股股票的过去的每天价格情况（日线）
klines："2021-12-03,17.64,17.65,17.70,17.41,707600,1242375056.00,1.65,0.34,0.06,0.36"
"trade_date,open,close,high,low,vol,amount,nil,pct_chg%,change,turnover%"

rt:21:获取某只股股票的过去的每天资金流东情况（分钟线）
"2021-12-03 15:00,-46770242.0,12655892.0,34114366.0,4332407.0,-51102649.0"

rt:22:获取某只股股票的过去的每天资金流动情况
klines："2021-12-03,-46770255.0,12655888.0,34114368.0,4332400.0,-51102655.0,-3.76,1.02,2.75,0.35,-4.11,17.65,0.34,0.00,0.00"
"trade_date,主力净流入/净额,小单净流入/净额,中单净流入/净额,大单净流入/净额,超大单净流入/净额,主力净流入/净占比%,
小单净流入/净占比%,中单净流入/净占比%,大单净流入/净占比%,超大单净流入/净占比%,close,pct_chg,,"
*/
type DayLineResponseData struct {
	Code      string   `json:"code,omitempty"`
	Market    string   `json:"code,omitempty"`
	Name      string   `json:"code,omitempty"`
	Decimal   int      `json:"code,omitempty"`    //小数位
	Dktotal   int      `json:"dktotal,omitempty"` //总记录数
	PreKPrice float64  `json:"preKPrice,omitempty"`
	Klines    []string `json:"klines,omitempty"` //数据
}

/*
*
rt:4:获取某只股票当天的价格情况
*/
type CurrentResponseData struct {
	F43  float64 `json:"f43"`  //close
	F44  float64 `json:"f44"`  //high
	F45  float64 `json:"f45"`  //low
	F46  float64 `json:"f46"`  //open
	F47  int     `json:"f47"`  //vol
	F48  float64 `json:"f48"`  //amount
	F50  float64 `json:"f50"`  //qrr
	F57  string  `json:"f57"`  //ts_code
	F58  string  `json:"f58"`  //name
	F59  int     `json:"f59"`  //
	F60  float64 `json:"f60"`  //preclose
	F107 int     `json:"f107"` //
	F152 int     `json:"f152"`
	F162 float64 `json:"f162"` //pe
	F168 float64 `json:"f168"` //turnover
	F169 float64 `json:"f169"` //pct_chg
	F170 float64 `json:"f170"` //close_chg
	F171 float64 `json:"f171"` //
	F292 int     `json:"f292"`
}

type DayLineResponseResult struct {
	Rc   int                  `json:"rc,omitempty"`
	Rt   int                  `json:"rt,omitempty"`
	Svr  int                  `json:"svr,omitempty"`
	Lt   int                  `json:"lt,omitempty"`
	Full int                  `json:"full,omitempty"`
	Data *DayLineResponseData `json:"data,omitempty"`
}

/*
*
获取从分钟到年的价格数据，获取过去的价格数据用此链接
*/
var dayline_url = "http://push2his.eastmoney.com/api/qt/stock/kline/get"
var dayline_callback = "jQuery112401201342267983887_1638513559390"
var dayline_token = "fa5fd1943c7b386f172d6893dbfba10b"
var dayline_type = "1638513559443"

func DayLineFastGet(requestParam DayLineRequestParam) ([]byte, error) {
	resp, err := FastGet(dayline_url, requestParam)
	if err != nil {
		fmt.Println("Get fail:", err.Error())
		return nil, err
	}
	respStr := string(resp)
	respStr = strings.TrimPrefix(respStr, requestParam.Cb)
	respStr = strings.TrimPrefix(respStr, "(")
	respStr = strings.TrimSuffix(respStr, ");")
	resp = []byte(respStr)

	return resp, nil
}

/*
*
获取按天的资金流动数据
*/
var day_fflow_url = "http://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get"
var day_fflow_callback = "jQuery1123020753842937846168_1638372426493"
var day_fflow_token = "b2884a393a59ad64002292a3e90d46a5"
var day_fflow_type = "1638372426494"

func FinanceFlowFastGet(requestParam DayLineRequestParam) ([]byte, error) {
	resp, err := FastGet(day_fflow_url, requestParam)
	if err != nil {
		fmt.Println("Get fail:", err.Error())
		return nil, err
	}
	respStr := string(resp)
	respStr = strings.TrimPrefix(respStr, requestParam.Cb)
	respStr = strings.TrimPrefix(respStr, "(")
	respStr = strings.TrimSuffix(respStr, ");")
	resp = []byte(respStr)

	return resp, nil
}

/*
*
获取当天实时的价格数据
*/
var today_url = "http://push2.eastmoney.com/api/qt/stock/trends2/get"
var today_callback = "cb_1638806285136_38811903"
var today_token = "e1e6871893c6386c5ff6967026016627"
var today_type = "1638371480346"

/*
*
获取当天实时的资金流动数据
*/
var today_ff_url = "http://push2.eastmoney.com/api/qt/stock/fflow/kline/get"
var today_ff_callback = "jQuery1123018372844125632604_1638371480345"
var today_ff_token = "b2884a393a59ad64002292a3e90d46a5"
var today_ff_type = "1638371480346"

func TodayFastGet(requestParam DayLineRequestParam) ([]byte, error) {
	resp, err := FastGet(today_url, requestParam)
	if err != nil {
		fmt.Println("Get fail:", err.Error())
		return nil, err
	}
	respStr := string(resp)
	respStr = strings.TrimPrefix(respStr, requestParam.Cb)
	respStr = strings.TrimPrefix(respStr, "(")
	respStr = strings.TrimSuffix(respStr, ");")
	resp = []byte(respStr)

	return resp, nil
}

func TodayFfFastGet(requestParam DayLineRequestParam) ([]byte, error) {
	resp, err := FastGet(today_ff_url, requestParam)
	if err != nil {
		fmt.Println("Get fail:", err.Error())
		return nil, err
	}
	respStr := string(resp)
	respStr = strings.TrimPrefix(respStr, requestParam.Cb)
	respStr = strings.TrimPrefix(respStr, "(")
	respStr = strings.TrimSuffix(respStr, ");")
	resp = []byte(respStr)

	return resp, nil
}
