package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

type ProcessLogController struct {
	controller.BaseController
}

var processLogController *ProcessLogController

func (this *ProcessLogController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.ProcessLog, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *ProcessLogController) Schedule(ctx iris.Context) {
	svc := this.BaseService.(*service.ProcessLogService)
	svc.Schedule()
}

func init() {
	processLogController = &ProcessLogController{
		BaseController: controller.BaseController{
			BaseService: service.GetProcessLogService(),
		},
	}
	processLogController.BaseController.ParseJSON = processLogController.ParseJSON
	container.RegistController("processlog", processLogController)
}
