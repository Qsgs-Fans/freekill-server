// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3
// Source: user.proto

package userclient

import (
	"context"

	"github.com/Qsgs-Fans/freekill-server/service/user/user"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ConnIdMsg    = user.ConnIdMsg
	LoginReply   = user.LoginReply
	LoginRequest = user.LoginRequest
	UserEmpty    = user.UserEmpty

	User interface {
		NewConn(ctx context.Context, in *ConnIdMsg, opts ...grpc.CallOption) (*UserEmpty, error)
		Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
		Logout(ctx context.Context, in *ConnIdMsg, opts ...grpc.CallOption) (*UserEmpty, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

func (m *defaultUser) NewConn(ctx context.Context, in *ConnIdMsg, opts ...grpc.CallOption) (*UserEmpty, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.NewConn(ctx, in, opts...)
}

func (m *defaultUser) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultUser) Logout(ctx context.Context, in *ConnIdMsg, opts ...grpc.CallOption) (*UserEmpty, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Logout(ctx, in, opts...)
}
