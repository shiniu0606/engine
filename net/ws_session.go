package net

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"strings"

	"github.com/gorilla/websocket"
	base "github.com/shiniu0606/engine/base"
)

type wsSession struct {
	Session
	conn     *websocket.Conn
	upgrader *websocket.Upgrader
	addr     string
	url      string
	wait       sync.WaitGroup
	listener *http.Server

	enablewss bool
	sslcrtpath string
	sslkeypath string
}

func (r *wsSession) GetNetType() NetType {
	return NetTypeWs
}

func (r *wsSession) LocalAddr() string {
	if r.conn != nil {
		return r.conn.LocalAddr().String()
	}
	return ""
}

func (r *wsSession) RemoteAddr() string {
	if r.realRemoteAddr != "" {
		return r.realRemoteAddr
	}
	if r.conn != nil {
		return r.conn.RemoteAddr().String()
	}
	return ""
}

func (r *wsSession) readMsg() {
	for {
		_, data, err := r.conn.ReadMessage()
		if err != nil {
			base.LogError("wsSession:%v recv data err:%v", r.id, err)
			break
		}
		if !r.processMsg(r, &Message{Data: data}) {
			break
		}
		r.lastTick = base.GetTimestamp()
	}
}

func (r *wsSession) read() {
	defer func() {
		if err := recover(); err != nil {
			base.LogError("wsSession read panic id:%v err:%v", r.id, err.(error))
			base.LogStack()
		}
		r.Stop()
	}()

	r.readMsg()
}

func (r *wsSession) writeMsg() {
	var m *Message
	tick := time.NewTimer(time.Second * time.Duration(r.timeout))
	defer tick.Stop()
	for {
		select {
		case m = <-r.sendChan:
			if m != nil {
				err := r.conn.WriteMessage(websocket.BinaryMessage, m.Data)
				if err != nil {
					base.LogError("wsSession write id:%v err:%v", r.id, err)
					break
				}
				r.lastTick = base.GetTimestamp()
			}
		case <-tick.C:
			base.LogInfo("wsSession write timeout id:%v ", r.id)
			if r.isTimeout(tick) {
				r.Stop()
			}
		case <- r.closeChan:
			r.handler.OnCloseHandle(r)
			base.LogInfo("wsSession write close id:%v ", r.id)
			return
		}
	}
}

func (r *wsSession) write() {
	defer func() {
		if err := recover(); err != nil {
			base.LogError("wsSession write panic id:%v err:%v", r.id, err.(error))
			base.LogStack()
		}
		if r.conn != nil {
			r.conn.Close()
		}
		r.Stop()
	}()

	r.writeMsg()
}

func (r *wsSession) listen() {
	defer func(){
		if err := recover(); err != nil {
			base.LogError("wsSession write panic id:%v err:%v", r.id, err.(error))
			base.LogStack()
		}
		if r.listener != nil {
			r.listener.Close()
		}
		r.Stop()
	}()

	r.upgrader = &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc(r.url, func(hw http.ResponseWriter, hr *http.Request) {
		c, err := r.upgrader.Upgrade(hw, hr, nil)
		if err != nil {
			base.LogError("wsSession accept failed :%v err:%v", r.id, err)
		} else {
			base.Go(func() {
				session := newWsAccept(c, r.handler)
				if r.handler.OnStartHandle(session) {
					base.Go(func() {
						base.LogInfo("process read for session:%d", session.id)
						session.read()
						base.LogInfo("process read end for session:%d", session.id)
					})
					base.Go(func() {
						base.LogInfo("process write for session:%d", session.id)
						session.write()
						base.LogInfo("process write end for session:%d", session.id)
					})
				} else {
					session.Stop()
				}
			})
		}
	})

	if r.enablewss {
		if r.sslcrtpath != "" && r.sslkeypath != "" {
			r.listener.ListenAndServeTLS(r.sslcrtpath, r.sslkeypath)
		} else {
			base.LogError("start wss failed ssl path not set please set now auto change to ws")
			r.listener.ListenAndServe()
		}
	} else {
		r.listener.ListenAndServe()
	}
}

