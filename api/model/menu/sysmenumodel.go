package menu

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysMenuModel = (*customSysMenuModel)(nil)

type (
	// SysMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysMenuModel.
	SysMenuModel interface {
		sysMenuModel
		List(ctx context.Context) ([]*SysMenu, error)
	}

	customSysMenuModel struct {
		*defaultSysMenuModel
	}
)

func (c customSysMenuModel) List(ctx context.Context) ([]*SysMenu, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE deleted_at IS NULL", sysMenuRows, c.table)
	var resp []*SysMenu
	err := c.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query)
	return resp, err
}

// NewSysMenuModel returns a model for the database table.
func NewSysMenuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysMenuModel {
	return &customSysMenuModel{
		defaultSysMenuModel: newSysMenuModel(conn, c, opts...),
	}
}
