package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsernameWhitelistModel = (*customUsernameWhitelistModel)(nil)

type (
	// UsernameWhitelistModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsernameWhitelistModel.
	UsernameWhitelistModel interface {
		usernameWhitelistModel
	}

	customUsernameWhitelistModel struct {
		*defaultUsernameWhitelistModel
	}
)

// NewUsernameWhitelistModel returns a model for the database table.
func NewUsernameWhitelistModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsernameWhitelistModel {
	return &customUsernameWhitelistModel{
		defaultUsernameWhitelistModel: newUsernameWhitelistModel(conn, c, opts...),
	}
}
