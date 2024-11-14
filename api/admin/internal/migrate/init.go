package migrate

import (
	"zerocms/api/model/menu"
	"zerocms/api/model/user"
)

type ModelContext struct {
	UserModel user.SysUserModel
	MenuModel menu.SysMenuModel
}

func Init(ctx *ModelContext) {
	createAdmin(ctx.UserModel)
	createMenus(ctx.MenuModel)
}
