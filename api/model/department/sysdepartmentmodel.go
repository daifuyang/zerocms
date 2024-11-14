package department

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SysDepartmentModel = (*customSysDepartmentModel)(nil)

type (
	// SysDepartmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDepartmentModel.
	SysDepartmentModel interface {
		sysDepartmentModel
		First(ctx context.Context, id int64) (*SysDepartment, error)
		Count(ctx context.Context) (int64, error)
		List(ctx context.Context, page, limit int64) ([]*SysDepartment, error)
	}

	customSysDepartmentModel struct {
		*defaultSysDepartmentModel
	}
)

func (c customSysDepartmentModel) Count(ctx context.Context) (int64, error) {
	// 查询角色表中的总记录数
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", c.table)
	var total int64
	err := c.CachedConn.QueryRowNoCacheCtx(ctx, &total, query)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (c customSysDepartmentModel) List(ctx context.Context, page, limit int64) ([]*SysDepartment, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE deleted_at IS NULL", sysDepartmentRows, c.table)
	var resp []*SysDepartment
	var args = make([]any, 0)
	if limit != 0 {
		query += " limit ?,?"
		args = append(args, (page-1)*limit)
		args = append(args, limit)
	}
	err := c.CachedConn.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (c customSysDepartmentModel) First(ctx context.Context, id int64) (*SysDepartment, error) {
	sysDepartmentIdKey := fmt.Sprintf("%s%v", cacheSysDepartmentIdPrefix, id)
	var resp SysDepartment
	err := c.QueryRowCtx(ctx, &resp, sysDepartmentIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? AND deleted_at IS NULL limit 1", sysDepartmentRows, c.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewSysDepartmentModel returns a model for the database table.
func NewSysDepartmentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysDepartmentModel {
	return &customSysDepartmentModel{
		defaultSysDepartmentModel: newSysDepartmentModel(conn, c, opts...),
	}
}
