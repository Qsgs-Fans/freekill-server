package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartGameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartGameLogic {
	return &StartGameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// rpc KickPlayer
func (l *StartGameLogic) StartGame(in *room.UidRequest) (*room.RoomEmptyReply, error) {
	// todo: add your logic here and delete this line

	return &room.RoomEmptyReply{}, nil
}
