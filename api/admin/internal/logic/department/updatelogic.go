package department

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"
	"zerocms/api/model/department"

	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.DepartmentReq) (resp *types.DepartmentResp, err error) {
	// 获取当前登录用户的ID
	userId := l.ctx.Value("userId").(int64)

	// 1. 查找要更新的部门
	existDepartment, err := l.svcCtx.DepartmentModel.First(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	var status int64 = 1
	if req.Status != nil {
		status = *req.Status
	}

	// 2. 检查部门名称是否重复（除非是同一部门）
	if req.Name != existDepartment.Name {
		_, err := l.svcCtx.DepartmentModel.FindOneByName(l.ctx, req.Name)
		if err != nil && !errors.Is(err, department.ErrNotFound) {
			return nil, errors.New("部门名称已存在")
		}
	}

	// 3. 检查父级部门是否有效，如果更新了 ParentId，需要更新 ancestors
	ancestors := existDepartment.Ancestors
	if req.ParentId > 0 && req.ParentId != existDepartment.ParentId {
		existParent, err := l.svcCtx.DepartmentModel.First(l.ctx, req.ParentId)
		if err != nil {
			return nil, err
		}
		ancestors = existParent.Ancestors + "-" + strconv.FormatInt(req.Id, 10)
	}

	// 4. 更新部门信息
	existDepartment.ParentId = req.ParentId
	existDepartment.Name = req.Name
	existDepartment.Sort = req.Sort
	existDepartment.Leader = sql.NullString{
		String: req.Leader,
		Valid:  req.Leader != "",
	}
	existDepartment.Phone = sql.NullString{
		String: req.Phone,
		Valid:  req.Phone != "",
	}
	existDepartment.Email = sql.NullString{
		String: req.Email,
		Valid:  req.Email != "",
	}
	existDepartment.Status = status
	existDepartment.UpdateBy = userId
	existDepartment.Ancestors = ancestors

	// 5. 更新部门数据
	err = l.svcCtx.DepartmentModel.Update(l.ctx, existDepartment)
	if err != nil {
		return nil, err
	}

	// 6. 造响应体
	resp = &types.DepartmentResp{
		Rest: new(types.Rest).Success("更新成功！"),
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
				UpdatedAt:   time.Now().Unix(),
				UpdatedTime: time.Unix(time.Now().Unix(), 0).Format(time.DateTime),
			},
			CreateBy: userId,
			UpdateBy: userId,
		},
	}
	return resp, nil
}
