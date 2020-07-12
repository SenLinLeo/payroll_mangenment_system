package controllers

import (
	"errors"
	"fmt"
	"geek-nebula/models"
	"geek-nebula/utils"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	"strconv"
)

func (self *SalaryController) CalcPriveSalary(tblSalaryDtl *models.TblSalaryDtl) (err error) {
	var (
		o               = orm.NewOrm()
		tblEmployeeInfo = models.TblEmployeeInfo{}
		tmpNum          float64
	)

	defer o.Commit()

	logger.InfoById(tblSalaryDtl.Id, "当前计算 [%s][%s]工资", tblSalaryDtl.CalcMonth, tblSalaryDtl.WorkerName)

	/*
		定岗工资 = 正常工作时间工资（基数表算得） + 岗位津贴 + 保密津贴
		pay_position = pay_general + pay_position_add + pay_secret_add
	*/
	var CurrentSalarystr string

	if err = o.Raw(`
	select 
		id,
		worker_no,
		worker_name,
		sex,
		mobile,
		en_mobile,
		idcard_type,
		idcard_no,
		en_idcard_no,
		birthday,
		bcard_no,
		en_bcard_no,
		bank_name,
		company,
		employee_type,
		work_status,
		is_regular_employee,
		department,
		jobs,
		tax_type,
		tax_payer,
		work_city,
		tax_city,
		email,
		  DATE_FORMAT(get_in_date,'%Y%m%d%H%i%s') ,
		get_in_salary,
		  DATE_FORMAT(change_formal_date,'%Y%m%d%H%i%s'),
		change_formal_salary,
		  DATE_FORMAT(get_out_date,'%Y%m%d%H%i%s'),
		practicec_salary,
		trial_salary,
		current_salary,
		pay_month,
		  DATE_FORMAT(change_salary_date,'%Y%m%d%H%i%s'),
		  DATE_FORMAT(change_job_date,'%Y%m%d%H%i%s'),
		  DATE_FORMAT(change_company_date,'%Y%m%d%H%i%s'),
		create_time,
		  DATE_FORMAT(update_time,'%Y%m%d%H%i%s'),
		import_time,
		resv
	from tbl_employee_info 
	where worker_no = ?
	`, tblSalaryDtl.WorkerNo).QueryRow(&tblEmployeeInfo); err != nil {
		logger.Errorf("%s", err.Error())
		return errors.New(tblSalaryDtl.WorkerNo + " 查询不到此员工")
	}

	logger.Infof("%s", tblEmployeeInfo.CurrentSalary)

	if tblEmployeeInfo.CurrentSalary == "0" {
		logger.Errorf("%s 定岗工资异常", tblSalaryDtl.WorkerNo)
		return errors.New(fmt.Sprintf("%s 定岗工资异常", tblSalaryDtl.WorkerNo))
	}

	CurrentSalarystr = tblEmployeeInfo.CurrentSalary
	CurrentSalarystr = models.GetDecData(CurrentSalarystr)

	tmpNum, _ = strconv.ParseFloat(CurrentSalarystr, 64)
	// 定岗工资
	tblSalaryDtl.PayPosition = tmpNum
	logger.Infof("定岗工资[%f]", tblSalaryDtl.PayPosition)


	// 正常工作时间工资
	if err := o.Raw(`
    SELECT value1+0
      FROM tbl_calc_conf 
     WHERE key1 + 0 <= ?
       AND key2 + 0 >= ?
       AND conf_type = "base_salary"`,
		CurrentSalarystr, CurrentSalarystr).QueryRow(&tblSalaryDtl.PayGeneral); err != nil {
		logger.Errorf("%s", err.Error())
		return errors.New(fmt.Sprintf("%s 查询配置异常", tblSalaryDtl.WorkerNo))
	}

	// 保密津贴
	tblSalaryDtl.PaySecretAdd = 200
	// 岗位津贴
	tblSalaryDtl.PayPositionAdd = tblSalaryDtl.PayPosition - tblSalaryDtl.PayGeneral - tblSalaryDtl.PaySecretAdd

	logger.Infof("正常工作时间工资[%f]", tblSalaryDtl.PayGeneral)
	logger.Infof("保密津贴[%f]", tblSalaryDtl.PaySecretAdd)
	logger.Infof("岗位津贴[%f]", tblSalaryDtl.PayPositionAdd)

	//工作日计算(只需算一次)
	//calc_work_days
	if self.CalcMonthStlmDays == 0 {
		if err = o.Raw(`
        select count(distinct stlm_date)
          from tbl_stlm_date 
         where substr(stlm_date,1,6) = ?
		`, tblSalaryDtl.CalcMonth).
			QueryRow(&self.CalcMonthStlmDays); err != nil {
			logger.Errorf("%s", err.Error())
			return errors.New(fmt.Sprintf("%s 查询清算日期异常", tblSalaryDtl.WorkerNo))
		}
	}

	logger.Infof("工作日计算[%f]", self.CalcMonthStlmDays)

	//当月固定工资合计
	//sum_pay_fix
	emyType := tblEmployeeInfo.EmployeeType
	trialSalary, _ := strconv.ParseFloat(models.GetDecData(tblEmployeeInfo.TrialSalary), 64)
	changeFormalSalary, _ := strconv.ParseFloat(models.GetDecData(tblEmployeeInfo.ChangeFormalSalary), 64)
	switch (true) {
	//  3-正式 2-试用
	case emyType == "2" || emyType == "3":
		{
			tblSalaryDtl.SumPayFix = tblSalaryDtl.PayPosition / 21.75 * float64(tblSalaryDtl.DaysCalcSalary)
		}

		// 当月调薪
	case len(tblEmployeeInfo.ChangeFormalDate) > 0 && tblEmployeeInfo.ChangeSalaryDate[:6] == tblSalaryDtl.CalcMonth:
		{
			tblSalaryDtl.SumPayFix = tblSalaryDtl.PayPosition / 21.75 * float64(tblSalaryDtl.DaysCalcSalary)
		}

		// 当月转正
	case len(tblEmployeeInfo.ChangeFormalDate) > 0 && tblEmployeeInfo.ChangeFormalDate[:6] == tblSalaryDtl.CalcMonth:
		{
			var daysBefore float64
			if err = o.Raw(`
                select count(distinct stlm_date)
                  from tbl_stlm_date 
                 where substr(stlm_date,1,6) = ?
                   and stlm_date < ?
				`, tblSalaryDtl.CalcMonth, tblEmployeeInfo.ChangeFormalDate).
				QueryRow(&daysBefore); err != nil {
				logger.Errorf("%s", err.Error())
				return errors.New(fmt.Sprintf("%s 查询清算日期异常", tblSalaryDtl.WorkerNo))
			}
			tblSalaryDtl.SumPayFix = trialSalary/21.75*float64(daysBefore) + changeFormalSalary/21.75*(self.CalcMonthStlmDays-daysBefore)

			if len(tblEmployeeInfo.GetOutDate) > 0 && tblEmployeeInfo.GetOutDate[:6] != tblSalaryDtl.CalcMonth {
				break
			}

			// 当月转正并离职
			// 转正前离职
			if utils.CalcDateSeek(tblEmployeeInfo.ChangeFormalDate, tblEmployeeInfo.GetOutDate) > 0 {
				tblSalaryDtl.SumPayFix = trialSalary / 21.75 * float64(tblSalaryDtl.DaysCalcSalary)
			}
		}
	default:
		logger.Errorf("%s 获取当月固定工资异常", tblSalaryDtl.WorkerNo)
		return errors.New(fmt.Sprintf("%s 获取当月固定工资异常", tblSalaryDtl.WorkerNo))
	}

	logger.Infof("当月固定工资 [%f]", tblSalaryDtl.SumPayFix)

	//加班费 = 正常工作时间工资/21.75 *加班天数
	tblSalaryDtl.PayOvertime = tblSalaryDtl.PayGeneral / 21.75 * float64(tblSalaryDtl.DaysOvertime)

	//晚餐餐补 = 补餐次数 *15
	tblSalaryDtl.PayFoodSubsidy = float64(tblSalaryDtl.DaysFoodSubsidy * 15)

	//事假扣款 = 当前工资/21.75 * 事假天数
	tblSalaryDtl.ReducePrivateLeave = tblSalaryDtl.PayPosition / 21.75 * float64(tblSalaryDtl.DaysPrivateLeave)

	//病假扣款 = 当前工资/21.75 * 病假天数 *40%
	tblSalaryDtl.ReduceSickLeave = tblSalaryDtl.PayPosition / 21.75 * float64(tblSalaryDtl.DaysSickLeave) * 0.4

	//	旷工扣款 = 当前工资/21.75 * 旷工天数
	tblSalaryDtl.ReduceSickLeave = tblSalaryDtl.PayPosition / 21.75 * float64(tblSalaryDtl.DaysSickLeave)

	/*  税前应发合计 = 当月固定工资合计 + 加班费 + 晚班补餐 + 绩效奖金 + 其他补发
	    事假扣款 + 病假扣款 + 迟到早退扣款 + 旷工扣款 + 其他扣款  */
	tblSalaryDtl.SumPayBeforeTax = tblSalaryDtl.SumPayFix + tblSalaryDtl.PayOvertime + tblSalaryDtl.PayFoodSubsidy + tblSalaryDtl.PayBonus + tblSalaryDtl.PayOther -
		tblSalaryDtl.ReducePrivateLeave - tblSalaryDtl.ReduceSickLeave - tblSalaryDtl.ReduceSickLeave - tblSalaryDtl.ReduceSickLeave
	logger.Infof("税前应发合计 %f", tblSalaryDtl.SumPayBeforeTax)

	// 应纳税所得额 = 正常工资 = 税前工资合计 - 公积金个人 - 社保个人
	tblSalaryDtl.SumTaxBase = tblSalaryDtl.SumPayBeforeTax - tblSalaryDtl.PrivateFund - tblSalaryDtl.PrivateSocial
	logger.Infof("应纳税所得额 %f", tblSalaryDtl.SumTaxBase)

	var changeCompany = false
	if len(tblEmployeeInfo.ChangeCompanyDate) > 0 && tblEmployeeInfo.ChangeCompanyDate[:6] == tblSalaryDtl.CalcMonth {
		changeCompany = true
	}

	// 累计应发工资合计 = 上一年12月到当月的税前应发工资 (如果公司主体改变清空)
	// sum_pay_calc
	//公积金累计 = 当月 + 上月历史累计(如果公司主体改变清空)
	//sum_fund = sum_fund + private_fund
	//社保累计 = 当月 + 上月历史累计(如果公司主体改变清空)
	//sum_social = sum_social + private_fund
	//	累计个税免征额
	//sum_tax_free
	//	累计已预扣预缴税额 = 2019年12月至上月
	//sum_tax_fact

	if !changeCompany {
		tblSalaryDtl.SumPayCalc = tblSalaryDtl.SumPayCalc + tblSalaryDtl.SumPayBeforeTax
		tblSalaryDtl.SumFund = tblSalaryDtl.SumFund + tblSalaryDtl.PrivateFund
		tblSalaryDtl.SumSocial = tblSalaryDtl.SumSocial + tblSalaryDtl.PrivateSocial
		tblSalaryDtl.SumTaxFree = tblSalaryDtl.SumTaxFree + 5000
		tblSalaryDtl.SumTaxFact = tblSalaryDtl.SumTaxFact
	} else {

		tblSalaryDtl.SumPayCalc = tblSalaryDtl.SumPayBeforeTax
		tblSalaryDtl.SumFund = tblSalaryDtl.PrivateFund
		tblSalaryDtl.SumSocial = tblSalaryDtl.PrivateSocial
		tblSalaryDtl.SumTaxFact = 0

		var CompanyBaseSettleAmt string

		if err := o.Raw(`
            SELECT resv from tbl_company_info where id = ? `, tblEmployeeInfo.Company).QueryRow(&CompanyBaseSettleAmt);
			err != nil {
			logger.Errorf("%s", err.Error())
			return errors.New(fmt.Sprintf("%s 查询公司配置异常", tblSalaryDtl.WorkerNo))
		}

		tblSalaryDtl.SumTaxFree = utils.Atof(CompanyBaseSettleAmt)
	}

	// 专项扣除累计
	tblSalaryDtl.SumTaxReduce = tblSalaryDtl.TaxReduceChild+
		tblSalaryDtl.TaxReduceEducation+
		tblSalaryDtl.TaxReduceHouseLoan+
		tblSalaryDtl.TaxReduceHouseRent+
		tblSalaryDtl.TaxReduceOld

	//累计预扣预缴应纳税所得额 = 累计应发工资合计 -公积金扣除累计-社保扣除累计-专项附加扣除累计(...)-累计个税免征额累计
	//sum_tax_calc = sum_pay_calc - sum_fund - sum_social - sum_tax_reduce - sum_tax_free
	tblSalaryDtl.SumTaxCalc = tblSalaryDtl.SumPayCalc - tblSalaryDtl.SumFund - tblSalaryDtl.SumSocial -
		tblSalaryDtl.SumTaxReduce - tblSalaryDtl.SumTaxFree
	logger.Infof("累计预扣预缴应纳税所得额 %f", tblSalaryDtl.SumTaxCalc)

	//累计应预扣预缴税额（公式）= 累计预扣预缴应纳税所得额*适用税率-速算扣除数
	//sum_tax_prepare =
	var taxReducsRate string
	var taxReducsAmt string
	if err := o.Raw(`
     SELECT value1+0, value2+0
      FROM tbl_calc_conf 
     WHERE key1 + 0 <= ?
       AND key2 + 0 >= ?
       AND conf_type = "personal_rate"`,
		tblSalaryDtl.SumTaxCalc, tblSalaryDtl.SumTaxCalc).QueryRow(&taxReducsRate, &taxReducsAmt); err != nil {
		logger.Errorf("%s", err.Error())
		return errors.New(fmt.Sprintf("%s 查询预扣税率配置异常", tblSalaryDtl.WorkerNo))
	}
	tblSalaryDtl.SumTaxPrepare = tblSalaryDtl.SumTaxCalc * utils.Atof(taxReducsRate) / 100 - utils.Atof(taxReducsAmt)
	logger.Infof("累计应预扣预缴税额 %f", tblSalaryDtl.SumTaxPrepare)

	//本期应预扣预缴税额(个税) = 累计应预扣预缴税额- 累计已预扣预缴税额
	//sum_tax_should	= sum_tax_prepare - sum_tax_fact
	tblSalaryDtl.SumTaxShould = tblSalaryDtl.SumTaxPrepare - tblSalaryDtl.SumTaxFact
	logger.Infof("个税 %f", tblSalaryDtl.SumTaxShould)

	var tmpNum1 float64
	if emyType == "4" {
		//非居民、劳务报酬 计算个税
		// = 应纳税所得额(税前工资收入金额 － 五险一金(个人缴纳部分) － 专项扣除 － 起征点(5000元)) x 税率 － 速算扣除数
		//sum_tax_other
		tmpNum1 = tblSalaryDtl.SumPayBeforeTax - tblSalaryDtl.PrivateFund - tblSalaryDtl.PrivateSocial -
			tblSalaryDtl.SumTaxReduce - 5000

		if err := o.Raw(`
        SELECT value1+0, value2+0
         FROM tbl_calc_conf 
        WHERE key1 + 0 <= ?
          AND key2 + 0 >= ?
          AND conf_type = "personal_rate"`,
			tmpNum1, tmpNum1).QueryRow(&taxReducsRate, &taxReducsAmt); err != nil {
			logger.Errorf("%s", err.Error())
			return errors.New(fmt.Sprintf("%s 查询预扣税率配置异常", tblSalaryDtl.WorkerNo))
		}
		tblSalaryDtl.SumTaxOther = tmpNum1 * utils.Atof(taxReducsRate) - utils.Atof(taxReducsAmt)

	} else if emyType == "5" || emyType == "6" {
		//个税计算：--劳务报酬
		//sum_tax_other
		//应纳税所得额 = 劳务报酬（少于4000元） - 800元
		//应纳税所得额 = 劳务报酬（超过4000元） × （1 - 20%）
		//应纳税额 = 应纳税所得额 × 适用税率 - 速算扣除数
		if tblSalaryDtl.PayPosition < 4000 {
			tmpNum1 = tblSalaryDtl.PayPosition - 800
		} else {
			tmpNum1 = tblSalaryDtl.PayPosition * 0.8
		}

		if err := o.Raw(`
        SELECT value1+0, value2+0
         FROM tbl_calc_conf 
        WHERE key1 + 0 <= ?
          AND key2 + 0 >= ?
          AND conf_type = "Labor_rate_tax"`,
			tmpNum1, tmpNum1).QueryRow(&taxReducsRate, &taxReducsAmt); err != nil {
			logger.Errorf("%s", err.Error())
			return errors.New(fmt.Sprintf("%s 查询预扣税率配置异常", tblSalaryDtl.WorkerNo))
		}
		tblSalaryDtl.SumTaxOther = tmpNum1 * utils.Atof(taxReducsRate) - utils.Atof(taxReducsAmt)
	}

	//其他不计税（免税工资） = 税后补款 =补偿金+代通知金
	//pay_other_after_tax	= pay_other_add + pay_wait
	tblSalaryDtl.PayOtherAfterTax  = tblSalaryDtl.PayOtherAdd + tblSalaryDtl.PayWait

	//实发工资 = 应纳税所得额 - 本期应预扣预缴税额（个税） - 其他税后扣款 + 免税工资
	//fact_pay = sum_tax_base - sum_tax_should - sum_after_tax_other + pay_other_after_tax
	tblSalaryDtl.FactPay = tblSalaryDtl.SumTaxBase - tblSalaryDtl.SumTaxShould - tblSalaryDtl.SumAfterTaxOther - tblSalaryDtl.PayOtherAfterTax
	logger.Infof("实发工资 %f", tblSalaryDtl.FactPay)

	//公司人力成本合计 = 工资合计+公司公积金 +公司社保
	//sum_company_cost = fact_pay + sum_company_fund + sum_company_social
	tblSalaryDtl.SumCompanyCost = tblSalaryDtl.FactPay + tblSalaryDtl.SumCompanyFund + tblSalaryDtl.SumCompanySocial
	logger.Infof("公司人力成本合计 %f", tblSalaryDtl.SumCompanyCost)

	if _, err := o.Raw(`
	 update
	 tbl_salary_dtl
	 set 
	   calc_work_days = ?,
	   days_calc_salary = ?,
	   days_overtime = ?,
	   days_private_leave = ?,
	   days_sick_leave = ?,
	   days_early_leave = ?,
	   days_strike_leave = ?,
	   days_food_subsidy = ?,
	   private_fund = ?,
	   private_social = ?,
	   private_old = ?,
	   private_treat = ?,
	   private_unemploy = ?,
	   private_sick = ?,
	   tax_reduce_child = ?,
	   tax_reduce_old = ?,
	   tax_reduce_education = ?,
	   tax_reduce_house_loan = ?,
	   tax_reduce_house_rent = ?,
	   company_fund = ?,
	   company_social = ?,
	   reduce_before_tax = ?,
	   reduce_afther_tax = ?,
	   pay_position = ?,
	   pay_bonus = ?,
	   pay_before_tax = ?,
	   pay_other_add = ?,
	   pay_wait = ?,
	   pay_food_subsidy = ?,
	   pay_overtime = ?,
	   pay_general = ?,
	   pay_position_add = ?,
	   pay_secret_add = ?,
	   reduce_private_leave = ?,
	   reduce_sick_leave = ?,
	   reduce_early_leave = ?,
	   reduce_strike_leave = ?,
	   pay_other = ?,
	   pay_other_after_tax = ?,
	   sum_pay_calc = ?,
	   sum_tax_free = ?,
	   sum_fund = ?,
	   sum_social = ?,
	   sum_tax_reduce = ?,
	   sum_tax_calc = ?,
	   sum_tax_prepare = ?,
	   sum_tax_fact = ?,
	   sum_tax_should = ?,
	   sum_pay_fix = ?,
	   sum_pay_before_tax = ?,
	   sum_fund_social = ?,
	   sum_tax_base = ?,
	   sum_tax_normal = ?,
	   sum_tax_other = ?,
	   sum_after_tax_other = ?,
	   fact_pay = ?,
	   sum_company_fund = ?,
	   sum_company_social = ?,
	   sum_company_cost = ?,
	   update_time = now()
	where calc_month = ?
	  and worker_no = ?
	`,
	tblSalaryDtl.CalcWorkDays,
	tblSalaryDtl.DaysCalcSalary,
	tblSalaryDtl.DaysOvertime,
	tblSalaryDtl.DaysPrivateLeave,
	tblSalaryDtl.DaysSickLeave,
	tblSalaryDtl.DaysEarlyLeave,
	tblSalaryDtl.DaysStrikeLeave,
	tblSalaryDtl.DaysFoodSubsidy,
	tblSalaryDtl.PrivateFund,
	tblSalaryDtl.PrivateSocial,
	tblSalaryDtl.PrivateOld,
	tblSalaryDtl.PrivateTreat,
	tblSalaryDtl.PrivateUnemploy,
	tblSalaryDtl.PrivateSick,
	tblSalaryDtl.TaxReduceChild,
	tblSalaryDtl.TaxReduceOld,
	tblSalaryDtl.TaxReduceEducation,
	tblSalaryDtl.TaxReduceHouseLoan,
	tblSalaryDtl.TaxReduceHouseRent,
	tblSalaryDtl.CompanyFund,
	tblSalaryDtl.CompanySocial,
	tblSalaryDtl.ReduceBeforeTax,
	tblSalaryDtl.ReduceAftherTax,
	tblSalaryDtl.PayPosition,
	tblSalaryDtl.PayBonus,
	tblSalaryDtl.PayBeforeTax,
	tblSalaryDtl.PayOtherAdd,
	tblSalaryDtl.PayWait,
	tblSalaryDtl.PayFoodSubsidy,
	tblSalaryDtl.PayOvertime,
	tblSalaryDtl.PayGeneral,
	tblSalaryDtl.PayOther,
	tblSalaryDtl.PayPositionAdd,
	tblSalaryDtl.PaySecretAdd,
	tblSalaryDtl.ReducePrivateLeave,
	tblSalaryDtl.ReduceSickLeave,
	tblSalaryDtl.ReduceEarlyLeave,
	tblSalaryDtl.ReduceStrikeLeave,
	tblSalaryDtl.PayOtherAfterTax,
	tblSalaryDtl.SumPayCalc,
	tblSalaryDtl.SumTaxFree,
	tblSalaryDtl.SumFund,
	tblSalaryDtl.SumSocial,
	tblSalaryDtl.SumTaxReduce,
	tblSalaryDtl.SumTaxCalc,
	tblSalaryDtl.SumTaxPrepare,
	tblSalaryDtl.SumTaxFact,
	tblSalaryDtl.SumTaxShould,
	tblSalaryDtl.SumPayFix,
	tblSalaryDtl.SumPayBeforeTax,
	tblSalaryDtl.SumFundSocial,
	tblSalaryDtl.SumTaxBase,
	tblSalaryDtl.SumTaxNormal,
	tblSalaryDtl.SumTaxOther,
	tblSalaryDtl.SumAfterTaxOther,
	tblSalaryDtl.FactPay,
	tblSalaryDtl.SumCompanyFund,
	tblSalaryDtl.SumCompanySocial,
	tblSalaryDtl.SumCompanyCost,
	tblSalaryDtl.CalcMonth,
	tblSalaryDtl.WorkerNo).Exec(); err != nil {
		logger.Errorf("%s", err.Error())
		return errors.New(fmt.Sprintf("%s 新增记录失败", tblSalaryDtl.WorkerNo))
	}

	return nil
}

