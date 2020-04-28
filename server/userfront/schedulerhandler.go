package userfront

import (
	base "github.com/shiniu0606/engine/core/base"
)

type SchedulerType int

const (
	SchedulerTypeConnectCenterServer 		SchedulerType = iota 	//连接中心服务
	SchedulerTypeCenterServerPing                   				//心跳
)

func CenterServerScheduler(args ...interface{}) {
	base.LogInfo("CenterServerScheduler type :%d",args[0].(SchedulerType))
	switch args[0].(SchedulerType) {
	case SchedulerTypeConnectCenterServer:
		if InitCenterServerClient() == false {
			ConnectCenterServerScheduler()  //reconnect server
		}
	case SchedulerTypeCenterServerPing:

	default:
		base.LogError("CenterServerScheduler type error :%d",args[0].(SchedulerType))
	}
}