package svc

import (
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type TcpServer struct {
	svcCtx *ServiceContext
	lister net.Listener
	connections sync.Map // k: connId, v: *TcpConn
	// TODO: rate limit
}

func NewTcpServer(svcCtx *ServiceContext) *TcpServer {
	return &TcpServer {
		svcCtx: svcCtx,
	}
}

func (self *TcpServer) Start() {
	var err error
	self.lister, err = net.Listen("tcp", self.svcCtx.Config.TcpListenOn)
	if err != nil {
		logx.Errorf("Tcp listen error: %v", err)
		return
	}
	logx.Infof("Tcp listening at %v", self.svcCtx.Config.TcpListenOn)

	for {
		conn, err := self.lister.Accept()
		if err != nil {
			logx.Errorf("Tcp Accept error: %v", err)
			continue
		}

		// TODO: SYN rate limiter

		connId := uuid.New().String()
		tcpConn := NewTcpConn(conn, connId, self)
		self.connections.Store(connId, tcpConn)

		go tcpConn.listen()
	}
}

func (self *TcpServer) GetConn(connId string) (*TcpConn) {
	val, ok := self.connections.Load(connId)
	if !ok {
		return nil
	}
	return val.(*TcpConn)
}

func (self *TcpServer) notifyLoginRequest(connId string) {
	// TODO: 等userRpc施工...
}
