package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

/**
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type QPerformanceController struct {
	controller.BaseController
}

var qperformanceController *QPerformanceController

func (this *QPerformanceController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.QPerformance, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type QPerformancePara struct {
	Terms         []int    `json:"terms,omitempty"`
	Term          int      `json:"term,omitempty"`
	SourceOptions []string `json:"source_options,omitempty"`
	From          int      `json:"from,omitempty"`
	Limit         int      `json:"limit,omitempty"`
	Orderby       string   `json:"orderby,omitempty"`
	Count         int64    `json:"count,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
	TsCode        string   `json:"ts_code,omitempty"`
	StartDate     string   `json:"start_date,omitempty"`
	EndDate       string   `json:"end_date,omitempty"`
	TradeDate     int64    `json:"trade_date,omitempty"`
	RankType      string   `json:"rank_type,omitempty"`
	StdType       int      `json:"std_type,omitempty"`
	Winsorize     bool     `json:"winsorize,omitempty"`
}

func (this *QPerformanceController) Search(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.QPerformanceService)
	ps, count, err := svc.Search(qperformancePara.Keyword, qperformancePara.TsCode, qperformancePara.Terms, qperformancePara.SourceOptions, qperformancePara.StartDate, qperformancePara.EndDate, qperformancePara.Orderby, qperformancePara.From, qperformancePara.Limit, qperformancePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *QPerformanceController) FindStdQPerformance(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.QPerformanceService)
	ps, err := svc.FindStdQPerformance(qperformancePara.TsCode, qperformancePara.Terms, qperformancePara.StartDate, qperformancePara.EndDate, service.StdType(qperformancePara.StdType), qperformancePara.Winsorize)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	ctx.JSON(ps)
}

func (this *QPerformanceController) FindStat(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := service.GetQPerformanceService()
	startDate := qperformancePara.StartDate
	term := qperformancePara.Term
	if term > 0 && startDate == "" {
		today := stock.GetQTradeDate(0)
		startDate, _ = stock.AddYear(today, -term)
	}

	ps := svc.FindAllQStatBySql(qperformancePara.TsCode, startDate, qperformancePara.EndDate)

	ctx.JSON(ps)
}

func (this *QPerformanceController) FindPercentRank(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	startDate := qperformancePara.StartDate
	term := qperformancePara.Term
	if term > 0 && startDate == "" {
		today := stock.GetQTradeDate(0)
		startDate, _ = stock.AddYear(today, -term)
	}
	svc := service.GetQPerformanceService()
	ps, err := svc.FindPercentRank(qperformancePara.RankType, qperformancePara.TsCode, qperformancePara.TradeDate, startDate, qperformancePara.EndDate, qperformancePara.From, qperformancePara.Limit, qperformancePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
	ctx.JSON(ps)
}

func (this *QPerformanceController) RefreshQPerformance(ctx iris.Context) {
	svc := this.BaseService.(*service.QPerformanceService)
	err := svc.RefreshWmqyQPerformance("")
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
	err = svc.RefreshDayQPerformance()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *QPerformanceController) GetUpdateWmqyQPerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.QPerformanceService)
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
	ps, err := svc.GetUpdateWmqyQPerformance(ts_code, "")
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *QPerformanceController) GetUpdateDayQPerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.QPerformanceService)
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
	ps, err := svc.GetUpdateDayQPerformance(ts_code)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

/**
注册bean管理器，注册序列
*/

func init() {
	qperformanceController = &QPerformanceController{
		BaseController: controller.BaseController{
			BaseService: service.GetQPerformanceService(),
		},
	}
	qperformanceController.BaseController.ParseJSON = qperformanceController.ParseJSON
	container.RegistController("qperformance", qperformanceController)
}
