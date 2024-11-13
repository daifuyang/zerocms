package role

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"
	"zerocms/api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type StoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoreLogic {
	return &StoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StoreLogic) Store(req *types.RoleReq) (resp *types.RoleResp, err error) {

	// 从请求中获取字段
	parentId := req.ParentId
	name := req.Name
	description := req.Description
	sort := req.Sort
	status := req.Status

	level := "0"
	// 如果上级parentId不等于0
	if parentId != 0 {
		parentRole, err := l.svcCtx.RoleModel.First(l.ctx, parentId)
		if err != nil {
			return nil, err
		}
		level = parentRole.Level.String
	}

	// 构造角色对象
	role := &model.SysRole{
		ParentId: parentId,
		Name:     name,
		Description: sql.NullString{
			String: description,
			Valid:  true,
		},
		Level: sql.NullString{
			String: level,
			Valid:  true,
		},
		Sort:   sort,
		Status: status,
	}

	// 插入角色
	insert, err := l.svcCtx.RoleModel.Insert(l.ctx, role)
	if err != nil {
		return nil, err
	}

	insertId, err := insert.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 更新path
	role.Id = insertId
	role.Level = sql.NullString{
		String: level + "-" + strconv.FormatInt(insertId, 10),
		Valid:  true,
	}
	err = l.svcCtx.RoleModel.Update(l.ctx, role)
	if err != nil {
		return nil, err
	}

	description = ""
	if role.Description.Valid {
		description = role.Description.String
	}

	level = ""
	if role.Level.Valid {
		level = role.Level.String
	}

	createdAt := time.Now().Unix()
	updatedAt := time.Now().Unix()

	createdTime := time.Unix(createdAt, 0).Format(time.DateTime)
	updatedTime := time.Unix(updatedAt, 0).Format(time.DateTime)

	resp = &types.RoleResp{
		Rest: new(types.Rest).Success("新增成功"),
		Data: types.Role{
			Id:          insertId,
			Name:        role.Name,
			Description: description,
			Level:       level,
			Sort:        role.Sort,
			Status:      role.Status,
			Date: types.Date{
				CreatedAt:   createdAt,
				CreatedTime: createdTime,
				UpdatedAt:   updatedAt,
				UpdatedTime: updatedTime,
			},
		},
	}
	return
}
