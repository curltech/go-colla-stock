package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/poem/entity"
)

var seqname = "seq_poem"

// PoemService 同步表结构，服务继承基本服务的方法
type PoemService struct {
	service.OrmBaseService
}

var poemService = &PoemService{}

func GetPoemService() *PoemService {
	return poemService
}

func (svc *PoemService) GetSeqName() string {
	return seqname
}

func (svc *PoemService) NewEntity(data []byte) (interface{}, error) {
	poem := &entity.Poem{}
	if data == nil {
		return poem, nil
	}
	err := message.Unmarshal(data, poem)
	if err != nil {
		return nil, err
	}

	return poem, err
}

func (svc *PoemService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Poem, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (svc *PoemService) Search(title string, author string, rhythmic string, paragraphs string, from int, limit int) ([]*entity.Poem, error) {
	conds := "1=1"
	paras := make([]interface{}, 0)
	if title != "" {
		conds = conds + " and title like ?"
		paras = append(paras, "%"+title+"%")
	}
	if author != "" {
		conds = conds + " and author like ?"
		paras = append(paras, "%"+author+"%")
	}
	if rhythmic != "" {
		conds = conds + " and rhythmic like ?"
		paras = append(paras, "%"+rhythmic+"%")
	}
	if paragraphs != "" {
		conds = conds + " and paragraphs like ?"
		paras = append(paras, "%"+paragraphs+"%")
	}
	poems := make([]*entity.Poem, 0)
	err := svc.Find(&poems, nil, "dynasty", from, limit, conds, paras...)
	if err != nil {
		return nil, err
	}

	return poems, nil
}

func init() {
	err := service.GetSession().Sync(new(entity.Poem))
	if err != nil {
		return
	}
	poemService.OrmBaseService.GetSeqName = poemService.GetSeqName
	poemService.OrmBaseService.FactNewEntity = poemService.NewEntity
	poemService.OrmBaseService.FactNewEntities = poemService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("poem", poemService)
}
