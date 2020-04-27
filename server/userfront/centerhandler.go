package servercenter

import (
	command "github.com/shiniu0606/engine/server/command"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)

type ServerClientHandler struct {
	net.DefMsgHandler
}

func (r *ServerClientHandler) OnCloseHandle(session net.ISession) bool {
	base.LogInfo("center tcp client close")
	return true
}

func (r *ServerClientHandler) OnConnectCompleteHandle(session ISession, ok bool) bool {
	req := &command.ServerRegisterReq{}
	req.Serverid = serverconfig.ServerId
	req.Servertype = ServerTypeUserFront
	req.Servername = serverconfig.ServerName
	req.Serveraddr = serverconfig.RemoteAddr
	req.Isopen = true
	//注册到服务中心
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

	handler.Register(command.CMD_TICKING,command.ACT_TICKING_RESP,pingServer)
	handler.Register(command.CMD_RIGSTER_SERVER,command.ACT_RIGSTER_SERVER_RESP,registerServer)

	return handler
}

func pingServer(session net.ISession, msg *net.Message) bool {
	base.LogInfo("pingServer:%v",msg)
	
	return true
}

func registerServer(session net.ISession, msg *net.Message) bool {
	base.LogInfo("registerServer:%v",msg)

	return true
}