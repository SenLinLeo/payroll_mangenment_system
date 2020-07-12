package models

import (
	"errors"
	"fmt"
	"geek-nebula/utils"
	"github.com/astaxie/beego/orm"
	"reflect"
	"regexp"
	"strings"
)

// 员工信息自定义查询条件
type EmployeeQueryParam struct {
	BaseQueryParam
	Date		 string `json:"Date"`
	Name         string `json:"Name"`
	PhoneNo      string `json:"PhoneNo"`
	WorkerNo     string `json:"WorkerNo"`
	Department   string `json:"Department"`
	EmployeeType string `json:"EmployeeType"`
	IsRegEmpl	 string `json:"IsRegularEmployee"`
	Condition	 string `json:"Condition"`
}

// 员工在职离职统计数据
type EmployeeInfoShow struct {
	OnTheJob		int
	LeaveJob		int
	WaitForJob		int
	GetIn			int
	ChangeFormal	int
	ChangeSalary	int
	ChangeJob		int
}

// 新员工导入
type NewEmployees struct {
	WorkerNo           string `match:"" pk:"1"` // 工号
	WorkerName         string `match:"" pk:"1"` // 员工姓名
	Mobile             string `match:"" pk:"1"` // 手机号
	IdcardType         string `match:"" pk:"1"` // 身份证类型
	IdcardNo           string `match:"" pk:"1"` // 身份证号
	Birthday           string `match:"" pk:"1"` // 出生日期
	Sex                string `match:"" pk:"1"` // 性别
	BankName           string `match:"" pk:"1"` // 开户行
	BcardNo            string `match:"" pk:"1"` // 银行卡号
	Company            string `match:"" pk:"1"` // 公司主体
	EmployeeType       string `match:"" pk:"1"` // 聘用形式
	Department         string `match:"" pk:"1"` // 部门
	Jobs               string `match:"" pk:"1"` // 岗位
	GetInDate          string `match:"" pk:"1"` // 入职日期
	TaxType            string `match:"" pk:"1"` // 计税方式
	TaxPayer           string `match:"" pk:"1"` // 税款负担方式
	WorkCity           string `match:"" pk:"1"` // 工作城市
	TaxCity            string `match:"" pk:"1"` // 纳税城市
	Email              string `match:"" pk:"1"` // 工作邮箱
	PracticecSalary    string `match:"" pk:"1"` // 实习薪资
	GetInSalary        string `match:"" pk:"1"` // 入职薪资
	TrialSalary        string `match:"" pk:"1"` // 试用薪资
	ChangeFormalSalary string `match:"" pk:"1"` // 转正薪资
}

// 离职员工
type LeaveEmployees struct {
	WorkerNo           string `match:"" pk:"1"` // 工号
	WorkerName         string `match:"" pk:"1"` // 员工姓名
	GetOutDate         string `match:"" pk:"1"` // 离职日期
}

// 历史员工导入
type HistoryEmployees struct {
	Date			   string `match:"" pk:"1"` // 导入日期
	WorkerNo           string `match:"" pk:"1"` // 工号
	WorkerName         string `match:"" pk:"1"` // 员工姓名
	Mobile             string `match:"" pk:"1"` // 手机号
	IdcardType         string `match:"" pk:"1"` // 身份证类型
	IdcardNo           string `match:"" pk:"1"` // 身份证号
	Birthday           string `match:"" pk:"1"` // 出生日期
	Sex                string `match:"" pk:"1"` // 性别
	BankName           string `match:"" pk:"1"` // 开户行
	BcardNo            string `match:"" pk:"1"` // 银行卡号
	Company            string `match:"" pk:"1"` // 公司主体
	EmployeeType       string `match:"" pk:"1"` // 聘用形式
	Department         string `match:"" pk:"1"` // 部门
	Jobs               string `match:"" pk:"1"` // 岗位
	GetInDate          string `match:"" pk:"1"` // 入职日期
	TaxType            string `match:"" pk:"1"` // 计税方式
	TaxPayer           string `match:"" pk:"1"` // 税款负担方式
	WorkCity           string `match:"" pk:"1"` // 工作城市
	TaxCity            string `match:"" pk:"1"` // 纳税城市
	Email              string `match:"" pk:"1"` // 工作邮箱
	PracticecSalary    string `match:"" pk:"1"` // 实习薪资
	GetInSalary        string `match:"" pk:"1"` // 入职薪资
	TrialSalary        string `match:"" pk:"1"` // 试用薪资
	ChangeFormalSalary string `match:"" pk:"1"` // 转正薪资
	CurrentSalary      string `match:"" pk:"1"` // 当前薪资
	PayMonth           string `match:"" pk:"1"` // 发薪月数
	ChangeFormalDate   string `match:""` // 转正日期
	GetOutDate         string `match:""` // 离职日期
	ChangeJobDate      string `match:""` // 调岗日期
	ChangeSalaryDate   string `match:""` // 调薪日期
}

// 调岗
type ChangeJobEmployees struct {
	WorkerNo           string `match:"" pk:"1"` // 工号
	WorkerName         string `match:"" pk:"1"` // 员工姓名
	BeforeCompany      string `match:"" pk:"1"` // 调岗前公司主体：关联公司表ID
	EndCompany         string `match:"" pk:"1"` // 调岗后公司主体：关联公司表ID
	BeforeDepartment   string `match:"" pk:"1"` // 调岗前部门
	EndDepartment      string `match:"" pk:"1"` // 调岗后部门
	JobEffectDate      string `match:"" pk:"1"` // 调岗生效日期
	BeforeJob          string `match:""` 		// 调岗前岗位
	EndJob             string `match:""` 		// 调岗后岗位
}

