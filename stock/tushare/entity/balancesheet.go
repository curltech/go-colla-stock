package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type BalanceSheet struct {
	entity.BaseEntity     `xorm:"extends"`
	TsCode                string  `xorm:"varchar(255)" json:"ts_code,omitempty"` // str	Y	TS代码
	AnnDate               string  `json:"ann_date,omitempty"`                    // str Y 公告日期
	FAnnDate              string  `json:"f_ann_date,omitempty"`                  // str Y 实际公告日期
	EndDate               string  `json:"end_date,omitempty"`                    // str Y 报告期
	ReportType            string  `json:"report_type,omitempty"`                 // str Y 报表类型
	CompType              string  `json:"comp_type,omitempty"`                   // str Y 公司类型
	TotalShare            float64 `json:"total_share,omitempty"`                 // float Y 期末总股本
	CapRese               float64 `json:"cap_rese,omitempty"`                    // float Y 资本公积金
	UndistrPorfit         float64 `json:"undistr_porfit,omitempty"`              // float Y 未分配利润
	SurplusRese           float64 `json:"surplus_rese,omitempty"`                // float Y 盈余公积金
	SpecialRese           float64 `json:"special_rese,omitempty"`                // float Y 专项储备
	MoneyCap              float64 `json:"money_cap,omitempty"`                   // float Y 货币资金
	TradAsset             float64 `json:"trad_asset,omitempty"`                  // float Y 交易性金融资产
	NotesReceiv           float64 `json:"notes_receiv,omitempty"`                // float Y 应收票据
	AccountsReceiv        float64 `json:"accounts_receiv,omitempty"`             // float Y 应收账款
	OthReceiv             float64 `json:"oth_receiv,omitempty"`                  // float Y 其他应收款
	Prepayment            float64 `json:"prepayment,omitempty"`                  // float Y 预付款项
	DivReceiv             float64 `json:"div_receiv,omitempty"`                  // float Y 应收股利
	IntReceiv             float64 `json:"int_receiv,omitempty"`                  // float Y 应收利息
	Inventories           float64 `json:"inventories,omitempty"`                 // float Y 存货
	AmorExp               float64 `json:"amor_exp,omitempty"`                    // float Y 待摊费用
	NcaWithin1y           float64 `json:"nca_within_1y,omitempty"`               // float Y 一年内到期的非流动资产
	SettRsrv              float64 `json:"sett_rsrv,omitempty"`                   // float Y 结算备付金
	LoantoOthBankFi       float64 `json:"loanto_oth_bank_fi,omitempty"`          // float Y 拆出资金
	PremiumReceiv         float64 `json:"premium_receiv,omitempty"`              // float Y 应收保费
	ReinsurReceiv         float64 `json:"reinsur_receiv,omitempty"`              // float Y 应收分保账款
	ReinsurResReceiv      float64 `json:"reinsur_res_receiv,omitempty"`          // float Y 应收分保合同准备金
	PurResaleFa           float64 `json:"pur_resale_fa,omitempty"`               // float Y 买入返售金融资产
	OthCurAssets          float64 `json:"oth_cur_assets,omitempty"`              // float Y 其他流动资产
	TotalCurAssets        float64 `json:"total_cur_assets,omitempty"`            // float Y 流动资产合计
	FaAvailForSale        float64 `json:"fa_avail_for_sale,omitempty"`           // float Y 可供出售金融资产
	HtmInvest             float64 `json:"htm_invest,omitempty"`                  // float Y 持有至到期投资
	LtEqtInvest           float64 `json:"lt_eqt_invest,omitempty"`               // float Y 长期股权投资
	InvestRealEstate      float64 `json:"invest_real_estate,omitempty"`          // float Y 投资性房地产
	TimeDeposits          float64 `json:"time_deposits,omitempty"`               // float Y 定期存款
	OthAssets             float64 `json:"oth_assets,omitempty"`                  // float Y 其他资产
	LtRec                 float64 `json:"lt_rec,omitempty"`                      // float Y 长期应收款
	FixAssets             float64 `json:"fix_assets,omitempty"`                  // float Y 固定资产
	Cip                   float64 `json:"cip,omitempty"`                         // float Y 在建工程
	ConstMaterials        float64 `json:"const_materials,omitempty"`             // float Y 工程物资
	FixedAssetsDisp       float64 `json:"fixed_assets_disp,omitempty"`           // float Y 固定资产清理
	ProducBioAssets       float64 `json:"produc_bio_assets,omitempty"`           // float Y 生产性生物资产
	OilAndGasAssets       float64 `json:"oil_and_gas_assets,omitempty"`          // float Y 油气资产
	IntanAssets           float64 `json:"intan_assets,omitempty"`                // float Y 无形资产
	RAndD                 float64 `json:"r_and_d,omitempty"`                     // float Y 研发支出
	Goodwill              float64 `json:"goodwill,omitempty"`                    // float Y 商誉
	LtAmorExp             float64 `json:"lt_amor_exp,omitempty"`                 // float Y 长期待摊费用
	DeferTaxAssets        float64 `json:"defer_tax_assets,omitempty"`            // float Y 递延所得税资产
	DecrInDisbur          float64 `json:"decr_in_disbur,omitempty"`              // float Y 发放贷款及垫款
	OthNca                float64 `json:"oth_nca,omitempty"`                     // float Y 其他非流动资产
	TotalNca              float64 `json:"total_nca,omitempty"`                   // float Y 非流动资产合计
	CashReserCb           float64 `json:"cash_reser_cb,omitempty"`               // float Y 现金及存放中央银行款项
	DeposInOthBfi         float64 `json:"depos_in_oth_bfi,omitempty"`            // float Y 存放同业和其它金融机构款项
	PrecMetals            float64 `json:"prec_metals,omitempty"`                 // float Y 贵金属
	DerivAssets           float64 `json:"deriv_assets,omitempty"`                // float Y 衍生金融资产
	RrReinsUnePrem        float64 `json:"rr_reins_une_prem,omitempty"`           // float Y 应收分保未到期责任准备金
	RrReinsOutstdCla      float64 `json:"rr_reins_outstd_cla,omitempty"`         // float Y 应收分保未决赔款准备金
	RrReinsLinsLiab       float64 `json:"rr_reins_lins_liab,omitempty"`          // float Y 应收分保寿险责任准备金
	RrReinsLthinsLiab     float64 `json:"rr_reins_lthins_liab,omitempty"`        // float Y 应收分保长期健康险责任准备金
	RefundDepos           float64 `json:"refund_depos,omitempty"`                // float Y 存出保证金
	PhPledgeLoans         float64 `json:"ph_pledge_loans,omitempty"`             // float Y 保户质押贷款
	RefundCapDepos        float64 `json:"refund_cap_depos,omitempty"`            // float Y 存出资本保证金
	IndepAcctAssets       float64 `json:"indep_acct_assets,omitempty"`           // float Y 独立账户资产
	ClientDepos           float64 `json:"client_depos,omitempty"`                // float Y 其中：客户资金存款
	ClientProv            float64 `json:"client_prov,omitempty"`                 // float Y 其中：客户备付金
	TransacSeatFee        float64 `json:"transac_seat_fee,omitempty"`            // float Y 其中:交易席位费
	InvestAsReceiv        float64 `json:"invest_as_receiv,omitempty"`            // float Y 应收款项类投资
	TotalAssets           float64 `json:"total_assets,omitempty"`                // float Y 资产总计
	LtBorr                float64 `json:"lt_borr,omitempty"`                     // float Y 长期借款
	StBorr                float64 `json:"st_borr,omitempty"`                     // float Y 短期借款
	CbBorr                float64 `json:"cb_borr,omitempty"`                     // float Y 向中央银行借款
	DeposIbDeposits       float64 `json:"depos_ib_deposits,omitempty"`           // float Y 吸收存款及同业存放
	LoanOthBank           float64 `json:"loan_oth_bank,omitempty"`               // float Y 拆入资金
	TradingFl             float64 `json:"trading_fl,omitempty"`                  // float Y 交易性金融负债
	NotesPayable          float64 `json:"notes_payable,omitempty"`               // float Y 应付票据
	AcctPayable           float64 `json:"acct_payable,omitempty"`                // float Y 应付账款
	AdvReceipts           float64 `json:"adv_receipts,omitempty"`                // float Y 预收款项
	SoldForRepurFa        float64 `json:"sold_for_repur_fa,omitempty"`           // float Y 卖出回购金融资产款
	CommPayable           float64 `json:"comm_payable,omitempty"`                // float Y 应付手续费及佣金
	PayrollPayable        float64 `json:"payroll_payable,omitempty"`             // float Y 应付职工薪酬
	TaxesPayable          float64 `json:"taxes_payable,omitempty"`               // float Y 应交税费
	IntPayable            float64 `json:"int_payable,omitempty"`                 // float Y 应付利息
	DivPayable            float64 `json:"div_payable,omitempty"`                 // float Y 应付股利
	OthPayable            float64 `json:"oth_payable,omitempty"`                 // float Y 其他应付款
	AccExp                float64 `json:"acc_exp,omitempty"`                     // float Y 预提费用
	DeferredInc           float64 `json:"deferred_inc,omitempty"`                // float Y 递延收益
	StBondsPayable        float64 `json:"st_bonds_payable,omitempty"`            // float Y 应付短期债券
	PayableToReinsurer    float64 `json:"payable_to_reinsurer,omitempty"`        // float Y 应付分保账款
	RsrvInsurCont         float64 `json:"rsrv_insur_cont,omitempty"`             // float Y 保险合同准备金
	ActingTradingSec      float64 `json:"acting_trading_sec,omitempty"`          // float Y 代理买卖证券款
	ActingUwSec           float64 `json:"acting_uw_sec,omitempty"`               // float Y 代理承销证券款
	NonCurLiabDue1y       float64 `json:"non_cur_liab_due_1y,omitempty"`         // float Y 一年内到期的非流动负债
	OthCurLiab            float64 `json:"oth_cur_liab,omitempty"`                // float Y 其他流动负债
	TotalCurLiab          float64 `json:"total_cur_liab,omitempty"`              // float Y 流动负债合计
	BondPayable           float64 `json:"bond_payable,omitempty"`                // float Y 应付债券
	LtPayable             float64 `json:"lt_payable,omitempty"`                  // float Y 长期应付款
	SpecificPayables      float64 `json:"specific_payables,omitempty"`           // float Y 专项应付款
	EstimatedLiab         float64 `json:"estimated_liab,omitempty"`              // float Y 预计负债
	DeferTaxLiab          float64 `json:"defer_tax_liab,omitempty"`              // float Y 递延所得税负债
	DeferIncNonCurLiab    float64 `json:"defer_inc_non_cur_liab,omitempty"`      // float Y 递延收益-非流动负债
	OthNcl                float64 `json:"oth_ncl,omitempty"`                     // float Y 其他非流动负债
	TotalNcl              float64 `json:"total_ncl,omitempty"`                   // float Y 非流动负债合计
	DeposOthBfi           float64 `json:"depos_oth_bfi,omitempty"`               // float Y 同业和其它金融机构存放款项
	DerivLiab             float64 `json:"deriv_liab,omitempty"`                  // float Y 衍生金融负债
	Depos                 float64 `json:"depos,omitempty"`                       // float Y 吸收存款
	AgencyBusLiab         float64 `json:"agency_bus_liab,omitempty"`             // float Y 代理业务负债
	OthLiab               float64 `json:"oth_liab,omitempty"`                    // float Y 其他负债
	PremReceivAdva        float64 `json:"prem_receiv_adva,omitempty"`            // float Y 预收保费
	DeposReceived         float64 `json:"depos_received,omitempty"`              // float Y 存入保证金
	PhInvest              float64 `json:"ph_invest,omitempty"`                   // float Y 保户储金及投资款
	ReserUnePrem          float64 `json:"reser_une_prem,omitempty"`              // float Y 未到期责任准备金
	ReserOutstdClaims     float64 `json:"reser_outstd_claims,omitempty"`         // float Y 未决赔款准备金
	ReserLinsLiab         float64 `json:"reser_lins_liab,omitempty"`             // float Y 寿险责任准备金
	ReserLthinsLiab       float64 `json:"reser_lthins_liab,omitempty"`           // float Y 长期健康险责任准备金
	IndeptAccLiab         float64 `json:"indept_acc_liab,omitempty"`             // float Y 独立账户负债
	PledgeBorr            float64 `json:"pledge_borr,omitempty"`                 // float Y 其中:质押借款
	IndemPayable          float64 `json:"indem_payable,omitempty"`               // float Y 应付赔付款
	PolicyDivPayable      float64 `json:"policy_div_payable,omitempty"`          // float Y 应付保单红利
	TotalLiab             float64 `json:"total_liab,omitempty"`                  // float Y 负债合计
	TreasuryShare         float64 `json:"treasury_share,omitempty"`              // float Y 减:库存股
	OrdinRiskReser        float64 `json:"ordin_risk_reser,omitempty"`            // float Y 一般风险准备
	ForexDiffer           float64 `json:"forex_differ,omitempty"`                // float Y 外币报表折算差额
	InvestLossUnconf      float64 `json:"invest_loss_unconf,omitempty"`          // float Y 未确认的投资损失
	MinorityInt           float64 `json:"minority_int,omitempty"`                // float Y 少数股东权益
	TotalHldrEqyExcMinInt float64 `json:"total_hldr_eqy_exc_min_int,omitempty"`  // float Y 股东权益合计(不含少数股东权益)
	TotalHldrEqyIncMinInt float64 `json:"total_hldr_eqy_inc_min_int,omitempty"`  // float Y 股东权益合计(含少数股东权益)
	TotalLiabHldrEqy      float64 `json:"total_liab_hldr_eqy,omitempty"`         // float Y 负债及股东权益总计
	LtPayrollPayable      float64 `json:"lt_payroll_payable,omitempty"`          // float Y 长期应付职工薪酬
	OthCompIncome         float64 `json:"oth_comp_income,omitempty"`             // float Y 其他综合收益
	OthEqtTools           float64 `json:"oth_eqt_tools,omitempty"`               // float Y 其他权益工具
	OthEqtToolsPShr       float64 `json:"oth_eqt_tools_p_shr,omitempty"`         // float Y 其他权益工具(优先股)
	LendingFunds          float64 `json:"lending_funds,omitempty"`               // float Y 融出资金
	AccReceivable         float64 `json:"acc_receivable,omitempty"`              // float Y 应收款项
	StFinPayable          float64 `json:"st_fin_payable,omitempty"`              // float Y 应付短期融资款
	Payables              float64 `json:"payables,omitempty"`                    // float Y 应付款项
	HfsAssets             float64 `json:"hfs_assets,omitempty"`                  // float Y 持有待售的资产
	HfsSales              float64 `json:"hfs_sales,omitempty"`                   // float Y 持有待售的负债
	UpdateFlag            string  `json:"update_flag,omitempty"`                 // str N 更新标识
}

func (BalanceSheet) TableName() string {
	return "stk_balancesheet"
}

func (BalanceSheet) KeyName() string {
	return entity.FieldName_Id
}

func (BalanceSheet) IdName() string {
	return entity.FieldName_Id
}
