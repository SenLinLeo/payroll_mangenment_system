package models

import (
	"encoding/hex"
	"github.com/astaxie/beego"
	"time"
)

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Page 	int    `json:"page"`
	Limit   int    `json:"limit"`
}

func GetEncData(data string) string {
	key := []byte(beego.AppConfig.String("secretKey"))
	return hex.EncodeToString(ThriDESEnCrypt([]byte(data), key))
}

func GetDecData(data string) string {
	key := []byte(beego.AppConfig.String("secretKey"))
	des, _ := hex.DecodeString(data)
	return string(ThriDESDeCrypt(des, key))
}

func (*TblCalcConf) TableName() string {
	return "tbl_calc_conf"
}

type TblCalcConf struct {
	Id       int    `orm:"column(id)" json:"Id"`
	ConfType string `orm:"column(conf_type)" json:"ConfType"` // 配置类型
	Key1     string `orm:"column(key1)" json:"Key1"`          // key1
	Key2     string `orm:"column(key2)" json:"Key2"`          // key2
	Value1   string `orm:"column(value1)" json:"Value1"`      // value1
	Value2   string `orm:"column(value2)" json:"Value2"`      // value2
	Resv1    string `orm:"column(resv1)" json:"Resv1"`        // 预留字段
}

func (*TblSalaryDtl) TableName() string {
	return "tbl_salary_dtl"
}

