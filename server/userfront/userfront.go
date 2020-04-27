package userfront

import (

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
	common "github.com/shiniu0606/engine/server/common"
)

//连接centerserver
func InitCenterServerClient() {
	handler := InitCenterClientHandler()
	parser  := InitCenterClientParser()
	net.StartTcpConnect("tcp://:"+base.Itoa(serverconfig.TcpPort),handler,parser)
}

//启动userfront
func InitUserServer() {

}