package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// EventController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type EventController struct {
	controller.BaseController
}

var eventController *EventController

func (ctl *EventController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Event, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *EventController) RefreshCacheEvent(ctx iris.Context) {
	svc := ctl.BaseService.(*service.EventService)
	svc.RefreshCacheEvent()
}

func (ctl *EventController) FindCacheEvent(ctx iris.Context) {
	svc := ctl.BaseService.(*service.EventService)
	eventMap := svc.GetCacheEvent()
	if eventMap == nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, "event nil")
		if err != nil {
			return
		}
		return
	}
	err := ctx.JSON(eventMap)
	if err != nil {
		return
	}
}

/*
*
注册bean管理器，注册序列
*/
func init() {
	eventController = &EventController{
		BaseController: controller.BaseController{
			BaseService: service.GetEventService(),
		},
	}
	eventController.BaseController.ParseJSON = eventController.ParseJSON
	container.RegistController("event", eventController)
}
