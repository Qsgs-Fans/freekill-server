package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserLoginHistoryModel = (*customUserLoginHistoryModel)(nil)

type (
	// UserLoginHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLoginHistoryModel.
	UserLoginHistoryModel interface {
		userLoginHistoryModel
	}

	customUserLoginHistoryModel struct {
		*defaultUserLoginHistoryModel
	}
)

// NewUserLoginHistoryModel returns a model for the database table.
func NewUserLoginHistoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserLoginHistoryModel {
	return &customUserLoginHistoryModel{
		defaultUserLoginHistoryModel: newUserLoginHistoryModel(conn, c, opts...),
	}
}
