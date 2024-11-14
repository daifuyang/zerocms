package department

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"
	"zerocms/api/model/department"

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

func (l *DestroyLogic) Destroy(req *types.ShowDepartmentReq) (resp *types.DepartmentResp, err error) {
	// 获取当前登录用户的ID
	userId := l.ctx.Value("userId").(int64)

	// 查询部门存在
	existDepartment, err := l.svcCtx.DepartmentModel.First(l.ctx, req.Id)
	if err != nil {
		if errors.Is(err, department.ErrNotFound) {
			return nil, errors.New("部门不存在")
		}
		return nil, err
	}

	// 删除部门
	existDepartment.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	existDepartment.UpdateBy = userId

	// 5. 更新部门数据
	err = l.svcCtx.DepartmentModel.Update(l.ctx, existDepartment)
	if err != nil {
		return nil, err
	}

	// 6. 造响应体
	resp = &types.DepartmentResp{
		Rest: new(types.Rest).Success("删除成功！"),
		Data: types.Department{
			Id:        existDepartment.Id,
			ParentId:  existDepartment.ParentId,
			Name:      existDepartment.Name,
			Ancestors: existDepartment.Ancestors,
			Sort:      existDepartment.Sort,
			Leader:    existDepartment.Leader.String,
			Phone:     existDepartment.Phone.String,
			Email:     existDepartment.Email.String,
			Status:    existDepartment.Status,
			Date: types.Date{
				CreatedAt:   existDepartment.CreatedAt.Unix(),
				CreatedTime: time.Unix(existDepartment.CreatedAt.Unix(), 0).Format(time.DateTime),
				UpdatedAt:   existDepartment.UpdatedAt.Unix(),
				UpdatedTime: time.Unix(existDepartment.UpdatedAt.Unix(), 0).Format(time.DateTime),
			},
			CreateBy:  existDepartment.CreateBy,
			UpdateBy:  existDepartment.UpdateBy,
			DeletedAt: existDepartment.DeletedAt.Time.Unix(),
		},
	}
	return resp, nil
}
