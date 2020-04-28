package userfront

import (
	"time"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
	//common "github.com/shiniu0606/engine/server/common"
)

var (
	centerSession 	net.ISession
	centerScheduler *base.TimerScheduler

	userSession 	net.ISession
)

func GetCenterSession() net.ISession {
	return centerSession
}

func GetCenterSchedluer() *base.TimerScheduler {
	return centerScheduler
}

func InitCenterServerClient() bool {
	handler := InitCenterClientHandler()
	parser  := InitCenterClientParser()
	centerSession = net.StartTcpConnect(serverconfig.CenterAddr,handler,parser)

	if centerSession == nil {
		base.LogError("InitCenterServerClient error")
		return false
	}
	base.LogInfo("InitCenterServerClient success")
	return true
}

func ConnectCenterServerScheduler() {
	if centerScheduler == nil {
		centerScheduler = base.NewAutoExecTimerScheduler()
	}

	tid, err := centerScheduler.NewTimerAfter(time.Duration(3) * time.Second, CenterServerScheduler, []interface{}{SchedulerTypeConnectCenterServer})
	if err != nil {
		base.LogInfo("InitScheduler error", tid, err)
	}
}


func InitUserFrontServer() bool {
	handler := InitUserHandler()
	parser  := InitUserParser()
	net.StartTcpServer("tcp://:"+base.Itoa(serverconfig.TcpPort),handler,parser)

	ConnectCenterServerScheduler()
	return true
}