// 调薪
type ChangeSalaryEmployees struct {
	WorkerNo           string `match:"" pk:"1"` // 工号
	WorkerName         string `match:"" pk:"1"` // 员工姓名
	ChangeSalaryAmt    string `match:"" pk:"1"` // 调薪金额
	ChangeSalaryRate   string `match:"" pk:"1"` // 调薪比例
	SalaryEffectDate   string `match:"" pk:"1"` // 调薪生效日期
	BeforeChangeSalary string `match:"" pk:"1"` // 调薪前薪资
	EndChangeSalary    string `match:"" pk:"1"` // 调薪后薪资
	ChangeSalaryReason string `match:"" pk:"1"` // 调薪理由
}

// 转正
type ChangeFormalEmployees struct {
	WorkerNo           string `match:"" pk:"1"` // 工号
	WorkerName         string `match:"" pk:"1"` // 员工姓名
	ChangeFormalDate   string `match:"" pk:"1"` // 转正日期
	ChangeFormalSalary string `match:"" pk:"1"` // 转正薪资
}

// 员工表查询
type TblEmployeeInfo struct {
	Id                 int     `orm:"column(id)" json:"Id"`
	WorkerNo           string  `orm:"column(worker_no)" json:"WorkerNo"`                      // 工号
	WorkerName         string  `orm:"column(worker_name)" json:"WorkerName"`                  // 员工姓名
	Sex                string  `orm:"column(sex)" json:"Sex"`                                 // 性别：1-男 2-女
	Mobile             string  `orm:"column(mobile)" json:"Mobile"`                           // 手机号(脱敏)
	EnMobile           string  `orm:"column(en_mobile)" json:"EnMobile"`                      // 手机号(密文)
	IdcardType         string  `orm:"column(idcard_type)" json:"IdcardType"`                  // 身份证类型：1-身份证 2-其他
	IdcardNo           string  `orm:"column(idcard_no)" json:"IdcardNo"`                      // 身份证号(脱敏)
	EnIdcardNo         string  `orm:"column(en_idcard_no)" json:"EnIdcardNo"`                 // 身份证号(密文)
	Birthday           string  `orm:"column(birthday)" json:"Birthday"`                       // 出生日期
	BcardNo            string  `orm:"column(bcard_no)" json:"BcardNo"`                        // 银行卡号(脱敏)
	EnBcardNo          string  `orm:"column(en_bcard_no)" json:"EnBcardNo"`                   // 银行卡号(密文)
	BankName           string  `orm:"column(bank_name)" json:"BankName"`                      // 开户行
	Company            string  `orm:"column(company)" json:"Company"`                         // 公司主体：关联公司表ID
	EmployeeType       string  `orm:"column(employee_type)" json:"EmployeeType"`              // 聘用形式：1-实习 2-试用 3-正式 4-外国人 5-劳务 6-顾问
	WorkStatus         string  `orm:"column(work_status)" json:"WorkStatus"`                  // 在职状态：1-在职 2-离职 3-待入职
	IsRegularEmployee  string  `orm:"column(is_regular_employee)" json:"IsRegularEmployee"`   // 转正状态：1-已转正 2-未转正
	Department         string  `orm:"column(department)" json:"Department"`                   // 部门
	Jobs               string  `orm:"column(jobs)" json:"Jobs"`                               // 岗位
	TaxType            string  `orm:"column(tax_type)" json:"TaxType"`                        // 计税类型：1-正常工资薪金 2-非居民工资 3-劳务报酬
	TaxPayer           string  `orm:"column(tax_payer)" json:"TaxPayer"`                      // 税款负担方式：1-公司支付 2-自行负担
	WorkCity           string  `orm:"column(work_city)" json:"WorkCity"`                      // 工作城市
	TaxCity            string  `orm:"column(tax_city)" json:"TaxCity"`                        // 纳税城市
	Email              string  `orm:"column(email)" json:"Email"`                             // 工作邮箱
	GetInDate          string  `orm:"column(get_in_date)" json:"GetInDate"`                   // 入职日期
	GetInSalary        string `orm:"column(get_in_salary)" json:"GetInSalary"`                // 入职薪资
	ChangeFormalDate   string  `orm:"column(change_formal_date)" json:"ChangeFormalDate"`     // 转正日期
	ChangeFormalSalary string `orm:"column(change_formal_salary)" json:"ChangeFormalSalary"`  // 转正薪资
	GetOutDate         string  `orm:"column(get_out_date)" json:"GetOutDate"`                 // 离职日期
	PracticecSalary    string `orm:"column(practicec_salary)" json:"PracticecSalary"`         // 实习薪资
	TrialSalary        string `orm:"column(trial_salary)" json:"TrialSalary"`                 // 试用薪资
	CurrentSalary      string `orm:"column(current_salary)" json:"CurrentSalary"`             // 当前薪资
	PayMonth           int     `orm:"column(pay_month)" json:"PayMonth"`                      // 发薪月数
	ChangeSalaryDate   string  `orm:"column(change_salary_date)" json:"ChangeSalaryDate"`     // 调薪日期
	ChangeJobDate      string  `orm:"column(change_job_date)" json:"ChangeJobDate"`           // 调岗日期
	CreateTime         string  `orm:"column(create_time)" json:"CreateTime"`                  // 创建时间
	UpdateTime         string  `orm:"column(update_time)" json:"UpdateTime"`                  // 创建时间
	Resv               string  `orm:"column(resv)" json:"Resv"`                               // 预留1
	ChangeCompanyDate  string `orm:"column(change_company_date)" json:"ChangeCompanyDate"`   // 公司主体变更时间

	TotalNum int64 `orm:"column(total_num)" json:"TotalNum"` // 查询条数
}

