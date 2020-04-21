package base

import (
	"runtime"
	"strings"
)

var DefLog *Log //日志

func init() {
	DefLog = NewLog(10000, &ConsoleLogger{true})
	DefLog.SetLevel(LogLevelDebug)
	DefLog.StartLoop()
}

func GetLog() *Log {
	return DefLog
}

func SetFileLog(filepath string,maxsize int) {
	if DefLog == nil {
		DefLog = NewLog(10000, &ConsoleLogger{true})
		DefLog.SetLevel(LogLevelDebug)
		DefLog.StartLoop()
	}
	DefLog.SetLogger(&FileLogger{
		FilePath:		filepath,
		Pln:			true,
		MaxSize:		maxsize,
	})
}

func LogInfo(v ...interface{}) {
	DefLog.Info(v...)
}

func LogDebug(v ...interface{}) {
	DefLog.Debug(v...)
}

func LogError(v ...interface{}) {
	DefLog.Error(v...)
}

func LogFatal(v ...interface{}) {
	DefLog.Fatal(v...)
}

func LogWarn(v ...interface{}) {
	DefLog.Warn(v...)
}

func LogStack() {
	buf := make([]byte, 1<<12)
	LogError(string(buf[:runtime.Stack(buf, false)]))
}

func LogSimpleStack() string {
	_, file, line, _ := runtime.Caller(2)
	i := strings.LastIndex(file, "/") + 1
	i = strings.LastIndex((string)(([]byte(file))[:i-1]), "/") + 1

	return Sprintf("%s:%d", (string)(([]byte(file))[i:]), line)
}