有点小bug没修 建议不用


# log
golang log logger日志管理

实现了对日志的文件写入或者终端输出2种功能。

文件写入：

var Config = map[string]string{
	"split-type": "nil",
	"hour":       "24",
	"size":       "104857600",
}

log.Config["split-type"] = "nil"/"time"/"size"  //指定文件写入是否切分 ：nil不切分，size按文件大小切分（默认大于100mb切分一次），time按照时间切分（默认24h切分一次）
log.Config["hour"] = "1" // 配合log.Config["split-type"] ="time" 使用

log.Config["size"] = "104857600"  // 配合log.Config["split-type"] ="size" 使用


快速开始（文件）:

filelog := log.CreatFileLog(log.LevelDebug, 9000000, "F:\\golang_project\\gin\\log\\logtest\\", "logtest")
// 参数：日志记录级别，chan的队列数长度，文件保存路径，文件名字
defer filelog.Close()
filelog.LogDebug("log a Debug")
filelog.LogError("log a Error")


快速开始（终端）：

filelog := log.CreatConsoleLog(log.LevelDebug)
// 参数：日志记录级别
defer filelog.Close()
filelog.LogDebug("log a Debug")
filelog.LogError("log a Error")
