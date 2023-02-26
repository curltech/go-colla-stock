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
type QStatController struct {
	controller.BaseController
}

var qstatController *QStatController

func (this *QStatController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.QStat, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

type QStatPara struct {
	Terms         []int    `json:"terms,omitempty"`
	SourceOptions []string `json:"source_options,omitempty"`
	From          int      `json:"from,omitempty"`
	Limit         int      `json:"limit,omitempty"`
	Orderby       string   `json:"orderby,omitempty"`
	Count         int64    `json:"count,omitempty"`
	Keyword       string   `json:"keyword,omitempty"`
}

func (this *QStatController) Search(ctx iris.Context) {
	qstatPara := &QStatPara{}
	err := ctx.ReadJSON(&qstatPara)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.QStatService)
	ps, count, err := svc.Search(qstatPara.Keyword, qstatPara.Terms, qstatPara.SourceOptions, qstatPara.From, qstatPara.Limit, qstatPara.Count)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["count"] = count
	result["data"] = ps
	ctx.JSON(result)
}

func (this *QStatController) FindQStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := service.GetQStatService()
	var ts_code string
	v, ok := params["ts_code"]
	if ok {
		ts_code = v.(string)
	}
	var term int
	v, ok = params["term"]
	if ok {
		f, ok := v.(float64)
		if ok && f > 0 {
			term = int(f)
		}
	}
	var source string
	v, ok = params["source"]
	if ok {
		source = v.(string)
	}
	var sourceName string
	v, ok = params["sourceName"]
	if ok {
		sourceName = v.(string)
	}
	terms := []int{term}
	ps, err := svc.FindQStat(ts_code, terms, source, sourceName)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	ctx.JSON(ps)
}

func (this *QStatController) RefreshQStat(ctx iris.Context) {
	svc := this.BaseService.(*service.QStatService)
	err := svc.RefreshQStat()
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}
}

func (this *QStatController) GetUpdateQStat(ctx iris.Context) {
	params := make(map[string]interface{})
	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	svc := this.BaseService.(*service.QStatService)
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
	ps, err := svc.GetUpdateQStat(ts_code)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
	}

	ctx.JSON(ps)
}

/**
注册bean管理器，注册序列
*/

func init() {
	qstatController = &QStatController{
		BaseController: controller.BaseController{
			BaseService: service.GetQStatService(),
		},
	}
	qstatController.BaseController.ParseJSON = qstatController.ParseJSON
	container.RegistController("qstat", qstatController)
}
