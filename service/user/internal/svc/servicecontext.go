package svc

import (
	"fmt"

	"github.com/Qsgs-Fans/freekill-server/common/cryptoutil"
	"github.com/Qsgs-Fans/freekill-server/service/room/room"
	"github.com/Qsgs-Fans/freekill-server/service/room/roomclient"
	"github.com/Qsgs-Fans/freekill-server/service/router/sender"
	"github.com/Qsgs-Fans/freekill-server/service/user/internal/config"
	"github.com/Qsgs-Fans/freekill-server/service/user/model"
	"github.com/Qsgs-Fans/freekill-server/service/user/user"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	Sender *sender.Sender
	RoomRpc room.RoomClient

	RsaKeyPair *cryptoutil.RSAKeyPair

	UsersModel model.UsersModel
	UserLoginHistoryModel model.UserLoginHistoryModel
	BannedIpsModel model.BannedIpsModel
	BannedDevicesModel model.BannedDevicesModel
	UsernameBlacklistModel model.UsernameBlacklistModel
  UsernameWhitelistModel model.UsernameWhitelistModel

	// redis设计：HASH k=user:${uid} v=HASH(Json Object) 用于取得user信息.
	//						SET online_users 用于取得在线用户（掉线用户不在此SET 但还是在HASH）
	//      这样就实现了C++版中QHash<int, ServerPlayer *>的功能了
	PlayerRedis *redis.Redis
	// TODO players 大约map[int]Player
}

func (s *ServiceContext) AddNewLoginUser(userInfo *model.Users, in *user.LoginRequest) error {
	rds := s.PlayerRedis
	// TODO 原子化与Error
	rds.Sadd("freekill:online_users", userInfo.Id)
	userKey := fmt.Sprintf("freekill:user:%v", userInfo.Id)
	rds.Hmset(userKey, map[string]string{
		// "id": strconv.FormatInt(userInfo.Id, 10),
		"username": userInfo.Username,
		"avatar": userInfo.Avatar,
		"connId": in.ConnId,
	})
	// TODO 处理掉线；没做掉线所有先自然expire
	// rds.Expire(userKey, 300)

	return nil
}

func (s *ServiceContext) HandleUserLogout(userId int64) error {
	rds := s.PlayerRedis
	// TODO 原子化与Error
	rds.Srem("freekill:online_users", userId)
	userKey := fmt.Sprintf("freekill:user:%v", userId)
	rds.Del(userKey)

	// TODO 人机房间or退出大厅

	return nil
}

func NewServiceContext(c config.Config) *ServiceContext {
	kpair, err := cryptoutil.LoadOrGenerateRSAKeys("rsakey/rsa", "rsakey/rsa_pub")
	if err != nil {
		panic(err)
	}

	dbconn := sqlx.NewMysql(c.MySqlSource)

	return &ServiceContext{
		Config: c,
		Sender: sender.MustNewSender(c.Amqp),
		RoomRpc: roomclient.NewRoom(zrpc.MustNewClient(c.RoomRpc)),

		RsaKeyPair: kpair,

		UsersModel: model.NewUsersModel(dbconn, c.CacheRedis),
		UserLoginHistoryModel: model.NewUserLoginHistoryModel(dbconn, c.CacheRedis),
		BannedIpsModel: model.NewBannedIpsModel(dbconn, c.CacheRedis),
		BannedDevicesModel: model.NewBannedDevicesModel(dbconn, c.CacheRedis),
		UsernameBlacklistModel: model.NewUsernameBlacklistModel(dbconn, c.CacheRedis),
		UsernameWhitelistModel: model.NewUsernameWhitelistModel(dbconn, c.CacheRedis),

		PlayerRedis: redis.MustNewRedis(c.PlayerRedis),
	}
}
