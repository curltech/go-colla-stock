package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
)

// FilterCondController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type FilterCondController struct {
	controller.BaseController
}

var filterCondController *FilterCondController

func (ctl *FilterCondController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.FilterCond, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

// 注册bean管理器，注册序列
func init() {
	filterCondController = &FilterCondController{
		BaseController: controller.BaseController{
			BaseService: service.GetFilterCondService(),
		},
	}
	filterCondController.BaseController.ParseJSON = filterCondController.ParseJSON
	container.RegistController("filtercond", filterCondController)
}
