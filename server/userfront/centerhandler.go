package userfront

import (
	command "github.com/shiniu0606/engine/server/command"
	common "github.com/shiniu0606/engine/server/common"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)

type ServerClientHandler struct {
	net.DefMsgHandler
}

func (r *ServerClientHandler) OnCloseHandle(session net.ISession) {
	base.LogInfo("center server tcp client closed")

	//reconnect
	ConnectCenterServerScheduler()
}

func (r *ServerClientHandler) OnConnectCompleteHandle(session net.ISession, ok bool) bool {
	req := &command.ServerRegisterReq{}
	req.Serverid = serverconfig.ServerId
	req.Servertype = int32(common.ServerTypeUserFront)
	req.Servername = serverconfig.ServerName
	req.Serveraddr = serverconfig.RemoteAddr
	req.Isopen = true
	//register userfront to centerserver
	msg := net.NewMsg(2,2,net.FlagNorlmal,session.GetParser().Pack(req))
	session.Send(msg)
	base.LogInfo("center tcp client connect complete")
	return true
}

func InitCenterClientParser() net.IParser {
	p := net.NewParser(net.ParserTypePB)

	p.Register(command.CMD_RIGSTER_SERVER,command.ACT_RIGSTER_SERVER_RESP,&command.ServerRegisterResp{})
	return p
}

func InitCenterClientHandler() net.IMsgHandler {
	var handler = &ServerClientHandler{}

	handler.Register(command.CMD_TICKING,command.ACT_TICKING_RESP,pongServer)
	handler.Register(command.CMD_RIGSTER_SERVER,command.ACT_RIGSTER_SERVER_RESP,registerServer)

	return handler
}

func pongServer(session net.ISession, msg *net.Message) bool {
	base.LogInfo("pongServer:%v",msg)
	
	return true
}

func registerServer(session net.ISession, msg *net.Message) bool {
	base.LogInfo("registerServer success:%v",msg)

	return true
}