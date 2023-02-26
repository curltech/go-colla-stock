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
type ExpressController struct {
	controller.BaseController
}

var expressController *ExpressController

func (this *ExpressController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Express, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *ExpressController) RefreshExpress(ctx iris.Context) {
	svc := this.BaseService.(*service.ExpressService)
	err := svc.RefreshExpress()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *ExpressController) GetUpdateExpress(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.ExpressService)
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
	ps, err := svc.GetUpdateExpress(ts_code)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *ExpressController) FindLatest(ctx iris.Context) {
	param := &PerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.ExpressService)
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

func (this *ExpressController) FindByQDate(ctx iris.Context) {
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
	svc := this.BaseService.(*service.ExpressService)
	es, err := svc.FindByQDate(param.SecurityCode, param.StartDate, param.EndDate, param.Orderby)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(es)
}

func (this *ExpressController) Search(ctx iris.Context) {
	param := &PerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.ExpressService)
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
	expressController = &ExpressController{
		BaseController: controller.BaseController{
			BaseService: service.GetExpressService(),
		},
	}
	expressController.BaseController.ParseJSON = expressController.ParseJSON
	container.RegistController("express", expressController)
}
