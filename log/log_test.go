package log

import (
	"testing"
)

func TestFile(t *testing.T) {
	Config["split-type"] = "time"
	Config["hour"] = "0.00001"
	filelog := CreatFileLog(LevelDebug, 9000000, "F:\\golang_project\\gin\\log\\logtest\\", "logtest")
	//filelog := CreatConsoleLog(LevelDebug)
	defer filelog.Close()
	for {
		filelog.LogWarn("warn warn")
		filelog.LogFatal("LogFatal  LogFatal  LogFatal")
		filelog.LogError("LogError LogError LogError")
		filelog.LogTrace("tttttttttttttttttttttttttt")
		filelog.LogDebug("username:%s debug", "sadfasf")
	}
	//time.Sleep(time.Second * 10)
}

//func TestConsole(t *testing.T) {
//	console := CreatConsoleLog(LevelDebug)
//	console.LogDebug("debug console console consle")
//	console.LogInfo("debug LogInfoLogInfo---------")
//	console.LogFatal("debug conLogFatalLogFatalLogFatalsole consle")
//	console.LogError("debug consErrore consle")
//	console.LogWarn("debwa----------------rnconsle")
//}
