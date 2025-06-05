package svc

import (
	"github.com/Qsgs-Fans/freekill-server/common/cryptoutil"
	"github.com/Qsgs-Fans/freekill-server/service/router/sender"
	"github.com/Qsgs-Fans/freekill-server/service/user/internal/config"
	"github.com/Qsgs-Fans/freekill-server/service/user/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	Sender *sender.Sender
	RsaKeyPair *cryptoutil.RSAKeyPair

	UsersModel model.UsersModel
	UserLoginHistoryModel model.UserLoginHistoryModel
	BannedIpsModel model.BannedIpsModel
	BannedDevicesModel model.BannedDevicesModel
	UsernameBlacklistModel model.UsernameBlacklistModel
  UsernameWhitelistModel model.UsernameWhitelistModel

	// PlayerRedis *redis.Redis

	// TODO players 大约map[int]Player
}

func NewServiceContext(c config.Config) *ServiceContext {
	kpair, err := cryptoutil.LoadOrGenerateRSAKeys("rsakey/rsa", "rsakey/rsa_pub")
	if err != nil {
		panic(err)
	}

	dbconn := sqlx.NewMysql(c.MySqlSource)

	return &ServiceContext{
		Config: c,
		Sender: sender.NewSender(c.RouterRpc),
		RsaKeyPair: kpair,

		UsersModel: model.NewUsersModel(dbconn, c.CacheRedis),
		UserLoginHistoryModel: model.NewUserLoginHistoryModel(dbconn, c.CacheRedis),
		BannedIpsModel: model.NewBannedIpsModel(dbconn, c.CacheRedis),
		BannedDevicesModel: model.NewBannedDevicesModel(dbconn, c.CacheRedis),
		UsernameBlacklistModel: model.NewUsernameBlacklistModel(dbconn, c.CacheRedis),
		UsernameWhitelistModel: model.NewUsernameWhitelistModel(dbconn, c.CacheRedis),
	}
}
