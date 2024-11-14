package department

import (
	"context"
	"fmt"
	"time"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLogic) All(req *types.DepartmentReq) (resp *types.AllDepartmentResp, err error) {
	// 查询部门列表
	departments, err := l.svcCtx.DepartmentModel.List(l.ctx, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("获取部门列表失败: %v", err)
	}

	var departmentList []types.Department
	for _, department := range departments {
		// 格式化部门字段
		leader := ""
		if department.Leader.Valid {
			leader = department.Leader.String
		}

		phone := ""
		if department.Phone.Valid {
			phone = department.Phone.String
		}

		email := ""
		if department.Email.Valid {
			email = department.Email.String
		}

		// 获取部门的创建和更新时间（Unix 时间戳）
		createdAt := department.CreatedAt.Unix()
		updatedAt := department.UpdatedAt.Unix()

		// 格式化时间戳为日期时间格式
		createdTime := time.Unix(createdAt, 0).Format(time.DateTime)
		updatedTime := time.Unix(updatedAt, 0).Format(time.DateTime)

		// 构建部门对象
		departmentList = append(departmentList, types.Department{
			Id:        department.Id,
			ParentId:  department.ParentId,
			Name:      department.Name,
			Ancestors: department.Ancestors,
			Sort:      department.Sort,
			Leader:    leader,
			Phone:     phone,
			Email:     email,
			Status:    department.Status,
			Date: types.Date{
				CreatedAt:   createdAt,
				CreatedTime: createdTime,
				UpdatedAt:   updatedAt,
				UpdatedTime: updatedTime,
			},
			CreateBy: department.CreateBy,
			UpdateBy: department.UpdateBy,
		})
	}

	// 将部门列表转换为树形结构
	data := listToTree(departmentList)

	resp = &types.AllDepartmentResp{
		Rest: new(types.Rest).Success("获取成功！"),
		Data: data,
	}
	return
}

func listToTree(departments []types.Department) []*types.Department {
	// 创建一个map来根据部门的 Id 快速查找对应的部门
	departmentMap := make(map[int64]*types.Department)

	// 生成树形结构的切片
	var tree []*types.Department

	// 遍历部门列表，将每个部门的指针放入map中
	for i := range departments {
		departmentMap[departments[i].Id] = &departments[i]
	}

	// 遍历部门列表，构建树形结构
	for i := range departments {
		department := &departments[i]
		// 如果部门的 ParentId 是 0 或者 ParentId 不在 departments 中，说明该部门是根节点
		if department.ParentId == 0 {
			// 将根节点加入树形结构中
			tree = append(tree, department)
		} else {
			// 将当前部门添加到其父部门的 Children 中
			parentDepartment, exists := departmentMap[department.ParentId]
			if exists {
				// 这里我们确保父部门的 Children 已初始化
				if parentDepartment.Children == nil {
					parentDepartment.Children = make([]*types.Department, 0)
				}
				// 将子部门添加到父部门的 Children
				parentDepartment.Children = append(parentDepartment.Children, department)
			}
		}
	}
	return tree
}
