package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BannedIpsModel = (*customBannedIpsModel)(nil)

type (
	// BannedIpsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBannedIpsModel.
	BannedIpsModel interface {
		bannedIpsModel
	}

	customBannedIpsModel struct {
		*defaultBannedIpsModel
	}
)

// NewBannedIpsModel returns a model for the database table.
func NewBannedIpsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BannedIpsModel {
	return &customBannedIpsModel{
		defaultBannedIpsModel: newBannedIpsModel(conn, c, opts...),
	}
}
