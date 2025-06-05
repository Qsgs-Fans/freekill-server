package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BannedDevicesModel = (*customBannedDevicesModel)(nil)

type (
	// BannedDevicesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBannedDevicesModel.
	BannedDevicesModel interface {
		bannedDevicesModel
	}

	customBannedDevicesModel struct {
		*defaultBannedDevicesModel
	}
)

// NewBannedDevicesModel returns a model for the database table.
func NewBannedDevicesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BannedDevicesModel {
	return &customBannedDevicesModel{
		defaultBannedDevicesModel: newBannedDevicesModel(conn, c, opts...),
	}
}
