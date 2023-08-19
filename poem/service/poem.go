package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/poem/entity"
	"io/ioutil"
	"strings"
)

var seqname = "seq_poem"

/*
*
同步表结构，服务继承基本服务的方法
*/
type PoemService struct {
	service.OrmBaseService
}

var poemService = &PoemService{}

func GetPoemService() *PoemService {
	return poemService
}

func (this *PoemService) GetSeqName() string {
	return seqname
}

func (this *PoemService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Poem{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *PoemService) NewEntities(data []byte) (interface{}, error) {
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

/*
*
读目录下的数据
*/
func (this *PoemService) ParsePath(src string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	routinePool := thread.CreateRoutinePool(10, this.AsyncParseFile, nil)
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

func (this *PoemService) AsyncParseFile(para interface{}) {
	src := (para.([]string))[0]
	filename := (para.([]string))[1]
	this.ParseFile(src, filename)
}

func (this *PoemService) ParseFile(src string, filename string) error {
	content, err := ioutil.ReadFile(src + "\\" + filename)
	if err != nil {
		return err
	}
	ps := strings.Split(string(content), "\n")
	i := 0
	poems := make([]*entity.Poem, 0)
	svc := GetRhythmicService()
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
						exist, _ := svc.Get(r, false, "", "")
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
	this.save(poems)

	return nil
}

func (this *PoemService) save(poems []*entity.Poem) error {
	batch := 1000
	ps := make([]interface{}, 0)
	for i := 0; i < len(poems); i = i + batch {
		for j := 0; j < batch; j++ {
			if i+j < len(poems) {
				poem := poems[i+j]
				ps = append(ps, poem)
			}
		}
		_, err := this.Insert(ps...)
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
	service.GetSession().Sync(new(entity.Poem))
	poemService.OrmBaseService.GetSeqName = poemService.GetSeqName
	poemService.OrmBaseService.FactNewEntity = poemService.NewEntity
	poemService.OrmBaseService.FactNewEntities = poemService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("poem", poemService)
}
