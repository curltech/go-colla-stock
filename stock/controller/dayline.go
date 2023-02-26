package controller

import (
	"errors"
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
type DayLineController struct {
	controller.BaseController
}

var dayLineController *DayLineController

func (this *DayLineController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.DayLine, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (this *DayLineController) ParsePath(ctx iris.Context) {
	svc := this.BaseService.(*service.DayLineService)
	err := svc.ParsePath("C:\\zd_zsone\\vipdoc\\sz\\lday", "C:\\stock\\data\\origin\\lday")
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *DayLineController) RefreshDayLine(ctx iris.Context) {
	svc := this.BaseService.(*service.DayLineService)
	err := svc.RefreshDayLine(-1)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *DayLineController) RefreshTodayLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
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
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *DayLineController) GetUpdateDayLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
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
	ps, err := svc.GetUpdateDayline(ts_code, -1, 10000)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *DayLineController) GetUpdateTodayLine(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
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
	ps, err := svc.GetUpdateDayline(ts_code, startDate, 10000)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *DayLineController) StdPath(ctx iris.Context) {
	svc := this.BaseService.(*service.DayLineService)
	err := svc.StdPath("", "C:\\stock\\data\\minmax\\lday", "C:\\stock\\data\\standard\\lday", 20210713, 20211215)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *DayLineController) RefreshStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

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
	svc := this.BaseService.(*service.DayLineService)
	err = svc.RefreshStat(startDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *DayLineController) UpdateStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	var ts_code string
	v, ok := params["ts_code"]
	if ok {
		ts_code, _ = v.(string)
	}
	var startDate int64
	v, ok = params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	if ts_code == "" {
		ps := make([]interface{}, 0)
		ctx.JSON(ps)
		return
	}
	ps, err := svc.UpdateStat(ts_code, startDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

func (this *DayLineController) RefreshBeforeMa(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

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
	svc := this.BaseService.(*service.DayLineService)
	err = svc.RefreshBeforeMa(startDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *DayLineController) UpdateBeforeMa(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	var ts_code string
	v, ok := params["ts_code"]
	if ok {
		ts_code, _ = v.(string)
	}
	var startDate int64
	v, ok = params["start_date"]
	if ok {
		f, ok := v.(float64)
		if ok {
			startDate = int64(f)
		}
	}
	if ts_code == "" {
		ps := make([]interface{}, 0)
		ctx.JSON(ps)
		return
	}
	ps, err := svc.UpdateBeforeMa(ts_code, startDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

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
	EventCode      string        `json:"event_code,omitempty"`
	EventContent   string        `json:"event_content,omitempty"`
	FilterParas    []interface{} `json:"filter_paras,omitempty"`
	CompareValue   float64       `json:"compare_value,omitempty"`
	CondNum        int           `json:"cond_num,omitempty"`
}

func (this *DayLineController) Search(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, count, err := svc.Search(daylinePara.TsCode, daylinePara.Industry, daylinePara.Sector, daylinePara.StartDate, daylinePara.EndDate, daylinePara.Orderby, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) FindPreceding(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindPreceding(daylinePara.TsCode, daylinePara.EndDate, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) FindFollowing(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindFollowing(daylinePara.TsCode, daylinePara.StartDate, daylinePara.EndDate, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) FindRange(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, err := svc.FindRange(daylinePara.TsCode, daylinePara.StartDate, daylinePara.EndDate, daylinePara.Limit)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) FindHighest(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindHighest(daylinePara.TsCode, daylinePara.DayCount, daylinePara.StartDate, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) FindLowestest(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindLowest(daylinePara.TsCode, daylinePara.DayCount, daylinePara.StartDate, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) FindMaCross(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindMaCross(daylinePara.TsCode, daylinePara.SrcDayCount, daylinePara.TargetDayCount, daylinePara.StartDate, daylinePara.Cross, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) FindFlexPoint(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	inOutPoint, err := svc.FindFlexPoint(daylinePara.TsCode, 0, nil, daylinePara.EventContent, daylinePara.FilterParas, daylinePara.StartDate, daylinePara.EndDate, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(inOutPoint)
}

func (this *DayLineController) FindInOutEvent(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	inOutPoint, err := svc.FindInOutEvent(daylinePara.TsCode, daylinePara.TradeDate, daylinePara.EventCode, daylinePara.FilterParas, daylinePara.StartDate, daylinePara.EndDate, daylinePara.CompareValue, daylinePara.CondNum, daylinePara.From, daylinePara.Limit, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(inOutPoint)
}

func (this *DayLineController) FindAllInOutEvent(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	inOutPoint := svc.FindAllInOutEvent(daylinePara.TsCode, daylinePara.EventCode, daylinePara.StartDate, daylinePara.EndDate, daylinePara.CompareValue, daylinePara.CondNum)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(inOutPoint)
}

func (this *DayLineController) FindCorr(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	if daylinePara.TsCode == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, errors.New("tscode is nil"))

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	ps, count, err := svc.FindCorr(daylinePara.TsCode, daylinePara.StartDate, daylinePara.From, daylinePara.Limit, daylinePara.Orderby, daylinePara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *DayLineController) WriteAllFile(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	err = svc.WriteAllFile(daylinePara.StartDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *DayLineController) WriteFile(ctx iris.Context) {
	daylinePara := &DayLinePara{}
	err := ctx.ReadJSON(&daylinePara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.DayLineService)
	err = svc.WriteFile("", daylinePara.TsCode, daylinePara.StartDate)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

/**
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
