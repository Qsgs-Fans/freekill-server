// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3
// Source: room.proto

package roomclient

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/room/room"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateRoomReply   = room.CreateRoomReply
	CreateRoomRequest = room.CreateRoomRequest
	RoomEmptyReply    = room.RoomEmptyReply
	UidAndRidRequest  = room.UidAndRidRequest
	UidRequest        = room.UidRequest

	Room interface {
		CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomReply, error)
		EnterRoom(ctx context.Context, in *UidAndRidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error)
		ObserveRoom(ctx context.Context, in *UidAndRidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error)
		QuitRoom(ctx context.Context, in *UidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error)
		// TODO list
		AddRobot(ctx context.Context, in *UidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error)
		// rpc KickPlayer
		StartGame(ctx context.Context, in *UidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error)
	}

	defaultRoom struct {
		cli zrpc.Client
	}
)

func NewRoom(cli zrpc.Client) Room {
	return &defaultRoom{
		cli: cli,
	}
}

func (m *defaultRoom) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomReply, error) {
	client := room.NewRoomClient(m.cli.Conn())
	return client.CreateRoom(ctx, in, opts...)
}

func (m *defaultRoom) EnterRoom(ctx context.Context, in *UidAndRidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error) {
	client := room.NewRoomClient(m.cli.Conn())
	return client.EnterRoom(ctx, in, opts...)
}

func (m *defaultRoom) ObserveRoom(ctx context.Context, in *UidAndRidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error) {
	client := room.NewRoomClient(m.cli.Conn())
	return client.ObserveRoom(ctx, in, opts...)
}

func (m *defaultRoom) QuitRoom(ctx context.Context, in *UidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error) {
	client := room.NewRoomClient(m.cli.Conn())
	return client.QuitRoom(ctx, in, opts...)
}

// TODO list
func (m *defaultRoom) AddRobot(ctx context.Context, in *UidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error) {
	client := room.NewRoomClient(m.cli.Conn())
	return client.AddRobot(ctx, in, opts...)
}

// rpc KickPlayer
func (m *defaultRoom) StartGame(ctx context.Context, in *UidRequest, opts ...grpc.CallOption) (*RoomEmptyReply, error) {
	client := room.NewRoomClient(m.cli.Conn())
	return client.StartGame(ctx, in, opts...)
}
