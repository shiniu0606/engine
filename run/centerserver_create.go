package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	server "github.com/shiniu0606/engine/server/servercenter"
	base "github.com/shiniu0606/engine/core/base"
)

var stopChanForSys = make(chan os.Signal, 1)

func main() {
	var cofpath string
	
	flag.StringVar(&cofpath, "cofpath", "", "配置文件路径")
	flag.Parse()

	server.InitConfig(cofpath)
	base.LogInfo("InitConfig ok")
	server.InitLog()
	base.LogInfo("InitLog ok")
	server.InitDB()
	base.LogInfo("InitDB ok")

	server.InitServer()

	signal.Notify(stopChanForSys, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-stopChanForSys:
		base.LogInfo("servercenter stop")
	}
}