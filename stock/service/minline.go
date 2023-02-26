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
type MinLineService struct {
	service.OrmBaseService
}

var minLineService = &MinLineService{}

func GetMinLineService() *MinLineService {
	return minLineService
}

func (this *MinLineService) GetSeqName() string {
	return seqname
}

func (this *MinLineService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.MinLine{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *MinLineService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.MinLine, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *MinLineService) AsyncParseFile(para interface{}) {
	src := (para.([]string))[0]
	target := (para.([]string))[1]
	filename := (para.([]string))[2]
	this.ParseFile(src, target, filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

/**
读目录下的数据
*/
func (this *MinLineService) ParsePath(src string, target string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	routinePool := thread.CreateRoutinePool(10, this.AsyncParseFile, nil)
	defer routinePool.Release()
	for _, file := range files {
		filename := file.Name()
		hasSuffix := strings.HasSuffix(filename, ".lc1")
		if hasSuffix {
			shareId := strings.TrimSuffix(filename, ".lc1")
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

func (this *MinLineService) ParseFile(src string, target string, filename string, flag int) error {
	shareId := strings.TrimSuffix(filename, ".lc5")
	logger.Sugar.Infof("shareId:%v", shareId)
	content, err := ioutil.ReadFile(src + string(os.PathSeparator) + filename)
	if err != nil {
		return err
	}
	targetFileName := target + string(os.PathSeparator) + shareId + ".csv"
	dayLines := this.ParseByte(shareId, content)
	raw := this.ToCsv(dayLines)
	file, err := os.OpenFile(targetFileName, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write([]byte(raw))
	logger.Sugar.Infof("Parse day file %v record %v completely!", targetFileName, len(dayLines))
	//this.save(dayLines)

	return nil
}

func (this *MinLineService) ToCsv(minLines []*entity.MinLine) string {
	raw := "id,trade_date,trade_minute,open,high,low,close,amount,vol\n"
	i := 0
	for _, minLine := range minLines {
		raw += fmt.Sprint(i) + ","
		raw += fmt.Sprint(minLine.TradeDate) + ","
		raw += fmt.Sprint(minLine.TradeMinute) + ","
		raw += fmt.Sprint(minLine.Open) + ","
		raw += fmt.Sprint(minLine.High) + ","
		raw += fmt.Sprint(minLine.Low) + ","
		raw += fmt.Sprint(minLine.Close) + ","
		raw += fmt.Sprint(minLine.Amount) + ","
		raw += fmt.Sprint(minLine.Vol) + "\n"
		i++
	}

	return raw
}

func (this *MinLineService) save(minLines []*entity.MinLine) error {
	batch := 1000
	mls := make([]interface{}, 0)
	for i := 0; i < len(minLines); i = i + batch {
		for j := 0; j < batch; j++ {
			if i+j < len(minLines) {
				minLine := minLines[i+j]
				mls = append(mls, minLine)
			}
		}
		_, err := this.Insert(mls...)
		if err != nil {
			logger.Sugar.Errorf("Insert database error:%v", err.Error())
			return err
		} else {
			logger.Sugar.Infof("Insert database record:%v", len(mls))
		}
		mls = make([]interface{}, 0)
	}

	return nil
}

func (this *MinLineService) ParseByte(shareId string, content []byte) []*entity.MinLine {
	minLines := make([]*entity.MinLine, 0)
	for i := 0; i < len(content); i = i + 32 {
		minLine := entity.MinLine{}
		minLine.TsCode = shareId
		num := stock.BytesToInt(content[i : i+2])
		year := int64(math.Floor(float64(num/2048))) + 2004
		month := int64(math.Floor(math.Mod(float64(num), 2048) / 100))
		day := int64(math.Mod(math.Mod(float64(num), 2048), 100))
		minLine.TradeDate = year*10000 + month*100 + day
		minLine.TradeMinute = int64(stock.BytesToInt(content[i+2 : i+4]))
		minLine.Open = stock.BytesToFloat64(content[i+4 : i+8])
		minLine.High = stock.BytesToFloat64(content[i+8 : i+12])
		minLine.Low = stock.BytesToFloat64(content[i+12 : i+16])
		minLine.Close = stock.BytesToFloat64(content[i+16 : i+20])
		minLine.Amount = stock.BytesToFloat64(content[i+20 : i+24])
		minLine.Vol = float64(stock.BytesToInt64(content[i+24 : i+28]))
		logger.Sugar.Infof("MinLine:%v", minLine)
		minLines = append(minLines, &minLine)
	}
	return minLines
}

/**
除非当天的分钟数据全部获取，每次访问都要重新获取网络的分钟数据
*/
func (this *MinLineService) FindMinLines(tscode string, tradeDate int64, tradeMinute int64) ([]*entity.MinLine, error) {
	minLines, err := this.findMinLines(tscode, tradeDate, tradeMinute)
	if err != nil {
		return minLines, err
	}
	if len(minLines) > 0 {
		lastMinute := minLines[len(minLines)-1].TradeMinute
		if lastMinute >= 900 {
			return minLines, err
		}
		currentMinute := stock.CurrentMinute()
		if currentMinute <= lastMinute {
			return minLines, err
		}
	}
	ps, err := this.GetUpdateTodayMinLine(tscode)
	if err == nil && len(ps) > 0 {
		today := stock.CurrentDate()
		daylines, err := GetDayLineService().GetUpdateDayline(tscode, today, 10000)
		if err == nil && len(daylines) > 0 {
			dayline := daylines[0]
			p := ps[len(ps)-1]
			dayline.MainNetInflow = p.MainNetInflow
			dayline.SuperNetInflow = p.SuperNetInflow
			dayline.SmallNetInflow = p.SmallNetInflow
			dayline.MiddleNetInflow = p.MiddleNetInflow
			dayline.LargeNetInflow = p.LargeNetInflow
			GetDayLineService().Upsert(dayline)
		}
		GetQPerformanceService().GetUpdateDayQPerformance(tscode)
		minLines, err = this.findMinLines(tscode, tradeDate, tradeMinute)
	}

	return minLines, err
}

func (this *MinLineService) findMinLines(tscode string, tradeDate int64, tradeMinute int64) ([]*entity.MinLine, error) {
	minLines := make([]*entity.MinLine, 0)
	condiBean := &entity.MinLine{}
	condiBean.TsCode = tscode
	if tradeDate == 0 {
		tradeDate = stock.CurrentDate()
	}
	condiBean.TradeDate = tradeDate
	conds := ""
	var paras []interface{} = nil
	if tradeMinute > 0 {
		conds = "trademinute>=?"
		paras = make([]interface{}, 0)
		paras = append(paras, tradeMinute)
	}
	err := this.Find(&minLines, condiBean, "trademinute", 0, 0, conds, paras)
	if err != nil {
		return nil, err
	}

	return minLines, nil
}

func init() {
	service.GetSession().Sync(new(entity.MinLine))
	minLineService.OrmBaseService.GetSeqName = minLineService.GetSeqName
	minLineService.OrmBaseService.FactNewEntity = minLineService.NewEntity
	minLineService.OrmBaseService.FactNewEntities = minLineService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("minline", minLineService)
}
