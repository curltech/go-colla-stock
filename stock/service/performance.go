package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
)

/**
同步表结构，服务继承基本服务的方法
*/
type PerformanceService struct {
	service.OrmBaseService
}

var performanceService = &PerformanceService{}

func GetPerformanceService() *PerformanceService {
	return performanceService
}

func (this *PerformanceService) GetSeqName() string {
	return seqname
}

func (this *PerformanceService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Performance{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *PerformanceService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Performance, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *PerformanceService) findMaxQDate(securityCode string) (string, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	ps := make([]*entity.Performance, 0)
	err := this.Find(&ps, nil, "qdate desc", 0, 1, conds, paras...)
	if err != nil {
		return "", err
	}
	if len(ps) > 0 {
		return ps[0].QDate, nil
	}

	return "", nil
}

func (this *PerformanceService) FindByQDate(securityCode string, startDate string, endDate string, orderby string) (map[string][]*entity.Performance, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	if startDate != "" {
		conds = conds + " and qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		conds = conds + " and qdate<=?"
		paras = append(paras, endDate)
	}
	if orderby == "" {
		orderby = "securitycode,qdate desc"
	}
	ps := make([]*entity.Performance, 0)
	err := this.Find(&ps, nil, orderby, 0, 0, conds, paras...)
	if err != nil {
		return nil, err
	}
	psMap := make(map[string][]*entity.Performance, 0)
	for _, p := range ps {
		qps, ok := psMap[p.SecurityCode]
		if !ok {
			qps = make([]*entity.Performance, 0)
		}
		qps = append(qps, p)
		psMap[p.SecurityCode] = qps
	}
	return psMap, nil
}

func (this *PerformanceService) FindLatest(securityCode string, latestDate string, orderby string, from int, limit int, count int64) ([]*entity.Performance, int64, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	if latestDate != "" {
		conds = conds + " and noticedate>=?"
		paras = append(paras, latestDate)
	}
	var err error
	condiBean := &entity.Performance{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "securitycode,noticedate desc"
	}
	ps := make([]*entity.Performance, 0)
	err = this.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	return ps, count, nil
}

func (this *PerformanceService) Search(securityCode string, startDate string, endDate string, orderby string, from int, limit int, count int64) ([]*entity.Performance, int64, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	if startDate != "" {
		conds = conds + " and qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		conds = conds + " and qdate<=?"
		paras = append(paras, endDate)
	}
	if orderby == "" {
		orderby = "securitycode,qdate desc"
	}
	var err error
	condiBean := &entity.Performance{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	ps := make([]*entity.Performance, 0)
	err = this.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	return ps, count, nil
}

func init() {
	service.GetSession().Sync(new(entity.Performance))
	performanceService.OrmBaseService.GetSeqName = performanceService.GetSeqName
	performanceService.OrmBaseService.FactNewEntity = performanceService.NewEntity
	performanceService.OrmBaseService.FactNewEntities = performanceService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("performance", performanceService)
}
