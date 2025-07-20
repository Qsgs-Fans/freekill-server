package svc

import (
	"context"
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/service/room/internal/config"
	"github.com/Qsgs-Fans/freekill-server/service/router/sender"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config

	Sender *sender.Sender

	// redis设计：
	// 			SET room_indexes (int) 房间id列表，获取所有房间用
	//			HASH room:${uid}:info 房间的基本信息（就用户大厅需要的那几个config）
	//			string room:${uid}:config JSON字符串 表示房间的复杂配置 用户加入时发过去以及发给lua
	//			SET room:${uid}:players 房间内玩家列表
	//			SET lobby:players 大厅内玩家列表
	//      这样就实现了C++版中QHash<int, Room *>的功能了
	//
	//      当然对于player->getRoom()这种东西 可以弄个player_get_room:${uid}来做
	//      Room 0表示大厅
	RoomRedis *redis.Redis
}

func (s *ServiceContext) CreateNewRoom() {
}

func (s *ServiceContext) AddPlayerToRoom(userId int64, roomId int64) {
	rds := s.RoomRedis
	// TODO 原子化与Error
	userKey := fmt.Sprintf("freekill:user:%v", userId)
	roomKey := fmt.Sprintf("freekill:room:%v:players", roomId)
	rds.Sadd(roomKey, userId)

	ctx := context.Background()

	if roomId == 0 {
		// FIXME Room服务假设RoomRedis和PlayerRedis是同一个并不妥当。
		connId, _ := rds.Hget(userKey, "connId")
		s.Sender.Notify(ctx, "EnterLobby", nil, connId)
		packet := &router.Packet{
			Command: "EnterLobby",
			Data: "",
			ConnectionId: connId,
		}
		s.Sender.Notify(ctx, packet)

>>>>>>> Stashed changes
		// TODO server->updateOnlineInfo
	} else {
	}
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Sender: sender.MustNewSender(c.Amqp),
		RoomRedis: redis.MustNewRedis(c.RoomRedis),
	}
}
