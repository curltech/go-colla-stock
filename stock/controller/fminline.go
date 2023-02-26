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
type FminLineController struct {
	controller.BaseController
}

var fminLineController *FminLineController

func (this *FminLineController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.FminLine, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *FminLineController) ParsePath(ctx iris.Context) {
	svc := this.BaseService.(*service.FminLineService)
	err := svc.ParsePath("C:\\zd_zsone\\vipdoc\\sz\\fzline", "C:\\stock\\data\\origin\\fzline")
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

/**
注册bean管理器，注册序列
*/
func init() {
	fminLineController = &FminLineController{
		BaseController: controller.BaseController{
			BaseService: service.GetFminLineService(),
		},
	}
	fminLineController.BaseController.ParseJSON = fminLineController.ParseJSON
	container.RegistController("fminline", fminLineController)
}