// 员工表查询方法-可拓展查询条件EmployeeQueryParam的内容
func EmployeeQuery(params *EmployeeQueryParam) ([]TblEmployeeInfo, int64, error) {
	var totalNum int64 = 0
	var startLine int = params.Page - 1
	var endLine int = startLine + params.Limit
	var filterQuery string = ""
	binds := make([]interface{}, 0)
	employeeList := make([]TblEmployeeInfo, 0)

	// 输入条件、设置绑定变量
	if params.Name != "" {
		filterQuery += " and worker_name like CONCAT('%', ?, '%') "
		binds = append(binds, params.Name)
	}
	if params.PhoneNo != "" {
		filterQuery += " and en_mobile = ? "
		binds = append(binds, GetEncData(params.PhoneNo))
	}
	if params.WorkerNo != "" {
		filterQuery += " and worker_no = ? "
		binds = append(binds, params.WorkerNo)
	}
	if params.Department != "" {
		filterQuery += " and department like CONCAT('%', ?, '%') "
		binds = append(binds, params.Department)
	}
	if params.EmployeeType != "" {
		filterQuery += " and employee_type = ? "
		binds = append(binds, params.EmployeeType)
	}
	if params.IsRegEmpl != "" {
		filterQuery += " and is_regular_employee = ? "
		binds = append(binds, params.IsRegEmpl)
	}
	if params.Condition != "" {
		if len(params.Condition) == 11 {
			filterQuery += " and en_mobile = ? "
			binds = append(binds, GetEncData(params.Condition))
		} else {
			pattern := "^\\d+(\\.\\d+)?$"
			result,_ := regexp.MatchString(pattern,params.Condition)
			if result {
				filterQuery += " and worker_no = ? "
			} else {
				filterQuery += " and worker_name like CONCAT('%', ?, '%')  "
			}
			binds = append(binds, params.Condition)
		}

	}
	fmt.Println(params.Date)
	switch params.Date {
	case "during":
		filterQuery += " and DATE_FORMAT( get_in_date, '%Y%m' ) = DATE_FORMAT( CURDATE( ), '%Y%m' ) "
	case "last":
		filterQuery += " and PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( get_in_date, '%Y%m' ) ) = 1 "
	case "history":
		filterQuery += " and PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( get_in_date, '%Y%m' ) ) > 1 "
	}

	db := orm.NewOrm()
	if err := db.Raw("select count(1) as total_num from tbl_employee_info WHERE 1 = 1 " + filterQuery, binds).
		QueryRow(&totalNum); err != nil {
		return nil, 0, err
	}

	binds = append(binds, startLine)
	binds = append(binds, endLine)

	if _, err := db.Raw(`
SELECT
		a.id,
		worker_no,
		worker_name,
		case sex when 1 then '男' when 2 then '女' end as sex,
		mobile,
		en_mobile,
		case idcard_type when 1 then '身份证' else '其他证件' end idcard_type,
		idcard_no,
		en_idcard_no,
		DATE_FORMAT(birthday,'%Y-%m-%d %H:%i:%s') birthday,
		bcard_no,
		en_bcard_no,
		bank_name,
		company_name as company,
		case employee_type 
		when 1 then '实习' 
		when 2 then '试用' 
		when 3 then '正式' 
		when 4 then '外国人' 
		when 5 then '劳务' 
		when 6 then '顾问' 
		end employee_type,
		case work_status when 1 then '在职' when 2 then '离职' when 3 then '待入职' end work_status,
		case is_regular_employee when 1 then '已转正' when 2 then '未转正' end is_regular_employee,
		department,
		jobs,
		case tax_type when 1 then '正常工资薪金' when 2 then '非居民工资' when 3 then '劳务报酬' end tax_type,
		case tax_payer when 1 then '公司支付' when 2 then '自行负担' end tax_payer,
		work_city,
		tax_city,
		email,
		DATE_FORMAT(get_in_date,'%Y-%m-%d') get_in_date,
		get_in_salary,
		DATE_FORMAT(change_formal_date,'%Y-%m-%d') change_formal_date,
		change_formal_salary,
		DATE_FORMAT(get_out_date,'%Y-%m-%d') get_out_date,
		practicec_salary,
		trial_salary,
		current_salary,
		pay_month,
		DATE_FORMAT(change_salary_date,'%Y-%m-%d') change_salary_date,
		DATE_FORMAT(change_job_date,'%Y-%m-%d') change_job_date,
		DATE_FORMAT(a.create_time,'%Y-%m-%d') create_time,
		DATE_FORMAT(a.update_time,'%Y-%m-%d') update_time,
		a.resv 
	FROM
		tbl_employee_info a LEFT JOIN tbl_company_info b on a.company = b.id
	WHERE
		1 = 1 `+filterQuery+` ORDER BY id LIMIT ?, ?`, binds).QueryRows(&employeeList); err != nil {
		return nil, 0, err
	}

	return employeeList, totalNum, nil
}

