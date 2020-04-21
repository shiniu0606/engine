package servercenter

import (
	command "github.com/shiniu0606/engine/server/command"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)

type ServerHandler struct {
	net.DefMsgHandler
}

func InitParser() net.IParser {
	p := net.NewParser(net.ParserTypePB)

	p.Register(CMD_RIGSTER_SERVER,ACT_RIGSTER_SERVER_REQ,&command.ServerRegisterReq{})
	return p
}

func InitHandler() net.IMsgHandler {
	var handler = &ServerHandler{}

	handler.Register(CMD_TICKING,ACT_TICKING_REQ,registerServer)
	handler.Register(CMD_RIGSTER_SERVER,ACT_RIGSTER_SERVER_REQ,registerServer)

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