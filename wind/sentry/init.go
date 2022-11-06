package sentry

import (
	"github.com/getsentry/sentry-go"
	"mond/wind/config"
	"sync"
)

var (
	sentryOnce sync.Once
)

func Init() {
	sentryOnce.Do(func() {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         config.GetSentry(),
			ServerName:  config.GetAppid(),
			Environment: config.GetEnv(),
		})
		if err != nil {
			panic(err)
		}
	})
}
