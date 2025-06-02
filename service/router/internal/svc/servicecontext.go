package svc

import (
	"github.com/Qsgs-Fans/freekill-server/service/router/internal/config"
	"github.com/Qsgs-Fans/freekill-server/service/router/internal/other"
)

type ServiceContext struct {
	Config config.Config

	TcpServer *other.TcpServer
}

func NewServiceContext(c config.Config) *ServiceContext {
	tcpServer := other.NewTcpServer(c.TcpListenOn)

	return &ServiceContext{
		Config: c,
		TcpServer: tcpServer,
	}
}
