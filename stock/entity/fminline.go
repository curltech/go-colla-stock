package entity

import "github.com/curltech/go-colla-core/entity"

/**
通达信5分钟线*.lc5文件和*.lc1文件
    文件名即股票代码
    每32个字节为一个5分钟数据，每字段内低字节在前
    00 ~ 01 字节：日期，整型，设其值为num，则日期计算方法为：
                  year=floor(num/2048)+2004;
                  month=floor(mod(num,2048)/100);
                  day=mod(mod(num,2048),100);
    02 ~ 03 字节： 从0点开始至目前的分钟数，整型
    04 ~ 07 字节：开盘价，float型
    08 ~ 11 字节：最高价，float型
    12 ~ 15 字节：最低价，float型
    16 ~ 19 字节：收盘价，float型
    20 ~ 23 字节：成交额，float型
    24 ~ 27 字节：成交量（股），整型
    28 ~ 31 字节：（保留）
*/
type FminLine struct {
	StockLine   `xorm:"extends"`
	TradeMinute int64 `json:"trade_minute,omitempty"`
}

func (FminLine) TableName() string {
	return "stk_fminline"
}

func (FminLine) KeyName() string {
	return entity.FieldName_Id
}

func (FminLine) IdName() string {
	return entity.FieldName_Id
}
