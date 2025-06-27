package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	Name string
	TcpListenOn string
	Amqp string

	UserRpc zrpc.RpcClientConf
}
