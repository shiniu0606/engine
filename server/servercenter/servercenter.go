package servercenter

import (
	jbp "github.com/shiniu0606/engine/core/db"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
	common "github.com/shiniu0606/engine/server/common"
)

func init() {
	InitLog()
	base.LogInfo("InitLog ok")
	InitConfig()
	base.LogInfo("InitConfig ok")
	InitDB()
	base.LogInfo("InitDB ok")
}

//创建表
func CreateDBTable() {
	jbp.GetDB().Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&common.Server{})
}

func InitServer() {
	handler := InitHandler()
	parser  := InitParser()
	net.StartTcpServer("tcp://:"+base.Itoa(serverconfig.TcpPort),handler,parser)
}