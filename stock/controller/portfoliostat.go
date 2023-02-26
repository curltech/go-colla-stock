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
type PortfolioStatController struct {
	controller.BaseController
}

var portfolioStatController *PortfolioStatController

func (this *PortfolioStatController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.PortfolioStat, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type PortfolioStatPara struct {
	TsCode  string `json:"ts_code,omitempty"`
	Term    int64  `json:"term,omitempty"`
	From    int    `json:"from,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Orderby string `json:"orderby,omitempty"`
	Count   int64  `json:"count,omitempty"`
}

func (this *PortfolioStatController) FindPortfolioStat(ctx iris.Context) {
	portfolioStatPara := &PortfolioStatPara{}
	err := ctx.ReadJSON(&portfolioStatPara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.PortfolioStatService)
	ps, count, err := svc.FindPortfolioStat(portfolioStatPara.TsCode, portfolioStatPara.Term, portfolioStatPara.From, portfolioStatPara.Limit, portfolioStatPara.Orderby, portfolioStatPara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *PortfolioStatController) RefreshPortfolioStat(ctx iris.Context) {
	svc := this.BaseService.(*service.PortfolioStatService)
	err := svc.RefreshPortfolioStat()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *PortfolioStatController) GetUpdatePortfolioStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.PortfolioStatService)
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
	ps, err := svc.GetUpdatePortfolioStat(ts_code)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

/**
注册bean管理器，注册序列
*/

func init() {
	portfolioStatController = &PortfolioStatController{
		BaseController: controller.BaseController{
			BaseService: service.GetPortfolioStatService(),
		},
	}
	portfolioStatController.BaseController.ParseJSON = portfolioStatController.ParseJSON
	container.RegistController("portfoliostat", portfolioStatController)
}
