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
type EventController struct {
	controller.BaseController
}

var eventController *EventController

func (this *EventController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Event, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *EventController) RefreshCacheEvent(ctx iris.Context) {
	svc := this.BaseService.(*service.EventService)
	svc.RefreshCacheEvent()
}

func (this *EventController) FindCacheEvent(ctx iris.Context) {
	svc := this.BaseService.(*service.EventService)
	eventMap := svc.GetCacheEvent()
	if eventMap == nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, "event nil")
		return
	}
	ctx.JSON(eventMap)
}

/**
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
