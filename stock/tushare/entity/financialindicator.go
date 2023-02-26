package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

type FinancialIndicator struct {
	entity.BaseEntity         `xorm:"extends"`
	TsCode                    string  `xorm:"varchar(255)" json:"ts_code,omitempty"`  // str	Y	TS代码
	AnnDate                   string  `json:"ann_date,omitempty"`                     //	str	Y	公告日期
	EndDate                   string  `json:"end_date,omitempty"`                     //	str	Y	报告期
	Eps                       float64 `json:"eps,omitempty"`                          //	float	Y	基本每股收益
	DtEps                     float64 `json:"dt_eps,omitempty"`                       //	float	Y	稀释每股收益
	TotalRevenuePs            float64 `json:"total_revenue_ps,omitempty"`             //	float	Y	每股营业总收入
	RevenuePs                 float64 `json:"revenue_ps,omitempty"`                   //	float	Y	每股营业收入
	CapitalResePs             float64 `json:"capital_rese_ps,omitempty"`              //	float	Y	每股资本公积
	SurplusResePs             float64 `json:"surplus_rese_ps,omitempty"`              //	float	Y	每股盈余公积
	UndistProfitPs            float64 `json:"undist_profit_ps,omitempty"`             //	float	Y	每股未分配利润
	ExtraItem                 float64 `json:"extra_item,omitempty"`                   //	float	Y	非经常性损益
	ProfitDedt                float64 `json:"profit_dedt,omitempty"`                  //	float	Y	扣除非经常性损益后的净利润
	GrossMargin               float64 `json:"gross_margin,omitempty"`                 //	float	Y	毛利
	CurrentRatio              float64 `json:"current_ratio,omitempty"`                //	float	Y	流动比率
	QuickRatio                float64 `json:"quick_ratio,omitempty"`                  //	float	Y	速动比率
	CashRatio                 float64 `json:"cash_ratio,omitempty"`                   //	float	Y	保守速动比率
	InvturnDays               float64 `json:"invturn_days,omitempty"`                 //	float	N	存货周转天数
	ArturnDays                float64 `json:"arturn_days,omitempty"`                  //	float	N	应收账款周转天数
	InvTurn                   float64 `json:"inv_turn,omitempty"`                     //	float	N	存货周转率
	ArTurn                    float64 `json:"ar_turn,omitempty"`                      //	float	Y	应收账款周转率
	CaTurn                    float64 `json:"ca_turn,omitempty"`                      //	float	Y	流动资产周转率
	FaTurn                    float64 `json:"fa_turn,omitempty"`                      //	float	Y	固定资产周转率
	AssetsTurn                float64 `json:"assets_turn,omitempty"`                  //	float	Y	总资产周转率
	OpIncome                  float64 `json:"op_income,omitempty"`                    //	float	Y	经营活动净收益
	ValuechangeIncome         float64 `json:"valuechange_income,omitempty"`           //	float	N	价值变动净收益
	InterstIncome             float64 `json:"interst_income,omitempty"`               //	float	N	利息费用
	Daa                       float64 `json:"daa,omitempty"`                          //	float	N	折旧与摊销
	Ebit                      float64 `json:"ebit,omitempty"`                         //	float	Y	息税前利润
	Ebitda                    float64 `json:"ebitda,omitempty"`                       //	float	Y	息税折旧摊销前利润
	Fcff                      float64 `json:"fcff,omitempty"`                         //	float	Y	企业自由现金流量
	Fcfe                      float64 `json:"fcfe,omitempty"`                         //	float	Y	股权自由现金流量
	CurrentExint              float64 `json:"current_exint,omitempty"`                //	float	Y	无息流动负债
	NoncurrentExint           float64 `json:"noncurrent_exint,omitempty"`             //	float	Y	无息非流动负债
	Interestdebt              float64 `json:"interestdebt,omitempty"`                 //	float	Y	带息债务
	Netdebt                   float64 `json:"netdebt,omitempty"`                      //	float	Y	净债务
	TangibleAsset             float64 `json:"tangible_asset,omitempty"`               //	float	Y	有形资产
	WorkingCapital            float64 `json:"working_capital,omitempty"`              //	float	Y	营运资金
	NetworkingCapital         float64 `json:"networking_capital,omitempty"`           //	float	Y	营运流动资本
	InvestCapital             float64 `json:"invest_capital,omitempty"`               //	float	Y	全部投入资本
	RetainedEarnings          float64 `json:"retained_earnings,omitempty"`            //	float	Y	留存收益
	Diluted2Eps               float64 `json:"diluted2_eps,omitempty"`                 //	float	Y	期末摊薄每股收益
	Bps                       float64 `json:"bps,omitempty"`                          //	float	Y	每股净资产
	Ocfps                     float64 `json:"ocfps,omitempty"`                        //	float	Y	每股经营活动产生的现金流量净额
	Retainedps                float64 `json:"retainedps,omitempty"`                   //	float	Y	每股留存收益
	Cfps                      float64 `json:"cfps,omitempty"`                         //	float	Y	每股现金流量净额
	EbitPs                    float64 `json:"ebit_ps,omitempty"`                      //	float	Y	每股息税前利润
	FcffPs                    float64 `json:"fcff_ps,omitempty"`                      //	float	Y	每股企业自由现金流量
	FcfePs                    float64 `json:"fcfe_ps,omitempty"`                      //	float	Y	每股股东自由现金流量
	NetprofitMargin           float64 `json:"netprofit_margin,omitempty"`             //	float	Y	销售净利率
	GrossprofitMargin         float64 `json:"grossprofit_margin,omitempty"`           //	float	Y	销售毛利率
	CogsOfSales               float64 `json:"cogs_of_sales,omitempty"`                //	float	Y	销售成本率
	ExpenseOfSales            float64 `json:"expense_of_sales,omitempty"`             //	float	Y	销售期间费用率
	ProfitToGr                float64 `json:"profit_to_gr,omitempty"`                 //	float	Y	净利润/营业总收入
	SaleexpToGr               float64 `json:"saleexp_to_gr,omitempty"`                //	float	Y	销售费用/营业总收入
	AdminexpOfGr              float64 `json:"adminexp_of_gr,omitempty"`               //	float	Y	管理费用/营业总收入
	FinaexpOfGr               float64 `json:"finaexp_of_gr,omitempty"`                //	float	Y	财务费用/营业总收入
	ImpaiTtm                  float64 `json:"impai_ttm,omitempty"`                    //	float	Y	资产减值损失/营业总收入
	GcOfGr                    float64 `json:"gc_of_gr,omitempty"`                     //	float	Y	营业总成本/营业总收入
	OpOfGr                    float64 `json:"op_of_gr,omitempty"`                     //	float	Y	营业利润/营业总收入
	EbitOfGr                  float64 `json:"ebit_of_gr,omitempty"`                   //	float	Y	息税前利润/营业总收入
	Roe                       float64 `json:"roe,omitempty"`                          //	float	Y	净资产收益率
	RoeWaa                    float64 `json:"roe_waa,omitempty"`                      //	float	Y	加权平均净资产收益率
	RoeDt                     float64 `json:"roe_dt,omitempty"`                       //	float	Y	净资产收益率(扣除非经常损益)
	Roa                       float64 `json:"roa,omitempty"`                          //	float	Y	总资产报酬率
	Npta                      float64 `json:"npta,omitempty"`                         //	float	Y	总资产净利润
	Roic                      float64 `json:"roic,omitempty"`                         //	float	Y	投入资本回报率
	RoeYearly                 float64 `json:"roe_yearly,omitempty"`                   //	float	Y	年化净资产收益率
	Roa2Yearly                float64 `json:"roa2_yearly,omitempty"`                  //	float	Y	年化总资产报酬率
	RoeAvg                    float64 `json:"roe_avg,omitempty"`                      //	float	N	平均净资产收益率(增发条件)
	OpincomeOfEbt             float64 `json:"opincome_of_ebt,omitempty"`              //	float	N	经营活动净收益/利润总额
	InvestincomeOfEbt         float64 `json:"investincome_of_ebt,omitempty"`          //	float	N	价值变动净收益/利润总额
	NOpProfitOfEbt            float64 `json:"n_op_profit_of_ebt,omitempty"`           //	float	N	营业外收支净额/利润总额
	TaxToEbt                  float64 `json:"tax_to_ebt,omitempty"`                   //	float	N	所得税/利润总额
	DtprofitToProfit          float64 `json:"dtprofit_to_profit,omitempty"`           //	float	N	扣除非经常损益后的净利润/净利润
	SalescashToOr             float64 `json:"salescash_to_or,omitempty"`              //	float	N	销售商品提供劳务收到的现金/营业收入
	OcfToOr                   float64 `json:"ocf_to_or,omitempty"`                    //	float	N	经营活动产生的现金流量净额/营业收入
	OcfToOpincome             float64 `json:"ocf_to_opincome,omitempty"`              //	float	N	经营活动产生的现金流量净额/经营活动净收益
	CapitalizedToDa           float64 `json:"capitalized_to_da,omitempty"`            //	float	N	资本支出/折旧和摊销
	DebtToAssets              float64 `json:"debt_to_assets,omitempty"`               //	float	Y	资产负债率
	AssetsToEqt               float64 `json:"assets_to_eqt,omitempty"`                //	float	Y	权益乘数
	DpAssetsToEqt             float64 `json:"dp_assets_to_eqt,omitempty"`             //	float	Y	权益乘数(杜邦分析)
	CaToAssets                float64 `json:"ca_to_assets,omitempty"`                 //	float	Y	流动资产/总资产
	NcaToAssets               float64 `json:"nca_to_assets,omitempty"`                //	float	Y	非流动资产/总资产
	TbassetsToTotalassets     float64 `json:"tbassets_to_totalassets,omitempty"`      //	float	Y	有形资产/总资产
	IntToTalcap               float64 `json:"int_to_talcap,omitempty"`                //	float	Y	带息债务/全部投入资本
	EqtToTalcapital           float64 `json:"eqt_to_talcapital,omitempty"`            //	float	Y	归属于母公司的股东权益/全部投入资本
	CurrentdebtToDebt         float64 `json:"currentdebt_to_debt,omitempty"`          //	float	Y	流动负债/负债合计
	LongdebToDebt             float64 `json:"longdeb_to_debt,omitempty"`              //	float	Y	非流动负债/负债合计
	OcfToShortdebt            float64 `json:"ocf_to_shortdebt,omitempty"`             //	float	Y	经营活动产生的现金流量净额/流动负债
	DebtToEqt                 float64 `json:"debt_to_eqt,omitempty"`                  //	float	Y	产权比率
	EqtToDebt                 float64 `json:"eqt_to_debt,omitempty"`                  //	float	Y	归属于母公司的股东权益/负债合计
	EqtToInterestdebt         float64 `json:"eqt_to_interestdebt,omitempty"`          //	float	Y	归属于母公司的股东权益/带息债务
	TangibleassetToDebt       float64 `json:"tangibleasset_to_debt,omitempty"`        //	float	Y	有形资产/负债合计
	TangassetToIntdebt        float64 `json:"tangasset_to_intdebt,omitempty"`         //	float	Y	有形资产/带息债务
	TangibleassetToNetdebt    float64 `json:"tangibleasset_to_netdebt,omitempty"`     //	float	Y	有形资产/净债务
	OcfToDebt                 float64 `json:"ocf_to_debt,omitempty"`                  //	float	Y	经营活动产生的现金流量净额/负债合计
	OcfToInterestdebt         float64 `json:"ocf_to_interestdebt,omitempty"`          //	float	N	经营活动产生的现金流量净额/带息债务
	OcfToNetdebt              float64 `json:"ocf_to_netdebt,omitempty"`               //	float	N	经营活动产生的现金流量净额/净债务
	EbitToInterest            float64 `json:"ebit_to_interest,omitempty"`             //	float	N	已获利息倍数(EBIT/利息费用)
	LongdebtToWorkingcapital  float64 `json:"longdebt_to_workingcapital,omitempty"`   //	float	N	长期债务与营运资金比率
	EbitdaToDebt              float64 `json:"ebitda_to_debt,omitempty"`               //	float	N	息税折旧摊销前利润/负债合计
	TurnDays                  float64 `json:"turn_days,omitempty"`                    //	float	Y	营业周期
	RoaYearly                 float64 `json:"roa_yearly,omitempty"`                   //	float	Y	年化总资产净利率
	RoaDp                     float64 `json:"roa_dp,omitempty"`                       //	float	Y	总资产净利率(杜邦分析)
	FixedAssets               float64 `json:"fixed_assets,omitempty"`                 //	float	Y	固定资产合计
	ProfitPrefinExp           float64 `json:"profit_prefin_exp,omitempty"`            //	float	N	扣除财务费用前营业利润
	NonOpProfit               float64 `json:"non_op_profit,omitempty"`                //	float	N	非营业利润
	OpToEbt                   float64 `json:"op_to_ebt,omitempty"`                    //	float	N	营业利润／利润总额
	NopToEbt                  float64 `json:"nop_to_ebt,omitempty"`                   //	float	N	非营业利润／利润总额
	OcfToProfit               float64 `json:"ocf_to_profit,omitempty"`                //	float	N	经营活动产生的现金流量净额／营业利润
	CashToLiqdebt             float64 `json:"cash_to_liqdebt,omitempty"`              //	float	N	货币资金／流动负债
	CashToLiqdebtWithinterest float64 `json:"cash_to_liqdebt_withinterest,omitempty"` //	float	N	货币资金／带息流动负债
	OpToLiqdebt               float64 `json:"op_to_liqdebt,omitempty"`                //	float	N	营业利润／流动负债
	OpToDebt                  float64 `json:"op_to_debt,omitempty"`                   //	float	N	营业利润／负债合计
	RoicYearly                float64 `json:"roic_yearly,omitempty"`                  //	float	N	年化投入资本回报率
	TotalFaTrun               float64 `json:"total_fa_trun,omitempty"`                //	float	N	固定资产合计周转率
	ProfitToOp                float64 `json:"profit_to_op,omitempty"`                 //	float	Y	利润总额／营业收入
	QOpincome                 float64 `json:"q_opincome,omitempty"`                   //	float	N	经营活动单季度净收益
	QInvestincome             float64 `json:"q_investincome,omitempty"`               //	float	N	价值变动单季度净收益
	QDtprofit                 float64 `json:"q_dtprofit,omitempty"`                   //	float	N	扣除非经常损益后的单季度净利润
	QEps                      float64 `json:"q_eps,omitempty"`                        //	float	N	每股收益(单季度)
	QNetprofitMargin          float64 `json:"q_netprofit_margin,omitempty"`           //	float	N	销售净利率(单季度)
	QGsprofitMargin           float64 `json:"q_gsprofit_margin,omitempty"`            //	float	N	销售毛利率(单季度)
	QExpToSales               float64 `json:"q_exp_to_sales,omitempty"`               //	float	N	销售期间费用率(单季度)
	QProfitToGr               float64 `json:"q_profit_to_gr,omitempty"`               //	float	N	净利润／营业总收入(单季度)
	QSaleexpToGr              float64 `json:"q_saleexp_to_gr,omitempty"`              //	float	Y	销售费用／营业总收入 (单季度)
	QAdminexpToGr             float64 `json:"q_adminexp_to_gr,omitempty"`             //	float	N	管理费用／营业总收入 (单季度)
	QFinaexpToGr              float64 `json:"q_finaexp_to_gr,omitempty"`              //	float	N	财务费用／营业总收入 (单季度)
	QImpairToGrTtm            float64 `json:"q_impair_to_gr_ttm,omitempty"`           //	float	N	资产减值损失／营业总收入(单季度)
	QGcToGr                   float64 `json:"q_gc_to_gr,omitempty"`                   //	float	Y	营业总成本／营业总收入 (单季度)
	QOpToGr                   float64 `json:"q_op_to_gr,omitempty"`                   //	float	N	营业利润／营业总收入(单季度)
	QRoe                      float64 `json:"q_roe,omitempty"`                        //	float	Y	净资产收益率(单季度)
	QDtRoe                    float64 `json:"q_dt_roe,omitempty"`                     //	float	Y	净资产单季度收益率(扣除非经常损益)
	QNpta                     float64 `json:"q_npta,omitempty"`                       //	float	Y	总资产净利润(单季度)
	QOpincomeToEbt            float64 `json:"q_opincome_to_ebt,omitempty"`            //	float	N	经营活动净收益／利润总额(单季度)
	QInvestincomeToEbt        float64 `json:"q_investincome_to_ebt,omitempty"`        //	float	N	价值变动净收益／利润总额(单季度)
	QDtprofitToProfit         float64 `json:"q_dtprofit_to_profit,omitempty"`         //	float	N	扣除非经常损益后的净利润／净利润(单季度)
	QSalescashToOr            float64 `json:"q_salescash_to_or,omitempty"`            //	float	N	销售商品提供劳务收到的现金／营业收入(单季度)
	QOcfToSales               float64 `json:"q_ocf_to_sales,omitempty"`               //	float	Y	经营活动产生的现金流量净额／营业收入(单季度)
	QOcfToOr                  float64 `json:"q_ocf_to_or,omitempty"`                  //	float	N	经营活动产生的现金流量净额／经营活动净收益(单季度)
	BasicEpsYoy               float64 `json:"basic_eps_yoy,omitempty"`                //	float	Y	基本每股收益同比增长率(%)
	DtEpsYoy                  float64 `json:"dt_eps_yoy,omitempty"`                   //	float	Y	稀释每股收益同比增长率(%)
	CfpsYoy                   float64 `json:"cfps_yoy,omitempty"`                     //	float	Y	每股经营活动产生的现金流量净额同比增长率(%)
	OpYoy                     float64 `json:"op_yoy,omitempty"`                       //	float	Y	营业利润同比增长率(%)
	EbtYoy                    float64 `json:"ebt_yoy,omitempty"`                      //	float	Y	利润总额同比增长率(%)
	NetprofitYoy              float64 `json:"netprofit_yoy,omitempty"`                //	float	Y	归属母公司股东的净利润同比增长率(%)
	DtNetprofitYoy            float64 `json:"dt_netprofit_yoy,omitempty"`             //	float	Y	归属母公司股东的净利润-扣除非经常损益同比增长率(%)
	OcfYoy                    float64 `json:"ocf_yoy,omitempty"`                      //	float	Y	经营活动产生的现金流量净额同比增长率(%)
	RoeYoy                    float64 `json:"roe_yoy,omitempty"`                      //	float	Y	净资产收益率(摊薄)同比增长率(%)
	BpsYoy                    float64 `json:"bps_yoy,omitempty"`                      //	float	Y	每股净资产相对年初增长率(%)
	AssetsYoy                 float64 `json:"assets_yoy,omitempty"`                   //	float	Y	资产总计相对年初增长率(%)
	EqtYoy                    float64 `json:"eqt_yoy,omitempty"`                      //	float	Y	归属母公司的股东权益相对年初增长率(%)
	TrYoy                     float64 `json:"tr_yoy,omitempty"`                       //	float	Y	营业总收入同比增长率(%)
	OrYoy                     float64 `json:"or_yoy,omitempty"`                       //	float	Y	营业收入同比增长率(%)
	QGrYoy                    float64 `json:"q_gr_yoy,omitempty"`                     //	float	N	营业总收入同比增长率(%)(单季度)
	QGrQoq                    float64 `json:"q_gr_qoq,omitempty"`                     //	float	N	营业总收入环比增长率(%)(单季度)
	QSalesYoy                 float64 `json:"q_sales_yoy,omitempty"`                  //	float	Y	营业收入同比增长率(%)(单季度)
	QSalesQoq                 float64 `json:"q_sales_qoq,omitempty"`                  //	float	N	营业收入环比增长率(%)(单季度)
	QOpYoy                    float64 `json:"q_op_yoy,omitempty"`                     //	float	N	营业利润同比增长率(%)(单季度)
	QOpQoq                    float64 `json:"q_op_qoq,omitempty"`                     //	float	Y	营业利润环比增长率(%)(单季度)
	QProfitYoy                float64 `json:"q_profit_yoy,omitempty"`                 //	float	N	净利润同比增长率(%)(单季度)
	QProfitQoq                float64 `json:"q_profit_qoq,omitempty"`                 //	float	N	净利润环比增长率(%)(单季度)
	QNetprofitYoy             float64 `json:"q_netprofit_yoy,omitempty"`              //	float	N	归属母公司股东的净利润同比增长率(%)(单季度)
	QNetprofitQoq             float64 `json:"q_netprofit_qoq,omitempty"`              //	float	N	归属母公司股东的净利润环比增长率(%)(单季度)
	EquityYoy                 float64 `json:"equity_yoy,omitempty"`                   //	float	Y	净资产同比增长率
	RdExp                     float64 `json:"rd_exp,omitempty"`                       //	float	N	研发费用
	UpdateFlag                string  `json:"update_flag,omitempty"`                  //	str	N	更新标识
}

func (FinancialIndicator) TableName() string {
	return "stk_financialindicator"
}

func (FinancialIndicator) KeyName() string {
	return entity.FieldName_Id
}

func (FinancialIndicator) IdName() string {
	return entity.FieldName_Id
}
