package net

const (
	TcpDialTimeout = 5
)

type NetType int

const (
	NetTypeTcp NetType = iota //TCP类型
	NetTypeUdp                //UDP类型
	NetTypeWs                 //websocket
)

type ConnType int

const (
	ConnTypeListen ConnType = iota //监听
	ConnTypeConn                   //连接产生的
	ConnTypeAccept                 //Accept产生的
)


type ISession interface {
	Id() uint64
	GetConnType() ConnType
	GetNetType() NetType
	GetParser() IParser

	LocalAddr() string
	RemoteAddr() string
	RealRemoteAddr() string
	SetRealRemoteAddr(addr string)

	Stop() error
	IsStop() bool
	SetTimeout(t int)
	GetTimeout() int

	SetUser(user interface{})
	GetUser() interface{}

	Send(m *Message) bool
	SendString(str string)  bool
}

