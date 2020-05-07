package userfront

import (
	//command "github.com/shiniu0606/engine/server/command"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)

type UserHandler struct {
	net.DefMsgHandler
}

func (r *UserHandler) OnStartHandle(session net.ISession) bool {
	base.LogInfo("user front tcp server start")
	return true
}

func (r *UserHandler) OnCloseHandle(session net.ISession) {
	base.LogInfo("user front tcp closed")
}

func (r *UserHandler) OnConnectCompleteHandle(session net.ISession, ok bool) bool {
	base.LogInfo("user front tcp client connect complete")
	return true
}

func InitUserParser() net.IParser {
	p := net.NewParser(net.ParserTypePB)
	return p
}

func InitUserHandler() net.IMsgHandler {
	var handler = &UserHandler{}
	return handler
}