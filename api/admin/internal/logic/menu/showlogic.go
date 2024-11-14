package menu

import (
	"context"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowMenuReq) (resp *types.MenuResp, err error) {
	// todo: add your logic here and delete this line

	return
}