// 详情查询
func EmployeeQueryById(id string) (TblEmployeeInfo, error) {
	employeeList := TblEmployeeInfo{}

	db := orm.NewOrm()
	if _, err := db.Raw(`
SELECT
		a.id,
		worker_no,
		worker_name,
		case sex when 1 then '男' when 2 then '女' end as sex,
		mobile,
		en_mobile,
		case idcard_type when 1 then '身份证' else '其他证件' end idcard_type,
		idcard_no,
		en_idcard_no,
		DATE_FORMAT(birthday,'%Y-%m-%d %H:%i:%s') birthday,
		bcard_no,
		en_bcard_no,
		bank_name,
		company_name as company,
		case employee_type 
		when 1 then '实习' 
		when 2 then '试用' 
		when 3 then '正式' 
		when 4 then '外国人' 
		when 5 then '劳务' 
		when 6 then '顾问' 
		end employee_type,
		case work_status when 1 then '在职' when 2 then '离职' when 3 then '待入职' end work_status,
		case is_regular_employee when 1 then '已转正' when 2 then '未转正' end is_regular_employee,
		department,
		jobs,
		case tax_type when 1 then '正常工资薪金' when 2 then '非居民工资' when 3 then '劳务报酬' end tax_type,
		case tax_payer when 1 then '公司支付' when 2 then '自行负担' end tax_payer,
		work_city,
		tax_city,
		email,
		DATE_FORMAT(get_in_date,'%Y-%m-%d') get_in_date,
		get_in_salary,
		DATE_FORMAT(change_formal_date,'%Y-%m-%d') change_formal_date,
		change_formal_salary,
		DATE_FORMAT(get_out_date,'%Y-%m-%d') get_out_date,
		practicec_salary,
		trial_salary,
		current_salary,
		pay_month,
		DATE_FORMAT(change_salary_date,'%Y-%m-%d') change_salary_date,
		DATE_FORMAT(change_job_date,'%Y-%m-%d') change_job_date,
		DATE_FORMAT(a.create_time,'%Y-%m-%d') create_time,
		DATE_FORMAT(a.update_time,'%Y-%m-%d') update_time,
		a.resv 
	FROM
		tbl_employee_info a LEFT JOIN tbl_company_info b on a.company = b.id
	WHERE
		id = ? `, id).QueryRows(&employeeList); err != nil {
		return employeeList, err
	}

	return employeeList, nil
}

// 统计在职离职等信息
func GetEmployeeInfo(date string) (EmployeeInfoShow, error) {
	var duringMonth string = "DATE_FORMAT( get_in_date, '%Y%m' ) = DATE_FORMAT( CURDATE( ), '%Y%m' )"
	var lastMonth   string = "PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( get_in_date, '%Y%m' ) ) = 1"
	var hisMonth	string = "PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( get_in_date, '%Y%m' ) ) > 1"
	var filterQuery string = ""

	showList := EmployeeInfoShow{0, 0, 0, 0, 0, 0,
		0}

	db := orm.NewOrm()

	switch date {
	case "during":
		filterQuery =  duringMonth
	case "last":
		filterQuery =  lastMonth
	case "history":
		filterQuery =  hisMonth
	}

	if err := db.Raw(`
select IFNULL(max(t.cnt), 0) from (
SELECT
	COUNT( 1 ) cnt 
FROM
	tbl_employee_info 
WHERE
	work_status = 1 and ` +filterQuery + `
GROUP BY
	work_status
) t`).
		QueryRow(&showList.OnTheJob);err != nil {
		utils.Logger.Errorf("%v", err)
		return showList, err
	}
	if err := db.Raw(`
select IFNULL(max(t.cnt), 0) from (
SELECT
	COUNT( 1 ) cnt 
FROM
	tbl_employee_info 
WHERE
	work_status = 2 and ` +filterQuery + `
GROUP BY
	work_status
) t`).
		QueryRow(&showList.LeaveJob);err != nil {
		utils.Logger.Errorf("%v", err)
		return showList, err
	}
	if err := db.Raw(`
select IFNULL(max(t.cnt), 0) from (
SELECT
	COUNT( 1 ) cnt 
FROM
	tbl_employee_info 
WHERE
	work_status = 3 and ` +filterQuery + `
GROUP BY
	work_status
) t`).
		QueryRow(&showList.WaitForJob);err != nil {
		utils.Logger.Errorf("%v", err)
		return showList, err
	}
	if err := db.Raw(`
select IFNULL(max(c.cnt),0) from (SELECT
	count(1) cnt
FROM
	(
	SELECT
	CASE
		WHEN DATE_FORMAT( change_salary_date, '%Y%m' ) = DATE_FORMAT( CURDATE( ), '%Y%m' ) THEN
				'本月' 
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( change_salary_date, '%Y%m' ) ) = 1 THEN
				'上月'
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( change_salary_date, '%Y%m' ) ) > 1 THEN
				'历史' 
	END change_salary_date 
FROM
	tbl_employee_info 
WHERE
	1 = 1 and ` +filterQuery + `
	) t 
GROUP BY
	t.change_salary_date
limit 1)c `).
		QueryRow(&showList.ChangeSalary);err != nil {
		utils.Logger.Errorf("%v", err)
		return showList, err
	}
	if err := db.Raw(`
select IFNULL(max(c.cnt),0) from (
SELECT
	count(1) cnt
FROM
	(
	SELECT
	CASE
		WHEN DATE_FORMAT( change_job_date, '%Y%m' ) = DATE_FORMAT( CURDATE( ), '%Y%m' ) THEN
				'本月' 
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( change_job_date, '%Y%m' ) ) = 1 THEN
				'上月'
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( change_job_date, '%Y%m' ) ) > 1 THEN
				'历史' 
	END change_job_date 
FROM
	tbl_employee_info 
WHERE
	1 = 1 and ` +filterQuery + `
	) t 
GROUP BY
	t.change_job_date
limit 1)c`).
		QueryRow(&showList.ChangeJob);err != nil {
		utils.Logger.Errorf("%v", err)
		return showList, err
	}
	if err := db.Raw(`
select IFNULL(max(c.cnt),0) from (
SELECT
	count(1)cnt
FROM
	(
	SELECT
	CASE
		WHEN DATE_FORMAT( get_in_date, '%Y%m' ) = DATE_FORMAT( CURDATE( ), '%Y%m' ) THEN
				'本月' 
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( get_in_date, '%Y%m' ) ) = 1 THEN
				'上月'
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( get_in_date, '%Y%m' ) ) > 1 THEN
				'历史' 
	END get_in_date 
FROM
	tbl_employee_info 
WHERE
	1 = 1 and ` +filterQuery + `
	) t 
GROUP BY
	t.get_in_date
limit 1)c`).
		QueryRow(&showList.GetIn);err != nil {
		utils.Logger.Errorf("%v", err)
		return showList, err
	}
	if err := db.Raw(`
select IFNULL(max(c.cnt),0) from (
SELECT
	count(1)cnt
FROM
	(
	SELECT
	CASE
		WHEN DATE_FORMAT( change_formal_date, '%Y%m' ) = DATE_FORMAT( CURDATE( ), '%Y%m' ) THEN
				'本月' 
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( change_formal_date, '%Y%m' ) ) = 1 THEN
				'上月'
		WHEN PERIOD_DIFF( date_format( now( ), '%Y%m' ), date_format( change_formal_date, '%Y%m' ) ) > 1 THEN
				'历史' 
	END change_formal_date 
FROM
	tbl_employee_info 
WHERE
	1 = 1 and ` +filterQuery + `
	) t 
GROUP BY
	t.change_formal_date
limit 1)c`).
		QueryRow(&showList.ChangeFormal);err != nil {
		utils.Logger.Errorf("%v", err)
		return showList, err
	}
	return showList, nil
}

