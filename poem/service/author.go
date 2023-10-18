package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/poem/entity"
	"os"
)

// AuthorService 同步表结构，服务继承基本服务的方法
type AuthorService struct {
	service.OrmBaseService
}

var authorService = &AuthorService{}

func GetAuthorService() *AuthorService {
	return authorService
}

func (svc *AuthorService) GetSeqName() string {
	return seqname
}

func (svc *AuthorService) NewEntity(data []byte) (interface{}, error) {
	author := &entity.Author{}
	if data == nil {
		return author, nil
	}
	err := message.Unmarshal(data, author)
	if err != nil {
		return nil, err
	}

	return author, err
}

func (svc *AuthorService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Author, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (svc *AuthorService) ParseFile(src string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	as := make([]map[string]string, 0)
	err = json.Unmarshal(content, &as)
	if err != nil {
		return err
	}
	authors := make([]*entity.Author, 0)
	for _, a := range as {
		author := &entity.Author{Name: a["name"], Notes: a["description"]}
		authors = append(authors, author)
	}
	err = svc.save(authors)
	if err != nil {
		return err
	}

	return nil
}

func (svc *AuthorService) save(authors []*entity.Author) error {
	batch := 1000
	as := make([]interface{}, 0)
	for i := 0; i < len(authors); i = i + batch {
		for j := 0; j < batch; j++ {
			if i+j < len(authors) {
				poem := authors[i+j]
				as = append(as, poem)
			}
		}
		_, err := svc.Insert(as...)
		if err != nil {
			logger.Sugar.Errorf("Insert database error:%v", err.Error())
			return err
		} else {
			logger.Sugar.Infof("Insert database record:%v", len(as))
		}
		as = make([]interface{}, 0)
	}

	return nil
}

func init() {
	err := service.GetSession().Sync(new(entity.Author))
	if err != nil {
		return
	}
	authorService.OrmBaseService.GetSeqName = authorService.GetSeqName
	authorService.OrmBaseService.FactNewEntity = authorService.NewEntity
	authorService.OrmBaseService.FactNewEntities = authorService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("author", authorService)
}