func QueryTblSalaryDtls(tbls *[]models.TblSalaryDtl, month string) (err error) {
	var (
		o       = orm.NewOrm()
		tblName models.TblSalaryDtl
	)

	qs := o.QueryTable(tblName.TableName()).Filter("calc_month", month)
	if num, err := qs.All(tbls); err != nil {
		logger.Errorf("%s", err.Error())
		return err
	} else {
		logger.Infof("试算总数 [%s][%d人]", month, num)
	}

	return nil
}

func loadFundSalarySql(o orm.Ormer, loadfundSalary LoadfundSalary, cells []*xlsx.Cell) (err error) {
	succ, msg := utils.SetStructField(&loadfundSalary, cells)
	if !succ {
		logger.Errorf("%s", msg)
		return errors.New(msg)
	}

	//utils.Logger.Infof("insert [%v]", loadCurrentSalary)
	if _, err = o.Raw(`
		INSERT INTO tbl_salary_dtl(
		  worker_no              ,
		  worker_name            ,
		  private_fund           ,
		  private_social         ,
		  private_old            ,
		  private_treat          ,
		  private_unemploy       ,
		  private_sick           ,
		  tax_reduce_child       ,
		  tax_reduce_old         ,
		  tax_reduce_education   ,
		  tax_reduce_house_loan  ,
		  tax_reduce_house_rent  ,
		  company_fund           ,
		  company_social         ,
		  create_time            ,
          update_time            ,
		  calc_month)
		VALUES(?,?,?,?,?,?,?,?,?,?, ?,?,?,?,?,now(),now(),?)
		ON DUPLICATE KEY UPDATE
          worker_no             = ?,
          worker_name           = ?,
		  private_fund          = ?,
		  private_social        = ?,
		  private_old           = ?,
		  private_treat         = ?,
		  private_unemploy      = ?,
		  private_sick          = ?,
		  tax_reduce_child      = ?,
		  tax_reduce_old        = ?,
		  tax_reduce_education  = ?,
		  tax_reduce_house_loan = ?,
		  tax_reduce_house_rent = ?,
		  company_fund          = ?,
		  company_social        = ?,
		  update_time           = now()
		`,
		loadfundSalary.WorkerNo,
		loadfundSalary.WorkerName,
		loadfundSalary.PrivateFund,
		loadfundSalary.PrivateSocial,
		loadfundSalary.PrivateOld,
		loadfundSalary.PrivateTreat,
		loadfundSalary.PrivateUnemploy,
		loadfundSalary.PrivateSick,
		loadfundSalary.TaxReduceChild,
		loadfundSalary.TaxReduceOld,
		loadfundSalary.TaxReduceEducation,
		loadfundSalary.TaxReduceHouseLoan,
		loadfundSalary.TaxReduceHouseRent,
		loadfundSalary.CompanyFund,
		loadfundSalary.CompanySocial,
		utils.GetDate()[:6],

		loadfundSalary.WorkerNo,
		loadfundSalary.WorkerName,
		loadfundSalary.PrivateFund,
		loadfundSalary.PrivateSocial,
		loadfundSalary.PrivateOld,
		loadfundSalary.PrivateTreat,
		loadfundSalary.PrivateUnemploy,
		loadfundSalary.PrivateSick,
		loadfundSalary.TaxReduceChild,
		loadfundSalary.TaxReduceOld,
		loadfundSalary.TaxReduceEducation,
		loadfundSalary.TaxReduceHouseLoan,
		loadfundSalary.TaxReduceHouseRent,
		loadfundSalary.CompanyFund,
		loadfundSalary.CompanySocial,
	).Exec(); err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}

