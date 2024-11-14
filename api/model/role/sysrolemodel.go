package role

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysRoleModel = (*customSysRoleModel)(nil)

type (
	// SysRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleModel.
	SysRoleModel interface {
		sysRoleModel
		WithSession(session sqlx.Session) SysRoleModel
		First(ctx context.Context, id int64) (*SysRole, error)
		Count(ctx context.Context) (int64, error)
		List(ctx context.Context, page, limit int64) ([]*SysRole, error)
	}

	customSysRoleModel struct {
		*defaultSysRoleModel
	}
)

func (c *customSysRoleModel) First(ctx context.Context, id int64) (*SysRole, error) {
	sysRole, err := c.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if sysRole.DeletedAt.Valid {
		return nil, sqlc.ErrNotFound
	}

	return sysRole, nil
}

func (c *customSysRoleModel) Count(ctx context.Context) (int64, error) {
	// 查询角色表中的总记录数
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", c.table)
	var total int64
	err := c.CachedConn.QueryRowNoCacheCtx(ctx, &total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (c *customSysRoleModel) List(ctx context.Context, page, limit int64) ([]*SysRole, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE deleted_at IS NULL", sysRoleRows, c.table)
	var resp []*SysRole
	var args = make([]any, 0)
	if limit != 0 {
		query += " limit ?,?"
		args = append(args, (page-1)*limit)
		args = append(args, limit)
	}
	err := c.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// NewSysRoleModel returns a model for the database table.
func NewSysRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysRoleModel {
	return &customSysRoleModel{
		defaultSysRoleModel: newSysRoleModel(conn, c, opts...),
	}
}

func (c *customSysRoleModel) WithSession(session sqlx.Session) SysRoleModel {
	return &customSysRoleModel{
		defaultSysRoleModel: &defaultSysRoleModel{
			CachedConn: c.CachedConn.WithSession(session),
			table:      "`sys_role`",
		}}
}
