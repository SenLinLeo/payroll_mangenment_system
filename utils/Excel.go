package utils

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"reflect"
	"regexp"
)

//将XSLX一行数据赋入结构体中，并进行正则校验
func SetStructField(ptr interface{}, containers []*xlsx.Cell) (bool, string) {
	var count int = 0
	cnt := 0
	v := reflect.ValueOf(ptr).Elem() // the struct variable

	for index, container := range containers {
		v.Field(index).Set(reflect.ValueOf(container.Value))
		count += 1
	}
	selectFlag := 0
	selectField := ""
	for i := 0; i < v.NumField(); i++ {

		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		getMatch := tag.Get("match")
		match := regexp.MustCompile(getMatch)
		pk := tag.Get("pk")

		if getMatch == "" || (pk != "1" && pk != "2") {
			// 跳过匹配
			continue
		}

		if pk == "1" {
			err := match.MatchString(v.Field(i).String())
			if !err {
				return false, fmt.Sprintf(`数据[%s]匹配失败[%s]`, reflect.TypeOf(ptr).Elem().Field(i).Name, v.Field(i).String())
			}
		}
		if pk == "1" && v.Field(i).String() == "" {
			return false, fmt.Sprintf(`必填项为空[%s]`, v.Field(i).String())
		}
		if pk == "2" && v.Field(i).String() == "" && cnt == 0 {
			selectFlag = 1
			selectField += reflect.TypeOf(ptr).Elem().Field(i).Name
		} else if pk == "2" && v.Field(i).String() != "" {
			cnt += 1
			selectFlag = 0
		}
	}
	if selectFlag != 0 {
		return false, fmt.Sprintf(`必填项为空[%s]`, selectField)
	}

	return true, "数据获取成功"
}