package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
)

/**
同步表结构，服务继承基本服务的方法
*/
type PinYinService struct {
	service.OrmBaseService
}

var pinyinService = &PinYinService{}

func GetPinYinService() *PinYinService {
	return pinyinService
}

func (this *PinYinService) GetSeqName() string {
	return seqname
}

func (this *PinYinService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.PinYin{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *PinYinService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.PinYin, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

var pinyinCache map[string]*entity.PinYin = nil

func GetCachePinYin() map[string]*entity.PinYin {
	if pinyinCache == nil {
		pinyinCache = make(map[string]*entity.PinYin, 0)
		pinyins := make([]*entity.PinYin, 0)
		svc := GetPinYinService()
		err := svc.Find(&pinyins, nil, "", 0, 0, "")
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
		ts_codes = make([]string, len(pinyins))
		i := 0
		for _, pinyin := range pinyins {
			pinyinCache[pinyin.ChineseChar] = pinyin
			i++
		}
	}

	return pinyinCache
}

func RefreshCachePinYin() {
	pinyinCache = nil
}

func init() {
	service.GetSession().Sync(new(entity.PinYin))
	pinyinService.OrmBaseService.GetSeqName = pinyinService.GetSeqName
	pinyinService.OrmBaseService.FactNewEntity = pinyinService.NewEntity
	pinyinService.OrmBaseService.FactNewEntities = pinyinService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("pinyin", pinyinService)
}
