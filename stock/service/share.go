package service

import (
	"errors"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/curltech/go-colla-stock/stock"
	"github.com/curltech/go-colla-stock/stock/entity"
	tushareSvc "github.com/curltech/go-colla-stock/stock/tushare/service"
	"strings"
)

type ShareRequest struct {
	TsCode     string `json:"ts_code,omitempty"`     // N	股票代码
	ListStatus string `json:"list_status,omitempty"` // N	上市状态: L上市 D退市 P暂停上市,默认L
	Exchange   string `json:"exchange,omitempty"`    // N	交易所 SSE上交所 SZSE深交所 HKEX港交所(未上线)
	IsHs       string `json:"is_hs,omitempty"`       // N  是否沪深港通标的,N否 H沪股通 S深股通
}

// ShareService 同步表结构，服务继承基本服务的方法
type ShareService struct {
	service.OrmBaseService
}

var shareService = &ShareService{}

func GetShareService() *ShareService {
	return shareService
}

var seqname = "seq_stock"

func (svc *ShareService) GetSeqName() string {
	return seqname
}

func (svc *ShareService) NewEntity(data []byte) (interface{}, error) {
	share := &entity.Share{}
	if data == nil {
		return share, nil
	}
	err := message.Unmarshal(data, share)
	if err != nil {
		return nil, err
	}

	return share, err
}

