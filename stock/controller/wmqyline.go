package controller

import (
	"errors"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

/**
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type WmqyLineController struct {
	controller.BaseController
}

var wmqyLineController *WmqyLineController

func (this *WmqyLineController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.WmqyLine, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *WmqyLineController) RefreshWmqyLine(ctx iris.Context) {
	svc := this.BaseService.(*service.WmqyLineService)
	err := svc.RefreshWmqyLine(-1)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *WmqyLineController) GetUpdateWmqyLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.WmqyLineService)
	var ts_code string
	v, ok := params["ts_code"]
	if ok {
		ts_code = v.(string)
	}
	if ts_code == "" {
		ps := make([]interface{}, 0)
		ctx.JSON(ps)
		return
	}
	ps, err := svc.GetUpdateWmqyLine(ts_code, -1, 10000, nil)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *WmqyLineController) StdPath(ctx iris.Context) {
	svc := this.BaseService.(*service.WmqyLineService)
	err := svc.StdPath("C:\\stock\\data\\minmax\\qline", "C:\\stock\\data\\standard\\qline", 19900101, 20211231)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *WmqyLineController) FindQPerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := service.GetQPerformanceService()
	var ts_code string
	var startDate string
	var endDate string
	v, ok := params["ts_code"]
	if ok {
		ts_code = v.(string)
	}
	if ts_code == "" {
		ps := make([]interface{}, 0)
		ctx.JSON(ps)
		return
	}
	v, ok = params["startDate"]
	if ok {
		startDate = v.(string)
	}
	v, ok = params["endDate"]
	if ok {
		endDate = v.(string)
	}
	var term int
	v, ok = params["term"]
	if ok {
		f, ok := v.(float64)
		if ok && f > 0 && startDate == "" {
			term = int(f)
			today := stock.GetQTradeDate(0)
			startDate, _ = stock.AddYear(today, -term)
		}
	}
	var isWinsorize bool
	v, ok = params["isWinsorize"]
	if ok {
		isWinsorize = v.(bool)
	}
	wmqyLineMap, err := svc.FindQPerformance(service.LineType_Wmqy, ts_code, startDate, endDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc.Compute(wmqyLineMap, nil)
	qps := svc.StdMap(wmqyLineMap, service.StdType_MinMax, isWinsorize)
	ctx.JSON(qps)
}

func (this *WmqyLineController) FindQExpress(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := service.GetQPerformanceService()
	var ts_code string
	var startDate string
	var endDate string
	v, ok := params["ts_code"]
	if ok {
		ts_code = v.(string)
	}
	if ts_code == "" {
		ps := make([]interface{}, 0)
		ctx.JSON(ps)
		return
	}
	v, ok = params["startDate"]
	if ok {
		startDate = v.(string)
	}
	v, ok = params["endDate"]
	if ok {
		endDate = v.(string)
	}
	var term int
	v, ok = params["term"]
	if ok {
		f, ok := v.(float64)
		if ok && f > 0 && startDate == "" {
			term = int(f)
			today := stock.GetQTradeDate(0)
			startDate, _ = stock.AddYear(today, -term)
		}
	}
	wmqyLineMap, err := svc.FindQExpress(service.LineType_Wmqy, ts_code, startDate, endDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc.Compute(wmqyLineMap, nil)
	var isWinsorize bool
	v, ok = params["isWinsorize"]
	if ok {
		isWinsorize = v.(bool)
	}
	qps := svc.StdMap(wmqyLineMap, service.StdType_MinMax, isWinsorize)
	ctx.JSON(qps)
}

func (this *WmqyLineController) FindQForecast(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := service.GetQPerformanceService()
	var ts_code string
	var startDate string
	var endDate string
	v, ok := params["ts_code"]
	if ok {
		ts_code = v.(string)
	}
	if ts_code == "" {
		ps := make([]interface{}, 0)
		ctx.JSON(ps)
		return
	}
	v, ok = params["startDate"]
	if ok {
		startDate = v.(string)
	}
	v, ok = params["endDate"]
	if ok {
		endDate = v.(string)
	}
	var term int
	v, ok = params["term"]
	if ok {
		f, ok := v.(float64)
		if ok && f > 0 && startDate == "" {
			term = int(f)
			today := stock.GetQTradeDate(0)
			startDate, _ = stock.AddYear(today, -term)
		}
	}
	wmqyLineMap, err := svc.FindQForecast(service.LineType_Wmqy, ts_code, startDate, endDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc.Compute(wmqyLineMap, nil)
	var isWinsorize bool
	v, ok = params["isWinsorize"]
	if ok {
		isWinsorize = v.(bool)
	}
	qps := svc.StdMap(wmqyLineMap, service.StdType_MinMax, isWinsorize)
	ctx.JSON(qps)
}

type WmqyLinePara struct {
	From      int    `json:"from,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Orderby   string `json:"orderby,omitempty"`
	Count     int64  `json:"count,omitempty"`
	TsCode    string `json:"ts_code,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	TradeDate int64  `json:"trade_date,omitempty"`
	LineType  int    `json:"line_type,omitempty"`
}

func (this *WmqyLineController) FindPreceding(ctx iris.Context) {
	wmqylinePara := &WmqyLinePara{}
	err := ctx.ReadJSON(&wmqylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if wmqylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	if wmqylinePara.LineType == 0 {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("LineType is 0"))

		return
	}
	svc := this.BaseService.(*service.WmqyLineService)
	ps, count, err := svc.FindPreceding(wmqylinePara.TsCode, wmqylinePara.LineType, wmqylinePara.EndDate, wmqylinePara.From, wmqylinePara.Limit, wmqylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *WmqyLineController) FindFollowing(ctx iris.Context) {
	wmqylinePara := &WmqyLinePara{}
	err := ctx.ReadJSON(&wmqylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if wmqylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	if wmqylinePara.LineType == 0 {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("LineType is 0"))

		return
	}
	svc := this.BaseService.(*service.WmqyLineService)
	ps, count, err := svc.FindFollowing(wmqylinePara.TsCode, wmqylinePara.LineType, wmqylinePara.StartDate, wmqylinePara.EndDate, wmqylinePara.From, wmqylinePara.Limit, wmqylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

/**
注册bean管理器，注册序列
*/
func init() {
	wmqyLineController = &WmqyLineController{
		BaseController: controller.BaseController{
			BaseService: service.GetWmqyLineService(),
		},
	}
	wmqyLineController.BaseController.ParseJSON = wmqyLineController.ParseJSON
	container.RegistController("wmqyline", wmqyLineController)
}
