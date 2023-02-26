package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type CashFlow struct {
	entity.BaseEntity       `xorm:"extends"`
	TsCode                  string  `xorm:"varchar(255)" json:"ts_code,omitempty"` // str	Y	TS代码
	AnnDate                 string  `json:"ann_date,omitempty"`                    // str	Y	公告日期
	FAnnDate                string  `json:"f_ann_date,omitempty"`                  // str	Y	实际公告日期
	EndDate                 string  `json:"end_date,omitempty"`                    // str	Y	报告期
	CompType                string  `json:"comp_type,omitempty"`                   // str	Y	公司类型
	ReportType              string  `json:"report_type,omitempty"`                 // str	Y	报表类型
	NetProfit               float64 `json:"net_profit,omitempty"`                  // float	Y	净利润
	FinanExp                float64 `json:"finan_exp,omitempty"`                   // float	Y	财务费用
	CFrSaleSg               float64 `json:"c_fr_sale_sg,omitempty"`                // float	Y	销售商品、提供劳务收到的现金
	RecpTaxRends            float64 `json:"recp_tax_rends,omitempty"`              // float	Y	收到的税费返还
	NDeposIncrFi            float64 `json:"n_depos_incr_fi,omitempty"`             // float	Y	客户存款和同业存放款项净增加额
	NIncrLoansCb            float64 `json:"n_incr_loans_cb,omitempty"`             // float	Y	向中央银行借款净增加额
	NIncBorrOthFi           float64 `json:"n_inc_borr_oth_fi,omitempty"`           // float	Y	向其他金融机构拆入资金净增加额
	PremFrOrigContr         float64 `json:"prem_fr_orig_contr,omitempty"`          // float	Y	收到原保险合同保费取得的现金
	NIncrInsuredDep         float64 `json:"n_incr_insured_dep,omitempty"`          // float	Y	保户储金净增加额
	NReinsurPrem            float64 `json:"n_reinsur_prem,omitempty"`              // float	Y	收到再保业务现金净额
	NIncrDispTfa            float64 `json:"n_incr_disp_tfa,omitempty"`             // float	Y	处置交易性金融资产净增加额
	IfcCashIncr             float64 `json:"ifc_cash_incr,omitempty"`               // float	Y	收取利息和手续费净增加额
	NIncrDispFaas           float64 `json:"n_incr_disp_faas,omitempty"`            // float	Y	处置可供出售金融资产净增加额
	NIncrLoansOthBank       float64 `json:"n_incr_loans_oth_bank,omitempty"`       // float	Y	拆入资金净增加额
	NCapIncrRepur           float64 `json:"n_cap_incr_repur,omitempty"`            // float	Y	回购业务资金净增加额
	CFrOthOperateA          float64 `json:"c_fr_oth_operate_a,omitempty"`          // float	Y	收到其他与经营活动有关的现金
	CInfFrOperateA          float64 `json:"c_inf_fr_operate_a,omitempty"`          // float	Y	经营活动现金流入小计
	CPaidGoodsS             float64 `json:"c_paid_goods_s,omitempty"`              // float	Y	购买商品、接受劳务支付的现金
	CPaidToForEmpl          float64 `json:"c_paid_to_for_empl,omitempty"`          // float	Y	支付给职工以及为职工支付的现金
	CPaidForTaxes           float64 `json:"c_paid_for_taxes,omitempty"`            // float	Y	支付的各项税费
	NIncrCltLoanAdv         float64 `json:"n_incr_clt_loan_adv,omitempty"`         // float	Y	客户贷款及垫款净增加额
	NIncrDepCbob            float64 `json:"n_incr_dep_cbob,omitempty"`             // float	Y	存放央行和同业款项净增加额
	CPayClaimsOrigInco      float64 `json:"c_pay_claims_orig_inco,omitempty"`      // float	Y	支付原保险合同赔付款项的现金
	PayHandlingChrg         float64 `json:"pay_handling_chrg,omitempty"`           // float	Y	支付手续费的现金
	PayCommInsurPlcy        float64 `json:"pay_comm_insur_plcy,omitempty"`         // float	Y	支付保单红利的现金
	OthCashPayOperAct       float64 `json:"oth_cash_pay_oper_act,omitempty"`       // float	Y	支付其他与经营活动有关的现金
	StCashOutAct            float64 `json:"st_cash_out_act,omitempty"`             // float	Y	经营活动现金流出小计
	NCashflowAct            float64 `json:"n_cashflow_act,omitempty"`              // float	Y	经营活动产生的现金流量净额
	OthRecpRalInvAct        float64 `json:"oth_recp_ral_inv_act,omitempty"`        // float	Y	收到其他与投资活动有关的现金
	CDispWithdrwlInvest     float64 `json:"c_disp_withdrwl_invest,omitempty"`      // float	Y	收回投资收到的现金
	CRecpReturnInvest       float64 `json:"c_recp_return_invest,omitempty"`        // float	Y	取得投资收益收到的现金
	NRecpDispFiolta         float64 `json:"n_recp_disp_fiolta,omitempty"`          // float	Y	处置固定资产、无形资产和其他长期资产收回的现金净额
	NRecpDispSobu           float64 `json:"n_recp_disp_sobu,omitempty"`            // float	Y	处置子公司及其他营业单位收到的现金净额
	StotInflowsInvAct       float64 `json:"stot_inflows_inv_act,omitempty"`        // float	Y	投资活动现金流入小计
	CPayAcqConstFiolta      float64 `json:"c_pay_acq_const_fiolta,omitempty"`      // float	Y	购建固定资产、无形资产和其他长期资产支付的现金
	CPaidInvest             float64 `json:"c_paid_invest,omitempty"`               // float	Y	投资支付的现金
	NDispSubsOthBiz         float64 `json:"n_disp_subs_oth_biz,omitempty"`         // float	Y	取得子公司及其他营业单位支付的现金净额
	OthPayRalInvAct         float64 `json:"oth_pay_ral_inv_act,omitempty"`         // float	Y	支付其他与投资活动有关的现金
	NIncrPledgeLoan         float64 `json:"n_incr_pledge_loan,omitempty"`          // float	Y	质押贷款净增加额
	StotOutInvAct           float64 `json:"stot_out_inv_act,omitempty"`            // float	Y	投资活动现金流出小计
	NCashflowInvAct         float64 `json:"n_cashflow_inv_act,omitempty"`          // float	Y	投资活动产生的现金流量净额
	CRecpBorrow             float64 `json:"c_recp_borrow,omitempty"`               // float	Y	取得借款收到的现金
	ProcIssueBonds          float64 `json:"proc_issue_bonds,omitempty"`            // float	Y	发行债券收到的现金
	OthCashRecpRalFncAct    float64 `json:"oth_cash_recp_ral_fnc_act,omitempty"`   // float	Y	收到其他与筹资活动有关的现金
	StotCashInFncAct        float64 `json:"stot_cash_in_fnc_act,omitempty"`        // float	Y	筹资活动现金流入小计
	FreeCashflow            float64 `json:"free_cashflow,omitempty"`               // float	Y	企业自由现金流量
	CPrepayAmtBorr          float64 `json:"c_prepay_amt_borr,omitempty"`           // float	Y	偿还债务支付的现金
	CPayDistDpcpIntExp      float64 `json:"c_pay_dist_dpcp_int_exp,omitempty"`     // float	Y	分配股利、利润或偿付利息支付的现金
	InclDvdProfitPaidScMs   float64 `json:"incl_dvd_profit_paid_sc_ms,omitempty"`  // float	Y	其中:子公司支付给少数股东的股利、利润
	OthCashpayRalFncAct     float64 `json:"oth_cashpay_ral_fnc_act,omitempty"`     // float	Y	支付其他与筹资活动有关的现金
	StotCashoutFncAct       float64 `json:"stot_cashout_fnc_act,omitempty"`        // float	Y	筹资活动现金流出小计
	NCashFlowsFncAct        float64 `json:"n_cash_flows_fnc_act,omitempty"`        // float	Y	筹资活动产生的现金流量净额
	EffFxFluCash            float64 `json:"eff_fx_flu_cash,omitempty"`             // float	Y	汇率变动对现金的影响
	NIncrCashCashEqu        float64 `json:"n_incr_cash_cash_equ,omitempty"`        // float	Y	现金及现金等价物净增加额
	CCashEquBegPeriod       float64 `json:"c_cash_equ_beg_period,omitempty"`       // float	Y	期初现金及现金等价物余额
	CCashEquEndPeriod       float64 `json:"c_cash_equ_end_period,omitempty"`       // float	Y	期末现金及现金等价物余额
	CRecpCapContrib         float64 `json:"c_recp_cap_contrib,omitempty"`          // float	Y	吸收投资收到的现金
	InclCashRecSaims        float64 `json:"incl_cash_rec_saims,omitempty"`         // float	Y	其中:子公司吸收少数股东投资收到的现金
	UnconInvestLoss         float64 `json:"uncon_invest_loss,omitempty"`           // float	Y	未确认投资损失
	ProvDeprAssets          float64 `json:"prov_depr_assets,omitempty"`            // float	Y	加:资产减值准备
	DeprFaCogaDpba          float64 `json:"depr_fa_coga_dpba,omitempty"`           // float	Y	固定资产折旧、油气资产折耗、生产性生物资产折旧
	AmortIntangAssets       float64 `json:"amort_intang_assets,omitempty"`         // float	Y	无形资产摊销
	LtAmortDeferredExp      float64 `json:"lt_amort_deferred_exp,omitempty"`       // float	Y	长期待摊费用摊销
	DecrDeferredExp         float64 `json:"decr_deferred_exp,omitempty"`           // float	Y	待摊费用减少
	IncrAccExp              float64 `json:"incr_acc_exp,omitempty"`                // float	Y	预提费用增加
	LossDispFiolta          float64 `json:"loss_disp_fiolta,omitempty"`            // float	Y	处置固定、无形资产和其他长期资产的损失
	LossScrFa               float64 `json:"loss_scr_fa,omitempty"`                 // float	Y	固定资产报废损失
	LossFvChg               float64 `json:"loss_fv_chg,omitempty"`                 // float	Y	公允价值变动损失
	InvestLoss              float64 `json:"invest_loss,omitempty"`                 // float	Y	投资损失
	DecrDefIncTaxAssets     float64 `json:"decr_def_inc_tax_assets,omitempty"`     // float	Y	递延所得税资产减少
	IncrDefIncTaxLiab       float64 `json:"incr_def_inc_tax_liab,omitempty"`       // float	Y	递延所得税负债增加
	DecrInventories         float64 `json:"decr_inventories,omitempty"`            // float	Y	存货的减少
	DecrOperPayable         float64 `json:"decr_oper_payable,omitempty"`           // float	Y	经营性应收项目的减少
	IncrOperPayable         float64 `json:"incr_oper_payable,omitempty"`           // float	Y	经营性应付项目的增加
	Others                  float64 `json:"others,omitempty"`                      // float	Y	其他
	ImNetCashflowOperAct    float64 `json:"im_net_cashflow_oper_act,omitempty"`    // float	Y	经营活动产生的现金流量净额(间接法)
	ConvDebtIntoCap         float64 `json:"conv_debt_into_cap,omitempty"`          // float	Y	债务转为资本
	ConvCopbondsDueWithin1y float64 `json:"conv_copbonds_due_within_1y,omitempty"` // float	Y	一年内到期的可转换公司债券
	FaFncLeases             float64 `json:"fa_fnc_leases,omitempty"`               // float	Y	融资租入固定资产
	EndBalCash              float64 `json:"end_bal_cash,omitempty"`                // float	Y	现金的期末余额
	BegBalCash              float64 `json:"beg_bal_cash,omitempty"`                // float	Y	减:现金的期初余额
	EndBalCashEqu           float64 `json:"end_bal_cash_equ,omitempty"`            // float	Y	加:现金等价物的期末余额
	BegBalCashEqu           float64 `json:"beg_bal_cash_equ,omitempty"`            // float	Y	减:现金等价物的期初余额
	ImNIncrCashEqu          float64 `json:"im_n_incr_cash_equ,omitempty"`          // float	Y	现金及现金等价物净增加额(间接法)
	UpdateFlag              string  `json:"update_flag,omitempty"`                 // str	N	更新标识
}

func (CashFlow) TableName() string {
	return "stk_cashflow"
}

func (CashFlow) KeyName() string {
	return entity.FieldName_Id
}

func (CashFlow) IdName() string {
	return entity.FieldName_Id
}
