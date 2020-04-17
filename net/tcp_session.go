package net

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"strings"

	base "github.com/shiniu0606/engine/base"
)

type tcpSession struct {
	Session
	conn       net.Conn     //连接
	listener   net.Listener //监听
	wait       sync.WaitGroup

	address    string
}

func (r *tcpSession) GetNetType() NetType {
	return NetTypeTcp
}

func (r *tcpSession) LocalAddr() string {
	if r.conn != nil {
		return r.conn.LocalAddr().String()
	} 
	return ""
}

func (r *tcpSession) RemoteAddr() string {
	if r.realRemoteAddr != "" {
		return r.realRemoteAddr
	}
	if r.conn != nil {
		return r.conn.RemoteAddr().String()
	}
	return ""
}

func (r *tcpSession) read() {
	defer func() {
		if err := recover(); err != nil {
			base.LogError("session read panic id:%v err:%v", r.id, err.(error))
			base.LogStack()
		}
		r.Stop()
	}()

	r.readMsg()
}

func (r *tcpSession) readMsg(){
	headData := make([]byte, MsgHeadSize)
	var data []byte
	var head *MessageHead

	for {
		_, err := io.ReadFull(r.conn, headData)
		if err != nil {
			if err != io.EOF {
				base.LogDebug("session:%v recv data err:%v", r.id, err)
			}
			break
		}
		if head = NewMessageHead(headData); head == nil {
			base.LogError("session:%v read msg head failed", r.id)
			break
		}

		if head.Ver != VerNormal && head.Ver != VerProxy {
			base.LogError("session:%v  msg head ver error:%v", r.id,head.Ver)
			break
		}
		//base.LogInfo("=========>read msg head:%v",head)
		if head.Len == 0 {
			if !r.processMsg(r,&Message{Head: head}) {
				base.LogError("session:%v process msg cmd:%v act:%v", r.id, head.Cmd, head.Act)
				break
			}
			head = nil
		} else {
			data = make([]byte, head.Len)

			_, err := io.ReadFull(r.conn, data)
			if err != nil {
				base.LogError("session:%v recv data err:%v", r.id, err)
				break
			}
			//base.LogInfo("=========>read msg data:%v",data)
			if !r.processMsg(r,&Message{Head: head, Data: data}) {
				base.LogError("session:%v process msg cmd:%v act:%v", r.id, head.Cmd, head.Act)
				break
			}

			head = nil
			data = nil
		}

		r.lastTick = base.GetTimestamp()
	}
}

func (r *tcpSession) write() {
	defer func() {
		if err := recover(); err != nil {
			base.LogError("session write panic id:%v err:%v", r.id, err.(error))
			base.LogStack()
		}
		if r.conn != nil {
			r.conn.Close()
		}
		r.Stop()
	}()

	r.writeMsg()
}

func (r *tcpSession) writeMsg() {
	var m *Message
	head := make([]byte, MsgHeadSize)
	writeCount := 0
	tick := time.NewTimer(time.Second * time.Duration(r.timeout))
	defer tick.Stop()
	for {
		select {
		case m = <-r.sendChan:
			if m != nil {
				m.Head.FastBytes(head)

				if writeCount < MsgHeadSize {
					n, err := r.conn.Write(head[writeCount:])
					if err != nil {
						base.LogError("session write id:%v err:%v", r.id, err)
						break
					}
					writeCount += n
				}

				if writeCount >= MsgHeadSize && m.Data != nil {
					n, err := r.conn.Write(m.Data[writeCount-MsgHeadSize : int(m.Head.Len)])
					if err != nil {
						base.LogError("session write id:%v err:%v", r.id, err)
						break
					}
					writeCount += n
				}

				if writeCount == int(m.Head.Len)+MsgHeadSize {
					writeCount = 0
				}
				r.lastTick = base.GetTimestamp()
			}
		case <-tick.C:
			base.LogInfo("session write timeout id:%v ", r.id)
			if r.isTimeout(tick) {
				r.Stop()
			}
		case <- r.closeChan:
			r.handler.OnCloseHandle(r)
			base.LogInfo("session write close id:%v ", r.id)
			return
		}
	}
}

