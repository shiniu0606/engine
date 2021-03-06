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
	
	flag.StringVar(&cofpath, "cofpath", "servercenter.json", "配置文件路径")
	flag.Parse()

	server.InitConfig(cofpath)
	server.InitLog()
	server.InitDB()

	server.InitServer()

	signal.Notify(stopChanForSys, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-stopChanForSys:
		base.LogInfo("servercenter stop")
	}
}