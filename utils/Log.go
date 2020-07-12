package utils

import (
	"bytes"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strconv"
)

type taskLog struct{
	log *logs.BeeLogger
}

var Logger *taskLog

func init() {
	Logger = new(taskLog)
	Logger.log = logs.NewLogger()

	Logger.log.EnableFuncCallDepth(true)
	Logger.log.SetLogFuncCallDepth(3)

	Logger.log.SetLogger("console", beego.AppConfig.DefaultString("LogConsoleConf", ""))
	Logger.log.SetLogger("file", beego.AppConfig.DefaultString("LogFileConf", `{"filename":"logs/employee_mng.log","daily":true}`))
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func (log *taskLog) Infof(f interface{}, v ...interface{}) {
	t := append([]interface{}{GetGID()}, v...)
	log.log.Info("[%d] "+f.(string), t...)
}

func (log *taskLog) Errorf(f interface{}, v ...interface{}) {
	t := append([]interface{}{GetGID()}, v...)
	log.log.Error("[%d] "+f.(string), t...)
}

func (log *taskLog) InfoById(id int, f interface{}, v ...interface{}) {
	t := append([]interface{}{GetGID(), id}, v...)
	log.log.Info("[%d] [id=%d] "+f.(string), t...)
}

func (log *taskLog) ErrorById(id int, f interface{}, v ...interface{}) {
	t := append([]interface{}{GetGID(), id}, v...)
	log.log.Error("[%d] [id=%d] "+f.(string), t...)
}