func loadTotalSalarySql(o orm.Ormer, loadTotalSalary LoadTotalSalary, cells []*xlsx.Cell) (err error) {
	succ, msg := utils.SetStructField(&loadTotalSalary, cells)
	if !succ {
		logger.Errorf("%s", msg)
		return errors.New(msg)
	}

	//utils.Logger.Infof("insert [%v]", loadCurrentSalary)
	if _, err = o.Raw(`
		INSERT INTO tbl_salary_dtl(
          worker_no,
		  worker_name,
          sum_pay_calc,
          sum_fund,
          sum_social,
          sum_tax_free,
          sum_tax_reduce,
          sum_tax_fact,
		  create_time,
          update_time,
		  calc_month)
		VALUES(?,?,?,?,?,?,?,?,now(),now(), ?)
		ON DUPLICATE KEY UPDATE
          worker_no = ?,
          worker_name = ?,
          sum_pay_calc = ?,
          sum_fund = ?,
          sum_social = ?,
          sum_tax_free = ?,
          sum_tax_reduce = ?,
          sum_tax_fact = ?,
          update_time = now()
		`,
		loadTotalSalary.WorkerNo,
		loadTotalSalary.WorkerName,
		loadTotalSalary.SumPayCalc,
		loadTotalSalary.SumFund,
		loadTotalSalary.SumSocial,
		loadTotalSalary.SumTaxFree,
		loadTotalSalary.SumTaxReduce,
		loadTotalSalary.SumTaxFact,
		utils.GetDate()[:6],

		loadTotalSalary.WorkerNo,
		loadTotalSalary.WorkerName,
		loadTotalSalary.SumPayCalc,
		loadTotalSalary.SumFund,
		loadTotalSalary.SumSocial,
		loadTotalSalary.SumTaxFree,
		loadTotalSalary.SumTaxReduce,
		loadTotalSalary.SumTaxFact,
	).Exec(); err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}

