package menu

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysMenuApiModel = (*customSysMenuApiModel)(nil)

type (
	defaultSysMenuApiModel struct {
		sqlc.CachedConn
		table string
	}

	SysMenuApi struct {
		MenuId int64 `db:"menu_id"` // 菜单ID
		ApiId  int64 `db:"api_id"`  // API ID
	}

	// SysMenuApiModel 是一个接口，可以根据需要自定义，添加更多方法，
	// 并在 customSysMenuApiModel 中实现新增的方法。
	SysMenuApiModel interface {
		Insert(ctx context.Context, data *SysMenuApi) (sql.Result, error)
		List(ctx context.Context, menuId int64) ([]*SysMenuApi, error)
		FindOneByMenuIdAndApiId(ctx context.Context, menuId int64, apiId int64) (*SysMenuApi, error)
		DeleteByMenuIdAndApiId(ctx context.Context, menuId int64, apiId int64) error
	}

	customSysMenuApiModel struct {
		*defaultSysMenuApiModel
	}
)

func (c customSysMenuApiModel) FindOneByMenuIdAndApiId(ctx context.Context, menuId int64, apiId int64) (*SysMenuApi, error) {
	query := fmt.Sprintf("SELECT menu_id, api_id FROM %s WHERE menu_id = ? AND api_id = ?", c.table)
	var resp SysMenuApi
	err := c.CachedConn.QueryRowNoCacheCtx(ctx, &resp, query, menuId, apiId)
	return &resp, err
}

func (c customSysMenuApiModel) DeleteByMenuIdAndApiId(ctx context.Context, menuId int64, apiId int64) error {
	_, err := c.FindOneByMenuIdAndApiId(ctx, menuId, apiId)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE menu_id = ? AND api_id = ?", c.table)
	_, err = c.ExecNoCacheCtx(ctx, query, menuId, apiId)
	return err
}

func (c customSysMenuApiModel) List(ctx context.Context, menuId int64) ([]*SysMenuApi, error) {
	query := fmt.Sprintf("SELECT menu_id, api_id FROM %s WHERE menu_id = ?", c.table)
	var resp []*SysMenuApi
	err := c.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, menuId)
	return resp, err
}

func (m *defaultSysMenuApiModel) Insert(ctx context.Context, data *SysMenuApi) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (`menu_id`, `api_id`) VALUES (?, ?)", m.table)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.MenuId, data.ApiId)
	return ret, err
}

func newSysMenuApiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultSysMenuApiModel {
	return &defaultSysMenuApiModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`sys_menu_api`", // 修改为 sys_menu_api 表
	}
}

// NewSysMenuApiModel 返回一个模型，用于操作数据库表 sys_menu_api。
func NewSysMenuApiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysMenuApiModel {
	return &customSysMenuApiModel{
		defaultSysMenuApiModel: newSysMenuApiModel(conn, c, opts...),
	}
}
