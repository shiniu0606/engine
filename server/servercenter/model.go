package servercenter

import (
	"errors"
	"sync"

	base "github.com/shiniu0606/engine/core/base"
	net "github.com/shiniu0606/engine/core/net"
	common "github.com/shiniu0606/engine/server/common"
)

type IServerSession interface {
	UID()							int64
	ServerId() 						int					
	SendMsg(msg *net.Message) 
}

type ServerSession struct {
	conn 	net.ISession
	server  common.Server
}

func (s *ServerSession) UID() int64 {
	return s.server.UID
}

func (s *ServerSession) ServerId() int {
	return s.server.ServerId
}

func (s *ServerSession) SendMsg(msg *net.Message)  {
	s.conn.Send(msg)
}

type IServerManager interface {
	Add(s IServerSession)                   
	Remove(s IServerSession)                
	Get(id uint32) (IServerSession, error)
	Len() int 
	GetHandler() net.IMsgHandler
}

type CenterServer struct {
	server 				IServerSession
	connections 		map[int]IServerSession
	connLock    		sync.RWMutex
}

func (s *CenterServer) Add(conn IServerSession) {
	s.connLock.Lock()
	defer s.connLock.Unlock()

	s.connections[conn.ServerId()] = conn

	base.LogInfo("connection add  successfully: conn num = %d", s.Len())
}

func (s *CenterServer) Remove(conn IServerSession) {
	s.connLock.Lock()
	defer s.connLock.Unlock()

	delete(s.connections, conn.ServerId())

	base.LogInfo("connection Remove ConnID = %d, successfully: conn num = %d", conn.ServerId(), s.Len())
}

func (s *CenterServer) Get(connID int) (IServerSession, error) {
	s.connLock.RLock()
	defer s.connLock.RUnlock()

	if conn, ok := s.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (s *CenterServer) Len() int {
	return len(s.connections)
}