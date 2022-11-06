package reload

import (
	"context"
	"go.uber.org/zap"
	"mond/wind/config"
	"mond/wind/env"
	"mond/wind/logger"
	"meta/service/walle.admin/proto"
	"time"
)

func Reload(ctx context.Context) {
	if env.GetAppId() == "MetaWalle" {
		return
	}
	ctx, _ = context.WithTimeout(ctx, time.Second*10)
	adminClient, err := WalleAdmin.GetGrpcWalleAdminServiceClient()
	if err != nil {
		logger.GetLogger().Error(ctx, "reload err", zap.Any("err", err))
		return
	}
	go func() {
		if env.InDevelop() || env.InLocal() {
			time.Sleep(time.Second * 5)
		} else {
			time.Sleep(time.Second)
		}
		resp, err := adminClient.ReloadApp(ctx, &WalleAdmin.ReloadAppReq{AppId: env.GetAppId(), HostName: config.GetHostName()})
		if err != nil {
			logger.GetLogger().Error(ctx, "reload ReloadApp err", zap.Any("err", err))
		} else {
			logger.GetLogger().Debug(ctx, "reload ReloadApp success", zap.Any("resp", resp))
		}
		err = adminClient.Close()
		if err != nil {
			logger.GetLogger().Error(ctx, "reload close err", zap.Any("err", err))
			return
		}
	}()
}
