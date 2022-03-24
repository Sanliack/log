package log

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

type ConsoleLog struct {
	LogLevel  int
	TargetLog *os.File
}

func CreatConsoleLog(level int) Loginterface {
	var alog = ConsoleLog{}
	alog.LogLevel = level
	err := alog.init()
	if err != nil {
		panic(err)
	}
	return &alog
}
func (c *ConsoleLog) SetLevel(level int) error {
	if LevelDebug > level || level > LevelFatal {
		return errors.New("level out of limit")
	}
	c.LogLevel = level
	return nil
}
func (c *ConsoleLog) init() error {
	logfile := os.Stdout
	c.TargetLog = logfile
	return nil
}

func (c *ConsoleLog) Writelog(level int, format string, args ...interface{}) {
	if c.LogLevel > level {
		return
	}
	appeartime := time.Now().Format("2006-01-02 15:04:05")
	logtext := fmt.Sprintf(format, args...)
	_, file, line, ok := runtime.Caller(1)
	file = path.Base(file)
	if ok == false {
		panic("ok error-----------------------")
	}
	appearindex := fmt.Sprintf("%s:%d", file, line)
	endtext := fmt.Sprintf("%s: %s (%s) : %s", appeartime, Levelmap[level], appearindex, logtext)
	_, _ = fmt.Fprintln(c.TargetLog, endtext)

}

func (c *ConsoleLog) LogDebug(format string, args ...interface{}) {
	c.Writelog(LevelDebug, format, args...)
}

func (c *ConsoleLog) LogTrace(format string, args ...interface{}) {
	c.Writelog(LevelTrace, format, args...)
}
func (c *ConsoleLog) LogInfo(format string, args ...interface{}) {
	c.Writelog(LevelInfo, format, args...)
}
func (c *ConsoleLog) LogWarn(format string, args ...interface{}) {
	c.Writelog(LevelWarn, format, args...)
}
func (c *ConsoleLog) LogError(format string, args ...interface{}) {
	c.Writelog(LevelError, format, args...)
}

func (c *ConsoleLog) LogFatal(format string, args ...interface{}) {
	c.Writelog(LevelFatal, format, args...)
}
func (c *ConsoleLog) Close() error {
	return nil
}
