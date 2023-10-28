package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
)

// PerformanceService 同步表结构，服务继承基本服务的方法
type PerformanceService struct {
	service.OrmBaseService
}

var performanceService = &PerformanceService{}

func GetPerformanceService() *PerformanceService {
	return performanceService
}

func (svc *PerformanceService) GetSeqName() string {
	return seqname
}

func (svc *PerformanceService) NewEntity(data []byte) (interface{}, error) {
	performance := &entity.Performance{}
	if data == nil {
		return performance, nil
	}
	err := message.Unmarshal(data, performance)
	if err != nil {
		return nil, err
	}

	return performance, err
}

func (svc *PerformanceService) NewEntities(data []byte) (interface{}, error) {
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

func (svc *PerformanceService) findMaxQDate(securityCode string) (string, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	ps := make([]*entity.Performance, 0)
	err := svc.Find(&ps, nil, "qdate desc", 0, 1, conds, paras...)
	if err != nil {
		return "", err
	}
	if len(ps) > 0 {
		return ps[0].QDate, nil
	}

	return "", nil
}

func (svc *PerformanceService) FindByQDate(securityCode string, startDate string, endDate string, orderby string, from int, limit int, count int64) ([]*entity.Performance, int64, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	if startDate != "" {
		conds = conds + " and qdate>=?"
		paras = append(paras, startDate)
	}
	if endDate != "" {
		conds = conds + " and qdate<=?"
		paras = append(paras, endDate)
	}
	var err error
	condiBean := &entity.Performance{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "securitycode,qdate desc"
	}
	ps := make([]*entity.Performance, 0)
	if limit == 0 {
		limit = 10
	}
	err = svc.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	return ps, count, nil
}

func (svc *PerformanceService) FindLatest(securityCode string, latestNoticeDate string, orderby string, from int, limit int, count int64) ([]*entity.Performance, int64, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	if latestNoticeDate != "" {
		conds = conds + " and noticedate>=?"
		paras = append(paras, latestNoticeDate)
	}
	var err error
	condiBean := &entity.Performance{}
	if count == 0 {
		count, err = svc.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "securitycode,noticedate desc"
	}
	ps := make([]*entity.Performance, 0)
	if limit == 0 {
		limit = 10
	}
	err = svc.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	return ps, count, nil
}

func init() {
	err := service.GetSession().Sync(new(entity.Performance))
	if err != nil {
		return
	}
	performanceService.OrmBaseService.GetSeqName = performanceService.GetSeqName
	performanceService.OrmBaseService.FactNewEntity = performanceService.NewEntity
	performanceService.OrmBaseService.FactNewEntities = performanceService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("performance", performanceService)
}
