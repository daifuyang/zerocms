package migrate

import (
	"context"
	"database/sql"
	"encoding/json"
	"zerocms/api/admin/internal/types"
	"zerocms/api/model/menu"
	"zerocms/utils"
)

// 创建默认菜单

func createMenus(model menu.SysMenuModel) {

	file, err := utils.ReadFile("public/install/menus.json")
	if err != nil {
		panic(err)
	}
	var menus []*types.Menu
	err = json.Unmarshal(file, &menus)
	if err != nil {
		panic(err)
	}

	m := menuModel{
		SysMenuModel: model,
	}

	m.recursionMenu(menus, 0)

}

type menuModel struct {
	menu.SysMenuModel
}

// 递归菜单
func (m *menuModel) recursionMenu(menus []*types.Menu, parentId int64) {
	for _, item := range menus {

		existMenu, err := m.FindOneByPerms(context.Background(), sql.NullString{String: item.Perms, Valid: true})
		if err != nil {
			return
		}

		var menuId int64 = existMenu.MenuId

		saveModel := &menu.SysMenu{
			MenuName: item.MenuName,
			ParentId: parentId,
			Order:    item.Order,
			Path:     item.Path,
			Component: sql.NullString{
				String: item.Component,
				Valid:  true,
			},
			Query: sql.NullString{
				String: item.Query,
				Valid:  true,
			},
			IsFrame:  item.IsFrame,
			IsCache:  item.IsCache,
			MenuType: item.MenuType,
			Visible:  item.Visible,
			Status:   item.Status,
			Perms: sql.NullString{
				String: item.Perms,
				Valid:  true,
			},
			Icon: item.Icon,
		}

		if existMenu != nil {
			saveModel.MenuId = existMenu.MenuId
			m.Update(context.Background(), saveModel)
		} else {
			insert, err := m.Insert(context.Background(), saveModel)
			if err != nil {
				panic(err)
			}

			menuId, err = insert.LastInsertId()
			if err != nil {
				panic(err)
			}
		}

		if item.Children != nil {
			m.recursionMenu(item.Children, menuId)
		}
	}
}
