package logic

import (
	"context"

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

func (l *RequestClientLogic) RequestClient(in *router.Packet) (*router.PacketSendResponse, error) {
	// todo: add your logic here and delete this line

	return &router.PacketSendResponse{}, nil
}
