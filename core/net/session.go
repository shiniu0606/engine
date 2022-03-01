package net

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	base "github.com/shiniu0606/engine/core/base"
)

var ErrSessionClosed = errors.New("session closes")
var ErrSessionBlocked = errors.New("session blocked")

type Session struct {
	id uint64

	msgTyp         NetType  //消息类型
	connTyp        ConnType //通道类型
	realRemoteAddr string

	sendChan  chan *Message
	sendMutex sync.RWMutex

	handler IMsgHandler
	parser  IParser
	timeout int //传输超时

	closeFlag int32
	closeChan chan int
	lastTick  int64

	user interface{} //用户
}

func (r *Session) Id() uint64 {
	return r.id
}

func (r *Session) GetNetType() NetType {
	return r.msgTyp
}

func (r *Session) GetConnType() ConnType {
	return r.connTyp
}

func (r *Session) GetParser() IParser {
	return r.parser
}

func (r *Session) SetRealRemoteAddr(addr string) {
	r.realRemoteAddr = addr
}

func (r *Session) RealRemoteAddr() string {
	return r.realRemoteAddr
}

func (r *Session) SetUser(user interface{}) {
	r.user = user
}

func (r *Session) GetUser() interface{} {
	return r.user
}

func (r *Session) SetTimeout(t int) {
	if t >= 0 {
		r.timeout = t
	}
}

func (r *Session) GetTimeout() int {
	return r.timeout
}

func (r *Session) isTimeout(tick *time.Timer) bool {
	left := int(base.GetTimestamp() - r.lastTick)
	if left < r.timeout || r.timeout == 0 {
		if r.timeout == 0 {
			tick.Reset(time.Second * time.Duration(MsgTimeout))
		} else {
			tick.Reset(time.Second * time.Duration(r.timeout-left))
		}
		return false
	}
	base.LogInfo("msgque close because timeout id:%v wait:%v timeout:%v", r.id, left, r.timeout)
	return true
}

func (r *Session) Send(m *Message) bool {
	if m == nil {
		return false
	}
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	r.sendMutex.RLock()
	if r.IsStop() {
		r.sendMutex.RUnlock()
		return false
	}

	select {
	case r.sendChan <- m:
		r.sendMutex.RUnlock()
		return true
	default:
		r.sendMutex.RUnlock()
		base.LogWarn("session write channel full msgque:%v", r.id)
		r.Stop()
		return false
	}
}

func (r *Session) SendString(str string) bool {
	return r.Send(&Message{Data: []byte(str)})
}

func (r *Session) IsStop() bool {
	return r.closeFlag == 1
}

func (r *Session) Stop() error {
	if atomic.CompareAndSwapInt32(&r.closeFlag, 0, 1) {
		close(r.closeChan)

		if r.sendChan != nil {
			r.sendMutex.Lock()
			close(r.sendChan)
			r.sendMutex.Unlock()
		}

		return nil
	}
	return ErrSessionClosed
}

func (r *Session) processMsg(s ISession, msg *Message) bool {
	f := r.handler.GetHandlerFunc(msg)
	if f == nil {
		f = r.handler.OnProcessMsgHandle
	}

	if r.parser != nil {
		//base.LogInfo("processMsg start:%v",msg.Data)
		r.parser.UnPack(msg)
		//base.LogInfo("processMsg start:%v",msg.UserData)
	}

	return f(s, msg)
}
