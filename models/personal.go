package models

import (
	"github.com/astaxie/beego/orm"
	"strings"
)
import "geek-nebula/utils"

func PersonalReportQuery(params *BaseQueryParam, date string) ([]TblSalaryDtl, int64, error) {
	var totalNum int64 = 0
	var startLine int = params.Page - 1
	var endLine int = startLine + params.Limit
	var filterQuery string = ""
	binds := make([]interface{}, 0)

	personalList := make([]TblSalaryDtl, 0)

	// 输入条件、设置绑定变量
	if date != "" {
		startDate := strings.Split(date, "-")[0]
		endDate := strings.Split(date, "-")[1]
		filterQuery += " and calc_month between ? and ? "
		binds = append(binds, startDate, endDate)
	}
	db := orm.NewOrm()
	if err := db.Raw("select count(1) as total_num from tbl_salary_dtl WHERE 1 = 1 " + filterQuery, binds).
		QueryRow(&totalNum); err != nil {
		return nil, 0, err
	}
	binds = append(binds, startLine)
	binds = append(binds, endLine)

	if _, err := db.Raw(`
SELECT
	worker_no,
	worker_name,
	calc_month,
	pay_general,
	pay_position_add,
	pay_secret_add,
	days_private_leave,
	days_sick_leave,
	days_overtime,
	pay_overtime,
	pay_food_subsidy,
	pay_bonus,
	pay_before_tax,
	reduce_before_tax,
	sum_pay_before_tax,
	company_social,
	company_fund,
	private_social,
	private_fund,
	sum_tax_reduce,
	pay_other_after_tax,
	sum_tax_normal,
	sum_after_tax_other,
	fact_pay 
FROM
	tbl_salary_dtl 
WHERE
	1 = 1 ` + filterQuery + ` limit ?, ?`, binds).QueryRows(&personalList);err != nil {
		utils.Logger.Errorf("%v", err)
		return personalList, 0, err
	}

	return personalList, totalNum, nil
}