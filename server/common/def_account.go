package common

import (
	"time"
)


type Account struct {
	UID 					int64  			`gorm:"primary_key;AUTO_INCREMENT;column:uid"`
	UserId  				int64  			`gorm:"column:user_id;not null"`				//游戏id(方便扩展为分服游戏)
	AccountName     		string			`gorm:"column:acc_name;size:64;"`				//账号名，0为游客
	AccountPassword      	string  		`gorm:"column:acc_password;size:32;"`			//密码
	IsEnable				bool			`gorm:"column:is_enabled;"`  					//是否封号
	Platform				int				`gorm:"column:platform;"`						//注册来源(第三方渠道账号，或者代理推广下载)
	Phone					string 			`gorm:"column:phone;size:18;"`  				//手机号	
	CreateTime   			time.Time 		`gorm:"column:create_time;"`					//注册时间
	CreateIP           		string      	`gorm:"column:create_ip;size:15;"`             //注册ip  
}

func (Account) TableName() string {
	return "accounts"
}
