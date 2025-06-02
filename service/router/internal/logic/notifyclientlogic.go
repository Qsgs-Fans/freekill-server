package logic

import (
	"context"
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/service/router/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/router/router"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyClientLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyClientLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyClientLogic {
	return &NotifyClientLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NotifyClientLogic) NotifyClient(in *router.Packet) (*router.RouterEmpty, error) {
	server := l.svcCtx.TcpServer
	conn := server.GetConn(in.ConnectionId)
	if conn == nil {
		return &empty,
		fmt.Errorf("Cannot find connection by connId %v", in.ConnectionId)
	}

	err := conn.Notify(in)
	if err != nil {
		return &empty,
		fmt.Errorf("conn.Notify error: %v", err)
	}

	return &empty, nil
}
