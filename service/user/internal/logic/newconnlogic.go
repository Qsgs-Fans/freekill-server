package logic

import (
	"context"
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/service/router/router"
	"github.com/Qsgs-Fans/freekill-server/service/user/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/user/model"
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
	// TODO 暂封ip（Redis）、服务器满员判定
	sender := l.svcCtx.Sender
	var errmsg string
	_, err := l.svcCtx.BannedIpsModel.FindOneByIpAddress(l.ctx, string(in.ConnIp))
	if err == nil {
		// TODO 这里可以发回ban reason和expire数据 但是需要等待客户端更新版本
		errmsg = "you have been banned!"
		packet := &router.Packet{
			Command: "ErrorDlg",
			Data: errmsg,
			ConnectionId: in.ConnId,
		}
		err = sender.Notify(l.ctx, packet)
		if err != nil {
			return &empty, fmt.Errorf("Error when sending banned message: %v", err)
		}
		return &empty, fmt.Errorf("Refused banned IP connId=%v", in.ConnId)
	} else if err != model.ErrNotFound {
		return &empty, err
	}

	packet := &router.Packet{
		Command: "NetworkDelayTest",
		Data: l.svcCtx.RsaKeyPair.PublicKeyString,
		ConnectionId: in.ConnId,
	}
	err = sender.Notify(l.ctx, packet)
	if err != nil {
		return &empty, err
	}

	return &empty, nil
}
