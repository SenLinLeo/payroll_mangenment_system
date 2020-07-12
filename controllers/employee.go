package controllers

import (
	"fmt"
	"geek-nebula/models"
	"geek-nebula/utils"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type EmployeeController struct {
	BaseController
}

func (self *EmployeeController) Index() {
	self.Data["pageTitle"] = "员工信息维护"
	//self.display()
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
	self.TplName = "employee/index.html"
	self.display()
}

// 表格数据加载
func (self *EmployeeController) List() {
	var params models.EmployeeQueryParam
	params.Page, _ = self.GetInt("page")
	params.Limit, _ = self.GetInt("limit")
	params.Date = self.GetString("month")
	params.EmployeeType = self.GetString("EmployeeType")
	params.Department = self.GetString("Department")
	params.IsRegEmpl = self.GetString("IsRegularEmployee")
	params.Condition = self.GetString("Condition")
	self.TplName = "employee/index.html"

	employeeList, totalNum, err := models.EmployeeQuery(&params)
	if err != nil {
		logger.ErrorById(self.userId, "%v", err)
		self.ajaxMsg("查询员工信息失败", MSG_ERR)
		return
	}

	list := make([]map[string]interface{}, len(employeeList))
	for k, employee := range employeeList {
		v := reflect.ValueOf(&employee).Elem()
		row := make(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			row[v.Type().Field(i).Tag.Get("json")] = v.Field(i).Interface()
		}
		list[k] = row
	}

	self.ajaxList("成功", MSG_OK, totalNum, list)
}

// 上传文件处理,根据参数区分接口
func (self *EmployeeController) Upload() {
	var newEmployee models.NewEmployees
	var leaveEmployees models.LeaveEmployees
	var historyEmployees models.HistoryEmployees
	var changeSalaryEmployees models.ChangeSalaryEmployees
	var changeJobEmployees models.ChangeJobEmployees
	var changeFormalEmployees models.ChangeFormalEmployees
	var tempDate string = ""
	LastMonthFirstDay := utils.GetLastMonthFirstDay()
	t1, err := time.Parse("20060102", LastMonthFirstDay)
	// 上个月
	lastMonth, _ := time.Parse("20060102", utils.GetAMonth(0, -1, 0))
	// 去年12月
	lastYear, _ := time.Parse("20060102", utils.GetAMonth(-1,0, 0)[:4] + "12")


	urlOption := self.Ctx.Request.Form.Get("uploadtype")

	// Excel保存
	filePath := self.saveFile()
	if filePath == "" {
		logger.ErrorById(self.userId, "%s", "保存文件失败")
		self.ajaxMsg("保存文件失败", MSG_ERR)
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		logger.ErrorById(self.userId, `file open failed [%s]`, err.Error())
	}

	sheet := xlFile.Sheets[0]
	for k, row := range sheet.Rows {
		// 跳过标题行
		if k == 0 {
			continue
		}
		switch urlOption {
		//新员工导入
		case "new":
			err, msg := SetStructField(&newEmployee ,row.Cells)
			if !err {
				self.ajaxMsg(msg, MSG_ERR)
			}
			// 离职日期只支持上月1号之后
			t2, _ := time.Parse("20060102", newEmployee.GetInDate)
			if t2.Before(t1) {
				//处理逻辑
				self.ajaxMsg(fmt.Sprintf("处理失败, 入职日期只支持上月1号之后[%s]", newEmployee.GetInDate), MSG_ERR)
				return
			}
			result := models.ImportNewEmployees(newEmployee)
			if result != nil {
				logger.ErrorById(self.userId, "新员工导入失败 [%s]", result)
			}
			// 批量离职
		case "leave":
			err, msg := SetStructField(&leaveEmployees ,row.Cells)
			if !err {
				self.ajaxMsg(msg, MSG_ERR)
			}

			// 离职日期只支持上月1号之后
			t2, _ := time.Parse("20060102", leaveEmployees.GetOutDate)
			if t2.Before(t1) {
				//处理逻辑
				self.ajaxMsg(fmt.Sprintf("处理失败, 离职日期只支持上月1号之后[%s]", leaveEmployees.GetOutDate), MSG_ERR)
				return
			}
			result := models.ImportLeaveEmployees(leaveEmployees)
			if result != nil {
				logger.ErrorById(self.userId, "批量离职修改失败 [%s]", result)
			}
			// 历史员工导入
		case "history":
			err, msg := SetStructField(&historyEmployees ,row.Cells)
			if !err {
				self.ajaxMsg(msg, MSG_ERR)
			}

			if tempDate != "" && tempDate != historyEmployees.Date {
				self.ajaxMsg(fmt.Sprintf("处理失败, 历史员工导入只支持单月导入[%s]", historyEmployees.Date), MSG_ERR)
				return
			}
			tempDate = historyEmployees.Date

			// 离职日期只支持上月1号之后
			t2, _ := time.Parse("200601", historyEmployees.Date)
			if t2.Before(lastYear) || t2.After(lastMonth) {
				//处理逻辑
				self.ajaxMsg(fmt.Sprintf("处理失败, 历史员工导入只支持本年度历史日期[%s]", historyEmployees.Date), MSG_ERR)
				return
			}
			result := models.ImportHistoryEmployees(historyEmployees)
			if result != nil {
				logger.ErrorById(self.userId, "插入历史员工信息失败 [%s]", result)
			}
			// 调薪
		case "changeSalary":
			err, msg := SetStructField(&changeSalaryEmployees ,row.Cells)
			if !err {
				self.ajaxMsg(msg, MSG_ERR)
			}

			result := models.ImportChangeSalary(changeSalaryEmployees)
			if result != nil {
				logger.ErrorById(self.userId, "调薪失败 [%s]", result)
			}
			// 调岗
		case "changeJob":
			err, msg := SetStructField(&changeJobEmployees ,row.Cells)
			if !err {
				self.ajaxMsg(msg, MSG_ERR)
			}

			result := models.ImportChangeJob(changeJobEmployees)
			if result != nil {
				logger.ErrorById(self.userId, "调岗失败 [%s]", result)
			}
			// 转正
		case "changeFormal":
			err, msg := SetStructField(&changeFormalEmployees ,row.Cells)
			if !err {
				self.ajaxMsg(msg, MSG_ERR)
			}
			// 转正日期只支持上月1号之后
			t2, _ := time.Parse("20060102", changeFormalEmployees.ChangeFormalDate)
			if t2.Before(t1) {
				//处理逻辑
				self.ajaxMsg(fmt.Sprintf("处理失败, 转正日期只支持上月1号之后[%s]",
					changeFormalEmployees.ChangeFormalDate), MSG_ERR)
				return
			}
			result := models.ImportChangeFormal(changeFormalEmployees)
			if result != nil {
				logger.ErrorById(self.userId, "转正失败 [%s]", result)
			}
		}

	}
	self.ajaxMsg("导入成功", MSG_OK)
}

// 文件保存封装
func (self *EmployeeController)saveFile() string {
	f, h, err := self.GetFile("file")
	if err != nil {
		self.ajaxMsg("获取文件失败", MSG_ERR)
		return ""
	}
	defer f.Close()

	logger.Infof("upload filename [%s]",  h.Filename)
	logger.Infof("upload filesize [%d]",  h.Size)

	fileparm := strings.Split(h.Filename, ".")
	filename := fmt.Sprintf("%s_%s.%s", fileparm[0], utils.GetTime(), fileparm[1])
	fileDir := "file"
	filePath := strings.Join([]string{fileDir, filename}, string(os.PathSeparator))
	logger.Infof("upload save to [%s]",  filePath)

	// 没有文件夹要先创建
	_, err = os.Stat(fileDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(fileDir, os.ModePerm); err != nil {
			logger.Errorf("create err [%s]", err.Error())
			self.ajaxMsg("生成文件夹失败", MSG_ERR)
			return ""
		}
	}

	if err = self.SaveToFile("file", filePath); err != nil {
		logger.Errorf("save err [%s]", err.Error())
	}
	return filePath
}

// 文件格式检查,可拓展检查内容
func (c *EmployeeController) FileCheck(filename string) (bool, string) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		log.Fatalln("err:", err.Error())
	}

	// 文件格式检查
	sheet := xlFile.Sheets[0]

	for k, row := range sheet.Rows {
		if k == 0 {
			continue
		}
		for _, cell := range row.Cells {
			if len(cell.Value) == 0 {
				return false, "必填项为空"
			}
		}
	}
	return  true, "表格匹配成功"

}

// 离职员工删除接口
func (self *EmployeeController)Delete() {
	id := self.GetString("id")

	err := models.Delete(id)
	if err != nil {
		self.ajaxMsg("删除员工失败", MSG_ERR)
	}
	self.ajaxMsg("删除员工成功", MSG_OK)
}