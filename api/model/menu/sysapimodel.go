package menu

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"zerocms/api/model"
)

var _ SysApiModel = (*customSysApiModel)(nil)

type (
	// SysApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysApiModel.
	SysApiModel interface {
		sysApiModel
		FindOneByPathAndMethod(ctx context.Context, path string, menthod string) (*SysApi, error)
		SyncAllApi(routes []rest.Route) error
	}

	customSysApiModel struct {
		*defaultSysApiModel
	}
)

func (c customSysApiModel) SyncAllApi(routes []rest.Route) error {
	for _, route := range routes {
		sqlCtx := context.Background()
		existApi, err := c.FindOneByPathAndMethod(sqlCtx, route.Path, route.Method)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return err
		}
		path := route.Path
		method := route.Method
		// 新增
		if errors.Is(err, model.ErrNotFound) {
			_, err := c.Insert(sqlCtx, &SysApi{
				Path:   path,
				Method: method,
			})
			if err != nil {
				return err
			}
		} else {
			err := c.Update(sqlCtx, &SysApi{
				Id:          existApi.Id,
				Path:        path,
				Method:      method,
				Description: existApi.Description,
				Group:       existApi.Group,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c customSysApiModel) FindOneByPathAndMethod(ctx context.Context, path string, method string) (*SysApi, error) {
	query := fmt.Sprintf("SELECT `id`, `path`, `description`, `group`, `method` FROM %s WHERE path = ? AND method = ?", c.table)
	var resp SysApi
	err := c.CachedConn.QueryRowNoCacheCtx(ctx, &resp, query, path, method)
	return &resp, err
}

// NewSysApiModel returns a model for the database table.
func NewSysApiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SysApiModel {
	return &customSysApiModel{
		defaultSysApiModel: newSysApiModel(conn, c, opts...),
	}
}
