package svc

import (
	"github.com/Qsgs-Fans/freekill-server/service/router/internal/config"
	"github.com/Qsgs-Fans/freekill-server/service/user/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	TcpServer *TcpServer
	UserRpc   userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx := &ServiceContext{
		Config:    c,
	}
	tcpServer := NewTcpServer(ctx)
	ctx.TcpServer = tcpServer

	ctx.UserRpc = userclient.NewUser(zrpc.MustNewClient(c.UserRpc))

	return ctx
}