// 导入新员工
func ImportNewEmployees(employees NewEmployees) error {
	var EncMobile string = ""
	var EnCIdcardNo string = ""
	var EncBcardNo string = ""
	var exist int = 0

	EncMobile, EnCIdcardNo, EncBcardNo = employees.DataExchange()

	db := orm.NewOrm()
	if err := db.Raw(`select id from tbl_company_info WHERE company_name like CONCAT('%', ?, '%')`,
		employees.Company).QueryRow(&employees.Company); err != nil {
		return err
	}

	binds := make([]interface{}, 0)
	v := reflect.ValueOf(&employees).Elem()
	for i := 1; i < v.NumField(); i++ {
		binds = append(binds, v.Field(i).String())
	}
	binds = append(binds, EncMobile)
	binds = append(binds, EnCIdcardNo)
	binds = append(binds, EncBcardNo)
	binds = append(binds, employees.GetInSalary)
	binds = append(binds, employees.WorkerNo)

	if err := db.Raw(`select count(1) as cnt from tbl_employee_info WHERE worker_no = ?`,
		employees.WorkerNo).QueryRow(&exist); err != nil {
		return err
	}

	if exist == 0 {
		if _, err := db.Raw(`
			update tbl_employee_info set
			worker_name				=	?,
			mobile					=	?,
			idcard_type				=	?,
			idcard_no				=	?,
			birthday				=	DATE(?),
			sex						=	?,
			bank_name				=	?,
			bcard_no				=	?,
			company					=	?,
			employee_type			=	?,
			department				=	?,
			jobs					=	?,
			get_in_date				=	DATE(?),
			tax_type				=	?,
			tax_payer				=	?,
			work_city				=	?,
			tax_city				=	?,
			email					=	?,
			practicec_salary		=	?,
			get_in_salary			=	?,
			trial_salary			=	?,
			change_formal_salary	=	?,
			pay_month				=   0,
			create_time				=	CURRENT_TIMESTAMP,
			update_time				=	CURRENT_TIMESTAMP,
			en_mobile				=	?,
			en_idcard_no			=	?,
			en_bcard_no				=	?,
			current_salary			=	?,
			import_time				=   DATE_FORMAT(CURRENT_TIMESTAMP,'%Y%m')
			where worker_no = ?`, binds).Exec(); err != nil {
			return err
		}
	} else {
		if _, err := db.Raw(`
		insert into tbl_employee_info 
		(
		worker_name,
		mobile, 
		idcard_type,
		idcard_no,
		birthday,
		sex,
		bank_name,
		bcard_no,
		company,
		employee_type,
		department,
		jobs,
		get_in_date,
		tax_type,
		tax_payer,
		work_city,
		tax_city,
		email,
		practicec_salary,
		get_in_salary,
		trial_salary,
		change_formal_salary,
		create_time,
        update_time,
		en_mobile,
		en_idcard_no,
		en_bcard_no,
		current_salary,
	    worker_no,
		import_time) VALUES
		(?, ?, ?, ?, ?, DATE(?), ?, ?, ?, ?, ?, ?, ?, DATE(?), 
		 ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, ?, ?,
		 DATE_FORMAT(CURRENT_TIMESTAMP,'%Y%m'))`, binds).
			Exec(); err != nil {
			db.Rollback()
			return err
		}
	}
	db.Commit()

	return nil
}

// 数据加密封装
func (employees *NewEmployees)DataExchange() (string, string, string) {
	EncMobile := GetEncData(employees.Mobile)
	EnCIdcardNo := GetEncData(employees.IdcardNo)
	EncBcardNo := GetEncData(employees.BcardNo)
	employees.Mobile = employees.Mobile[0:3] + "****" + employees.Mobile[6:10]
	employees.IdcardNo = employees.IdcardNo[0:6] + "****" + employees.IdcardNo[len(employees.IdcardNo)-4:]
	employees.BcardNo = employees.BcardNo[0:6] + "****" + employees.BcardNo[len(employees.BcardNo)-4:]
	employees.PracticecSalary = GetEncData(employees.PracticecSalary)
	employees.GetInSalary = GetEncData(employees.GetInSalary)
	employees.TrialSalary = GetEncData(employees.TrialSalary)
	employees.ChangeFormalSalary = GetEncData(employees.ChangeFormalSalary)
	employees.Department = strings.Split(employees.Department, "||")[0]

	switch employees.TaxPayer {
	case "公司支付":
		employees.TaxPayer = "1"
	case "自行负担":
		employees.TaxPayer = "2"
	}
	switch employees.TaxType {
	case "正常工资薪金":
		employees.TaxType = "1"
	case "非居民工资":
		employees.TaxType = "2"
	case "劳务报酬":
		employees.TaxType = "3"
	}
	switch employees.EmployeeType {
	case "实习":
		employees.EmployeeType = "1"
	case "试用":
		employees.EmployeeType = "2"
	case "正式":
		employees.EmployeeType = "3"
	case "外国人":
		employees.EmployeeType = "4"
	case "劳务":
		employees.EmployeeType = "5"
	case "顾问":
		employees.EmployeeType = "3"

	}
	switch employees.Sex {
	case "男":
		employees.Sex = "1"
	case "女":
		employees.Sex = "2"
	}
	switch employees.IdcardType {
	case "身份证":
		employees.IdcardType = "1"
	case "其他":
		employees.IdcardType = "2"
	}

	return EncMobile, EnCIdcardNo, EncBcardNo
}

