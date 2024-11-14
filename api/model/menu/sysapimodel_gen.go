// Code generated by goctl. DO NOT EDIT.

package menu

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	sysApiFieldNames          = builder.RawFieldNames(&SysApi{})
	sysApiRows                = strings.Join(sysApiFieldNames, ",")
	sysApiRowsExpectAutoSet   = strings.Join(stringx.Remove(sysApiFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	sysApiRowsWithPlaceHolder = strings.Join(stringx.Remove(sysApiFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheSysApiIdPrefix = "cache:sysApi:id:"
)

type (
	sysApiModel interface {
		Insert(ctx context.Context, data *SysApi) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*SysApi, error)
		Update(ctx context.Context, data *SysApi) error
		Delete(ctx context.Context, id int64) error
	}

	defaultSysApiModel struct {
		sqlc.CachedConn
		table string
	}

	SysApi struct {
		Id          int64  `db:"id"`          // 接口id
		Path        string `db:"path"`        // 路径
		Description string `db:"description"` // 描述
		Group       string `db:"group"`       // 分组
		Method      string `db:"method"`      // 请求方法
	}
)

func newSysApiModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultSysApiModel {
	return &defaultSysApiModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`sys_api`",
	}
}

func (m *defaultSysApiModel) withSession(session sqlx.Session) *defaultSysApiModel {
	return &defaultSysApiModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`sys_api`",
	}
}

func (m *defaultSysApiModel) Delete(ctx context.Context, id int64) error {
	sysApiIdKey := fmt.Sprintf("%s%v", cacheSysApiIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, sysApiIdKey)
	return err
}

func (m *defaultSysApiModel) FindOne(ctx context.Context, id int64) (*SysApi, error) {
	sysApiIdKey := fmt.Sprintf("%s%v", cacheSysApiIdPrefix, id)
	var resp SysApi
	err := m.QueryRowCtx(ctx, &resp, sysApiIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", sysApiRows, m.table)
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

func (m *defaultSysApiModel) Insert(ctx context.Context, data *SysApi) (sql.Result, error) {
	sysApiIdKey := fmt.Sprintf("%s%v", cacheSysApiIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, sysApiRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Path, data.Description, data.Group, data.Method)
	}, sysApiIdKey)
	return ret, err
}

func (m *defaultSysApiModel) Update(ctx context.Context, data *SysApi) error {
	sysApiIdKey := fmt.Sprintf("%s%v", cacheSysApiIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, sysApiRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Path, data.Description, data.Group, data.Method, data.Id)
	}, sysApiIdKey)
	return err
}

func (m *defaultSysApiModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheSysApiIdPrefix, primary)
}

func (m *defaultSysApiModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", sysApiRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultSysApiModel) tableName() string {
	return m.table
}
