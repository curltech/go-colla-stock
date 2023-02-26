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
type StatScoreController struct {
	controller.BaseController
}

var statScoreController *StatScoreController

func (this *StatScoreController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.StatScore, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type StatScorePara struct {
	Terms        []int    `json:"terms,omitempty"`
	ScoreOptions []string `json:"score_options,omitempty"`
	From         int      `json:"from,omitempty"`
	Limit        int      `json:"limit,omitempty"`
	Orderby      string   `json:"orderby,omitempty"`
	Count        int64    `json:"count,omitempty"`
	Keyword      string   `json:"keyword,omitempty"`
	TsCode       string   `json:"tscode,omitempty"`
}

func (this *StatScoreController) Search(ctx iris.Context) {
	statScorePara := &StatScorePara{}
	err := ctx.ReadJSON(&statScorePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.StatScoreService)
	ps, count, err := svc.Search(statScorePara.Keyword, statScorePara.TsCode, statScorePara.Terms, statScorePara.ScoreOptions, statScorePara.Orderby, statScorePara.From, statScorePara.Limit, statScorePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *StatScoreController) RefreshStatScore(ctx iris.Context) {
	svc := this.BaseService.(*service.StatScoreService)
	err := svc.RefreshStatScore()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *StatScoreController) CreateScorePercentile(ctx iris.Context) {
	svc := this.BaseService.(*service.StatScoreService)
	_, err := svc.CreateScorePercentile()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *StatScoreController) GetUpdateStatScore(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.StatScoreService)
	var ts_code string
	v, ok := params["ts_code"]
	if ok {
		ts_code = v.(string)
	}
	if ts_code == "" {
		ps := make([]interface{}, 0)
		ctx.JSON(ps)
		return
	}
	ps, err := svc.GetUpdateStatScore(ts_code)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

/**
注册bean管理器，注册序列
*/

func init() {
	statScoreController = &StatScoreController{
		BaseController: controller.BaseController{
			BaseService: service.GetStatScoreService(),
		},
	}
	statScoreController.BaseController.ParseJSON = statScoreController.ParseJSON
	container.RegistController("statscore", statScoreController)
}
