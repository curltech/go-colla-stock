package controller

import (
	"errors"
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
type ShareController struct {
	controller.BaseController
}

var shareController *ShareController

func (this *ShareController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Share, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *ShareController) GetMine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode, _ = v.(string)
	}
	if tsCode == "" {
		userShares := make([]interface{}, 0)
		ctx.JSON(userShares)
		return
	}
	svc := this.BaseService.(*service.ShareService)
	userShares, err := svc.GetShares(tsCode)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	ctx.JSON(userShares)

	return
}

func (this *ShareController) Search(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	var keyword string
	v, ok := params["keyword"]
	if ok && v != nil {
		keyword = v.(string)
	}
	svc := this.BaseService.(*service.ShareService)
	shares, err := svc.Search(keyword, 0, 0)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	ctx.JSON(shares)

	return
}

func (this *ShareController) UpdateSector(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	var tscode string
	v, ok := params["ts_code"]
	if ok && v != nil {
		tscode = v.(string)
	}
	var sector string
	v, ok = params["sector"]
	if ok && v != nil {
		sector = v.(string)
	}
	if tscode == "" || sector == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("NoTsCodeOrSector"))

		return
	}
	svc := this.BaseService.(*service.ShareService)
	svc.UpdateSector(tscode, sector)

	return
}

func (this *ShareController) UpdateShares(ctx iris.Context) {
	svc := this.BaseService.(*service.ShareService)
	svc.UpdateShares()

	return
}

/**
注册bean管理器，注册序列
*/
func init() {
	shareController = &ShareController{
		BaseController: controller.BaseController{
			BaseService: service.GetShareService(),
		},
	}
	shareController.BaseController.ParseJSON = shareController.ParseJSON
	container.RegistController("share", shareController)
}
