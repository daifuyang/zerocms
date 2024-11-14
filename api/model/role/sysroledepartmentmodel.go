package role

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRoleDepartmentModel = (*customSysRoleDepartmentModel)(nil)

type (
	defaultSysRoleDepartmentModel struct {
		sqlc.CachedConn
		table string
	}

	SysRoleDepartment struct {
		RoleId       int64 `db:"role_id"`       // 角色ID
		DepartmentId int64 `db:"department_id"` // 部门ID
	}

	// SysRoleDepartmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleDepartmentModel.
	SysRoleDepartmentModel interface {
		Insert(ctx context.Context, data *SysRoleDepartment) (sql.Result, error)
		List(ctx context.Context, roleId int64) ([]*SysRoleDepartment, error)
		First(ctx context.Context, roleId int64, departmentId int64) (*SysRoleDepartment, error)
		DeleteByRoleIdAndDepartmentId(ctx context.Context, roleId int64, departmentId int64) error
	}

	customSysRoleDepartmentModel struct {
		*defaultSysRoleDepartmentModel
	}
)

func (c customSysRoleDepartmentModel) First(ctx context.Context, roleId int64, departmentId int64) (*SysRoleDepartment, error) {
	query := fmt.Sprintf("SELECT role_id, department_id FROM %s WHERE role_id = ? AND department_id = ?", c.table)
	var resp SysRoleDepartment
	err := c.CachedConn.QueryRowNoCacheCtx(ctx, &resp, query, roleId, departmentId)
	return &resp, err
}

func (c customSysRoleDepartmentModel) DeleteByRoleIdAndDepartmentId(ctx context.Context, roleId int64, departmentId int64) error {
	_, err := c.First(ctx, roleId, departmentId)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("delete from %s  WHERE role_id = ? AND department_id = ?", c.table)
	_, err = c.ExecNoCacheCtx(ctx, query, roleId, departmentId)
	return err
}

func (c customSysRoleDepartmentModel) List(ctx context.Context, roleId int64) ([]*SysRoleDepartment, error) {
	query := fmt.Sprintf("SELECT role_id, department_id FROM %s WHERE role_id = ?", c.table)
	var resp []*SysRoleDepartment
	err := c.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, roleId)
	return resp, err
}

func (m *defaultSysRoleDepartmentModel) Insert(ctx context.Context, data *SysRoleDepartment) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`role_id`, `department_id`) values (?, ?)", m.table)
	ret, err := m.ExecNoCacheCtx(ctx, query, data.RoleId, data.DepartmentId)
	return ret, err
}
func newSysRoleDepartmentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultSysRoleDepartmentModel {
	return &defaultSysRoleDepartmentModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`sys_role_department`",
	}
}

// NewSysRoleDepartmentModel returns a model for the database table.
func NewSysRoleDepartmentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysRoleDepartmentModel {
	return &customSysRoleDepartmentModel{
		defaultSysRoleDepartmentModel: newSysRoleDepartmentModel(conn, c, opts...),
	}
}
