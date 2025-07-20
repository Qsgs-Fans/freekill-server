package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnterRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEnterRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnterRoomLogic {
	return &EnterRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EnterRoomLogic) EnterRoom(in *room.UidAndRidRequest) (*room.RoomEmptyReply, error) {
	l.svcCtx.AddPlayerToRoom(in.UserId, in.RoomId)

	return &room.RoomEmptyReply{}, nil
}
