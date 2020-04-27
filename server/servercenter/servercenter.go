package servercenter

import (
	jbp "github.com/shiniu0606/engine/core/db"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
	common "github.com/shiniu0606/engine/server/common"
)

//创建表
func CreateDBTable() {
	if jbp.GetDB().HasTable(&common.Server{}) {
		jbp.GetDB().AutoMigrate(&common.Server{})
	}else {
		jbp.GetDB().Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&common.Server{})
	}
	
	if jbp.GetDB().HasTable(&common.Account{}) {
        jbp.GetDB().AutoMigrate(&common.Account{})
    } else {
		jbp.GetDB().Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&common.Account{})
    }
	
	jbp.GetDB().Model(&common.Account{}).AddUniqueIndex("idx_user_id", "user_id")
	jbp.GetDB().Model(&common.Account{}).AddIndex("idx_acc_name", "acc_name")
}

func InitServer() {
	handler := InitHandler()
	parser  := InitParser()
	net.StartTcpServer("tcp://:"+base.Itoa(serverconfig.TcpPort),handler,parser)
}