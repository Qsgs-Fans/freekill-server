package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/router/router"
	"github.com/Qsgs-Fans/freekill-server/service/user/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewConnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewConnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewConnLogic {
	return &NewConnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewConnLogic) NewConn(in *user.ConnIdMsg) (*user.UserEmpty, error) {
	sender := l.svcCtx.Sender
	packet := &router.Packet{
		Command: "NetworkDelayTest",
		Data: l.svcCtx.RsaKeyPair.PublicKeyString,
		ConnectionId: in.ConnId,
	}
	err := sender.Notify(l.ctx, packet)
	if err != nil {
		return &empty, err
	}

	return &empty, nil
}
