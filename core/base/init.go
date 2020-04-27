package base

import (
	"runtime"
	"time"
	"math/rand"
	"sync/atomic"
)

type GlobalObj struct {
	MaxMsgDataSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度
	MsgSessionId     uint64 //消息队列id
}


var GlobalObject *GlobalObj

func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		MaxConn:          12000,
		MaxMsgDataSize:    40960,
		MaxMsgChanLen:    64,	
		MsgSessionId:     1,	
	}
	//随机种子
	rand.Seed(time.Now().UnixNano())
	//多核
	runtime.GOMAXPROCS(runtime.NumCPU())
	//内置timer
	timerTick()
}

func GetGlobal() *GlobalObj{
	return GlobalObject
}

func GetMsgSessionId() uint64 {
	return atomic.AddUint64(&GlobalObject.MsgSessionId, 1)
}

func GetMaxMsgChanLen() uint32 {
	return GlobalObject.MaxMsgChanLen
}

