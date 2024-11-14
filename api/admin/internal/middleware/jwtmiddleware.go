package middleware

import (
	"context"
	"net/http"
	"strings"
	"zerocms/api/model/user"
)

type JwtMiddleware struct {
	UserTokenModel user.SysUserTokenModel
}

func NewJwtMiddleware(userTokenModel user.SysUserTokenModel) *JwtMiddleware {
	return &JwtMiddleware{
		UserTokenModel: userTokenModel,
	}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 从请求头中获取 Authorization 字段
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// 检查 Authorization 格式，通常是 Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// 获取 token
		tokenString := parts[1]

		token, err := m.UserTokenModel.FindOneByAccessToken(r.Context(), tokenString)
		if err != nil {
			http.Error(w, "查询失败："+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", token.UserId)
		r = r.WithContext(ctx)
		// Passthrough to next handler if need
		next(w, r)
	}
}
