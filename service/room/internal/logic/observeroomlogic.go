package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/core/logx"
)

type ObserveRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewObserveRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ObserveRoomLogic {
	return &ObserveRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ObserveRoomLogic) ObserveRoom(in *room.UidAndRidRequest) (*room.RoomEmptyReply, error) {
	// todo: add your logic here and delete this line

	return &room.RoomEmptyReply{}, nil
}
