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

// WmqyLineController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，
// 因此每个实体（外部rest方式访问）的控制层都需要写一遍
type WmqyLineController struct {
	controller.BaseController
}

var wmqyLineController *WmqyLineController

func (ctl *WmqyLineController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.WmqyLine, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *WmqyLineController) RefreshWmqyLine(ctx iris.Context) {
	svc := ctl.BaseService.(*service.WmqyLineService)
	err := svc.RefreshWmqyLine(-1)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *WmqyLineController) GetUpdateWmqyLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.WmqyLineService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err = ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps, err := svc.GetUpdateWmqyLine(tsCode, -1, 10000, nil)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}

	err = ctx.JSON(ps)
	if err != nil {
		return
	}
}

func (ctl *WmqyLineController) StdPath(ctx iris.Context) {
	svc := ctl.BaseService.(*service.WmqyLineService)
	err := svc.StdPath("C:\\stock\\data\\minmax\\qline", "C:\\stock\\data\\standard\\qline", 19900101, 20211231)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *WmqyLineController) FindQPerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := service.GetQPerformanceService()
	var tsCode string
	var startDate string
	var endDate string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err = ctx.JSON(ps)
		if err != nil {
			return
		}
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
	wmqyLineMap, err := svc.FindQPerformance(service.LinetypeWmqy, tsCode, startDate, endDate)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc.Compute(wmqyLineMap, nil)
	qps := svc.StdMap(wmqyLineMap, service.StdtypeMinmax, isWinsorize)
	err = ctx.JSON(qps)
	if err != nil {
		return
	}
}

func (ctl *WmqyLineController) FindQExpress(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := service.GetQPerformanceService()
	var tsCode string
	var startDate string
	var endDate string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err = ctx.JSON(ps)
		if err != nil {
			return
		}
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
	wmqyLineMap, err := svc.FindQExpress(service.LinetypeWmqy, tsCode, startDate, endDate)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc.Compute(wmqyLineMap, nil)
	var isWinsorize bool
	v, ok = params["isWinsorize"]
	if ok {
		isWinsorize = v.(bool)
	}
	qps := svc.StdMap(wmqyLineMap, service.StdtypeMinmax, isWinsorize)
	err = ctx.JSON(qps)
	if err != nil {
		return
	}
}

func (ctl *WmqyLineController) FindQForecast(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := service.GetQPerformanceService()
	var tsCode string
	var startDate string
	var endDate string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err = ctx.JSON(ps)
		if err != nil {
			return
		}
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
	wmqyLineMap, err := svc.FindQForecast(service.LinetypeWmqy, tsCode, startDate, endDate)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc.Compute(wmqyLineMap, nil)
	var isWinsorize bool
	v, ok = params["isWinsorize"]
	if ok {
		isWinsorize = v.(bool)
	}
	qps := svc.StdMap(wmqyLineMap, service.StdtypeMinmax, isWinsorize)
	err = ctx.JSON(qps)
	if err != nil {
		return
	}
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

func (ctl *WmqyLineController) FindPreceding(ctx iris.Context) {
	wmqylinePara := &WmqyLinePara{}
	err := ctx.ReadJSON(&wmqylinePara)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if wmqylinePara.TsCode == "" {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	if wmqylinePara.LineType == 0 {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("LineType is 0"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.WmqyLineService)
	ps, count, err := svc.FindPreceding(wmqylinePara.TsCode, wmqylinePara.LineType, wmqylinePara.EndDate, wmqylinePara.From, wmqylinePara.Limit, wmqylinePara.Count)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *WmqyLineController) FindFollowing(ctx iris.Context) {
	wmqylinePara := &WmqyLinePara{}
	err := ctx.ReadJSON(&wmqylinePara)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if wmqylinePara.TsCode == "" {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	if wmqylinePara.LineType == 0 {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("LineType is 0"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.WmqyLineService)
	ps, count, err := svc.FindFollowing(wmqylinePara.TsCode, wmqylinePara.LineType, wmqylinePara.StartDate, wmqylinePara.EndDate, wmqylinePara.From, wmqylinePara.Limit, wmqylinePara.Count)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

/*
*
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
