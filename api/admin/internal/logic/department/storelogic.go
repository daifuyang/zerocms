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

type StoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoreLogic {
	return &StoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StoreLogic) Store(req *types.DepartmentReq) (resp *types.DepartmentResp, err error) {

	userId := l.ctx.Value("userId").(int64)

	// 检查是否存在相同名称的部门（避免重复）
	_, err = l.svcCtx.DepartmentModel.FindOneByName(l.ctx, req.Name)
	if err != nil && !errors.Is(err, department.ErrNotFound) {
		return nil, errors.New("部门名称已存在")
	}

	// 检查parentId是否存在
	ancestors := "0"
	if req.ParentId > 0 {
		existParent, err := l.svcCtx.DepartmentModel.First(l.ctx, req.ParentId)
		if err != nil {
			return nil, err
		}
		ancestors = existParent.Ancestors
	}

	// 构造部门模型并插入数据

	var status int64 = 1
	if req.Status != nil {
		status = *req.Status
	}

	department := &department.SysDepartment{
		ParentId:  req.ParentId,
		Ancestors: ancestors,
		Name:      req.Name,
		Sort:      req.Sort,
		Leader: sql.NullString{
			String: req.Leader,
			Valid:  true,
		},
		Phone: sql.NullString{
			String: req.Phone,
			Valid:  true,
		},
		Email: sql.NullString{
			String: req.Email,
			Valid:  true,
		},
		Status:   status,
		CreateBy: userId,
	}

	// 插入部门数据
	insert, err := l.svcCtx.DepartmentModel.Insert(l.ctx, department)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	insertId, err := insert.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 更新ancestors

	ancestors = ancestors + "-" + strconv.FormatInt(insertId, 10)
	department.Id = insertId
	department.Ancestors = ancestors
	err = l.svcCtx.DepartmentModel.Update(l.ctx, department)
	if err != nil {
		return nil, err
	}

	// 造响应体
	resp = &types.DepartmentResp{
		Rest: new(types.Rest).Success("新增成功！"),
		Data: types.Department{
			Id:        department.Id,
			ParentId:  department.ParentId,
			Name:      department.Name,
			Ancestors: department.Ancestors,
			Sort:      department.Sort,
			Leader:    department.Leader.String,
			Phone:     department.Phone.String,
			Email:     department.Email.String,
			Status:    department.Status,
			Date: types.Date{
				CreatedAt:   time.Now().Unix(),
				CreatedTime: time.Unix(time.Now().Unix(), 0).Format(time.DateTime),
				UpdatedAt:   time.Now().Unix(),
				UpdatedTime: time.Unix(time.Now().Unix(), 0).Format(time.DateTime),
			},
			CreateBy: userId,
			UpdateBy: userId,
		},
	}
	return resp, nil
}