func loadAddSalarySql(o orm.Ormer, loadAddSalary LoadAddSalary, cells []*xlsx.Cell) (err error) {
	succ, msg := utils.SetStructField(&loadAddSalary, cells)
	if !succ {
		logger.Errorf("%s", msg)
		return errors.New(msg)
	}

	//utils.Logger.Infof("insert [%v]", loadCurrentSalary)
	if _, err = o.Raw(`
		INSERT INTO tbl_salary_dtl(
			worker_no,
			worker_name,
			reduce_before_tax,
			pay_bonus,
			pay_before_tax,
			sum_after_tax_other,
			pay_other_add,
			pay_wait,
            create_time,
			update_time,
		    calc_month)
		VALUES(?,?,?,?,?,?,?,?,now(),now(), ?)
		ON DUPLICATE KEY UPDATE
			worker_no = ?,
			worker_name = ?,
			reduce_before_tax = ?,
			pay_bonus = ?,
			pay_before_tax = ?,
			sum_after_tax_other = ?,
			pay_other_add = ?,
			pay_wait = ?,
			update_time = now()
		`,
		loadAddSalary.WorkerNo,
		loadAddSalary.WorkerName,
		loadAddSalary.ReduceBeforeTax,
		loadAddSalary.PayBonus,
		loadAddSalary.PayBeforeTax,
		loadAddSalary.SumAfterTaxOther,
		loadAddSalary.PayOtherAdd,
		loadAddSalary.PayWait,
		utils.GetDate()[:6],

		loadAddSalary.WorkerNo,
		loadAddSalary.WorkerName,
		loadAddSalary.ReduceBeforeTax,
		loadAddSalary.PayBonus,
		loadAddSalary.PayBeforeTax,
		loadAddSalary.SumAfterTaxOther,
		loadAddSalary.PayOtherAdd,
		loadAddSalary.PayWait,
	).Exec(); err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}

