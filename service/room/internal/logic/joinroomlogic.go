package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinRoomLogic {
	return &JoinRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 房间应当有自动过期机制
func (l *JoinRoomLogic) JoinRoom(in *room.JoinRoomRequest) (*room.RoomEmptyReply, error) {
	l.svcCtx.AddPlayerToRoom(in.UserId, in.RoomId)

	return &room.RoomEmptyReply{}, nil
}
