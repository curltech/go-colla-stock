package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// StatScoreController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type StatScoreController struct {
	controller.BaseController
}

var statScoreController *StatScoreController

func (ctl *StatScoreController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.StatScore, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type StatScorePara struct {
	Terms   []int  `json:"terms,omitempty"`
	From    int    `json:"from,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Orderby string `json:"orderby,omitempty"`
	Count   int64  `json:"count,omitempty"`
	Keyword string `json:"keyword,omitempty"`
	TsCode  string `json:"ts_code,omitempty"`
}

func (ctl *StatScoreController) Search(ctx iris.Context) {
	statScorePara := &StatScorePara{}
	err := ctx.ReadJSON(&statScorePara)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.StatScoreService)
	ps, count, err := svc.Search(statScorePara.Keyword, statScorePara.TsCode, statScorePara.Terms, statScorePara.Orderby, statScorePara.From, statScorePara.Limit, statScorePara.Count)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *StatScoreController) RefreshStatScore(ctx iris.Context) {
	svc := ctl.BaseService.(*service.StatScoreService)
	err := svc.RefreshStatScore()
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *StatScoreController) CreateScorePercentile(ctx iris.Context) {
	svc := ctl.BaseService.(*service.StatScoreService)
	_, err := svc.CreateScorePercentile()
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *StatScoreController) GetUpdateStatScore(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.StatScoreService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err = ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps, err := svc.GetUpdateStatScore(tsCode)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}

	err = ctx.JSON(ps)
	if err != nil {
		return
	}
}

// 注册bean管理器，注册序列
func init() {
	statScoreController = &StatScoreController{
		BaseController: controller.BaseController{
			BaseService: service.GetStatScoreService(),
		},
	}
	statScoreController.BaseController.ParseJSON = statScoreController.ParseJSON
	container.RegistController("statscore", statScoreController)
}
