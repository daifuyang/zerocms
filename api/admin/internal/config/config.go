package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Debug          bool
	LockFile       string
	DataSource     string
	TokenSecretKey string
	Cache          cache.CacheConf
}
