package other

import (
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type TcpServer struct {
	listenOn string
	lister net.Listener
	connections sync.Map // k: connId, v: *TcpConn
	// TODO: rate limit
}

func NewTcpServer(listenOn string) *TcpServer {
	return &TcpServer {
		listenOn: listenOn,
	}
}

func (self *TcpServer) Start() {
	var err error
	self.lister, err = net.Listen("tcp", self.listenOn)
	if err != nil {
		logx.Errorf("Tcp listen error: %v", err)
		return
	}
	logx.Infof("Tcp listening at %v", self.listenOn)

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

func (self *TcpServer) notifyLoginRequest(connId string) {
	// TODO: 等userRpc施工...
}
