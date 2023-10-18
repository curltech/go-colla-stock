package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// MinLineController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type MinLineController struct {
	controller.BaseController
}

var minLineController *MinLineController

func (ctl *MinLineController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.MinLine, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *MinLineController) ParsePath(ctx iris.Context) {
	svc := ctl.BaseService.(*service.MinLineService)
	err := svc.ParsePath("C:\\zd_zsone\\vipdoc\\sz\\minline", "C:\\stock\\data\\origin\\minline")
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *MinLineController) RefreshMinLine(ctx iris.Context) {
	svc := ctl.BaseService.(*service.MinLineService)
	err := svc.RefreshMinLine(-1)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *MinLineController) GetUpdateMinLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.MinLineService)
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
	ps, err := svc.GetUpdateMinLine(tsCode, -1, 10000)
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

func (ctl *MinLineController) RefreshTodayMinLine(ctx iris.Context) {
	svc := ctl.BaseService.(*service.MinLineService)
	err := svc.RefreshTodayMinLine()
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *MinLineController) GetUpdateTodayMinLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.MinLineService)
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
	ps, err := svc.GetUpdateTodayMinLine(tsCode)
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

func (ctl *MinLineController) FindMinLines(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.MinLineService)
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
	var tradeDate int64
	v, ok = params["trade_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			tradeDate = int64(f)
		}
	}
	var tradeMinute int64
	v, ok = params["trade_minute"]
	if ok {
		f, ok := v.(float64)
		if ok {
			tradeMinute = int64(f)
		}
	}
	ps, err := svc.FindMinLines(tsCode, tradeDate, tradeMinute)
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

// 注册bean管理器，注册序列
func init() {
	minLineController = &MinLineController{
		BaseController: controller.BaseController{
			BaseService: service.GetMinLineService(),
		},
	}
	minLineController.BaseController.ParseJSON = minLineController.ParseJSON
	container.RegistController("minline", minLineController)
}
