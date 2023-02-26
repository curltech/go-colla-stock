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
type MinLineController struct {
	controller.BaseController
}

var minLineController *MinLineController

func (this *MinLineController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.MinLine, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *MinLineController) ParsePath(ctx iris.Context) {
	svc := this.BaseService.(*service.MinLineService)
	err := svc.ParsePath("C:\\zd_zsone\\vipdoc\\sz\\minline", "C:\\stock\\data\\origin\\minline")
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *MinLineController) RefreshMinLine(ctx iris.Context) {
	svc := this.BaseService.(*service.MinLineService)
	err := svc.RefreshMinLine(-1)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *MinLineController) GetUpdateMinLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.MinLineService)
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
	ps, err := svc.GetUpdateMinLine(ts_code, -1, 10000)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *MinLineController) RefreshTodayMinLine(ctx iris.Context) {
	svc := this.BaseService.(*service.MinLineService)
	err := svc.RefreshTodayMinLine()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *MinLineController) GetUpdateTodayMinLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.MinLineService)
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
	ps, err := svc.GetUpdateTodayMinLine(ts_code)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *MinLineController) FindMinLines(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.MinLineService)
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
	var trade_date int64
	v, ok = params["trade_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			trade_date = int64(f)
		}
	}
	var trade_minute int64
	v, ok = params["trade_minute"]
	if ok {
		f, ok := v.(float64)
		if ok {
			trade_minute = int64(f)
		}
	}
	ps, err := svc.FindMinLines(ts_code, trade_date, trade_minute)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

/**
注册bean管理器，注册序列
*/
func init() {
	minLineController = &MinLineController{
		BaseController: controller.BaseController{
			BaseService: service.GetMinLineService(),
		},
	}
	minLineController.BaseController.ParseJSON = minLineController.ParseJSON
	container.RegistController("minline", minLineController)
}
