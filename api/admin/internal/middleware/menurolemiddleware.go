package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"net/http"
	"strconv"
	"zerocms/api/model/menu"
	"zerocms/api/model/role"
)

type MenuRoleMiddleware struct {
	MenuModel     menu.SysMenuModel
	RoleMenuModel role.SysRoleMenuModel
	Enforcer      *casbin.Enforcer
}

func NewMenuRoleMiddleware(menuModel menu.SysMenuModel, roleMenuModel role.SysRoleMenuModel, enforcer *casbin.Enforcer) *MenuRoleMiddleware {
	return &MenuRoleMiddleware{
		MenuModel:     menuModel,
		RoleMenuModel: roleMenuModel,
		Enforcer:      enforcer,
	}
}

func (m *MenuRoleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取当前用户信息
		userId := r.Context().Value("userId").(int64)
		userIdStr := strconv.FormatInt(userId, 10)

		fmt.Println("userIdStr", userIdStr)
		currentURL := r.URL.String()
		method := r.Method
		fmt.Println("currentURL", currentURL, method)

		// Passthrough to next handler if need
		next(w, r)
	}
}
