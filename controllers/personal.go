package controllers

import (
	"geek-nebula/models"
	"reflect"
)

type PersonalController struct {
	BaseController
}

func (self *PersonalController) Index() {
	self.Data["pageTitle"] = "个人报表"
	//self.display()
	self.TplName = "personal/index.html"
	self.display()
}

func (self *PersonalController) List() {
	var params models.BaseQueryParam
	params.Page, _ = self.GetInt("page")
	params.Limit, _ = self.GetInt("limit")
	Date := self.GetString("month")

	employeeList, totalNum, err := models.PersonalReportQuery(&params, Date)
	if err != nil {
		logger.ErrorById(self.userId, "%v", err)
		self.ajaxMsg("查询个人报表信息失败", MSG_ERR)
		return
	}

	list := make([]interface{}, len(employeeList))
	for k, employee := range employeeList {
		v := reflect.ValueOf(&employee).Elem()
		row := make(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			row[v.Type().Field(i).Tag.Get("json")] = v.Field(i).Interface()
		}
		list[k] = row
	}

	self.ajaxList("成功", MSG_OK, totalNum, list)
	self.TplName = "personal/index.html"
}