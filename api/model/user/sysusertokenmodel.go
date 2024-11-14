package user

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysUserTokenModel = (*customSysUserTokenModel)(nil)

type (
	// SysUserTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysUserTokenModel.
	SysUserTokenModel interface {
		sysUserTokenModel
	}

	customSysUserTokenModel struct {
		*defaultSysUserTokenModel
	}
)

// NewSysUserTokenModel returns a model for the database table.
func NewSysUserTokenModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysUserTokenModel {
	return &customSysUserTokenModel{
		defaultSysUserTokenModel: newSysUserTokenModel(conn, c, opts...),
	}
}
