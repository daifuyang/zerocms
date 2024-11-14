package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"zerocms/api/admin/internal/config"
	"zerocms/api/admin/internal/handler"
	"zerocms/api/admin/internal/svc"
	"zerocms/utils"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	debug := c.Debug

	if !utils.FileExists(c.LockFile) || debug {
		err := ctx.ApiModel.SyncAllApi(server.Routes())
		if err != nil {
			panic(err)
		}
		if !debug {
			utils.CreateFile(c.LockFile)
		}
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
