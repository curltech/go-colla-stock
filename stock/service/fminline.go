package service

import (
	"fmt"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/thread"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

/**
同步表结构，服务继承基本服务的方法
*/
type FminLineService struct {
	service.OrmBaseService
}

var fminLineService = &FminLineService{}

func GetFminLineService() *FminLineService {
	return fminLineService
}

func (this *FminLineService) GetSeqName() string {
	return seqname
}

func (this *FminLineService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.FminLine{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *FminLineService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.FminLine, 0)
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
读目录下的数据
*/
func (this *FminLineService) ParsePath(src string, target string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	routinePool := thread.CreateRoutinePool(10, this.AsyncParseFile, nil)
	defer routinePool.Release()
	for _, file := range files {
		filename := file.Name()
		hasSuffix := strings.HasSuffix(filename, ".lc5")
		if hasSuffix {
			shareId := strings.TrimSuffix(filename, ".lc5")
			logger.Sugar.Infof("shareId:", shareId)
			para := make([]string, 0)
			para = append(para, src)
			para = append(para, target)
			para = append(para, filename)
			routinePool.Invoke(para)
		}
	}
	routinePool.Wait(nil)
	stock.Rename(src, src+"-"+fmt.Sprint(stock.CurrentDate()))
	stock.Mkdir(src)
	return nil
}

func (this *FminLineService) AsyncParseFile(para interface{}) {
	src := (para.([]string))[0]
	target := (para.([]string))[1]
	filename := (para.([]string))[2]
	this.ParseFile(src, target, filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

func (this *FminLineService) ParseFile(src string, target string, filename string, flag int) error {
	shareId := strings.TrimSuffix(filename, ".lc5")
	logger.Sugar.Infof("shareId:%v", shareId)
	content, err := ioutil.ReadFile(src + string(os.PathSeparator) + filename)
	if err != nil {
		return err
	}
	targetFileName := target + string(os.PathSeparator) + shareId + ".csv"
	fminLines := this.ParseByte(shareId, content)
	raw := this.ToCsv(fminLines)
	file, err := os.OpenFile(targetFileName, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write([]byte(raw))
	logger.Sugar.Infof("Parse day file %v record %v completely!", targetFileName, len(fminLines))
	//this.save(fminLines)

	return nil
}

func (this *FminLineService) ToCsv(fminLines []*entity.FminLine) string {
	raw := "id,trade_date,trade_minute,open,high,low,close,amount,vol\n"
	i := 0
	for _, fminLine := range fminLines {
		raw += fmt.Sprint(i) + ","
		raw += fmt.Sprint(fminLine.TradeDate) + ","
		raw += fmt.Sprint(fminLine.TradeMinute) + ","
		raw += fmt.Sprint(fminLine.Open) + ","
		raw += fmt.Sprint(fminLine.High) + ","
		raw += fmt.Sprint(fminLine.Low) + ","
		raw += fmt.Sprint(fminLine.Close) + ","
		raw += fmt.Sprint(fminLine.Amount) + ","
		raw += fmt.Sprint(fminLine.Vol) + "\n"
		i++
	}

	return raw
}

func (this *FminLineService) save(fminLines []*entity.FminLine) error {
	batch := 1000
	fls := make([]interface{}, 0)
	for i := 0; i < len(fminLines); i = i + batch {
		for j := 0; j < batch; j++ {
			if i+j < len(fminLines) {
				fminLine := fminLines[i+j]
				fls = append(fls, fminLine)
			}
		}
		_, err := this.Insert(fls...)
		if err != nil {
			logger.Sugar.Errorf("Insert database error:%v", err.Error())
			return err
		} else {
			logger.Sugar.Infof("Insert database record:%v", len(fls))
		}
		fls = make([]interface{}, 0)
	}

	return nil
}

func (this *FminLineService) ParseByte(shareId string, content []byte) []*entity.FminLine {
	fminLines := make([]*entity.FminLine, 0)
	for i := 0; i < len(content); i = i + 32 {
		fminLine := entity.FminLine{}
		fminLine.TsCode = shareId
		num := stock.BytesToInt(content[i : i+2])
		year := int64(math.Floor(float64(num/2048))) + 2004
		month := int64(math.Floor(math.Mod(float64(num), 2048) / 100))
		day := int64(math.Mod(math.Mod(float64(num), 2048), 100))
		fminLine.TradeDate = year*10000 + month*100 + day
		fminLine.TradeMinute = int64(stock.BytesToInt(content[i+2 : i+4]))
		fminLine.Open = stock.BytesToFloat64(content[i+4 : i+8])
		fminLine.High = stock.BytesToFloat64(content[i+8 : i+12])
		fminLine.Low = stock.BytesToFloat64(content[i+12 : i+16])
		fminLine.Close = stock.BytesToFloat64(content[i+16 : i+20])
		fminLine.Amount = stock.BytesToFloat64(content[i+20 : i+24])
		fminLine.Vol = float64(stock.BytesToInt64(content[i+24 : i+28]))
		logger.Sugar.Infof("FminLine:%v", fminLine)
		fminLines = append(fminLines, &fminLine)
	}

	return fminLines
}

func init() {
	service.GetSession().Sync(new(entity.FminLine))
	fminLineService.OrmBaseService.GetSeqName = fminLineService.GetSeqName
	fminLineService.OrmBaseService.FactNewEntity = fminLineService.NewEntity
	fminLineService.OrmBaseService.FactNewEntities = fminLineService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("fminline", fminLineService)
}
