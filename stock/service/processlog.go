package service

import (
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-stock/stock/entity"
	"github.com/robfig/cron"
	"time"
)

type ProcessLogService struct {
	service.OrmBaseService
}

var processLogService = &ProcessLogService{}

func GetProcessLogService() *ProcessLogService {
	return processLogService
}

func (this *ProcessLogService) GetSeqName() string {
	return seqname
}

func (this *ProcessLogService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.ProcessLog{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ProcessLogService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.ProcessLog, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *ProcessLogService) StartLog(name string, methodName string, bizCode string) *entity.ProcessLog {
	startDate := time.Now()
	processLog := &entity.ProcessLog{StartDate: &startDate}
	processLog.SchedualDate = &startDate
	processLog.Name = name
	processLog.BizCode = bizCode
	processLog.MethodName = methodName
	log, _ := json.TextMarshal(processLog)
	logger.Sugar.Infof("Start Log:%v", log)

	return processLog
}

func (this *ProcessLogService) EndLog(processLog *entity.ProcessLog, errorCode string, errorMsg string) {
	endDate := time.Now()
	processLog.EndDate = &endDate
	elapse := time.Since(*processLog.StartDate)
	processLog.Elapse = elapse.Milliseconds()
	processLog.ErrorCode = errorCode
	processLog.ErrorMsg = errorMsg
	log, _ := json.TextMarshal(processLog)
	logger.Sugar.Infof("End Log:%v", log)
	//go this.Insert(processLog)
}

func (this *ProcessLogService) Schedule() {
	processLog := this.StartLog("Schedule", "Schedule", "")
	err := dayLineService.RefreshDayLine(-1)
	if err != nil {
		logger.Sugar.Errorf("RefreshDayLine Error:%v", err.Error())
	}

	err = minLineService.RefreshMinLine(-1)
	if err != nil {
		logger.Sugar.Errorf("RefreshMinLine Error:%v", err.Error())
	}

	err = wmqyLineService.RefreshWmqyLine(-1)
	if err != nil {
		logger.Sugar.Errorf("RefreshWmqyLine Error:%v", err.Error())
	}

	err = forecastService.RefreshForecast()
	if err != nil {
		logger.Sugar.Errorf("RefreshForecast Error:%v", err.Error())
	}

	err = expressService.RefreshExpress()
	if err != nil {
		logger.Sugar.Errorf("RefreshExpress Error:%v", err.Error())
	}

	err = performanceService.RefreshPerformance()
	if err != nil {
		logger.Sugar.Errorf("RefreshPerformance Error:%v", err.Error())
	}

	err = qperformanceService.RefreshDayQPerformance()
	if err != nil {
		logger.Sugar.Errorf("RefreshDayQPerformance Error:%v", err.Error())
	}

	_, err = GetStatScoreService().CreateScorePercentile()
	if err != nil {
		logger.Sugar.Errorf("CreateScoreMedian Error:%v", err.Error())
	}

	this.EndLog(processLog, "", "")

	return
}

func (this *ProcessLogService) Cron() *cron.Cron {
	c := cron.New()
	c.AddFunc("0 0 16 * * 1-5", this.Schedule)
	c.Start()

	return c
}

func init() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	service.GetSession().Sync(new(entity.ProcessLog))
	processLogService.OrmBaseService.GetSeqName = processLogService.GetSeqName
	processLogService.OrmBaseService.FactNewEntity = processLogService.NewEntity
	processLogService.OrmBaseService.FactNewEntities = processLogService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("processlog", processLogService)

	go processLogService.Cron()
}
