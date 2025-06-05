package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsernameBlacklistModel = (*customUsernameBlacklistModel)(nil)

type (
	// UsernameBlacklistModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsernameBlacklistModel.
	UsernameBlacklistModel interface {
		usernameBlacklistModel
	}

	customUsernameBlacklistModel struct {
		*defaultUsernameBlacklistModel
	}
)

// NewUsernameBlacklistModel returns a model for the database table.
func NewUsernameBlacklistModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsernameBlacklistModel {
	return &customUsernameBlacklistModel{
		defaultUsernameBlacklistModel: newUsernameBlacklistModel(conn, c, opts...),
	}
}