type TblSalaryDtl struct {
	Id                 int    `orm:"column(id)" json:"Id" match:"工号"`
	WorkerNo           string    `orm:"column(worker_no)" json:"WorkerNo" match:"姓名"`                           // 工号
	WorkerName         string `orm:"column(worker_name)" json:"WorkerName"`                                  // 员工姓名
	CalcMonth          string `orm:"column(calc_month)" json:"CalcMonth" match:"月份"`                         // 年月
	CalcWorkDays       int    `orm:"column(calc_work_days)" json:"CalcWorkDays"`                             // 工作日计算
	DaysCalcSalary     int    `orm:"column(days_calc_salary)" json:"DaysCalcSalary"`                         // 计薪天数
	DaysOvertime       int    `orm:"column(days_overtime)" json:"DaysOvertime" match:"加班天数"`                 // 加班天数
	DaysPrivateLeave   int    `orm:"column(days_private_leave)" json:"DaysPrivateLeave" match:"事假天数"`        // 事假天数
	DaysSickLeave      int    `orm:"column(days_sick_leave)" json:"DaysSickLeave" match:"病假天数"`              // 病假天数
	DaysEarlyLeave     int    `orm:"column(days_early_leave)" json:"DaysEarlyLeave"`                         // 迟到早退天数
	DaysStrikeLeave    int    `orm:"column(days_strike_leave)" json:"DaysStrikeLeave"`                       // 旷工天数
	DaysFoodSubsidy    int    `orm:"column(days_food_subsidy)" json:"DaysFoodSubsidy"`                       // 晚班餐补天数
	PrivateFund        float64    `orm:"column(private_fund)" json:"PrivateFund"`                                // 公积金个人
	PrivateSocial      float64    `orm:"column(private_social)" json:"PrivateSocial"`                            // 社保个人
	PrivateOld         float64    `orm:"column(private_old)" json:"PrivateOld"`                                  // 养老个人
	PrivateTreat       float64    `orm:"column(private_treat)" json:"PrivateTreat"`                              // 医疗个人
	PrivateUnemploy    float64    `orm:"column(private_unemploy)" json:"PrivateUnemploy"`                        // 失业个人
	PrivateSick        float64    `orm:"column(private_sick)" json:"PrivateSick"`                                // 大病个人
	TaxReduceChild     float64    `orm:"column(tax_reduce_child)" json:"TaxReduceChild"`                         // 累计子女教育支出扣除
	TaxReduceOld       float64    `orm:"column(tax_reduce_old)" json:"TaxReduceOld"`                             // 累计赡养老人支出扣除
	TaxReduceEducation float64    `orm:"column(tax_reduce_education)" json:"TaxReduceEducation"`                 // 累计继续教育支出扣除
	TaxReduceHouseLoan float64    `orm:"column(tax_reduce_house_loan)" json:"TaxReduceHouseLoan"`                // 累计住房贷款利息支出扣除
	TaxReduceHouseRent float64    `orm:"column(tax_reduce_house_rent)" json:"TaxReduceHouseRent"`                // 累计住房租金支出扣除
	CompanyFund        float64    `orm:"column(company_fund)" json:"CompanyFund"`                                // 公司公积金
	CompanySocial      float64    `orm:"column(company_social)" json:"CompanySocial"`                            // 公司社保
	ReduceBeforeTax    float64    `orm:"column(reduce_before_tax)" json:"ReduceBeforeTax" match:"税前扣款"`          // 税前其他扣款
	ReduceAftherTax    float64    `orm:"column(reduce_afther_tax)" json:"ReduceAftherTax"`                       // 其他税后扣除
	PayPosition        float64    `orm:"column(pay_position)" json:"PayPosition" match:"岗位工资"`                   // 定岗工资
	PayBonus           float64    `orm:"column(pay_bonus)" json:"PayBonus" match:"绩效奖金"`                         // 绩效奖金
	PayBeforeTax       float64    `orm:"column(pay_before_tax)" json:"PayBeforeTax" match:"税前补发"`                // 税前其他补发
	PayOtherAdd        float64    `orm:"column(pay_other_add)" json:"PayOtherAdd"`                               // 补偿金
	PayWait            float64    `orm:"column(pay_wait)" json:"PayWait"`                                        // 代通知金/
	PayFoodSubsidy     float64    `orm:"column(pay_food_subsidy)" json:"PayFoodSubsidy" match:"加班餐补"`            // 晚班补餐
	PayOvertime        float64    `orm:"column(pay_overtime)" json:"PayOvertime" match:"加班费"`                    // 加班费
	PayGeneral         float64    `orm:"column(pay_general)" json:"PayGeneral" match:"正常工作时间工资"`                 // 正常工作时间工资计算
	PayOther         float64    `orm:"column(pay_other)" json:"PayOther" match:"其他补发"`                 // 其他补发
	PayPositionAdd     float64    `orm:"column(pay_position_add)" json:"PayPositionAdd"`                         // 岗位津贴
	PaySecretAdd       float64    `orm:"column(pay_secret_add)" json:"PaySecretAdd" match:"保密工资"`                // 保密津贴
	ReducePrivateLeave float64    `orm:"column(reduce_private_leave)" json:"ReducePrivateLeave"`                 // 事假扣款
	ReduceSickLeave    float64    `orm:"column(reduce_sick_leave)" json:"ReduceSickLeave"`                       // 病假扣款
	ReduceEarlyLeave   float64    `orm:"column(reduce_early_leave)" json:"ReduceEarlyLeave"`                     // 迟到早退扣款
	ReduceStrikeLeave  float64    `orm:"column(reduce_strike_leave)" json:"ReduceStrikeLeave"`                   // 旷工扣款
	PayOtherAfterTax        float64    `orm:"column(pay_other_after_tax)" json:"PayOtherAfterTax" match:"其他不计税"`       // 其他不计税=税后补款 =补偿金+代通知金
	SumPayCalc         float64    `orm:"column(sum_pay_calc)" json:"SumPayCalc" `                                // 累计应发工资合计
	SumTaxFree         float64    `orm:"column(sum_tax_free)" json:"SumTaxFree"`                                 // 累计个税免征额
	SumFund            float64    `orm:"column(sum_fund)" json:"SumFund"`                                        // 公积金累计
	SumSocial          float64    `orm:"column(sum_social)" json:"SumSocial"`                                    // 社保累计
	SumTaxReduce       float64    `orm:"column(sum_tax_reduce)" json:"SumTaxReduce" match:"累计专项扣除"`              // 专项扣除累计
	SumTaxCalc         float64    `orm:"column(sum_tax_calc)" json:"SumTaxCalc"`                                 // 累计预扣预缴应纳税所得额
	SumTaxPrepare      float64    `orm:"column(sum_tax_prepare)" json:"SumTaxPrepare"`                           // 累计应预扣预缴税额（公式）
	SumTaxFact         float64    `orm:"column(sum_tax_fact)" json:"SumTaxFact"`                                 // 累计已预扣预缴税额
	SumTaxShould       float64    `orm:"column(sum_tax_should)" json:"SumTaxShould" match:"个税"`                  // 本期应预扣预缴税额(个税)
	SumPayFix          float64    `orm:"column(sum_pay_fix)" json:"SumPayFix"`                                   // 当月固定工资合计
	SumPayBeforeTax    float64    `orm:"column(sum_pay_before_tax)" json:"SumPayBeforeTax" match:"应发工资"`         // 税前应发合计
	SumFundSocial      float64    `orm:"column(sum_fund_social)" json:"SumFundSocial" match:"代扣社保（个人）代扣公积金（个人）"` // 公积金 个人 |社保 个人
	SumTaxBase         float64    `orm:"column(sum_tax_base)" json:"SumTaxBase" `                                // 应纳税所得额
	SumTaxNormal       float64    `orm:"column(sum_tax_normal)" json:"SumTaxNormal" `                            // 正常工资个税
	SumTaxOther        float64    `orm:"column(sum_tax_other)" json:"SumTaxOther"`                               // 非居民、劳务报酬 计算个税
	SumAfterTaxOther   float64    `orm:"column(sum_after_tax_other)" json:"SumAfterTaxOther" match:"税后扣款"`       // 其他税后扣除
	FactPay            float64    `orm:"column(fact_pay)" json:"FactPay" match:"实发工资"`                           // 实发工资
	SumCompanyFund     float64    `orm:"column(sum_company_fund)" json:"SumCompanyFund" match:"公司公积金"`           // 公司公积金
	SumCompanySocial   float64    `orm:"column(sum_company_social)" json:"SumCompanySocial" match:"公司社保"`        // 公司社保
	SumCompanyCost     float64    `orm:"column(sum_company_cost)" json:"SumCompanyCost"`                         // 公司人力成本合计
	IsShow             int `orm:"column(is_show)" json:"IsShow"`                         // 公司人力成本合计
	CreateTime         time.Time        `orm:"column(create_time)" json:"CreateTime"` // 创建时间
	UpdateTime         time.Time        `orm:"column(update_time)" json:"UpdateTime"` // 修改时间
}