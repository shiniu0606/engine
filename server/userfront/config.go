package userfront

import (
	"fmt"
	"path/filepath"

	base "github.com/shiniu0606/engine/core/base"
	jbp "github.com/shiniu0606/engine/core/db"
)

type ServerConfig struct {
	ServerId 	int32 	`json:"ServerId"`
	ServerName 	string 	`json:"ServerName"`
	TcpPort 	int 	`json:"TcpPort"`
	RemoteAddr  string	`json:"RemoteAddr"`
	CenterAddr  string  `json:"CenterAddr"`
	LogPath 	string 	`json:"LogPath"`
	DataBase   DBConfig  `json:"DataBase"`
}

type DBConfig struct {
	User			string 		`json:"User"`
	Password        string		`json:"Password"`
	Name            string		`json:"Name"`
	Port            int			`json:"Port"`
	IP              string		`json:"IP"`
	MaxPoll         int			`json:"MaxPoll"`
	IdlePoll        int			`json:"IdlePoll"`
}

var serverconfig ServerConfig

func InitConfig(path string) {
	//初始化配置文件
	configFilePath, _ := filepath.Abs(path)
	err := base.InitJsonConfigFile(configFilePath, &serverconfig)
	if err != nil {
		panic(err)
	}
}

func InitDB() {
	username := serverconfig.DataBase.User  //账号
	password := serverconfig.DataBase.Password //密码
	host := serverconfig.DataBase.IP //数据库地址，可以是Ip或者域名
	port := serverconfig.DataBase.Port //数据库端口
	Dbname := serverconfig.DataBase.Name //数据库名
	timeout := "10s" //连接超时，10秒

	//拼接dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

	jbp.InitDBForPoll(jbp.DEFAULT_DB_DRIVER,dsn,serverconfig.DataBase.MaxPoll,serverconfig.DataBase.IdlePoll)
}

func InitLog() {
	filepath := serverconfig.LogPath
	base.SetFileLog(filepath,1024*1024)
}