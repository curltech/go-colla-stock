package controller

import (
	"errors"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	rbac "github.com/curltech/go-colla-web/rbac/controller"
	"github.com/kataras/iris/v12"
)

/**
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type ShareEventController struct {
	controller.BaseController
}

var shareEventController *ShareEventController

func (this *ShareEventController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.ShareEvent, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *ShareEventController) GetMine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	userName := rbac.GetUserController().GetCurrentUserName(ctx)
	if userName != "" {
		svc := this.BaseService.(*service.ShareEventService)
		userShares, err := svc.GetMine(userName, tsCode)
		if err != nil {
			ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

			return
		}
		ctx.JSON(userShares)
	} else {
		ctx.StopWithJSON(iris.StatusUnauthorized, errors.New("UserName is not exist"))

		return
	}

	return
}

/**
注册bean管理器，注册序列
*/
func init() {
	shareEventController = &ShareEventController{
		BaseController: controller.BaseController{
			BaseService: service.GetShareEventService(),
		},
	}
	shareEventController.BaseController.ParseJSON = shareEventController.ParseJSON
	container.RegistController("shareevent", shareEventController)
}
