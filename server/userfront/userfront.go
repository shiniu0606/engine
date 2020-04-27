package userfront

import (

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
	common "github.com/shiniu0606/engine/server/common"
)

var centerSession * net.ISession


func InitCenterServerClient() bool {
	handler := InitCenterClientHandler()
	parser  := InitCenterClientParser()
	centerSession = net.StartTcpConnect("tcp://:"+base.Itoa(serverconfig.TcpPort),handler,parser)

	if centerSession == nil {
		base.LogError("InitCenterServerClient error")
		return false
	}
	return true
}


func InitUserServer() {

}