package logic

import (
	"context"
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/service/router/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/router/router"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestClientLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRequestClientLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RequestClientLogic {
	return &RequestClientLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RequestClientLogic) RequestClient(in *router.RequestPacket) (*router.RouterEmpty, error) {
	server := l.svcCtx.TcpServer
	conn := server.GetConn(in.ConnectionId)
	if conn == nil {
		return &empty,
		fmt.Errorf("Cannot find connection by connId %v", in.ConnectionId)
	}

	err := conn.Request(in)
	if err != nil {
		return &empty,
		fmt.Errorf("conn.Request error: %v", err)
	}

	return &empty, nil
}
