package service

import (
	"github.com/curltech/go-colla-core/util/json"
	"github.com/curltech/go-colla-stock/stock/entity"
)

type ShareRequest struct {
	TsCode     string `json:"ts_code,omitempty"`     // N	股票代码
	ListStatus string `json:"list_status,omitempty"` // N	上市状态: L上市 D退市 P暂停上市,默认L
	Exchange   string `json:"exchange,omitempty"`    // N	交易所 SSE上交所 SZSE深交所 HKEX港交所(未上线)
	IsHs       string `json:"is_hs,omitempty"`       // N  是否沪深港通标的,N否 H沪股通 S深股通
}

// StockBasic 获取基础信息数据,包括股票代码、名称、上市日期、退市日期等
func StockBasic(params ShareRequest) (tsData []*entity.Share, err error) {
	resp, err := FastPost(&TushareRequest{
		ApiName: "stock_basic",
		Token:   token,
		Params:  params,
		Fields:  reflectFields(entity.Share{}),
	})
	tsData = assembleStockBasic(resp)

	return tsData, err
}

func assembleStockBasic(tsRsp *TushareResponse) []*entity.Share {
	var tsData []*entity.Share
	if tsRsp.Data == nil {
		return nil
	}
	for _, data := range tsRsp.Data.Items {
		body, err := ReflectResponseData(tsRsp.Data.Fields, data)
		if err == nil {
			n := new(entity.Share)
			err = json.Unmarshal(body, &n)
			if err == nil {
				tsData = append(tsData, n)
			}
		}
	}
	return tsData
}
