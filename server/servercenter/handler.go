package servercenter

import (
	command "github.com/shiniu0606/engine/server/command"
	common "github.com/shiniu0606/engine/server/common"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)

type ServerHandler struct {
	net.DefMsgHandler
}

func (r *ServerHandler) OnStartHandle(session net.ISession) bool {
	base.LogInfo("center tcp server start")

	serversession := ServerSession{
		conn : session,
		server : common.Server{
			ServerId : serverconfig.ServerId,
			ServerName: serverconfig.ServerName,
			ServerType: common.ServerTypeServerCenter,
			IsEnable: true,
			ServerAddr: serverconfig.RemoteAddr,
		},
	}

	common.CreateServer(&serversession.server)
	NewCenterServer(&serversession)
	return true
}

func InitParser() net.IParser {
	p := net.NewParser(net.ParserTypePB)

	p.Register(command.CMD_RIGSTER_SERVER,command.ACT_RIGSTER_SERVER_REQ,&command.ServerRegisterReq{})
	return p
}

func InitHandler() net.IMsgHandler {
	var handler = &ServerHandler{}

	handler.Register(command.CMD_TICKING,command.ACT_TICKING_REQ,pingServer)
	handler.Register(command.CMD_RIGSTER_SERVER,command.ACT_RIGSTER_SERVER_REQ,registerServer)

	return handler
}

func pingServer(session net.ISession, msg *net.Message) bool {
	base.LogInfo("pingServer:%v",msg)
	
	return true
}

func registerServer(session net.ISession, msg *net.Message) bool {
	req := msg.UserData.(*command.ServerRegisterReq)

	serversession := ServerSession{
		conn : session,
		server : common.Server{
			ServerId : int(req.Serverid),
			ServerName: req.Servername,
			ServerType: common.ServerType(req.Servertype),
			IsEnable: req.Isopen,
			ServerAddr: req.Serveraddr,
		},
	}

	GetCenterServer().Add(&serversession)

	//base.LogInfo("registerServer:%v",req)
	return true
}