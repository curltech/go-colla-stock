package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/poem/entity"
	"github.com/curltech/go-colla-stock/poem/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// AuthorController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type AuthorController struct {
	controller.BaseController
}

var authorController *AuthorController

func (ctl *AuthorController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Author, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *AuthorController) ParsePath(ctx iris.Context) {
	authorService := ctl.BaseService.(*service.AuthorService)
	err := authorService.ParseFile("C:\\go-workspace\\chinese-poetry-master\\ci\\author.song.json")
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

/*
*
注册bean管理器，注册序列
*/
func init() {
	authorController = &AuthorController{
		BaseController: controller.BaseController{
			BaseService: service.GetAuthorService(),
		},
	}
	authorController.BaseController.ParseJSON = authorController.ParseJSON
	container.RegistController("author", authorController)
}
