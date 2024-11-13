package role

import (
	"context"
	"fmt"
	"time"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRoleReq) (resp *types.RoleResp, err error) {
	// 获取请求中的角色 ID
	id := req.Id

	// 查询角色数据
	role, err := l.svcCtx.RoleModel.First(l.ctx, id)
	if err != nil {
		// 如果查询角色失败，可以返回角色不存在的错误
		return nil, fmt.Errorf("角色不存在或查询失败: %v", err)
	}

	// 格式化角色的描述字段（如果有）
	description := ""
	if role.Description.Valid {
		description = role.Description.String
	}

	// 格式化角色的 level 字段（如果有）
	level := ""
	if role.Level.Valid {
		level = role.Level.String
	}

	// 获取角色的创建和更新时间（Unix 时间戳）
	createdAt := role.CreatedAt.Unix()
	updatedAt := role.UpdatedAt.Unix()

	// 格式化时间戳为日期时间格式
	createdTime := time.Unix(createdAt, 0).Format(time.DateTime)
	updatedTime := time.Unix(updatedAt, 0).Format(time.DateTime)

	// 返回角色的响应
	resp = &types.RoleResp{
		Rest: new(types.Rest).Success("查询成功"),
		Data: types.Role{
			Id:          role.Id,
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
