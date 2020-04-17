package net

import (

)

type HandlerFunc func(msgque ISession, msg *Message) bool

type IMsgHandler interface {
	OnStartHandle(session ISession) bool                         //新的消息队列
	OnCloseHandle(session ISession)                              //消息队列关闭
	OnProcessMsgHandle(session ISession, msg *Message) bool          //默认的消息处理函数
	OnConnectCompleteHandle(session ISession, ok bool) bool          //连接成功
	GetHandlerFunc(msg *Message) HandlerFunc
}

type DefMsgHandler struct {
	msgMap  map[int]HandlerFunc
}

func (r *DefMsgHandler) OnStartHandle(session ISession) bool                { return true }
func (r *DefMsgHandler) OnCloseHandle(session ISession)                     {}
func (r *DefMsgHandler) OnProcessMsgHandle(session ISession, msg *Message) bool { return true }
func (r *DefMsgHandler) OnConnectCompleteHandle(session ISession, ok bool) bool { return true }
func (r *DefMsgHandler) GetHandlerFunc(msg *Message) HandlerFunc {
	if r.msgMap != nil {
		if f, ok := r.msgMap[msg.CmdAct()]; ok {
			return f
		}
	}

	return nil
}

func (r *DefMsgHandler) Register(cmd, act uint8, fun HandlerFunc) {
	if r.msgMap == nil {
		r.msgMap = map[int]HandlerFunc{}
	}
	r.msgMap[CmdAct(cmd, act)] = fun
}