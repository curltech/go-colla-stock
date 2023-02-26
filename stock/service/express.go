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
type ExpressService struct {
	service.OrmBaseService
}

var expressService = &ExpressService{}

func GetExpressService() *ExpressService {
	return expressService
}

func (this *ExpressService) GetSeqName() string {
	return seqname
}

func (this *ExpressService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Express{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ExpressService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Express, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *ExpressService) findMaxQDate(securityCode string) (string, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	ps := make([]*entity.Express, 0)
	err := this.Find(&ps, nil, "qdate desc", 0, 1, conds, paras...)
	if err != nil {
		return "", err
	}
	if len(ps) > 0 {
		return ps[0].QDate, nil
	}

	return "", nil
}

func (this *ExpressService) FindByQDate(securityCode string, startDate string, endDate string, orderby string) (map[string][]*entity.Express, error) {
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
	ps := make([]*entity.Express, 0)
	err := this.Find(&ps, nil, orderby, 0, 0, conds, paras...)
	if err != nil {
		return nil, err
	}
	psMap := make(map[string][]*entity.Express, 0)
	for _, p := range ps {
		qps, ok := psMap[p.SecurityCode]
		if !ok {
			qps = make([]*entity.Express, 0)
		}
		qps = append(qps, p)
		psMap[p.SecurityCode] = qps
	}

	return psMap, nil
}

func (this *ExpressService) FindLatest(securityCode string, latestDate string, orderby string, from int, limit int, count int64) ([]*entity.Express, int64, error) {
	conds, paras := stock.InBuildStr("securitycode", securityCode, ",")
	if latestDate != "" {
		conds = conds + " and noticedate>=?"
		paras = append(paras, latestDate)
	}
	var err error
	condiBean := &entity.Express{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	if orderby == "" {
		orderby = "securitycode,noticedate desc"
	}
	ps := make([]*entity.Express, 0)
	err = this.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	return ps, count, nil
}

func (this *ExpressService) Search(securityCode string, startDate string, endDate string, orderby string, from int, limit int, count int64) ([]*entity.Express, int64, error) {
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
	condiBean := &entity.Express{}
	if count == 0 {
		count, err = this.Count(condiBean, conds, paras...)
		if err != nil {
			return nil, count, err
		}
	}
	ps := make([]*entity.Express, 0)
	err = this.Find(&ps, nil, orderby, from, limit, conds, paras...)
	if err != nil {
		return nil, count, err
	}
	return ps, count, nil
}

func init() {
	service.GetSession().Sync(new(entity.Express))
	expressService.OrmBaseService.GetSeqName = expressService.GetSeqName
	expressService.OrmBaseService.FactNewEntity = expressService.NewEntity
	expressService.OrmBaseService.FactNewEntities = expressService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("express", expressService)
}
