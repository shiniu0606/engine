package base

import (
	"os"
	"fmt"
	"sync/atomic"
	"path"
	"path/filepath"
	"strings"
	"runtime"
)

type ILogger interface {
	Write(msg string)
}

type ConsoleLogger struct {
	Pln bool
}

type FileLogger struct {
	FilePath 	string
	Pln			bool
	MaxSize		int

	size 		int
	writer		*os.File
	filename	string
	extname  string
	dirname  string
 }

func (l *ConsoleLogger) Write(msg string) {
	if l.Pln {
		fmt.Println(msg)
	} else {
		fmt.Print(msg)
	}
}

func (l *FileLogger) Write(msg string) {
	if l.writer == nil {
		return
	}

	addsize := l.size + len(msg)
	if l.Pln {
		addsize += 1
	}

	if l.MaxSize > 0 && addsize >= l.MaxSize {
		l.writer.Close()
		l.writer = nil

		newpath := l.dirname + "/" + l.filename + fmt.Sprintf("_%v", GetDate()) + l.extname
		os.Rename(l.FilePath,newpath)
		file,err := os.OpenFile(l.FilePath,os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			l.writer = nil
			return
		}
		l.writer = file
		l.size = 0
	}

	if l.Pln {
		l.writer.WriteString(msg+"\n")
	}else {
		l.writer.WriteString(msg)
	}

	l.size = addsize
}

type LogLevel int

const (
	LogLevelDebug  LogLevel = iota //调试信息
	LogLevelInfo                   //资讯讯息
	LogLevelWarn                   //警告状况发生
	LogLevelError                  //一般错误，可能导致功能不正常
	LogLevelFatal                  //严重错误，会导致进程退出
	LogLevelAllOff                 //关闭所有日志
)

var LogLevelNameMap = map[string]LogLevel{
	"debug": LogLevelDebug,
	"info":  LogLevelInfo,
	"warn":  LogLevelWarn,
	"error": LogLevelError,
	"fatal": LogLevelFatal,
	"off":   LogLevelAllOff,
}

type Log struct {
	logger         [32]ILogger
	cwrite         chan string
	cstop		   chan bool
	bufsize        int
	stop           int32
	preLoggerCount int32
	loggerCount    int32
	level          LogLevel
}

func NewLog(bufsize int, logger ...ILogger) *Log {
	log := &Log{
		bufsize:        bufsize,
		cwrite:         make(chan string, bufsize),
		cstop:			make(chan bool),
		level:          LogLevelDebug,
		preLoggerCount: -1,
	}
	for _, l := range logger {
		log.SetLogger(l)
	}
	return log
}

func (r *Log) SetLogger(logger ILogger) bool {
	if r.preLoggerCount >= 31 {
		return false
	}
	if f, ok := logger.(*FileLogger); ok {
		if r.initFileLogger(f) == nil {
			return false
		}
	}
	r.logger[atomic.AddInt32(&r.preLoggerCount, 1)] = logger
	atomic.AddInt32(&r.loggerCount, 1)
	return true
}

func (r *Log) initFileLogger(f *FileLogger) *FileLogger {
	if f.writer == nil {
		f.FilePath, _ = filepath.Abs(f.FilePath)
		f.FilePath = strings.Replace(f.FilePath, "\\", "/",-1)
		f.dirname = path.Dir(f.FilePath)
		f.extname = path.Ext(f.FilePath)
		f.filename = filepath.Base(f.FilePath[0 : len(f.FilePath)-len(f.extname)])
		os.MkdirAll(f.dirname, 0666)
		file, err := os.OpenFile(f.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err == nil {
			f.writer = file
			info, err := f.writer.Stat()
			if err != nil {
				return nil
			}

			f.size = int(info.Size())

			return f
		}
	}
	return nil
}

func (r *Log) Stop() {
	if atomic.CompareAndSwapInt32(&r.stop, 0, 1) {
		close(r.cwrite)
		close(r.cstop)
	}
}

func (r *Log) IsStop() bool {
	return r.stop == 1
}

func (r *Log) StartLoop() {
	Go(func (){
		var i int32
		for {
			select {
			case s,ok := <-r.cwrite:
				if ok {
					for i = 0; i < r.loggerCount; i++ {
						r.logger[i].Write(s)
					}
				}
			case <-r.cstop:
				for i = 0; i < r.loggerCount; i++ {
					if f, ok := r.logger[i].(*FileLogger); ok {
						if f.writer != nil {
							f.writer.Close()
							f.writer = nil
						}
					}
				}
				return
			}
		}
	})
}

func (r *Log) GetLevel() LogLevel {
	return r.level
}
func (r *Log) SetLevel(level LogLevel) {
	r.level = level
}

func (r *Log) SetLevelByName(name string) bool {
	level, ok := LogLevelNameMap[name]
	if ok {
		r.SetLevel(level)
	}
	return ok
}

func (r *Log) write(levstr string, v ...interface{}) {
	defer func() { recover() }()
	if r.IsStop() {
		return
	}
	prefix := levstr
	_, file, line, ok := runtime.Caller(3)
	if ok {
		i := strings.LastIndex(file, "/") + 1
		prefix = fmt.Sprintf("[%s][%s][%s:%d]:", levstr, GetDate(), (string)(([]byte(file))[i:]), line)
	}
	if len(v) > 1 {
		r.cwrite <- prefix + fmt.Sprintf(v[0].(string), v[1:]...)
	} else {
		r.cwrite <- prefix + fmt.Sprint(v[0])
	}
}

func (r *Log) Debug(v ...interface{}) {
	if r.level <= LogLevelDebug {
		r.write("D", v...)
	}
}

func (r *Log) Info(v ...interface{}) {
	if r.level <= LogLevelInfo {
		r.write("I", v...)
	}
}

func (r *Log) Warn(v ...interface{}) {
	if r.level <= LogLevelWarn {
		r.write("W", v...)
	}
}

func (r *Log) Error(v ...interface{}) {
	if r.level <= LogLevelError {
		r.write("E", v...)
	}
}

func (r *Log) Fatal(v ...interface{}) {
	if r.level <= LogLevelFatal {
		r.write("FATAL", v...)
	}
}

func (r *Log) Write(v ...interface{}) {
	defer func() { recover() }()
	if r.IsStop() {
		return
	}

	if len(v) > 1 {
		r.cwrite <- fmt.Sprintf(v[0].(string), v[1:]...)
	} else if len(v) > 0 {
		r.cwrite <- fmt.Sprint(v[0])
	}
}