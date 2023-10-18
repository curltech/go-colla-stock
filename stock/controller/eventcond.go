package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// EventCondController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type EventCondController struct {
	controller.BaseController
}

var eventCondController *EventCondController

func (ctl *EventCondController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.EventCond, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

func (ctl *EventCondController) RefreshEventCond(ctx iris.Context) {
	svc := ctl.BaseService.(*service.EventCondService)
	err := svc.RefreshEventCond()
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *EventCondController) GetUpdateEventCond(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.EventCondService)
	var tsCode string
	v, ok := params["ts_code"]
	if ok {
		tsCode, _ = v.(string)
	}
	if tsCode == "" {
		ps := make([]interface{}, 0)
		err = ctx.JSON(ps)
		if err != nil {
			return
		}
		return
	}
	ps := svc.GetUpdateEventCond(tsCode)
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

type EventCondPara struct {
	TsCode    string `json:"ts_code,omitempty"`
	TradeDate int64  `json:"trade_date,omitempty"`
	StartDate int64  `json:"start_date,omitempty"`
	EndDate   int64  `json:"end_date,omitempty"`
	EventCode string `json:"event_code,omitempty"`
	EventType string `json:"event_type,omitempty"`
	Orderby   string `json:"orderby,omitempty"`
	From      int    `json:"from"`
	Limit     int    `json:"limit"`
	Count     int64  `json:"count"`
}

func (ctl *EventCondController) Search(ctx iris.Context) {
	param := &EventCondPara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if param.TsCode == "" {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.EventCondService)
	es, count, err := svc.Search(param.TsCode, param.StartDate, param.EndDate, param.EventCode, param.EventType, param.Orderby, param.From, param.Limit, param.Count)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{})
	result["data"] = es
	result["count"] = count
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *EventCondController) FindGroupby(ctx iris.Context) {
	param := &EventCondPara{}
	err := ctx.ReadJSON(param)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.EventCondService)
	es, count, err := svc.FindGroupby(param.TsCode, param.StartDate, param.EndDate, param.EventCode, param.EventType, param.Orderby, param.From, param.Limit, param.Count)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	result := make(map[string]interface{})
	result["data"] = es
	result["count"] = count

	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

/*
*
注册bean管理器，注册序列
*/
func init() {
	eventCondController = &EventCondController{
		BaseController: controller.BaseController{
			BaseService: service.GetEventCondService(),
		},
	}
	eventCondController.BaseController.ParseJSON = eventCondController.ParseJSON
	container.RegistController("eventcond", eventCondController)
}
