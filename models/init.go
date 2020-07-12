package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"net/url"
	"geek-nebula/utils"
)

func Init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	// fmt.Println(dsn)

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(Auth), new(Role), new(RoleAuth), new(Admin),
		new(Group), new(Env), new(Code), new(ApiSource), new(ApiDetail), new(ApiPublic), new(Template), new(TblSalaryDtl))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	// 定时任务,每分钟巡检一次在职离职调薪调岗日期，修改员工基本信息
	m := cron.New()
	m.AddFunc("1 * * * * *", func() {
		BatchChangeFormal()
		BatchChangeSalary()
		BatchChangeJobs()
		BatchLeave()
		BatchIn()
	})
	m.Start()
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}

// BatchChangeFormal 批量转正
func BatchChangeFormal() (error) {
	utils.Logger.Infof("%s", "start BatchChangeFormal")
	db := orm.NewOrm()
	if _, err := db.Raw(`
			update tbl_employee_info set is_regular_employee = 1, current_salary = change_formal_salary 
				WHERE change_formal_date >= CURRENT_DATE`).Exec(); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

// BatchChangeSalary 批量调薪
func BatchChangeSalary() (error) {
	var salary int = 0
	var workerNo int = 0
	db := orm.NewOrm()
	if err := db.Raw(`
			SELECT b.end_change_salary, a.worker_no from tbl_employee_info a, tbl_employee_change_salary b  
					WHERE a.worker_no = b.worker_no ORDER BY b.change_salary_date desc limit 1`).
		QueryRow(&salary, &workerNo); err != nil {
		return err
	}

	if _, err := db.Raw(`
			update tbl_employee_info set current_salary = ? 
				WHERE worker_no = ?`, salary, workerNo).Exec(); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

// BatchChangeJobs 批量调岗
func BatchChangeJobs() (error) {
	var endDepartment string = ""
	var endCompany string = ""
	var endJob string = ""
	var workerNo string = ""

	db := orm.NewOrm()
	if err := db.Raw(`
			SELECT b.end_department, b.end_company, b.end_job, a.worker_no from tbl_employee_info a, 
				tbl_employee_change_jobs b  WHERE a.worker_no = b.worker_no ORDER BY b.change_job_date desc limit 1`).
		QueryRow(&endDepartment, &endCompany, &endJob, &workerNo); err != nil {
		return err
	}

	if _, err := db.Raw(`
			update tbl_employee_info set department = ? , company = ? , jobs = ?, 
					change_company_date = date_add(curdate(), interval - day(curdate()) + 1 day)
				WHERE worker_no = ?`, endDepartment, endCompany, endJob, workerNo).Exec(); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}
// BatchLeave 批量离职
func BatchLeave() (error) {
	db := orm.NewOrm()
	if _, err := db.Raw(`
			update tbl_employee_info set work_status = 2
				WHERE get_out_date >= CURRENT_DATE`).Exec(); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}

// BatchIn 批量入职
func BatchIn() (error) {
	db := orm.NewOrm()
	if _, err := db.Raw(`
			update tbl_employee_info set work_status = 1
				WHERE get_in_date >= CURRENT_DATE and work_status = 3`).Exec(); err != nil {
		db.Rollback()
		return err
	}
	db.Commit()
	return nil
}
