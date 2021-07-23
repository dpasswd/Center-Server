package logger

import (
	"dh-passwd/global"
	"dh-passwd/tools"
	"dh-passwd/tools/config"
	"github.com/gogf/gf/os/glog"
)

var Logger *glog.Logger
var RequestLogger *glog.Logger

func Setup() {
	Logger = glog.New()
	_ = Logger.SetPath(config.LoggerConfig.Path + "/dh-passwd")
	Logger.SetStdoutPrint(config.LoggerConfig.EnabledWeb && config.LoggerConfig.Stdout)
	Logger.SetFile("dh-passwd-{Ymd}.log")
	_ = Logger.SetLevelStr(config.LoggerConfig.Level)

	RequestLogger = glog.New()
	_ = RequestLogger.SetPath(config.LoggerConfig.Path + "/request")
	RequestLogger.SetStdoutPrint(false)
	RequestLogger.SetFile("access-{Ymd}.log")
	_ = RequestLogger.SetLevelStr(config.LoggerConfig.Level)

	Logger.Info(tools.Green("Logger init success!"))

	global.Logger = Logger.Line()
	global.RequestLogger = RequestLogger.Line()
}
