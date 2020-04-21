package main

import (
	"os"
	"os/signal"
	"syscall"
	"flag"

	"net/url"
	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
)

var addr = flag.String("addr", "localhost:6666", "http service address")


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

func (r *ClientHandler) OnProcessMsgHandle(session net.ISession, msg *net.Message) bool {
	//session.Send(msg)
	base.LogInfo("ClientHandler ProcessMsgHandle :%v",string(msg.Data))
	session.Stop()
	return true
}

var sh = ServerHandler{}

func main() {
	flag.Parse()

	go net.StartWebscoketServer("ws://:6666",&sh,"","")

	n := 10
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
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	base.LogInfo("connecting to %s", u.String())

	h := &ClientHandler{}
	session := net.StartWebsocketConnect(u.String(),h)

	if session != nil {
		msg := net.NewMsg(2,2,net.FlagNorlmal,[]byte("hello word"))
		session.Send(msg)
	}
}