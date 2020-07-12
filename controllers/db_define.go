package controllers


type SalaryController struct {
	BaseController
	CalcMonthStlmDays float64
}

type BaseResp struct {
	Code int           `json:"code"`
	Total int64           `json:"total"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
}

type QuerySumResp struct {
	Code int           `json:"code"`
	Total int64           `json:"total"`
	Msg  string        `json:"msg"`
	Data []QuerySumRespData `json:"data"`
}

type QuerySumRespData struct {
	CalcMonth         string   `json:"calc_month"`
	TotalPaypeaple   int  `json:"total_pay_peaple"`
	SumPayFix        float64  `json:"sum_pay_fix"`
	PayBonus          float64  `json:"pay_bonus"`
	SumPayBeforeTax float64  `json:"sum_pay_before_tax"`
	SumFund           float64  `json:"sum_fund"`
	SumSocial         float64  `json:"sum_social"`
	SumTaxShould     float64  `json:"sum_tax_should"`
	FactPay           float64  `json:"fact_pay"`
	SumCompanyFund   float64  `json:"sum_company_fund"`
	SumCompanySocial float64  `json:"sum_company_social"`
	SumCompanyCost    float64  `json:"sum_company_cost"`
}

type LoadCurrentSalary struct {
	WorkerNo         string `orm:"column(worker_no)" json:"WorkerNo" `                               // 工号
	WorkerName       string `orm:"column(worker_name)" json:"WorkerName"`                            // 员工姓名
	DaysCalcSalary   string   `orm:"column(days_calc_salary)" json:"DaysCalcSalary"`                         // 计薪天数
	DaysPrivateLeave string `orm:"column(days_private_leave)" json:"DaysPrivateLeave" match:"事假天数"`  // 事假天数
	DaysSickLeave    string `orm:"column(days_sick_leave)" json:"DaysSickLeave" match:"病假天数"`        // 病假天数
	DaysEarlyLeave   string   `orm:"column(days_early_leave)" json:"DaysEarlyLeave"`                         // 迟到早退天数
	DaysStrikeLeave  string   `orm:"column(days_strike_leave)" json:"DaysStrikeLeave"`                       // 旷工天数
	DaysOvertime     string   `orm:"column(days_overtime)" json:"DaysOvertime" match:"加班天数"`                 // 加班天数
	DaysFoodSubsidy  string `orm:"column(days_food_subsidy)" json:"DaysFoodSubsidy"`                 // 晚班餐补天数
}

type LoadHistorySalary struct {
	CalcMonth        string `orm:"column(calc_month)" json:"CalcMonth" `                             // 年月
	WorkerNo         string `orm:"column(worker_no)" json:"WorkerNo" `                               // 工号
	WorkerName       string `orm:"column(worker_name)" json:"WorkerName"`                            // 员工姓名
	PayGeneral       string `orm:"column(pay_general)" json:"PayGeneral" match:"正常工作时间工资"`           // 正常工作时间工资计算
	PayPositionAdd   string `orm:"column(pay_position_add)" json:"PayPositionAdd"`                   // 岗位津贴
	PaySecretAdd     string `orm:"column(pay_secret_add)" json:"PaySecretAdd" match:"保密工资"`          // 保密津贴
	DaysPrivateLeave string `orm:"column(days_private_leave)" json:"DaysPrivateLeave" match:"事假天数"`  // 事假天数
	DaysSickLeave    string `orm:"column(days_sick_leave)" json:"DaysSickLeave" match:"病假天数"`        // 病假天数
	DaysOvertime     string `orm:"column(days_overtime)" json:"DaysOvertime" match:"加班天数"`           // 加班天数
	PayOvertime      string `orm:"column(pay_overtime)" json:"PayOvertime" match:"加班费"`              // 加班费
	DaysFoodSubsidy  string `orm:"column(days_food_subsidy)" json:"DaysFoodSubsidy"`                 // 晚班餐补天数
	PayBonus         string `orm:"column(pay_bonus)" json:"PayBonus" match:"绩效奖金"`                   // 绩效奖金
	PayOther         string `orm:"column(pay_other)" json:"PayAfterTax" match:"其他补发"`                // 其他补发
	SumAfterTaxOther string `orm:"column(sum_after_tax_other)" json:"SumAfterTaxOther" match:"税后扣款"` // 其他税后扣除
	SumPayBeforeTax  string `orm:"column(sum_pay_before_tax)" json:"SumPayBeforeTax" match:"应发工资"`   // 税前应发合计
	SumCompanySocial string `orm:"column(sum_company_social)" json:"SumCompanySocial" match:"公司社保"`  // 公司社保
	SumCompanyFund   string `orm:"column(sum_company_fund)" json:"SumCompanyFund" match:"公司公积金"`     // 公司公积金
	SumFund          string `orm:"column(sum_fund)" json:"SumFund"`                                  // 公积金累计
	SumSocial        string `orm:"column(sum_social)" json:"SumSocial"`                              // 社保累计
	SumTaxReduce     string `orm:"column(sum_tax_reduce)" json:"SumTaxReduce" match:"累计专项扣除"`        // 专项扣除累计
	PayOtherAfterTax string `orm:"column(pay_other_after_tax)" json:"PayAfterTax" match:"其他不计税"`     // 其他不计税=税后补款 =补偿金+代通知金
	SumTaxShould     string `orm:"column(sum_tax_should)" json:"SumTaxShould" match:"个税"`            // 本期应预扣预缴税额(个税)
	FactPay          string `orm:"column(fact_pay)" json:"FactPay" match:"实发工资"`                     // 实发工资
}

type LoadAddSalary struct {
	WorkerNo         string `orm:"column(worker_no)" json:"WorkerNo" `                               // 工号
	WorkerName       string `orm:"column(worker_name)" json:"WorkerName"`                            // 员工姓名
	ReduceBeforeTax  string   `orm:"column(reduce_before_tax)" json:"ReduceBeforeTax" match:"税前扣款"`          // 税前其他扣款
	PayBonus         string `orm:"column(pay_bonus)" json:"PayBonus" match:"绩效奖金"`                   // 绩效奖金
	PayBeforeTax     string   `orm:"column(pay_before_tax)" json:"PayBeforeTax" match:"税前补发"`                // 税前其他补发
	SumAfterTaxOther string `orm:"column(sum_after_tax_other)" json:"SumAfterTaxOther" match:"税后扣款"` // 其他税后扣除
	PayOtherAdd      string   `orm:"column(pay_other_add)" json:"PayOtherAdd"`                               // 补偿金
	PayWait          string   `orm:"column(pay_wait)" json:"PayWait"`                                        // 代通知金/
}

type LoadfundSalary struct {
	WorkerNo           string `orm:"column(worker_no)" json:"WorkerNo" `                               // 工号
	WorkerName         string `orm:"column(worker_name)" json:"WorkerName"`                            // 员工姓名
	PrivateFund        string    `orm:"column(private_fund)" json:"PrivateFund"`                                // 公积金个人
	PrivateSocial      string    `orm:"column(private_social)" json:"PrivateSocial"`                            // 社保个人
	PrivateOld         string    `orm:"column(private_old)" json:"PrivateOld"`                                  // 养老个人
	PrivateTreat       string    `orm:"column(private_treat)" json:"PrivateTreat"`                              // 医疗个人
	PrivateUnemploy    string    `orm:"column(private_unemploy)" json:"PrivateUnemploy"`                        // 失业个人
	PrivateSick        string    `orm:"column(private_sick)" json:"PrivateSick"`                                // 大病个人
	TaxReduceChild     string    `orm:"column(tax_reduce_child)" json:"TaxReduceChild"`                         // 累计子女教育支出扣除
	TaxReduceOld       string    `orm:"column(tax_reduce_old)" json:"TaxReduceOld"`                             // 累计赡养老人支出扣除
	TaxReduceEducation string    `orm:"column(tax_reduce_education)" json:"TaxReduceEducation"`                 // 累计继续教育支出扣除
	TaxReduceHouseLoan string    `orm:"column(tax_reduce_house_loan)" json:"TaxReduceHouseLoan"`                // 累计住房贷款利息支出扣除
	TaxReduceHouseRent string    `orm:"column(tax_reduce_house_rent)" json:"TaxReduceHouseRent"`                // 累计住房租金支出扣除
	CompanyFund        string    `orm:"column(company_fund)" json:"CompanyFund"`                                // 公司公积金
	CompanySocial      string    `orm:"column(company_social)" json:"CompanySocial"`                            // 公司社保
}

type LoadTotalSalary struct {
	WorkerNo           string     `orm:"column(worker_no)" json:"WorkerNo" `                                        // 工号
	WorkerName         string `orm:"column(worker_name)" json:"WorkerName"`                                     // 员工姓名
	SumPayCalc         string     `orm:"column(sum_pay_calc)" json:"SumPayCalc" `                                // 累计应发工资合计
	SumFund            string     `orm:"column(sum_fund)" json:"SumFund"`                                        // 公积金累计
	SumSocial          string     `orm:"column(sum_social)" json:"SumSocial"`                                    // 社保累计
	SumTaxFree         string     `orm:"column(sum_tax_free)" json:"SumTaxFree"`                                 // 累计个税免征额
	SumTaxReduce       string     `orm:"column(sum_tax_reduce)" json:"SumTaxReduce" match:"累计专项扣除"`              // 专项扣除累计
	SumTaxFact         string     `orm:"column(sum_tax_fact)" json:"SumTaxFact"`                                 // 累计已预扣预缴税额
}