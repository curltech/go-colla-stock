package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// PerformanceController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type PerformanceController struct {
	controller.BaseController
}

var performanceController *PerformanceController

func (ctl *PerformanceController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Performance, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *PerformanceController) RefreshPerformance(ctx iris.Context) {
	svc := ctl.BaseService.(*service.PerformanceService)
	err := svc.RefreshPerformance()
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *PerformanceController) GetUpdatePerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.PerformanceService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err := ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps, err := svc.GetUpdatePerformance(tsCode)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}

	err = ctx.JSON(ps)
	if err != nil {
		return
	}
}

type PerformancePara struct {
	SecurityCode string `json:"security_code,omitempty"`
	StartDate    string `json:"start_date,omitempty"`
	EndDate      string `json:"end_date,omitempty"`
	Orderby      string `json:"orderby,omitempty"`
	From         int    `json:"from"`
	Limit        int    `json:"limit"`
	Count        int64  `json:"count"`
}

func (ctl *PerformanceController) FindLatest(ctx iris.Context) {
	param := &PerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.PerformanceService)
	es, count, err := svc.FindLatest(param.SecurityCode, param.StartDate, param.Orderby, param.From, param.Limit, param.Count)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{}, 0)
	result["data"] = es
	result["count"] = count
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *PerformanceController) FindByQDate(ctx iris.Context) {
	param := &PerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.PerformanceService)
	es, count, err := svc.FindByQDate(param.SecurityCode, param.StartDate, param.EndDate, param.Orderby, param.From, param.Limit, param.Count)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{}, 0)
	result["data"] = es
	result["count"] = count
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
	performanceController = &PerformanceController{
		BaseController: controller.BaseController{
			BaseService: service.GetPerformanceService(),
		},
	}
	performanceController.BaseController.ParseJSON = performanceController.ParseJSON
	container.RegistController("performance", performanceController)
}
