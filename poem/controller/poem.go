package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/poem/entity"
	"github.com/curltech/go-colla-stock/poem/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

/*
*
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type PoemController struct {
	controller.BaseController
}

var poemController *PoemController

func (this *PoemController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Poem, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *PoemController) ParsePath(ctx iris.Context) {
	service := this.BaseService.(*service.PoemService)
	err := service.ParsePath("C:\\go-workspace\\Poetry-master")
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

/*
*
注册bean管理器，注册序列
*/
func init() {
	poemController = &PoemController{
		BaseController: controller.BaseController{
			BaseService: service.GetPoemService(),
		},
	}
	poemController.BaseController.ParseJSON = poemController.ParseJSON
	container.RegistController("poem", poemController)
}
