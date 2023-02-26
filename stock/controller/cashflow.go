package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
)

/**
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type CashFlowController struct {
	controller.BaseController
}

var cashFlowController *CashFlowController

func (this *CashFlowController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.CashFlow, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

/**
注册bean管理器，注册序列
*/
func init() {
	cashFlowController = &CashFlowController{
		BaseController: controller.BaseController{
			BaseService: service.GetCashFlowService(),
		},
	}
	cashFlowController.BaseController.ParseJSON = cashFlowController.ParseJSON
	container.RegistController("cashFlow", cashFlowController)
}
