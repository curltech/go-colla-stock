package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/poem/entity"
	"os"
	"strings"
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

// ParsePath 读目录下的数据
func (svc *PoemService) ParsePath(src string) error {
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	routinePool := thread.CreateRoutinePool(10, svc.AsyncParseFile, nil)
	defer routinePool.Release()
	for _, file := range files {
		filename := file.Name()
		hasSuffix := strings.HasSuffix(filename, ".csv")
		if hasSuffix {
			para := make([]string, 0)
			para = append(para, src)
			para = append(para, filename)
			routinePool.Invoke(para)
		}
	}
	routinePool.Wait(nil)
	return nil
}

func (svc *PoemService) AsyncParseFile(para interface{}) {
	src := (para.([]string))[0]
	filename := (para.([]string))[1]
	err := svc.ParseFile(src, filename)
	if err != nil {
		return
	}
}

func (svc *PoemService) ParseFile(src string, filename string) error {
	content, err := os.ReadFile(src + "\\" + filename)
	if err != nil {
		return err
	}
	ps := strings.Split(string(content), "\n")
	i := 0
	poems := make([]*entity.Poem, 0)
	rhythmicService := GetRhythmicService()
	for _, p := range ps {
		if i > 0 {
			fs := strings.Split(p, "\",\"")
			if len(fs) == 4 {
				paragraphs := strings.Trim(fs[3], "\"")
				paragraphs = strings.Trim(paragraphs, "\r")
				paragraphs = strings.Trim(paragraphs, "\"")
				paragraphs = strings.ReplaceAll(paragraphs, "。", "。\n")
				title := strings.Trim(fs[0], "\"")
				poem := entity.Poem{Title: title, Dynasty: fs[1], Author: fs[2], Paragraphs: paragraphs}
				tss := strings.Split(title, " ")
				for _, ts := range tss {
					ss := strings.Split(ts, "/")
					for _, t := range ss {
						r := &entity.Rhythmic{Name: t}
						exist, _ := rhythmicService.Get(r, false, "", "")
						if exist {
							poem.Rhythmic = t
							break
						}
					}
				}
				poems = append(poems, &poem)
			} else {
				logger.Sugar.Errorf("error content:%s", fs)
			}
		}
		i++
	}
	err = svc.save(poems)
	if err != nil {
		return err
	}

	return nil
}

func (svc *PoemService) save(poems []*entity.Poem) error {
	batch := 1000
	ps := make([]interface{}, 0)
	for i := 0; i < len(poems); i = i + batch {
		for j := 0; j < batch; j++ {
			if i+j < len(poems) {
				poem := poems[i+j]
				ps = append(ps, poem)
			}
		}
		_, err := svc.Insert(ps...)
		if err != nil {
			logger.Sugar.Errorf("Insert database error:%v", err.Error())
			return err
		} else {
			logger.Sugar.Infof("Insert database record:%v", len(ps))
		}
		ps = make([]interface{}, 0)
	}

	return nil
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