func (svc *ShareService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Share, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (svc *ShareService) save(shares []*entity.Share) error {
	batch := 1000
	dls := make([]interface{}, 0)
	for i := 0; i < len(shares); i = i + batch {
		for j := 0; j < batch; j++ {
			if i+j < len(shares) {
				share := shares[i+j]
				dls = append(dls, share)
			}
		}
		_, err := svc.Insert(dls...)
		if err != nil {
			logger.Sugar.Errorf("Insert database error:%v", err.Error())
			return err
		} else {
			logger.Sugar.Infof("Insert database record:%v", len(dls))
		}
		dls = make([]interface{}, 0)
	}

	return nil
}

type UserShare struct {
	Id                 int64   `json:"id,omitempty"`
	TsCode             string  `json:"ts_code,omitempty"`
	Name               string  `json:"name,omitempty"`
	Industry           string  `json:"industry,omitempty"`
	Sector             string  `json:"sector,omitempty"`
	TradeDate          int64   `json:"trade_date"`
	QDate              string  `json:"qdate"`
	Source             string  `json:"source"`
	Close              float64 `json:"close"`
	PctChgClose        float64 `json:"pct_chg_close"`
	PctChgVol          float64 `json:"pct_chg_vol"`
	Turnover           float64 `json:"turnover"`
	Pe                 float64 `json:"pe"`
	Peg                float64 `json:"peg"`
	Percent3Close      float64 `json:"percent3_close"`
	Percent5Close      float64 `json:"percent5_close"`
	Percent10Close     float64 `json:"percent10_close"`
	Percent13Close     float64 `json:"percent13_close"`
	Percent20Close     float64 `json:"percent20_close"`
	Percent21Close     float64 `json:"percent21_close"`
	Percent30Close     float64 `json:"percent30_close"`
	Percent34Close     float64 `json:"percent34_close"`
	Percent55Close     float64 `json:"percent55_close"`
	Percent60Close     float64 `json:"percent60_close"`
	Percent90Close     float64 `json:"percent90_close"`
	Percent120Close    float64 `json:"percent120_close"`
	Percent144Close    float64 `json:"percent144_close"`
	Percent233Close    float64 `json:"percent233_close"`
	Percent240Close    float64 `json:"percent240_close"`
	PercentPe          float64 `json:"percent_pe"`
	PercentPeg         float64 `json:"percent_peg"`
	IndustryPercentPe  float64 `json:"industry_percent_pe"`
	IndustryPercentPeg float64 `json:"industry_percent_peg"`
}

func (svc *ShareService) GetShares(tsCode string) ([]interface{}, error) {
	tscode := "000001"
	if tsCode != "" {
		tscode = tsCode
	}
	dayLines, err := GetDayLineService().findMaxTradeDate(tscode)
	if err != nil {
		return nil, err
	}
	if dayLines == nil || len(dayLines) == 0 || dayLines[0] == nil {
		return nil, errors.New("NoTradeDate")
	}
	sql := "select distinct on(s.tscode) s.tscode as ts_code,s.name,s.industry,s.sector,sd.tradedate as trade_date" +
		",qp.qdate as qdate,qp.source as source,sd.close,sd.pctchgclose as pct_chg_close" +
		",sd.pctchgvol as pct_chg_vol,sd.turnover as turnover,qp.pe as pe,qp.peg as peg" +
		",(sd.close-sd.min3close)/(sd.max3close-sd.min3close) as percent3_close,(sd.close-sd.min5close)/(sd.max5close-sd.min5close) as percent5_close" +
		",(sd.close-sd.min10close)/(sd.max10close-sd.min10close) as percent10_close,(sd.close-sd.min13close)/(sd.max13close-sd.min13close) as percent13_close" +
		",(sd.close-sd.min20close)/(sd.max20close-sd.min20close) as percent20_close,(sd.close-sd.min21close)/(sd.max21close-sd.min21close) as percent21_close" +
		",(sd.close-sd.min30close)/(sd.max30close-sd.min30close) as percent30_close,(sd.close-sd.min34close)/(sd.max34close-sd.min34close) as percent34_close" +
		",(sd.close-sd.min55close)/(sd.max55close-sd.min55close) as percent55_close,(sd.close-sd.min60close)/(sd.max60close-sd.min60close) as percent60_close" +
		",(sd.close-sd.min90close)/(sd.max90close-sd.min90close) as percent90_close,(sd.close-sd.min120close)/(sd.max120close-sd.min120close) as percent120_close" +
		",(sd.close-sd.min144close)/(sd.max144close-sd.min144close) as percent144_close,(sd.close-sd.min233close)/(sd.max233close-sd.min233close) as percent233_close" +
		",(sd.close-sd.min240close)/(sd.max240close-sd.min240close) as percent240_close" +
		",pqp.percent_pe as percent_pe,pqp.percent_peg as percent_peg" +
		",spqp.industry_percent_pe as industry_percent_pe,spqp.industry_percent_peg as industry_percent_peg" +
		" from stk_share s join stk_dayline sd on sd.tscode=s.tscode" +
		" join stk_qperformance qp on qp.tscode=s.tscode"
	percentSql := "(select tscode as tscode,tradedate as tradedate" +
		",percent_rank() over (partition by tscode order by pe asc) as percent_pe" +
		",percent_rank() over (partition by tscode order by peg asc) as percent_peg"
	in, paras := stock.InBuildStr("tscode", tsCode, ",")
	percentSql = percentSql + " from stk_qperformance" +
		" where " + in + ") pqp"
	industryPercentSql := "(select distinct on(tscode) tscode as tscode,tradedate as tradedate" +
		",percent_rank() over (partition by industry order by pe asc) as industry_percent_pe" +
		",percent_rank() over (partition by industry order by peg asc) as industry_percent_peg"
	industryPercentSql = industryPercentSql + " from stk_qperformance" +
		" order by tscode,tradedate desc) spqp"
	sql = sql + " join " + percentSql + " on pqp.tscode=s.tscode"
	sql = sql + " join " + industryPercentSql + " on spqp.tscode=s.tscode"
	sql = sql + " where sd.tradedate = qp.tradedate and pqp.tradedate = qp.tradedate and spqp.tradedate = qp.tradedate"
	sql = sql + " order by s.tscode,sd.tradedate desc"
	results, err := svc.Query(sql, paras...)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
		return nil, err
	}
	ps := make([]interface{}, 0)
	if results == nil || len(results) == 0 {
		err = errors.New("results is nil")
		logger.Sugar.Errorf("Error:%v", err.Error())
		return ps, nil
	}
	jsonMap, _, _ := stock.GetJsonMap(UserShare{})
	var i int64
	for _, result := range results {
		qp := &UserShare{}
		for colname, v := range result {
			err = reflect.Set(qp, jsonMap[colname], string(v))
			if err != nil {
				logger.Sugar.Errorf("Set colname %v value %v error", colname, string(v))
			}
		}
		i++
		qp.Id = i
		ps = append(ps, qp)
	}

	return ps, nil
}

