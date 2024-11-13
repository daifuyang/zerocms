package svc

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"zerocms/api/admin/internal/config"
	"zerocms/api/admin/internal/middleware"
	"zerocms/api/admin/internal/migrate"
	"zerocms/api/model"
	"zerocms/utils"
)

type ServiceContext struct {
	Config         config.Config
	UserModel      model.SysUserModel
	UserTokenModel model.SysUserTokenModel
	RoleModel      model.SysRoleModel
	JwtMiddleware  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	dataSource := c.DataSource
	lockFile := "install.lock"
	debug := c.Debug
	if utils.FileExists(lockFile) == false || debug == true {
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
		utils.CreateFile(lockFile)
	}

	conn := sqlx.NewMysql(dataSource)

	userTokenModel := model.NewSysUserTokenModel(conn, c.Cache)

	ctx := &ServiceContext{
		Config:         c,
		UserModel:      model.NewSysUserModel(conn, c.Cache),
		UserTokenModel: userTokenModel,
		RoleModel:      model.NewSysRoleModel(conn, c.Cache),
		JwtMiddleware:  middleware.NewJwtMiddleware(userTokenModel).Handle,
	}

	migrate.Init(&migrate.ModelContext{
		UserModel: ctx.UserModel,
	})
	return ctx
}
