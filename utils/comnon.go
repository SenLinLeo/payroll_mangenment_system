package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Atof(s string) float64 {
	tmp, _ := strconv.ParseFloat(s, 64)
	return tmp
}

func Atoi(s string) int {
	tmp, _ := strconv.Atoi(s)
	return tmp
}


func GetDate() string {
	return time.Now().Format("20060102")
}

func GetTime() string {
	return time.Now().Format("20060102150405")
}

func GetFormatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func CalcDateSeek(t1, t2 string) int {
	t1Num, _ := strconv.Atoi(t1)
	t2Num, _ := strconv.Atoi(t2)

	return t1Num - t2Num
}

func GetSeekDate(seek int, date string) string {
	settleTime, _ := time.ParseInLocation("20060102", date, time.Local)
	settleTime.Add(time.Duration(seek) * time.Hour * 24)

	return settleTime.Format("20060102")
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func FormatTimeStr(t string) string {
	stamp, err := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
	if err != nil {
		fmt.Println(err.Error())
	}

	return stamp.Format("20060102")
}

func ChangeTimeStr(s string) string {
	stamp, _ := time.ParseInLocation("20060102150405", s, time.Local)
	return FormatTime(stamp)
}

func Strip(s string) string {
	str := strings.Replace(s, "\n", "", -1)
	return strings.Replace(str, "\r", "", -1)
}

func ConvertToFormatDay(excelDaysString string)string{
	// 2006-01-02 距离 1900-01-01的天数
	baseDiffDay := 38719 //在网上工具计算的天数需要加2天，什么原因没弄清楚
	curDiffDay := excelDaysString
	b,_ := strconv.Atoi(curDiffDay)
	// 获取excel的日期距离2006-01-02的天数
	realDiffDay := b - baseDiffDay
	//fmt.Println("realDiffDay:",realDiffDay)
	// 距离2006-01-02 秒数
	realDiffSecond := realDiffDay * 24 * 3600
	//fmt.Println("realDiffSecond:",realDiffSecond)
	// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond + realDiffSecond), 0).Format("2006-01-02")
	return resultTime
}

func GetLastMonthFirstDay() string  {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return thisMonth.AddDate(0, -1, 0).Format("20060102")
}

func GetAMonth(year int, month int, day int) string {
	return time.Now().AddDate(year, month, day).Format("20060102")
}
