package controllers

import (
	"fmt"
	"geek-nebula/models"
	"geek-nebula/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
)

var logger = utils.Logger

func (self *SalaryController) Index() {
	self.Data["pageTitle"] = "薪酬计算"

	// 统计员工在职离职等信息,利用模板加载
	showList, err := models.GetEmployeeInfo("history")
	if err != nil {
		logger.ErrorById(self.userId, "%v", err)
		self.ajaxMsg("统计员工信息失败", MSG_ERR)
		return
	}
	self.Data["ChangeJob"] = showList.ChangeJob
	self.Data["ChangeFormal"] = showList.ChangeFormal
	self.Data["ChangeSalary"] = showList.ChangeSalary
	self.Data["GetIn"] = showList.GetIn
	self.Data["WaitForJob"] = showList.WaitForJob
	self.Data["LeaveJob"] = showList.LeaveJob
	self.Data["OnTheJob"] = showList.OnTheJob
	self.TplName = "salary/index.html"
	self.display()
}
func (self *SalaryController) Detail() {
	self.Data["pageTitle"] = "薪酬计算明细"
	self.TplName = "salary/datail.html"
	self.display()
}

func (self *SalaryController) Again()  error {
	month := beego.AppConfig.DefaultString("TEST_MONTH", utils.GetAMonth(0,-1,0)[:6])
	o := orm.NewOrm()

	if _, err := o.Raw(`delete from tbl_salary_dtl where calc_month = ?`, month).Exec(); err != nil {
		self.jsonResult(MSG_ERR, "删除上月记录失败")
		return nil
	}

	logger.Infof("%s", "上月记录已删除")
	self.jsonResult(MSG_OK, "上月记录已删除")
	return nil
}

func (self *SalaryController) Sum() error {
	var (
		err error
		rowsum int64
		o = orm.NewOrm()
		querySumRespDatas []QuerySumRespData
	)

	selectTime := self.Ctx.Request.Form.Get("selectTime")
	if len(selectTime) == 0 {
		self.jsonResult(MSG_OK, "未填写条件")
		return nil
	}

	selectTime = strings.Replace(selectTime, "-", " ", -1)
	beginMonth := strings.Replace(selectTime[:7], " ", "", -1)
	endMonth := strings.Replace(selectTime[7:], " ", "", -1)

	logger.Infof("%s - %s", beginMonth, endMonth)

	if rowsum, err = o.Raw(`
	SELECT
		calc_month,
		count( 1 ),
		sum( sum_pay_fix ),
		sum( pay_bonus ),
		sum( sum_pay_before_tax ),
		sum( sum_fund ),
		sum( sum_social ),
		sum( sum_tax_should ),
		sum( fact_pay ),
		sum( sum_company_fund ),
		sum( sum_company_social ),
		sum( sum_company_cost ) 
	FROM
		tbl_salary_dtl 
	WHERE
		fact_pay >= 0 
		and calc_month >= ?
        and calc_month <= ?
	GROUP BY
		calc_month
	`, beginMonth, endMonth).QueryRows(&querySumRespDatas); err != nil {
		logger.Errorf("%s", err.Error())
		self.jsonResult(MSG_ERR, err.Error())
	}

	resp := QuerySumResp{Code: 0, Msg: "OK", Total:rowsum, Data: querySumRespDatas}
	self.Data["json"] = resp
	logger.Infof("[%+v]", resp)
	self.ServeJSON()

	return nil
}

