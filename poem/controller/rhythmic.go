package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/poem/entity"
	"github.com/curltech/go-colla-stock/poem/service"
	"github.com/curltech/go-colla-web/controller"
)

/**
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type RhythmicController struct {
	controller.BaseController
}

var rhythmicController *RhythmicController

func (this *RhythmicController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Rhythmic, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

/**
注册bean管理器，注册序列
*/
func init() {
	rhythmicController = &RhythmicController{
		BaseController: controller.BaseController{
			BaseService: service.GetRhythmicService(),
		},
	}
	rhythmicController.BaseController.ParseJSON = rhythmicController.ParseJSON
	container.RegistController("rhythmic", rhythmicController)
}
