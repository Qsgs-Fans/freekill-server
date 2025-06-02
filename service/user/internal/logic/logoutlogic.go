package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/user/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *user.ConnIdMsg) (*user.UserEmpty, error) {
	// todo: add your logic here and delete this line

	return &empty, nil
}
