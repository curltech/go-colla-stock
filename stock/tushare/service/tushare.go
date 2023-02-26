package service

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/valyala/fasthttp"
	"io"
	"mime"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

type TushareRequest struct {
	ApiName string      `json:"api_name,omitempty"`
	Token   string      `json:"token,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Fields  string      `json:"fields,omitempty"`
}

/**
code： 接口返回码，2002表示权限问题。

msg：错误信息，比如“系统内部错误”，“没有权限”等

data：数据，data里包含fields和items字段，分别为字段和数据内容
*/
type TushareData struct {
	Fields []string        `json:"fields,omitempty"`
	Items  [][]interface{} `json:"items,omitempty"`
}

type TushareResponse struct {
	Code int          `json:"code,omitempty"`
	Msg  string       `json:"msg,omitempty"`
	Data *TushareData `json:"data,omitempty"`
}

var url = "https://api.waditu.com/"
var token = "6723fd77a56192845e75cc1366e434c937d54300801d7081315e9b1c"

/**
api_name：接口名称，比如stock_basic

token ：用户唯一标识，可通过登录pro网站获取

params：接口参数，如daily接口中start_date和end_date

fields：字段列表，用于接口获取指定的字段，以逗号分隔，如"open,high,low,close"
*/

func Post(tushareRequest *TushareRequest) (tsRsp *TushareResponse, err error) {
	var req *http.Request
	var resp *http.Response
	var bodyJSON []byte
	bodyJSON, err = json.Marshal(tushareRequest)
	if err != nil {
		return
	}

	// Build send data
	senddata := io.NopCloser(bytes.NewReader(bodyJSON))
	req, err = http.NewRequest("POST", url, senddata)
	if err != nil {
		return
	}

	// Set http content type
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// Check mime type of response
	var mimeType string
	mimeType, _, err = mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return
	}
	if mimeType != "application/json" {
		err = fmt.Errorf("Could not execute request (%s)", fmt.Sprintf("Response Content-Type is '%s', but should be 'application/json'.", mimeType))
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("code:%d return:%s", resp.StatusCode, string(body))
		return
	}
	tsRsp = new(TushareResponse)

	err = json.Unmarshal(body, &tsRsp)
	if err != nil {
		return
	}
	err = tsRsp.CheckValid()
	return
}

func FastPost(tushareRequest *TushareRequest) (*TushareResponse, error) {
	client := fasthttp.Client{TLSConfig: &tls.Config{InsecureSkipVerify: true}}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	bodyjson, err := json.TextMarshal(tushareRequest)
	if err != nil {
		return nil, err
	}
	req.SetBodyString(bodyjson) //设置请求参数
	//req.SetBody([]byte(`{"body": "sjon_str"}`)) //设置[]byte

	resp := fasthttp.AcquireResponse()
	if err := client.Do(req, resp); err != nil {
		logger.Sugar.Errorf("loan list fail to do request. appID=%s. [err=%v]\n", req.Header, err)
		return nil, err
	}
	b := resp.Body()
	if resp.StatusCode() != fasthttp.StatusOK {
		logger.Sugar.Errorf("loan list failed code=%d. [err=%v]\n", resp.StatusCode(), string(b))
		return nil, err
	}
	r := &TushareResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func reflectFields(i interface{}) (fieldParam string) {
	fields := []string{}
	tp := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if f.Name == "BaseEntity" || f.Name == "PinYin" || f.Name == "Sector" {
			continue
		}
		v.Field(i).Interface()
		fields = append(fields, removeOmitEmpty(f.Tag.Get("json")))
	}
	fieldParam = strings.Join(fields, ",")
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

func (ts *TushareResponse) CheckValid() error {
	switch ts.Code {
	case -2001:
		logger.Sugar.Errorf("[TushareAPI Response] code:%d,argument error:%s", ts.Code, ts.Msg)
	case -2002:
		logger.Sugar.Errorf("[TushareAPI Response] code:%d,privilege error:%s", ts.Code, ts.Msg)
	case 0:
		if len(ts.Data.Items) == 0 {
			logger.Sugar.Errorf("[TushareAPI Response] code:%d,msg:%s but empty data", ts.Code, ts.Msg)
		}
		return nil
	default:
		logger.Sugar.Errorf("[TushareAPI Response] code:%d,msg:%s", ts.Code, ts.Msg)
	}

	return errors.New(fmt.Sprint(ts.Code))
}

func (ts *TushareRequest) CheckValid() error {
	if ts.Token == "" {
		return errors.New("[TushareAPI Request] must set user token")
	}
	if ts.ApiName == "" {
		return errors.New("[TushareAPI Request] must set api_name")
	}
	return nil
}

// 重组数据为[]byte以方便映射struct
func ReflectResponseData(fields []string, data []interface{}) (body []byte, err error) {
	m := make(map[string]interface{})
	if len(fields) != len(data) {
		err = fmt.Errorf("fields(len %d) not fit on data(len %d)", len(fields), len(data))
		return
	}
	if len(fields) == 0 {
		err = errors.New("empty data and fields")
		return
	}
	for n, f := range fields {
		m[f] = data[n]
	}
	body, err = json.Marshal(m)
	return
}

func init() {
	os.Setenv("LANG", "zh_CN.UTF8")
}
