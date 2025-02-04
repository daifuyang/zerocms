package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.RoleReq) (resp *types.IndexRoleResp, err error) {

	// 获取请求中的分页参数
	current := req.Current   // 当前页，默认值为1
	pageSize := req.PageSize // 每页数量，默认值为10

	if current < 0 {
		return nil, errors.New("分页参数错误！")
	}

	// 查询角色总数
	totalCount, err := l.svcCtx.RoleModel.Count(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("获取角色总数失败: %v", err)
	}

	// 查询角色列表
	roles, err := l.svcCtx.RoleModel.List(l.ctx, current, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取角色列表失败: %v", err)
	}

	var roleList []types.Role
	for _, role := range roles {
		// 格式化角色字段
		description := ""
		if role.Description.Valid {
			description = role.Description.String
		}

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

		// 构建角色对象
		roleList = append(roleList, types.Role{
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
		})
	}

	resp = &types.IndexRoleResp{
		Rest: new(types.Rest).Success("获取成功！"),
		Data: types.PaginateRole{
			Paginate: types.Paginate{
				Current:  current,
				PageSize: pageSize,
				Total:    totalCount,
			},
			Data: roleList,
		},
	}
	return
}
