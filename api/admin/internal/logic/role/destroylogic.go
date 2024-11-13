package role

import (
	"context"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DestroyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDestroyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DestroyLogic {
	return &DestroyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DestroyLogic) Destroy() (resp *types.RoleResp, err error) {
	// todo: add your logic here and delete this line

	return
}
