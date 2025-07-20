package logic

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/svc"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRobotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRobotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRobotLogic {
	return &AddRobotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TODO list
func (l *AddRobotLogic) AddRobot(in *room.UidRequest) (*room.RoomEmptyReply, error) {
	// todo: add your logic here and delete this line

	return &room.RoomEmptyReply{}, nil
}
