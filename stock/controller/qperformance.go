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

// QPerformanceController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type QPerformanceController struct {
	controller.BaseController
}

var qperformanceController *QPerformanceController

func (ctl *QPerformanceController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.QPerformance, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type QPerformancePara struct {
	Terms         []int         `json:"terms,omitempty"`
	Term          int           `json:"term,omitempty"`
	SourceOptions []string      `json:"source_options,omitempty"`
	From          int           `json:"from,omitempty"`
	Limit         int           `json:"limit,omitempty"`
	Orderby       string        `json:"orderby,omitempty"`
	Count         int64         `json:"count,omitempty"`
	Keyword       string        `json:"keyword,omitempty"`
	TsCode        string        `json:"ts_code,omitempty"`
	QDate         string        `json:"qdate,omitempty"`
	StartDate     string        `json:"start_date,omitempty"`
	EndDate       string        `json:"end_date,omitempty"`
	TradeDate     int64         `json:"trade_date,omitempty"`
	RankType      string        `json:"rank_type,omitempty"`
	StdType       int           `json:"std_type,omitempty"`
	Winsorize     bool          `json:"winsorize,omitempty"`
	CondContent   string        `json:"cond_content,omitempty"`
	CondParas     []interface{} `json:"cond_paras,omitempty"`
}

func (ctl *QPerformanceController) FindByQDate(ctx iris.Context) {
	param := &QPerformancePara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.QPerformanceService)
	es, count, err := svc.FindByCondContent(param.TsCode, param.QDate, param.TradeDate, param.CondContent, param.CondParas, param.Orderby, param.From, param.Limit, param.Count)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{}, 0)
	result["data"] = es
	result["count"] = count
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *QPerformanceController) Search(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.QPerformanceService)
	ps, count, err := svc.Search(qperformancePara.Keyword, qperformancePara.Terms, qperformancePara.SourceOptions, qperformancePara.StartDate, qperformancePara.EndDate, qperformancePara.CondContent, qperformancePara.CondParas, qperformancePara.Orderby, qperformancePara.From, qperformancePara.Limit, qperformancePara.Count)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
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

func (ctl *QPerformanceController) FindStdQPerformance(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.QPerformanceService)
	ps, err := svc.FindStdQPerformance(qperformancePara.TsCode, qperformancePara.Terms, qperformancePara.StartDate, qperformancePara.EndDate, service.StdType(qperformancePara.StdType), qperformancePara.Winsorize)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	err = ctx.JSON(ps)
	if err != nil {
		return
	}
}

func (ctl *QPerformanceController) FindStat(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

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

	err = ctx.JSON(ps)
	if err != nil {
		return
	}
}

func (ctl *QPerformanceController) FindPercentRank(ctx iris.Context) {
	qperformancePara := &QPerformancePara{}
	err := ctx.ReadJSON(&qperformancePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

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
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
	err = ctx.JSON(ps)
	if err != nil {
		return
	}
}

func (ctl *QPerformanceController) RefreshQPerformance(ctx iris.Context) {
	svc := ctl.BaseService.(*service.QPerformanceService)
	err := svc.RefreshWmqyQPerformance("")
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
	err = svc.RefreshDayQPerformance()
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *QPerformanceController) GetUpdateWmqyQPerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.QPerformanceService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err := ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps, err := svc.GetUpdateWmqyQPerformance(tsCode, "")
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}

	err = ctx.JSON(ps)
	if err != nil {
		return
	}
}

func (ctl *QPerformanceController) GetUpdateDayQPerformance(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.QPerformanceService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err := ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps, err := svc.GetUpdateDayQPerformance(tsCode)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
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
	qperformanceController = &QPerformanceController{
		BaseController: controller.BaseController{
			BaseService: service.GetQPerformanceService(),
		},
	}
	qperformanceController.BaseController.ParseJSON = qperformanceController.ParseJSON
	container.RegistController("qperformance", qperformanceController)
}
