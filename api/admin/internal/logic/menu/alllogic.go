package menu

import (
	"context"
	"strconv"
	"time"
	"zerocms/api/model/menu"
	"zerocms/api/model/user"

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

func (l *AllLogic) All(req *types.MenuReq) (resp *types.AllMenuResp, err error) {

	userId := l.ctx.Value("userId").(int64)
	userIdStr := strconv.FormatInt(userId, 10)

	// 获取角色
	roles, err := l.svcCtx.UserRoleModel.List(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	super := user.SuperAdmin(roles)

	var menus []*menu.SysMenu
	if super || userId == 1 {
		menus, err = l.svcCtx.MenuModel.List(l.ctx)
		if err != nil {
			return nil, err
		}
	} else {
		// 获取权限
		policy, err := l.svcCtx.Enforcer.GetFilteredPolicy(0, userIdStr)
		if err != nil {
			return nil, err
		}
		for _, access := range policy {
			menuIdStr := access[1]
			menuId, err := strconv.ParseInt(menuIdStr, 10, 64)
			if err == nil {
				menu, err := l.svcCtx.MenuModel.FindOne(l.ctx, menuId)
				if err != nil {
					return nil, err
				}
				menus = append(menus, menu)
			}
		}
	}

	var menuList []types.Menu
	for _, menu := range menus {
		// 格式化菜单字段
		createdAt := menu.CreatedAt.Unix()
		updatedAt := menu.UpdatedAt.Unix()

		// 格式化时间戳为日期时间格式
		createdTime := time.Unix(createdAt, 0).Format(time.DateTime)
		updatedTime := time.Unix(updatedAt, 0).Format(time.DateTime)

		// 构建菜单对象
		menuList = append(menuList, types.Menu{
			MenuId:    menu.MenuId,
			MenuName:  menu.MenuName,
			ParentId:  menu.ParentId,
			Order:     menu.Order,
			Path:      menu.Path,
			Component: menu.Component.String,
			Query:     menu.Query.String,
			IsFrame:   menu.IsFrame,
			IsCache:   menu.IsCache,
			MenuType:  menu.MenuType,
			Visible:   menu.Visible,
			Status:    menu.Status,
			Perms:     menu.Perms.String,
			Icon:      menu.Icon,
			CreatedId: menu.CreatedId.Int64,
			CreatedBy: menu.CreatedBy,
			UpdatedId: menu.UpdatedId.Int64,
			UpdatedBy: menu.UpdatedBy,
			Remark:    menu.Remark,
			Date: types.Date{
				CreatedAt:   createdAt,
				CreatedTime: createdTime,
				UpdatedAt:   updatedAt,
				UpdatedTime: updatedTime,
			},
		})
	}

	// 将菜单列表转换为树形结构
	data := listToTree(menuList)

	// 返回响应
	resp = &types.AllMenuResp{
		Rest: new(types.Rest).Success("获取成功！"),
		Data: data,
	}
	return
}

// 构建菜单树形结构
func listToTree(menus []types.Menu) []*types.Menu {
	// 创建一个map来根据菜单的 MenuId 快速查找对应的菜单
	menuMap := make(map[int64]*types.Menu)

	// 生成树形结构的切片
	var tree []*types.Menu

	for i := range menus {
		menuMap[menus[i].MenuId] = &menus[i]
	}

	// 遍历菜单列表，构建树形结构
	for i := range menus {
		menu := &menus[i]
		// 如果菜单的 ParentId 是 0 或者 ParentId 不在 menus 中，说明该菜单是根节点
		if menu.ParentId == 0 {
			// 将根节点加入树形结构中
			tree = append(tree, menu)
		} else {
			// 将当前菜单添加到其父节点的 Children 中
			parentMenu, exists := menuMap[menu.ParentId]
			if exists {
				// 确保父菜单的 Children 已初始化
				if parentMenu.Children == nil {
					parentMenu.Children = make([]*types.Menu, 0)
				}
				// 将子菜单添加到父菜单的 Children
				parentMenu.Children = append(parentMenu.Children, menu)
			}
		}
	}
	return tree
}
