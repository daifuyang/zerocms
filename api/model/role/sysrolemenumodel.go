package role

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRoleMenuModel = (*customSysRoleMenuModel)(nil)

type (
	defaultSysRoleMenuModel struct {
		sqlc.CachedConn
		table string
	}

	SysRoleMenu struct {
		RoleId int64 `db:"role_id"` // 角色ID
		MenuId int64 `db:"menu_id"` // 菜单ID
	}

	// SysRoleMenuModel 是一个接口，可以根据需要自定义，添加更多方法，
	// 并在 customSysRoleMenuModel 中实现新增的方法。
	SysRoleMenuModel interface {
		Insert(ctx context.Context, data *SysRoleMenu) (sql.Result, error)
		List(ctx context.Context, roleId int64) ([]*SysRoleMenu, error)
		FindOneByRoleIdAndMenuId(ctx context.Context, roleId int64, menuId int64) (*SysRoleMenu, error)
		DeleteByRoleIdAndMenuId(ctx context.Context, roleId int64, menuId int64) error
	}

	customSysRoleMenuModel struct {
		*defaultSysRoleMenuModel
	}
)

func (c customSysRoleMenuModel) FindOneByRoleIdAndMenuId(ctx context.Context, roleId int64, menuId int64) (*SysRoleMenu, error) {
	query := fmt.Sprintf("SELECT role_id, menu_id FROM %s WHERE role_id = ? AND menu_id = ?", c.table)
	var resp SysRoleMenu
	err := c.CachedConn.QueryRowNoCacheCtx(ctx, &resp, query, roleId, menuId)
	return &resp, err
}

func (c customSysRoleMenuModel) DeleteByRoleIdAndMenuId(ctx context.Context, roleId int64, menuId int64) error {
	_, err := c.FindOneByRoleIdAndMenuId(ctx, roleId, menuId)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE role_id = ? AND menu_id = ?", c.table)
	_, err = c.ExecNoCacheCtx(ctx, query, roleId, menuId)
	return err
}

func (c customSysRoleMenuModel) List(ctx context.Context, roleId int64) ([]*SysRoleMenu, error) {
	query := fmt.Sprintf("SELECT role_id, menu_id FROM %s WHERE role_id = ?", c.table)
	var resp []*SysRoleMenu
	err := c.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, roleId)
	return resp, err
}

func (m *defaultSysRoleMenuModel) Insert(ctx context.Context, data *SysRoleMenu) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (`role_id`, `menu_id`) VALUES (?, ?)", m.table)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.RoleId, data.MenuId)
	return ret, err
}

func newSysRoleMenuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultSysRoleMenuModel {
	return &defaultSysRoleMenuModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`sys_role_menu`",
	}
}

// NewSysRoleMenuModel 返回一个模型，用于操作数据库表 sys_role_menu。
func NewSysRoleMenuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysRoleMenuModel {
	return &customSysRoleMenuModel{
		defaultSysRoleMenuModel: newSysRoleMenuModel(conn, c, opts...),
	}
}
