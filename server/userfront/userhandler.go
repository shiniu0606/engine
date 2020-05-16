package userfront

import (
	command "github.com/shiniu0606/engine/server/command"
	common "github.com/shiniu0606/engine/server/common"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)

type UserHandler struct {
	net.DefMsgHandler
}

func (r *UserHandler) OnStartHandle(session net.ISession) bool {
	base.LogInfo("user front tcp server start")

	userSession = session
	return true
}

func (r *UserHandler) OnCloseHandle(session net.ISession) {
	base.LogInfo("user front tcp closed")

	userSession = nil
}

func (r *UserHandler) OnConnectCompleteHandle(session net.ISession, ok bool) bool {
	base.LogInfo("user front tcp client connect complete")
	return true
}

func InitUserParser() net.IParser {
	p := net.NewParser(net.ParserTypePB)

	p.Register(command.CMD_USER_MAIN,command.ACT_USER_REGISTER_REQ,&command.AccountRegisterReq{})
	return p
}

func InitUserHandler() net.IMsgHandler {
	var handler = &UserHandler{}

	handler.Register(command.CMD_USER_MAIN,command.ACT_USER_REGISTER_REQ,userRegister)
	handler.Register(command.CMD_USER_MAIN,command.ACT_USER_LOGIN_REQ,userLogin)

	return handler
}

func userRegister(session net.ISession, msg *net.Message) bool {
	req := msg.UserData.(*command.AccountRegisterReq)

	Accname := req.Accname
	//Accpassword := base.MD5WithSalt(req.Accpassword)

	acc,err := common.GetAccountByAccountName(Accname)
	if(err != nil){
		base.LogInfo("userRegister err :%v",err)
		return false
	}

	if(acc != nil){
		
		return true
	}

	return true
}

func userLogin(session net.ISession, msg *net.Message) bool {
	base.LogInfo("userRegister:%v",msg)
	
	return true
}