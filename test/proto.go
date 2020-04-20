package main

import (
	stproto "github.com/shiniu0606/engine/proto"
	//"github.com/golang/protobuf/proto"

	"os"
	"os/signal"
	"syscall"
	base "github.com/shiniu0606/engine/base"
	net "github.com/shiniu0606/engine/net"
)

var stopChanForSys = make(chan os.Signal, 1)

type ServerHandler struct {
	net.DefMsgHandler
}

type ClientHandler struct {
	net.DefMsgHandler
}

func (r *ServerHandler) OnProcessMsgHandle(session net.ISession, msg *net.Message) bool {
	//session.Send(msg)
	base.LogInfo("ServerHandler ProcessMsgHandle :%v",msg.Head.CmdAct())
	base.LogInfo("ServerHandler ProcessMsgHandle :%v",msg.UserData.(*stproto.UserInfoReq).Message)
	return true
}

func (r *ClientHandler) OnProcessMsgHandle(session net.ISession, msg *net.Message) bool {
	//session.Send(msg)
	base.LogInfo("ClientHandler ProcessMsgHandle :%v",string(msg.Data))
	session.Stop()
	return true
}

var sh = ServerHandler{}

func main() {
	p := net.NewParser(net.ParserTypePB)
	p.Register(2,2,&stproto.UserInfoReq{})
	go net.StartTcpServer("tcp://:6666",&sh,p)

	n := 1
	for i := 0; i < n; i++ {
		testTcpClient()
	}

	signal.Notify(stopChanForSys, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-stopChanForSys:
		base.LogInfo("sys stop")
	}
}

func testTcpClient() {
	h := &ClientHandler{}
	p := net.NewParser(net.ParserTypePB)
	p.Register(2,1,&stproto.UserInfoResp{})
	session := net.StartTcpConnect("127.0.0.1:6666",h,p)

	if session != nil {
		req := &stproto.UserInfoReq{}
		req.Message = "hello proto"
		req.Length = 11
		req.Cnt = 1

		msg := net.NewMsg(2,2,net.FlagNorlmal,session.GetParser().Pack(req))

		//base.LogInfo("============ ========== :%v",msg.Head)
		session.Send(msg)
	}
}