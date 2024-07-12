package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/poem/entity"
	"github.com/curltech/go-colla-stock/poem/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// PoemController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type PoemController struct {
	controller.BaseController
}

var poemController *PoemController

func (ctl *PoemController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Poem, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type PoemPara struct {
	entity.Poem
	From  int `json:"from,omitempty"`
	Limit int `json:"limit,omitempty"`
}

func (ctl *PoemController) Search(ctx iris.Context) {
	param := &PoemPara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.PoemService)
	poems, err := svc.Search(param.Title, param.Author, param.Rhythmic, param.Dynasty, param.Paragraphs, param.From, param.Limit)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	err = ctx.JSON(poems)
	if err != nil {
		return
	}

	return
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