func (svc *ShareService) Search(keyword string, from int, limit int) ([]*entity.Share, error) {
	conds := "name like ? or tscode like ? or pinyin like ?"
	paras := make([]interface{}, 0)
	paras = append(paras, keyword+"%")
	paras = append(paras, keyword+"%")
	paras = append(paras, strings.ToLower(keyword)+"%")
	shares := make([]*entity.Share, 0)
	err := svc.Find(&shares, nil, "industry,sector", from, limit, conds, paras...)
	if err != nil {
		return nil, err
	}

	return shares, nil
}

func (svc *ShareService) UpdatePinYin() {
	_, shareMap := svc.GetCacheShare()
	ps := make([]interface{}, 0)
	for _, share := range shareMap {
		if share.PinYin != "" {
			continue
		}
		py := getPinYin(share.Name)
		share.PinYin = py

		ps = append(ps, share)
	}

	_, err := svc.Upsert(ps...)
	if err != nil {
		return
	}
}

func (svc *ShareService) UpdateShares() {
	paras := tushareSvc.ShareRequest{}
	shares, err := tushareSvc.StockBasic(paras)
	if err != nil {
		logger.Sugar.Errorf("Error:%v", err.Error())
	}
	for _, share := range shares {
		sh := &entity.Share{TsCode: share.Symbol}
		ok, _ := svc.Get(sh, false, "", "", nil)
		if !ok {
			symbol := share.TsCode
			share.TsCode = share.Symbol
			share.Symbol = symbol
			py := getPinYin(share.Name)
			share.PinYin = py
			_, err := svc.Insert(share)
			if err != nil {
				logger.Sugar.Errorf("Error:%v", err.Error())
			}
		}
	}
}

func (svc *ShareService) UpdateSector(tsCode string, sector string) {
	tsCodes := strings.Split(tsCode, ",")
	sectors := strings.Split(sector, ",")
	_, shareMap := svc.GetCacheShare()
	ps := make([]interface{}, 0)

	for k, t := range tsCodes {
		share, ok := shareMap[t]
		if !ok {
			continue
		}
		if k >= len(sectors) {
			continue
		}
		sector := sectors[k]
		if sector == "" {
			continue
		}
		share.Sector = sector
		ps = append(ps, share)
	}

	_, err := svc.Upsert(ps...)
	if err != nil {
		return
	}
}

var shareCache map[string]*entity.Share = nil
var cacheTsCodes []string = nil

func (svc *ShareService) GetCacheShare() ([]string, map[string]*entity.Share) {
	if shareCache == nil {
		shareCache = make(map[string]*entity.Share, 0)
		shares := make([]*entity.Share, 0)
		err := GetShareService().Find(&shares, nil, "tscode", 0, 0, "")
		if err != nil {
			logger.Sugar.Errorf("Error: %s", err.Error())
		}
		cacheTsCodes = make([]string, len(shares))
		i := 0
		for _, share := range shares {
			shareCache[share.TsCode] = share
			cacheTsCodes[i] = share.TsCode
			i++
		}
	}

	return cacheTsCodes, shareCache
}

func (svc *ShareService) RefreshCacheShare() {
	shareCache = nil
	cacheTsCodes = nil
}

func getPinYin(name string) string {
	pinyinMap := GetCachePinYin()
	cs := []rune(name)
	l := len(cs)
	py := ""
	for i := 0; i <= l-1; i++ {
		c := string(cs[i : i+1])
		pinyin, ok := pinyinMap[c]
		if ok {
			py += pinyin.FirstChar
		}
	}
	return py
}

func init() {
	err := service.GetSession().Sync(new(entity.Share))
	if err != nil {
		return
	}
	shareService.OrmBaseService.GetSeqName = shareService.GetSeqName
	shareService.OrmBaseService.FactNewEntity = shareService.NewEntity
	shareService.OrmBaseService.FactNewEntities = shareService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("share", shareService)
}
