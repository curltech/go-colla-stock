package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type Income struct {
	entity.BaseEntity `xorm:"extends"`
	TsCode            string  `xorm:"varchar(255)" json:"ts_code,omitempty"` // str	Y	TS代码
	AnnDate           string  `json:"ann_date,omitempty"`                    // str	Y	公告日期
	FAnnDate          string  `json:"f_ann_date,omitempty"`                  // str	Y	实际公告日期
	EndDate           string  `json:"end_date,omitempty"`                    // str	Y	报告期
	ReportType        string  `json:"report_type,omitempty"`                 // str	Y	报告类型 1合并报表 2单季合并 3调整单季合并表 4调整合并报表 5调整前合并报表 6母公司报表 7母公司单季表 8 母公司调整单季表 9母公司调整表 10母公司调整前报表 11调整前合并报表 12母公司调整前报表
	CompType          string  `json:"comp_type,omitempty"`                   // str	Y	公司类型(1一般工商业2银行3保险4证券)
	BasicEps          float64 `json:"basic_eps,omitempty"`                   // float	Y	基本每股收益
	DilutedEps        float64 `json:"diluted_eps,omitempty"`                 // float	Y	稀释每股收益
	TotalRevenue      float64 `json:"total_revenue,omitempty"`               // float	Y	营业总收入
	Revenue           float64 `json:"revenue,omitempty"`                     // float	Y	营业收入
	IntIncome         float64 `json:"int_income,omitempty"`                  // float	Y	利息收入
	PremEarned        float64 `json:"prem_earned,omitempty"`                 // float	Y	已赚保费
	CommIncome        float64 `json:"comm_income,omitempty"`                 // float	Y	手续费及佣金收入
	NCommisIncome     float64 `json:"n_commis_income,omitempty"`             // float	Y	手续费及佣金净收入
	NOthIncome        float64 `json:"n_oth_income,omitempty"`                // float	Y	其他经营净收益
	NOthBIncome       float64 `json:"n_oth_b_income,omitempty"`              // float	Y	加:其他业务净收益
	PremIncome        float64 `json:"prem_income,omitempty"`                 // float	Y	保险业务收入
	OutPrem           float64 `json:"out_prem,omitempty"`                    // float	Y	减:分出保费
	UnePremReser      float64 `json:"une_prem_reser,omitempty"`              // float	Y	提取未到期责任准备金
	ReinsIncome       float64 `json:"reins_income,omitempty"`                // float	Y	其中:分保费收入
	NSecTbIncome      float64 `json:"n_sec_tb_income,omitempty"`             // float	Y	代理买卖证券业务净收入
	NSecUwIncome      float64 `json:"n_sec_uw_income,omitempty"`             // float	Y	证券承销业务净收入
	NAssetMgIncome    float64 `json:"n_asset_mg_income,omitempty"`           // float	Y	受托客户资产管理业务净收入
	OthBIncome        float64 `json:"oth_b_income,omitempty"`                // float	Y	其他业务收入
	FvValueChgGain    float64 `json:"fv_value_chg_gain,omitempty"`           // float	Y	加:公允价值变动净收益
	InvestIncome      float64 `json:"invest_income,omitempty"`               // float	Y	加:投资净收益
	AssInvestIncome   float64 `json:"ass_invest_income,omitempty"`           // float	Y	其中:对联营企业和合营企业的投资收益
	ForexGain         float64 `json:"forex_gain,omitempty"`                  // float	Y	加:汇兑净收益
	TotalCogs         float64 `json:"total_cogs,omitempty"`                  // float	Y	营业总成本
	OperCost          float64 `json:"oper_cost,omitempty"`                   // float	Y	减:营业成本
	IntExp            float64 `json:"int_exp,omitempty"`                     // float	Y	减:利息支出
	CommExp           float64 `json:"comm_exp,omitempty"`                    // float	Y	减:手续费及佣金支出
	BizTaxSurchg      float64 `json:"biz_tax_surchg,omitempty"`              // float	Y	减:营业税金及附加
	SellExp           float64 `json:"sell_exp,omitempty"`                    // float	Y	减:销售费用
	AdminExp          float64 `json:"admin_exp,omitempty"`                   // float	Y	减:管理费用
	FinExp            float64 `json:"fin_exp,omitempty"`                     // float	Y	减:财务费用
	AssetsImpairLoss  float64 `json:"assets_impair_loss,omitempty"`          // float	Y	减:资产减值损失
	PremRefund        float64 `json:"prem_refund,omitempty"`                 // float	Y	退保金
	CompensPayout     float64 `json:"compens_payout,omitempty"`              // float	Y	赔付总支出
	ReserInsurLiab    float64 `json:"reser_insur_liab,omitempty"`            // float	Y	提取保险责任准备金
	DivPayt           float64 `json:"div_payt,omitempty"`                    // float	Y	保户红利支出
	ReinsExp          float64 `json:"reins_exp,omitempty"`                   // float	Y	分保费用
	OperExp           float64 `json:"oper_exp,omitempty"`                    // float	Y	营业支出
	CompensPayoutRefu float64 `json:"compens_payout_refu,omitempty"`         // float	Y	减:摊回赔付支出
	InsurReserRefu    float64 `json:"insur_reser_refu,omitempty"`            // float	Y	减:摊回保险责任准备金
	ReinsCostRefund   float64 `json:"reins_cost_refund,omitempty"`           // float	Y	减:摊回分保费用
	OtherBusCost      float64 `json:"other_bus_cost,omitempty"`              // float	Y	其他业务成本
	OperateProfit     float64 `json:"operate_profit,omitempty"`              // float	Y	营业利润
	NonOperIncome     float64 `json:"non_oper_income,omitempty"`             // float	Y	加:营业外收入
	NonOperExp        float64 `json:"non_oper_exp,omitempty"`                // float	Y	减:营业外支出
	NcaDisploss       float64 `json:"nca_disploss,omitempty"`                // float	Y	其中:减:非流动资产处置净损失
	TotalProfit       float64 `json:"total_profit,omitempty"`                // float	Y	利润总额
	IncomeTax         float64 `json:"income_tax,omitempty"`                  // float	Y	所得税费用
	NIncome           float64 `json:"n_income,omitempty"`                    // float	Y	净利润(含少数股东损益)
	NIncomeAttrP      float64 `json:"n_income_attr_p,omitempty"`             // float	Y	净利润(不含少数股东损益)
	MinorityGain      float64 `json:"minority_gain,omitempty"`               // float	Y	少数股东损益
	OthComprIncome    float64 `json:"oth_compr_income,omitempty"`            // float	Y	其他综合收益
	TComprIncome      float64 `json:"t_compr_income,omitempty"`              // float	Y	综合收益总额
	ComprIncAttrP     float64 `json:"compr_inc_attr_p,omitempty"`            // float	Y	归属于母公司(或股东)的综合收益总额
	ComprIncAttrMS    float64 `json:"compr_inc_attr_m_s,omitempty"`          // float	Y	归属于少数股东的综合收益总额
	Ebit              float64 `json:"ebit,omitempty"`                        // float	Y	息税前利润
	Ebitda            float64 `json:"ebitda,omitempty"`                      // float	Y	息税折旧摊销前利润
	InsuranceExp      float64 `json:"insurance_exp,omitempty"`               // float	Y	保险业务支出
	UndistProfit      float64 `json:"undist_profit,omitempty"`               // float	Y	年初未分配利润
	DistableProfit    float64 `json:"distable_profit,omitempty"`             // float	Y	可分配利润
	UpdateFlag        string  `json:"update_flag,omitempty"`                 // str	N	更新标识,0未修改1更正过
}

func (Income) TableName() string {
	return "stk_income"
}

func (Income) KeyName() string {
	return entity.FieldName_Id
}

func (Income) IdName() string {
	return entity.FieldName_Id
}
