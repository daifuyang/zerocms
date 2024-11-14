package role

import (
	"context"
	"errors"
	"strings"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"
	"zerocms/api/model/department"
	"zerocms/api/model/role"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataScopeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataScopeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataScopeLogic {
	return &DataScopeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataScopeLogic) DataScope(req *types.DataScopeReq) (resp *types.DataScopeResp, err error) {
	roleId := req.Id

	roleOne, err := l.svcCtx.RoleModel.First(l.ctx, roleId)
	if err != nil {
		return nil, err
	}

	dataScope := req.DataScope
	departmentIds := req.DepartmentIds
	if dataScope == 2 {
		if len(departmentIds) == 0 {
			return nil, errors.New("必须选择部门")
		}

		// 获取所有选择的关系
		list, err := l.svcCtx.RoleDepartmentModel.List(l.ctx, roleId)
		if err != nil {
			return nil, err
		}

		inputMap := make(map[int64]bool)
		for _, item := range departmentIds {
			inputMap[*item] = true
		}

		existMap := make(map[int64]bool)
		for _, item := range list {
			existMap[item.DepartmentId] = true
		}

		var toRemove = make([]*role.SysRoleDepartment, 0)
		var toAdd = make([]*role.SysRoleDepartment, 0)

		// 如果本地不存在，则为全部添加
		if len(list) == 0 {
			for _, departmentId := range departmentIds {
				toAdd = append(toAdd, &role.SysRoleDepartment{
					RoleId:       roleId,
					DepartmentId: *departmentId,
				}) // 新增的元素
			}
		} else {
			// 查找要添加的元素（添加数组中有，但现有数组中没有的）
			for _, departmentId := range departmentIds {
				if _, exists := existMap[*departmentId]; !exists {
					toAdd = append(toAdd, &role.SysRoleDepartment{
						RoleId:       roleId,
						DepartmentId: *departmentId,
					})
				}
			}
		}

		// 如果传入的为空则全部删除
		if len(departmentIds) == 0 {
			for _, item := range list {
				toRemove = append(toRemove, item)
			}
		} else {
			for _, item := range list {
				if _, exists := inputMap[item.DepartmentId]; !exists {
					toRemove = append(toRemove, item)
				}
			}
		}

		// 遍历部门
		notFoundDepartments := make([]string, 0)
		for _, departmentId := range departmentIds {
			// 获取部门信息
			departmentOne, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, *departmentId)
			if err != nil {
				if errors.Is(err, department.ErrNotFound) {
					notFoundDepartments = append(notFoundDepartments, departmentOne.Name)
				} else {
					return nil, err
				}
			}
		}

		if len(notFoundDepartments) > 0 {
			return nil, errors.New("部门不存在：" + strings.Join(notFoundDepartments, ","))
		}

		for _, roleDepartment := range toAdd {
			_, err := l.svcCtx.RoleDepartmentModel.Insert(l.ctx, &role.SysRoleDepartment{
				RoleId:       roleId,
				DepartmentId: roleDepartment.DepartmentId,
			})
			if err != nil {
				return nil, err
			}
		}

		for _, roleDepartment := range toRemove {
			err := l.svcCtx.RoleDepartmentModel.DeleteByRoleIdAndDepartmentId(l.ctx, roleDepartment.RoleId, roleDepartment.DepartmentId)
			if err != nil {
				return nil, err
			}
		}
	}

	resp = &types.DataScopeResp{
		Rest: new(types.Rest).Success("操作成功！"),
		Data: types.DataScope{
			Id:            roleId,
			DataScope:     dataScope,
			DepartmentIds: departmentIds,
			RoleName:      roleOne.Name,
		},
	}
	return
}
