package logic

import (
	"context"
	"fmt"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"zerocms/api/admin/internal/svc"
)

type HelloLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
	return &HelloLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloLogic) Hello() (resp *types.HelloResp, err error) {
	userId := l.ctx.Value("userId").(int64)
	fmt.Println("userId", userId)
	return &types.HelloResp{
		Rest: new(types.Rest).Success("获取成功！"),
		Data: types.HelloData{
			Id: userId,
		},
	}, nil
}
