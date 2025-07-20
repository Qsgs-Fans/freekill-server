package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuitRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuitRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitRoomLogic {
	return &QuitRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuitRoomLogic) QuitRoom(in *room.UidRequest) (*room.RoomEmptyReply, error) {
	// todo: add your logic here and delete this line

	return &room.RoomEmptyReply{}, nil
}
