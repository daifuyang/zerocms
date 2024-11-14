package svc

import (
	"database/sql"
	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"zerocms/api/admin/internal/config"
	"zerocms/api/admin/internal/middleware"
	"zerocms/api/admin/internal/migrate"
	"zerocms/api/model/department"
	"zerocms/api/model/menu"
	"zerocms/api/model/role"
	"zerocms/api/model/user"
	"zerocms/utils"
)

type ServiceContext struct {
	Config              config.Config
	Conn                sqlx.SqlConn
	UserModel           user.SysUserModel
	UserTokenModel      user.SysUserTokenModel
	RoleModel           role.SysRoleModel
	MenuModel           menu.SysMenuModel
	DepartmentModel     department.SysDepartmentModel
	RoleDepartmentModel role.SysRoleDepartmentModel
	RoleMenuModel       role.SysRoleMenuModel
	UserRoleModel       user.SysUserRoleModel
	ApiModel            menu.SysApiModel
	JwtMiddleware       rest.Middleware
	MenuRoleMiddleware  rest.Middleware
	Enforcer            *casbin.Enforcer
}

func NewServiceContext(c config.Config) *ServiceContext {
	dataSource := c.DataSource
	lockFile := c.LockFile
	debug := c.Debug
	if !utils.FileExists(lockFile) || debug {
		dbName, err := utils.ParseDatabaseName(dataSource)
		if err != nil {
			panic(err.Error())
		}
		dsn := utils.RemoveDatabaseName(dataSource)
		// 连接数据库
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		utils.CreateDatabase(db, dbName)
		_, err = db.Exec("USE " + dbName)
		if err != nil {
			panic(err)
		}

		err = utils.ExecuteSQLFilesInDirectory(db, "../model/sql/")
		if err != nil {
			panic(err)
		}
	}

	conn := sqlx.NewMysql(dataSource)
	db, err := conn.RawDB()
	if err != nil {
		panic(err)
	}

	// Initialize an adapter and use it in a Casbin enforcer:
	// The adapter will use the SQLite3 table name "casbin_rule_test",
	// the default table name is "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	a, err := sqladapter.NewAdapter(db, "mysql", "casbin_rule")
	if err != nil {
		panic(err)
	}

	e, err := casbin.NewEnforcer("public/install/rbac_model.conf", a)
	if err != nil {
		panic(err)
	}

	userTokenModel := user.NewSysUserTokenModel(conn, c.Cache)
	menuModel := menu.NewSysMenuModel(conn, c.Cache)
	roleMenuModel := role.NewSysRoleMenuModel(conn, c.Cache)

	ctx := &ServiceContext{
		Config:              c,
		Conn:                conn,
		UserModel:           user.NewSysUserModel(conn, c.Cache),
		MenuModel:           menuModel,
		UserTokenModel:      userTokenModel,
		RoleModel:           role.NewSysRoleModel(conn, c.Cache),
		DepartmentModel:     department.NewSysDepartmentModel(conn, c.Cache),
		RoleDepartmentModel: role.NewSysRoleDepartmentModel(conn, c.Cache),
		RoleMenuModel:       roleMenuModel,
		UserRoleModel:       user.NewSysUserRoleModel(conn, c.Cache),
		ApiModel:            menu.NewSysApiModel(conn, c.Cache),
		JwtMiddleware:       middleware.NewJwtMiddleware(userTokenModel).Handle,
		MenuRoleMiddleware:  middleware.NewMenuRoleMiddleware(menuModel, roleMenuModel, e).Handle,
		Enforcer:            e,
	}

	if !utils.FileExists(lockFile) || debug {
		migrate.Init(&migrate.ModelContext{
			UserModel: ctx.UserModel,
			MenuModel: ctx.MenuModel,
		})
	}
	return ctx
}
