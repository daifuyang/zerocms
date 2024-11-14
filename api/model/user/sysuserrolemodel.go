package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysUserRoleModel = (*customSysUserRoleModel)(nil)

type (
	defaultSysUserRoleModel struct {
		sqlc.CachedConn
		table string
	}

	SysUserRole struct {
		UserId int64 `db:"user_id"` // 用户ID
		RoleId int64 `db:"role_id"` // 角色ID
	}

	// SysUserRoleModel 是一个接口，可以根据需要自定义，添加更多方法，
	// 并在 customSysUserRoleModel 中实现新增的方法。
	SysUserRoleModel interface {
		Insert(ctx context.Context, data *SysUserRole) (sql.Result, error)
		List(ctx context.Context, userId int64) ([]*SysUserRole, error)
		First(ctx context.Context, userId int64, roleId int64) (*SysUserRole, error)
		DeleteByUserIdAndRoleId(ctx context.Context, userId int64, roleId int64) error
	}

	customSysUserRoleModel struct {
		*defaultSysUserRoleModel
	}
)

func SuperAdmin(userRole []*SysUserRole) bool {
	for _, role := range userRole {
		if role.RoleId == 1 {
			return true
		}
	}
	return false
}

func (c customSysUserRoleModel) First(ctx context.Context, userId int64, roleId int64) (*SysUserRole, error) {
	query := fmt.Sprintf("SELECT user_id, role_id FROM %s WHERE user_id = ? AND role_id = ?", c.table)
	var resp SysUserRole
	err := c.CachedConn.QueryRowNoCacheCtx(ctx, &resp, query, userId, roleId)
	return &resp, err
}

func (c customSysUserRoleModel) DeleteByUserIdAndRoleId(ctx context.Context, userId int64, roleId int64) error {
	_, err := c.First(ctx, userId, roleId)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND role_id = ?", c.table)
	_, err = c.ExecNoCacheCtx(ctx, query, userId, roleId)
	return err
}

func (c customSysUserRoleModel) List(ctx context.Context, userId int64) ([]*SysUserRole, error) {
	query := fmt.Sprintf("SELECT user_id, role_id FROM %s WHERE user_id = ?", c.table)
	var resp []*SysUserRole
	err := c.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	return resp, err
}

func (m *defaultSysUserRoleModel) Insert(ctx context.Context, data *SysUserRole) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (`user_id`, `role_id`) VALUES (?, ?)", m.table)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.UserId, data.RoleId)
	return ret, err
}

func newSysUserRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultSysUserRoleModel {
	return &defaultSysUserRoleModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`sys_user_role`", // 修改为 sys_user_role 表
	}
}

// NewSysUserRoleModel 返回一个模型，用于操作数据库表 sys_user_role。
func NewSysUserRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysUserRoleModel {
	return &customSysUserRoleModel{
		defaultSysUserRoleModel: newSysUserRoleModel(conn, c, opts...),
	}
}

func (m *defaultSysUserRoleModel) withSession(session sqlx.Session) *defaultSysUserRoleModel {
	return &defaultSysUserRoleModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`sys_user`",
	}
}