func (employees *TblEmployeeInfo)DataExchange() {
	employees.GetInSalary = GetDecData(employees.GetInSalary)
	employees.ChangeFormalSalary = GetDecData(employees.ChangeFormalSalary)
	employees.CurrentSalary = GetDecData(employees.CurrentSalary)
	employees.PracticecSalary = GetDecData(employees.PracticecSalary)
	employees.TrialSalary = GetDecData(employees.TrialSalary)
}

func (employees *HistoryEmployees)DataExchange() (string, string, string) {
	EncMobile := GetEncData(employees.Mobile)
	EnCIdcardNo := GetEncData(employees.IdcardNo)
	EncBcardNo := GetEncData(employees.BcardNo)
	employees.Mobile = employees.Mobile[0:3] + "****" + employees.Mobile[6:10]
	employees.IdcardNo = employees.IdcardNo[0:2] + "****" + employees.IdcardNo[len(employees.IdcardNo)-2:]
	employees.BcardNo = employees.BcardNo[0:2] + "****" + employees.BcardNo[len(employees.BcardNo)-2:]
	employees.PracticecSalary = GetEncData(employees.PracticecSalary)
	employees.GetInSalary = GetEncData(employees.GetInSalary)
	employees.TrialSalary = GetEncData(employees.TrialSalary)
	employees.ChangeFormalSalary = GetEncData(employees.ChangeFormalSalary)
	employees.CurrentSalary = GetEncData(employees.CurrentSalary)
	employees.Department = strings.Split(employees.Department, "|")[0]

	if employees.GetInDate != "" {
		employees.GetInDate = utils.ConvertToFormatDay(employees.GetInDate)
	}
	if employees.ChangeFormalDate != "" {
		employees.ChangeFormalDate = utils.ConvertToFormatDay(employees.ChangeFormalDate)
	}
	if employees.GetOutDate != "" {
		employees.GetOutDate = utils.ConvertToFormatDay(employees.GetOutDate)
	}
	if employees.ChangeJobDate != "" {
		employees.ChangeJobDate = utils.ConvertToFormatDay(employees.ChangeJobDate)
	}
	if employees.ChangeSalaryDate != "" {
		employees.ChangeSalaryDate = utils.ConvertToFormatDay(employees.ChangeSalaryDate)
	}

	switch employees.TaxPayer {
	case "公司支付":
		employees.TaxPayer = "1"
	case "自行负担":
		employees.TaxPayer = "2"
	}
	switch employees.TaxType {
	case "正常工资薪金":
		employees.TaxType = "1"
	case "非居民工资":
		employees.TaxType = "2"
	case "劳务报酬":
		employees.TaxType = "3"
	case "外国人薪金":
		employees.TaxType = "4"
	}
	switch employees.EmployeeType {
	case "实习":
		employees.EmployeeType = "1"
	case "试用":
		employees.EmployeeType = "2"
	case "正式":
		employees.EmployeeType = "3"
	case "外国人":
		employees.EmployeeType = "4"
	case "劳务":
		employees.EmployeeType = "5"
	case "顾问":
		employees.EmployeeType = "3"

	}
	switch employees.Sex {
	case "男":
		employees.Sex = "1"
	case "女":
		employees.Sex = "2"
	}
	switch employees.IdcardType {
	case "身份证":
		employees.IdcardType = "1"
	case "其他":
		employees.IdcardType = "2"
	}

	return EncMobile, EnCIdcardNo, EncBcardNo
}

func (employees *ChangeSalaryEmployees)DataExchange() {
	employees.BeforeChangeSalary = GetEncData(employees.BeforeChangeSalary)
	employees.ChangeSalaryAmt = GetEncData(employees.ChangeSalaryAmt)
	employees.EndChangeSalary = GetEncData(employees.EndChangeSalary)
}

