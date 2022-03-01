package base

import (
	"sync/atomic"
)

var gocount int32 //goroutine数量
var goid uint32

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			LogStack()
			LogError("error catch:%v", err)
			if handler != nil {
				handler(err)
			}
		}
	}()
	fun()
}

func Go(fn func()) {
	var debugStr string
	id := atomic.AddUint32(&goid, 1)
	c := atomic.AddInt32(&gocount, 1)
	if DefLog.GetLevel() <= LogLevelDebug {
		debugStr = LogSimpleStack()
		LogDebug("goroutine start id:%d count:%d from:%s", id, c, debugStr)
	}

	go func() {
		Try(fn, nil)

		if DefLog.GetLevel() <= LogLevelDebug {
			c = atomic.AddInt32(&gocount, -1)
			LogDebug("goroutine end id:%d count:%d from:%s", id, c, debugStr)
		}
	}()
}
