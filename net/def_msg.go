package net

import (
	"errors"
	"unsafe"
)

const (
	MsgHeadSize = 8
	MsgNormalVersion = 0x80	   //普通协议头
	MsgRealIpVersion = 0x90	   //协议头插入realip
	MaxMsgDataSize = 40960 
	MsgTimeout = 300
)

var ErrMsgLenTooShort = errors.New("message len too short")
var ErrMsgLenTooLong = errors.New("message len too long Blocked")

const (
	VerNormal    = 0x80	 //原始协议版本
	VerProxy     = 0x81  //添加client ip协议版本
	FlagNorlmal  = 0x10	 //原始协议流
	FlagEncrypt  = 0x11 //数据是经过加密的
	FlagCompress = 0x12 //数据是经过压缩的
)

type MessageHead struct {
	Ver 	uint8 //协议版本
	Fla 	uint8 //标记
	Cmd   	uint8
	Act	  	uint8
	Len   	uint32 //数据长度
}


type Message struct {
	Head       *MessageHead //消息头，可能为nil
	Data       []byte       //消息数据

	UserData   interface{}  //协议解包数据
}

func (r *Message) CmdAct() int {
	if r.Head != nil {
		return CmdAct(r.Head.Cmd, r.Head.Act)
	}
	return 0
}

func (r *Message) Cmd() uint8 {
	if r.Head != nil {
		return r.Head.Cmd
	}
	return 0
}

func (r *Message) Act() uint8 {
	if r.Head != nil {
		return r.Head.Act
	}
	return 0
}

func (r *MessageHead) Bytes() []byte {
	data := make([]byte, MsgHeadSize)
	phead := (*MessageHead)(unsafe.Pointer(&data[0]))
	r.Ver = phead.Ver
	r.Fla = phead.Fla
	r.Len = phead.Len
	r.Cmd = phead.Cmd
	r.Act = phead.Act
	return data
}

func (r *MessageHead) FastBytes(data []byte) []byte {
	phead := (*MessageHead)(unsafe.Pointer(&data[0]))
	phead.Ver = r.Ver
	phead.Fla = r.Fla
	phead.Len = r.Len
	phead.Cmd = r.Cmd
	phead.Act = r.Act
	return data
}

func (r *MessageHead) FromBytes(data []byte) error {
	if len(data) < MsgHeadSize {
		return ErrMsgLenTooShort
	}
	phead := (*MessageHead)(unsafe.Pointer(&data[0]))
	r.Ver = phead.Ver
	r.Fla = phead.Fla
	r.Len = phead.Len
	r.Cmd = phead.Cmd
	r.Act = phead.Act
	if r.Len > MaxMsgDataSize {
		return ErrMsgLenTooLong
	}
	return nil
}

func (r *MessageHead) CmdAct() int {
	return CmdAct(r.Cmd, r.Act)
}


func NewMessageHead(data []byte) *MessageHead {
	head := &MessageHead{}
	if err := head.FromBytes(data); err != nil {
		return nil
	}
	return head
}

func NewMsg(cmd, act, flag uint8, data []byte) *Message {
	return &Message{
		Head: &MessageHead{
			Ver:   MsgNormalVersion,
			Fla:   flag,
			Len:   uint32(len(data)),
			Cmd:   cmd,
			Act:   act,
		},
		Data: data,
	}
}

func CmdAct(cmd, act uint8) int {
	return int(cmd)<<8 + int(act)
}