func loadCurrentSalarySql(o orm.Ormer, loadCurrentSalary LoadCurrentSalary, cells []*xlsx.Cell) (err error) {
	succ, msg := utils.SetStructField(&loadCurrentSalary, cells)
	if !succ {
		logger.Errorf("%s", msg)
		return errors.New(msg)
	}

	//utils.Logger.Infof("insert [%v]", loadCurrentSalary)
	if _, err = o.Raw(`
		INSERT INTO tbl_salary_dtl(
			worker_no,
			worker_name,
			days_calc_salary,
			days_private_leave,
			days_sick_leave,
			days_early_leave,
			days_strike_leave,
			days_overtime,
			days_food_subsidy,
		  create_time,
			update_time,
		  calc_month)
		VALUES(?,?,?,?,?,?,?,?,?,now(), now(),?)
		ON DUPLICATE KEY UPDATE
			worker_no = ?,
			worker_name = ?,
			days_calc_salary = ?,
			days_private_leave = ?,
			days_sick_leave = ?,
			days_early_leave = ?,
			days_strike_leave = ?,
			days_overtime = ?,
			days_food_subsidy = ?,
			update_time = now()
		`,
		loadCurrentSalary.WorkerNo,
		loadCurrentSalary.WorkerName,
		loadCurrentSalary.DaysCalcSalary,
		loadCurrentSalary.DaysPrivateLeave,
		loadCurrentSalary.DaysSickLeave,
		loadCurrentSalary.DaysEarlyLeave,
		loadCurrentSalary.DaysStrikeLeave,
		loadCurrentSalary.DaysOvertime,
		loadCurrentSalary.DaysFoodSubsidy,
		utils.GetDate()[:6],

		loadCurrentSalary.WorkerNo,
		loadCurrentSalary.WorkerName,
		loadCurrentSalary.DaysCalcSalary,
		loadCurrentSalary.DaysPrivateLeave,
		loadCurrentSalary.DaysSickLeave,
		loadCurrentSalary.DaysEarlyLeave,
		loadCurrentSalary.DaysStrikeLeave,
		loadCurrentSalary.DaysOvertime,
		loadCurrentSalary.DaysFoodSubsidy,
	).Exec(); err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}

