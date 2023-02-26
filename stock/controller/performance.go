package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

/**
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type PerformanceController struct {
	controller.BaseController
}

var performanceController *PerformanceController

func (this *PerformanceController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Performance, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *PerformanceController) RefreshPerformance(ctx iris.Context) {
	svc := this.BaseService.(*service.PerformanceService)
	err := svc.RefreshPerformance()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *PerformanceController) GetUpdatePerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.PerformanceService)
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
	ps, err := svc.GetUpdatePerformance(ts_code)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
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

func (this *PerformanceController) FindLatest(ctx iris.Context) {
	param := &PerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.PerformanceService)
	es, count, err := svc.FindLatest(param.SecurityCode, param.StartDate, param.Orderby, param.From, param.Limit, param.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{}, 0)
	result["data"] = es
	result["count"] = count
	ctx.JSON(result)
}

func (this *PerformanceController) FindByQDate(ctx iris.Context) {
	param := &PerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if param.SecurityCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.PerformanceService)
	es, err := svc.FindByQDate(param.SecurityCode, param.StartDate, param.EndDate, param.Orderby)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(es)
}

func (this *PerformanceController) Search(ctx iris.Context) {
	param := &PerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.PerformanceService)
	es, count, err := svc.Search(param.SecurityCode, param.StartDate, param.EndDate, param.Orderby, param.From, param.Limit, param.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{}, 0)
	result["data"] = es
	result["count"] = count

	ctx.JSON(result)
}

/**
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
