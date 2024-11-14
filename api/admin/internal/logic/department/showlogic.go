package department

import (
	"context"
	"errors"
	"time"
	"zerocms/api/model/department"

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

func (l *ShowLogic) Show(req *types.ShowDepartmentReq) (resp *types.DepartmentResp, err error) {
	one, err := l.svcCtx.DepartmentModel.First(l.ctx, req.Id)
	if err != nil {
		if errors.Is(err, department.ErrNotFound) {
			return nil, errors.New("部门不存在")
		}
		return nil, err
	}
	resp = &types.DepartmentResp{
		Rest: new(types.Rest).Success("获取成功！"),
		Data: types.Department{
			Id:        one.Id,
			ParentId:  one.ParentId,
			Name:      one.Name,
			Ancestors: one.Ancestors,
			Sort:      one.Sort,
			Leader:    one.Leader.String,
			Phone:     one.Phone.String,
			Email:     one.Email.String,
			Status:    one.Status,
			Date: types.Date{
				CreatedAt:   one.CreatedAt.Unix(),
				CreatedTime: one.CreatedAt.Format(time.DateTime),
				UpdatedAt:   one.UpdatedAt.Unix(),
				UpdatedTime: one.UpdatedAt.Format(time.DateTime),
			},
			CreateBy: one.CreateBy,
			UpdateBy: one.UpdateBy,
		},
	}
	return resp, nil
}
