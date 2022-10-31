package app

import (
	"mond/service/biz.demo/domain/demo"
	mredis "mond/wind/cache/redis"
	"mond/wind/logger"
)

type App struct {
	rdb     *mredis.Client
	demoSvc *demo.Service
	_log    logger.Logger
}

func NewApp(rdb *mredis.Client, demoSvc *demo.Service) *App {
	return &App{
		rdb:     rdb,
		demoSvc: demoSvc,
		_log:    logger.GetLogger(),
	}
}
