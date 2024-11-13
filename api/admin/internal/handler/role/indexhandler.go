package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zerocms/api/admin/internal/logic/role"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"
)

func IndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewIndexLogic(r.Context(), svcCtx)
		resp, err := l.Index(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
