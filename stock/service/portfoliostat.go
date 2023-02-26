package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	entity "github.com/curltech/go-colla-stock/stock/entity"
)

type PortfolioStatService struct {
	service.OrmBaseService
}

var portfolioStatService = &PortfolioStatService{}

func GetPortfolioStatService() *PortfolioStatService {
	return portfolioStatService
}

func (this *PortfolioStatService) GetSeqName() string {
	return seqname
}

func (this *PortfolioStatService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.PortfolioStat{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *PortfolioStatService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.PortfolioStat, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

/**
删除股票季度业绩统计数据
*/
func (this *PortfolioStatService) deletePortfolioStat(ts_code string) error {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	portfolioStat := &entity.PortfolioStat{}
	_, err := this.Delete(portfolioStat, conds, paras...)
	if err != nil {
		return err
	}

	return nil
}

/**
刷新所有股票的季度业绩统计数据
*/
func (this *PortfolioStatService) RefreshPortfolioStat() error {
	processLog := GetProcessLogService().StartLog("portfolioStat", "RefreshPortfolioStat", "")
	routinePool := thread.CreateRoutinePool(10, this.AsyncUpdatePortfolioStat, nil)
	defer routinePool.Release()
	ts_codes, _ := GetShareService().GetCacheShare()
	for _, ts_code := range ts_codes {
		para := make([]interface{}, 0)
		para = append(para, ts_code)
		routinePool.Invoke(para)
	}
	routinePool.Wait(nil)

	GetProcessLogService().EndLog(processLog, "", "")

	return nil
}

func (this *PortfolioStatService) AsyncUpdatePortfolioStat(para interface{}) {
	tscode := (para.([]interface{}))[0].(string)
	this.GetUpdatePortfolioStat(tscode)
}

/**
更新股票季度业绩统计数据，并返回结果
*/
func (this *PortfolioStatService) GetUpdatePortfolioStat(tscode string) ([]interface{}, error) {
	var ps []interface{}
	var err error
	processLog := GetProcessLogService().StartLog("portfolioStat", "GetUpdatePortfolioStat", tscode)
	this.deletePortfolioStat(tscode)
	//for _, term := range this.Terms {
	//	ps, err = this.UpdatePortfolioStatBySql(tscode, term)
	//}
	err = this.UpdatePortfolioStat(tscode)
	if err != nil {
		GetProcessLogService().EndLog(processLog, "", err.Error())
		return ps, err
	}
	GetStatScoreService().GetUpdateStatScore(tscode)

	return ps, err
}

/**
通过内存更新股票季度业绩统计数据，并返回结果
*/
func (this *PortfolioStatService) UpdatePortfolioStat(tscode string) error {
	ps := make([]interface{}, 0)
	pstats, _, err := GetDayLineService().FindCorr(tscode, 0, 0, 100, "desc", 0)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return err
	}
	for _, pstat := range pstats {
		pstat.Term = 0
		ps = append(ps, pstat)
	}
	current := stock.CurrentDate()
	pstats, _, err = GetDayLineService().FindCorr(tscode, current-10000, 0, 100, "desc", 0)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return err
	}
	for _, pstat := range pstats {
		pstat.Term = 1
		ps = append(ps, pstat)
	}
	pstats, _, err = GetDayLineService().FindCorr(tscode, current-30000, 0, 100, "desc", 0)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return err
	}
	for _, pstat := range pstats {
		pstat.Term = 3
		ps = append(ps, pstat)
	}
	pstats, _, err = GetDayLineService().FindCorr(tscode, 0, 0, 100, "asc", 0)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return err
	}
	for _, pstat := range pstats {
		pstat.Term = 0
		ps = append(ps, pstat)
	}
	pstats, _, err = GetDayLineService().FindCorr(tscode, current-10000, 0, 100, "asc", 0)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return err
	}
	for _, pstat := range pstats {
		pstat.Term = 1
		ps = append(ps, pstat)
	}
	pstats, _, err = GetDayLineService().FindCorr(tscode, current-30000, 0, 100, "asc", 0)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return err
	}
	for _, pstat := range pstats {
		pstat.Term = 3
		ps = append(ps, pstat)
	}
	_, err = this.Insert(ps...)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return err
	}
	return nil
}

func (this *PortfolioStatService) FindPortfolioStat(ts_code string, term int64, from int, limit int, orderby string, count int64) ([]*entity.PortfolioStat, int64, error) {
	conds, paras := stock.InBuildStr("tscode", ts_code, ",")
	pstats := make([]*entity.PortfolioStat, 0)
	condiBean := &entity.PortfolioStat{}
	conds = conds + " and term=?"
	paras = append(paras, term)
	err := this.Find(&pstats, condiBean, orderby, from, limit, conds, paras...)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return nil, 0, err
	}
	if len(pstats) > 0 {
		return pstats, 200, nil
	}
	this.UpdatePortfolioStat(ts_code)
	err = this.Find(&pstats, condiBean, orderby, from, limit, conds, paras...)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return nil, 0, err
	}
	return pstats, 200, nil
}

func init() {
	service.GetSession().Sync(new(entity.PortfolioStat))
	portfolioStatService.OrmBaseService.GetSeqName = portfolioStatService.GetSeqName
	portfolioStatService.OrmBaseService.FactNewEntity = portfolioStatService.NewEntity
	portfolioStatService.OrmBaseService.FactNewEntities = portfolioStatService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("portfoliostat", portfolioStatService)
}