func (r *tcpSession) connect() {
	base.LogDebug("connect to addr:%s session:%d", r.address, r.id)
	c, err := net.DialTimeout("tcp", r.address, time.Second*time.Duration(TcpDialTimeout))
	if err != nil {
		base.LogError("connect to addr:%s failed session:%d err:%v", r.address, r.id, err)
		r.handler.OnConnectCompleteHandle(r, false)
		r.Stop()
	} else {
		r.conn = c
		base.LogDebug("connect to addr:%s ok session:%d", r.address, r.id)
		if r.handler.OnConnectCompleteHandle(r, true) {
			base.Go(func() {
				base.LogInfo("process read for session:%d", r.id)
				r.read()
				base.LogInfo("process read end for session:%d", r.id)
			})
			base.Go(func() {
				base.LogInfo("process write for session:%d", r.id)
				r.write()
				base.LogInfo("process write end for session:%d", r.id)
			})
		} else {
			r.Stop()
		}
	}
}

func (r *tcpSession) listen() {
	defer func(){
		if err := recover(); err != nil {
			base.LogError("session write panic id:%v err:%v", r.id, err.(error))
			base.LogStack()
		}
		if r.listener != nil {
			r.listener.Close()
		}
		r.Stop()
	}()
	for {
		c, err := r.listener.Accept()
		if err != nil {
			base.LogError("accept failed session:%v err:%v", r.id, err)
			break
		} else {
			base.Go(func(){
				session := newTcpAccept(c, r.handler)
				if r.handler.OnStartHandle(session) {
					base.Go(func(){
						base.LogInfo("process read for session:%d", session.id)
						session.read()
						base.LogInfo("process read end for session:%d", session.id)
					})

					base.Go(func(){
						base.LogInfo("process write for session:%d", session.id)
						session.write()
						base.LogInfo("process write end for session:%d", session.id)
					})
				} else {
					if r.conn != nil {
						r.conn.Close()
					}
					session.Stop()
				}
			})
		}
	}
}

func newTcpAccept(conn net.Conn, handler IMsgHandler) *tcpSession {
	tcpsession := tcpSession{
		Session: Session{
			id:            atomic.AddUint64(&base.GetGlobal().MsgSessionId, 1),
			sendChan:      make(chan *Message, 64),
			closeChan:	   make(chan int),
			msgTyp:        NetTypeTcp,
			handler:       handler,
			timeout:       MsgTimeout,
			connTyp:       ConnTypeAccept,
			lastTick:      base.GetTimestamp(),
		},
		conn: conn,
	}

	base.LogInfo("new accept session id:%d from addr:%s", tcpsession.Id(), conn.RemoteAddr().String())
	return &tcpsession
}

func newTcpListen(listener net.Listener, handler IMsgHandler, addr string) *tcpSession {
	tcpsession := tcpSession{
		Session: Session{
			id:            atomic.AddUint64(&base.GetGlobal().MsgSessionId, 1),
			msgTyp:        NetTypeTcp,
			handler:       handler,
			connTyp:       ConnTypeListen,
		},
		listener: listener,
	}

	base.LogInfo("new tcp listen id:%d addr:%s", tcpsession.Id(), addr)
	return &tcpsession
}

func newTcpConn(addr string, conn net.Conn, handler IMsgHandler) *tcpSession {
	tcpsession := tcpSession{
		Session: Session{
			id:            atomic.AddUint64(&base.GetGlobal().MsgSessionId, 1),
			sendChan:      make(chan *Message, base.GetGlobal().MaxMsgChanLen),
			closeChan:	   make(chan int),
			msgTyp:        NetTypeTcp,
			handler:       handler,
			timeout:       MsgTimeout,
			connTyp:       ConnTypeConn,
			lastTick:      base.GetTimestamp(),
		},
		conn:    conn,
		address: addr,
	}

	base.LogDebug("new session id:%d connect to addr:%s", tcpsession.Id(), addr)
	return &tcpsession
}

func StartTcpServer(addr string, handler IMsgHandler) error {
	addrs := strings.Split(addr, "://")
	if addrs[0] == "tcp" || addrs[0] == "all" {
		listen, err := net.Listen("tcp", addrs[1])
		if err == nil {
			session := newTcpListen(listen,handler, addr)
			base.Go(func() {
				base.LogDebug("process listen for tcp session:%d", session.id)
				session.listen()
				base.LogDebug("process listen end for tcp session:%d", session.id)
			})
		} else {
			base.LogError("listen on %s failed, errstr:%s", addr, err)
			return err
		}
	}

	return nil
}

func StartTcpConnect(addr string, handler IMsgHandler) ISession {
	session := newTcpConn(addr, nil,handler)
	
	if handler.OnStartHandle(session) {
		session.connect()
		return session
	} else {
		session.Stop()
	}
	return nil
}


