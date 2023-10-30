package controller

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/curltech/go-colla-stock/stock/service"
	"github.com/curltech/go-colla-web/controller"
	"github.com/kataras/iris/v12"
)

// QStatController 控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
type QStatController struct {
	controller.BaseController
}

var qstatController *QStatController

func (ctl *QStatController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.QStat, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type QStatPara struct {
	TsCode        string   `json:"ts_code,omitempty"`
	Source        []string `json:"source,omitempty"`
	Terms         []int    `json:"terms,omitempty"`
	SourceOptions []string `json:"source_options,omitempty"`
	From          int      `json:"from,omitempty"`
	Limit         int      `json:"limit,omitempty"`
	Orderby       string   `json:"orderby,omitempty"`
	Count         int64    `json:"count,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
}

func (ctl *QStatController) Search(ctx iris.Context) {
	qstatPara := &QStatPara{}
	err := ctx.ReadJSON(&qstatPara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.QStatService)
	ps, count, err := svc.Search(qstatPara.Keyword, qstatPara.Terms, qstatPara.SourceOptions, qstatPara.From, qstatPara.Limit, qstatPara.Count)
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

func (ctl *QStatController) FindQStatBy(ctx iris.Context) {
	qstatPara := &QStatPara{}
	err := ctx.ReadJSON(&qstatPara)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	if qstatPara.TsCode == "" {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := service.GetQStatService()
	ps, count, err := svc.FindQStatBy(qstatPara.TsCode, qstatPara.Terms, qstatPara.Source, qstatPara.Orderby, qstatPara.From, qstatPara.Limit, qstatPara.Count)
	if err != nil {
		err = ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	result := make(map[string]interface{})
	result["data"] = ps
	result["count"] = count
	err = ctx.JSON(result)
	if err != nil {
		return
	}
}

func (ctl *QStatController) RefreshQStat(ctx iris.Context) {
	svc := ctl.BaseService.(*service.QStatService)
	err := svc.RefreshQStat()
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}

func (ctl *QStatController) GetUpdateQStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		err := ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}

		return
	}
	svc := ctl.BaseService.(*service.QStatService)
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
	ps, err := svc.GetUpdateQStat(tsCode)
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
	qstatController = &QStatController{
		BaseController: controller.BaseController{
			BaseService: service.GetQStatService(),
		},
	}
	qstatController.BaseController.ParseJSON = qstatController.ParseJSON
	container.RegistController("qstat", qstatController)
}
