package role

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"time"
	"zerocms/api/admin/internal/svc"
	"zerocms/api/admin/internal/types"
	"zerocms/api/model/role"
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

func (l *StoreLogic) Store(req *types.RoleReq) (resp *types.RoleResp, err error) {

	userId := l.ctx.Value("userId").(int64)

	userIdStr := strconv.FormatInt(userId, 10)

	// 获取当前用户角色
	//userRoles, err := l.svcCtx.UserRoleModel.List(l.ctx, userId)
	//if err != nil {
	//	return nil, err
	//}

	//super := user.SuperAdmin(userRoles)

	// 从请求中获取字段
	parentId := req.ParentId
	name := req.Name
	description := req.Description
	sort := req.Sort
	status := req.Status
	menuIds := req.MenuIds

	// 设置菜单权限
	// 判断当前用户角色，如果用户id为1或者角色id为1，则为超级管理员

	level := "0"
	// 如果上级parentId不等于0
	if parentId != 0 {
		parentRole, err := l.svcCtx.RoleModel.First(l.ctx, parentId)
		if err != nil {
			return nil, err
		}
		level = parentRole.Level.String
	}

	// 构造角色对象
	role := &role.SysRole{
		ParentId: parentId,
		Name:     name,
		Description: sql.NullString{
			String: description,
			Valid:  true,
		},
		Level: sql.NullString{
			String: level,
			Valid:  true,
		},
		Sort:   sort,
		Status: status,
	}

	var insertId int64

	err = l.svcCtx.Conn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 插入角色
		insert, err := l.svcCtx.RoleModel.WithSession(session).Insert(l.ctx, role)
		if err != nil {
			return err
		}

		insertId, err = insert.LastInsertId()
		if err != nil {
			return err
		}

		// 更新path
		role.Id = insertId
		role.Level = sql.NullString{
			String: level + "-" + strconv.FormatInt(insertId, 10),
			Valid:  true,
		}
		err = l.svcCtx.RoleModel.WithSession(session).Update(l.ctx, role)
		if err != nil {
			return err
		}

		// 插入角色菜单关联
		rules := make([][]string, 0)

		for _, menuId := range menuIds {
			menuIdStr := strconv.FormatInt(*menuId, 10)
			rules = append(rules, []string{userIdStr, menuIdStr})
		}

		_, err = l.svcCtx.Enforcer.AddPolicies(rules)
		if err != nil {
			return err
		}

		if err = l.svcCtx.Enforcer.SavePolicy(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	description = ""
	if role.Description.Valid {
		description = role.Description.String
	}

	level = ""
	if role.Level.Valid {
		level = role.Level.String
	}

	createdAt := time.Now().Unix()
	updatedAt := time.Now().Unix()

	createdTime := time.Unix(createdAt, 0).Format(time.DateTime)
	updatedTime := time.Unix(updatedAt, 0).Format(time.DateTime)

	resp = &types.RoleResp{
		Rest: new(types.Rest).Success("新增成功"),
		Data: types.Role{
			Id:          insertId,
			Name:        role.Name,
			Description: description,
			Level:       level,
			Sort:        role.Sort,
			Status:      role.Status,
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
