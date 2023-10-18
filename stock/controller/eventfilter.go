package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// EventFilterController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type EventFilterController struct {
	controller.BaseController
}

var eventFilterController *EventFilterController

func (ctl *EventFilterController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.EventFilter, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *EventFilterController) RefreshCacheEventFilter(ctx iris.Context) {
	svc := ctl.BaseService.(*service.EventFilterService)
	svc.RefreshCacheEventFilter()
}

// 注册bean管理器，注册序列
func init() {
	eventFilterController = &EventFilterController{
		BaseController: controller.BaseController{
			BaseService: service.GetEventFilterService(),
		},
	}
	eventFilterController.BaseController.ParseJSON = eventFilterController.ParseJSON
	container.RegistController("eventfilter", eventFilterController)
}
