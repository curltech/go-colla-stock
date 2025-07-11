package controller

import (
	"errors"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
	"time"
)

// DayLineController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，
// 因此每个实体（外部rest方式访问）的控制层都需要写一遍
type DayLineController struct {
	controller.BaseController
}

var dayLineController *DayLineController

func (ctl *DayLineController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.DayLine, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *DayLineController) ParsePath(ctx iris.Context) {
	svc := ctl.BaseService.(*service.DayLineService)
	err := svc.ParsePath("C:\\zd_zsone\\vipdoc\\sz\\lday", "C:\\stock\\data\\origin\\lday")
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *DayLineController) RefreshDayLine(ctx iris.Context) {
	svc := ctl.BaseService.(*service.DayLineService)
	err := svc.RefreshDayLine(-1)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *DayLineController) RefreshTodayLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	var startDate int64
	v, ok := params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	if startDate == 0 {
		startDate = stock.CurrentDate()
	}
	err = svc.RefreshDayLine(startDate)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *DayLineController) GetUpdateDayLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
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
	ps, err := svc.GetUpdateDayline(tsCode, -1, 10000)
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

func (ctl *DayLineController) GetUpdateTodayLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
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
	var startDate int64
	v, ok = params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	if startDate == 0 {
		startDate = stock.CurrentDate()
	}
	ps, err := svc.GetUpdateDayline(tsCode, startDate, 10000)
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

func (ctl *DayLineController) StdPath(ctx iris.Context) {
	svc := ctl.BaseService.(*service.DayLineService)
	err := svc.StdPath("", "C:\\stock\\data\\minmax\\lday", "C:\\stock\\data\\standard\\lday", 20210713, 20211215)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *DayLineController) RefreshStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	var startDate int64
	v, ok := params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	svc := ctl.BaseService.(*service.DayLineService)
	err = svc.RefreshStat(startDate)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *DayLineController) UpdateStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode, _ = v.(string)
	}
	var startDate int64
	v, ok = params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err := ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps, err := svc.UpdateStat(tsCode, startDate)
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

func (ctl *DayLineController) RefreshBeforeMa(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	var startDate int64
	v, ok := params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	svc := ctl.BaseService.(*service.DayLineService)
	err = svc.RefreshBeforeMa(startDate)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *DayLineController) UpdateBeforeMa(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode, _ = v.(string)
	}
	var startDate int64
	v, ok = params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err := ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps, err := svc.UpdateBeforeMa(tsCode, startDate)
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

// DayLinePara 查询DayLine的参数结构，主要包括tsCode,tradeDate,startDate,endDate,industry,eventCode,condContent
// 还有分页参数
type DayLinePara struct {
	From           int           `json:"from,omitempty"`
	Limit          int           `json:"limit,omitempty"`
	Orderby        string        `json:"orderby,omitempty"`
	Count          int64         `json:"count,omitempty"`
	TsCode         string        `json:"ts_code,omitempty"`
	Industry       string        `json:"industry,omitempty"`
	Sector         string        `json:"sector,omitempty"`
	StartDate      int64         `json:"start_date,omitempty"`
	EndDate        int64         `json:"end_date,omitempty"`
	TradeDate      int64         `json:"trade_date,omitempty"`
	DayCount       string        `json:"day_count,omitempty"`
	SrcDayCount    string        `json:"src_day_count,omitempty"`
	TargetDayCount string        `json:"target_day_count,omitempty"`
	Cross          string        `json:"cross,omitempty"`
	CondContent    string        `json:"cond_content,omitempty"`
	CondParas      []interface{} `json:"cond_paras,omitempty"`
}

// Search 主要的查询参数进行查询
func (ctl *DayLineController) Search(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, count, err := svc.Search(dayLinePara.TsCode, dayLinePara.Industry, dayLinePara.Sector, dayLinePara.StartDate, dayLinePara.EndDate, dayLinePara.Orderby, dayLinePara.From, dayLinePara.Limit, dayLinePara.Count)
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

// FindNewest 查询最新的日期的股票日线数据
func (ctl *DayLineController) FindLatest(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, err := svc.FindLatest(dayLinePara.TsCode)
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

func (ctl *DayLineController) FindPreceding(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	start := time.Now()
	ps, count, err := svc.FindPreceding(dayLinePara.TsCode, dayLinePara.EndDate, dayLinePara.From, dayLinePara.Limit, dayLinePara.Count)
	end := time.Now()

	logger.Sugar.Infof("FindPreceding duration:%v", end.UnixMilli()-start.UnixMilli())
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

func (ctl *DayLineController) FindFollowing(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindFollowing(dayLinePara.TsCode, dayLinePara.StartDate, dayLinePara.EndDate, dayLinePara.From, dayLinePara.Limit, dayLinePara.Count)
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

func (ctl *DayLineController) FindRange(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, err := svc.FindRange(dayLinePara.TsCode, dayLinePara.StartDate, dayLinePara.EndDate, dayLinePara.Limit)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{})
	result["data"] = ps
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *DayLineController) FindHighest(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindHighest(dayLinePara.TsCode, dayLinePara.DayCount, dayLinePara.StartDate, dayLinePara.From, dayLinePara.Limit, dayLinePara.Count)
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

func (ctl *DayLineController) FindLowest(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindLowest(dayLinePara.TsCode, dayLinePara.DayCount, dayLinePara.StartDate, dayLinePara.From, dayLinePara.Limit, dayLinePara.Count)
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

func (ctl *DayLineController) FindMaCross(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindMaCross(dayLinePara.TsCode, dayLinePara.SrcDayCount, dayLinePara.TargetDayCount, dayLinePara.StartDate, dayLinePara.Cross, dayLinePara.From, dayLinePara.Limit, dayLinePara.Count)
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

// FindFlexPoint 最基本的查询买卖点的方法，最为灵活
// 条件包括tsCode，tradeDate，filterContent（filterParas），startDate，endDate
// 如果filterContent中有?，则filterParas中必须有对应的参数值
func (ctl *DayLineController) FindFlexPoint(ctx iris.Context) {
	//解析查询参数
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.CondContent == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("filterContent is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	dayLines, count, err := svc.FindFlexPoint(dayLinePara.TsCode, dayLinePara.TradeDate, dayLinePara.CondContent, dayLinePara.CondParas, dayLinePara.StartDate, dayLinePara.EndDate, dayLinePara.From, dayLinePara.Limit, dayLinePara.Count)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = dayLines
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *DayLineController) FindCorr(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if dayLinePara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindCorr(dayLinePara.TsCode, dayLinePara.StartDate, dayLinePara.From, dayLinePara.Limit, dayLinePara.Orderby, dayLinePara.Count)
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

func (ctl *DayLineController) WriteAllFile(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	err = svc.WriteAllFile(dayLinePara.StartDate)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *DayLineController) WriteFile(ctx iris.Context) {
	dayLinePara := &DayLinePara{}
	err := ctx.ReadJSON(&dayLinePara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.DayLineService)
	err = svc.WriteFile("", dayLinePara.TsCode, dayLinePara.StartDate)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
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
	dayLineController = &DayLineController{
		BaseController: controller.BaseController{
			BaseService: service.GetDayLineService(),
		},
	}
	dayLineController.BaseController.ParseJSON = dayLineController.ParseJSON
	container.RegistController("dayline", dayLineController)
}
