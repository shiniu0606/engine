package main

import (
	"flag"
	server "github.com/shiniu0606/engine/server/servercenter"
	//base "github.com/shiniu0606/engine/core/base"
)

func main() {
	var cofpath string
	
	flag.StringVar(&cofpath, "cofpath", "", "配置文件路径")
	flag.Parse()

	server.InitConfig(cofpath)
	server.InitLog()
	server.InitDB()

	server.CreateDBTable()
}