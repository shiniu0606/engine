package main

import (
	"os"
	"os/signal"
	"syscall"
	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)
var stopChanForSys = make(chan os.Signal, 1)

type ServerHandler struct {
	net.DefMsgHandler
}

type ClientHandler struct {
	net.DefMsgHandler
}

func (r *ServerHandler) OnProcessMsgHandle(session net.ISession, msg *net.Message) bool {
	session.Send(msg)
	base.LogInfo("ServerHandler ProcessMsgHandle :%v",string(msg.Data))
	return true
}

func (r *ServerHandler) OnCloseHandle(session net.ISession)  {
	base.LogInfo("ServerHandler OnCloseHandle")
}

func (r *ClientHandler) OnProcessMsgHandle(session net.ISession, msg *net.Message) bool {
	//session.Send(msg)
	base.LogInfo("ClientHandler ProcessMsgHandle :%v",string(msg.Data))
	session.Stop()
	return true
}

var sh = ServerHandler{}

func main() {

	go net.StartTcpServer("tcp://:6666",&sh,nil)

	n := 10
	for i := 0; i < n; i++ {
		testTcpClient()
	}

	for i := 0; i < n; i++ {
		testTcpClientClose()
	}

	signal.Notify(stopChanForSys, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-stopChanForSys:
		base.LogInfo("sys stop")
	}
}

func testTcpClient() {
	h := &ClientHandler{}
	session := net.StartTcpConnect("127.0.0.1:6666",h,nil)

	if session != nil {
		msg := net.NewMsg(2,2,net.FlagNorlmal,[]byte("hello word"))
		session.Send(msg)
	}
}

func testTcpClientClose() {
	
	h := &ClientHandler{}
	session := net.StartTcpConnect("127.0.0.1:6666",h,nil)
	defer session.Stop()
	if session != nil {
		msg := net.NewMsg(2,2,net.FlagNorlmal,[]byte("hello word"))
		session.Send(msg)
	}
}