func (self *SalaryController) Upload() (err error) {
	var (
		o = orm.NewOrm()
		errNum = 0
		errDtl = ""
	)

	urlOption := self.Ctx.Request.Form.Get("salarytype")

	f, h, err := self.GetFile("file")
	if err != nil {
		self.jsonResult(MSG_GET_FILE_ERR, "获取文件失败")
		self.ServeJSON()
		return
	}

	logger.Infof("upload filename [%s]", h.Filename)
	logger.Infof("upload filesize [%d]", h.Size)

	defer f.Close()

	fileparm := strings.Split(h.Filename, ".")
	filename := fmt.Sprintf("%s_%s.%s", fileparm[0], utils.GetDate(), fileparm[1])
	fileDir := "file"
	filePath := strings.Join([]string{fileDir, filename}, string(os.PathSeparator))
	logger.Infof("upload save to [%s]", filePath)

	// 没有文件夹要先创建
	_, err = os.Stat(fileDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(fileDir, os.ModePerm); err != nil {
			logger.Errorf("create err [%s]", err.Error())
			self.jsonResult(CREATE_FILE_ERR, err.Error())
			return
		}
	}

	if err = self.SaveToFile("file", filePath); err != nil {
		logger.Errorf("save file err [%s]", err.Error())
		self.jsonResult(CREATE_FILE_ERR, err.Error())
		return
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		logger.Errorf("open file err [%s]", err.Error())
		self.jsonResult(CREATE_FILE_ERR, err.Error())
		return
	}

	sheet := xlFile.Sheets[0]
	for _, row := range sheet.Rows {
		if (row.Cells[0].Value == "月份" || row.Cells[0].Value == "工号") {
			continue
		}

		logger.Infof("loading [%s]", filePath)
		utils.Logger.Infof("loading [%v]", row.Cells)

		switch urlOption {
		case "now":{
			loadCurrentSalary := LoadCurrentSalary{}
			if err = loadCurrentSalarySql(o, loadCurrentSalary, row.Cells); err != nil {
				errNum += 1
				errDtl += fmt.Sprintf("员工号%s-员工姓名%s</br>", row.Cells[0].Value, row.Cells[1].Value)
			}
		}
		case "history":{
			loadHistorySalary := LoadHistorySalary{}
			if err = loadHistorySalarySql(o, loadHistorySalary, row.Cells); err != nil {
				errNum += 1
				errDtl += fmt.Sprintf("员工号%s-员工姓名%s</br>", row.Cells[0].Value, row.Cells[1].Value)
			}
		}
		case "add":{
			loadAddSalary := LoadAddSalary{}
			if err = loadAddSalarySql(o, loadAddSalary, row.Cells); err != nil {
				errNum += 1
				errDtl += fmt.Sprintf("员工号%s-员工姓名%s</br>", row.Cells[0].Value, row.Cells[1].Value)
			}
		}
		case "total":{
			loadTotalSalary := LoadTotalSalary{}
			if err = loadTotalSalarySql(o, loadTotalSalary, row.Cells); err != nil {
				errNum += 1
				errDtl += fmt.Sprintf("员工号%s-员工姓名%s</br>", row.Cells[0].Value, row.Cells[1].Value)
			}
		}
		case "fund":{
			loadfundSalary := LoadfundSalary{}
			if err = loadFundSalarySql(o, loadfundSalary, row.Cells); err != nil {
				errNum += 1
				errDtl += fmt.Sprintf("员工号%s-员工姓名%s</br>", row.Cells[0].Value, row.Cells[1].Value)
			}
		}
		}

		o.Commit()
		if errNum > 20 {
			self.jsonResult(MSG_ERR, "数据导入失败")
			return
		}
	}

	if errNum > 0 {
		self.jsonResult(MSG_ERR, fmt.Sprintf("导入失败 :%d 条</br>%s", errNum, errDtl))
	} else {
		self.jsonResult(MSG_OK, "OK")
	}

	logger.Infof("err num : %d", errNum)

	return nil
}

func (self *SalaryController) Calc() (err error) {
	var (
		errMsg string
		tblSalaryDtls []models.TblSalaryDtl
	)

	month := beego.AppConfig.DefaultString("TEST_MONTH", utils.GetAMonth(0,-1,0)[:6])

	// 获取当月所有流水
	if err = QueryTblSalaryDtls(&tblSalaryDtls, month); err != nil {
		self.jsonResult(MSG_ERR, err.Error())
		return
	}

	for _, tblSalaryDtl := range tblSalaryDtls {
		self.CalcMonthStlmDays = 0
		if err = self.CalcPriveSalary(&tblSalaryDtl); err != nil {
			errMsg += err.Error() + "</br>"
		}
	}

	if len(errMsg) > 0 {
		self.jsonResult(MSG_ERR, errMsg)
		return nil
	}

	self.jsonResult(MSG_OK, "OK")
	return nil
}