package department

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zerocms/api/admin/internal/logic/department"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"
)

func AllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DepartmentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := department.NewAllLogic(r.Context(), svcCtx)
		resp, err := l.All(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
