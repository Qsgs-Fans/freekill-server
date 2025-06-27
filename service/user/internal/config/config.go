package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Amqp string
	RoomRpc zrpc.RpcClientConf

	MySqlSource string
	CacheRedis cache.CacheConf
	PlayerRedis redis.RedisConf

}
