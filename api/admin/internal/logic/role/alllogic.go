package role

import (
	"context"
	"fmt"
	"time"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLogic) All(req *types.RoleReq) (resp *types.AllRoleResp, err error) {
	// 查询角色列表
	roles, err := l.svcCtx.RoleModel.List(l.ctx, 1, 0)
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
			ParentId:    role.ParentId,
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

	data := listToTree(roleList)

	resp = &types.AllRoleResp{
		Rest: new(types.Rest).Success("获取成功！"),
		Data: data,
	}
	return
}

func listToTree(roles []types.Role) []*types.Role {
	// 创建一个map来根据角色的 Id 快速查找对应的角色
	roleMap := make(map[int64]*types.Role)

	// 生成树形结构的切片
	var tree []*types.Role

	for i := range roles {
		roleMap[roles[i].Id] = &roles[i]
	}

	// 遍历角色列表，构建树形结构

	for i := range roles {
		role := &roles[i]
		// 如果角色的 ParentId 是 0 或者 ParentId 不在 roles 中，说明该角色是根节点
		if role.ParentId == 0 {
			// 将根节点加入树形结构中
			tree = append(tree, role)
		} else {
			// 将当前角色添加到其父节点的Children中
			parentRole, exists := roleMap[role.ParentId]
			if exists {
				// 这里我们确保父角色的Children已初始化
				if parentRole.Children == nil {
					parentRole.Children = make([]*types.Role, 0)
				}
				// 将子角色添加到父角色的Children
				parentRole.Children = append(parentRole.Children, role)
			}
		}
	}
	return tree
}
