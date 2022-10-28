package cmd

import "text/template"

var appTemplate, _ = template.New("").Parse(`
package app

import (
	mredis "meta/frame/cache/redis"
	"meta/frame/logger"
	"meta/service/{{.FolderPath}}/domain/demo"
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


`)
