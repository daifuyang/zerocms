package role

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.RoleReq) (resp *types.RoleResp, err error) {
	// 从请求中获取字段
	id := req.Id // 角色ID
	parentId := req.ParentId
	name := req.Name
	description := req.Description
	sort := req.Sort
	status := req.Status

	// 查找角色是否存在
	role, err := l.svcCtx.RoleModel.First(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("角色不存在或无法查找: %v", err)
	}

	// 如果上级parentId不等于0，更新父角色的level
	level := "0"
	if parentId != 0 {
		parentRole, err := l.svcCtx.RoleModel.FindOne(l.ctx, parentId)
		if err != nil {
			return nil, fmt.Errorf("父角色查找失败: %v", err)
		}
		level = parentRole.Level.String
	}

	// 更新角色字段
	role.ParentId = parentId
	role.Name = name
	role.Description = sql.NullString{
		String: description,
		Valid:  true,
	}
	role.Sort = sort
	role.Status = status

	// 更新角色的 Level 字段
	role.Level = sql.NullString{
		String: level,
		Valid:  true,
	}

	// 更新角色
	err = l.svcCtx.RoleModel.Update(l.ctx, role)
	if err != nil {
		return nil, fmt.Errorf("角色更新失败: %v", err)
	}

	// 如果角色更新成功，获取角色的更新时间
	updatedAt := time.Now().Unix()
	updatedTime := time.Unix(updatedAt, 0).Format(time.DateTime)

	// 生成响应数据
	description = ""
	if role.Description.Valid {
		description = role.Description.String
	}

	level = ""
	if role.Level.Valid {
		level = role.Level.String
	}

	resp = &types.RoleResp{
		Rest: new(types.Rest).Success("更新成功"),
		Data: types.Role{
			Id:          role.Id,
			Name:        role.Name,
			Description: description,
			Level:       level,
			Sort:        role.Sort,
			Status:      role.Status,
			Date: types.Date{
				UpdatedAt:   updatedAt,
				UpdatedTime: updatedTime,
			},
		},
	}

	return
}
