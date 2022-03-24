package log

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type FileLog struct {
	LogLevel       int
	FilePath       string
	FileName       string
	TargetLog      *os.File
	TargetLogWn    *os.File
	LogMessageChan chan *LogMessage
	wg             sync.WaitGroup
	LenOfChan      int
	SpiltType      string
	Hour           string
	Size           int64
	LastFileName   string
	LastFileWnName string
	timetot        int64
}

func CreatFileLog(level, lenofchan int, path, name string) Loginterface {
	var alog = FileLog{}
	alog.LogLevel = level
	alog.FileName = name
	alog.FilePath = path
	alog.LenOfChan = lenofchan
	alog.ReadConfig()
	err := alog.init()
	if err != nil {
		panic(err)
	}
	return &alog
}

func (c *FileLog) ReadConfig() {
	c.SpiltType = Config["split-type"]
	c.Hour = Config["hour"]
	ss, _ := strconv.Atoi(Config["size"])
	c.Size = int64(ss)
}

func (c *FileLog) SetLevel(level int) error {
	if LevelDebug > level || level > LevelFatal {
		return errors.New("level out of limit")
	}
	c.LogLevel = level
	return nil
}
func (c *FileLog) init() error {
	logaddr := fmt.Sprintf("%s%s.log", c.FilePath, c.FileName)
	logfile, err := os.OpenFile(logaddr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	logaddrwn := fmt.Sprintf("%s%s.log.wn", c.FilePath, c.FileName)
	logfilewn, err := os.OpenFile(logaddrwn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	c.LastFileName = "nil"
	c.LastFileWnName = "nil"
	c.TargetLog = logfile
	c.TargetLogWn = logfilewn
	c.timetot = 0
	c.LogMessageChan = make(chan *LogMessage, c.LenOfChan)
	go c.RunLog()
	c.wg.Add(1)
	return nil
}

func (c *FileLog) Writelog(level int, format string, args ...interface{}) {
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
	logmsg := CreatLogMessage(level, endtext)
	if c.checkLogMsgChan() {
		c.LogMessageChan <- logmsg
	} else {
		fmt.Println("chan is too small")
	}
}

func (c *FileLog) LogDebug(format string, args ...interface{}) {
	c.Writelog(LevelDebug, format, args...)
}

func (c *FileLog) LogTrace(format string, args ...interface{}) {
	c.Writelog(LevelTrace, format, args...)

}
func (c *FileLog) LogInfo(format string, args ...interface{}) {
	c.Writelog(LevelInfo, format, args...)
}

func (c *FileLog) LogWarn(format string, args ...interface{}) {
	c.Writelog(LevelWarn, format, args...)

}
func (c *FileLog) LogError(format string, args ...interface{}) {
	c.Writelog(LevelError, format, args...)

}

func (c *FileLog) LogFatal(format string, args ...interface{}) {
	c.Writelog(LevelFatal, format, args...)

}
func (c *FileLog) Close() error {
	c.wg.Wait()
	err := c.TargetLogWn.Close()
	if err != nil {
		return err
	}
	err = c.TargetLog.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *FileLog) RunLog() {
	for i := range c.LogMessageChan {
		c.CheckSpilt()
		c.WriteInFile(i)
		if len(c.LogMessageChan) == 0 {
			c.wg.Done()
		}
	}
}

func (c *FileLog) CheckSpilt() {
	if c.SpiltType == "nil" {
		return
	} else if c.SpiltType == "time" {
		nowtime := time.Now()
		timetest := nowtime.Unix()
		timestr := fmt.Sprintf("%04d%02d%02d%02d", nowtime.Year(), nowtime.Month(), nowtime.Day(), nowtime.Hour())
		h, _ := strconv.Atoi(c.Hour)
		if c.timetot == 0 || c.timetot-timetest >= int64(h)*3600 {
			logaddr := fmt.Sprintf("%s%s.log"+timestr, c.FilePath, c.FileName)
			logfile, err := os.OpenFile(logaddr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				panic("file open error ")
			}
			logaddrwn := fmt.Sprintf("%s%s.log.wn"+timestr, c.FilePath, c.FileName)
			logfilewn, err := os.OpenFile(logaddrwn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				panic("file open error ")
			}
			_ = c.TargetLog.Close()
			_ = c.TargetLogWn.Close()
			c.TargetLog = logfile
			c.TargetLogWn = logfilewn
			c.LastFileName = timestr
		} else if c.timetot-timetest < int64(h)*3600 {
			return
		}
	} else if c.SpiltType == "size" {
		nowtime := time.Now()
		timestr := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", nowtime.Year(), nowtime.Month(), nowtime.Day(), nowtime.Hour(), nowtime.Minute(), nowtime.Second())
		filestat, err := c.TargetLog.Stat()
		if err != nil {
			panic(err)
		}
		filestatwn, err := c.TargetLogWn.Stat()
		if err != nil {
			panic(err)
		}
		if c.LastFileName == "nil" || filestat.Size() >= c.Size {
			logaddr := fmt.Sprintf("%s%s.log"+timestr, c.FilePath, c.FileName)
			logfile, _ := os.OpenFile(logaddr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			_ = c.TargetLog.Close()
			c.TargetLog = logfile
			c.LastFileName = timestr
		}
		if c.LastFileWnName == "nil" || filestatwn.Size() >= c.Size {
			logaddr := fmt.Sprintf("%s%s.log.wn"+timestr, c.FilePath, c.FileName)
			logfile, _ := os.OpenFile(logaddr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			_ = c.TargetLogWn.Close()
			c.TargetLogWn = logfile
			c.LastFileWnName = timestr
		}

	} else {
		_ = c.HighClose()
		panic("SpiltType should is \"time\",\"size\" or \"nil\"")
	}
}

func (c *FileLog) WriteInFile(i *LogMessage) {
	if i.Warnflag == false {
		_, err := fmt.Fprintln(c.TargetLog, i.Message)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := fmt.Fprintln(c.TargetLog, i.Message)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Fprintln(c.TargetLogWn, i.Message)
		if err != nil {
			panic(err)
		}
	}
}

func (c *FileLog) checkLogMsgChan() bool {
	return len(c.LogMessageChan) < cap(c.LogMessageChan)
}

func (c *FileLog) HighClose() error {
	err := c.TargetLogWn.Close()
	if err != nil {
		return err
	}
	err = c.TargetLog.Close()
	if err != nil {
		return err
	}
	return nil
}
