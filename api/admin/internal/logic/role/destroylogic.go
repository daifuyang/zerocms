package role

import (
	"context"
	"database/sql"
	"time"
	"zerocms/api/model"

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

func (l *DestroyLogic) Destroy(req *types.ShowRoleReq) (resp *types.RoleResp, err error) {
	id := req.Id
	existRole, err := l.svcCtx.RoleModel.First(l.ctx, id)
	if err != nil {
		return nil, err
	}

	deleteRole := &model.SysRole{
		Id:          id,
		ParentId:    existRole.ParentId,
		Name:        existRole.Name,
		Description: existRole.Description,
		Level:       existRole.Level,
		Sort:        existRole.Sort,
		Status:      existRole.Status,
		DeletedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err = l.svcCtx.RoleModel.Update(l.ctx, deleteRole)
	if err != nil {
		return nil, err
	}

	// 获取角色的创建和更新时间（Unix 时间戳）
	createdAt := existRole.CreatedAt.Unix()
	updatedAt := existRole.UpdatedAt.Unix()

	// 格式化时间戳为日期时间格式
	createdTime := time.Unix(createdAt, 0).Format(time.DateTime)
	updatedTime := time.Unix(updatedAt, 0).Format(time.DateTime)

	resp = &types.RoleResp{
		Rest: new(types.Rest).Success("删除成功！"),
		Data: types.Role{
			Id:          id,
			ParentId:    existRole.ParentId,
			Name:        existRole.Name,
			Description: existRole.Description.String, // 这里需要从数据库字段获取，确保正确处理NullString
			Level:       existRole.Level.String,       // 处理NullString
			Sort:        existRole.Sort,
			Status:      existRole.Status,
			Date: types.Date{
				CreatedAt:   createdAt,
				CreatedTime: createdTime,
				UpdatedAt:   updatedAt,
				UpdatedTime: updatedTime,
			},
		},
	}
	return
}