func loadHistorySalarySql(o orm.Ormer, loadHistorySalary LoadHistorySalary, cells []*xlsx.Cell) (err error) {
	succ, msg := utils.SetStructField(&loadHistorySalary, cells)
	if !succ {
		logger.Errorf("%s", msg)
		return errors.New(msg)
	}

	utils.Logger.Infof("insert [%v]", loadHistorySalary)
	if _, err = o.Raw(`
		insert into tbl_salary_dtl(
			   calc_month,
			   worker_no,
			   worker_name,
			   pay_general,
			   pay_position_add,
			   pay_secret_add,
			   days_private_leave,
			   days_sick_leave,
			   days_overtime,
			   pay_overtime,
			   days_food_subsidy,
			   pay_bonus,
			   pay_other,
			   sum_after_tax_other,
			   sum_pay_before_tax,
			   sum_company_social,
			   sum_company_fund,
			   sum_fund,
			   sum_social,
			   sum_tax_reduce,
			   pay_other_after_tax,
			   sum_tax_should,
			   fact_pay,
              create_time,
	          update_time
		)
		values(
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, now(), now()
		)
		`, loadHistorySalary.CalcMonth,
		loadHistorySalary.WorkerNo,
		loadHistorySalary.WorkerName,
		loadHistorySalary.PayGeneral,
		loadHistorySalary.PayPositionAdd,
		loadHistorySalary.PaySecretAdd,
		loadHistorySalary.DaysPrivateLeave,
		loadHistorySalary.DaysSickLeave,
		loadHistorySalary.DaysOvertime,
		loadHistorySalary.PayOvertime,
		loadHistorySalary.DaysFoodSubsidy,
		loadHistorySalary.PayBonus,
		loadHistorySalary.PayOther,
		loadHistorySalary.SumAfterTaxOther,
		loadHistorySalary.SumPayBeforeTax,
		loadHistorySalary.SumCompanySocial,
		loadHistorySalary.SumCompanyFund,
		loadHistorySalary.SumFund,
		loadHistorySalary.SumSocial,
		loadHistorySalary.SumTaxReduce,
		loadHistorySalary.PayOtherAfterTax,
		loadHistorySalary.SumTaxShould,
		loadHistorySalary.FactPay,
	).Exec(); err != nil {
		logger.Errorf(err.Error())
		return err
	}

	return nil
}

func (self *SalaryController) jsonResult(code int, msg string) {
	resp := BaseResp{Code: code, Msg: msg}
	self.Data["json"] = resp
	self.ServeJSON()
}