// 导入离职员工
func ImportLeaveEmployees(employees LeaveEmployees) error {
	var  exist int

	db := orm.NewOrm()
	binds := make([]interface{}, 0)
	binds = append(binds, employees.GetOutDate)
	binds = append(binds, employees.WorkerNo)
	if err := db.Raw("select count(1) as cnt from tbl_employee_info where worker_no = ?",employees.WorkerNo).
		QueryRow(&exist); err != nil {
		return err
	}
	if exist == 0 {
		return errors.New("工号不存在")
	}

	if _, err := db.Raw("update tbl_employee_info set get_out_date = date(?) where worker_no = ?",binds).Exec();
		err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

// 导入历史员工
func ImportHistoryEmployees(employees HistoryEmployees) error {
	var EncMobile string = ""
	var EnCIdcardNo string = ""
	var EncBcardNo string = ""
	var exist int = 0

	// 数据加密
	EncMobile, EnCIdcardNo, EncBcardNo = employees.DataExchange()

	db := orm.NewOrm()
	if err := db.Raw(`select id from tbl_company_info WHERE company_name like CONCAT('%', ?, '%')`,
		employees.Company).QueryRow(&employees.Company); err != nil {
		return err
	}

	binds := make([]interface{}, 0)
	v := reflect.ValueOf(&employees).Elem()
	for i := 0; i < v.NumField(); i++ {
		if i == 1 {
			continue
		}
		if v.Field(i).String() == "" {
			continue
		}
		binds = append(binds, v.Field(i).String())
	}
	binds = append(binds, EncMobile)
	binds = append(binds, EnCIdcardNo)
	binds = append(binds, EncBcardNo)
	binds = append(binds, employees.WorkerNo)

	var filterInsert string = ""
	var bindVar		 string = ""
	var filterUpdate string = ""
	if employees.ChangeFormalDate != "" {
		filterInsert += " change_formal_date, "
		bindVar += " DATE(?), "
		filterUpdate += " change_formal_date = DATE(?), "
	}
	if employees.GetOutDate != "" {
		filterInsert += " get_out_date, "
		bindVar += " DATE(?), "
		filterUpdate += " get_out_date = DATE(?), "
	}
	if employees.ChangeJobDate != "" {
		filterInsert += " change_job_date, "
		bindVar += " DATE(?), "
		filterUpdate += " change_job_date = DATE(?), "
	}
	if employees.ChangeSalaryDate != "" {
		filterInsert += " change_salary_date, "
		bindVar += " DATE(?), "
		filterUpdate += " change_salary_date = DATE(?), "
	}

	if err := db.Raw(`select count(1) as cnt from tbl_employee_info WHERE worker_no = ?`,
		employees.WorkerNo).QueryRow(&exist); err != nil {
		return err
	}

	if exist != 0 {
		if _, err := db.Raw(`
			update tbl_employee_info set
			import_time				=   ?,
			worker_name				=	?,
			mobile					=	?,
			idcard_type				=	?,
			idcard_no				=	?,
			birthday				=	DATE(?),
			sex						=	?,
			bank_name				=	?,
			bcard_no				=	?,
			company					=	?,
			employee_type			=	?,
			department				=	?,
			jobs					=	?,
			get_in_date				=	DATE(?),
			tax_type				=	?,
			tax_payer				=	?,
			work_city				=	?,
			tax_city				=	?,
			email					=	?,
			practicec_salary		=	?,
			get_in_salary			=	?,
			trial_salary			=	?,
			change_formal_salary	=	?,
			current_salary			=   ?,
			pay_month				= 	pay_month + ?,
			` + filterUpdate +`
			create_time				=	CURRENT_TIMESTAMP,
			update_time				=	CURRENT_TIMESTAMP,
			en_mobile				=	?,
			en_idcard_no			=	?,
			en_bcard_no				=	?
			where worker_no = ?`, binds).Exec(); err != nil {
			return err
		}
	} else {
		if _, err := db.Raw(`
		insert into tbl_employee_info 
		(
		import_time,
		worker_name,
		mobile, 
		idcard_type,
		idcard_no,
		birthday,
		sex,
		bank_name,
		bcard_no,
		company,
		employee_type,
		department,
		jobs,
		get_in_date,
		tax_type,
		tax_payer,
		work_city,
		tax_city,
		email,
		practicec_salary,
		get_in_salary,
		trial_salary,
		change_formal_salary,
		current_salary,
		pay_month,
		` + filterInsert + `
		create_time,
        update_time,
		en_mobile,
		en_idcard_no,
		en_bcard_no,
	    worker_no) VALUES
		(?, ?, ?, ?, ?, DATE(?), ?, ?, ?, ?, ?, ?, ?, DATE(?), 
		 ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
		 ` + bindVar + `
		 CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, ?, ?)`, binds).
			Exec(); err != nil {
			db.Rollback()
			return err
		}
	}
	db.Commit()

	return nil
}

// 调薪
func ImportChangeSalary(employees ChangeSalaryEmployees) error {
	db := orm.NewOrm()
	// 判断工号是否存在
	if err := existWorker(employees.WorkerNo); err != nil {
		return err
	}

	employees.DataExchange()

	binds := make([]interface{}, 0)
	v := reflect.ValueOf(&employees).Elem()
	for i := 0; i < v.NumField(); i++ {
		if i == 0 {
			continue
		}
		if v.Field(i).String() == "" {
			continue
		}
		binds = append(binds, v.Field(i).String())
	}
	binds = append(binds, employees.WorkerNo)

	if _, err := db.Raw(`
			INSERT INTO tbl_employee_change_salary
(worker_name, change_salary_date, change_salary_amt, change_salary_rate, salary_effect_date, before_change_salary, 
 end_change_salary, change_salary_reason, worker_no, create_time, update_time)
VALUES
(?, CURRENT_TIMESTAMP, ?, ?, date(?), ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		binds).Exec(); err != nil {
		db.Rollback()
		return err
	}

	if _, err := db.Raw(`
			update tbl_employee_info set change_salary_date = date(?) where worker_no = ?`,
		employees.SalaryEffectDate, employees.WorkerNo).Exec(); err != nil {
		db.Rollback()
		return err
	}

	db.Commit()

	return nil
}

// 判断工号是否存在
func existWorker(WorkerNo string) error {
	var exist int = 0

	db := orm.NewOrm()
	// 判断工号是否存在
	if err := db.Raw(`select count(1) from tbl_employee_info WHERE worker_no = ?`,
		WorkerNo).QueryRow(&exist); err != nil {
		return err
	}
	if exist == 0 {
		return errors.New("工号不存在")
	}
	return nil
}

// 调岗
func ImportChangeJob(employees ChangeJobEmployees) error {
	var nowCompany string = ""
	var nowJobs string = ""
	var exist int = 0
	var ChangeJobDate string = ""

	db := orm.NewOrm()
	// 判断工号是否存在
	if err := existWorker(employees.WorkerNo); err != nil {
		return err
	}

	// 获取员工当前公司主体和岗位
	if err := db.Raw(`select a.company, a.jobs from tbl_employee_info a LEFT JOIN tbl_company_info b 
								on a.company = b.id where a.worker_no = ?`,
		employees.WorkerNo).QueryRow(&nowCompany, &nowJobs); err != nil {
		return err
	}

	// 获取公司主体ID
	if err := db.Raw(`select id from tbl_company_info WHERE company_name like CONCAT('%', ?, '%')`,
		employees.BeforeCompany).QueryRow(&employees.BeforeCompany); err != nil {
		return err
	}

	if err := db.Raw(`select id from tbl_company_info WHERE company_name like CONCAT('%', ?, '%')`,
		employees.EndCompany).QueryRow(&employees.EndCompany); err != nil {
		return err
	}
	// 获取调岗前后部门、职位
	beforeDepartment := strings.Split(employees.BeforeDepartment, "||")
	employees.BeforeDepartment = beforeDepartment[0]
	if len(beforeDepartment) > 1 {
		employees.BeforeJob = beforeDepartment[1]
	}
	endDepartment := strings.Split(employees.EndDepartment, "||")
	employees.EndDepartment = endDepartment[0]
	if len(endDepartment) > 1 {
		employees.EndJob = endDepartment[1]
	}

	binds := make([]interface{}, 0)
	v := reflect.ValueOf(&employees).Elem()
	for i := 0; i < v.NumField(); i++ {
		if i == 0 || i == 6 {
			continue
		}
		binds = append(binds, v.Field(i).String())
	}

	if nowCompany != employees.EndCompany {
		ChangeJobDate = " date_add(curdate(), interval - day(curdate()) + 1 day) "
		if err := db.Raw(`select count(1) from tbl_employee_change_jobs 
								WHERE worker_no = ? and job_effect_date = 
								date_add(curdate(), interval - day(curdate()) + 1 day)`,
			employees.WorkerNo).QueryRow(&exist); err != nil {
			return err
		}
	} else {
		ChangeJobDate = " date(?) "
		binds = append(binds, employees.JobEffectDate)
		if err := db.Raw(`select count(1) from tbl_employee_change_jobs 
								WHERE worker_no = ? and job_effect_date = date(?)`,
			employees.WorkerNo, employees.JobEffectDate).QueryRow(&exist); err != nil {
			return err
		}
	}
	binds = append(binds, employees.WorkerNo)

	if exist != 0 {
		if _, err := db.Raw(`
			UPDATE tbl_employee_change_jobs 
				SET worker_name = ?,
					before_company = ?,
					end_company = ?,
					before_department = ?,
					end_department = ?,
					before_job = ?,
					end_job = ?,
					job_effect_date = ` + ChangeJobDate + `,
					change_job_date = CURRENT_TIMESTAMP
			WHERE
				worker_no = ?`,
			binds).Exec(); err != nil {
			db.Rollback()
			return err
		}
	} else {
		if _, err := db.Raw(`
			INSERT INTO tbl_employee_change_jobs 
(worker_name, before_company, end_company, before_department, end_department, before_job, end_job, job_effect_date, 
 worker_no, change_job_date, create_time, update_time)
VALUES
(?, ?, ?, ?, ?, ?, ?, ` + ChangeJobDate + `, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
			binds).Exec(); err != nil {
			db.Rollback()
			return err
		}
	}

	if _, err := db.Raw(`
			update tbl_employee_info set change_job_date = date(?) where worker_no = ?`,
		employees.JobEffectDate, employees.WorkerNo).Exec(); err != nil {
		db.Rollback()
		return err
	}

	db.Commit()

	return nil
}

// 转正
func ImportChangeFormal(employees ChangeFormalEmployees) error  {
	var IsFormal int = 0

	db := orm.NewOrm()
	// 判断工号是否存在
	if err := existWorker(employees.WorkerNo); err != nil {
		return err
	}

	employees.ChangeFormalSalary = GetEncData(employees.ChangeFormalSalary)

	// 判断是否转正
	if err := db.Raw(`select is_regular_employee from tbl_employee_info WHERE worker_no = ?`,
		employees.WorkerNo).QueryRow(&IsFormal); err != nil {
		return err
	}
	fmt.Println(IsFormal)
	fmt.Println(employees)
	if IsFormal == 1 {// 已经转正
		if _, err := db.Raw(`
			UPDATE tbl_employee_info
				SET change_formal_date = ?,
					change_formal_salary = ?,
					change_formal_salary = ?,
					update_time = CURRENT_TIMESTAMP
			WHERE
				worker_no = ?`,
			employees.ChangeFormalDate, employees.ChangeFormalSalary, employees.ChangeFormalSalary, employees.WorkerNo).
			Exec(); err != nil {
			db.Rollback()
			return err
		}
	} else {
		if _, err := db.Raw(`
			UPDATE tbl_employee_info
				SET change_formal_date = ?,
					change_formal_salary = ?,
					update_time = CURRENT_TIMESTAMP
			WHERE
				worker_no = ?`,
			employees.ChangeFormalDate, employees.ChangeFormalSalary, employees.WorkerNo).Exec(); err != nil {
			db.Rollback()
			return err
		}
	}

	return nil
}

// 单个删除离职员工
func Delete(id string) error {
	db := orm.NewOrm()
	// 判断工号是否存在
	if _, err := db.Raw(`delete from tbl_employee_info where id = ? and work_status = 2`,
		id).Exec(); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}