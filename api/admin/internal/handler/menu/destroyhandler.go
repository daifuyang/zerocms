package menu

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zerocms/api/admin/internal/logic/menu"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"
)

func DestroyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShowMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := menu.NewDestroyLogic(r.Context(), svcCtx)
		resp, err := l.Destroy(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