func (r *wsSession) connect() {
	base.LogInfo("connect to addr:%s wsSession:%d",r.addr, r.id)
	c, _, err := websocket.DefaultDialer.Dial(r.addr, nil)
	if err != nil {
		base.LogInfo("connect to addr:%s failed wsSession:%d err:%v ", r.addr, r.id, err)
		r.handler.OnConnectCompleteHandle(r, false)
		r.Stop()
	} else {
		r.conn = c
		base.LogInfo("connect to addr:%s ok wsSession:%d", r.addr, r.id)
		if r.handler.OnConnectCompleteHandle(r, true) {
			base.Go(func() {
				base.LogInfo("process read for wsSession:%d", r.id)
				r.read()
				base.LogInfo("process read end for wsSession:%d", r.id)
			})
			base.Go(func() {
				base.LogInfo("process write for wsSession:%d", r.id)
				r.write()
				base.LogInfo("process write end for wsSession:%d", r.id)
			})
		} else {
			r.Stop()
		}
	}
}

func newWsAccept(conn *websocket.Conn, handler IMsgHandler) *wsSession {
	wssession := wsSession{
		Session: Session{
			id:            atomic.AddUint64(&base.GetGlobal().MsgSessionId, 1),
			sendChan:      make(chan *Message, 64),
			closeChan:	   make(chan int),
			msgTyp:        NetTypeWs,
			handler:       handler,
			timeout:       MsgTimeout,
			connTyp:       ConnTypeAccept,
			lastTick:      base.GetTimestamp(),
		},
		conn: conn,
	}

	base.LogInfo("new wsSession id:%d from addr:%s", wssession.id, conn.RemoteAddr().String())
	return &wssession
}

func newWsListen(addr, url string,enableWss bool,wssCrtPath, wssKeyPath string, handler IMsgHandler) *wsSession {
	wssession := wsSession{
		Session: Session{
			id:            	atomic.AddUint64(&base.GetGlobal().MsgSessionId, 1),
			sendChan:      	make(chan *Message, 64),
			closeChan:	   	make(chan int),
			msgTyp:        	NetTypeWs,
			handler:       	handler,
			timeout:       	MsgTimeout,
			connTyp:       	ConnTypeListen,
			lastTick:      	base.GetTimestamp(),
		},
		addr:     addr,
		url:      url,
		enablewss:	   	enableWss,
		sslcrtpath:		wssCrtPath,
		sslkeypath:		wssKeyPath,
		listener: &http.Server{Addr: addr},
	}

	base.LogInfo("new wssession listen id:%d addr:%s url:%s", wssession.id, addr, url)
	return &wssession
}

func newWsConn(addr string, conn *websocket.Conn, handler IMsgHandler) *wsSession {
	wssession := wsSession{
		Session: Session{
			id:            atomic.AddUint64(&base.GetGlobal().MsgSessionId, 1),
			sendChan:      make(chan *Message, 64),
			closeChan:	   make(chan int),
			msgTyp:        NetTypeWs,
			handler:       handler,
			timeout:       MsgTimeout,
			connTyp:       ConnTypeConn,
			lastTick:      base.GetTimestamp(),
		},
		conn:    conn,
		addr: addr,
	}

	base.LogInfo("new wssession conn id:%d connect to addr:%s", wssession.id, addr)
	return &wssession
}

func StartWebscoketServer(addr string, handler IMsgHandler,wsscrtpath,wsskeypath string) error {
	addrs := strings.Split(addr, "://")
	enablewss := false
	if addrs[0] == "ws" || addrs[0] == "wss" {
		naddr := strings.SplitN(addrs[1], "/", 2)
		url := "/"
		if len(naddr) > 1 {
			url = "/" + naddr[1]
		}
		if addrs[0] == "wss" {
			enablewss = true
		}

		wssession := newWsListen(naddr[0], url,enablewss, wsscrtpath,wsskeypath, handler)
		base.Go(func(){
			base.LogDebug("process listen for wssession:%d", wssession.id)
			wssession.listen()
			base.LogDebug("process listen end for wssession:%d", wssession.id)
		})
	}

	return nil
}

func StartWebsocketConnect(addr string, handler IMsgHandler) ISession {
	session := newWsConn(addr, nil,handler)
	
	if handler.OnStartHandle(session) {
		session.connect()
		return session
	} else {
		session.Stop()
	}
	return nil
}
