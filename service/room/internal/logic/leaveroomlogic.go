package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveRoomLogic {
	return &LeaveRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveRoomLogic) LeaveRoom(in *room.LeaveRoomRequest) (*room.RoomEmptyReply, error) {
	// todo: add your logic here and delete this line

	return &room.RoomEmptyReply{}, nil
}
