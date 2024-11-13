package migrate

import (
	"zerocms/api/model"
)

type ModelContext struct {
	UserModel model.SysUserModel
}

func Init(ctx *ModelContext) {
	createAdmin(ctx.UserModel)
}
