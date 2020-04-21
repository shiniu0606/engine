package common

import (
	
)

type ServerType int

const (
	ServerTypeUserFront ServerType = iota 	//登录服务
	ServerTypeUserCenter                   	//游戏大厅服务
	ServerTypeGameCell                 		//游戏交互场景服务单元
	ServerTypeServerCenter					//服务管理中心
)

type Server struct {
	UID 			int64  			`gorm:"primary_key;AUTO_INCREMENT;column:uid"`
	ServerId  		int  			`gorm:"column:server_id;"`
	ServerName      string			`gorm:"column:server_name;size:30;"`
	ServerType      ServerType 		`gorm:"column:server_type;"`
	IsEnable		bool			`gorm:"column:is_enabled;"`  //是否开放
	OnlineNum		int 			`gorm:"column:online_num;"`  //在线人数
	ServerAddr      string  		`gorm:"column:server_addr;size:255;"`
}

func (Server) TableName() string {
	return "servers"
}
